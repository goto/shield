package servicedata

import (
	"context"

	"github.com/goto/shield/core/action"
	"github.com/goto/shield/core/namespace"
	"github.com/goto/shield/core/project"
	"github.com/goto/shield/core/relation"
	"github.com/goto/shield/core/resource"
	"github.com/goto/shield/core/user"
	"github.com/goto/shield/internal/schema"
)

const (
	keyNamespace = schema.ServiceDataKeyNamespace
	editActionID = schema.EditPermission
)

type ResourceService interface {
	Create(ctx context.Context, res resource.Resource) (resource.Resource, error)
	GetByURN(ctx context.Context, urn string) (resource.Resource, error)
}

type RelationService interface {
	Create(ctx context.Context, rel relation.RelationV2) (relation.RelationV2, error)
	CheckPermission(ctx context.Context, usr user.User, resourceNS namespace.Namespace, resourceIdxa string, action action.Action) (bool, error)
}

type ProjectService interface {
	Get(ctx context.Context, idOrSlug string) (project.Project, error)
}

type UserService interface {
	FetchCurrentUser(ctx context.Context) (user.User, error)
}

type Service struct {
	repository      Repository
	resourceService ResourceService
	relationService RelationService
	projectService  ProjectService
	userService     UserService
}

func NewService(repository Repository, resourceService ResourceService, relationService RelationService, projectService ProjectService, userService UserService) *Service {
	return &Service{
		repository:      repository,
		resourceService: resourceService,
		relationService: relationService,
		projectService:  projectService,
		userService:     userService,
	}
}

func (s Service) CreateKey(ctx context.Context, key Key) (Key, error) {
	// check if key contains ':'
	if key.Key == "" {
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
	key.URN = key.CreateURN()

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

	return createdServiceDataKey, nil
}

func (s Service) Upsert(ctx context.Context, servicedata ServiceData) (ServiceData, error) {
	if servicedata.Key.Key == "" {
		return ServiceData{}, ErrInvalidDetail
	}

	currentUser, err := s.userService.FetchCurrentUser(ctx)
	if err != nil {
		return ServiceData{}, err
	}

	prj, err := s.projectService.Get(ctx, servicedata.Key.ProjectID)
	if err != nil {
		return ServiceData{}, err
	}
	servicedata.Key.ProjectID = prj.ID
	servicedata.Key.ProjectSlug = prj.Slug

	servicedata.Key.URN = servicedata.Key.CreateURN()

	res, err := s.resourceService.GetByURN(ctx, servicedata.Key.URN)
	if err != nil {
		return ServiceData{}, err
	}

	permission, err := s.relationService.CheckPermission(ctx, currentUser, namespace.Namespace{ID: schema.ServiceDataKeyNamespace},
		res.Idxa, action.Action{ID: editActionID})
	if err != nil {
		return ServiceData{}, err
	}
	if !permission {
		return ServiceData{}, user.ErrInvalidEmail
	}

	returnedServiceData, err := s.repository.Upsert(ctx, servicedata)
	if err != nil {
		return ServiceData{}, err
	}

	return returnedServiceData, nil
}
