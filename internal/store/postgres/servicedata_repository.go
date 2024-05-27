package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/goto/shield/core/servicedata"
	"github.com/goto/shield/pkg/db"
	newrelic "github.com/newrelic/go-agent"
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
	query, params, err := dialect.From().
		With("key", dialect.From(TABLE_SERVICE_DATA_KEYS).Select("id", "urn").Where(goqu.Ex{"urn": data.Key.URN})).
		With("data", dialect.Insert(TABLE_SERVICE_DATA).
			Rows(
				goqu.Record{
					"namespace_id": data.NamespaceID,
					"entity_id":    data.EntityID,
					"key_id":       dialect.Select("id").From("key"),
					"value":        data.Value,
				},
			).OnConflict(goqu.DoUpdate(
			"ON CONSTRAINT servicedata_namespace_id_entity_id_key_id_key", goqu.Record{
				"key_id": dialect.Select("id").From("key"),
				"value":  data.Value,
			},
		)).Returning("id", "entity_id", "namespace_id", "key_id", "value", "created_at", "updated_at")).
		Select(&UpsertServiceData{}).From("data").Join(goqu.T("key"), goqu.On(goqu.I("key.id").Eq(goqu.I("data.key_id")))).ToSQL()
	if err != nil {
		return servicedata.ServiceData{}, err
	}

	var serviceDataModel UpsertServiceData
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
