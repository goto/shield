package project

import (
	"context"
	"fmt"

	"github.com/goto/salt/log"
	"github.com/goto/shield/core/action"
	"github.com/goto/shield/core/activity"
	"github.com/goto/shield/core/namespace"
	"github.com/goto/shield/core/organization"
	"github.com/goto/shield/core/relation"
	"github.com/goto/shield/core/user"
	"github.com/goto/shield/internal/schema"
	pkgctx "github.com/goto/shield/pkg/context"
	"github.com/goto/shield/pkg/uuid"
)

const (
	auditKeyProjectCreate = "project.create"
	auditKeyProjectUpdate = "project.update"
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
	Log(ctx context.Context, action string, actor activity.Actor, data any) error
}

type Service struct {
	logger          log.Logger
	repository      Repository
	relationService RelationService
	userService     UserService
	activityService ActivityService
}

func NewService(logger log.Logger, repository Repository, relationService RelationService, userService UserService, activityService ActivityService) *Service {
	return &Service{
		logger:          logger,
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
		return Project{}, err
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

	go func() {
		ctx := pkgctx.WithoutCancel(ctx)
		projectLogData := newProject.ToProjectLogData()
		actor := activity.Actor{ID: currentUser.ID, Email: currentUser.Email}
		if err := s.activityService.Log(ctx, auditKeyProjectCreate, actor, projectLogData); err != nil {
			s.logger.Error(fmt.Sprintf("%s: %s", ErrLogActivity.Error(), err.Error()))
		}
	}()

	return newProject, nil
}

func (s Service) List(ctx context.Context) ([]Project, error) {
	return s.repository.List(ctx)
}

func (s Service) Update(ctx context.Context, prj Project) (Project, error) {
	currentUser, err := s.userService.FetchCurrentUser(ctx)
	if err != nil {
		return Project{}, err
	}

	if prj.ID != "" {
		return s.repository.UpdateByID(ctx, prj)
	}

	updatedProject, err := s.repository.UpdateBySlug(ctx, prj)
	if err != nil {
		return Project{}, err
	}

	go func() {
		ctx := pkgctx.WithoutCancel(ctx)
		projectLogData := updatedProject.ToProjectLogData()
		actor := activity.Actor{ID: currentUser.ID, Email: currentUser.Email}
		if err := s.activityService.Log(ctx, auditKeyProjectUpdate, actor, projectLogData); err != nil {
			s.logger.Error(fmt.Sprintf("%s: %s", ErrLogActivity.Error(), err.Error()))
		}
	}()

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
