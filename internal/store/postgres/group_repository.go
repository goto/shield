package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/goto/shield/core/group"
	"github.com/goto/shield/core/namespace"
	"github.com/goto/shield/core/organization"
	"github.com/goto/shield/core/relation"
	"github.com/goto/shield/internal/schema"
	"github.com/goto/shield/pkg/db"
	newrelic "github.com/newrelic/go-agent/v3/newrelic"
	"go.nhat.io/otelsql"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type GroupRepository struct {
	dbc *db.Client
}

type joinGroupMetadata struct {
	ID        string         `db:"id"`
	Name      string         `db:"name"`
	Slug      string         `db:"slug"`
	OrgId     string         `db:"org_id"`
	Key       any            `db:"key"`
	Value     sql.NullString `db:"value"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
}

func NewGroupRepository(dbc *db.Client) *GroupRepository {
	return &GroupRepository{
		dbc: dbc,
	}
}

func (r GroupRepository) GetByID(ctx context.Context, id string) (group.Group, error) {
	if strings.TrimSpace(id) == "" {
		return group.Group{}, group.ErrInvalidID
	}

	query, params, err := dialect.From(TABLE_GROUPS).Where(
		goqu.Ex{
			"id": id,
		}).ToSQL()
	if err != nil {
		return group.Group{}, fmt.Errorf("%w: %s", queryErr, err)
	}

	ctx = otelsql.WithCustomAttributes(
		ctx,
		[]attribute.KeyValue{
			attribute.String("db.repository.method", "GetByID"),
			attribute.String(string(semconv.DBSQLTableKey), TABLE_GROUPS),
		}...,
	)

	var groupModel Group
	if err = r.dbc.WithTimeout(ctx, func(ctx context.Context) error {
		nrCtx := newrelic.FromContext(ctx)
		if nrCtx != nil {
			nr := newrelic.DatastoreSegment{
				Product:    newrelic.DatastorePostgres,
				Collection: TABLE_GROUPS,
				Operation:  "GetByID",
				StartTime:  nrCtx.StartSegmentNow(),
			}
			defer nr.End()
		}

		return r.dbc.GetContext(ctx, &groupModel, query, params...)
	}); err != nil {
		err = checkPostgresError(err)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return group.Group{}, group.ErrNotExist
		case errors.Is(err, errInvalidTexRepresentation):
			return group.Group{}, group.ErrInvalidUUID
		default:
			return group.Group{}, err
		}
	}

	transformedGroup, err := groupModel.transformToGroup()
	if err != nil {
		return group.Group{}, fmt.Errorf("%w: %s", parseErr, err)
	}

	return transformedGroup, nil
}

func (r GroupRepository) GetBySlug(ctx context.Context, slug string) (group.Group, error) {
	if strings.TrimSpace(slug) == "" {
		return group.Group{}, group.ErrInvalidID
	}

	query, params, err := dialect.From(TABLE_GROUPS).Where(goqu.Ex{
		"slug": slug,
	}).ToSQL()
	if err != nil {
		return group.Group{}, fmt.Errorf("%w: %s", queryErr, err)
	}

	ctx = otelsql.WithCustomAttributes(
		ctx,
		[]attribute.KeyValue{
			attribute.String("db.repository.method", "GetBySlug"),
			attribute.String(string(semconv.DBSQLTableKey), TABLE_GROUPS),
		}...,
	)

	var groupModel Group
	if err = r.dbc.WithTimeout(ctx, func(ctx context.Context) error {
		nrCtx := newrelic.FromContext(ctx)
		if nrCtx != nil {
			nr := newrelic.DatastoreSegment{
				Product:    newrelic.DatastorePostgres,
				Collection: TABLE_GROUPS,
				Operation:  "GetBySlug",
				StartTime:  nrCtx.StartSegmentNow(),
			}
			defer nr.End()
		}

		return r.dbc.GetContext(ctx, &groupModel, query, params...)
	}); err != nil {
		err = checkPostgresError(err)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return group.Group{}, group.ErrNotExist
		default:
			return group.Group{}, err
		}
	}

	transformedGroup, err := groupModel.transformToGroup()
	if err != nil {
		return group.Group{}, fmt.Errorf("%w: %s", parseErr, err)
	}

	return transformedGroup, nil
}

func (r GroupRepository) GetByIDs(ctx context.Context, groupIDs []string) ([]group.Group, error) {
	var fetchedGroups []Group

	for _, id := range groupIDs {
		if strings.TrimSpace(id) == "" {
			return []group.Group{}, group.ErrInvalidID
		}
	}

	query, params, err := dialect.From(TABLE_GROUPS).Where(
		goqu.Ex{
			"id": goqu.Op{"in": groupIDs},
		}).ToSQL()
	if err != nil {
		return []group.Group{}, fmt.Errorf("%w: %s", queryErr, err)
	}

	ctx = otelsql.WithCustomAttributes(
		ctx,
		[]attribute.KeyValue{
			attribute.String("db.repository.method", "GetByIDs"),
			attribute.String(string(semconv.DBSQLTableKey), TABLE_GROUPS),
		}...,
	)

	// this query will return empty array of groups if no UUID is not matched
	// TODO: check and fox what to do in this scenerio
	if err = r.dbc.WithTimeout(ctx, func(ctx context.Context) error {
		nrCtx := newrelic.FromContext(ctx)
		if nrCtx != nil {
			nr := newrelic.DatastoreSegment{
				Product:    newrelic.DatastorePostgres,
				Collection: TABLE_GROUPS,
				Operation:  "GetByIDs",
				StartTime:  nrCtx.StartSegmentNow(),
			}
			defer nr.End()
		}

		return r.dbc.SelectContext(ctx, &fetchedGroups, query, params...)
	}); err != nil {
		err = checkPostgresError(err)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return []group.Group{}, group.ErrNotExist
		case errors.Is(err, errInvalidTexRepresentation):
			return []group.Group{}, group.ErrInvalidUUID
		default:
			return []group.Group{}, err
		}
	}

	var transformedGroups []group.Group
	for _, g := range fetchedGroups {
		transformedGroup, err := g.transformToGroup()
		if err != nil {
			return []group.Group{}, fmt.Errorf("%w: %s", parseErr, err)
		}

		transformedGroups = append(transformedGroups, transformedGroup)
	}

	return transformedGroups, nil
}

func (r GroupRepository) Create(ctx context.Context, grp group.Group) (group.Group, error) {
	if strings.TrimSpace(grp.Name) == "" || strings.TrimSpace(grp.Slug) == "" {
		return group.Group{}, group.ErrInvalidDetail
	}

	ctx = otelsql.WithCustomAttributes(
		ctx,
		[]attribute.KeyValue{
			attribute.String("db.repository.method", "Create"),
			attribute.String(string(semconv.DBSQLTableKey), TABLE_GROUPS),
		}...,
	)

	query, params, err := dialect.Insert(TABLE_GROUPS).Rows(
		goqu.Record{
			"name":     grp.Name,
			"slug":     grp.Slug,
			"org_id":   grp.OrganizationID,
			"metadata": nil,
		}).Returning(&Group{}).ToSQL()
	if err != nil {
		return group.Group{}, fmt.Errorf("%w: %s", queryErr, err)
	}

	var groupModel Group
	if err = r.dbc.WithTimeout(ctx, func(ctx context.Context) error {
		nrCtx := newrelic.FromContext(ctx)
		if nrCtx != nil {
			nr := newrelic.DatastoreSegment{
				Product:    newrelic.DatastorePostgres,
				Collection: TABLE_GROUPS,
				Operation:  "Create",
				StartTime:  nrCtx.StartSegmentNow(),
			}
			defer nr.End()
		}

		return r.dbc.QueryRowxContext(ctx, query, params...).StructScan(&groupModel)
	}); err != nil {
		err = checkPostgresError(err)
		switch {
		case errors.Is(err, errForeignKeyViolation):
			return group.Group{}, organization.ErrNotExist
		case errors.Is(err, errInvalidTexRepresentation):
			return group.Group{}, organization.ErrInvalidUUID
		case errors.Is(err, errDuplicateKey):
			return group.Group{}, group.ErrConflict
		default:
			return group.Group{}, err
		}
	}

	transformedGroup, err := groupModel.transformToGroup()
	if err != nil {
		return group.Group{}, fmt.Errorf("%w: %s", parseErr, err)
	}

	return transformedGroup, nil
}

func (r GroupRepository) List(ctx context.Context, flt group.Filter) ([]group.Group, error) {
	sqlStatement := dialect.From(TABLE_GROUPS).Select(
		goqu.I("id"),
		goqu.I("name"),
		goqu.I("slug"),
		goqu.I("org_id"),
		goqu.I("created_at"),
		goqu.I("updated_at"),
	)

	if len(flt.ServicedataKeyResourceIDs) > 0 {
		subquery := dialect.Select(
			goqu.I("sd.namespace_id"),
			goqu.I("sd.entity_id"),
			goqu.I("sk.name").As("name"),
			goqu.I("sd.value"),
			goqu.I("sk.resource_id"),
		).From(goqu.T(TABLE_SERVICE_DATA_KEYS).As("sk")).
			RightJoin(goqu.T(TABLE_SERVICE_DATA).As("sd"), goqu.On(
				goqu.I("sk.id").Eq(goqu.I("sd.key_id")))).
			Where(goqu.Ex{"sd.namespace_id": schema.GroupPrincipal},
				goqu.Ex{"sk.project_id": flt.ProjectID},
				goqu.L(
					"sk.resource_id",
				).In(flt.ServicedataKeyResourceIDs))

		sqlStatement = dialect.Select(
			goqu.I("g.id"),
			goqu.I("g.name"),
			goqu.I("g.slug"),
			goqu.I("g.org_id"),
			goqu.I("sd.name").As("key"),
			goqu.I("sd.value"),
			goqu.I("g.created_at"),
			goqu.I("g.updated_at"),
		).From(goqu.T(TABLE_GROUPS).As("g")).LeftJoin(subquery.As("sd"), goqu.On(
			goqu.Cast(goqu.C("id"), "TEXT").Eq(goqu.I("sd.entity_id"))))
	}

	if flt.OrganizationID != "" {
		sqlStatement = sqlStatement.Where(goqu.Ex{"org_id": flt.OrganizationID})
	}
	query, params, err := sqlStatement.ToSQL()
	if err != nil {
		return []group.Group{}, fmt.Errorf("%w: %s", queryErr, err)
	}

	ctx = otelsql.WithCustomAttributes(
		ctx,
		[]attribute.KeyValue{
			attribute.String("db.repository.method", "List"),
			attribute.String(string(semconv.DBSQLTableKey), TABLE_GROUPS),
		}...,
	)

	var fetchedJoinGroupMetadata []joinGroupMetadata
	if err = r.dbc.WithTimeout(ctx, func(ctx context.Context) error {
		nrCtx := newrelic.FromContext(ctx)
		if nrCtx != nil {
			nr := newrelic.DatastoreSegment{
				Product:    newrelic.DatastorePostgres,
				Collection: TABLE_GROUPS,
				Operation:  "List",
				StartTime:  nrCtx.StartSegmentNow(),
			}
			defer nr.End()
		}

		return r.dbc.SelectContext(ctx, &fetchedJoinGroupMetadata, query, params...)
	}); err != nil {
		err = checkPostgresError(err)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return []group.Group{}, nil
		case errors.Is(err, errInvalidTexRepresentation):
			return []group.Group{}, nil
		default:
			return []group.Group{}, fmt.Errorf("%w: %s", dbErr, err)
		}
	}

	groupedMetadataByGroup := make(map[string]group.Group)
	for _, g := range fetchedJoinGroupMetadata {
		if _, ok := groupedMetadataByGroup[g.ID]; !ok {
			groupedMetadataByGroup[g.ID] = group.Group{}
		}
		currentGroup := groupedMetadataByGroup[g.ID]
		currentGroup.ID = g.ID
		currentGroup.Slug = g.Slug
		currentGroup.Name = g.Name
		currentGroup.OrganizationID = g.OrgId
		currentGroup.CreatedAt = g.CreatedAt
		currentGroup.UpdatedAt = g.UpdatedAt

		if currentGroup.Metadata == nil {
			currentGroup.Metadata = make(map[string]any)
		}

		if g.Key != nil {
			var value any
			err := json.Unmarshal([]byte(g.Value.String), &value)
			if err != nil {
				continue
			}

			currentGroup.Metadata[g.Key.(string)] = value
		}

		groupedMetadataByGroup[g.ID] = currentGroup
	}

	var transformedGroups []group.Group
	for _, group := range groupedMetadataByGroup {
		transformedGroups = append(transformedGroups, group)
	}

	return transformedGroups, nil
}

func (r GroupRepository) UpdateByID(ctx context.Context, grp group.Group) (group.Group, error) {
	if strings.TrimSpace(grp.ID) == "" {
		return group.Group{}, group.ErrInvalidID
	}

	if strings.TrimSpace(grp.Name) == "" || strings.TrimSpace(grp.Slug) == "" {
		return group.Group{}, group.ErrInvalidDetail
	}

	query, params, err := dialect.Update(TABLE_GROUPS).Set(
		goqu.Record{
			"name":       grp.Name,
			"slug":       grp.Slug,
			"org_id":     grp.OrganizationID,
			"updated_at": goqu.L("now()"),
		}).Where(goqu.ExOr{
		"id": grp.ID,
	}).Returning(&Group{}).ToSQL()
	if err != nil {
		return group.Group{}, fmt.Errorf("%w: %s", queryErr, err)
	}

	var groupModel Group
	if err = r.dbc.WithTimeout(ctx, func(ctx context.Context) error {
		nrCtx := newrelic.FromContext(ctx)
		if nrCtx != nil {
			nr := newrelic.DatastoreSegment{
				Product:    newrelic.DatastorePostgres,
				Collection: TABLE_GROUPS,
				Operation:  "UpdateByID",
				StartTime:  nrCtx.StartSegmentNow(),
			}
			defer nr.End()
		}

		return r.dbc.QueryRowxContext(ctx, query, params...).StructScan(&groupModel)
	}); err != nil {
		err = checkPostgresError(err)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return group.Group{}, group.ErrNotExist
		case errors.Is(err, errInvalidTexRepresentation):
			return group.Group{}, group.ErrInvalidUUID
		case errors.Is(err, errDuplicateKey):
			return group.Group{}, group.ErrConflict
		case errors.Is(err, errForeignKeyViolation):
			return group.Group{}, organization.ErrNotExist
		default:
			return group.Group{}, fmt.Errorf("%w: %s", dbErr, err)
		}
	}

	updated, err := groupModel.transformToGroup()
	if err != nil {
		return group.Group{}, fmt.Errorf("%s: %w", parseErr, err)
	}

	return updated, nil
}

func (r GroupRepository) UpdateBySlug(ctx context.Context, grp group.Group) (group.Group, error) {
	if strings.TrimSpace(grp.Slug) == "" {
		return group.Group{}, group.ErrInvalidID
	}

	if strings.TrimSpace(grp.Name) == "" {
		return group.Group{}, group.ErrInvalidDetail
	}

	query, params, err := dialect.Update(TABLE_GROUPS).Set(
		goqu.Record{
			"name":       grp.Name,
			"org_id":     grp.OrganizationID,
			"updated_at": goqu.L("now()"),
		}).Where(goqu.Ex{
		"slug": grp.Slug,
	}).Returning(&Group{}).ToSQL()
	if err != nil {
		return group.Group{}, fmt.Errorf("%w: %s", queryErr, err)
	}

	var groupModel Group
	if err = r.dbc.WithTimeout(ctx, func(ctx context.Context) error {
		nrCtx := newrelic.FromContext(ctx)
		if nrCtx != nil {
			nr := newrelic.DatastoreSegment{
				Product:    newrelic.DatastorePostgres,
				Collection: TABLE_GROUPS,
				Operation:  "GetBySlug",
				StartTime:  nrCtx.StartSegmentNow(),
			}
			defer nr.End()
		}

		return r.dbc.QueryRowxContext(ctx, query, params...).StructScan(&groupModel)
	}); err != nil {
		err = checkPostgresError(err)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return group.Group{}, group.ErrNotExist
		case errors.Is(err, errInvalidTexRepresentation):
			return group.Group{}, organization.ErrInvalidUUID
		case errors.Is(err, errDuplicateKey):
			return group.Group{}, group.ErrConflict
		case errors.Is(err, errForeignKeyViolation):
			return group.Group{}, organization.ErrNotExist
		default:
			return group.Group{}, fmt.Errorf("%w: %s", dbErr, err)
		}
	}

	updated, err := groupModel.transformToGroup()
	if err != nil {
		return group.Group{}, fmt.Errorf("%s: %w", parseErr, err)
	}

	return updated, nil
}

func (r GroupRepository) ListUserGroups(ctx context.Context, userID string, roleID string) ([]group.Group, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, group.ErrInvalidID
	}

	sqlStatement := dialect.Select(
		goqu.I("g.id").As("id"),
		goqu.I("g.metadata").As("metadata"),
		goqu.I("g.name").As("name"),
		goqu.I("g.slug").As("slug"),
		goqu.I("g.updated_at").As("updated_at"),
		goqu.I("g.created_at").As("created_at"),
		goqu.I("g.org_id").As("org_id"),
	).
		From(goqu.T(TABLE_RELATIONS).As("r")).
		Join(goqu.T(TABLE_GROUPS).As("g"), goqu.On(
			goqu.I("g.id").Cast("VARCHAR").
				Eq(goqu.I("r.object_id")),
		)).
		Where(goqu.Ex{
			"r.object_namespace_id": namespace.DefinitionTeam.ID,
			"subject_namespace_id":  namespace.DefinitionUser.ID,
			"subject_id":            userID,
		})

	if strings.TrimSpace(roleID) != "" {
		sqlStatement = sqlStatement.Where(goqu.Ex{
			"role_id": roleID,
		})
	}

	query, params, err := sqlStatement.ToSQL()
	if err != nil {
		return []group.Group{}, fmt.Errorf("%w: %s", queryErr, err)
	}

	var fetchedGroups []Group
	if err = r.dbc.WithTimeout(ctx, func(ctx context.Context) error {
		nrCtx := newrelic.FromContext(ctx)
		if nrCtx != nil {
			nr := newrelic.DatastoreSegment{
				Product:    newrelic.DatastorePostgres,
				Collection: TABLE_GROUPS,
				Operation:  "ListUserGroups",
				StartTime:  nrCtx.StartSegmentNow(),
			}
			defer nr.End()
		}

		return r.dbc.SelectContext(ctx, &fetchedGroups, query, params...)
	}); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []group.Group{}, nil
		}
		return []group.Group{}, fmt.Errorf("%w: %s", dbErr, err)
	}

	var transformedGroups []group.Group
	for _, v := range fetchedGroups {
		transformedGroup, err := v.transformToGroup()
		if err != nil {
			return []group.Group{}, fmt.Errorf("%w: %s", parseErr, err)
		}
		transformedGroups = append(transformedGroups, transformedGroup)
	}

	return transformedGroups, nil
}

func (r GroupRepository) ListGroupRelations(ctx context.Context, objectId string, subject_type string, role string) ([]relation.RelationV2, error) {
	whereClauseExp := goqu.Ex{}
	whereClauseExp["object_id"] = objectId
	whereClauseExp["object_namespace_id"] = schema.GroupNamespace

	if subject_type != "" {
		if subject_type == "user" {
			whereClauseExp["subject_namespace_id"] = schema.UserPrincipal
		} else if subject_type == "group" {
			whereClauseExp["subject_namespace_id"] = schema.GroupPrincipal
		}
	}

	if role != "" {
		like := "%:" + role
		whereClauseExp["role_id"] = goqu.Op{"like": like}
	}

	query, params, err := dialect.Select(&relationCols{}).From(TABLE_RELATIONS).Where(whereClauseExp).ToSQL()
	if err != nil {
		return []relation.RelationV2{}, fmt.Errorf("%w: %s", queryErr, err)
	}

	var fetchedRelations []Relation
	if err = r.dbc.WithTimeout(ctx, func(ctx context.Context) error {
		nrCtx := newrelic.FromContext(ctx)
		if nrCtx != nil {
			nr := newrelic.DatastoreSegment{
				Product:    newrelic.DatastorePostgres,
				Collection: TABLE_GROUPS,
				Operation:  "ListGroupRelations",
				StartTime:  nrCtx.StartSegmentNow(),
			}
			defer nr.End()
		}

		return r.dbc.SelectContext(ctx, &fetchedRelations, query, params...)
	}); err != nil {
		// List should return empty list and no error instead
		if errors.Is(err, sql.ErrNoRows) {
			return []relation.RelationV2{}, nil
		}
		return []relation.RelationV2{}, fmt.Errorf("%w: %s", dbErr, err)
	}

	var transformedRelations []relation.RelationV2
	for _, r := range fetchedRelations {
		transformedRelations = append(transformedRelations, r.transformToRelationV2())
	}

	return transformedRelations, nil
}
