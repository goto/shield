package postgres

import (
	"context"
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
		err = checkPostgresError(err)
		switch {
		case errors.Is(err, errInvalidTexRepresentation):
			return activity.ErrInvalidUUID
		default:
			return err
		}
	}

	return nil
}
