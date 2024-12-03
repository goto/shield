package resource

import (
	"context"
	"errors"

	"github.com/raystack/shield/internal/permission"
	"github.com/raystack/shield/model"
	"github.com/raystack/shield/utils"
)

type Service struct {
	Store       Store
	Permissions permission.Permissions
}

var (
	ResourceDoesntExist = errors.New("resource doesn't exist")
	InvalidUUID         = errors.New("invalid syntax of uuid")
)

type Store interface {
	GetResource(ctx context.Context, id string) (model.Resource, error)
	CreateResource(ctx context.Context, resource model.Resource) (model.Resource, error)
	ListResources(ctx context.Context, filters model.ResourceFilters) ([]model.Resource, error)
	UpdateResource(ctx context.Context, id string, resource model.Resource) (model.Resource, error)
}

func (s Service) Get(ctx context.Context, id string) (model.Resource, error) {
	return s.Store.GetResource(ctx, id)
}

func (s Service) Create(ctx context.Context, resource model.Resource) (model.Resource, error) {
	urn := utils.CreateResourceURN(resource)

	user, err := s.Permissions.FetchCurrentUser(ctx)

	if err != nil {
		return model.Resource{}, err
	}

	userId := resource.UserId

	if userId == "" {
		userId = user.Id
	}

	newResource, err := s.Store.CreateResource(ctx, model.Resource{
		Urn:            urn,
		Name:           resource.Name,
		OrganizationId: resource.OrganizationId,
		ProjectId:      resource.ProjectId,
		GroupId:        resource.GroupId,
		NamespaceId:    resource.NamespaceId,
		UserId:         userId,
	})

	if err != nil {
		return model.Resource{}, err
	}

	if err = s.Permissions.DeleteSubjectRelations(ctx, newResource); err != nil {
		return model.Resource{}, err
	}

	if newResource.GroupId != "" {
		err = s.Permissions.AddTeamToResource(ctx, model.Group{Id: resource.GroupId}, newResource)
		if err != nil {
			return model.Resource{}, err
		}
	}

	if userId != "" {
		err = s.Permissions.AddOwnerToResource(ctx, model.User{Id: userId}, newResource)
		if err != nil {
			return model.Resource{}, err
		}
	}

	err = s.Permissions.AddProjectToResource(ctx, model.Project{Id: resource.ProjectId}, newResource)

	if err != nil {
		return model.Resource{}, err
	}

	err = s.Permissions.AddOrgToResource(ctx, model.Organization{Id: resource.OrganizationId}, newResource)

	if err != nil {
		return model.Resource{}, err
	}

	return newResource, nil
}

func (s Service) List(ctx context.Context, filters model.ResourceFilters) ([]model.Resource, error) {
	return s.Store.ListResources(ctx, filters)
}

func (s Service) Update(ctx context.Context, id string, resource model.Resource) (model.Resource, error) {
	return s.Store.UpdateResource(ctx, id, resource)
}
