package organization

import (
	"context"
	"fmt"

	"github.com/goto/salt/log"
	"github.com/goto/shield/core/action"
	"github.com/goto/shield/core/namespace"
	"github.com/goto/shield/core/relation"
	"github.com/goto/shield/core/user"
	"github.com/goto/shield/internal/schema"
	pkgctx "github.com/goto/shield/pkg/context"
	"github.com/goto/shield/pkg/uuid"
)

const (
	auditKeyOrganizationCreate = "organization.create"
	auditKeyOrganizationUpdate = "organization.update"
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
	Log(ctx context.Context, action string, actor string, data any) error
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

func (s Service) Get(ctx context.Context, idOrSlug string) (Organization, error) {
	if uuid.IsValid(idOrSlug) {
		return s.repository.GetByID(ctx, idOrSlug)
	}
	return s.repository.GetBySlug(ctx, idOrSlug)
}

func (s Service) Create(ctx context.Context, org Organization) (Organization, error) {
	currentUser, err := s.userService.FetchCurrentUser(ctx)
	if err != nil {
		return Organization{}, fmt.Errorf("%w: %s", user.ErrInvalidEmail, err.Error())
	}

	newOrg, err := s.repository.Create(ctx, Organization{
		Name:     org.Name,
		Slug:     org.Slug,
		Metadata: org.Metadata,
	})
	if err != nil {
		return Organization{}, err
	}

	if err = s.addAdminToOrg(ctx, currentUser, newOrg); err != nil {
		return Organization{}, err
	}

	go func() {
		ctx := pkgctx.WithoutCancel(ctx)
		organizationLogData := newOrg.ToOrganizationLogData()
		if err := s.activityService.Log(ctx, auditKeyOrganizationCreate, currentUser.ID, organizationLogData); err != nil {
			s.logger.Error(fmt.Sprintf("%s: %s", ErrLogActivity.Error(), err.Error()))
		}
	}()

	return newOrg, nil
}

func (s Service) List(ctx context.Context) ([]Organization, error) {
	return s.repository.List(ctx)
}

func (s Service) Update(ctx context.Context, org Organization) (Organization, error) {
	currentUser, err := s.userService.FetchCurrentUser(ctx)
	if err != nil {
		s.logger.Error(fmt.Sprintf("%s: %s", user.ErrInvalidEmail.Error(), err.Error()))
	}

	var updatedOrg Organization

	if org.ID != "" {
		updatedOrg, err = s.repository.UpdateByID(ctx, org)
	} else {
		updatedOrg, err = s.repository.UpdateBySlug(ctx, org)
	}

	if err != nil {
		return Organization{}, err
	}

	go func() {
		ctx := pkgctx.WithoutCancel(ctx)
		organizationLogData := updatedOrg.ToOrganizationLogData()
		if err := s.activityService.Log(ctx, auditKeyOrganizationUpdate, currentUser.ID, organizationLogData); err != nil {
			s.logger.Error(fmt.Sprintf("%s: %s", ErrLogActivity.Error(), err.Error()))
		}
	}()

	return updatedOrg, nil
}

func (s Service) ListAdmins(ctx context.Context, idOrSlug string) ([]user.User, error) {
	var org Organization
	var err error
	if uuid.IsValid(idOrSlug) {
		return s.repository.ListAdminsByOrgID(ctx, idOrSlug)
	}
	org, err = s.repository.GetBySlug(ctx, idOrSlug)
	if err != nil {
		return []user.User{}, err
	}
	return s.repository.ListAdminsByOrgID(ctx, org.ID)
}

func (s Service) addAdminToOrg(ctx context.Context, user user.User, org Organization) error {
	rel := relation.RelationV2{
		Object: relation.Object{
			ID:          org.ID,
			NamespaceID: schema.OrganizationNamespace,
		},
		Subject: relation.Subject{
			ID:        user.ID,
			Namespace: schema.UserPrincipal,
			RoleID:    schema.OwnerRole,
		},
	}

	if _, err := s.relationService.Create(ctx, rel); err != nil {
		return err
	}
	return nil
}
