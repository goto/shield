package postgres_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/goto/salt/log"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/suite"

	"github.com/goto/shield/core/namespace"
	"github.com/goto/shield/core/policy"
	"github.com/goto/shield/internal/store/postgres"
	"github.com/goto/shield/pkg/db"
)

type PolicyRepositoryTestSuite struct {
	suite.Suite
	ctx        context.Context
	client     *db.Client
	pool       *dockertest.Pool
	resource   *dockertest.Resource
	repository *postgres.PolicyRepository
	policyIDs  []string
}

func (s *PolicyRepositoryTestSuite) SetupSuite() {
	var err error

	logger := log.NewZap()
	s.client, s.pool, s.resource, err = newTestClient(logger)
	if err != nil {
		s.T().Fatal(err)
	}

	s.ctx = context.TODO()
	s.repository = postgres.NewPolicyRepository(s.client)

	_, err = bootstrapNamespace(s.client)
	if err != nil {
		s.T().Fatal(err)
	}

	_, err = bootstrapAction(s.client)
	if err != nil {
		s.T().Fatal(err)
	}

	_, err = bootstrapRole(s.client)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *PolicyRepositoryTestSuite) SetupTest() {
	var err error
	s.policyIDs, err = bootstrapPolicy(s.client)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *PolicyRepositoryTestSuite) TearDownSuite() {
	// Clean tests
	if err := purgeDocker(s.pool, s.resource); err != nil {
		s.T().Fatal(err)
	}
}

func (s *PolicyRepositoryTestSuite) TearDownTest() {
	if err := s.cleanup(); err != nil {
		s.T().Fatal(err)
	}
}

func (s *PolicyRepositoryTestSuite) cleanup() error {
	queries := []string{
		fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", postgres.TABLE_POLICIES),
	}
	return execQueries(context.TODO(), s.client, queries)
}

func (s *PolicyRepositoryTestSuite) TestGet() {
	type testCase struct {
		Description    string
		SelectedID     string
		ExpectedPolicy policy.Policy
		ErrString      string
	}

	testCases := []testCase{
		{
			Description: "should get a policy",
			SelectedID:  s.policyIDs[0],
			ExpectedPolicy: policy.Policy{
				RoleID:      "ns1:role1",
				NamespaceID: "ns1",
				ActionID:    "action1",
			},
		},
		{
			Description: "should return error invalid id if empty",
			ErrString:   policy.ErrInvalidID.Error(),
		},
		{
			Description: "should return error no exist if can't found policy",
			SelectedID:  uuid.NewString(),
			ErrString:   policy.ErrNotExist.Error(),
		},
		{
			Description: "should return error no exist if id is not uuid",
			SelectedID:  "some-id",
			ErrString:   policy.ErrInvalidUUID.Error(),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			got, err := s.repository.Get(s.ctx, tc.SelectedID)
			if tc.ErrString != "" {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
			if !cmp.Equal(got, tc.ExpectedPolicy, cmpopts.IgnoreFields(policy.Policy{},
				"ID")) {
				s.T().Fatalf("got result %+v, expected was %+v", got, tc.ExpectedPolicy)
			}
		})
	}
}

func (s *PolicyRepositoryTestSuite) TestUpsert() {
	type testCase struct {
		Description string
		Policy      *policy.Policy
		Err         error
	}

	testCases := []testCase{
		{
			Description: "should create a policy",
			Policy: &policy.Policy{
				RoleID:      "ns1:role2",
				ActionID:    "action4",
				NamespaceID: "ns1",
			},
		},
		{
			Description: "should return error if role id does not exist",
			Policy: &policy.Policy{
				RoleID:      "role2-random",
				ActionID:    "action4",
				NamespaceID: "ns1",
			},
			Err: policy.ErrInvalidDetail,
		},
		{
			Description: "should return error if policy id does not exist",
			Policy: &policy.Policy{
				RoleID:      "role2",
				ActionID:    "action4-random",
				NamespaceID: "ns1",
			},
			Err: policy.ErrInvalidDetail,
		},
		{
			Description: "should return error if namespace id does not exist",
			Policy: &policy.Policy{
				RoleID:      "role2",
				ActionID:    "action4",
				NamespaceID: "ns1-random",
			},
			Err: policy.ErrInvalidDetail,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			got, err := s.repository.Upsert(s.ctx, tc.Policy)
			if tc.Err != nil {
				if errors.Is(tc.Err, err) {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.Err.Error())
				}
			} else {
				if len(got) != len(uuid.NewString()) {
					s.T().Fatalf("got result %s, expected was a uuid", got)
				}
			}
		})
	}
}

func (s *PolicyRepositoryTestSuite) TestList() {
	type testCase struct {
		Description     string
		ExpectedPolicys []policy.Policy
		Filter          policy.Filters
		ErrString       string
	}

	testCases := []testCase{
		{
			Description: "should get all policys",
			ExpectedPolicys: []policy.Policy{
				{
					RoleID:      "ns1:role1",
					NamespaceID: "ns1",
					ActionID:    "action1",
				},
				{
					RoleID:      "ns2:role2",
					NamespaceID: "ns2",
					ActionID:    "action2",
				},
				{
					RoleID:      "ns1:role2",
					NamespaceID: "ns1",
					ActionID:    "action3",
				},
			},
		},
		{
			Description: "should get all policys with filter",
			Filter:      policy.Filters{NamespaceID: "ns2"},
			ExpectedPolicys: []policy.Policy{
				{
					RoleID:      "ns2:role2",
					NamespaceID: "ns2",
					ActionID:    "action2",
				},
			},
		},
		{
			Description: "should get none policys with filter",
			Filter:      policy.Filters{NamespaceID: "ns3"},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			got, err := s.repository.List(s.ctx, tc.Filter)
			if tc.ErrString != "" {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
			// TODO figure out how to compare metadata map[string]any
			if !cmp.Equal(got, tc.ExpectedPolicys, cmpopts.IgnoreFields(policy.Policy{},
				"ID")) {
				s.T().Fatalf("got result %+v, expected was %+v", got, tc.ExpectedPolicys)
			}
		})
	}
}

func (s *PolicyRepositoryTestSuite) TestUpdate() {
	type testCase struct {
		Description      string
		Policy           *policy.Policy
		ExpectedPolicyID string
		ErrString        string
	}

	testCases := []testCase{
		{
			Description: "should update an policy",
			Policy: &policy.Policy{
				ID:          s.policyIDs[0],
				RoleID:      "ns1:role1",
				ActionID:    "action4",
				NamespaceID: "ns1",
			},
			ExpectedPolicyID: s.policyIDs[0],
		},
		{
			Description: "should return error if namespace id does not exist",
			Policy: &policy.Policy{
				ID:          s.policyIDs[1],
				RoleID:      "ns2:role2",
				ActionID:    "action2",
				NamespaceID: "random-ns2",
			},
			ErrString:        namespace.ErrNotExist.Error(),
			ExpectedPolicyID: "",
		},
		{
			Description: "should return error if policy not found",
			Policy: &policy.Policy{
				ID:          uuid.NewString(),
				RoleID:      "ns1:role2",
				NamespaceID: "ns1",
				ActionID:    "action3",
			},
			ErrString:        policy.ErrNotExist.Error(),
			ExpectedPolicyID: "",
		},
		{
			Description:      "should return error if policy id is empty",
			ErrString:        policy.ErrInvalidDetail.Error(),
			ExpectedPolicyID: "",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			got, err := s.repository.Update(s.ctx, tc.Policy)
			if tc.ErrString != "" {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
			if !cmp.Equal(got, tc.ExpectedPolicyID) {
				s.T().Fatalf("got result %+v, expected was %+v", got, tc.ExpectedPolicyID)
			}
		})
	}
}

func TestPolicyRepository(t *testing.T) {
	suite.Run(t, new(PolicyRepositoryTestSuite))
}
