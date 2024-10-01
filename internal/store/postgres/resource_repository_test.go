package postgres_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/goto/salt/log"
	"github.com/goto/shield/core/namespace"
	"github.com/goto/shield/core/organization"
	"github.com/goto/shield/core/project"
	"github.com/goto/shield/core/resource"
	"github.com/goto/shield/core/user"
	"github.com/goto/shield/internal/schema"
	"github.com/goto/shield/internal/store/postgres"
	"github.com/goto/shield/pkg/db"
	"github.com/goto/shield/pkg/uuid"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/suite"
)

type ResourceRepositoryTestSuite struct {
	suite.Suite
	ctx            context.Context
	client         *db.Client
	pool           *dockertest.Pool
	resource       *dockertest.Resource
	repository     *postgres.ResourceRepository
	resources      []resource.Resource
	projects       []project.Project
	orgs           []organization.Organization
	namespaces     []namespace.Namespace
	users          []user.User
	resourceConfig []resource.ResourceConfig
}

func (s *ResourceRepositoryTestSuite) SetupSuite() {
	var err error

	logger := log.NewZap()
	s.client, s.pool, s.resource, err = newTestClient(logger)
	if err != nil {
		s.T().Fatal(err)
	}

	s.ctx = context.TODO()
	s.repository = postgres.NewResourceRepository(s.client)

	s.namespaces, err = bootstrapNamespace(s.client)
	if err != nil {
		s.T().Fatal(err)
	}

	s.orgs, err = bootstrapOrganization(s.client)
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

	s.projects, err = bootstrapProject(s.client, s.orgs)
	if err != nil {
		s.T().Fatal(err)
	}

	s.resourceConfig, err = bootstrapResourceConfig(s.client)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *ResourceRepositoryTestSuite) SetupTest() {
	var err error
	s.resources, err = bootstrapResource(s.client, s.projects, s.orgs, s.namespaces, s.users)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *ResourceRepositoryTestSuite) TearDownSuite() {
	// Clean tests
	if err := purgeDocker(s.pool, s.resource); err != nil {
		s.T().Fatal(err)
	}
}

func (s *ResourceRepositoryTestSuite) TearDownTest() {
	if err := s.cleanup(); err != nil {
		s.T().Fatal(err)
	}
}

func (s *ResourceRepositoryTestSuite) cleanup() error {
	queries := []string{
		fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", postgres.TABLE_RESOURCES),
	}
	return execQueries(context.TODO(), s.client, queries)
}

func (s *ResourceRepositoryTestSuite) TestGetByID() {
	type testCase struct {
		Description      string
		SelectedID       string
		ExpectedResource resource.Resource
		ErrString        string
	}

	testCases := []testCase{
		{
			Description: "should get a resource",
			SelectedID:  s.resources[0].Idxa,
			ExpectedResource: resource.Resource{
				Idxa:           s.resources[0].Idxa,
				URN:            s.resources[0].URN,
				Name:           s.resources[0].Name,
				ProjectID:      s.resources[0].ProjectID,
				OrganizationID: s.resources[0].OrganizationID,
				NamespaceID:    s.resources[0].NamespaceID,
				UserID:         s.resources[0].UserID,
			},
		},
		{
			Description: "should return error if id is empty",
			ErrString:   resource.ErrInvalidID.Error(),
		},
		{
			Description: "should return error no exist if can't found resource",
			SelectedID:  uuid.NewString(),
			ErrString:   resource.ErrNotExist.Error(),
		},
		{
			Description: "should return error if id is not uuid",
			SelectedID:  "10000",
			ErrString:   resource.ErrInvalidUUID.Error(),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			got, err := s.repository.GetByID(s.ctx, tc.SelectedID)
			if tc.ErrString != "" {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
			if !cmp.Equal(got, tc.ExpectedResource, cmpopts.IgnoreFields(resource.Resource{},
				"CreatedAt",
				"UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", got, tc.ExpectedResource)
			}
		})
	}
}

func (s *ResourceRepositoryTestSuite) TestGetByURN() {
	type testCase struct {
		Description      string
		SelectedURN      string
		ExpectedResource resource.Resource
		ErrString        string
	}

	testCases := []testCase{
		{
			Description: "should get a resource",
			SelectedURN: s.resources[0].URN,
			ExpectedResource: resource.Resource{
				Idxa:           s.resources[0].Idxa,
				URN:            s.resources[0].URN,
				Name:           s.resources[0].Name,
				ProjectID:      s.resources[0].ProjectID,
				OrganizationID: s.resources[0].OrganizationID,
				NamespaceID:    s.resources[0].NamespaceID,
				UserID:         s.resources[0].UserID,
			},
		},
		{
			Description: "should return error if urn is empty",
			ErrString:   resource.ErrInvalidURN.Error(),
		},
		{
			Description: "should return error no exist if can't found resource",
			SelectedURN: "some-urn",
			ErrString:   resource.ErrNotExist.Error(),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			got, err := s.repository.GetByURN(s.ctx, tc.SelectedURN)
			if tc.ErrString != "" {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
			if !cmp.Equal(got, tc.ExpectedResource, cmpopts.IgnoreFields(resource.Resource{},
				"CreatedAt",
				"UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", got, tc.ExpectedResource)
			}
		})
	}
}

func (s *ResourceRepositoryTestSuite) TestUpsert() {
	type testCase struct {
		Description      string
		ResourceToCreate resource.Resource
		ExpectedResource resource.Resource
		ErrString        string
	}

	testCases := []testCase{
		{
			Description: "should create a resource",
			ResourceToCreate: resource.Resource{
				URN:            "new-urn-4",
				Name:           "resource4",
				ProjectID:      s.resources[0].ProjectID,
				OrganizationID: s.resources[0].OrganizationID,
				NamespaceID:    s.resources[0].NamespaceID,
				UserID:         s.resources[0].UserID,
			},
			ExpectedResource: resource.Resource{
				URN:            "new-urn-4",
				Name:           "resource4",
				ProjectID:      s.resources[0].ProjectID,
				OrganizationID: s.resources[0].OrganizationID,
				NamespaceID:    s.resources[0].NamespaceID,
				UserID:         s.resources[0].UserID,
			},
		},
		{
			Description: "should return error if namespace id does not exist",
			ResourceToCreate: resource.Resource{
				URN:            "new-urn-notexist",
				Name:           "resource4",
				ProjectID:      s.resources[0].ProjectID,
				OrganizationID: s.resources[0].OrganizationID,
				NamespaceID:    "some-ns",
				UserID:         s.resources[0].UserID,
			},
			ErrString: resource.ErrInvalidDetail.Error(),
		},
		{
			Description: "should return error if org id does not exist",
			ResourceToCreate: resource.Resource{
				URN:            "new-urn-notexist",
				Name:           "resource4",
				ProjectID:      s.resources[0].ProjectID,
				OrganizationID: uuid.NewString(),
				NamespaceID:    s.resources[0].NamespaceID,
				UserID:         s.resources[0].UserID,
			},
			ErrString: resource.ErrInvalidDetail.Error(),
		},
		{
			Description: "should return error if org id is not uuid",
			ResourceToCreate: resource.Resource{
				URN:            "new-urn-notexist",
				Name:           "resource4",
				ProjectID:      s.resources[0].ProjectID,
				OrganizationID: "some-str",
				NamespaceID:    s.resources[0].NamespaceID,
				UserID:         s.resources[0].UserID,
			},
			ErrString: resource.ErrInvalidUUID.Error(),
		},
		//{
		//	Description: "should return error if group id does not exist",
		//	ResourceToCreate: resource.Resource{
		//		URN:            "new-urn-notexist",
		//		Name:           "resource4",
		//		ProjectID:      s.resources[0].ProjectID,
		//		OrganizationID: s.resources[0].OrganizationID,
		//		NamespaceID:    s.resources[0].NamespaceID,
		//		UserID:         s.resources[0].UserID,
		//	},
		//	ErrString: resource.ErrInvalidDetail.Error(),
		//},
		//{
		//	Description: "should return error if group id is not uuid",
		//	ResourceToCreate: resource.Resource{
		//		URN:            "new-urn-notexist",
		//		Name:           "resource4",
		//		ProjectID:      s.resources[0].ProjectID,
		//		OrganizationID: s.resources[0].OrganizationID,
		//		NamespaceID:    s.resources[0].NamespaceID,
		//		UserID:         s.resources[0].UserID,
		//	},
		//	ErrString: resource.ErrInvalidUUID.Error(),
		//},
		{
			Description: "should return error if project id does not exist",
			ResourceToCreate: resource.Resource{
				URN:            "new-urn-notexist",
				Name:           "resource4",
				ProjectID:      uuid.NewString(),
				OrganizationID: s.resources[0].OrganizationID,
				NamespaceID:    s.resources[0].NamespaceID,
				UserID:         s.resources[0].UserID,
			},
			ErrString: resource.ErrInvalidDetail.Error(),
		},
		{
			Description: "should return error if project id is not uuid",
			ResourceToCreate: resource.Resource{
				URN:            "new-urn-notexist",
				Name:           "resource4",
				ProjectID:      "some-id",
				OrganizationID: s.resources[0].OrganizationID,
				NamespaceID:    s.resources[0].NamespaceID,
				UserID:         s.resources[0].UserID,
			},
			ErrString: resource.ErrInvalidUUID.Error(),
		},
		{
			Description: "should return error if resource urn is empty",
			ErrString:   resource.ErrInvalidURN.Error(),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			got, err := s.repository.Upsert(s.ctx, tc.ResourceToCreate)
			if tc.ErrString != "" {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
			if !cmp.Equal(got, tc.ExpectedResource, cmpopts.IgnoreFields(resource.Resource{},
				"Idxa",
				"CreatedAt",
				"UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", got, tc.ExpectedResource)
			}
		})
	}
}

func (s *ResourceRepositoryTestSuite) TestCreate() {
	type testCase struct {
		Description      string
		ResourceToCreate resource.Resource
		ExpectedResource resource.Resource
		ErrString        string
	}

	testCases := []testCase{
		{
			Description: "should create a resource",
			ResourceToCreate: resource.Resource{
				URN:            "new-urn-4",
				Name:           "resource4",
				ProjectID:      s.resources[0].ProjectID,
				OrganizationID: s.resources[0].OrganizationID,
				NamespaceID:    s.resources[0].NamespaceID,
				UserID:         s.resources[0].UserID,
			},
			ExpectedResource: resource.Resource{
				URN:            "new-urn-4",
				Name:           "resource4",
				ProjectID:      s.resources[0].ProjectID,
				OrganizationID: s.resources[0].OrganizationID,
				NamespaceID:    s.resources[0].NamespaceID,
				UserID:         s.resources[0].UserID,
			},
		},
		{
			Description: "should return error if namespace id does not exist",
			ResourceToCreate: resource.Resource{
				URN:            "new-urn-notexist",
				Name:           "resource4",
				ProjectID:      s.resources[0].ProjectID,
				OrganizationID: s.resources[0].OrganizationID,
				NamespaceID:    "some-ns",
				UserID:         s.resources[0].UserID,
			},
			ErrString: resource.ErrInvalidDetail.Error(),
		},
		{
			Description: "should return error if org id does not exist",
			ResourceToCreate: resource.Resource{
				URN:            "new-urn-notexist",
				Name:           "resource4",
				ProjectID:      s.resources[0].ProjectID,
				OrganizationID: uuid.NewString(),
				NamespaceID:    s.resources[0].NamespaceID,
				UserID:         s.resources[0].UserID,
			},
			ErrString: resource.ErrInvalidDetail.Error(),
		},
		{
			Description: "should return error if org id is not uuid",
			ResourceToCreate: resource.Resource{
				URN:            "new-urn-notexist",
				Name:           "resource4",
				ProjectID:      s.resources[0].ProjectID,
				OrganizationID: "some-str",
				NamespaceID:    s.resources[0].NamespaceID,
				UserID:         s.resources[0].UserID,
			},
			ErrString: resource.ErrInvalidUUID.Error(),
		},
		{
			Description: "should return error if project id does not exist",
			ResourceToCreate: resource.Resource{
				URN:            "new-urn-notexist",
				Name:           "resource4",
				ProjectID:      uuid.NewString(),
				OrganizationID: s.resources[0].OrganizationID,
				NamespaceID:    s.resources[0].NamespaceID,
				UserID:         s.resources[0].UserID,
			},
			ErrString: resource.ErrInvalidDetail.Error(),
		},
		{
			Description: "should return error if project id is not uuid",
			ResourceToCreate: resource.Resource{
				URN:            "new-urn-notexist",
				Name:           "resource4",
				ProjectID:      "some-id",
				OrganizationID: s.resources[0].OrganizationID,
				NamespaceID:    s.resources[0].NamespaceID,
				UserID:         s.resources[0].UserID,
			},
			ErrString: resource.ErrInvalidUUID.Error(),
		},
		{
			Description: "should return error if resource urn is empty",
			ErrString:   resource.ErrInvalidURN.Error(),
		},
		{
			Description: "should return error if urn already exist",
			ResourceToCreate: resource.Resource{
				URN:            s.resources[0].URN,
				Name:           s.resources[0].Name,
				ProjectID:      s.resources[0].ProjectID,
				OrganizationID: s.resources[0].OrganizationID,
				NamespaceID:    s.resources[0].NamespaceID,
				UserID:         s.resources[0].UserID,
			},
			ErrString: resource.ErrConflict.Error(),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			got, err := s.repository.Create(s.ctx, tc.ResourceToCreate)
			if tc.ErrString != "" {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
			if !cmp.Equal(got, tc.ExpectedResource, cmpopts.IgnoreFields(resource.Resource{},
				"Idxa",
				"CreatedAt",
				"UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", got, tc.ExpectedResource)
			}
		})
	}
}

func (s *ResourceRepositoryTestSuite) TestList() {
	type testCase struct {
		Description       string
		Filter            resource.Filter
		ExpectedResources []resource.Resource
		ErrString         string
	}

	testCases := []testCase{
		{
			Description:       "should get all resources",
			ExpectedResources: s.resources,
		},
		{
			Description: "should get filtered resources",
			Filter: resource.Filter{
				ProjectID:      s.projects[1].ID,
				OrganizationID: s.orgs[1].ID,
				NamespaceID:    s.namespaces[1].ID,
			},
			ExpectedResources: []resource.Resource{
				{
					Idxa:           s.resources[1].Idxa,
					URN:            s.resources[1].URN,
					Name:           s.resources[1].Name,
					ProjectID:      s.resources[1].ProjectID,
					OrganizationID: s.resources[1].OrganizationID,
					NamespaceID:    s.resources[1].NamespaceID,
					UserID:         s.resources[1].UserID,
				},
			},
		},
		{
			Description: "should return the page given if filter page given is 1 or greater",
			Filter: resource.Filter{
				Page:  1,
				Limit: 2,
			},
			ExpectedResources: []resource.Resource{
				{
					Idxa:           s.resources[0].Idxa,
					URN:            s.resources[0].URN,
					Name:           s.resources[0].Name,
					ProjectID:      s.resources[0].ProjectID,
					OrganizationID: s.resources[0].OrganizationID,
					NamespaceID:    s.resources[0].NamespaceID,
					UserID:         s.resources[0].UserID,
				},
				{
					Idxa:           s.resources[1].Idxa,
					URN:            s.resources[1].URN,
					Name:           s.resources[1].Name,
					ProjectID:      s.resources[1].ProjectID,
					OrganizationID: s.resources[1].OrganizationID,
					NamespaceID:    s.resources[1].NamespaceID,
					UserID:         s.resources[1].UserID,
				},
			},
		},
		{
			Description: "should return 1st page if filter page given is 0 or less",
			Filter: resource.Filter{
				Page:  0,
				Limit: 2,
			},
			ExpectedResources: []resource.Resource{
				{
					Idxa:           s.resources[0].Idxa,
					URN:            s.resources[0].URN,
					Name:           s.resources[0].Name,
					ProjectID:      s.resources[0].ProjectID,
					OrganizationID: s.resources[0].OrganizationID,
					NamespaceID:    s.resources[0].NamespaceID,
					UserID:         s.resources[0].UserID,
				},
				{
					Idxa:           s.resources[1].Idxa,
					URN:            s.resources[1].URN,
					Name:           s.resources[1].Name,
					ProjectID:      s.resources[1].ProjectID,
					OrganizationID: s.resources[1].OrganizationID,
					NamespaceID:    s.resources[1].NamespaceID,
					UserID:         s.resources[1].UserID,
				},
			},
		},
		{
			Description: "should return list of users with maximum 50 data if limit given is 0 or less",
			Filter: resource.Filter{
				Limit: 0,
				Page:  1,
			},
			ExpectedResources: s.resources,
		},
		{
			Description: "should return first page of filtered resources based on search filters",
			Filter: resource.Filter{
				Page:           1,
				Limit:          2,
				ProjectID:      s.projects[1].ID,
				OrganizationID: s.orgs[1].ID,
				NamespaceID:    s.namespaces[1].ID,
			},
			ExpectedResources: []resource.Resource{
				{
					Idxa:           s.resources[1].Idxa,
					URN:            s.resources[1].URN,
					Name:           s.resources[1].Name,
					ProjectID:      s.resources[1].ProjectID,
					OrganizationID: s.resources[1].OrganizationID,
					NamespaceID:    s.resources[1].NamespaceID,
					UserID:         s.resources[1].UserID,
				},
			},
		},
		{
			Description: "should return second page of filtered resources based on search filters",
			Filter: resource.Filter{
				Page:           2,
				Limit:          2,
				ProjectID:      s.projects[1].ID,
				OrganizationID: s.orgs[1].ID,
				NamespaceID:    s.namespaces[1].ID,
			},
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
			if !cmp.Equal(got, tc.ExpectedResources, cmpopts.IgnoreFields(resource.Resource{},
				"CreatedAt",
				"UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", got, tc.ExpectedResources)
			}
		})
	}
}

func (s *ResourceRepositoryTestSuite) TestUpdate() {
	type testCase struct {
		Description      string
		ResourceID       string
		ResourceToUpdate resource.Resource
		ExpectedResource resource.Resource
		ErrString        string
	}

	testCases := []testCase{
		{
			Description: "should update a resource",
			ResourceID:  s.resources[0].Idxa,
			ResourceToUpdate: resource.Resource{
				Name:           "resource-1",
				ProjectID:      s.resources[0].ProjectID,
				OrganizationID: s.resources[0].OrganizationID,
				NamespaceID:    s.resources[0].NamespaceID,
			},
			ExpectedResource: resource.Resource{
				Idxa:           s.resources[0].Idxa,
				URN:            "resource-1-urn",
				Name:           "resource-1",
				ProjectID:      s.resources[0].ProjectID,
				OrganizationID: s.resources[0].OrganizationID,
				NamespaceID:    s.resources[0].NamespaceID,
				UserID:         s.resources[0].UserID,
			},
			ErrString: "",
		},
		{
			Description: "should return error if namespace id does not exist",
			ResourceID:  s.resources[0].Idxa,
			ResourceToUpdate: resource.Resource{
				URN:            "new-urn-notexist",
				Name:           "resource4",
				ProjectID:      s.resources[0].ProjectID,
				OrganizationID: s.resources[0].OrganizationID,
				NamespaceID:    "some-ns",
				UserID:         s.resources[0].UserID,
			},
			ErrString: resource.ErrNotExist.Error(),
		},
		{
			Description: "should return error if org id does not exist",
			ResourceID:  s.resources[0].Idxa,
			ResourceToUpdate: resource.Resource{
				URN:            "new-urn-notexist",
				Name:           "resource4",
				ProjectID:      s.resources[0].ProjectID,
				OrganizationID: uuid.NewString(),
				NamespaceID:    s.resources[0].NamespaceID,
				UserID:         s.resources[0].UserID,
			},
			ErrString: resource.ErrNotExist.Error(),
		},
		{
			Description: "should return error if org id is not uuid",
			ResourceID:  s.resources[0].Idxa,
			ResourceToUpdate: resource.Resource{
				URN:            "new-urn-notexist",
				Name:           "resource4",
				ProjectID:      s.resources[0].ProjectID,
				OrganizationID: "some-str",
				NamespaceID:    s.resources[0].NamespaceID,
				UserID:         s.resources[0].UserID,
			},
			ErrString: resource.ErrInvalidDetail.Error(),
		},
		{
			Description: "should return error if project id does not exist",
			ResourceID:  s.resources[0].Idxa,
			ResourceToUpdate: resource.Resource{
				URN:            "new-urn-notexist",
				Name:           "resource4",
				ProjectID:      uuid.NewString(),
				OrganizationID: s.resources[0].OrganizationID,
				NamespaceID:    s.resources[0].NamespaceID,
				UserID:         s.resources[0].UserID,
			},
			ErrString: resource.ErrNotExist.Error(),
		},
		{
			Description: "should return error if project id is not uuid",
			ResourceID:  s.resources[0].Idxa,
			ResourceToUpdate: resource.Resource{
				URN:            "new-urn-notexist",
				Name:           "resource4",
				ProjectID:      "some-id",
				OrganizationID: s.resources[0].OrganizationID,
				NamespaceID:    s.resources[0].NamespaceID,
				UserID:         s.resources[0].UserID,
			},
			ErrString: resource.ErrInvalidDetail.Error(),
		},
		{
			Description: "should return error if resource id is empty",
			ErrString:   resource.ErrInvalidID.Error(),
		},
		{
			Description: "should return error if resource id is invalid",
			ResourceID:  "abc",
			ErrString:   resource.ErrInvalidUUID.Error(),
		},
		{
			Description: "should return error if resource urn is empty",
			ResourceID:  uuid.NewString(),
			ErrString:   resource.ErrInvalidDetail.Error(),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			got, err := s.repository.Update(s.ctx, tc.ResourceID, tc.ResourceToUpdate)
			if tc.ErrString != "" {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
			if !cmp.Equal(got, tc.ExpectedResource, cmpopts.IgnoreFields(resource.Resource{},
				"CreatedAt",
				"UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", got, tc.ExpectedResource)
			}
		})
	}
}

func (s *ResourceRepositoryTestSuite) TestGetByNamespace() {
	type testCase struct {
		Description string
		Name        string
		Namespace   string
		Expected    resource.Resource
		ErrString   string
	}

	testCases := []testCase{
		{
			Description: "should get a resource",
			Name:        s.resources[0].Name,
			Namespace:   s.resources[0].NamespaceID,
			Expected:    s.resources[0],
		},
		{
			Description: "should return error no exist if can't found resource",
			Name:        "some-urn",
			ErrString:   resource.ErrNotExist.Error(),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			got, err := s.repository.GetByNamespace(s.ctx, tc.Name, tc.Namespace)
			if tc.ErrString != "" {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
			if !cmp.Equal(got, tc.Expected, cmpopts.IgnoreFields(resource.Resource{},
				"CreatedAt",
				"UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", got, tc.Expected)
			}
		})
	}
}

func (s *ResourceRepositoryTestSuite) TestGetSchema() {
	expected := schema.NamespaceConfigMapType{}
	for _, rc := range s.resourceConfig {
		var config schema.NamespaceConfigMapType
		if err := json.Unmarshal([]byte(rc.Config), &config); err != nil {
			s.T().Fatal(err)
		}
		expected = schema.MergeNamespaceConfigMap(expected, config)
	}

	type testCase struct {
		Description string
		Expected    schema.NamespaceConfigMapType
		ErrString   string
	}

	testCases := []testCase{
		{
			Description: "should get all resources configs",
			Expected:    expected,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			got, err := s.repository.GetSchema(s.ctx)
			if tc.ErrString != "" {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
			if !cmp.Equal(got, tc.Expected, cmpopts.IgnoreFields(resource.ResourceConfig{},
				"CreatedAt",
				"UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", got, tc.Expected)
			}
		})
	}
}

func (s *ResourceRepositoryTestSuite) TestUpsertResourceConfigs() {
	type testCase struct {
		Description string
		Name        string
		Config      schema.NamespaceConfigMapType
		Expected    resource.ResourceConfig
		ErrString   string
	}

	testCases := []testCase{
		{
			Description: "should create a resource config",
			Name:        "test",
			Config: schema.NamespaceConfigMapType{
				"test/resource": schema.NamespaceConfig{
					Type:        "resource_group_namespace",
					Roles:       map[string][]string{},
					Permissions: map[string][]string{},
				},
			},
			Expected: resource.ResourceConfig{
				ID:     3,
				Name:   "test",
				Config: "{\"test/resource\": {\"Type\": \"resource_group_namespace\", \"Roles\": {}, \"Permissions\": {}, \"InheritedNamespaces\": null}}",
			},
		},
		{
			Description: "should update a resource config",
			Name:        s.resourceConfig[0].Name,
			Config: schema.NamespaceConfigMapType{
				"test/resource": schema.NamespaceConfig{
					Type:        "resource_group_namespace",
					Roles:       map[string][]string{},
					Permissions: map[string][]string{},
				},
			},
			Expected: resource.ResourceConfig{
				ID:     s.resourceConfig[0].ID,
				Name:   s.resourceConfig[0].Name,
				Config: "{\"test/resource\": {\"Type\": \"resource_group_namespace\", \"Roles\": {}, \"Permissions\": {}, \"InheritedNamespaces\": null}}",
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			got, err := s.repository.UpsertResourceConfigs(s.ctx, tc.Name, tc.Config)
			if tc.ErrString != "" {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
			if !cmp.Equal(got, tc.Expected, cmpopts.IgnoreFields(resource.ResourceConfig{},
				"CreatedAt",
				"UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", got, tc.Expected)
			}
		})
	}
}

func TestResourceRepository(t *testing.T) {
	suite.Run(t, new(ResourceRepositoryTestSuite))
}
