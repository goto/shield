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
		return servicedata.Key{}, err
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
