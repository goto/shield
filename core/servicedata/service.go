package servicedata

import (
	"context"
	"fmt"
	"slices"

	"github.com/goto/salt/log"
	"github.com/goto/shield/core/action"
	"github.com/goto/shield/core/activity"
	"github.com/goto/shield/core/namespace"
	"github.com/goto/shield/core/project"
	"github.com/goto/shield/core/relation"
	"github.com/goto/shield/core/resource"
	"github.com/goto/shield/core/user"
	"github.com/goto/shield/internal/schema"
	"github.com/goto/shield/pkg/errors"
)

const (
	keyNamespace         = schema.ServiceDataKeyNamespace
	userNamespace        = schema.UserPrincipal
	groupNamespace       = schema.GroupPrincipal
	viewActionID         = schema.ViewPermission
	editActionID         = schema.EditPermission
	membershipPermission = schema.MembershipPermission

	auditKeyServiceDataKeyCreate = "service_data_key.create"
)

type ResourceService interface {
	Create(ctx context.Context, res resource.Resource) (resource.Resource, error)
	GetByURN(ctx context.Context, urn string) (resource.Resource, error)
}

type RelationService interface {
	Create(ctx context.Context, rel relation.RelationV2) (relation.RelationV2, error)
	CheckPermission(ctx context.Context, usr user.User, resourceNS namespace.Namespace, resourceIdxa string, action action.Action) (bool, error)
	LookupResources(ctx context.Context, resourceType, permission, subjectType, subjectID string) ([]string, error)
}

type ProjectService interface {
	Get(ctx context.Context, idOrSlug string) (project.Project, error)
}

type UserService interface {
	FetchCurrentUser(ctx context.Context) (user.User, error)
}

type ActivityService interface {
	Log(ctx context.Context, action string, actor activity.Actor, data any) error
}

type Service struct {
	logger          log.Logger
	repository      Repository
	resourceService ResourceService
	relationService RelationService
	projectService  ProjectService
	userService     UserService
	activityService ActivityService
}

func NewService(logger log.Logger, repository Repository, resourceService ResourceService, relationService RelationService, projectService ProjectService, userService UserService, activityService ActivityService) *Service {
	return &Service{
		logger:          logger,
		repository:      repository,
		resourceService: resourceService,
		relationService: relationService,
		projectService:  projectService,
		userService:     userService,
		activityService: activityService,
	}
}

func (s Service) CreateKey(ctx context.Context, key Key) (Key, error) {
	// check if key contains ':'
	if key.Name == "" {
		return Key{}, ErrInvalidDetail
	}

	// fetch current user
	currentUser, err := s.userService.FetchCurrentUser(ctx)
	if err != nil {
		return Key{}, err
	}

	// Get Project
	prj, err := s.projectService.Get(ctx, key.ProjectID)
	if err != nil {
		return Key{}, err
	}
	key.ProjectID = prj.ID
	key.ProjectSlug = prj.Slug

	// create URN
	key.URN = CreateURN(key.ProjectSlug, key.Name)

	// Transaction for postgres repository
	// TODO find way to use transaction for spicedb
	ctx = s.repository.WithTransaction(ctx)

	// insert the service data key
	res, err := s.resourceService.Create(ctx, resource.Resource{
		Name:        key.URN,
		NamespaceID: keyNamespace,
		ProjectID:   key.ProjectID,
		UserID:      currentUser.ID,
	})
	if err != nil {
		if err := s.repository.Rollback(ctx, err); err != nil {
			return Key{}, err
		}
		return Key{}, err
	}
	key.ResourceID = res.Idxa

	// insert service data key to the servicedata_keys table
	createdServiceDataKey, err := s.repository.CreateKey(ctx, key)
	if err != nil {
		if err := s.repository.Rollback(ctx, err); err != nil {
			return Key{}, err
		}
		return Key{}, err
	}

	// create relation
	_, err = s.relationService.Create(ctx, relation.RelationV2{
		Object: relation.Object{
			ID:          res.Idxa,
			NamespaceID: schema.ServiceDataKeyNamespace,
		},
		Subject: relation.Subject{
			ID:        currentUser.ID,
			RoleID:    schema.OwnerRole,
			Namespace: schema.UserPrincipal,
		},
	})
	if err != nil {
		if err := s.repository.Rollback(ctx, err); err != nil {
			return Key{}, err
		}
		return Key{}, err
	}

	if err := s.repository.Commit(ctx); err != nil {
		return Key{}, err
	}

	go func() {
		ctx = context.WithoutCancel(ctx)
		actor := activity.Actor{ID: currentUser.ID, Email: currentUser.Email}
		if err := s.activityService.Log(ctx, auditKeyServiceDataKeyCreate, actor, key.ToKeyLogData()); err != nil {
			s.logger.Error(fmt.Sprintf("%s: %s", ErrLogActivity.Error(), err.Error()))
		}
	}()

	return createdServiceDataKey, nil
}

func (s Service) GetKeyByURN(ctx context.Context, urn string) (Key, error) {
	return s.repository.GetKeyByURN(ctx, urn)
}

func (s Service) Upsert(ctx context.Context, sd ServiceData) (ServiceData, error) {
	if sd.Key.Name == "" {
		return ServiceData{}, ErrInvalidDetail
	}

	currentUser, err := s.userService.FetchCurrentUser(ctx)
	if err != nil {
		return ServiceData{}, err
	}

	prj, err := s.projectService.Get(ctx, sd.Key.ProjectID)
	if err != nil {
		return ServiceData{}, err
	}
	sd.Key.ProjectSlug = prj.Slug

	sd.Key.URN = CreateURN(sd.Key.ProjectSlug, sd.Key.Name)

	sd.Key, err = s.repository.GetKeyByURN(ctx, sd.Key.URN)
	if err != nil {
		return ServiceData{}, err
	}

	permission, err := s.relationService.CheckPermission(ctx, currentUser, namespace.Namespace{ID: schema.ServiceDataKeyNamespace},
		sd.Key.ResourceID, action.Action{ID: editActionID})
	if err != nil {
		return ServiceData{}, err
	}
	if !permission {
		return ServiceData{}, errors.ErrForbidden
	}

	returnedServiceData, err := s.repository.Upsert(ctx, sd)
	if err != nil {
		return ServiceData{}, err
	}

	return returnedServiceData, nil
}

func (s Service) Get(ctx context.Context, filter Filter) ([]ServiceData, error) {
	// fetch current user
	currentUser, err := s.userService.FetchCurrentUser(ctx)
	if err != nil {
		return []ServiceData{}, err
	}

	// validate project and get project id from project slug
	if filter.Project != "" {
		prj, err := s.projectService.Get(ctx, filter.Project)
		if err != nil {
			return []ServiceData{}, err
		}
		filter.Project = prj.ID
	}

	// build entity ID filter
	filter.EntityIDs = [][]string{}
	if filter.Namespace == groupNamespace {
		filter.EntityIDs = append(filter.EntityIDs, []string{groupNamespace, filter.ID})
	}
	if filter.Namespace == userNamespace && slices.Contains(filter.Entities, userNamespace) {
		filter.EntityIDs = append(filter.EntityIDs, []string{userNamespace, filter.ID})
	}
	if filter.Namespace == userNamespace && slices.Contains(filter.Entities, groupNamespace) {
		entityGroup, err := s.relationService.LookupResources(ctx, groupNamespace, membershipPermission, userNamespace, filter.ID)
		if err != nil {
			return []ServiceData{}, err
		}
		for _, ent := range entityGroup {
			filter.EntityIDs = append(filter.EntityIDs, []string{groupNamespace, ent})
		}
	}

	if len(filter.EntityIDs) == 0 {
		return []ServiceData{}, nil
	}

	// get all service data key resources that visible by current user
	viewSD, err := s.relationService.LookupResources(ctx, keyNamespace, viewActionID, userNamespace, currentUser.ID)
	if err != nil {
		return []ServiceData{}, err
	}

	serviceData, err := s.repository.Get(ctx, filter)
	if err != nil {
		return []ServiceData{}, err
	}

	resultSD := []ServiceData{}
	for _, sd := range serviceData {
		if slices.Contains(viewSD, sd.Key.ResourceID) {
			resultSD = append(resultSD, sd)
		}
	}

	return resultSD, nil
}
