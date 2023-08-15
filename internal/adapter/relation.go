package adapter

import (
	"context"
	"fmt"

	"github.com/goto/shield/core/group"
	"github.com/goto/shield/core/relation"
	"github.com/goto/shield/core/user"
	"github.com/goto/shield/internal/schema"
	"github.com/goto/shield/pkg/uuid"
)

type Relation struct {
	groupService    *group.Service
	userService     *user.Service
	relationService *relation.Service
}

func NewRelation(
	groupService *group.Service,
	userService *user.Service,
	relationService *relation.Service,
) *Relation {
	return &Relation{
		groupService:    groupService,
		userService:     userService,
		relationService: relationService,
	}
}

func (a Relation) TransformRelation(ctx context.Context, rlt relation.RelationV2) (relation.RelationV2, error) {
	rel := rlt

	// If Principal is a user, then we will get ID for that user as Subject.ID
	if rel.Subject.Namespace == schema.UserPrincipal || rel.Subject.Namespace == "user" {
		fetchedUser, err := a.userService.GetByEmail(ctx, rel.Subject.ID)
		if err != nil {
			return relation.RelationV2{}, fmt.Errorf("%w: %s", relation.ErrFetchingUser, err.Error())
		}

		rel.Subject.Namespace = schema.UserPrincipal
		rel.Subject.ID = fetchedUser.ID
	} else if rel.Subject.Namespace == schema.GroupPrincipal || rel.Subject.Namespace == "group" {
		// If Principal is a group, then we will get ID for that group as Subject.ID
		groupID := rel.Subject.ID

		if !uuid.IsValid(groupID) {
			fetchedGroup, err := a.groupService.Get(ctx, rel.Subject.ID)
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

// // If Principal is a user, then we will get ID for that user as Subject.ID
// if rel.Subject.Namespace == schema.UserPrincipal || rel.Subject.Namespace == "user" {
// 	fetchedUser, err := s.userService.GetByEmail(ctx, rel.Subject.ID)
// 	if err != nil {
// 		return RelationV2{}, fmt.Errorf("%w: %s", ErrFetchingUser, err.Error())
// 	}

// 	rel.Subject.Namespace = schema.UserPrincipal
// 	rel.Subject.ID = fetchedUser.ID
// }
