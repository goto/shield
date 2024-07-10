package e2e_test

import (
	"context"
	"testing"

	"github.com/goto/shield/config"
	shieldv1beta1 "github.com/goto/shield/proto/v1beta1"
	"github.com/goto/shield/test/e2e_test/testbench"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

type EndToEndAPIRegressionTestSuite struct {
	suite.Suite
	client                  shieldv1beta1.ShieldServiceClient
	serviceDataClient       shieldv1beta1.ServiceDataServiceClient
	cancelClient            func()
	cancelServiceDataClient func()
	testBench               *testbench.TestBench
	appConfig               *config.Shield
}

func (s *EndToEndAPIRegressionTestSuite) SetupTest() {
	ctx := context.Background()
	ctxOrgAdminAuth := metadata.NewOutgoingContext(ctx, metadata.New(map[string]string{
		testbench.IdentityHeader: testbench.OrgAdminEmail,
	}))
	s.client, s.serviceDataClient, s.appConfig, s.cancelClient, s.cancelServiceDataClient, _ = testbench.SetupTests(s.T())

	// validate
	// list user length is 10 because there are 8 mock data, 1 system email, and 1 admin email created in test setup
	uRes, err := s.client.ListUsers(ctxOrgAdminAuth, &shieldv1beta1.ListUsersRequest{})
	s.Require().NoError(err)
	s.Require().Equal(10, len(uRes.GetUsers()))

	oRes, err := s.client.ListOrganizations(ctx, &shieldv1beta1.ListOrganizationsRequest{})
	s.Require().NoError(err)
	s.Require().Equal(1, len(oRes.GetOrganizations()))

	pRes, err := s.client.ListProjects(ctx, &shieldv1beta1.ListProjectsRequest{})
	s.Require().NoError(err)
	s.Require().Equal(2, len(pRes.GetProjects()))

	gRes, err := s.client.ListGroups(ctxOrgAdminAuth, &shieldv1beta1.ListGroupsRequest{})
	s.Require().NoError(err)
	s.Require().Equal(3, len(gRes.GetGroups()))

	rRes, err := s.client.ListResources(ctx, &shieldv1beta1.ListResourcesRequest{})
	s.Require().NoError(err)
	s.Require().Equal(5, len(rRes.GetResources()))
}

func (s *EndToEndAPIRegressionTestSuite) TearDownTest() {
	s.cancelClient()
	s.cancelServiceDataClient()
	// Clean tests
	err := s.testBench.CleanUp()
	s.Require().NoError(err)
}

func (s *EndToEndAPIRegressionTestSuite) TestProjectAPI() {
	var newProject *shieldv1beta1.Project

	ctxOrgAdminAuth := metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{
		testbench.IdentityHeader: testbench.OrgAdminEmail,
	}))

	// get my org
	res, err := s.client.ListOrganizations(context.Background(), &shieldv1beta1.ListOrganizationsRequest{})
	s.Require().NoError(err)
	s.Require().Greater(len(res.GetOrganizations()), 0)
	myOrg := res.GetOrganizations()[0]

	s.Run("1. org admin create a new project with empty auth email should return unauthenticated error", func() {
		_, err := s.client.CreateProject(context.Background(), &shieldv1beta1.CreateProjectRequest{
			Body: &shieldv1beta1.ProjectRequestBody{
				Name:  "new project",
				Slug:  "new-project",
				OrgId: myOrg.GetId(),
				Metadata: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"foo": structpb.NewBoolValue(true),
					},
				},
			},
		})
		s.Assert().Equal(codes.Unauthenticated, status.Convert(err).Code())
	})

	s.Run("2. org admin create a new project with empty name should return invalid argument", func() {
		_, err := s.client.CreateProject(ctxOrgAdminAuth, &shieldv1beta1.CreateProjectRequest{
			Body: &shieldv1beta1.ProjectRequestBody{
				Slug:  "new-project",
				OrgId: myOrg.GetId(),
			},
		})
		s.Assert().Equal(codes.InvalidArgument, status.Convert(err).Code())
	})

	s.Run("3. org admin create a new project with wrong org id should return invalid argument", func() {
		_, err := s.client.CreateProject(ctxOrgAdminAuth, &shieldv1beta1.CreateProjectRequest{
			Body: &shieldv1beta1.ProjectRequestBody{
				Name:  "new project",
				Slug:  "new-project",
				OrgId: "not-uuid",
			},
		})
		s.Assert().Equal(codes.InvalidArgument, status.Convert(err).Code())
	})

	s.Run("4. org admin create a new project with same name and org-id should conflict", func() {
		res, err := s.client.CreateProject(ctxOrgAdminAuth, &shieldv1beta1.CreateProjectRequest{
			Body: &shieldv1beta1.ProjectRequestBody{
				Name:  "new project",
				Slug:  "new-project",
				OrgId: myOrg.GetId(),
			},
		})
		s.Assert().NoError(err)
		newProject = res.GetProject()
		s.Assert().NotNil(newProject)

		_, err = s.client.CreateProject(ctxOrgAdminAuth, &shieldv1beta1.CreateProjectRequest{
			Body: &shieldv1beta1.ProjectRequestBody{
				Name:  "new project",
				Slug:  "new-project",
				OrgId: myOrg.GetId(),
			},
		})
		s.Assert().Equal(codes.AlreadyExists, status.Convert(err).Code())
	})

	s.Run("5. org admin update a new project with empty body should return invalid argument", func() {
		_, err := s.client.UpdateProject(ctxOrgAdminAuth, &shieldv1beta1.UpdateProjectRequest{
			Id:   newProject.GetId(),
			Body: nil,
		})
		s.Assert().Equal(codes.InvalidArgument, status.Convert(err).Code())
	})

	s.Run("6. org admin update a new project with empty group id should return not found", func() {
		_, err := s.client.UpdateProject(ctxOrgAdminAuth, &shieldv1beta1.UpdateProjectRequest{
			Id: "",
			Body: &shieldv1beta1.ProjectRequestBody{
				Name:  "new project",
				Slug:  "new-project",
				OrgId: myOrg.GetId(),
			},
		})
		s.Assert().Equal(codes.NotFound, status.Convert(err).Code())
	})

	s.Run("7. org admin update a new project with unknown project id and not uuid should return not found", func() {
		_, err := s.client.UpdateProject(ctxOrgAdminAuth, &shieldv1beta1.UpdateProjectRequest{
			Id: "random",
			Body: &shieldv1beta1.ProjectRequestBody{
				Name:  "new project",
				Slug:  "new-project",
				OrgId: myOrg.GetId(),
			},
		})
		s.Assert().Equal(codes.NotFound, status.Convert(err).Code())
	})

	s.Run("8. org admin update a new project with same name and org-id but different id should return conflict", func() {
		_, err := s.client.UpdateProject(ctxOrgAdminAuth, &shieldv1beta1.UpdateProjectRequest{
			Id: newProject.GetId(),
			Body: &shieldv1beta1.ProjectRequestBody{
				Name:  "project 1",
				Slug:  "project-1",
				OrgId: myOrg.GetId(),
			},
		})
		s.Assert().Equal(codes.AlreadyExists, status.Convert(err).Code())
	})
}

func (s *EndToEndAPIRegressionTestSuite) TestGroupAPI() {
	var newGroup *shieldv1beta1.Group

	ctxOrgAdminAuth := metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{
		testbench.IdentityHeader: testbench.OrgAdminEmail,
	}))

	// get my org
	res, err := s.client.ListOrganizations(context.Background(), &shieldv1beta1.ListOrganizationsRequest{})
	s.Require().NoError(err)
	s.Require().Greater(len(res.GetOrganizations()), 0)
	myOrg := res.GetOrganizations()[0]

	s.Run("1. org admin create a new team with empty auth email should return unauthenticated error", func() {
		_, err := s.client.CreateGroup(context.Background(), &shieldv1beta1.CreateGroupRequest{
			Body: &shieldv1beta1.GroupRequestBody{
				Name:  "new-group",
				Slug:  "new-group",
				OrgId: myOrg.GetId(),
			},
		})
		s.Assert().Equal(codes.Unauthenticated, status.Convert(err).Code())
	})

	s.Run("2. org admin create a new team with empty name should return invalid argument", func() {
		_, err := s.client.CreateGroup(ctxOrgAdminAuth, &shieldv1beta1.CreateGroupRequest{
			Body: &shieldv1beta1.GroupRequestBody{
				Slug:  "new-group",
				OrgId: myOrg.GetId(),
			},
		})
		s.Assert().Equal(codes.InvalidArgument, status.Convert(err).Code())
	})

	s.Run("3. org admin create a new team with wrong org id should return invalid argument", func() {
		_, err := s.client.CreateGroup(ctxOrgAdminAuth, &shieldv1beta1.CreateGroupRequest{
			Body: &shieldv1beta1.GroupRequestBody{
				Name:  "new group",
				Slug:  "new-group",
				OrgId: "not-uuid",
			},
		})
		s.Assert().Equal(codes.InvalidArgument, status.Convert(err).Code())
	})

	s.Run("4. org admin create a new team with same name and org-id should conflict", func() {
		res, err := s.client.CreateGroup(ctxOrgAdminAuth, &shieldv1beta1.CreateGroupRequest{
			Body: &shieldv1beta1.GroupRequestBody{
				Name:  "new group",
				Slug:  "new-group",
				OrgId: myOrg.GetId(),
			},
		})
		s.Assert().NoError(err)
		newGroup = res.GetGroup()
		s.Assert().NotNil(newGroup)

		_, err = s.client.CreateGroup(ctxOrgAdminAuth, &shieldv1beta1.CreateGroupRequest{
			Body: &shieldv1beta1.GroupRequestBody{
				Name:  "new group",
				Slug:  "new-group",
				OrgId: myOrg.GetId(),
			},
		})
		s.Assert().Equal(codes.AlreadyExists, status.Convert(err).Code())
	})

	s.Run("5. group admin update a new team with empty body should return invalid argument", func() {
		_, err := s.client.UpdateGroup(ctxOrgAdminAuth, &shieldv1beta1.UpdateGroupRequest{
			Id:   newGroup.GetId(),
			Body: nil,
		})
		s.Assert().Equal(codes.InvalidArgument, status.Convert(err).Code())
	})

	s.Run("6. group admin update a new team with empty group id should return not found", func() {
		_, err := s.client.UpdateGroup(ctxOrgAdminAuth, &shieldv1beta1.UpdateGroupRequest{
			Id: "",
			Body: &shieldv1beta1.GroupRequestBody{
				Name:  "new group",
				Slug:  "new-group",
				OrgId: myOrg.GetId(),
			},
		})
		s.Assert().Equal(codes.NotFound, status.Convert(err).Code())
	})

	s.Run("7. group admin update a new team with unknown group id and not uuid should return not found", func() {
		_, err := s.client.UpdateGroup(ctxOrgAdminAuth, &shieldv1beta1.UpdateGroupRequest{
			Id: "random",
			Body: &shieldv1beta1.GroupRequestBody{
				Name:  "new group",
				Slug:  "new-group",
				OrgId: myOrg.GetId(),
			},
		})
		s.Assert().Equal(codes.NotFound, status.Convert(err).Code())
	})

	s.Run("8. group admin update a new team with same name and org-id but different id should return conflict", func() {
		_, err := s.client.UpdateGroup(ctxOrgAdminAuth, &shieldv1beta1.UpdateGroupRequest{
			Id: newGroup.GetId(),
			Body: &shieldv1beta1.GroupRequestBody{
				Name:  "org1 group1",
				Slug:  "org1-group1",
				OrgId: myOrg.GetId(),
			},
		})
		s.Assert().Equal(codes.AlreadyExists, status.Convert(err).Code())
	})
}

func (s *EndToEndAPIRegressionTestSuite) TestUserAPI() {
	var newUser *shieldv1beta1.User

	ctxOrgAdminAuth := metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{
		testbench.IdentityHeader: testbench.OrgAdminEmail,
	}))

	s.Run("1. org admin create a new user with empty auth email should return unauthenticated error", func() {
		_, err := s.client.CreateUser(context.Background(), &shieldv1beta1.CreateUserRequest{
			Body: &shieldv1beta1.UserRequestBody{
				Name:  "new user a",
				Email: "new-user-a@gotocompany.com",
				Metadata: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"foo": structpb.NewBoolValue(true),
					},
				},
			},
		})
		s.Assert().Equal(codes.Unauthenticated, status.Convert(err).Code())
	})

	s.Run("2. org admin create a new user with unparsable metadata should return invalid argument error", func() {
		_, err := s.client.CreateUser(ctxOrgAdminAuth, &shieldv1beta1.CreateUserRequest{
			Body: &shieldv1beta1.UserRequestBody{
				Name:  "new user a",
				Email: "new-user-a@gotocompany.com",
				Metadata: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"foo": structpb.NewNullValue(),
					},
				},
			},
		})
		s.Assert().Equal(codes.InvalidArgument, status.Convert(err).Code())
	})

	s.Run("4. org admin create a new user with same email should return conflict error", func() {
		res, err := s.client.CreateUser(ctxOrgAdminAuth, &shieldv1beta1.CreateUserRequest{
			Body: &shieldv1beta1.UserRequestBody{
				Name:  "new user a",
				Email: "new-user-a@gotocompany.com",
				Metadata: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"foo": structpb.NewBoolValue(true),
					},
				},
			},
		})
		s.Assert().NoError(err)
		newUser = res.GetUser()

		_, err = s.client.CreateUser(ctxOrgAdminAuth, &shieldv1beta1.CreateUserRequest{
			Body: &shieldv1beta1.UserRequestBody{
				Name:  "new user a",
				Email: "new-user-a@gotocompany.com",
				Metadata: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"foo": structpb.NewBoolValue(true),
					},
				},
			},
		})
		s.Assert().Equal(codes.AlreadyExists, status.Convert(err).Code())
	})

	s.Run("5. org admin update user with conflicted detail should return conflict error", func() {
		_, err := s.client.UpdateUser(ctxOrgAdminAuth, &shieldv1beta1.UpdateUserRequest{
			Id: newUser.GetId(),
			Body: &shieldv1beta1.UserRequestBody{
				Name:  "new user a",
				Email: "admin1-group1-org1@gotocompany.com",
				Metadata: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"foo": structpb.NewBoolValue(true),
					},
				},
			},
		})
		s.Assert().Equal(codes.AlreadyExists, status.Convert(err).Code())
	})

	ctxCurrentUser := metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{
		testbench.IdentityHeader: newUser.GetEmail(),
	}))

	s.Run("6. update current user with empty email should return invalid argument error", func() {
		_, err := s.client.UpdateCurrentUser(ctxCurrentUser, &shieldv1beta1.UpdateCurrentUserRequest{
			Body: &shieldv1beta1.UserRequestBody{
				Name:  "new user a",
				Email: "",
				Metadata: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"foo": structpb.NewBoolValue(true),
					},
				},
			},
		})
		s.Assert().Equal(codes.InvalidArgument, status.Convert(err).Code())
	})

	s.Run("7. update current user with different email in header and body should return invalid argument error", func() {
		_, err := s.client.UpdateCurrentUser(ctxCurrentUser, &shieldv1beta1.UpdateCurrentUserRequest{
			Body: &shieldv1beta1.UserRequestBody{
				Name:  "new user a",
				Email: "admin1-group1-org1@gotocompany.com",
				Metadata: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"foo": structpb.NewBoolValue(true),
					},
				},
			},
		})
		s.Assert().Equal(codes.InvalidArgument, status.Convert(err).Code())
	})
}

func (s *EndToEndAPIRegressionTestSuite) TestRelationAPI() {
	ctxOrgAdminAuth := metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{
		testbench.IdentityHeader: testbench.OrgAdminEmail,
	}))

	// get my res
	res, err := s.client.ListResources(context.Background(), &shieldv1beta1.ListResourcesRequest{})
	s.Require().NoError(err)
	s.Require().Greater(len(res.GetResources()), 0)
	resource := res.GetResources()[1]

	s.Run("1. should return not found error when resource name is send as object id, when name is uuid", func() {
		_, err := s.client.CreateRelation(ctxOrgAdminAuth, &shieldv1beta1.CreateRelationRequest{
			Body: &shieldv1beta1.RelationRequestBody{
				ObjectId:        "47c412cf-d223-40ba-b8b3-895062980221", // appeal name
				ObjectNamespace: "guardian/appeal",
				Subject:         "shield/user:member2-group1@gotocompany.com",
				RoleName:        "owner",
			},
		})
		s.Assert().Equal(codes.NotFound, status.Convert(err).Code())
	})

	s.Run("2. should return not found error when resource name is send as object id, when name is not uuid", func() {
		_, err := s.client.CreateRelation(ctxOrgAdminAuth, &shieldv1beta1.CreateRelationRequest{
			Body: &shieldv1beta1.RelationRequestBody{
				ObjectId:        resource.GetName(),
				ObjectNamespace: resource.GetNamespace().GetId(),
				Subject:         "shield/user:member2-group1@gotocompany.com",
				RoleName:        "owner",
			},
		})
		s.Assert().Equal(codes.InvalidArgument, status.Convert(err).Code())
	})

	s.Run("3. should return success when object id is resource id", func() {
		_, err := s.client.CreateRelation(ctxOrgAdminAuth, &shieldv1beta1.CreateRelationRequest{
			Body: &shieldv1beta1.RelationRequestBody{
				ObjectId:        resource.GetId(),
				ObjectNamespace: resource.GetNamespace().GetId(),
				Subject:         "shield/user:member2-group1@gotocompany.com",
				RoleName:        "owner",
			},
		})
		s.Assert().Equal(codes.OK, status.Convert(err).Code())
	})
}

func (s *EndToEndAPIRegressionTestSuite) TestServiceDataAPI() {
	ctxOrgAdminAuth := metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{
		testbench.IdentityHeader: testbench.OrgAdminEmail,
	}))

	res, err := s.client.ListProjects(context.Background(), &shieldv1beta1.ListProjectsRequest{})
	s.Require().NoError(err)
	s.Require().Greater(len(res.GetProjects()), 0)
	myProject := res.GetProjects()[0]

	usr, err := s.client.ListUsers(ctxOrgAdminAuth, &shieldv1beta1.ListUsersRequest{})
	s.Require().NoError(err)
	s.Require().Greater(len(usr.GetUsers()), 0)
	myUser := usr.GetUsers()[0]

	s.Run("1. org admin create a new service data key with empty auth email should return unauthenticated error", func() {
		_, err := s.serviceDataClient.CreateServiceDataKey(context.Background(), &shieldv1beta1.CreateServiceDataKeyRequest{
			Body: &shieldv1beta1.ServiceDataKeyRequestBody{
				Project:     myProject.Id,
				Key:         "new-key-01",
				Description: "description for key 01",
			},
		})
		s.Assert().Equal(codes.Unauthenticated, status.Convert(err).Code())
	})

	s.Run("2. org admin create a new service data key with blank key should return invalid argument error", func() {
		_, err := s.serviceDataClient.CreateServiceDataKey(ctxOrgAdminAuth, &shieldv1beta1.CreateServiceDataKeyRequest{
			Body: &shieldv1beta1.ServiceDataKeyRequestBody{
				Project:     myProject.Id,
				Key:         "",
				Description: "description for key 01",
			},
		})
		s.Assert().Equal(codes.InvalidArgument, status.Convert(err).Code())
	})

	s.Run("3. org admin create a new service data key with invalid project id should return invalid argument error", func() {
		_, err := s.serviceDataClient.CreateServiceDataKey(ctxOrgAdminAuth, &shieldv1beta1.CreateServiceDataKeyRequest{
			Body: &shieldv1beta1.ServiceDataKeyRequestBody{
				Project:     "invalid-project-id",
				Key:         "new-key-01",
				Description: "description for key 01",
			},
		})
		s.Assert().Equal(codes.InvalidArgument, status.Convert(err).Code())
	})

	s.Run("4. org admin create a new service data key in same project with same name should conflict", func() {
		res, err := s.serviceDataClient.CreateServiceDataKey(ctxOrgAdminAuth, &shieldv1beta1.CreateServiceDataKeyRequest{
			Body: &shieldv1beta1.ServiceDataKeyRequestBody{
				Project:     myProject.Id,
				Key:         "new-key-01",
				Description: "description for key 01",
			},
		})
		s.Assert().NoError(err)
		newServiceDataKey := res.GetServiceDataKey()
		s.Assert().NotNil(newServiceDataKey)

		_, err = s.serviceDataClient.CreateServiceDataKey(ctxOrgAdminAuth, &shieldv1beta1.CreateServiceDataKeyRequest{
			Body: &shieldv1beta1.ServiceDataKeyRequestBody{
				Project:     myProject.Id,
				Key:         "new-key-01",
				Description: "description for key 01",
			},
		})
		s.Assert().Equal(codes.AlreadyExists, status.Convert(err).Code())
	})

	s.Run("5. org admin update a user service data with invalid user id should return invalid argument error", func() {
		_, err := s.serviceDataClient.UpsertUserServiceData(ctxOrgAdminAuth, &shieldv1beta1.UpsertUserServiceDataRequest{
			UserId: "invalid-user-id",
			Body: &shieldv1beta1.UpsertServiceDataRequestBody{
				Project: myProject.Id,
				Data: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"update-key": structpb.NewStringValue("update value"),
					},
				},
			},
		})
		s.Assert().Equal(codes.InvalidArgument, status.Convert(err).Code())
	})

	s.Run("6. org admin update a user service data with invalid project id should return invalid argument error", func() {
		_, err := s.serviceDataClient.UpsertUserServiceData(ctxOrgAdminAuth, &shieldv1beta1.UpsertUserServiceDataRequest{
			UserId: testbench.OrgAdminEmail,
			Body: &shieldv1beta1.UpsertServiceDataRequestBody{
				Project: "invalid-project-id",
				Data: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"update-key": structpb.NewStringValue("update value"),
					},
				},
			},
		})
		s.Assert().Equal(codes.InvalidArgument, status.Convert(err).Code())
	})

	s.Run("7. org admin update multiple service data return invalid argument error", func() {
		_, err := s.serviceDataClient.UpsertUserServiceData(ctxOrgAdminAuth, &shieldv1beta1.UpsertUserServiceDataRequest{
			UserId: testbench.OrgAdminEmail,
			Body: &shieldv1beta1.UpsertServiceDataRequestBody{
				Project: myProject.Id,
				Data: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"update-key-1": structpb.NewStringValue("update value-1"),
						"update-key-2": structpb.NewStringValue("update value-2"),
					},
				},
			},
		})
		s.Assert().Equal(codes.InvalidArgument, status.Convert(err).Code())
	})

	s.Run("8. update service data org admin created without edit permission should return unauthenticated error", func() {
		res, err := s.serviceDataClient.CreateServiceDataKey(ctxOrgAdminAuth, &shieldv1beta1.CreateServiceDataKeyRequest{
			Body: &shieldv1beta1.ServiceDataKeyRequestBody{
				Project:     myProject.Id,
				Key:         "new-key",
				Description: "description for key",
			},
		})
		s.Assert().NoError(err)
		newServiceDataKey := res.GetServiceDataKey()
		s.Assert().NotNil(newServiceDataKey)

		ctxTestUser := metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{
			testbench.IdentityHeader: myUser.Email,
		}))

		_, err = s.serviceDataClient.UpsertUserServiceData(ctxTestUser, &shieldv1beta1.UpsertUserServiceDataRequest{
			UserId: testbench.OrgAdminEmail,
			Body: &shieldv1beta1.UpsertServiceDataRequestBody{
				Project: myProject.Id,
				Data: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"new-key": structpb.NewStringValue("new-value"),
					},
				},
			},
		})
		s.Assert().Equal(codes.Unauthenticated, status.Convert(err).Code())
	})

	s.Run("9. org admin update a group service data with invalid group id should return invalid argument error", func() {
		_, err := s.serviceDataClient.UpsertGroupServiceData(ctxOrgAdminAuth, &shieldv1beta1.UpsertGroupServiceDataRequest{
			GroupId: "invalid-group-id",
			Body: &shieldv1beta1.UpsertServiceDataRequestBody{
				Project: myProject.Id,
				Data: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"update-key": structpb.NewStringValue("update value"),
					},
				},
			},
		})
		s.Assert().Equal(codes.InvalidArgument, status.Convert(err).Code())
	})
}

func TestEndToEndAPIRegressionTestSuite(t *testing.T) {
	suite.Run(t, new(EndToEndAPIRegressionTestSuite))
}
