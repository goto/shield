package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/goto/shield/core/namespace"
	"github.com/goto/shield/core/policy"
	"github.com/goto/shield/pkg/db"
	newrelic "github.com/newrelic/go-agent/v3/newrelic"
	"go.nhat.io/otelsql"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type PolicyRepository struct {
	dbc *db.Client
}

func NewPolicyRepository(dbc *db.Client) *PolicyRepository {
	return &PolicyRepository{
		dbc: dbc,
	}
}

func (r PolicyRepository) buildListQuery() *goqu.SelectDataset {
	return dialect.Select(
		"p.id",
		"p.namespace_id",
		"p.action_id",
		"p.role_id",
	).From(goqu.T(TABLE_POLICIES).As("p"))
}

func (r PolicyRepository) Get(ctx context.Context, id string) (policy.Policy, error) {
	if strings.TrimSpace(id) == "" {
		return policy.Policy{}, policy.ErrInvalidID
	}

	query, params, err := r.buildListQuery().
		Where(
			goqu.Ex{
				"p.id": id,
			},
		).ToSQL()
	if err != nil {
		return policy.Policy{}, fmt.Errorf("%w: %s", queryErr, err)
	}

	ctx = otelsql.WithCustomAttributes(
		ctx,
		[]attribute.KeyValue{
			attribute.String("db.repository.method", "Get"),
			attribute.String(string(semconv.DBSQLTableKey), TABLE_POLICIES),
		}...,
	)

	var policyModel Policy
	if err = r.dbc.WithTimeout(ctx, func(ctx context.Context) error {
		nrCtx := newrelic.FromContext(ctx)
		if nrCtx != nil {
			nr := newrelic.DatastoreSegment{
				Product:    newrelic.DatastorePostgres,
				Collection: TABLE_POLICIES,
				Operation:  "Get",
				StartTime:  nrCtx.StartSegmentNow(),
			}
			defer nr.End()
		}

		return r.dbc.GetContext(ctx, &policyModel, query, params...)
	}); err != nil {
		err = checkPostgresError(err)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return policy.Policy{}, policy.ErrNotExist
		case errors.Is(err, errInvalidTexRepresentation):
			return policy.Policy{}, policy.ErrInvalidUUID
		default:
			return policy.Policy{}, err
		}
	}

	transformedPolicy, err := policyModel.transformToPolicy()
	if err != nil {
		return policy.Policy{}, fmt.Errorf("%w: %s", parseErr, err)
	}

	return transformedPolicy, nil
}

func (r PolicyRepository) List(ctx context.Context, filter policy.Filters) ([]policy.Policy, error) {
	var fetchedPolicies []Policy

	listQuery := r.buildListQuery()

	if filter.NamespaceID != "" {
		listQuery = listQuery.Where(goqu.Ex{"namespace_id": filter.NamespaceID})
	}

	query, params, err := listQuery.ToSQL()
	if err != nil {
		return []policy.Policy{}, fmt.Errorf("%w: %s", queryErr, err)
	}

	ctx = otelsql.WithCustomAttributes(
		ctx,
		[]attribute.KeyValue{
			attribute.String("db.repository.method", "List"),
			attribute.String(string(semconv.DBSQLTableKey), TABLE_POLICIES),
		}...,
	)

	if err = r.dbc.WithTimeout(ctx, func(ctx context.Context) error {
		nrCtx := newrelic.FromContext(ctx)
		if nrCtx != nil {
			nr := newrelic.DatastoreSegment{
				Product:    newrelic.DatastorePostgres,
				Collection: TABLE_POLICIES,
				Operation:  "List",
				StartTime:  nrCtx.StartSegmentNow(),
			}
			defer nr.End()
		}
		return r.dbc.SelectContext(ctx, &fetchedPolicies, query, params...)
	}); err != nil {
		err = checkPostgresError(err)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return []policy.Policy{}, nil
		default:
			return []policy.Policy{}, err
		}
	}

	var transformedPolicies []policy.Policy
	for _, p := range fetchedPolicies {
		transformedPolicy, err := p.transformToPolicy()
		if err != nil {
			return []policy.Policy{}, fmt.Errorf("%w: %s", parseErr, err)
		}
		transformedPolicies = append(transformedPolicies, transformedPolicy)
	}

	return transformedPolicies, nil
}

func (r PolicyRepository) Upsert(ctx context.Context, pol *policy.Policy) (string, error) {
	if pol == nil {
		return "", policy.ErrInvalidDetail
	}
	// TODO(krtkvrm) | IMP: need to find a way to deprecate this
	// This is required by bootstrap, which will be changed in this PR
	roleID := pol.RoleID
	actionID := pol.ActionID
	nsID := pol.NamespaceID

	if strings.TrimSpace(actionID) == "" {
		return "", policy.ErrInvalidDetail
	}

	ctx = otelsql.WithCustomAttributes(
		ctx,
		[]attribute.KeyValue{
			attribute.String("db.repository.method", "Upsert"),
			attribute.String(string(semconv.DBSQLTableKey), TABLE_POLICIES),
		}...,
	)

	query, params, err := dialect.Insert(TABLE_POLICIES).Rows(
		goqu.Record{
			"namespace_id": nsID,
			"role_id":      roleID,
			"action_id":    sql.NullString{String: actionID, Valid: actionID != ""},
		}).OnConflict(goqu.DoUpdate("role_id, namespace_id, action_id", goqu.Record{
		"namespace_id": nsID,
	})).Returning("id").ToSQL()
	if err != nil {
		return "", fmt.Errorf("%w: %s", queryErr, err)
	}

	var policyID string
	if err = r.dbc.WithTimeout(ctx, func(ctx context.Context) error {
		nrCtx := newrelic.FromContext(ctx)
		if nrCtx != nil {
			nr := newrelic.DatastoreSegment{
				Product:    newrelic.DatastorePostgres,
				Collection: TABLE_POLICIES,
				Operation:  "Upsert",
				StartTime:  nrCtx.StartSegmentNow(),
			}
			defer nr.End()
		}
		return r.dbc.QueryRowxContext(ctx, query, params...).Scan(&policyID)
	}); err != nil {
		err = checkPostgresError(err)
		switch {
		case errors.Is(err, errForeignKeyViolation):
			return "", fmt.Errorf("%w: %s", policy.ErrInvalidDetail, err)
		default:
			return "", fmt.Errorf("%w: %s", dbErr, err)
		}
	}

	return policyID, nil
}

func (r PolicyRepository) Update(ctx context.Context, toUpdate *policy.Policy) (string, error) {
	if toUpdate == nil {
		return "", policy.ErrInvalidDetail
	}

	if strings.TrimSpace(toUpdate.ID) == "" {
		return "", policy.ErrInvalidID
	}

	if strings.TrimSpace(toUpdate.ActionID) == "" {
		return "", policy.ErrInvalidDetail
	}

	ctx = otelsql.WithCustomAttributes(
		ctx,
		[]attribute.KeyValue{
			attribute.String("db.repository.method", "Update"),
			attribute.String(string(semconv.DBSQLTableKey), TABLE_POLICIES),
		}...,
	)

	query, params, err := dialect.Update(TABLE_POLICIES).Set(
		goqu.Record{
			"namespace_id": toUpdate.NamespaceID,
			"role_id":      toUpdate.RoleID,
			"action_id":    sql.NullString{String: toUpdate.ActionID, Valid: toUpdate.ActionID != ""},
			"updated_at":   goqu.L("now()"),
		}).Where(goqu.Ex{
		"id": toUpdate.ID,
	}).Returning("id").ToSQL()
	if err != nil {
		return "", fmt.Errorf("%w: %s", queryErr, err)
	}

	var policyID string
	if err = r.dbc.WithTimeout(ctx, func(ctx context.Context) error {
		nrCtx := newrelic.FromContext(ctx)
		if nrCtx != nil {
			nr := newrelic.DatastoreSegment{
				Product:    newrelic.DatastorePostgres,
				Collection: TABLE_POLICIES,
				Operation:  "Update",
				StartTime:  nrCtx.StartSegmentNow(),
			}
			defer nr.End()
		}
		return r.dbc.QueryRowxContext(ctx, query, params...).Scan(&policyID)
	}); err != nil {
		err = checkPostgresError(err)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return "", policy.ErrNotExist
		case errors.Is(err, errDuplicateKey):
			return "", policy.ErrConflict
		case errors.Is(err, errInvalidTexRepresentation):
			return "", policy.ErrInvalidUUID
		case errors.Is(err, errForeignKeyViolation):
			return "", namespace.ErrNotExist
		default:
			return "", err
		}
	}

	return policyID, nil
}
