package postgres_test

import (
	"context"
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
				Config: "{\"Rules\": [{\"Hooks\": null, \"Backend\": {\"URL\": \"\", \"Prefix\": \"\", \"Namespace\": \"\"}, \"Frontend\": {\"URL\": \"\", \"URLRx\": null, \"Method\": \"\"}, \"Middlewares\": null}]}",
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
				Config: "{\"Rules\": [{\"Hooks\": null, \"Backend\": {\"URL\": \"\", \"Prefix\": \"\", \"Namespace\": \"\"}, \"Frontend\": {\"URL\": \"\", \"URLRx\": null, \"Method\": \"\"}, \"Middlewares\": null}]}",
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
