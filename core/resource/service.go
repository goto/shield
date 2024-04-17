package resource

import (
	"context"
	"fmt"
	"strings"

	"github.com/goto/shield/core/action"
	"github.com/goto/shield/core/group"
	"github.com/goto/shield/core/namespace"
	"github.com/goto/shield/core/organization"
	"github.com/goto/shield/core/project"
	"github.com/goto/shield/core/relation"
	"github.com/goto/shield/core/user"
	"github.com/goto/shield/internal/schema"
	pkgctx "github.com/goto/shield/pkg/context"
	"github.com/goto/shield/pkg/uuid"
	"go.uber.org/zap"
)

const (
	auditKeyResourceCreate = "resource.create"
	auditKeyResourceUpdate = "resource.update"
)

type RelationService interface {
	Create(ctx context.Context, rel relation.RelationV2) (relation.RelationV2, error)
	Delete(ctx context.Context, rel relation.Relation) error
	CheckPermission(ctx context.Context, usr user.User, resourceNS namespace.Namespace, resourceIdxa string, action action.Action) (bool, error)
	DeleteSubjectRelations(ctx context.Context, resourceType, optionalResourceID string) error
}

type UserService interface {
	FetchCurrentUser(ctx context.Context) (user.User, error)
}

type ProjectService interface {
	Get(ctx context.Context, id string) (project.Project, error)
}

type OrganizationService interface {
	Get(ctx context.Context, id string) (organization.Organization, error)
}

type GroupService interface {
	Get(ctx context.Context, id string) (group.Group, error)
}

type ActivityService interface {
	Log(ctx context.Context, action string, actor string, data any) error
}

type Service struct {
	logger              *zap.SugaredLogger
	repository          Repository
	configRepository    ConfigRepository
	relationService     RelationService
	userService         UserService
	projectService      ProjectService
	organizationService OrganizationService
	groupService        GroupService
	activityService     ActivityService
}

func NewService(logger *zap.SugaredLogger, repository Repository, configRepository ConfigRepository, relationService RelationService, userService UserService, projectService ProjectService, organizationService OrganizationService, groupService GroupService, activityService ActivityService) *Service {
	return &Service{
		logger:              logger,
		repository:          repository,
		configRepository:    configRepository,
		relationService:     relationService,
		userService:         userService,
		projectService:      projectService,
		organizationService: organizationService,
		groupService:        groupService,
		activityService:     activityService,
	}
}

func (s Service) Get(ctx context.Context, id string) (Resource, error) {
	return s.repository.GetByID(ctx, id)
}

func (s Service) Create(ctx context.Context, res Resource) (Resource, error) {
	currentUser, err := s.userService.FetchCurrentUser(ctx)
	if err != nil {
		return Resource{}, fmt.Errorf("%w: %s", user.ErrInvalidEmail, err.Error())
	}

	urn := res.CreateURN()

	if err != nil {
		return Resource{}, err
	}

	fetchedProject, err := s.projectService.Get(ctx, res.ProjectID)
	if err != nil {
		return Resource{}, err
	}

	userId := res.UserID
	if strings.TrimSpace(userId) == "" {
		userId = currentUser.ID
	}

	newResource, err := s.repository.Create(ctx, Resource{
		URN:            urn,
		Name:           res.Name,
		OrganizationID: fetchedProject.Organization.ID,
		ProjectID:      fetchedProject.ID,
		NamespaceID:    res.NamespaceID,
		UserID:         userId,
	})
	if err != nil {
		return Resource{}, err
	}

	if err = s.relationService.DeleteSubjectRelations(ctx, newResource.NamespaceID, newResource.Idxa); err != nil {
		return Resource{}, err
	}

	if err = s.AddProjectToResource(ctx, project.Project{ID: res.ProjectID}, newResource); err != nil {
		return Resource{}, err
	}

	if err = s.AddOrgToResource(ctx, organization.Organization{ID: newResource.OrganizationID}, newResource); err != nil {
		return Resource{}, err
	}

	go func() {
		ctx := pkgctx.WithoutCancel(ctx)
		resourceLogData := newResource.ToResourceLogData()
		if err := s.activityService.Log(ctx, auditKeyResourceCreate, currentUser.ID, resourceLogData); err != nil {
			s.logger.Errorf("%s: %s", ErrLogActivity.Error(), err.Error())
		}
	}()

	return newResource, nil
}

func (s Service) List(ctx context.Context, flt Filter) (PagedResources, error) {
	resources, err := s.repository.List(ctx, flt)
	if err != nil {
		return PagedResources{}, err
	}
	return PagedResources{
		Count:     int32(len(resources)),
		Resources: resources,
	}, nil
}

func (s Service) Update(ctx context.Context, id string, resource Resource) (Resource, error) {
	currentUser, err := s.userService.FetchCurrentUser(ctx)
	if err != nil {
		return Resource{}, fmt.Errorf("%w: %s", user.ErrInvalidEmail, err.Error())
	}

	// TODO there should be an update logic like create here
	updatedResource, err := s.repository.Update(ctx, id, resource)
	if err != nil {
		return Resource{}, err
	}

	go func() {
		ctx := pkgctx.WithoutCancel(ctx)
		resourceLogData := updatedResource.ToResourceLogData()
		if err := s.activityService.Log(ctx, auditKeyResourceUpdate, currentUser.ID, resourceLogData); err != nil {
			s.logger.Errorf("%s: %s", ErrLogActivity.Error(), err.Error())
		}
	}()

	return updatedResource, nil
}

func (s Service) AddProjectToResource(ctx context.Context, project project.Project, res Resource) error {
	rel := relation.RelationV2{
		Object: relation.Object{
			ID:          res.Idxa,
			NamespaceID: res.NamespaceID,
		},
		Subject: relation.Subject{
			RoleID:    schema.ProjectRelationName,
			ID:        project.ID,
			Namespace: schema.ProjectNamespace,
		},
	}

	if _, err := s.relationService.Create(ctx, rel); err != nil {
		return err
	}

	return nil
}

func (s Service) AddOrgToResource(ctx context.Context, org organization.Organization, res Resource) error {
	rel := relation.RelationV2{
		Object: relation.Object{
			ID:          res.Idxa,
			NamespaceID: res.NamespaceID,
		},
		Subject: relation.Subject{
			RoleID:    schema.OrganizationRelationName,
			ID:        org.ID,
			Namespace: schema.OrganizationNamespace,
		},
	}

	if _, err := s.relationService.Create(ctx, rel); err != nil {
		return err
	}
	return nil
}

func (s Service) GetAllConfigs(ctx context.Context) ([]YAML, error) {
	return s.configRepository.GetAll(ctx)
}

// TODO(krkvrm): Separate Authz for Resources & System Namespaces
func (s Service) CheckAuthz(ctx context.Context, res Resource, act action.Action) (bool, error) {
	currentUser, err := s.userService.FetchCurrentUser(ctx)
	if err != nil {
		return false, err
	}

	isSystemNS := namespace.IsSystemNamespaceID(res.NamespaceID)
	fetchedResource := res

	if isSystemNS {
		if !uuid.IsValid(res.Name) {
			switch res.NamespaceID {
			case namespace.DefinitionProject.ID:
				project, err := s.projectService.Get(ctx, res.Name)
				if err != nil {
					return false, err
				}
				res.Name = project.ID
			case namespace.DefinitionOrg.ID:
				organization, err := s.organizationService.Get(ctx, res.Name)
				if err != nil {
					return false, err
				}
				res.Name = organization.ID
			case namespace.DefinitionTeam.ID:
				group, err := s.groupService.Get(ctx, res.Name)
				if err != nil {
					return false, err
				}
				res.Name = group.ID
			}
		}
		fetchedResource.Idxa = res.Name
	} else {
		fetchedResource, err = s.repository.GetByNamespace(ctx, res.Name, res.NamespaceID)
		if err != nil {
			return false, ErrNotExist
		}
	}

	fetchedResourceNS := namespace.Namespace{ID: fetchedResource.NamespaceID}
	return s.relationService.CheckPermission(ctx, currentUser, fetchedResourceNS, fetchedResource.Idxa, act)
}
