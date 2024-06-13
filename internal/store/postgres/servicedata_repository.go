package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/goto/shield/core/servicedata"
	"github.com/goto/shield/core/user"
	"github.com/goto/shield/internal/schema"
	"github.com/goto/shield/pkg/db"
	newrelic "github.com/newrelic/go-agent/v3/newrelic"
	"go.nhat.io/otelsql"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type ServiceDataRepository struct {
	dbc *db.Client
}

func NewServiceDataRepository(dbc *db.Client) *ServiceDataRepository {
	return &ServiceDataRepository{
		dbc: dbc,
	}
}

func (r ServiceDataRepository) CreateKey(ctx context.Context, key servicedata.Key) (servicedata.Key, error) {
	if len(key.Key) == 0 {
		return servicedata.Key{}, servicedata.ErrInvalidDetail
	}

	query, params, err := dialect.Insert(TABLE_SERVICE_DATA_KEYS).Rows(
		goqu.Record{
			"urn":         key.URN,
			"project_id":  key.ProjectID,
			"key":         key.Key,
			"description": key.Description,
			"resource_id": key.ResourceID,
		}).Returning(&Key{}).ToSQL()
	if err != nil {
		return servicedata.Key{}, queryErr
	}

	ctx = otelsql.WithCustomAttributes(
		ctx,
		[]attribute.KeyValue{
			attribute.String("db.repository.method", "CreateKey"),
			attribute.String(string(semconv.DBSQLTableKey), TABLE_SERVICE_DATA_KEYS),
		}...,
	)

	var serviceDataKeyModel Key
	if err = r.dbc.WithTimeout(ctx, func(ctx context.Context) error {
		nrCtx := newrelic.FromContext(ctx)
		if nrCtx != nil {
			nr := newrelic.DatastoreSegment{
				Product:    newrelic.DatastorePostgres,
				Collection: TABLE_SERVICE_DATA_KEYS,
				Operation:  "Create",
				StartTime:  nrCtx.StartSegmentNow(),
			}
			defer nr.End()
		}

		return r.dbc.QueryRowxContext(ctx, query, params...).StructScan(&serviceDataKeyModel)
	}); err != nil {
		err = checkPostgresError(err)
		switch {
		case errors.Is(err, errForeignKeyViolation):
			return servicedata.Key{}, servicedata.ErrInvalidDetail
		case errors.Is(err, errDuplicateKey):
			return servicedata.Key{}, servicedata.ErrConflict
		default:
			return servicedata.Key{}, err
		}
	}
	return serviceDataKeyModel.transformToServiceDataKey(), nil
}

func (r ServiceDataRepository) Upsert(ctx context.Context, data servicedata.ServiceData) (servicedata.ServiceData, error) {
	query, params, err := dialect.Insert(TABLE_SERVICE_DATA).Rows(
		goqu.Record{
			"namespace_id": data.NamespaceID,
			"entity_id":    data.EntityID,
			"key_id":       data.Key.ID,
			"value":        data.Value,
		},
	).OnConflict(goqu.DoUpdate(
		"ON CONSTRAINT servicedata_namespace_id_entity_id_key_id_key", goqu.Record{
			"key_id": data.Key.ID,
			"value":  data.Value,
		},
	)).Returning("value", goqu.L(`?`, data.Key.Key).As("key")).ToSQL()
	if err != nil {
		return servicedata.ServiceData{}, queryErr
	}

	ctx = otelsql.WithCustomAttributes(
		ctx,
		[]attribute.KeyValue{
			attribute.String("db.repository.method", "Upsert"),
			attribute.String(string(semconv.DBSQLTableKey), TABLE_SERVICE_DATA),
		}...,
	)

	var serviceDataModel ServiceData
	if err = r.dbc.WithTimeout(ctx, func(ctx context.Context) error {
		nrCtx := newrelic.FromContext(ctx)
		if nrCtx != nil {
			nr := newrelic.DatastoreSegment{
				Product:    newrelic.DatastorePostgres,
				Collection: TABLE_SERVICE_DATA,
				Operation:  "Upsert",
				StartTime:  nrCtx.StartSegmentNow(),
			}
			defer nr.End()
		}

		return r.dbc.QueryRowxContext(ctx, query, params...).StructScan(&serviceDataModel)
	}); err != nil {
		err = checkPostgresError(err)
		switch {
		default:
			return servicedata.ServiceData{}, err
		}
	}

	return serviceDataModel.transformToServiceData(), nil
}

func (r ServiceDataRepository) GetKeyByURN(ctx context.Context, URN string) (servicedata.Key, error) {
	query, params, err := dialect.From(TABLE_SERVICE_DATA_KEYS).Select().Where(goqu.Ex{
		"urn": URN,
	}).ToSQL()
	if err != nil {
		return servicedata.Key{}, queryErr
	}

	ctx = otelsql.WithCustomAttributes(
		ctx,
		[]attribute.KeyValue{
			attribute.String("db.repository.method", "GetKeyByURN"),
			attribute.String(string(semconv.DBSQLTableKey), TABLE_SERVICE_DATA_KEYS),
		}...,
	)

	var keyModel Key
	if err = r.dbc.WithTimeout(ctx, func(ctx context.Context) error {
		nrCtx := newrelic.FromContext(ctx)
		if nrCtx != nil {
			nr := newrelic.DatastoreSegment{
				Product:    newrelic.DatastorePostgres,
				Collection: TABLE_SERVICE_DATA_KEYS,
				Operation:  "GetKeyByURN",
				StartTime:  nrCtx.StartSegmentNow(),
			}
			defer nr.End()
		}

		return r.dbc.QueryRowxContext(ctx, query, params...).StructScan(&keyModel)
	}); err != nil {
		err = checkPostgresError(err)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return servicedata.Key{}, servicedata.ErrNotExist
		default:
			return servicedata.Key{}, err
		}
	}

	return keyModel.transformToServiceDataKey(), nil
}

func (r ServiceDataRepository) Get(ctx context.Context, filter servicedata.Filter) ([]servicedata.ServiceData, error) {
	if len(filter.EntityIDs) == 0 {
		return []servicedata.ServiceData{}, nil
	}

	sqlStatement := dialect.Select(
		goqu.I("sk.urn"),
		goqu.I("sk.project_id"),
		goqu.I("sk.resource_id"),
		goqu.I("sd.namespace_id"),
		goqu.I("sd.entity_id"),
		goqu.I("sk.key"),
		goqu.I("sd.value"),
	).From(goqu.T(TABLE_SERVICE_DATA).As("sd")).
		Join(goqu.T(TABLE_SERVICE_DATA_KEYS).As("sk"), goqu.On(
			goqu.I("sk.id").Eq(goqu.I("sd.key_id")))).
		Where(goqu.L(
			"(sd.namespace_id, sd.entity_id)",
		).In(filter.EntityIDs))

	if filter.Project != "" {
		sqlStatement = sqlStatement.Where(goqu.Ex{"sk.project_id": filter.Project})
	}

	query, params, err := sqlStatement.ToSQL()
	if err != nil {
		return []servicedata.ServiceData{}, err
	}

	ctx = otelsql.WithCustomAttributes(
		ctx,
		[]attribute.KeyValue{
			attribute.String("db.repository.method", "Get"),
			attribute.String(string(semconv.DBSQLTableKey), TABLE_SERVICE_DATA),
		}...,
	)

	var serviceDataModel []ServiceData
	if err = r.dbc.WithTimeout(ctx, func(ctx context.Context) error {
		nrCtx := newrelic.FromContext(ctx)
		if nrCtx != nil {
			nr := newrelic.DatastoreSegment{
				Product:    newrelic.DatastorePostgres,
				Collection: TABLE_SERVICE_DATA,
				Operation:  "Get",
				StartTime:  nrCtx.StartSegmentNow(),
			}
			defer nr.End()
		}

		return r.dbc.SelectContext(ctx, &serviceDataModel, query, params...)
	}); err != nil {
		err = checkPostgresError(err)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return []servicedata.ServiceData{}, nil
		case errors.Is(err, errInvalidTexRepresentation):
			return []servicedata.ServiceData{}, servicedata.ErrInvalidDetail
		default:
			return []servicedata.ServiceData{}, err
		}
	}

	var transformedServiceData []servicedata.ServiceData
	for _, sdm := range serviceDataModel {
		sd := sdm.transformToServiceData()
		transformedServiceData = append(transformedServiceData, sd)
	}

	return transformedServiceData, nil
}

func (r ServiceDataRepository) ListUsersHavingData(ctx context.Context, filter servicedata.ListUsersFilter, serviceData [][]string) ([]user.User, error) {
	var defaultLimit int32 = 50
	var defaultPage int32 = 1
	if filter.Limit < 1 {
		filter.Limit = defaultLimit
	}
	if filter.Page < 1 {
		filter.Page = defaultPage
	}

	_ = (filter.Page - 1) * filter.Limit

	subquery := dialect.Select("sd.entity_id").From(goqu.T(TABLE_SERVICE_DATA).As("sd")).
		Join(goqu.T(TABLE_SERVICE_DATA_KEYS).As("sk"), goqu.On(
			goqu.I("sk.id").Eq(goqu.I("sd.key_id")))).
		Where(goqu.Ex{"sk.project_id": filter.Project},
			goqu.Ex{"sd.namespace_id": schema.UserPrincipal},
			goqu.L(
				"(sk.key, sd.value)",
			).In(serviceData)).
		GroupBy(goqu.I("sd.entity_id"))

	// Main query
	query, params, err := dialect.Select(
		goqu.I("u.id"),
		goqu.I("u.name"),
		goqu.I("u.email"),
		goqu.I("u.created_at"),
		goqu.I("u.updated_at"),
	).From(goqu.T(TABLE_USERS).As("u")).Join(subquery.As("sd"), goqu.On(
		goqu.Cast(goqu.C("id"), "TEXT").Eq(goqu.I("sd.entity_id")))).ToSQL()
	if err != nil {
		return []user.User{}, nil
	}

	var users []User
	if err = r.dbc.WithTimeout(ctx, func(ctx context.Context) error {
		return r.dbc.SelectContext(ctx, &users, query, params...)
	}); err != nil {
		err = checkPostgresError(err)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return []user.User{}, nil
		case errors.Is(err, errInvalidTexRepresentation):
			return []user.User{}, servicedata.ErrInvalidDetail
		default:
			return []user.User{}, err
		}
	}

	return []user.User{}, nil
}

func (r ServiceDataRepository) WithTransaction(ctx context.Context) context.Context {
	return r.dbc.WithTransaction(ctx, sql.TxOptions{})
}

func (r ServiceDataRepository) Rollback(ctx context.Context, err error) error {
	if txErr := r.dbc.Rollback(ctx); txErr != nil {
		return fmt.Errorf("rollback error %s with error: %w", txErr.Error(), err)
	}
	return nil
}

func (r ServiceDataRepository) Commit(ctx context.Context) error {
	return r.dbc.Commit(ctx)
}
