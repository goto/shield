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
	"github.com/goto/shield/core/user"
	"github.com/goto/shield/internal/schema"
	"github.com/goto/shield/internal/store/postgres"
	"github.com/goto/shield/pkg/db"
	"github.com/goto/shield/pkg/uuid"
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
	data       []servicedata.ServiceData
	users      []user.User
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

	s.users, err = bootstrapUser(s.client)
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

	s.resources, err = bootstrapResource(s.client, s.projects, organizations, namespaces, s.users)
	if err != nil {
		s.T().Fatal(err)
	}

	s.keys, err = bootstrapServiceDataKey(s.client, s.resources, s.projects)
	if err != nil {
		s.T().Fatal(err)
	}

	s.data, err = bootstrapServiceData(s.client, s.users, s.keys)
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
		fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", postgres.TABLE_SERVICE_DATA),
		fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", postgres.TABLE_USERS),
		fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", postgres.TABLE_ORGANIZATIONS),
		fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", postgres.TABLE_PROJECTS),
		fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", postgres.TABLE_RESOURCES),
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

	testCases := []testCase{
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

func (s *ServiceDataRepositoryTestSuite) TestUpsert() {
	type testCase struct {
		Description         string
		ServiceDataToCreate servicedata.ServiceData
		ExpectedServiceData servicedata.ServiceData
		ErrString           string
	}

	testNamespaceID := uuid.NewString()
	testEntityID := uuid.NewString()
	testValue := "test-value"

	testCases := []testCase{
		{
			Description: "should create a service data",
			ServiceDataToCreate: servicedata.ServiceData{
				NamespaceID: testNamespaceID,
				EntityID:    testEntityID,
				Key:         s.keys[0],
				Value:       testValue,
			},
			ExpectedServiceData: servicedata.ServiceData{
				Key: servicedata.Key{
					Key: s.keys[0].Key,
				},
				Value: testValue,
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			got, err := s.repository.Upsert(s.ctx, tc.ServiceDataToCreate)
			if tc.ErrString != "" {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
			if !cmp.Equal(got, tc.ExpectedServiceData, cmpopts.IgnoreFields(servicedata.ServiceData{},
				"ID",
			)) {
				s.T().Fatalf("got result %+v, expected was %+v", got, tc.ExpectedServiceData)
			}
		})
	}
}

func (s *ServiceDataRepositoryTestSuite) TestGetKeyByURN() {
	type testCase struct {
		Description string
		URN         string
		ExpectedKey servicedata.Key
		ErrString   string
	}

	testCases := []testCase{
		{
			Description: "should create a key",
			URN:         s.keys[0].URN,
			ExpectedKey: s.keys[0],
		},
		{
			Description: "should return not exist error if key not found",
			URN:         "invalid-urn",
			ErrString:   servicedata.ErrNotExist.Error(),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			got, err := s.repository.GetKeyByURN(s.ctx, tc.URN)
			if tc.ErrString != "" {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
			if !cmp.Equal(got, tc.ExpectedKey, cmpopts.IgnoreFields(servicedata.ServiceData{},
				"ID",
			)) {
				s.T().Fatalf("got result %+v, expected was %+v", got, tc.ExpectedKey)
			}
		})
	}
}

func (s *ServiceDataRepositoryTestSuite) TestGet() {
	type testCase struct {
		Description  string
		filter       servicedata.Filter
		ExpectedData []servicedata.ServiceData
		ErrString    string
	}

	expected := servicedata.ServiceData{
		NamespaceID: schema.UserPrincipal,
		EntityID:    s.users[0].ID,
		Key: servicedata.Key{
			URN:        s.keys[0].URN,
			ProjectID:  s.keys[0].ProjectID,
			Key:        s.keys[0].Key,
			ResourceID: s.keys[0].ResourceID,
		},
		Value: s.data[0].Value,
	}

	testCases := []testCase{
		{
			Description: "should get a service data",
			filter: servicedata.Filter{
				EntityIDs: [][]string{{schema.UserPrincipal, s.users[0].ID}},
			},
			ExpectedData: []servicedata.ServiceData{expected},
		},
		{
			Description: "should get none service data if no data queried",
			filter: servicedata.Filter{
				EntityIDs: [][]string{{schema.UserPrincipal, s.users[0].ID}},
				Project:   s.projects[1].ID,
			},
		},
		{
			Description: "should return none service data if filter entity ids is empty",
			filter: servicedata.Filter{
				EntityIDs: [][]string{},
				Project:   s.projects[1].ID,
			},
			ExpectedData: []servicedata.ServiceData{},
		},
		{
			Description: "should get err invalid detail",
			filter: servicedata.Filter{
				EntityIDs: [][]string{{schema.UserPrincipal, "invalid-entity-id"}},
				Project:   "invalid-project-uuid",
			},
			ErrString:    servicedata.ErrInvalidDetail.Error(),
			ExpectedData: []servicedata.ServiceData{},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			got, err := s.repository.Get(s.ctx, tc.filter)
			if tc.ErrString != "" {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
			if !cmp.Equal(got, tc.ExpectedData, cmpopts.IgnoreFields(servicedata.Key{},
				"ID",
			)) {
				s.T().Fatalf("got result %+v, expected was %+v", got, tc.ExpectedData)
			}
		})
	}
}

func TestServiceDataRepository(t *testing.T) {
	suite.Run(t, new(ServiceDataRepositoryTestSuite))
}
