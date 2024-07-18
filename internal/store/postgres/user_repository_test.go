package postgres_test

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/suite"

	"github.com/goto/salt/log"
	"github.com/goto/shield/core/project"
	"github.com/goto/shield/core/resource"
	"github.com/goto/shield/core/servicedata"
	"github.com/goto/shield/core/user"
	"github.com/goto/shield/internal/store/postgres"
	"github.com/goto/shield/pkg/db"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	ctx        context.Context
	client     *db.Client
	pool       *dockertest.Pool
	resource   *dockertest.Resource
	repository *postgres.UserRepository
	keys       []servicedata.Key
	projects   []project.Project
	resources  []resource.Resource
	data       []servicedata.ServiceData
	users      []user.User
}

func (s *UserRepositoryTestSuite) SetupSuite() {
	var err error

	logger := log.NewZap()
	s.client, s.pool, s.resource, err = newTestClient(logger)
	if err != nil {
		s.T().Fatal(err)
	}

	s.ctx = context.TODO()
	s.repository = postgres.NewUserRepository(s.client)
}

func (s *UserRepositoryTestSuite) SetupTest() {
	var err error

	_, err = bootstrapMetadataKeys(s.client)
	if err != nil {
		s.T().Fatal(err)
	}
	s.users, err = bootstrapUser(s.client)
	if err != nil {
		s.T().Fatal(err)
	}

	namespaces, err := bootstrapNamespace(s.client)
	if err != nil {
		s.T().Fatal(err)
	}

	_, err = bootstrapMetadataKeys(s.client)
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

func (s *UserRepositoryTestSuite) TearDownSuite() {
	// Clean tests
	if err := purgeDocker(s.pool, s.resource); err != nil {
		s.T().Fatal(err)
	}
}

func (s *UserRepositoryTestSuite) TearDownTest() {
	if err := s.cleanup(); err != nil {
		s.T().Fatal(err)
	}
}

func (s *UserRepositoryTestSuite) cleanup() error {
	queries := []string{
		fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", postgres.TABLE_METADATA),
		fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", postgres.TABLE_USERS),
		fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", postgres.TABLE_METADATA_KEYS),
		fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", postgres.TABLE_SERVICE_DATA_KEYS),
		fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", postgres.TABLE_SERVICE_DATA),
		fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", postgres.TABLE_ORGANIZATIONS),
		fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", postgres.TABLE_PROJECTS),
		fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", postgres.TABLE_RESOURCES),
	}
	return execQueries(context.TODO(), s.client, queries)
}

func (s *UserRepositoryTestSuite) TestGetByID() {
	type testCase struct {
		Description  string
		SelectedID   string
		ExpectedUser user.User
		ErrString    string
	}

	testCases := []testCase{
		{
			Description: "should get a user",
			SelectedID:  s.users[0].ID,
			ExpectedUser: user.User{
				ID:       s.users[0].ID,
				Name:     s.users[0].Name,
				Email:    s.users[0].Email,
				Metadata: s.users[0].Metadata,
			},
		},
		{
			Description: "should return error if id is empty",
			SelectedID:  "",
			ErrString:   user.ErrInvalidID.Error(),
		},
		{
			Description: "should return error no exist if can't found user",
			SelectedID:  uuid.NewString(),
			ErrString:   user.ErrNotExist.Error(),
		},
		{
			Description: "should return error if id is not uuid",
			SelectedID:  "not-uuid",
			ErrString:   user.ErrInvalidUUID.Error(),
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
			if !cmp.Equal(got, tc.ExpectedUser, cmpopts.IgnoreFields(user.User{}, "ID", "Metadata", "CreatedAt", "UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", got, tc.ExpectedUser)
			}
		})
	}
}

func (s *UserRepositoryTestSuite) TestGetByEmail() {
	type testCase struct {
		Description   string
		SelectedEmail string
		ExpectedUser  user.User
		ErrString     string
	}

	testCases := []testCase{
		{
			Description:   "should get a user",
			SelectedEmail: s.users[0].Email,
			ExpectedUser: user.User{
				ID:       s.users[0].ID,
				Name:     s.users[0].Name,
				Email:    s.users[0].Email,
				Metadata: s.users[0].Metadata,
			},
		},
		{
			Description:   "should get a user with metadata",
			SelectedEmail: s.users[1].Email,
			ExpectedUser: user.User{
				ID:       s.users[1].ID,
				Name:     s.users[1].Name,
				Email:    s.users[1].Email,
				Metadata: s.users[1].Metadata,
			},
		},
		{
			Description:   "should return error if email is empty",
			SelectedEmail: "",
			ErrString:     user.ErrInvalidEmail.Error(),
		},
		{
			Description:   "should return error no exist if can't found user",
			SelectedEmail: "random",
			ErrString:     user.ErrNotExist.Error(),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			got, err := s.repository.GetByEmail(s.ctx, tc.SelectedEmail)
			if tc.ErrString != "" {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
			if !cmp.Equal(got, tc.ExpectedUser, cmpopts.IgnoreFields(user.User{}, "ID", "Metadata", "CreatedAt", "UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", got, tc.ExpectedUser)
			}
		})
	}
}

func (s *UserRepositoryTestSuite) TestCreate() {
	type testCase struct {
		Description   string
		UserToCreate  user.User
		ExpectedEmail string
		ErrString     string
	}

	testCases := []testCase{
		{
			Description: "should create a user",
			UserToCreate: user.User{
				Name:  "new user",
				Email: "new.user@gotocompany.com",
			},
			ExpectedEmail: "new.user@gotocompany.com",
		},
		{
			Description: "should return error if user already exist",
			UserToCreate: user.User{
				Name:  "new user",
				Email: "new.user@gotocompany.com",
			},
			ErrString: user.ErrConflict.Error(),
		},
		{
			Description: "should return error if email is empty",
			UserToCreate: user.User{
				Name:  "new user",
				Email: "",
			},
			ErrString: user.ErrInvalidEmail.Error(),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			got, err := s.repository.Create(s.ctx, tc.UserToCreate)
			if err != nil && tc.ErrString != "" {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
			if tc.ExpectedEmail != "" && (got.Email != tc.ExpectedEmail) {
				s.T().Fatalf("got result %+v, expected was %+v", got.ID, tc.ExpectedEmail)
			}
		})
	}
}

func (s *UserRepositoryTestSuite) TestList() {
	type testCase struct {
		Description   string
		Filter        user.Filter
		ExpectedUsers []user.User
		ErrString     string
	}

	testCases := []testCase{
		{
			Description:   "should get all users",
			ExpectedUsers: s.users,
		},
		{
			Description: "should return empty users if keyword not match any",
			Filter: user.Filter{
				Keyword: "random=keyword",
			},
		},
		{
			Description: "should return list of users if keyword match",
			Filter: user.Filter{
				Keyword: "alex",
			},
			ExpectedUsers: []user.User{
				{
					Name:  s.users[2].Name,
					Email: s.users[2].Email,
				},
				{
					Name:  s.users[3].Name,
					Email: s.users[3].Email,
				},
				{
					Name:  s.users[5].Name,
					Email: s.users[5].Email,
				},
			},
		},
		{
			Description: "should return 1 if filter with page",
			Filter: user.Filter{
				Limit: 1,
				Page:  1,
			},
			ExpectedUsers: []user.User{
				{
					Name:  s.users[0].Name,
					Email: s.users[0].Email,
				},
			},
		},
		{
			Description: "should return 1st page after filtering the users based on keywords",
			Filter: user.Filter{
				Keyword: "alex",
				Page:    1,
				Limit:   2,
			},
			ExpectedUsers: []user.User{
				{
					Name:  s.users[2].Name,
					Email: s.users[2].Email,
				},
				{
					Name:  s.users[3].Name,
					Email: s.users[3].Email,
				},
			},
		},
		{
			Description: "should return 2nd page after filtering the users based on keywords",
			Filter: user.Filter{
				Keyword: "alex",
				Page:    2,
				Limit:   2,
			},
			ExpectedUsers: []user.User{
				{
					Name:  s.users[5].Name,
					Email: s.users[5].Email,
				},
			},
		},
		{
			Description: "should return all users with keyword matching email or name",
			Filter: user.Filter{
				Keyword: "xu",
			},
			ExpectedUsers: []user.User{
				{
					Name:  s.users[2].Name,
					Email: s.users[2].Email,
				},
				{
					Name:  s.users[3].Name,
					Email: s.users[3].Email,
				},
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
			if !(len(got) == len(tc.ExpectedUsers)) {
				s.T().Fatalf("got result %+v, expected was %+v", got, tc.ExpectedUsers)
			}

			sort.Slice(got, func(i, j int) bool {
				return got[i].Email < got[j].Email
			})

			sort.Slice(tc.ExpectedUsers, func(i, j int) bool {
				return tc.ExpectedUsers[i].Email < tc.ExpectedUsers[j].Email
			})

			for idx := 0; idx < len(got); idx++ {
				if got[idx].Name != tc.ExpectedUsers[idx].Name || got[idx].Email != tc.ExpectedUsers[idx].Email {
					s.T().Fatalf("got user %+v, expected was %+v", got[idx], tc.ExpectedUsers[idx])
				}
			}
		})
	}
}

func (s *UserRepositoryTestSuite) TestGetByIDs() {
	type testCase struct {
		Description   string
		IDs           []string
		ExpectedUsers []user.User
		ErrString     string
	}

	testCases := []testCase{
		{
			Description:   "should get all users with ids",
			IDs:           []string{s.users[0].ID, s.users[0].ID},
			ExpectedUsers: []user.User{s.users[0]},
		},
		{
			Description: "should return empty users if ids not exist",
			IDs:         []string{uuid.NewString(), uuid.NewString()},
		},
		{
			Description:   "should return error if ids not uuid",
			IDs:           []string{"a", "b"},
			ExpectedUsers: []user.User{},
			ErrString:     user.ErrInvalidUUID.Error(),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			got, err := s.repository.GetByIDs(s.ctx, tc.IDs)
			if tc.ErrString != "" {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
			if !cmp.Equal(got, tc.ExpectedUsers, cmpopts.IgnoreFields(user.User{}, "ID", "Metadata", "CreatedAt", "UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", got, tc.ExpectedUsers)
			}
		})
	}
}

func (s *UserRepositoryTestSuite) TestUpdateByEmail() {
	type testCase struct {
		Description  string
		UserToUpdate user.User
		ExpectedUser user.User
		Err          error
	}

	testCases := []testCase{
		{
			Description: "should update a user",
			UserToUpdate: user.User{
				Name:  "Doe John",
				Email: s.users[0].Email,
			},
			ExpectedUser: user.User{
				Name:  "Doe John",
				Email: s.users[0].Email,
			},
		},
		{
			Description: "should return error if user not found",
			UserToUpdate: user.User{
				Email: "random@email.com",
			},
			Err: user.ErrNotExist,
		},
		{
			Description: "should return error if user email is empty",
			UserToUpdate: user.User{
				Email: "",
			},
			Err: user.ErrInvalidEmail,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			got, err := s.repository.UpdateByEmail(s.ctx, tc.UserToUpdate)
			if tc.Err != nil && tc.Err.Error() != "" {
				if errors.Unwrap(err) == tc.Err {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.Err)
				}
			}

			if !cmp.Equal(got, tc.ExpectedUser, cmpopts.IgnoreFields(user.User{},
				"ID", "CreatedAt", "UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", got, tc.ExpectedUser)
			}
		})
	}
}

func (s *UserRepositoryTestSuite) TestUpdateByID() {
	type testCase struct {
		Description  string
		UserToUpdate user.User
		ExpectedUser user.User
		Err          error
	}

	testCases := []testCase{
		{
			Description: "should update a user",
			UserToUpdate: user.User{
				ID:    s.users[0].ID,
				Name:  "Doe John",
				Email: s.users[0].Email,
			},
			ExpectedUser: user.User{
				ID:    s.users[0].ID,
				Name:  "Doe John",
				Email: s.users[0].Email,
			},
		},
		{
			Description: "should return error if user not found",
			UserToUpdate: user.User{
				ID:    uuid.NewString(),
				Name:  "Doe John",
				Email: "john.doe@gotocompany.com",
			},
			Err: user.ErrNotExist,
		},
		{
			Description: "should return error if user email already exist",
			UserToUpdate: user.User{
				ID:    s.users[0].ID,
				Name:  "Doe John",
				Email: s.users[1].Email,
			},
			Err: user.ErrConflict,
		},
		{
			Description: "should return error if user id is empty",
			Err:         user.ErrInvalidID,
		},
		{
			Description: "should return error if user id is not uuid",
			UserToUpdate: user.User{
				ID:    "abc",
				Name:  "Doe John",
				Email: s.users[1].Email,
			},
			Err: user.ErrInvalidID,
		},
		{
			Description: "should return error if email is empty",
			UserToUpdate: user.User{
				ID:    s.users[0].ID,
				Name:  "Doe John",
				Email: "",
			},
			Err: user.ErrInvalidEmail,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			got, err := s.repository.UpdateByID(s.ctx, tc.UserToUpdate)
			if tc.Err != nil && tc.Err.Error() != "" {
				if errors.Unwrap(err) == tc.Err {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.Err)
				}
			}
			if !cmp.Equal(got, tc.ExpectedUser, cmpopts.IgnoreFields(user.User{},
				"ID", "CreatedAt", "UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", got, tc.ExpectedUser)
			}
		})
	}
}

func TestUserRepository(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}
