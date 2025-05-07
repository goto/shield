package postgres_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

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
	config     []rule.Config
}

func (s *RuleRepositoryTestSuite) SetupSuite() {
	var err error

	logger := log.NewZap()
	s.client, s.pool, s.resource, err = newTestClient(logger)
	if err != nil {
		s.T().Fatal(err)
	}

	s.ctx = context.TODO()

	s.config, err = bootstrapRuleConfig(s.client)
	if err != nil {
		s.T().Fatal(err)
	}

	s.repository = postgres.NewRuleRepository(s.client)
	err = s.repository.InitCache(s.ctx)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *RuleRepositoryTestSuite) mergeRules() ([]rule.Ruleset, error) {
	rules := []rule.Ruleset{}
	for _, ruleConfig := range s.config {
		var targetRuleset rule.Ruleset
		if err := json.Unmarshal([]byte(ruleConfig.Config), &targetRuleset); err != nil {
			return []rule.Ruleset{}, err
		}
		rules = append(rules, targetRuleset)
	}

	return rules, nil
}

func (s *RuleRepositoryTestSuite) TestGetAll() {
	expected, err := s.mergeRules()
	if err != nil {
		s.T().Fatal(err)
	}

	type testCase struct {
		Description string
		Expected    []rule.Ruleset
		ErrString   string
	}

	testCases := []testCase{
		{
			Description: "should get all rules from repository cache",
			Expected:    expected,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			got, err := s.repository.GetAll(s.ctx)
			if tc.ErrString != "" {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
			if !cmp.Equal(got, tc.Expected) {
				s.T().Fatalf("got result %+v, expected was %+v", got, tc.Expected)
			}
		})
	}
}

func (s *RuleRepositoryTestSuite) TestFetch() {
	expected, err := s.mergeRules()
	if err != nil {
		s.T().Fatal(err)
	}

	type testCase struct {
		Description string
		Expected    []rule.Ruleset
		ErrString   string
	}

	testCases := []testCase{
		{
			Description: "should get all rules from repository cache",
			Expected:    expected,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			got, err := s.repository.Fetch(s.ctx)
			if tc.ErrString != "" {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
			if !cmp.Equal(got, tc.Expected) {
				s.T().Fatalf("got result %+v, expected was %+v", got, tc.Expected)
			}
		})
	}
}

func (s *RuleRepositoryTestSuite) TestIsUpdated() {
	type testCase struct {
		Description string
		Since       time.Time
		Expected    bool
	}

	testCases := []testCase{
		{
			Description: "should get true if since before last updated",
			Since:       time.Time{},
			Expected:    true,
		},
		{
			Description: "should get false if since after last updated",
			Since:       s.config[0].UpdatedAt.Add(10 * time.Hour),
			Expected:    false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			got := s.repository.IsUpdated(s.ctx, tc.Since)
			if !cmp.Equal(got, tc.Expected) {
				s.T().Fatalf("got result %+v, expected was %+v", got, tc.Expected)
			}
		})
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
			Name:        s.config[0].Name,
			Config: rule.Ruleset{
				Rules: []rule.Rule{{}},
			},
			Expected: rule.Config{
				ID:     s.config[0].ID,
				Name:   s.config[0].Name,
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
