package postgres_test

import (
	"context"
	"fmt"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/goto/salt/log"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/suite"

	"github.com/goto/shield/core/group"
	"github.com/goto/shield/core/organization"
	"github.com/goto/shield/core/relation"
	"github.com/goto/shield/core/user"
	"github.com/goto/shield/internal/schema"
	"github.com/goto/shield/internal/store/postgres"
	"github.com/goto/shield/pkg/db"
)

type SortByName []group.Group

func (a SortByName) Len() int           { return len(a) }
func (a SortByName) Less(i, j int) bool { return a[i].Name < a[j].Name }
func (a SortByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type GroupRepositoryTestSuite struct {
	suite.Suite
	ctx                 context.Context
	client              *db.Client
	pool                *dockertest.Pool
	resource            *dockertest.Resource
	repository          *postgres.GroupRepository
	relationRepository  *postgres.RelationRepository
	namespaceRepository *postgres.NamespaceRepository
	roleRepository      *postgres.RoleRepository
	orgs                []organization.Organization
	groups              []group.Group
	users               []user.User
}

func (s *GroupRepositoryTestSuite) SetupSuite() {
	var err error

	logger := log.NewZap()
	s.client, s.pool, s.resource, err = newTestClient(logger)
	if err != nil {
		s.T().Fatal(err)
	}

	s.ctx = context.TODO()
	s.repository = postgres.NewGroupRepository(s.client)

	_, err = bootstrapMetadataKeys(s.client)
	if err != nil {
		s.T().Fatal(err)
	}
	s.users, err = bootstrapUser(s.client)
	if err != nil {
		s.T().Fatal(err)
	}

	s.relationRepository = postgres.NewRelationRepository(s.client)
	s.namespaceRepository = postgres.NewNamespaceRepository(s.client)
	s.roleRepository = postgres.NewRoleRepository(s.client)

	s.orgs, err = bootstrapOrganization(s.client)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *GroupRepositoryTestSuite) SetupTest() {
	var err error
	s.groups, err = bootstrapGroup(s.client, s.orgs)
	if err != nil {
		s.T().Fatal(err)
	}

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

	for _, group := range s.groups {
		_, err = s.relationRepository.Create(context.Background(), relation.RelationV2{
			Subject: relation.Subject{
				ID:        s.users[0].ID,
				Namespace: schema.UserPrincipal,
				RoleID:    schema.MemberRole,
			},
			Object: relation.Object{
				ID:          group.ID,
				NamespaceID: schema.GroupNamespace,
			},
		})
		if err != nil {
			s.T().Fatal(err)
		}
	}

	for _, user := range s.users {
		_, err = s.relationRepository.Create(context.Background(), relation.RelationV2{
			Subject: relation.Subject{
				ID:        user.ID,
				Namespace: schema.UserPrincipal,
				RoleID:    schema.MemberRole,
			},
			Object: relation.Object{
				ID:          s.groups[0].ID,
				NamespaceID: schema.GroupNamespace,
			},
		})
		if err != nil {
			s.T().Fatal(err)
		}
	}
}

func (s *GroupRepositoryTestSuite) TearDownSuite() {
	// Clean tests
	if err := purgeDocker(s.pool, s.resource); err != nil {
		s.T().Fatal(err)
	}
}

func (s *GroupRepositoryTestSuite) TearDownTest() {
	if err := s.cleanup(); err != nil {
		s.T().Fatal(err)
	}
}

func (s *GroupRepositoryTestSuite) cleanup() error {
	queries := []string{
		fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", postgres.TABLE_GROUPS),
		fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", postgres.TABLE_RELATIONS),
		fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", postgres.TABLE_ROLES),
		fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", postgres.TABLE_NAMESPACES),
	}
	return execQueries(context.TODO(), s.client, queries)
}

func (s *GroupRepositoryTestSuite) TestGetByID() {
	type testCase struct {
		Description   string
		SelectedID    string
		ExpectedGroup group.Group
		ErrString     string
	}

	testCases := []testCase{
		{
			Description: "should get a group",
			SelectedID:  s.groups[0].ID,
			ExpectedGroup: group.Group{
				Name:           "group1",
				Slug:           "group-1",
				OrganizationID: s.groups[0].OrganizationID,
			},
		},
		{
			Description: "should return error no exist if can't found group",
			SelectedID:  uuid.NewString(),
			ErrString:   group.ErrNotExist.Error(),
		},
		{
			Description: "should return error if id empty",
			ErrString:   group.ErrInvalidID.Error(),
		},
		{
			Description: "should return error if id is not uuid",
			SelectedID:  "10000",
			ErrString:   group.ErrInvalidUUID.Error(),
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
			if !cmp.Equal(got, tc.ExpectedGroup, cmpopts.IgnoreFields(group.Group{},
				"ID",
				"Metadata",
				"CreatedAt",
				"UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", got, tc.ExpectedGroup)
			}
		})
	}
}

func (s *GroupRepositoryTestSuite) TestGetByIDs() {
	type testCase struct {
		Description    string
		SelectedIDs    []string
		ExpectedGroups []group.Group
		ErrString      string
	}

	testCases := []testCase{
		{
			Description: "should get a group",
			SelectedIDs: []string{s.groups[0].ID, s.groups[1].ID},
			ExpectedGroups: []group.Group{
				{
					Name:           "group1",
					Slug:           "group-1",
					OrganizationID: s.groups[0].OrganizationID,
				}, {
					Name:           "group2",
					Slug:           "group-2",
					OrganizationID: s.groups[1].OrganizationID,
				},
			},
		},
		{
			Description: "should return error if id empty",
			SelectedIDs: []string{s.groups[0].ID, ""},
			ErrString:   group.ErrInvalidID.Error(),
		},
		{
			Description: "should return error if id is not uuid",
			SelectedIDs: []string{s.groups[0].ID, "10000"},
			ErrString:   group.ErrInvalidUUID.Error(),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			got, err := s.repository.GetByIDs(s.ctx, tc.SelectedIDs)
			if tc.ErrString != "" && err != nil {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}

			for i, grp := range got {
				if !cmp.Equal(grp, tc.ExpectedGroups[i], cmpopts.IgnoreFields(group.Group{},
					"ID",
					"Metadata",
					"CreatedAt",
					"UpdatedAt")) {
					s.T().Fatalf("got result %+v, expected was %+v", grp, tc.ExpectedGroups[i])
				}
			}
		})
	}
}

func (s *GroupRepositoryTestSuite) TestGetBySlug() {
	type testCase struct {
		Description   string
		SelectedSlug  string
		ExpectedGroup group.Group
		ErrString     string
	}

	testCases := []testCase{
		{
			Description:  "should get a group",
			SelectedSlug: "group-1",
			ExpectedGroup: group.Group{
				Name:           "group1",
				Slug:           "group-1",
				OrganizationID: s.groups[0].OrganizationID,
			},
		},
		{
			Description:  "should return error no exist if can't found group",
			SelectedSlug: "randomslug",
			ErrString:    group.ErrNotExist.Error(),
		},
		{
			Description: "should return error if slug empty",
			ErrString:   group.ErrInvalidID.Error(),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			got, err := s.repository.GetBySlug(s.ctx, tc.SelectedSlug)
			if tc.ErrString != "" {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
			if !cmp.Equal(got, tc.ExpectedGroup, cmpopts.IgnoreFields(group.Group{}, "ID", "Metadata", "CreatedAt", "UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", got, tc.ExpectedGroup)
			}
		})
	}
}

func (s *GroupRepositoryTestSuite) TestCreate() {
	type testCase struct {
		Description   string
		GroupToCreate group.Group
		ExpectedGroup group.Group
		ErrString     string
	}

	testCases := []testCase{
		{
			Description: "should create a group",
			GroupToCreate: group.Group{
				Name:           "new-group",
				Slug:           "new-group-slug",
				OrganizationID: s.orgs[0].ID,
			},
			ExpectedGroup: group.Group{
				Name:           "new-group",
				Slug:           "new-group-slug",
				OrganizationID: s.orgs[0].ID,
			},
		},
		{
			Description: "should return error if group name already exist",
			GroupToCreate: group.Group{
				Name:           "group2",
				Slug:           "new-slug",
				OrganizationID: s.orgs[0].ID,
			},
			ErrString: group.ErrConflict.Error(),
		},
		{
			Description: "should return error if group slug already exist",
			GroupToCreate: group.Group{
				Name:           "newslug",
				Slug:           "group-2",
				OrganizationID: s.orgs[0].ID,
			},
			ErrString: group.ErrConflict.Error(),
		},
		{
			Description: "should return error if org id not an uuid",
			GroupToCreate: group.Group{
				Name:           "newslug",
				Slug:           "groupnewslug",
				OrganizationID: "some-id",
			},
			ErrString: group.ErrInvalidUUID.Error(),
		},
		{
			Description: "should return error if org id does not exist",
			GroupToCreate: group.Group{
				Name:           "newslug",
				Slug:           "groupnewslug",
				OrganizationID: uuid.NewString(),
			},
			ErrString: organization.ErrNotExist.Error(),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			got, err := s.repository.Create(s.ctx, tc.GroupToCreate)
			if tc.ErrString != "" {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
			if !cmp.Equal(got, tc.ExpectedGroup, cmpopts.IgnoreFields(group.Group{}, "ID", "Metadata", "CreatedAt", "UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", got, tc.ExpectedGroup)
			}
		})
	}
}

func (s *GroupRepositoryTestSuite) TestList() {
	type testCase struct {
		Description    string
		Filter         group.Filter
		ExpectedGroups []group.Group
		ErrString      string
	}

	testCases := []testCase{
		{
			Description: "should get all groups",
			ExpectedGroups: []group.Group{
				{
					Name:           "group1",
					Slug:           "group-1",
					OrganizationID: s.orgs[0].ID,
				},
				{
					Name:           "group2",
					Slug:           "group-2",
					OrganizationID: s.orgs[0].ID,
				},
				{
					Name:           "group3",
					Slug:           "group-3",
					OrganizationID: s.orgs[1].ID,
				},
			},
		},
		{
			Description: "should get filtered groups",
			Filter: group.Filter{
				OrganizationID: s.orgs[1].ID,
			},
			ExpectedGroups: []group.Group{
				{
					Name:           "group3",
					Slug:           "group-3",
					OrganizationID: s.orgs[1].ID,
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

			sort.Sort(SortByName(got))

			if !cmp.Equal(got, tc.ExpectedGroups, cmpopts.IgnoreFields(group.Group{}, "ID", "Metadata", "CreatedAt", "UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", got, tc.ExpectedGroups)
			}
		})
	}
}

func (s *GroupRepositoryTestSuite) TestUpdateByID() {
	type testCase struct {
		Description   string
		GroupToUpdate group.Group
		ExpectedGroup group.Group
		ErrString     string
	}

	testCases := []testCase{
		{
			Description: "should update a group",
			GroupToUpdate: group.Group{
				ID:             s.groups[0].ID,
				Name:           "new group update",
				Slug:           "new-group-update",
				OrganizationID: s.orgs[0].ID,
			},
			ExpectedGroup: group.Group{
				Name:           "new group update",
				Slug:           "new-group-update",
				OrganizationID: s.orgs[0].ID,
			},
		},
		{
			Description: "should return error if group name already exist",
			GroupToUpdate: group.Group{
				ID:             s.groups[0].ID,
				Name:           "group2",
				Slug:           "new-slug",
				OrganizationID: s.orgs[0].ID,
			},
			ErrString: group.ErrConflict.Error(),
		},
		{
			Description: "should return error if group slug already exist",
			GroupToUpdate: group.Group{
				ID:             s.groups[0].ID,
				Name:           "new-group-2",
				Slug:           "group-2",
				OrganizationID: s.orgs[0].ID,
			},
			ErrString: group.ErrConflict.Error(),
		},
		{
			Description: "should return error if group not found",
			GroupToUpdate: group.Group{
				ID:             uuid.NewString(),
				Name:           "not-exist",
				Slug:           "some-slug",
				OrganizationID: s.orgs[0].ID,
			},
			ErrString: group.ErrNotExist.Error(),
		},
		{
			Description: "should return error if group id is not uuid",
			GroupToUpdate: group.Group{
				ID:             "12345",
				Name:           "not-exist",
				Slug:           "some-slug",
				OrganizationID: s.orgs[0].ID,
			},
			ErrString: group.ErrInvalidUUID.Error(),
		},
		{
			Description: "should return error if org id is not uuid",
			GroupToUpdate: group.Group{
				ID:             s.groups[0].ID,
				Slug:           "new-prj",
				Name:           "not-exist",
				OrganizationID: "not-uuid",
			},
			ErrString: organization.ErrInvalidUUID.Error(),
		},
		{
			Description: "should return error if org id not exist",
			GroupToUpdate: group.Group{
				ID:             s.groups[0].ID,
				Slug:           "new-prj",
				Name:           "not-exist",
				OrganizationID: uuid.NewString(),
			},
			ErrString: organization.ErrNotExist.Error(),
		},
		{
			Description: "should return error if group id is empty",
			ErrString:   group.ErrInvalidID.Error(),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			got, err := s.repository.UpdateByID(s.ctx, tc.GroupToUpdate)
			if tc.ErrString != "" {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
			if !cmp.Equal(got, tc.ExpectedGroup, cmpopts.IgnoreFields(group.Group{}, "ID", "Metadata", "CreatedAt", "UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", got, tc.ExpectedGroup)
			}
		})
	}
}

func (s *GroupRepositoryTestSuite) TestUpdateBySlug() {
	type testCase struct {
		Description   string
		GroupToUpdate group.Group
		ExpectedGroup group.Group
		ErrString     string
	}

	testCases := []testCase{
		{
			Description: "should update a group",
			GroupToUpdate: group.Group{
				Name:           "new group update",
				Slug:           "group-1",
				OrganizationID: s.orgs[0].ID,
			},
			ExpectedGroup: group.Group{
				Name:           "new group update",
				Slug:           "group-1",
				OrganizationID: s.orgs[0].ID,
			},
		},
		{
			Description: "should return error if group name already exist",
			GroupToUpdate: group.Group{
				Name:           "group2",
				Slug:           "group-1",
				OrganizationID: s.orgs[0].ID,
			},
			ErrString: group.ErrConflict.Error(),
		},
		{
			Description: "should return error if group not found",
			GroupToUpdate: group.Group{
				Slug:           "slug",
				Name:           "not-exist",
				OrganizationID: s.orgs[0].ID,
			},
			ErrString: group.ErrNotExist.Error(),
		},
		{
			Description: "should return error if org id is not uuid",
			GroupToUpdate: group.Group{
				Slug:           "group-1",
				Name:           "not-exist",
				OrganizationID: "not-uuid",
			},
			ErrString: organization.ErrInvalidUUID.Error(),
		},
		{
			Description: "should return error if org id not exist",
			GroupToUpdate: group.Group{
				Slug:           "group-1",
				Name:           "not-exist",
				OrganizationID: uuid.NewString(),
			},
			ErrString: organization.ErrNotExist.Error(),
		},
		{
			Description: "should return error if group slug is empty",
			ErrString:   group.ErrInvalidID.Error(),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			got, err := s.repository.UpdateBySlug(s.ctx, tc.GroupToUpdate)
			if tc.ErrString != "" {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
			if !cmp.Equal(got, tc.ExpectedGroup, cmpopts.IgnoreFields(group.Group{}, "ID", "Metadata", "CreatedAt", "UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", got, tc.ExpectedGroup)
			}
		})
	}
}

func (s *GroupRepositoryTestSuite) TestListUserGroups() {
	type testCase struct {
		Description    string
		UserID         string
		RoleID         string
		ExpectedGroups []group.Group
		ErrString      string
	}

	testCases := []testCase{
		{
			Description: "should get a list of group",
			UserID:      s.users[0].ID,
			RoleID:      "shield/group:member",
			ExpectedGroups: []group.Group{
				{
					Name:           "group1",
					Slug:           "group-1",
					OrganizationID: s.groups[0].OrganizationID,
				},
				{
					Name:           "group2",
					Slug:           "group-2",
					OrganizationID: s.groups[1].OrganizationID,
				},
				{
					Name:           "group3",
					Slug:           "group-3",
					OrganizationID: s.groups[2].OrganizationID,
				},
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			got, err := s.repository.ListUserGroups(s.ctx, tc.UserID, tc.RoleID)
			if tc.ErrString != "" && err != nil {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}

			for i, grp := range got {
				if !cmp.Equal(grp, tc.ExpectedGroups[i], cmpopts.IgnoreFields(group.Group{},
					"ID",
					"Metadata",
					"CreatedAt",
					"UpdatedAt")) {
					s.T().Fatalf("got result %+v, expected was %+v", grp, tc.ExpectedGroups[i])
				}
			}
		})
	}
}

func (s *GroupRepositoryTestSuite) TestListGroupRelations() {
	type testCase struct {
		Description       string
		ObjectID          string
		SubjectType       string
		Role              string
		ExpectedRelations []relation.RelationV2
		ErrString         string
	}

	testCases := []testCase{
		{
			Description: "should get a list of relations",
			ObjectID:    s.groups[0].ID,
			SubjectType: "user",
			Role:        "member",
			ExpectedRelations: []relation.RelationV2{
				{
					Object: relation.Object{
						ID:          s.groups[0].ID,
						NamespaceID: schema.GroupNamespace,
					},
					Subject: relation.Subject{
						ID:        s.users[0].ID,
						Namespace: schema.UserPrincipal,
						RoleID:    "shield/group:member",
					},
				},
				{
					Object: relation.Object{
						ID:          s.groups[0].ID,
						NamespaceID: schema.GroupNamespace,
					},
					Subject: relation.Subject{
						ID:        s.users[1].ID,
						Namespace: schema.UserPrincipal,
						RoleID:    "shield/group:member",
					},
				},
				{
					Object: relation.Object{
						ID:          s.groups[0].ID,
						NamespaceID: schema.GroupNamespace,
					},
					Subject: relation.Subject{
						ID:        s.users[2].ID,
						Namespace: schema.UserPrincipal,
						RoleID:    "shield/group:member",
					},
				},
				{
					Object: relation.Object{
						ID:          s.groups[0].ID,
						NamespaceID: schema.GroupNamespace,
					},
					Subject: relation.Subject{
						ID:        s.users[3].ID,
						Namespace: schema.UserPrincipal,
						RoleID:    "shield/group:member",
					},
				},
				{
					Object: relation.Object{
						ID:          s.groups[0].ID,
						NamespaceID: schema.GroupNamespace,
					},
					Subject: relation.Subject{
						ID:        s.users[4].ID,
						Namespace: schema.UserPrincipal,
						RoleID:    "shield/group:member",
					},
				},
				{
					Object: relation.Object{
						ID:          s.groups[0].ID,
						NamespaceID: schema.GroupNamespace,
					},
					Subject: relation.Subject{
						ID:        s.users[5].ID,
						Namespace: schema.UserPrincipal,
						RoleID:    "shield/group:member",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			got, err := s.repository.ListGroupRelations(s.ctx, tc.ObjectID, tc.SubjectType, tc.Role)
			if tc.ErrString != "" && err != nil {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}

			for _, rel := range got {
				found := false
				for _, expectedRel := range tc.ExpectedRelations {
					if cmp.Equal(rel, expectedRel, cmpopts.IgnoreFields(relation.RelationV2{},
						"ID",
						"CreatedAt",
						"UpdatedAt")) {
						found = true
						break
					}
				}
				if !found {
					s.T().Fatalf("can't find relation %+v", rel)
				}
			}
		})
	}
}

func TestGroupRepository(t *testing.T) {
	suite.Run(t, new(GroupRepositoryTestSuite))
}
