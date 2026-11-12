package postgres_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/goto/salt/log"
	"github.com/goto/shield/core/rule"
	"github.com/goto/shield/internal/store/postgres"
	"github.com/goto/shield/pkg/db"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/suite"
)

type RuleRepositoryTestSuite struct {
	suite.Suite
	ctx        context.Context
	client     *db.Client
	pool       *dockertest.Pool
	resource   *dockertest.Resource
	repository *postgres.RuleRepository
	Config     []rule.Config
}

func (s *RuleRepositoryTestSuite) SetupSuite() {
	var err error

	logger := log.NewZap()
	s.client, s.pool, s.resource, err = newTestClient(logger)
	if err != nil {
		s.T().Fatal(err)
	}

	s.ctx = context.TODO()
	s.repository = postgres.NewRuleRepository(s.client)

	s.Config, err = bootstrapRuleConfig(s.client)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *RuleRepositoryTestSuite) TestUpsert() {
	type testCase struct {
		Description string
		Name        string
		Config      rule.Ruleset
		Expected    rule.Config
		ErrString   string
	}

	testCases := []testCase{
		{
			Description: "should create a resource config",
			Name:        "test",
			Config: rule.Ruleset{
				Rules: []rule.Rule{{}},
			},
			Expected: rule.Config{
				ID:     2,
				Name:   "test",
				Config: "{\"Rules\": [{\"Hooks\": null, \"Backend\": {\"URL\": \"\", \"Prefix\": \"\", \"Namespace\": \"\"}, \"Frontend\": {\"URL\": \"\", \"URLRx\": null, \"Method\": \"\"}, \"Middlewares\": null, \"SkipReadBody\": false}]}",
			},
		},
		{
			Description: "should update a resource config",
			Name:        s.Config[0].Name,
			Config: rule.Ruleset{
				Rules: []rule.Rule{{}},
			},
			Expected: rule.Config{
				ID:     s.Config[0].ID,
				Name:   s.Config[0].Name,
				Config: "{\"Rules\": [{\"Hooks\": null, \"Backend\": {\"URL\": \"\", \"Prefix\": \"\", \"Namespace\": \"\"}, \"Frontend\": {\"URL\": \"\", \"URLRx\": null, \"Method\": \"\"}, \"Middlewares\": null, \"SkipReadBody\": false}]}",
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			got, err := s.repository.Upsert(s.ctx, tc.Name, tc.Config)
			if tc.ErrString != "" {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
			if !cmp.Equal(got, tc.Expected, cmpopts.IgnoreFields(rule.Config{},
				"CreatedAt",
				"UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", got, tc.Expected)
			}
		})
	}
}

func TestRuleRepository(t *testing.T) {
	suite.Run(t, new(RuleRepositoryTestSuite))
}

func TestWeakBasicAuthRejectedBeforeDatabaseWrite(t *testing.T) {
	repo := postgres.NewRuleRepository(nil)
	_, err := repo.Upsert(context.Background(), "synthetic", rule.Ruleset{Rules: []rule.Rule{{Middlewares: rule.MiddlewareSpecs{{Name: "basic_auth", Config: map[string]interface{}{"users": []map[string]interface{}{{"user": "synthetic", "password": "{SHA}synthetic"}}}}}}}})
	if !errors.Is(err, rule.ErrInvalidRuleConfig) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func (s *RuleRepositoryTestSuite) TestInvalidReloadRetainsPreviousCache() {
	before, err := s.repository.GetAll(s.ctx)
	s.Require().NoError(err)
	invalid := rule.Ruleset{Rules: []rule.Rule{{Middlewares: rule.MiddlewareSpecs{{Name: "basic_auth", Config: map[string]interface{}{"users": []map[string]interface{}{{"user": "synthetic", "password": "$apr1$synthetic"}}}}}}}}
	payload, err := json.Marshal(invalid)
	s.Require().NoError(err)
	_, err = s.client.ExecContext(s.ctx, "INSERT INTO "+postgres.TABLE_RULE_CONFIGS+" (id, name, config) VALUES ($1, $2, $3)", 1000000, "invalid-security-fixture", string(payload))
	s.Require().NoError(err)
	defer func() {
		_, err := s.client.ExecContext(s.ctx, "DELETE FROM "+postgres.TABLE_RULE_CONFIGS+" WHERE name=$1", "invalid-security-fixture")
		s.Require().NoError(err)
	}()
	s.ErrorIs(s.repository.InitCache(s.ctx), rule.ErrInvalidRuleConfig)
	after, err := s.repository.GetAll(s.ctx)
	s.Require().NoError(err)
	s.Equal(before, after)
}

func TestEmptyRulesetRetiresPersistedRoute(t *testing.T) {
	client, pool, container, err := newTestClient(log.NewNoop())
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := purgeDocker(pool, container); err != nil {
			t.Error(err)
		}
	}()
	configs, err := bootstrapRuleConfig(client)
	if err != nil {
		t.Fatal(err)
	}
	if len(configs) != 1 {
		t.Fatal("expected one synthetic named configuration")
	}
	repo := postgres.NewRuleRepository(client)
	if err := repo.InitCache(context.Background()); err != nil {
		t.Fatal(err)
	}
	before, err := repo.GetAll(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(before) != 1 || len(before[0].Rules) == 0 {
		t.Fatal("expected a registered route")
	}
	_, err = repo.Upsert(context.Background(), configs[0].Name, rule.Ruleset{})
	if err != nil {
		t.Fatal(err)
	}
	after, err := repo.GetAll(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(after) != 1 || len(after[0].Rules) != 0 {
		t.Fatal("retired frontend still cached")
	}
}
