package project

import (
	"context"
	"fmt"
	"maps"

	"github.com/goto/shield/core/action"
	"github.com/goto/shield/core/namespace"
	"github.com/goto/shield/core/organization"
	"github.com/goto/shield/core/relation"
	"github.com/goto/shield/core/user"
	"github.com/goto/shield/internal/schema"
	"github.com/goto/shield/pkg/uuid"
	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
)

const (
	AuditKeyProjectCreate = "project.create"
	AuditKeyProjectUpdate = "project.update"

	AuditEntity = "project"
)

type RelationService interface {
	Create(ctx context.Context, rel relation.RelationV2) (relation.RelationV2, error)
	Delete(ctx context.Context, rel relation.Relation) error
	CheckPermission(ctx context.Context, usr user.User, resourceNS namespace.Namespace, resourceIdxa string, action action.Action) (bool, error)
}

type UserService interface {
	FetchCurrentUser(ctx context.Context) (user.User, error)
	GetByID(ctx context.Context, id string) (user.User, error)
	GetByIDs(ctx context.Context, userIDs []string) ([]user.User, error)
}

type ActivityService interface {
	Log(ctx context.Context, action string, actor string, data map[string]string) error
}

type Service struct {
	repository      Repository
	relationService RelationService
	userService     UserService
	activityService ActivityService
}

func NewService(repository Repository, relationService RelationService, userService UserService, activityService ActivityService) *Service {
	return &Service{
		repository:      repository,
		relationService: relationService,
		userService:     userService,
		activityService: activityService,
	}
}

func (s Service) Get(ctx context.Context, idOrSlug string) (Project, error) {
	if uuid.IsValid(idOrSlug) {
		return s.repository.GetByID(ctx, idOrSlug)
	}
	return s.repository.GetBySlug(ctx, idOrSlug)
}

func (s Service) Create(ctx context.Context, prj Project) (Project, error) {
	currentUser, err := s.userService.FetchCurrentUser(ctx)
	if err != nil {
		return Project{}, fmt.Errorf("%w: %s", user.ErrInvalidEmail, err.Error())
	}

	newProject, err := s.repository.Create(ctx, Project{
		Name:         prj.Name,
		Slug:         prj.Slug,
		Metadata:     prj.Metadata,
		Organization: prj.Organization,
	})
	if err != nil {
		return Project{}, err
	}

	if err = s.addProjectToOrg(ctx, newProject, prj.Organization); err != nil {
		return Project{}, err
	}

	logData := map[string]string{
		"entity": AuditEntity,
		"id":     newProject.ID,
		"name":   newProject.Name,
		"slug":   newProject.Slug,
		"orgId":  newProject.Organization.ID,
	}
	maps.Copy(logData, newProject.Metadata.ToStringValueMap())
	if err := s.activityService.Log(ctx, AuditKeyProjectCreate, currentUser.ID, logData); err != nil {
		logger := grpczap.Extract(ctx)
		logger.Error(ErrLogActivity.Error())
	}

	return newProject, nil
}

func (s Service) List(ctx context.Context) ([]Project, error) {
	return s.repository.List(ctx)
}

func (s Service) Update(ctx context.Context, prj Project) (Project, error) {
	if prj.ID != "" {
		return s.repository.UpdateByID(ctx, prj)
	}

	updatedProject, err := s.repository.UpdateBySlug(ctx, prj)
	if err != nil {
		return Project{}, err
	}

	logger := grpczap.Extract(ctx)
	currentUser, _ := s.userService.FetchCurrentUser(ctx)
	logData := map[string]string{
		"entity": AuditEntity,
		"id":     updatedProject.ID,
		"name":   updatedProject.Name,
		"slug":   updatedProject.Slug,
		"orgId":  updatedProject.Organization.ID,
	}
	maps.Copy(logData, updatedProject.Metadata.ToStringValueMap())
	if err := s.activityService.Log(ctx, AuditKeyProjectUpdate, currentUser.ID, logData); err != nil {
		logger.Error(ErrLogActivity.Error())
	}

	return updatedProject, err
}

func (s Service) AddAdmins(ctx context.Context, idOrSlug string, userIds []string) ([]user.User, error) {
	// TODO(discussion): can be done with create relations
	return []user.User{}, nil
}

func (s Service) ListAdmins(ctx context.Context, id string) ([]user.User, error) {
	return s.repository.ListAdmins(ctx, id)
}

func (s Service) RemoveAdmin(ctx context.Context, idOrSlug string, userId string) ([]user.User, error) {
	// TODO(discussion): can be done with delete relations
	return []user.User{}, nil
}

func (s Service) addProjectToOrg(ctx context.Context, prj Project, org organization.Organization) error {
	rel := relation.RelationV2{
		Object: relation.Object{
			ID:          prj.ID,
			NamespaceID: schema.ProjectNamespace,
		},
		Subject: relation.Subject{
			ID:        org.ID,
			Namespace: schema.OrganizationNamespace,
			RoleID:    schema.OrganizationRelationName,
		},
	}

	if _, err := s.relationService.Create(ctx, rel); err != nil {
		return err
	}

	return nil
}
