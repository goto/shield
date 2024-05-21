package postgres_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/goto/salt/log"
	"github.com/goto/shield/core/project"
	"github.com/goto/shield/core/resource"
	"github.com/goto/shield/core/servicedata"
	"github.com/goto/shield/internal/store/postgres"
	"github.com/goto/shield/pkg/db"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/suite"
)

type ServiceDataRepositoryTestSuite struct {
	suite.Suite
	ctx        context.Context
	client     *db.Client
	pool       *dockertest.Pool
	resource   *dockertest.Resource
	repository *postgres.ServiceDataRepository
	keys       []servicedata.Key
	projects   []project.Project
	resources  []resource.Resource
}

func (s *ServiceDataRepositoryTestSuite) SetupSuite() {
	var err error

	logger := log.NewZap()
	s.client, s.pool, s.resource, err = newTestClient(logger)
	if err != nil {
		s.T().Fatal(err)
	}

	s.ctx = context.TODO()
	s.repository = postgres.NewServiceDataRepository(s.client)
}

func (s *ServiceDataRepositoryTestSuite) SetupTest() {
	var err error

	namespaces, err := bootstrapNamespace(s.client)
	if err != nil {
		s.T().Fatal(err)
	}

	_, err = bootstrapMetadataKeys(s.client)
	if err != nil {
		s.T().Fatal(err)
	}

	users, err := bootstrapUser(s.client)
	if err != nil {
		s.T().Fatal(err)
	}

	organizations, err := bootstrapOrganization(s.client)
	if err != nil {
		s.T().Fatal(err)
	}

	s.projects, err = bootstrapProject(s.client, organizations)
	if err != nil {
		s.T().Fatal(err)
	}

	s.resources, err = bootstrapResource(s.client, s.projects, organizations, namespaces, users)
	if err != nil {
		s.T().Fatal(err)
	}

	s.keys, err = bootstrapServiceDataKey(s.client, s.resources, s.projects)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *ServiceDataRepositoryTestSuite) TearDownSuite() {
	// Clean tests
	if err := purgeDocker(s.pool, s.resource); err != nil {
		s.T().Fatal(err)
	}
}

func (s *ServiceDataRepositoryTestSuite) TearDownTest() {
	if err := s.cleanup(); err != nil {
		s.T().Fatal(err)
	}
}

func (s *ServiceDataRepositoryTestSuite) cleanup() error {
	queries := []string{
		fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", postgres.TABLE_SERVICE_DATA_KEYS),
	}
	return execQueries(context.TODO(), s.client, queries)
}

func (s *ServiceDataRepositoryTestSuite) TestCreateKey() {
	type testCase struct {
		Description string
		KeyToCreate servicedata.Key
		ExpectedKey servicedata.Key
		ErrString   string
	}

	var testCases = []testCase{
		{
			Description: "should create a key",
			KeyToCreate: servicedata.Key{
				URN:         "test-urn",
				ProjectID:   s.projects[0].ID,
				Key:         "test-key",
				Description: "description for test-key",
				ResourceID:  s.resources[0].Idxa,
			},
			ExpectedKey: servicedata.Key{
				URN:         "test-urn",
				ProjectID:   s.projects[0].ID,
				Key:         "test-key",
				Description: "description for test-key",
				ResourceID:  s.resources[0].Idxa,
			},
		},
		{
			Description: "should return conflict error if key urn already exist",
			KeyToCreate: servicedata.Key{
				URN:         s.keys[0].URN,
				ProjectID:   s.projects[0].ID,
				Key:         s.keys[0].Key,
				Description: s.keys[0].Key,
				ResourceID:  s.resources[0].Idxa,
			},
			ErrString: servicedata.ErrConflict.Error(),
		},
		{
			Description: "should return invalid detail if project id not exist",
			KeyToCreate: servicedata.Key{
				URN:         "test-urn-00",
				ProjectID:   "00000000-0000-0000-0000-000000000000",
				Key:         "test-key",
				Description: "description for test-key",
				ResourceID:  s.resources[0].Idxa,
			},
			ErrString: servicedata.ErrInvalidDetail.Error(),
		},
		{
			Description: "should return invalid detail if resource id not exist",
			KeyToCreate: servicedata.Key{
				URN:         "test-urn-00",
				ProjectID:   s.projects[0].ID,
				Key:         "test-key",
				Description: "description for test-key",
				ResourceID:  "00000000-0000-0000-0000-000000000000",
			},
			ErrString: servicedata.ErrInvalidDetail.Error(),
		},
		{
			Description: "should return invalid detail if key is empty",
			KeyToCreate: servicedata.Key{
				URN:         "test-urn-00",
				ProjectID:   s.projects[0].ID,
				Key:         "",
				Description: "description for test-key",
				ResourceID:  s.resources[0].Idxa,
			},
			ErrString: servicedata.ErrInvalidDetail.Error(),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			got, err := s.repository.CreateKey(s.ctx, tc.KeyToCreate)
			if tc.ErrString != "" {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
			if !cmp.Equal(got, tc.ExpectedKey, cmpopts.IgnoreFields(servicedata.Key{},
				"ID",
			)) {
				s.T().Fatalf("got result %+v, expected was %+v", got, tc.ExpectedKey)
			}
		})
	}
}

func TestServiceDataRepository(t *testing.T) {
	suite.Run(t, new(ServiceDataRepositoryTestSuite))
}
