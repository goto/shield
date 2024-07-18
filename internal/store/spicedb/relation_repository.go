package spicedb

import (
	"context"
	"fmt"
	"io"

	"github.com/goto/shield/core/action"
	"github.com/goto/shield/core/relation"
	"github.com/goto/shield/internal/schema"
	"github.com/goto/shield/internal/store/spicedb/schema_generator"

	authzedpb "github.com/authzed/authzed-go/proto/authzed/api/v1"
	newrelic "github.com/newrelic/go-agent/v3/newrelic"
)

type RelationRepository struct {
	spiceDB *SpiceDB
}

const nrProductName = "spicedb"

func NewRelationRepository(spiceDB *SpiceDB) *RelationRepository {
	return &RelationRepository{
		spiceDB: spiceDB,
	}
}

func (r RelationRepository) Add(ctx context.Context, rel relation.Relation) error {
	relationship, err := schema_generator.TransformRelation(rel)
	if err != nil {
		return err
	}
	request := &authzedpb.WriteRelationshipsRequest{
		Updates: []*authzedpb.RelationshipUpdate{
			{
				Operation:    authzedpb.RelationshipUpdate_OPERATION_TOUCH,
				Relationship: relationship,
			},
		},
	}

	if _, err = r.spiceDB.client.WriteRelationships(ctx, request); err != nil {
		return err
	}

	return nil
}

func getRelation(a string) string {
	if a == schema.GroupPrincipal {
		return "membership"
	}

	return ""
}

func (r RelationRepository) AddV2(ctx context.Context, rel relation.RelationV2) error {
	relationship := &authzedpb.Relationship{
		Resource: &authzedpb.ObjectReference{
			ObjectType: rel.Object.NamespaceID,
			ObjectId:   rel.Object.ID,
		},
		Relation: schema.GetRoleName(rel.Subject.RoleID),
		Subject: &authzedpb.SubjectReference{
			Object: &authzedpb.ObjectReference{
				ObjectType: rel.Subject.Namespace,
				ObjectId:   rel.Subject.ID,
			},
			OptionalRelation: getRelation(rel.Subject.Namespace),
		},
	}
	request := &authzedpb.WriteRelationshipsRequest{
		Updates: []*authzedpb.RelationshipUpdate{
			{
				Operation:    authzedpb.RelationshipUpdate_OPERATION_TOUCH,
				Relationship: relationship,
			},
		},
	}

	nrCtx := newrelic.FromContext(ctx)
	if nrCtx != nil {
		nr := newrelic.DatastoreSegment{
			Product: nrProductName,
			QueryParameters: map[string]interface{}{
				"relation":          rel.Subject.RoleID,
				"subject_namespace": rel.Subject.Namespace,
				"object_namespace":  rel.Object.NamespaceID,
			},
			Operation: "Upsert_Relation",
			StartTime: nrCtx.StartSegmentNow(),
		}
		defer nr.End()
	}

	if _, err := r.spiceDB.client.WriteRelationships(ctx, request); err != nil {
		return err
	}

	return nil
}

func (r RelationRepository) Check(ctx context.Context, rel relation.Relation, act action.Action) (bool, error) {
	relationship, err := schema_generator.TransformCheckRelation(rel)
	if err != nil {
		return false, err
	}

	request := &authzedpb.CheckPermissionRequest{
		Consistency: &authzedpb.Consistency{
			Requirement: &authzedpb.Consistency_FullyConsistent{
				FullyConsistent: true,
			},
		},
		Resource:   relationship.Resource,
		Subject:    relationship.Subject,
		Permission: act.ID,
	}

	nrCtx := newrelic.FromContext(ctx)
	if nrCtx != nil {
		nr := newrelic.DatastoreSegment{
			Product:    nrProductName,
			Collection: fmt.Sprintf("object:%s::subject:%s", request.Resource.ObjectType, request.Subject.Object.ObjectType),
			Operation:  "Check",
			StartTime:  nrCtx.StartSegmentNow(),
		}
		defer nr.End()
	}

	response, err := r.spiceDB.client.CheckPermission(ctx, request)
	if err != nil {
		return false, err
	}

	return response.Permissionship == authzedpb.CheckPermissionResponse_PERMISSIONSHIP_HAS_PERMISSION, nil
}

func (r RelationRepository) BulkCheck(ctx context.Context, rels []relation.Relation, acts []action.Action) ([]relation.Permission, error) {
	var requestItems []*authzedpb.CheckBulkPermissionsRequestItem

	if len(rels) != len(acts) {
		return []relation.Permission{}, relation.ErrInvalidDetail
	}

	for i, rel := range rels {
		relationship, err := schema_generator.TransformCheckRelation(rel)
		if err != nil {
			return []relation.Permission{}, err
		}
		requestItems = append(requestItems, &authzedpb.CheckBulkPermissionsRequestItem{
			Resource:   relationship.Resource,
			Subject:    relationship.Subject,
			Permission: acts[i].ID,
		})
	}

	request := &authzedpb.CheckBulkPermissionsRequest{
		Consistency: &authzedpb.Consistency{
			Requirement: &authzedpb.Consistency_FullyConsistent{
				FullyConsistent: true,
			},
		},
		Items: requestItems,
	}

	nrCtx := newrelic.FromContext(ctx)
	if nrCtx != nil {
		nr := newrelic.DatastoreSegment{
			Product:   nrProductName,
			Operation: "Bulk_Check",
			StartTime: nrCtx.StartSegmentNow(),
		}
		defer nr.End()
	}

	response, err := r.spiceDB.client.CheckBulkPermissions(ctx, request)
	if err != nil {
		return []relation.Permission{}, err
	}

	var result []relation.Permission
	for _, res := range response.Pairs {
		switch res.Response.(type) {
		case *authzedpb.CheckBulkPermissionsPair_Error:
			result = append(result, relation.Permission{
				ObjectID:        res.Request.Resource.ObjectId,
				ObjectNamespace: res.Request.Resource.ObjectType,
				Permission:      res.Request.Permission,
				Allowed:         false,
			})
		case *authzedpb.CheckBulkPermissionsPair_Item:
			result = append(result, relation.Permission{
				ObjectID:        res.Request.Resource.ObjectId,
				ObjectNamespace: res.Request.Resource.ObjectType,
				Permission:      res.Request.Permission,
				Allowed:         res.Response.(*authzedpb.CheckBulkPermissionsPair_Item).Item.Permissionship == 2,
			})
		}
	}

	return result, nil
}

func (r RelationRepository) Delete(ctx context.Context, rel relation.Relation) error {
	relationship, err := schema_generator.TransformRelation(rel)
	if err != nil {
		return err
	}
	request := &authzedpb.DeleteRelationshipsRequest{
		RelationshipFilter: &authzedpb.RelationshipFilter{
			ResourceType:       relationship.Resource.ObjectType,
			OptionalResourceId: relationship.Resource.ObjectId,
			OptionalRelation:   relationship.Relation,
			OptionalSubjectFilter: &authzedpb.SubjectFilter{
				SubjectType:       relationship.Subject.Object.ObjectType,
				OptionalSubjectId: relationship.Subject.Object.ObjectId,
			},
		},
	}

	if _, err = r.spiceDB.client.DeleteRelationships(ctx, request); err != nil {
		return err
	}

	return nil
}

func (r RelationRepository) DeleteV2(ctx context.Context, rel relation.RelationV2) error {
	relationship, err := schema_generator.TransformRelationV2(rel)
	if err != nil {
		return err
	}
	request := &authzedpb.DeleteRelationshipsRequest{
		RelationshipFilter: &authzedpb.RelationshipFilter{
			ResourceType:       relationship.Resource.ObjectType,
			OptionalResourceId: relationship.Resource.ObjectId,
			OptionalRelation:   relationship.Relation,
			OptionalSubjectFilter: &authzedpb.SubjectFilter{
				SubjectType:       relationship.Subject.Object.ObjectType,
				OptionalSubjectId: relationship.Subject.Object.ObjectId,
			},
		},
	}

	nrCtx := newrelic.FromContext(ctx)
	if nrCtx != nil {
		nr := newrelic.DatastoreSegment{
			Product: nrProductName,
			QueryParameters: map[string]interface{}{
				"relation":          rel.Subject.RoleID,
				"subject_namespace": rel.Subject.Namespace,
				"object_namespace":  rel.Object.NamespaceID,
			},
			Operation: "Delete_Relation",
			StartTime: nrCtx.StartSegmentNow(),
		}
		defer nr.End()
	}
	_, err = r.spiceDB.client.DeleteRelationships(ctx, request)
	if err != nil {
		return err
	}

	return nil
}

func (r RelationRepository) DeleteSubjectRelations(ctx context.Context, resourceType, optionalResourceID string) error {
	request := &authzedpb.DeleteRelationshipsRequest{
		RelationshipFilter: &authzedpb.RelationshipFilter{
			ResourceType:       resourceType,
			OptionalResourceId: optionalResourceID,
		},
	}

	nrCtx := newrelic.FromContext(ctx)
	if nrCtx != nil {
		nr := newrelic.DatastoreSegment{
			Product: nrProductName,
			QueryParameters: map[string]interface{}{
				"object_namespace": resourceType,
				"object_id":        optionalResourceID,
			},
			Operation: "Delete_Subject_Relations",
			StartTime: nrCtx.StartSegmentNow(),
		}
		defer nr.End()
	}

	if _, err := r.spiceDB.client.DeleteRelationships(ctx, request); err != nil {
		return err
	}

	return nil
}

func (r RelationRepository) LookupResources(ctx context.Context, resourceType, permission, subjectType, subjectID string) ([]string, error) {
	request := &authzedpb.LookupResourcesRequest{
		Consistency: &authzedpb.Consistency{
			Requirement: &authzedpb.Consistency_FullyConsistent{
				FullyConsistent: true,
			},
		},
		ResourceObjectType: resourceType,
		Permission:         permission,
		Subject: &authzedpb.SubjectReference{
			Object: &authzedpb.ObjectReference{
				ObjectType: subjectType,
				ObjectId:   subjectID,
			},
		},
	}

	nrCtx := newrelic.FromContext(ctx)
	if nrCtx != nil {
		nr := newrelic.DatastoreSegment{
			Product:    nrProductName,
			Collection: fmt.Sprintf("object:%s::subject:%s", resourceType, subjectType),
			Operation:  "Lookup_Resources",
			StartTime:  nrCtx.StartSegmentNow(),
		}
		defer nr.End()
	}

	response, err := r.spiceDB.client.LookupResources(ctx, request)
	if err != nil {
		return []string{}, err
	}

	var res []string
	for {
		resp, err := response.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return []string{}, err
		}

		res = append(res, resp.ResourceObjectId)
	}

	return res, nil
}
