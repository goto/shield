package e2e_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/goto/shield/config"
	shieldv1beta1 "github.com/goto/shield/proto/v1beta1"
	"github.com/goto/shield/test/e2e_test/testbench"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/metadata"
)

type EndToEndAPISmokeTestSuite struct {
	suite.Suite
	client       shieldv1beta1.ShieldServiceClient
	cancelClient func()
	testBench    *testbench.TestBench
	appConfig    *config.Shield
	users        []*shieldv1beta1.User
}

func (s *EndToEndAPISmokeTestSuite) SetupTest() {
	ctx := context.Background()
	s.client, s.appConfig, s.cancelClient, _ = testbench.SetupTests(s.T())

	// validate
	uRes, err := s.client.ListUsers(ctx, &shieldv1beta1.ListUsersRequest{})
	s.Require().NoError(err)
	s.Require().Equal(9, len(uRes.GetUsers()))
	s.users = uRes.GetUsers()

	oRes, err := s.client.ListOrganizations(ctx, &shieldv1beta1.ListOrganizationsRequest{})
	s.Require().NoError(err)
	s.Require().Equal(1, len(oRes.GetOrganizations()))

	pRes, err := s.client.ListProjects(ctx, &shieldv1beta1.ListProjectsRequest{})
	s.Require().NoError(err)
	s.Require().Equal(1, len(pRes.GetProjects()))

	gRes, err := s.client.ListGroups(ctx, &shieldv1beta1.ListGroupsRequest{})
	s.Require().NoError(err)
	s.Require().Equal(3, len(gRes.GetGroups()))
}

func (s *EndToEndAPISmokeTestSuite) TearDownTest() {
	s.cancelClient()
	// Clean tests
	err := s.testBench.CleanUp()
	s.Require().NoError(err)
}

func (s *EndToEndAPISmokeTestSuite) TestUserAPI() {
	ctxOrgAdminAuth := metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{
		testbench.IdentityHeader: testbench.OrgAdminEmail,
	}))

	s.Run("1. get user with uuid should return a correct user detail", func() {
		res, err := s.client.GetUser(ctxOrgAdminAuth, &shieldv1beta1.GetUserRequest{
			Id: s.users[0].GetId(),
		})

		s.Assert().Empty(cmp.Diff(s.users[0], res.GetUser(),
			cmpopts.IgnoreUnexported(shieldv1beta1.User{}),
			cmpopts.IgnoreFields(shieldv1beta1.User{}, "Metadata", "CreatedAt", "UpdatedAt"),
		))
		s.Assert().NoError(err)
	})

	s.Run("2. get user with e-mail should return a correct user detail", func() {
		res, err := s.client.GetUser(ctxOrgAdminAuth, &shieldv1beta1.GetUserRequest{
			Id: s.users[1].GetEmail(),
		})

		s.Assert().Empty(cmp.Diff(s.users[1], res.GetUser(),
			cmpopts.IgnoreUnexported(shieldv1beta1.User{}),
			cmpopts.IgnoreFields(shieldv1beta1.User{}, "Metadata", "CreatedAt", "UpdatedAt"),
		))
		s.Assert().NoError(err)
	})
}

func TestEndToEndAPISmokeTestSuite(t *testing.T) {
	suite.Run(t, new(EndToEndAPISmokeTestSuite))
}
