package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/goto/shield/core/resource"
	"github.com/goto/shield/internal/schema"
	"github.com/goto/shield/pkg/db"
	"github.com/goto/shield/pkg/uuid"
	newrelic "github.com/newrelic/go-agent/v3/newrelic"
	"go.nhat.io/otelsql"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type ResourceRepository struct {
	dbc *db.Client
}

func NewResourceRepository(dbc *db.Client) *ResourceRepository {
	return &ResourceRepository{
		dbc: dbc,
	}
}

func (r ResourceRepository) Upsert(ctx context.Context, res resource.Resource) (resource.Resource, error) {
	if strings.TrimSpace(res.URN) == "" {
		return resource.Resource{}, resource.ErrInvalidURN
	}

	userID := sql.NullString{String: res.UserID, Valid: res.UserID != ""}

	query, params, err := dialect.Insert(TABLE_RESOURCES).Rows(
		goqu.Record{
			"urn":          res.URN,
			"name":         res.Name,
			"project_id":   res.ProjectID,
			"org_id":       res.OrganizationID,
			"namespace_id": res.NamespaceID,
			"user_id":      userID,
		}).OnConflict(
		goqu.DoUpdate("ON CONSTRAINT resources_urn_unique", goqu.Record{
			"name":         res.Name,
			"project_id":   res.ProjectID,
			"org_id":       res.OrganizationID,
			"namespace_id": res.NamespaceID,
			"user_id":      userID,
		})).Returning(&ResourceCols{}).ToSQL()
	if err != nil {
		return resource.Resource{}, fmt.Errorf("%w: %s", queryErr, err)
	}

	ctx = otelsql.WithCustomAttributes(
		ctx,
		[]attribute.KeyValue{
			attribute.String("db.repository.method", "Upsert"),
			attribute.String(string(semconv.DBSQLTableKey), TABLE_RESOURCES),
		}...,
	)

	var resourceModel Resource
	if err = r.dbc.WithTimeout(ctx, func(ctx context.Context) error {
		nrCtx := newrelic.FromContext(ctx)
		if nrCtx != nil {
			nr := newrelic.DatastoreSegment{
				Product:    newrelic.DatastorePostgres,
				Collection: TABLE_RESOURCES,
				Operation:  "Create",
				StartTime:  nrCtx.StartSegmentNow(),
			}
			defer nr.End()
		}

		return r.dbc.QueryRowxContext(ctx, query, params...).StructScan(&resourceModel)
	}); err != nil {
		err = checkPostgresError(err)
		switch {
		case errors.Is(err, errForeignKeyViolation):
			return resource.Resource{}, resource.ErrInvalidDetail
		case errors.Is(err, errInvalidTexRepresentation):
			return resource.Resource{}, resource.ErrInvalidUUID
		default:
			return resource.Resource{}, err
		}
	}

	return resourceModel.transformToResource(), nil
}

func (r ResourceRepository) Create(ctx context.Context, res resource.Resource) (resource.Resource, error) {
	if strings.TrimSpace(res.URN) == "" {
		return resource.Resource{}, resource.ErrInvalidURN
	}

	userID := sql.NullString{String: res.UserID, Valid: res.UserID != ""}

	query, params, err := dialect.Insert(TABLE_RESOURCES).Rows(
		goqu.Record{
			"urn":          res.URN,
			"name":         res.Name,
			"project_id":   res.ProjectID,
			"org_id":       res.OrganizationID,
			"namespace_id": res.NamespaceID,
			"user_id":      userID,
		}).Returning(&ResourceCols{}).ToSQL()
	if err != nil {
		return resource.Resource{}, fmt.Errorf("%w: %s", queryErr, err)
	}

	ctx = otelsql.WithCustomAttributes(
		ctx,
		[]attribute.KeyValue{
			attribute.String("db.repository.method", "Create"),
			attribute.String(string(semconv.DBSQLTableKey), TABLE_RESOURCES),
		}...,
	)

	var resourceModel Resource
	if err = r.dbc.WithTimeout(ctx, func(ctx context.Context) error {
		nrCtx := newrelic.FromContext(ctx)
		if nrCtx != nil {
			nr := newrelic.DatastoreSegment{
				Product:    newrelic.DatastorePostgres,
				Collection: TABLE_RESOURCES,
				Operation:  "Create",
				StartTime:  nrCtx.StartSegmentNow(),
			}
			defer nr.End()
		}

		return r.dbc.QueryRowxContext(ctx, query, params...).StructScan(&resourceModel)
	}); err != nil {
		err = checkPostgresError(err)
		switch {
		case errors.Is(err, errForeignKeyViolation):
			return resource.Resource{}, resource.ErrInvalidDetail
		case errors.Is(err, errInvalidTexRepresentation):
			return resource.Resource{}, resource.ErrInvalidUUID
		case errors.Is(err, errDuplicateKey):
			return resource.Resource{}, resource.ErrConflict
		default:
			return resource.Resource{}, err
		}
	}

	return resourceModel.transformToResource(), nil
}

func (r ResourceRepository) List(ctx context.Context, flt resource.Filter) ([]resource.Resource, error) {
	var fetchedResources []Resource

	var defaultLimit int32 = 50
	var defaultPage int32 = 1
	if flt.Limit < 1 {
		flt.Limit = defaultLimit
	}
	if flt.Page < 1 {
		flt.Page = defaultPage
	}

	offset := (flt.Page - 1) * flt.Limit

	sqlStatement := dialect.From(TABLE_RESOURCES)
	if flt.ProjectID != "" {
		sqlStatement = sqlStatement.Where(goqu.Ex{"project_id": flt.ProjectID})
	}
	if flt.GroupID != "" {
		sqlStatement = sqlStatement.Where(goqu.Ex{"group_id": flt.GroupID})
	}
	if flt.OrganizationID != "" {
		sqlStatement = sqlStatement.Where(goqu.Ex{"org_id": flt.OrganizationID})
	}
	if flt.NamespaceID != "" {
		sqlStatement = sqlStatement.Where(goqu.Ex{"namespace_id": flt.NamespaceID})
	}
	sqlStatement = sqlStatement.Limit(uint(flt.Limit)).Offset(uint(offset))
	query, params, err := sqlStatement.ToSQL()
	if err != nil {
		return nil, err
	}

	ctx = otelsql.WithCustomAttributes(
		ctx,
		[]attribute.KeyValue{
			attribute.String("db.repository.method", "List"),
			attribute.String(string(semconv.DBSQLTableKey), TABLE_RESOURCES),
		}...,
	)

	if err = r.dbc.WithTimeout(ctx, func(ctx context.Context) error {
		nrCtx := newrelic.FromContext(ctx)
		if nrCtx != nil {
			nr := newrelic.DatastoreSegment{
				Product:    newrelic.DatastorePostgres,
				Collection: TABLE_RESOURCES,
				Operation:  "List",
				StartTime:  nrCtx.StartSegmentNow(),
			}
			defer nr.End()
		}

		return r.dbc.SelectContext(ctx, &fetchedResources, query, params...)
	}); err != nil {
		err = checkPostgresError(err)
		if errors.Is(err, sql.ErrNoRows) {
			return []resource.Resource{}, nil
		}
		if errors.Is(err, errInvalidTexRepresentation) {
			return []resource.Resource{}, nil
		}
		return []resource.Resource{}, fmt.Errorf("%w: %s", dbErr, err)
	}

	var transformedResources []resource.Resource
	for _, r := range fetchedResources {
		transformedResources = append(transformedResources, r.transformToResource())
	}

	return transformedResources, nil
}

func (r ResourceRepository) GetByID(ctx context.Context, id string) (resource.Resource, error) {
	if strings.TrimSpace(id) == "" {
		return resource.Resource{}, resource.ErrInvalidID
	}

	query, params, err := dialect.From(TABLE_RESOURCES).Where(goqu.Ex{
		"id": id,
	}).ToSQL()
	if err != nil {
		return resource.Resource{}, fmt.Errorf("%w: %s", queryErr, err)
	}

	ctx = otelsql.WithCustomAttributes(
		ctx,
		[]attribute.KeyValue{
			attribute.String("db.repository.method", "GetByID"),
			attribute.String(string(semconv.DBSQLTableKey), TABLE_RESOURCES),
		}...,
	)

	var resourceModel Resource
	if err = r.dbc.WithTimeout(ctx, func(ctx context.Context) error {
		nrCtx := newrelic.FromContext(ctx)
		if nrCtx != nil {
			nr := newrelic.DatastoreSegment{
				Product:    newrelic.DatastorePostgres,
				Collection: TABLE_RESOURCES,
				Operation:  "GetByID",
				StartTime:  nrCtx.StartSegmentNow(),
			}
			defer nr.End()
		}

		return r.dbc.GetContext(ctx, &resourceModel, query, params...)
	}); err != nil {
		err = checkPostgresError(err)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return resource.Resource{}, resource.ErrNotExist
		case errors.Is(err, errInvalidTexRepresentation):
			return resource.Resource{}, resource.ErrInvalidUUID
		default:
			return resource.Resource{}, err
		}
	}

	return resourceModel.transformToResource(), nil
}

func (r ResourceRepository) Update(ctx context.Context, id string, res resource.Resource) (resource.Resource, error) {
	if strings.TrimSpace(id) == "" {
		return resource.Resource{}, resource.ErrInvalidID
	}

	if !uuid.IsValid(id) {
		return resource.Resource{}, resource.ErrInvalidUUID
	}

	query, params, err := dialect.Update(TABLE_RESOURCES).Set(
		goqu.Record{
			"name":         res.Name,
			"project_id":   res.ProjectID,
			"org_id":       res.OrganizationID,
			"namespace_id": res.NamespaceID,
		},
	).Where(goqu.Ex{
		"id": id,
	}).Returning(&ResourceCols{}).ToSQL()
	if err != nil {
		return resource.Resource{}, fmt.Errorf("%w: %s", queryErr, err)
	}

	ctx = otelsql.WithCustomAttributes(
		ctx,
		[]attribute.KeyValue{
			attribute.String("db.repository.method", "Update"),
			attribute.String(string(semconv.DBSQLTableKey), TABLE_RESOURCES),
		}...,
	)

	var resourceModel Resource
	if err = r.dbc.WithTimeout(ctx, func(ctx context.Context) error {
		nrCtx := newrelic.FromContext(ctx)
		if nrCtx != nil {
			nr := newrelic.DatastoreSegment{
				Product:    newrelic.DatastorePostgres,
				Collection: TABLE_RESOURCES,
				Operation:  "Update",
				StartTime:  nrCtx.StartSegmentNow(),
			}
			defer nr.End()
		}

		return r.dbc.QueryRowxContext(ctx, query, params...).StructScan(&resourceModel)
	}); err != nil {
		err = checkPostgresError(err)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return resource.Resource{}, resource.ErrNotExist
		case errors.Is(err, errDuplicateKey):
			return resource.Resource{}, resource.ErrConflict
		case errors.Is(err, errForeignKeyViolation):
			return resource.Resource{}, resource.ErrNotExist
		case errors.Is(err, errInvalidTexRepresentation):
			return resource.Resource{}, resource.ErrInvalidDetail
		default:
			return resource.Resource{}, err
		}
	}

	return resourceModel.transformToResource(), nil
}

func (r ResourceRepository) GetByURN(ctx context.Context, urn string) (resource.Resource, error) {
	if strings.TrimSpace(urn) == "" {
		return resource.Resource{}, resource.ErrInvalidURN
	}

	query, params, err := dialect.Select(&ResourceCols{}).From(TABLE_RESOURCES).Where(
		goqu.Ex{
			"urn": urn,
		}).ToSQL()
	if err != nil {
		return resource.Resource{}, fmt.Errorf("%w: %s", queryErr, err)
	}

	ctx = otelsql.WithCustomAttributes(
		ctx,
		[]attribute.KeyValue{
			attribute.String("db.repository.method", "GetByURN"),
			attribute.String(string(semconv.DBSQLTableKey), TABLE_RESOURCES),
		}...,
	)

	var resourceModel Resource
	if err = r.dbc.WithTimeout(ctx, func(ctx context.Context) error {
		nrCtx := newrelic.FromContext(ctx)
		if nrCtx != nil {
			nr := newrelic.DatastoreSegment{
				Product:    newrelic.DatastorePostgres,
				Collection: TABLE_RESOURCES,
				Operation:  "GetByURN",
				StartTime:  nrCtx.StartSegmentNow(),
			}
			defer nr.End()
		}

		return r.dbc.GetContext(ctx, &resourceModel, query, params...)
	}); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return resource.Resource{}, resource.ErrNotExist
		}
		return resource.Resource{}, err
	}

	return resourceModel.transformToResource(), nil
}

func buildGetResourcesByNamespaceQuery(dialect goqu.DialectWrapper, name string, namespace string) (string, interface{}, error) {
	getResourcesByURNQuery, params, err := dialect.Select(&ResourceCols{}).From(TABLE_RESOURCES).Where(goqu.Ex{
		"name":         name,
		"namespace_id": namespace,
	}).ToSQL()

	return getResourcesByURNQuery, params, err
}

func (r ResourceRepository) GetByNamespace(ctx context.Context, name string, ns string) (resource.Resource, error) {
	var fetchedResource Resource

	query, _, err := buildGetResourcesByNamespaceQuery(dialect, name, ns)
	if err != nil {
		return resource.Resource{}, fmt.Errorf("%w: %s", queryErr, err)
	}

	ctx = otelsql.WithCustomAttributes(
		ctx,
		[]attribute.KeyValue{
			attribute.String("db.repository.method", "GetByNamespace"),
			attribute.String(string(semconv.DBSQLTableKey), TABLE_RESOURCES),
		}...,
	)

	err = r.dbc.WithTimeout(ctx, func(ctx context.Context) error {
		nrCtx := newrelic.FromContext(ctx)
		if nrCtx != nil {
			nr := newrelic.DatastoreSegment{
				Product:    newrelic.DatastorePostgres,
				Collection: TABLE_RESOURCES,
				Operation:  "GetByNamespace",
				StartTime:  nrCtx.StartSegmentNow(),
			}
			defer nr.End()
		}

		return r.dbc.GetContext(ctx, &fetchedResource, query)
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return resource.Resource{}, resource.ErrNotExist
		}

		return resource.Resource{}, fmt.Errorf("%w: %s", dbErr, err)
	}

	return fetchedResource.transformToResource(), nil
}

func (r ResourceRepository) UpsertConfig(ctx context.Context, name string, config schema.NamespaceConfigMapType) (resource.Config, error) {
	configJson, err := json.Marshal(config)
	if err != nil {
		return resource.Config{}, resource.ErrMarshal
	}

	query, params, err := goqu.Insert(TABLE_RESOURCE_CONFIGS).Rows(
		goqu.Record{"name": name, "config": configJson},
	).OnConflict(
		goqu.DoUpdate("name", goqu.Record{"name": name, "config": configJson})).Returning(&RuleConfig{}).ToSQL()
	if err != nil {
		return resource.Config{}, fmt.Errorf("%w: %s", queryErr, err)
	}

	ctx = otelsql.WithCustomAttributes(
		ctx,
		[]attribute.KeyValue{
			attribute.String("db.repository.method", "Upsert"),
			attribute.String(string(semconv.DBSQLTableKey), TABLE_RESOURCE_CONFIGS),
		}...,
	)

	var resourceConfigModel ResourceConfig
	if err = r.dbc.WithTimeout(ctx, func(ctx context.Context) error {
		nrCtx := newrelic.FromContext(ctx)
		if nrCtx != nil {
			nr := newrelic.DatastoreSegment{
				Product:    newrelic.DatastorePostgres,
				Collection: TABLE_RESOURCE_CONFIGS,
				Operation:  "Upsert",
				StartTime:  nrCtx.StartSegmentNow(),
			}
			defer nr.End()
		}
		return r.dbc.QueryRowxContext(ctx, query, params...).StructScan(&resourceConfigModel)
	}); err != nil {
		err = checkPostgresError(err)
		return resource.Config{}, err
	}

	return resourceConfigModel.transformToResourceConfig(), nil
}

func (r ResourceRepository) GetSchema(ctx context.Context) (schema.NamespaceConfigMapType, error) {
	query, params, err := dialect.From(TABLE_RESOURCE_CONFIGS).ToSQL()
	if err != nil {
		return schema.NamespaceConfigMapType{}, err
	}

	ctx = otelsql.WithCustomAttributes(
		ctx,
		[]attribute.KeyValue{
			attribute.String("db.repository.method", "List"),
			attribute.String(string(semconv.DBSQLTableKey), TABLE_RULE_CONFIGS),
		}...,
	)

	var resourceConfigModel []ResourceConfig
	if err = r.dbc.WithTimeout(ctx, func(ctx context.Context) error {
		nrCtx := newrelic.FromContext(ctx)
		if nrCtx != nil {
			nr := newrelic.DatastoreSegment{
				Product:    newrelic.DatastorePostgres,
				Collection: TABLE_RESOURCE_CONFIGS,
				Operation:  "List",
				StartTime:  nrCtx.StartSegmentNow(),
			}
			defer nr.End()
		}

		return r.dbc.SelectContext(ctx, &resourceConfigModel, query, params...)
	}); err != nil {
		err = checkPostgresError(err)
		if !errors.Is(err, sql.ErrNoRows) {
			return schema.NamespaceConfigMapType{}, err
		}
	}

	configMap := make(schema.NamespaceConfigMapType)
	for _, resourceConfig := range resourceConfigModel {
		var targetConfig schema.NamespaceConfigMapType
		if err := json.Unmarshal([]byte(resourceConfig.Config), &targetConfig); err != nil {
			return schema.NamespaceConfigMapType{}, err
		}
		configMap = schema.MergeNamespaceConfigMap(configMap, targetConfig)
	}

	return configMap, nil
}

func GetAll(ctx context.Context) ([]resource.YAML, error) {
	return []resource.YAML{}, nil
}

func (r *ResourceRepository) WithTransaction(ctx context.Context) context.Context {
	return r.dbc.WithTransaction(ctx, sql.TxOptions{})
}

func (r *ResourceRepository) Rollback(ctx context.Context, err error) error {
	if txErr := r.dbc.Rollback(ctx); txErr != nil {
		return fmt.Errorf("rollback error %s with error: %w", txErr.Error(), err)
	}
	return nil
}

func (r *ResourceRepository) Commit(ctx context.Context) error {
	return r.dbc.Commit(ctx)
}
