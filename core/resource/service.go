package resource

import (
	"context"
	"fmt"
	"strings"

	"github.com/goto/salt/log"
	"github.com/goto/shield/core/action"
	"github.com/goto/shield/core/activity"
	"github.com/goto/shield/core/group"
	"github.com/goto/shield/core/namespace"
	"github.com/goto/shield/core/organization"
	"github.com/goto/shield/core/policy"
	"github.com/goto/shield/core/project"
	"github.com/goto/shield/core/relation"
	"github.com/goto/shield/core/user"
	"github.com/goto/shield/internal/schema"
	"github.com/goto/shield/pkg/db"
	"github.com/goto/shield/pkg/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	auditKeyResourceCreate = "resource.create"
	auditKeyResourceUpdate = "resource.update"

	userNamespace = schema.UserPrincipal
)

type RelationService interface {
	Create(ctx context.Context, rel relation.RelationV2) (relation.RelationV2, error)
	Delete(ctx context.Context, rel relation.Relation) error
	CheckPermission(ctx context.Context, usr user.User, resourceNS namespace.Namespace, resourceIdxa string, action action.Action) (bool, error)
	BulkCheckPermission(ctx context.Context, rels []relation.Relation, acts []action.Action) ([]relation.Permission, error)
	DeleteSubjectRelations(ctx context.Context, resourceType, optionalResourceID string) error
	LookupResources(ctx context.Context, resourceType, permission, subjectType, subjectID string) ([]string, error)
}

type UserService interface {
	FetchCurrentUser(ctx context.Context) (user.User, error)
	Get(ctx context.Context, userID string) (user.User, error)
}

type ProjectService interface {
	Get(ctx context.Context, id string) (project.Project, error)
}

type OrganizationService interface {
	Get(ctx context.Context, id string) (organization.Organization, error)
}

type GroupService interface {
	GetBySlug(ctx context.Context, id string) (group.Group, error)
}

type ActivityService interface {
	Log(ctx context.Context, action string, actor activity.Actor, data any) error
}

type PolicyService interface {
	List(ctx context.Context, filter policy.Filters) ([]policy.Policy, error)
}

type NamespaceService interface {
	List(ctx context.Context) ([]namespace.Namespace, error)
}

type Service struct {
	logger              log.Logger
	repository          Repository
	configRepository    ConfigRepository
	relationService     RelationService
	userService         UserService
	projectService      ProjectService
	organizationService OrganizationService
	groupService        GroupService
	policyService       PolicyService
	namespaceService    NamespaceService
	activityService     ActivityService
}

func NewService(logger log.Logger, repository Repository, configRepository ConfigRepository, relationService RelationService, userService UserService, projectService ProjectService, organizationService OrganizationService, groupService GroupService, policyService PolicyService, namespaceService NamespaceService, activityService ActivityService) *Service {
	return &Service{
		logger:              logger,
		repository:          repository,
		configRepository:    configRepository,
		relationService:     relationService,
		userService:         userService,
		projectService:      projectService,
		organizationService: organizationService,
		groupService:        groupService,
		policyService:       policyService,
		namespaceService:    namespaceService,
		activityService:     activityService,
	}
}

func (s Service) GetByURN(ctx context.Context, id string) (Resource, error) {
	return s.repository.GetByURN(ctx, id)
}

func (s Service) Get(ctx context.Context, id string) (Resource, error) {
	return s.repository.GetByID(ctx, id)
}

func (s Service) Upsert(ctx context.Context, res Resource) (Resource, error) {
	currentUser, err := s.userService.FetchCurrentUser(ctx)
	if err != nil {
		return Resource{}, err
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

	newResource, err := s.repository.Upsert(ctx, Resource{
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
		ctx := context.WithoutCancel(ctx)
		resourceLogData := newResource.ToLogData()
		actor := activity.Actor{ID: currentUser.ID, Email: currentUser.Email}
		if err := s.activityService.Log(ctx, auditKeyResourceCreate, actor, resourceLogData); err != nil {
			s.logger.Error(fmt.Sprintf("%s: %s", ErrLogActivity.Error(), err.Error()))
		}
	}()

	return newResource, nil
}

func (s Service) Create(ctx context.Context, res Resource) (Resource, error) {
	currentUser, err := s.userService.FetchCurrentUser(ctx)
	if err != nil {
		return Resource{}, err
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
		ctx = db.WithoutTx(ctx)
		ctx = context.WithoutCancel(ctx)
		resourceLogData := newResource.ToLogData()
		actor := activity.Actor{ID: currentUser.ID, Email: currentUser.Email}
		if err := s.activityService.Log(ctx, auditKeyResourceCreate, actor, resourceLogData); err != nil {
			s.logger.Error(fmt.Sprintf("%s: %s", ErrLogActivity.Error(), err.Error()))
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
		return Resource{}, err
	}

	// TODO there should be an update logic like create here
	updatedResource, err := s.repository.Update(ctx, id, resource)
	if err != nil {
		return Resource{}, err
	}

	go func() {
		ctx := context.WithoutCancel(ctx)
		resourceLogData := updatedResource.ToLogData()
		actor := activity.Actor{ID: currentUser.ID, Email: currentUser.Email}
		if err := s.activityService.Log(ctx, auditKeyResourceUpdate, actor, resourceLogData); err != nil {
			s.logger.Error(fmt.Sprintf("%s: %s", ErrLogActivity.Error(), err.Error()))
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
				group, err := s.groupService.GetBySlug(ctx, res.Name)
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
			fetchedResource, err = s.repository.GetByID(ctx, res.Name)
			if err != nil {
				return false, ErrNotExist
			}
		}
	}
	fetchedResourceNS := namespace.Namespace{ID: fetchedResource.NamespaceID}
	return s.relationService.CheckPermission(ctx, currentUser, fetchedResourceNS, fetchedResource.Idxa, act)
}

func (s Service) BulkCheckAuthz(ctx context.Context, resources []Resource, actions []action.Action) ([]relation.Permission, error) {
	currentUser, err := s.userService.FetchCurrentUser(ctx)
	if err != nil {
		return []relation.Permission{}, err
	}

	var relations []relation.Relation
	for _, res := range resources {
		isSystemNS := namespace.IsSystemNamespaceID(res.NamespaceID)
		fetchedResource := res

		if isSystemNS {
			if !uuid.IsValid(res.Name) {
				switch res.NamespaceID {
				case namespace.DefinitionProject.ID:
					project, err := s.projectService.Get(ctx, res.Name)
					if err != nil {
						return []relation.Permission{}, err
					}
					res.Name = project.ID
				case namespace.DefinitionOrg.ID:
					organization, err := s.organizationService.Get(ctx, res.Name)
					if err != nil {
						return []relation.Permission{}, err
					}
					res.Name = organization.ID
				case namespace.DefinitionTeam.ID:
					group, err := s.groupService.GetBySlug(ctx, res.Name)
					if err != nil {
						return []relation.Permission{}, err
					}
					res.Name = group.ID
				}
			}
			fetchedResource.Idxa = res.Name
		} else {
			fetchedResource, err = s.repository.GetByNamespace(ctx, res.Name, res.NamespaceID)
			if err != nil {
				fetchedResource, err = s.repository.GetByID(ctx, res.Name)
				if err != nil {
					return []relation.Permission{}, ErrNotExist
				}
			}
		}
		fetchedResourceNS := namespace.Namespace{ID: fetchedResource.NamespaceID}

		relations = append(relations, relation.Relation{
			SubjectID:        currentUser.ID,
			SubjectNamespace: namespace.DefinitionUser,
			ObjectID:         fetchedResource.Idxa,
			ObjectNamespace:  fetchedResourceNS,
		})
	}
	return s.relationService.BulkCheckPermission(ctx, relations, actions)
}

func (s Service) ListUserResourcesByType(ctx context.Context, userID string, resourceType string, permissions []string) (ResourcePermissions, error) {
	user, err := s.userService.Get(ctx, userID)
	if err != nil {
		return ResourcePermissions{}, err
	}

	res, err := s.listUserResources(ctx, resourceType, user, permissions)
	if err != nil {
		switch status.Code(err) {
		case codes.FailedPrecondition:
			s.logger.Warn(err.Error())
			return ResourcePermissions{}, ErrInvalidDetail
		default:
			return ResourcePermissions{}, err
		}
	}

	return res, nil
}

func (s Service) ListAllUserResources(ctx context.Context, userID string, resourceTypes []string, permissions []string) (map[string]ResourcePermissions, error) {
	user, err := s.userService.Get(ctx, userID)
	if err != nil {
		return map[string]ResourcePermissions{}, err
	}

	if len(resourceTypes) == 0 {
		namespaces, err := s.namespaceService.List(ctx)
		if err != nil {
			return map[string]ResourcePermissions{}, err
		}

		for _, ns := range namespaces {
			if namespace.IsSystemNamespaceID(ns.ID) {
				continue
			}
			resourceTypes = append(resourceTypes, ns.ID)
		}
	}

	result := make(map[string]ResourcePermissions)
	for _, res := range resourceTypes {
		if _, ok := result[res]; !ok {
			list, err := s.listUserResources(ctx, res, user, permissions)
			if err != nil {
				switch status.Code(err) {
				case codes.FailedPrecondition:
					s.logger.Warn(err.Error())
					continue
				default:
					return map[string]ResourcePermissions{}, err
				}
			}
			if len(list) != 0 {
				result[res] = list
			}
		}
	}

	return result, nil
}

func (s Service) listUserResources(ctx context.Context, resourceType string, user user.User, permissions []string) (ResourcePermissions, error) {
	if len(permissions) == 0 {
		policies, err := s.policyService.List(ctx, policy.Filters{NamespaceID: resourceType})
		if err != nil {
			return ResourcePermissions{}, err
		}

		for _, p := range policies {
			permissions = append(permissions, strings.Split(p.ActionID, ".")[0])
		}
	}

	resPermissionsMap := make(ResourcePermissions)
	actMap := make(map[string]bool)
	for _, p := range permissions {
		if _, ok := actMap[p]; ok {
			continue
		}
		actMap[p] = true
		resources, err := s.relationService.LookupResources(ctx, resourceType, p, userNamespace, user.ID)
		if err != nil {
			// continue if permission under a namespace is not found
			// https://github.com/authzed/spicedb/blob/main/internal/dispatch/graph/errors.go#L73
			if strings.Contains(err.Error(), "not found under definition") {
				s.logger.Warn(err.Error())
				continue
			}
			return ResourcePermissions{}, err
		}

		for _, r := range resources {
			if _, ok := resPermissionsMap[r]; !ok {
				resPermissionsMap[r] = []string{p}
			} else {
				resPermissionsMap[r] = append(resPermissionsMap[r], p)
			}
		}
	}

	return resPermissionsMap, nil
}
