package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/goto/salt/audit"
	"github.com/goto/shield/core/activity"
	"github.com/goto/shield/pkg/db"
	newrelic "github.com/newrelic/go-agent"
)

type ActivityRepository struct {
	dbc *db.Client
}

func NewActivityRepository(dbc *db.Client) *ActivityRepository {
	return &ActivityRepository{
		dbc: dbc,
	}
}

func (r ActivityRepository) Insert(ctx context.Context, log *audit.Log) error {
	marshaledMetadata, err := json.Marshal(log.Metadata)
	if err != nil {
		return fmt.Errorf("%w: %s", parseErr, err)
	}

	marshaledData, err := json.Marshal(log.Data)
	if err != nil {
		return fmt.Errorf("%w: %s", parseErr, err)
	}

	query, params, err := dialect.Insert(TABLE_ACTIVITY).Rows(
		goqu.Record{
			"actor":     log.Actor,
			"action":    log.Action,
			"data":      marshaledData,
			"metadata":  marshaledMetadata,
			"timestamp": log.Timestamp,
		}).ToSQL()
	if err != nil {
		return fmt.Errorf("%w: %s", queryErr, err)
	}

	if err = r.dbc.WithTimeout(ctx, func(ctx context.Context) error {
		nrCtx := newrelic.FromContext(ctx)
		if nrCtx != nil {
			nr := newrelic.DatastoreSegment{
				Product:    newrelic.DatastorePostgres,
				Collection: TABLE_ACTIVITY,
				Operation:  "Create",
				StartTime:  nrCtx.StartSegmentNow(),
			}
			defer nr.End()
		}
		_, err = r.dbc.ExecContext(ctx, query, params...)
		return err
	}); err != nil {
		return checkPostgresError(err)
	}

	return nil
}

func (r ActivityRepository) List(ctx context.Context, filter activity.Filter) ([]audit.Log, error) {
	var fetchedActivity []audit.Log

	var defaultLimit int32 = 50
	var defaultPage int32 = 1
	if filter.Limit < 1 {
		filter.Limit = defaultLimit
	}
	if filter.Page < 1 {
		filter.Page = defaultPage
	}

	offset := (filter.Page - 1) * filter.Limit

	sqlStatement := dialect.From(TABLE_ACTIVITY)
	if filter.Actor != "" {
		sqlStatement = sqlStatement.Where(goqu.Ex{"actor": filter.Actor})
	}
	if filter.Action != "" {
		sqlStatement = sqlStatement.Where(goqu.Ex{"action": goqu.Op{"like": fmt.Sprintf("%%%s%%", filter.Action)}})
	}
	if filter.Data != nil {
		for key, value := range filter.Data {
			dataQuery := fmt.Sprintf("data->>'%s' = '%s'", key, value)
			sqlStatement = sqlStatement.Where(goqu.L(dataQuery))
		}
	}
	if filter.Metadata != nil {
		for key, value := range filter.Metadata {
			dataQuery := fmt.Sprintf("metadata->>'%s' = '%s'", key, value)
			sqlStatement = sqlStatement.Where(goqu.L(dataQuery))
		}
	}

	sqlStatement = sqlStatement.Where(
		goqu.Ex{"timestamp": goqu.Op{"between": goqu.Range(filter.StartTime, filter.EndTime)}},
	).Limit(uint(filter.Limit)).Offset(uint(offset)).Order(goqu.I("timestamp").Desc())
	query, params, err := sqlStatement.ToSQL()
	if err != nil {
		return nil, err
	}

	if err = r.dbc.WithTimeout(ctx, func(ctx context.Context) error {
		nrCtx := newrelic.FromContext(ctx)
		if nrCtx != nil {
			nr := newrelic.DatastoreSegment{
				Product:    newrelic.DatastorePostgres,
				Collection: TABLE_ACTIVITY,
				Operation:  "List",
				StartTime:  nrCtx.StartSegmentNow(),
			}
			defer nr.End()
		}

		return r.dbc.SelectContext(ctx, &fetchedActivity, query, params...)
	}); err != nil {
		err = checkPostgresError(err)
		if errors.Is(err, sql.ErrNoRows) {
			return []audit.Log{}, nil
		}

		return []audit.Log{}, fmt.Errorf("%w: %s", dbErr, err)
	}

	return fetchedActivity, nil
}
