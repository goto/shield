package group

import (
	"context"
	"fmt"
	"strings"

	"github.com/goto/shield/core/action"
	"github.com/goto/shield/core/namespace"
	"github.com/goto/shield/core/relation"
	"github.com/goto/shield/core/user"
	"github.com/goto/shield/internal/schema"
	pkgctx "github.com/goto/shield/pkg/context"
	"github.com/goto/shield/pkg/str"
	"github.com/goto/shield/pkg/uuid"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
)

const (
	AuditKeyGroupCreate = "group.create"
	AuditKeyGroupUpdate = "group.update"
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
	Log(ctx context.Context, action string, actor string, data map[string]interface{}) error
}

type Service struct {
	repository      Repository
	relationService RelationService
	userService     UserService
	activityService ActivityService
	logger          *zap.SugaredLogger
}

func NewService(repository Repository, relationService RelationService, userService UserService, activityService ActivityService, logger *zap.SugaredLogger) *Service {
	return &Service{
		repository:      repository,
		relationService: relationService,
		userService:     userService,
		activityService: activityService,
		logger:          logger,
	}
}

func (s Service) Create(ctx context.Context, grp Group) (Group, error) {
	currentUser, err := s.userService.FetchCurrentUser(ctx)
	if err != nil {
		return Group{}, fmt.Errorf("%w: %s", user.ErrInvalidEmail, err.Error())
	}

	newGroup, err := s.repository.Create(ctx, grp)
	if err != nil {
		return Group{}, err
	}

	if err = s.addTeamToOrg(ctx, newGroup); err != nil {
		return Group{}, err
	}

	go func() {
		ctx := pkgctx.WithoutCancel(ctx)
		groupLogData := newGroup.ToGroupLogData()
		var logDataMap map[string]interface{}
		if err := mapstructure.Decode(groupLogData, &logDataMap); err != nil {
			s.logger.Errorf("%s: %s", ErrLogActivity.Error(), err.Error())
		}
		if err := s.activityService.Log(ctx, AuditKeyGroupCreate, currentUser.ID, logDataMap); err != nil {
			s.logger.Errorf("%s: %s", ErrLogActivity.Error(), err.Error())
		}
	}()

	return newGroup, nil
}

func (s Service) Get(ctx context.Context, idOrSlug string) (Group, error) {
	if uuid.IsValid(idOrSlug) {
		return s.repository.GetByID(ctx, idOrSlug)
	}
	return s.repository.GetBySlug(ctx, idOrSlug)
}

func (s Service) GetBySlug(ctx context.Context, slug string) (Group, error) {
	return s.repository.GetBySlug(ctx, slug)
}

func (s Service) GetByIDs(ctx context.Context, groupIDs []string) ([]Group, error) {
	return s.repository.GetByIDs(ctx, groupIDs)
}

func (s Service) List(ctx context.Context, flt Filter) ([]Group, error) {
	return s.repository.List(ctx, flt)
}

func (s Service) Update(ctx context.Context, grp Group) (Group, error) {
	currentUser, err := s.userService.FetchCurrentUser(ctx)
	if err != nil {
		return Group{}, fmt.Errorf("%w: %s", user.ErrInvalidEmail, err.Error())
	}

	if strings.TrimSpace(grp.ID) != "" {
		return s.repository.UpdateByID(ctx, grp)
	}

	updatedGroup, err := s.repository.UpdateBySlug(ctx, grp)
	if err != nil {
		return Group{}, err
	}

	go func() {
		ctx := pkgctx.WithoutCancel(ctx)
		groupLogData := updatedGroup.ToGroupLogData()
		var logDataMap map[string]interface{}
		if err := mapstructure.Decode(groupLogData, &logDataMap); err != nil {
			s.logger.Errorf("%s: %s", ErrLogActivity.Error(), err.Error())
		}
		if err := s.activityService.Log(ctx, AuditKeyGroupUpdate, currentUser.ID, logDataMap); err != nil {
			s.logger.Errorf("%s: %s", ErrLogActivity.Error(), err.Error())
		}
	}()

	return updatedGroup, nil
}

func (s Service) ListUserGroups(ctx context.Context, userId string, roleId string) ([]Group, error) {
	return s.repository.ListUserGroups(ctx, userId, roleId)
}

func (s Service) ListGroupRelations(ctx context.Context, objectId, subjectType, role string) ([]user.User, []Group, map[string][]string, map[string][]string, error) {
	relationList, err := s.repository.ListGroupRelations(ctx, objectId, subjectType, role)
	if err != nil {
		return []user.User{}, []Group{}, map[string][]string{}, map[string][]string{}, fmt.Errorf("%w: %s", ErrListingGroupRelations, err.Error())
	}

	userIDs := []string{}
	groupIDs := []string{}
	userIDRoleMap := map[string][]string{}
	groupIDRoleMap := map[string][]string{}
	users := []user.User{}
	groups := []Group{}

	for _, relation := range relationList {
		if relation.Subject.Namespace == schema.UserPrincipal {
			userIDs = append(userIDs, relation.Subject.ID)
			userIDRoleMap[relation.Subject.ID] = append(userIDRoleMap[relation.Subject.ID], relation.Subject.RoleID)
		} else if relation.Subject.Namespace == schema.GroupPrincipal {
			groupIDs = append(groupIDs, relation.Subject.ID)
			groupIDRoleMap[relation.Subject.ID] = append(groupIDRoleMap[relation.Subject.ID], relation.Subject.RoleID)
		}
	}

	if len(userIDs) > 0 {
		userList, err := s.userService.GetByIDs(ctx, userIDs)
		if err != nil {
			return []user.User{}, []Group{}, map[string][]string{}, map[string][]string{}, fmt.Errorf("%w: %s", ErrFetchingUsers, err.Error())
		}

		users = append(users, userList...)
	}

	if len(groupIDs) > 0 {
		groupList, err := s.repository.GetByIDs(ctx, groupIDs)
		if err != nil {
			return []user.User{}, []Group{}, map[string][]string{}, map[string][]string{}, fmt.Errorf("%w: %s", ErrFetchingGroups, err.Error())
		}

		groups = append(groups, groupList...)
	}

	return users, groups, userIDRoleMap, groupIDRoleMap, nil
}

func (s Service) addTeamToOrg(ctx context.Context, team Group) error {
	orgId := str.DefaultStringIfEmpty(team.OrganizationID, team.OrganizationID)
	rel := relation.RelationV2{
		Object: relation.Object{
			ID:          team.ID,
			NamespaceID: schema.GroupNamespace,
		},
		Subject: relation.Subject{
			ID:        orgId,
			Namespace: schema.OrganizationNamespace,
			RoleID:    schema.OrganizationRelationName,
		},
	}

	_, err := s.relationService.Create(ctx, rel)
	if err != nil {
		return err
	}

	return nil
}
