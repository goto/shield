package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/goto/shield/core/rule"
	"github.com/goto/shield/pkg/db"
	jsoniter "github.com/json-iterator/go"
	newrelic "github.com/newrelic/go-agent/v3/newrelic"
	"go.nhat.io/otelsql"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type RuleRepository struct {
	dbc    *db.Client
	cached []rule.Ruleset
}

func NewRuleRepository(dbc *db.Client) *RuleRepository {
	return &RuleRepository{
		dbc: dbc,
	}
}

func (r *RuleRepository) Upsert(ctx context.Context, name string, config rule.Ruleset) (rule.Config, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	configJson, err := json.Marshal(config)
	if err != nil {
		return rule.Config{}, rule.ErrMarshal
	}

	ctx = r.WithTransaction(ctx)

	query, params, err := goqu.Insert(TABLE_RULE_CONFIGS).Rows(
		goqu.Record{"name": name, "config": configJson},
	).OnConflict(
		goqu.DoUpdate("name", goqu.Record{"name": name, "config": configJson, "updated_at": goqu.L("now()")})).Returning(&RuleConfig{}).ToSQL()
	if err != nil {
		return rule.Config{}, fmt.Errorf("%w: %s", queryErr, err)
	}

	ctx = otelsql.WithCustomAttributes(
		ctx,
		[]attribute.KeyValue{
			attribute.String("db.repository.method", "Upsert"),
			attribute.String(string(semconv.DBSQLTableKey), TABLE_RULE_CONFIGS),
		}...,
	)

	var ruleConfigModel RuleConfig
	if err = r.dbc.WithTimeout(ctx, func(ctx context.Context) error {
		nrCtx := newrelic.FromContext(ctx)
		if nrCtx != nil {
			nr := newrelic.DatastoreSegment{
				Product:    newrelic.DatastorePostgres,
				Collection: TABLE_RULE_CONFIGS,
				Operation:  "Upsert",
				StartTime:  nrCtx.StartSegmentNow(),
			}
			defer nr.End()
		}
		return r.dbc.QueryRowxContext(ctx, query, params...).StructScan(&ruleConfigModel)
	}); err != nil {
		if txErr := r.Rollback(ctx, err); txErr != nil {
			return rule.Config{}, err
		}
		err = checkPostgresError(err)
		return rule.Config{}, err
	}

	err = r.InitCache(ctx)
	if err != nil {
		if txErr := r.Rollback(ctx, err); txErr != nil {
			return rule.Config{}, err
		}
		return rule.Config{}, err
	}

	err = r.Commit(ctx)
	if err != nil {
		return rule.Config{}, err
	}

	return ruleConfigModel.transformToRuleConfig(), nil
}

func (r *RuleRepository) InitCache(ctx context.Context) error {
	query, params, err := dialect.From(TABLE_RULE_CONFIGS).ToSQL()
	if err != nil {
		return err
	}
	ctx = otelsql.WithCustomAttributes(
		ctx,
		[]attribute.KeyValue{
			attribute.String("db.repository.method", "List"),
			attribute.String(string(semconv.DBSQLTableKey), TABLE_RULE_CONFIGS),
		}...,
	)

	var ruleConfigModel []RuleConfig
	if err = r.dbc.WithTimeout(ctx, func(ctx context.Context) error {
		nrCtx := newrelic.FromContext(ctx)
		if nrCtx != nil {
			nr := newrelic.DatastoreSegment{
				Product:    newrelic.DatastorePostgres,
				Collection: TABLE_RULE_CONFIGS,
				Operation:  "List",
				StartTime:  nrCtx.StartSegmentNow(),
			}
			defer nr.End()
		}

		return r.dbc.SelectContext(ctx, &ruleConfigModel, query, params...)
	}); err != nil {
		err = checkPostgresError(err)
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
	}

	newCache := []rule.Ruleset{}
	for _, ruleConfig := range ruleConfigModel {
		rc := ruleConfig.transformToRuleConfig()
		var targetRuleset rule.Ruleset
		if err := json.Unmarshal([]byte(rc.Config), &targetRuleset); err != nil {
			return err
		}
		newCache = append(newCache, targetRuleset)
	}

	r.cached = newCache
	return nil
}

func (r *RuleRepository) IsUpdated(ctx context.Context, since time.Time) bool {
	query, params, err := dialect.From(TABLE_RULE_CONFIGS).Select(goqu.C("updated_at").Gt(since)).Order(goqu.C("updated_at").Desc()).Limit(1).ToSQL()
	if err != nil {
		return false
	}

	ctx = otelsql.WithCustomAttributes(
		ctx,
		[]attribute.KeyValue{
			attribute.String("db.repository.method", "List"),
			attribute.String(string(semconv.DBSQLTableKey), TABLE_RULE_CONFIGS),
		}...,
	)

	var isUpdated bool
	if err = r.dbc.WithTimeout(ctx, func(ctx context.Context) error {
		nrCtx := newrelic.FromContext(ctx)
		if nrCtx != nil {
			nr := newrelic.DatastoreSegment{
				Product:    newrelic.DatastorePostgres,
				Collection: TABLE_RULE_CONFIGS,
				Operation:  "List",
				StartTime:  nrCtx.StartSegmentNow(),
			}
			defer nr.End()
		}

		return r.dbc.GetContext(ctx, &isUpdated, query, params...)
	}); err != nil {
		err = checkPostgresError(err)
		if !errors.Is(err, sql.ErrNoRows) {
			return false
		}
	}

	isUpdatedTest := isUpdated
	return isUpdatedTest
}

func (r *RuleRepository) GetAll(ctx context.Context) ([]rule.Ruleset, error) {
	return r.cached, nil
}

func (r *RuleRepository) Fetch(ctx context.Context) ([]rule.Ruleset, error) {
	query, params, err := dialect.From(TABLE_RULE_CONFIGS).ToSQL()
	if err != nil {
		return []rule.Ruleset{}, err
	}
	ctx = otelsql.WithCustomAttributes(
		ctx,
		[]attribute.KeyValue{
			attribute.String("db.repository.method", "List"),
			attribute.String(string(semconv.DBSQLTableKey), TABLE_RULE_CONFIGS),
		}...,
	)

	var ruleConfigModel []RuleConfig
	if err = r.dbc.WithTimeout(ctx, func(ctx context.Context) error {
		nrCtx := newrelic.FromContext(ctx)
		if nrCtx != nil {
			nr := newrelic.DatastoreSegment{
				Product:    newrelic.DatastorePostgres,
				Collection: TABLE_RULE_CONFIGS,
				Operation:  "List",
				StartTime:  nrCtx.StartSegmentNow(),
			}
			defer nr.End()
		}

		return r.dbc.SelectContext(ctx, &ruleConfigModel, query, params...)
	}); err != nil {
		err = checkPostgresError(err)
		if !errors.Is(err, sql.ErrNoRows) {
			return []rule.Ruleset{}, err
		}
	}

	rules := []rule.Ruleset{}
	for _, ruleConfig := range ruleConfigModel {
		rc := ruleConfig.transformToRuleConfig()
		var targetRuleset rule.Ruleset
		if err := json.Unmarshal([]byte(rc.Config), &targetRuleset); err != nil {
			return []rule.Ruleset{}, err
		}
		rules = append(rules, targetRuleset)
	}

	return rules, nil
}

func (r *RuleRepository) WithTransaction(ctx context.Context) context.Context {
	return r.dbc.WithTransaction(ctx, sql.TxOptions{})
}

func (r *RuleRepository) Rollback(ctx context.Context, err error) error {
	if txErr := r.dbc.Rollback(ctx); txErr != nil {
		return fmt.Errorf("rollback error %s with error: %w", txErr.Error(), err)
	}
	return nil
}

func (r *RuleRepository) Commit(ctx context.Context) error {
	return r.dbc.Commit(ctx)
}
