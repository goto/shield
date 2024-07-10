package adapter

import (
	"context"
	"fmt"
	"slices"

	"github.com/goto/shield/core/group"
	"github.com/goto/shield/core/relation"
	"github.com/goto/shield/core/role"
	"github.com/goto/shield/core/user"
	"github.com/goto/shield/internal/schema"
	"github.com/goto/shield/pkg/uuid"
)

const WILDCARD = "*"

type Relation struct {
	groupService    *group.Service
	userService     *user.Service
	relationService *relation.Service
	roleService     *role.Service
}

func NewRelation(
	groupService *group.Service,
	userService *user.Service,
	relationService *relation.Service,
	roleService *role.Service,
) *Relation {
	return &Relation{
		groupService:    groupService,
		userService:     userService,
		relationService: relationService,
		roleService:     roleService,
	}
}

func (a Relation) TransformRelation(ctx context.Context, rlt relation.RelationV2) (relation.RelationV2, error) {
	rel := rlt

	// If Principal is a user, then we will get ID for that user as Subject.ID
	if rel.Subject.Namespace == schema.UserPrincipal || rel.Subject.Namespace == "user" {
		userID := rel.Subject.ID

		if userID == WILDCARD {
			err := a.isWildCardAllowed(ctx, rel)
			if err != nil {
				return relation.RelationV2{}, err
			}
		} else if !uuid.IsValid(userID) {
			fetchedUser, err := a.userService.GetByEmail(ctx, rel.Subject.ID)
			if err != nil {
				return relation.RelationV2{}, fmt.Errorf("%w: %s", relation.ErrFetchingUser, err.Error())
			}
			userID = fetchedUser.ID
		}

		rel.Subject.Namespace = schema.UserPrincipal
		rel.Subject.ID = userID
	} else if rel.Subject.Namespace == schema.GroupPrincipal || rel.Subject.Namespace == "group" {
		// If Principal is a group, then we will get ID for that group as Subject.ID
		groupID := rel.Subject.ID

		if !uuid.IsValid(groupID) {
			fetchedGroup, err := a.groupService.GetBySlug(ctx, rel.Subject.ID)
			if err != nil {
				return relation.RelationV2{}, fmt.Errorf("%w on subject conversion: %s", relation.ErrFetchingGroup, err.Error())
			}
			groupID = fetchedGroup.ID
		}
		rel.Subject.Namespace = schema.GroupPrincipal
		rel.Subject.ID = groupID
	}

	// Group
	if rel.Object.NamespaceID == schema.GroupNamespace || rel.Object.NamespaceID == "group" {
		// If object is a group, then we will get ID for that group as Object.ID
		groupID := rel.Object.ID

		if !uuid.IsValid(groupID) {
			fetchedGroup, err := a.groupService.Get(ctx, rel.Subject.ID)
			if err != nil {
				return relation.RelationV2{}, fmt.Errorf("%w on object conversion:: %s", relation.ErrFetchingGroup, err.Error())
			}
			groupID = fetchedGroup.ID
		}
		rel.Object.NamespaceID = schema.GroupNamespace
		rel.Object.ID = groupID
	}

	return rel, nil
}

func (a Relation) isWildCardAllowed(ctx context.Context, rlt relation.RelationV2) error {
	roleID := rlt.Object.NamespaceID + ":" + rlt.Subject.RoleID
	role, err := a.roleService.Get(ctx, roleID)
	if err != nil {
		return fmt.Errorf("error fetching role: %s", err.Error())
	}
	if !slices.Contains(role.Types, schema.UserPrincipalWildcard) {
		return fmt.Errorf("%s does not allow wildcard for subject %s", rlt.Object.NamespaceID, rlt.Subject.Namespace)
	}

	return nil
}
