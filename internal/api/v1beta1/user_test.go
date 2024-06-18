package v1beta1

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/goto/shield/core/group"
	"github.com/goto/shield/core/project"
	"github.com/goto/shield/core/servicedata"
	"github.com/goto/shield/core/user"
	"github.com/goto/shield/internal/api/v1beta1/mocks"
	"github.com/goto/shield/internal/schema"
	errorsPkg "github.com/goto/shield/pkg/errors"
	"github.com/goto/shield/pkg/metadata"
	"github.com/goto/shield/pkg/uuid"
	"golang.org/x/exp/maps"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	shieldv1beta1 "github.com/goto/shield/proto/v1beta1"
)

var (
	testUserID  = "9f256f86-31a3-11ec-8d3d-0242ac130003"
	testUserMap = map[string]user.User{
		"9f256f86-31a3-11ec-8d3d-0242ac130003": {
			ID:    "9f256f86-31a3-11ec-8d3d-0242ac130003",
			Name:  "User 1",
			Email: "test@test.com",
			Metadata: metadata.Metadata{
				"foo":    "bar",
				"age":    21,
				"intern": true,
			},
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		},
	}
)

func TestListUsers(t *testing.T) {
	table := []struct {
		title string
		setup func(us *mocks.UserService, rs *mocks.RelationService, ps *mocks.ProjectService)
		req   *shieldv1beta1.ListUsersRequest
		want  *shieldv1beta1.ListUsersResponse
		err   error
	}{
		{
			title: "should return internal error in if user service return some error",
			setup: func(us *mocks.UserService, rs *mocks.RelationService, ps *mocks.ProjectService) {
				us.EXPECT().FetchCurrentUser(mock.AnythingOfType("context.todoCtx")).Return(user.User{
					ID: "083a77a2-ab14-40d2-a06d-f6d9f80c6378",
				}, nil)

				rs.EXPECT().LookupResources(mock.AnythingOfType("context.todoCtx"), schema.ServiceDataKeyNamespace, schema.ViewPermission, schema.UserPrincipal, "083a77a2-ab14-40d2-a06d-f6d9f80c6378").Return([]string{}, nil)
				ps.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), "").Return(project.Project{
					ID: "11a58737-f366-4d05-b925-6f7bded29257",
				}, nil)
				us.EXPECT().List(mock.Anything, mock.Anything).Return(user.PagedUsers{}, errors.New("some error"))
			},
			req: &shieldv1beta1.ListUsersRequest{
				PageSize: 50,
				PageNum:  1,
				Keyword:  "",
			},
			want: nil,
			err:  status.Errorf(codes.Internal, ErrInternalServer.Error()),
		}, {
			title: "should return all users if user service return all users",
			setup: func(us *mocks.UserService, rs *mocks.RelationService, ps *mocks.ProjectService) {
				var testUserList []user.User
				for _, u := range testUserMap {
					testUserList = append(testUserList, u)
				}
				us.EXPECT().FetchCurrentUser(mock.AnythingOfType("context.todoCtx")).Return(user.User{
					ID: "083a77a2-ab14-40d2-a06d-f6d9f80c6378",
				}, nil)

				rs.EXPECT().LookupResources(mock.AnythingOfType("context.todoCtx"), schema.ServiceDataKeyNamespace, schema.ViewPermission, schema.UserPrincipal, "083a77a2-ab14-40d2-a06d-f6d9f80c6378").Return([]string{}, nil)
				ps.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), "").Return(project.Project{
					ID: "11a58737-f366-4d05-b925-6f7bded29257",
				}, nil)
				us.EXPECT().List(mock.Anything, mock.Anything).Return(
					user.PagedUsers{
						Users: testUserList,
						Count: int32(len(testUserList)),
					}, nil)
			},
			req: &shieldv1beta1.ListUsersRequest{
				PageSize: 50,
				PageNum:  1,
				Keyword:  "",
			},
			want: &shieldv1beta1.ListUsersResponse{
				Count: 1,
				Users: []*shieldv1beta1.User{
					{
						Id:    "9f256f86-31a3-11ec-8d3d-0242ac130003",
						Name:  "User 1",
						Email: "test@test.com",
						Metadata: &structpb.Struct{
							Fields: map[string]*structpb.Value{
								"foo":    structpb.NewStringValue("bar"),
								"age":    structpb.NewNumberValue(21),
								"intern": structpb.NewBoolValue(true),
							},
						},
						CreatedAt: timestamppb.New(time.Time{}),
						UpdatedAt: timestamppb.New(time.Time{}),
					},
				},
			},
			err: nil,
		},
	}

	for _, tt := range table {
		t.Run(tt.title, func(t *testing.T) {
			mockUserSrv := new(mocks.UserService)
			mockRelationSrv := new(mocks.RelationService)
			mockProjectSvc := new(mocks.ProjectService)
			if tt.setup != nil {
				tt.setup(mockUserSrv, mockRelationSrv, mockProjectSvc)
			}
			mockDep := Handler{
				userService:     mockUserSrv,
				relationService: mockRelationSrv,
				projectService:  mockProjectSvc,
			}
			req := tt.req
			resp, err := mockDep.ListUsers(context.TODO(), req)
			assert.EqualValues(t, resp, tt.want)
			assert.EqualValues(t, err, tt.err)
		})
	}
}

func TestCreateUser(t *testing.T) {
	email := "user@gotocompany.com"
	table := []struct {
		title string
		setup func(ctx context.Context, us *mocks.UserService, sds *mocks.ServiceDataService) context.Context
		req   *shieldv1beta1.CreateUserRequest
		want  *shieldv1beta1.CreateUserResponse
		err   error
	}{
		{
			title: "should return unauthenticated error if no auth email header in context",
			req: &shieldv1beta1.CreateUserRequest{Body: &shieldv1beta1.UserRequestBody{
				Name:     "some user",
				Email:    "abc@test.com",
				Metadata: &structpb.Struct{},
			}},
			want: nil,
			err:  grpcUnauthenticated,
		},
		{
			title: "should return bad request error if metadata is not parsable",
			setup: func(ctx context.Context, us *mocks.UserService, sds *mocks.ServiceDataService) context.Context {
				return user.SetContextWithEmail(ctx, email)
			},
			req: &shieldv1beta1.CreateUserRequest{Body: &shieldv1beta1.UserRequestBody{
				Name:  "some user",
				Email: "abc@test.com",
				Metadata: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"foo": structpb.NewNullValue(),
					},
				},
			}},
			want: nil,
			err:  grpcBadBodyError,
		},
		{
			title: "should return bad request error if email is empty",
			setup: func(ctx context.Context, us *mocks.UserService, sds *mocks.ServiceDataService) context.Context {
				us.EXPECT().Create(mock.AnythingOfType("*context.valueCtx"), user.User{
					Name: "some user",
				}).Return(user.User{}, user.ErrInvalidEmail)
				return user.SetContextWithEmail(ctx, email)
			},
			req: &shieldv1beta1.CreateUserRequest{Body: &shieldv1beta1.UserRequestBody{
				Name:  "some user",
				Email: "",
				Metadata: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"foo": structpb.NewNullValue(),
					},
				},
			}},
			want: nil,
			err:  grpcBadBodyError,
		},
		{
			title: "should return invalid email error if email is invalid",
			setup: func(ctx context.Context, us *mocks.UserService, sds *mocks.ServiceDataService) context.Context {
				return user.SetContextWithEmail(ctx, email)
			},
			req: &shieldv1beta1.CreateUserRequest{Body: &shieldv1beta1.UserRequestBody{
				Name:  "some user",
				Email: "invalid email",
				Metadata: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"foo": structpb.NewNullValue(),
					},
				},
			}},
			want: nil,
			err:  user.ErrInvalidEmail,
		},
		{
			title: "should return already exist error if user service return error conflict",
			setup: func(ctx context.Context, us *mocks.UserService, sds *mocks.ServiceDataService) context.Context {
				ctx = user.SetContextWithEmail(ctx, email)
				us.EXPECT().Create(ctx, user.User{
					Name:     "some user",
					Email:    "abc@test.com",
					Metadata: nil,
				}).Return(user.User{}, user.ErrConflict)
				return ctx
			},
			req: &shieldv1beta1.CreateUserRequest{Body: &shieldv1beta1.UserRequestBody{
				Name:     "some user",
				Email:    "abc@test.com",
				Metadata: &structpb.Struct{},
			}},
			want: nil,
			err:  grpcConflictError,
		},
		{
			title: "should return success if user email contain whitespace but still valid service return nil error",
			setup: func(ctx context.Context, us *mocks.UserService, sds *mocks.ServiceDataService) context.Context {
				ctx = user.SetContextWithEmail(ctx, email)
				us.EXPECT().Create(mock.AnythingOfType("*context.valueCtx"), user.User{
					Name:     "some user",
					Email:    "abc@test.com",
					Metadata: nil,
				}).Return(
					user.User{
						ID:       "new-abc",
						Name:     "some user",
						Email:    "abc@test.com",
						Metadata: nil,
					}, nil)

				sds.EXPECT().Upsert(ctx, servicedata.ServiceData{
					EntityID:    "new-abc",
					NamespaceID: schema.UserPrincipal,
					Key: servicedata.Key{
						Name:      "foo",
						ProjectID: "",
					},
					Value: "bar",
				}).Return(servicedata.ServiceData{
					Key: servicedata.Key{
						Name: "foo",
					},
					Value: "bar",
				}, nil)
				return ctx
			},
			req: &shieldv1beta1.CreateUserRequest{Body: &shieldv1beta1.UserRequestBody{
				Name:  "some user",
				Email: "  abc@test.com  ",
				Metadata: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"foo": structpb.NewStringValue("bar"),
					},
				},
			}},
			want: &shieldv1beta1.CreateUserResponse{User: &shieldv1beta1.User{
				Id:    "new-abc",
				Name:  "some user",
				Email: "abc@test.com",
				Metadata: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"foo": structpb.NewStringValue("bar"),
					},
				},
				CreatedAt: timestamppb.New(time.Time{}),
				UpdatedAt: timestamppb.New(time.Time{}),
			}},
			err: nil,
		},
		{
			title: "should return success if user service return nil error",
			setup: func(ctx context.Context, us *mocks.UserService, sds *mocks.ServiceDataService) context.Context {
				ctx = user.SetContextWithEmail(ctx, email)
				us.EXPECT().Create(mock.AnythingOfType("*context.valueCtx"), user.User{
					Name:     "some user",
					Email:    "abc@test.com",
					Metadata: nil,
				}).Return(
					user.User{
						ID:       "new-abc",
						Name:     "some user",
						Email:    "abc@test.com",
						Metadata: nil,
					}, nil)

				sds.EXPECT().Upsert(ctx, servicedata.ServiceData{
					EntityID:    "new-abc",
					NamespaceID: schema.UserPrincipal,
					Key: servicedata.Key{
						Name:      "foo",
						ProjectID: "",
					},
					Value: "bar",
				}).Return(servicedata.ServiceData{
					Key: servicedata.Key{
						Name: "foo",
					},
					Value: "bar",
				}, nil)
				return ctx
			},
			req: &shieldv1beta1.CreateUserRequest{Body: &shieldv1beta1.UserRequestBody{
				Name:  "some user",
				Email: "abc@test.com",
				Metadata: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"foo": structpb.NewStringValue("bar"),
					},
				},
			}},
			want: &shieldv1beta1.CreateUserResponse{User: &shieldv1beta1.User{
				Id:    "new-abc",
				Name:  "some user",
				Email: "abc@test.com",
				Metadata: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"foo": structpb.NewStringValue("bar"),
					},
				},
				CreatedAt: timestamppb.New(time.Time{}),
				UpdatedAt: timestamppb.New(time.Time{}),
			}},
			err: nil,
		},
	}

	for _, tt := range table {
		t.Run(tt.title, func(t *testing.T) {
			var resp *shieldv1beta1.CreateUserResponse
			var err error

			ctx := context.TODO()
			mockUserSrv := new(mocks.UserService)
			mockServiceDataSvc := new(mocks.ServiceDataService)
			if tt.setup != nil {
				ctx = tt.setup(ctx, mockUserSrv, mockServiceDataSvc)
			}
			mockDep := Handler{
				userService:        mockUserSrv,
				serviceDataService: mockServiceDataSvc,
			}
			resp, err = mockDep.CreateUser(ctx, tt.req)
			assert.EqualValues(t, tt.want, resp)
			assert.EqualValues(t, tt.err, err)
		})
	}
}

func TestGetUser(t *testing.T) {
	randomID := uuid.NewString()
	table := []struct {
		title string
		req   *shieldv1beta1.GetUserRequest
		setup func(us *mocks.UserService, sd *mocks.ServiceDataService)
		want  *shieldv1beta1.GetUserResponse
		err   error
	}{
		{
			title: "should return not found error if user does not exist",
			setup: func(us *mocks.UserService, sd *mocks.ServiceDataService) {
				us.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), randomID).Return(user.User{}, user.ErrNotExist)
			},
			req: &shieldv1beta1.GetUserRequest{
				Id: randomID,
			},
			want: nil,
			err:  grpcUserNotFoundError,
		},
		{
			title: "should return not found error if user id is not uuid",
			setup: func(us *mocks.UserService, sd *mocks.ServiceDataService) {
				us.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), "some-id").Return(user.User{}, user.ErrInvalidUUID)
			},
			req: &shieldv1beta1.GetUserRequest{
				Id: "some-id",
			},
			want: nil,
			err:  grpcUserNotFoundError,
		},
		{
			title: "should return not found error if user id is invalid",
			setup: func(us *mocks.UserService, sd *mocks.ServiceDataService) {
				us.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), "").Return(user.User{}, user.ErrInvalidID)
			},
			req:  &shieldv1beta1.GetUserRequest{},
			want: nil,
			err:  grpcUserNotFoundError,
		},
		{
			title: "should return user if user service return nil error",
			setup: func(us *mocks.UserService, sd *mocks.ServiceDataService) {
				us.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), randomID).Return(
					user.User{
						ID:        randomID,
						Name:      "some user",
						Email:     "someuser@test.com",
						CreatedAt: time.Time{},
						UpdatedAt: time.Time{},
					}, nil)

				sd.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), servicedata.Filter{
					ID:        randomID,
					Namespace: userNamespaceID,
					Entities: maps.Values(map[string]string{
						"user": userNamespaceID,
					}),
				}).Return([]servicedata.ServiceData{
					{
						Key: servicedata.Key{
							Name: "foo",
						},
						Value: "bar",
					},
				}, nil)
			},
			req: &shieldv1beta1.GetUserRequest{
				Id: randomID,
			},
			want: &shieldv1beta1.GetUserResponse{User: &shieldv1beta1.User{
				Id:    randomID,
				Name:  "some user",
				Email: "someuser@test.com",
				Metadata: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"foo": structpb.NewStringValue("bar"),
					},
				},
				CreatedAt: timestamppb.New(time.Time{}),
				UpdatedAt: timestamppb.New(time.Time{}),
			}},
			err: nil,
		},
	}

	for _, tt := range table {
		t.Run(tt.title, func(t *testing.T) {
			mockUserSrv := new(mocks.UserService)
			mockServiceDataSrv := new(mocks.ServiceDataService)

			if tt.setup != nil {
				tt.setup(mockUserSrv, mockServiceDataSrv)
			}
			mockDep := Handler{userService: mockUserSrv, serviceDataService: mockServiceDataSrv}
			resp, err := mockDep.GetUser(context.TODO(), tt.req)
			assert.EqualValues(t, resp, tt.want)
			assert.EqualValues(t, err, tt.err)
		})
	}
}

func TestGetCurrentUser(t *testing.T) {
	email := "user@gotocompany.com"
	table := []struct {
		title  string
		setup  func(ctx context.Context, us *mocks.UserService) context.Context
		header string
		want   *shieldv1beta1.GetCurrentUserResponse
		err    error
	}{
		{
			title: "should return unauthenticated error if no auth email header in context",
			want:  nil,
			err:   grpcUnauthenticated,
		},
		{
			title: "should return not found error if user does not exist",
			setup: func(ctx context.Context, us *mocks.UserService) context.Context {
				us.EXPECT().GetByEmail(mock.AnythingOfType("*context.valueCtx"), email).Return(user.User{}, user.ErrNotExist)
				return user.SetContextWithEmail(ctx, email)
			},
			want: nil,
			err:  grpcUserNotFoundError,
		},
		{
			title: "should return error if user service return some error",
			setup: func(ctx context.Context, us *mocks.UserService) context.Context {
				us.EXPECT().GetByEmail(mock.AnythingOfType("*context.valueCtx"), email).Return(user.User{}, errors.New("some error"))
				return user.SetContextWithEmail(ctx, email)
			},
			want: nil,
			err:  grpcInternalServerError,
		},
		{
			title: "should return user if user service return nil error",
			setup: func(ctx context.Context, us *mocks.UserService) context.Context {
				us.EXPECT().GetByEmail(mock.AnythingOfType("*context.valueCtx"), email).Return(
					user.User{
						ID:        "user-id-1",
						Name:      "some user",
						Email:     "someuser@test.com",
						CreatedAt: time.Time{},
						UpdatedAt: time.Time{},
					}, nil)
				return user.SetContextWithEmail(ctx, email)
			},
			want: &shieldv1beta1.GetCurrentUserResponse{User: &shieldv1beta1.User{
				Id:    "user-id-1",
				Name:  "some user",
				Email: "someuser@test.com",
				Metadata: &structpb.Struct{
					Fields: map[string]*structpb.Value{},
				},
				CreatedAt: timestamppb.New(time.Time{}),
				UpdatedAt: timestamppb.New(time.Time{}),
			}},
			err: nil,
		},
	}

	for _, tt := range table {
		t.Run(tt.title, func(t *testing.T) {
			mockUserSrv := new(mocks.UserService)
			ctx := context.TODO()
			if tt.setup != nil {
				ctx = tt.setup(ctx, mockUserSrv)
			}
			mockDep := Handler{userService: mockUserSrv}
			resp, err := mockDep.GetCurrentUser(ctx, nil)
			assert.EqualValues(t, resp, tt.want)
			assert.EqualValues(t, err, tt.err)
		})
	}
}

func TestUpdateUser(t *testing.T) {
	someID := uuid.NewString()
	table := []struct {
		title  string
		setup  func(us *mocks.UserService, sds *mocks.ServiceDataService)
		req    *shieldv1beta1.UpdateUserRequest
		header string
		want   *shieldv1beta1.UpdateUserResponse
		err    error
	}{
		{
			title: "should return internal error if user service return some error",
			setup: func(us *mocks.UserService, sds *mocks.ServiceDataService) {
				us.EXPECT().UpdateByID(mock.AnythingOfType("context.todoCtx"), user.User{
					ID:    someID,
					Name:  "abc user",
					Email: "user@gotocompany.com",
				}).Return(user.User{}, errors.New("some error"))
			},
			req: &shieldv1beta1.UpdateUserRequest{
				Id: someID,
				Body: &shieldv1beta1.UserRequestBody{
					Name:  "abc user",
					Email: "user@gotocompany.com",
					Metadata: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							"foo": structpb.NewStringValue("bar"),
						},
					},
				},
			},
			want: nil,
			err:  grpcInternalServerError,
		},
		{
			title: "should return not found error if id is invalid",
			setup: func(us *mocks.UserService, sds *mocks.ServiceDataService) {
				us.EXPECT().UpdateByID(mock.AnythingOfType("context.todoCtx"), user.User{
					Name:  "abc user",
					Email: "user@gotocompany.com",
					Metadata: metadata.Metadata{
						"foo": "bar",
					},
				}).Return(user.User{}, user.ErrInvalidID)
			},
			req: &shieldv1beta1.UpdateUserRequest{
				Body: &shieldv1beta1.UserRequestBody{
					Name:  "abc user",
					Email: "user@gotocompany.com",
					Metadata: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							"foo": structpb.NewStringValue("bar"),
						},
					},
				},
			},
			want: nil,
			err:  grpcUserNotFoundError,
		},
		{
			title: "should return already exist error if user service return error conflict",
			setup: func(us *mocks.UserService, sds *mocks.ServiceDataService) {
				us.EXPECT().UpdateByID(mock.AnythingOfType("context.todoCtx"), user.User{
					ID:    someID,
					Name:  "abc user",
					Email: "user@gotocompany.com",
				}).Return(user.User{}, user.ErrConflict)
			},
			req: &shieldv1beta1.UpdateUserRequest{
				Id: someID,
				Body: &shieldv1beta1.UserRequestBody{
					Name:  "abc user",
					Email: "user@gotocompany.com",
					Metadata: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							"foo": structpb.NewStringValue("bar"),
						},
					},
				},
			},
			want: nil,
			err:  grpcConflictError,
		},
		{
			title: "should return bad request error if email in request empty",
			setup: func(us *mocks.UserService, sds *mocks.ServiceDataService) {
				us.EXPECT().UpdateByID(mock.AnythingOfType("context.todoCtx"), user.User{
					ID:   someID,
					Name: "abc user",
					Metadata: metadata.Metadata{
						"foo": "bar",
					},
				}).Return(user.User{}, user.ErrInvalidEmail)
			},
			req: &shieldv1beta1.UpdateUserRequest{
				Id: someID,
				Body: &shieldv1beta1.UserRequestBody{
					Name: "abc user",
					Metadata: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							"foo": structpb.NewStringValue("bar"),
						},
					},
				},
			},
			want: nil,
			err:  grpcBadBodyError,
		},
		{
			title: "should return invalid email error if email is invalid",
			setup: func(us *mocks.UserService, sds *mocks.ServiceDataService) {
			},
			req: &shieldv1beta1.UpdateUserRequest{
				Id: someID,
				Body: &shieldv1beta1.UserRequestBody{
					Name:  "abc user",
					Email: "invalid email",
					Metadata: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							"foo": structpb.NewStringValue("bar"),
						},
					},
				},
			},
			want: nil,
			err:  user.ErrInvalidEmail,
		},
		{
			title: "should return bad request error if empty request body",
			req:   &shieldv1beta1.UpdateUserRequest{Id: someID, Body: nil},
			want:  nil,
			err:   grpcBadBodyError,
		},
		{
			title: "should return error if servicedata service return error",
			setup: func(us *mocks.UserService, sds *mocks.ServiceDataService) {
				us.EXPECT().UpdateByID(mock.AnythingOfType("context.todoCtx"), user.User{
					ID:    someID,
					Name:  "abc user",
					Email: "user@gotocompany.com",
				}).Return(
					user.User{
						ID:        someID,
						Name:      "abc user",
						Email:     "user@gotocompany.com",
						CreatedAt: time.Time{},
						UpdatedAt: time.Time{},
					}, nil)

				sds.EXPECT().Upsert(mock.AnythingOfType("context.todoCtx"), servicedata.ServiceData{
					EntityID:    someID,
					NamespaceID: userNamespaceID,
					Key: servicedata.Key{
						Name:      "foo",
						ProjectID: "system",
					},
					Value: "bar",
				}).Return(servicedata.ServiceData{}, errorsPkg.ErrForbidden)
			},
			req: &shieldv1beta1.UpdateUserRequest{
				Id: someID,
				Body: &shieldv1beta1.UserRequestBody{
					Name:  "abc user",
					Email: "user@gotocompany.com",
					Metadata: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							"foo": structpb.NewStringValue("bar"),
						},
					},
				},
			},
			want: nil,
			err:  status.Error(codes.PermissionDenied, "you are not authorized to update foo key"),
		},
		{
			title: "should be successful if user and servicedata service return nil error",
			setup: func(us *mocks.UserService, sds *mocks.ServiceDataService) {
				us.EXPECT().UpdateByID(mock.AnythingOfType("context.todoCtx"), user.User{
					ID:    someID,
					Name:  "abc user",
					Email: "user@gotocompany.com",
				}).Return(
					user.User{
						ID:    someID,
						Name:  "abc user",
						Email: "user@gotocompany.com",
						Metadata: metadata.Metadata{
							"foo": "bar",
						},
						CreatedAt: time.Time{},
						UpdatedAt: time.Time{},
					}, nil)

				sds.EXPECT().Upsert(mock.AnythingOfType("context.todoCtx"), servicedata.ServiceData{
					EntityID:    someID,
					NamespaceID: userNamespaceID,
					Key: servicedata.Key{
						Name:      "foo",
						ProjectID: "system",
					},
					Value: "bar",
				}).Return(servicedata.ServiceData{
					Key: servicedata.Key{
						Name: "foo",
					},
					Value: "bar",
				}, nil)
			},
			req: &shieldv1beta1.UpdateUserRequest{
				Id: someID,
				Body: &shieldv1beta1.UserRequestBody{
					Name:  "abc user",
					Email: "user@gotocompany.com",
					Metadata: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							"foo": structpb.NewStringValue("bar"),
						},
					},
				},
			},
			want: &shieldv1beta1.UpdateUserResponse{User: &shieldv1beta1.User{
				Id:    someID,
				Name:  "abc user",
				Email: "user@gotocompany.com",
				Metadata: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"foo": structpb.NewStringValue("bar"),
					},
				},
				CreatedAt: timestamppb.New(time.Time{}),
				UpdatedAt: timestamppb.New(time.Time{}),
			}},
			err: nil,
		},
		{
			title: "should return success even though name is empty",
			setup: func(us *mocks.UserService, sds *mocks.ServiceDataService) {
				us.EXPECT().UpdateByID(mock.AnythingOfType("context.todoCtx"), user.User{
					ID:    someID,
					Email: "user@gotocompany.com",
				}).Return(
					user.User{
						ID:    someID,
						Email: "user@gotocompany.com",
						Metadata: metadata.Metadata{
							"foo": "bar",
						},
						CreatedAt: time.Time{},
						UpdatedAt: time.Time{},
					}, nil)

				sds.EXPECT().Upsert(mock.AnythingOfType("context.todoCtx"), servicedata.ServiceData{
					EntityID:    someID,
					NamespaceID: userNamespaceID,
					Key: servicedata.Key{
						Name:      "foo",
						ProjectID: "system",
					},
					Value: "bar",
				}).Return(servicedata.ServiceData{
					Key: servicedata.Key{
						Name: "foo",
					},
					Value: "bar",
				}, nil)
			},
			req: &shieldv1beta1.UpdateUserRequest{
				Id: someID,
				Body: &shieldv1beta1.UserRequestBody{
					Email: "user@gotocompany.com",
					Metadata: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							"foo": structpb.NewStringValue("bar"),
						},
					},
				},
			},
			want: &shieldv1beta1.UpdateUserResponse{User: &shieldv1beta1.User{
				Id:    someID,
				Email: "user@gotocompany.com",
				Metadata: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"foo": structpb.NewStringValue("bar"),
					},
				},
				CreatedAt: timestamppb.New(time.Time{}),
				UpdatedAt: timestamppb.New(time.Time{}),
			}},
			err: nil,
		},
	}

	for _, tt := range table {
		t.Run(tt.title, func(t *testing.T) {
			mockUserSrv := new(mocks.UserService)
			mockServiceDataSrv := new(mocks.ServiceDataService)
			ctx := context.TODO()
			if tt.setup != nil {
				tt.setup(mockUserSrv, mockServiceDataSrv)
			}
			mockDep := Handler{userService: mockUserSrv, serviceDataService: mockServiceDataSrv, serviceDataConfig: ServiceDataConfig{
				DefaultServiceDataProject: "system",
			}}
			resp, err := mockDep.UpdateUser(ctx, tt.req)
			assert.EqualValues(t, resp, tt.want)
			assert.EqualValues(t, tt.err, err)
		})
	}
}

func TestUpdateCurrentUser(t *testing.T) {
	email := "user@gotocompany.com"
	table := []struct {
		title  string
		setup  func(ctx context.Context, us *mocks.UserService) context.Context
		req    *shieldv1beta1.UpdateCurrentUserRequest
		header string
		want   *shieldv1beta1.UpdateCurrentUserResponse
		err    error
	}{
		{
			title: "should return unauthenticated error if auth email header not exist",
			req: &shieldv1beta1.UpdateCurrentUserRequest{Body: &shieldv1beta1.UserRequestBody{
				Name:  "abc user",
				Email: "abcuser123@test.com",
				Metadata: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"foo": structpb.NewStringValue("bar"),
					},
				},
			}},
			want: nil,
			err:  grpcUnauthenticated,
		},
		{
			title: "should return internal error if user service return some error",
			setup: func(ctx context.Context, us *mocks.UserService) context.Context {
				us.EXPECT().UpdateByEmail(mock.AnythingOfType("*context.valueCtx"), user.User{
					Name:  "abc user",
					Email: "user@gotocompany.com",
					Metadata: metadata.Metadata{
						"foo": "bar",
					},
				}).Return(user.User{}, errors.New("some error"))
				return user.SetContextWithEmail(ctx, email)
			},
			req: &shieldv1beta1.UpdateCurrentUserRequest{Body: &shieldv1beta1.UserRequestBody{
				Name:  "abc user",
				Email: "user@gotocompany.com",
				Metadata: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"foo": structpb.NewStringValue("bar"),
					},
				},
			}},
			want: nil,
			err:  grpcInternalServerError,
		},
		{
			title: "should return not found error if user service return err not exist",
			setup: func(ctx context.Context, us *mocks.UserService) context.Context {
				us.EXPECT().UpdateByEmail(mock.AnythingOfType("*context.valueCtx"), user.User{
					Name:  "abc user",
					Email: "user@gotocompany.com",
					Metadata: metadata.Metadata{
						"foo": "bar",
					},
				}).Return(user.User{}, user.ErrNotExist)
				return user.SetContextWithEmail(ctx, email)
			},
			req: &shieldv1beta1.UpdateCurrentUserRequest{Body: &shieldv1beta1.UserRequestBody{
				Name:  "abc user",
				Email: "user@gotocompany.com",
				Metadata: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"foo": structpb.NewStringValue("bar"),
					},
				},
			}},
			want: nil,
			err:  grpcUserNotFoundError,
		},
		{
			title: "should return bad request error if diff emails in header and body",
			setup: func(ctx context.Context, us *mocks.UserService) context.Context {
				return user.SetContextWithEmail(ctx, email)
			},
			req: &shieldv1beta1.UpdateCurrentUserRequest{Body: &shieldv1beta1.UserRequestBody{
				Name:  "abc user",
				Email: "abcuser123@test.com",
				Metadata: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"foo": structpb.NewStringValue("bar"),
					},
				},
			}},
			want: nil,
			err:  grpcBadBodyError,
		},
		{
			title: "should return bad request error if empty request body",
			setup: func(ctx context.Context, us *mocks.UserService) context.Context {
				return user.SetContextWithEmail(ctx, email)
			},
			req:  &shieldv1beta1.UpdateCurrentUserRequest{Body: nil},
			want: nil,
			err:  grpcBadBodyError,
		},
		{
			title: "should return success if user service return nil error",
			setup: func(ctx context.Context, us *mocks.UserService) context.Context {
				us.EXPECT().UpdateByEmail(mock.Anything, mock.Anything).Return(
					user.User{
						ID:    "user-id-1",
						Name:  "abc user",
						Email: "user@gotocompany.com",
						Metadata: metadata.Metadata{
							"foo": "bar",
						},
						CreatedAt: time.Time{},
						UpdatedAt: time.Time{},
					}, nil)
				return user.SetContextWithEmail(ctx, email)
			},
			req: &shieldv1beta1.UpdateCurrentUserRequest{Body: &shieldv1beta1.UserRequestBody{
				Name:  "abc user",
				Email: "user@gotocompany.com",
				Metadata: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"foo": structpb.NewStringValue("bar"),
					},
				},
			}},
			want: &shieldv1beta1.UpdateCurrentUserResponse{User: &shieldv1beta1.User{
				Id:    "user-id-1",
				Name:  "abc user",
				Email: "user@gotocompany.com",
				Metadata: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"foo": structpb.NewStringValue("bar"),
					},
				},
				CreatedAt: timestamppb.New(time.Time{}),
				UpdatedAt: timestamppb.New(time.Time{}),
			}},
			err: nil,
		},
	}

	for _, tt := range table {
		t.Run(tt.title, func(t *testing.T) {
			mockUserSrv := new(mocks.UserService)
			ctx := context.TODO()
			if tt.setup != nil {
				ctx = tt.setup(ctx, mockUserSrv)
			}
			mockDep := Handler{userService: mockUserSrv}
			resp, err := mockDep.UpdateCurrentUser(ctx, tt.req)
			assert.EqualValues(t, resp, tt.want)
			assert.EqualValues(t, err, tt.err)
		})
	}
}

func TestHandler_ListUserGroups(t *testing.T) {
	someUserID := uuid.NewString()
	someRoleID := "role-id"
	tests := []struct {
		name    string
		setup   func(gs *mocks.GroupService)
		request *shieldv1beta1.ListUserGroupsRequest
		want    *shieldv1beta1.ListUserGroupsResponse
		wantErr error
	}{
		{
			name: "should return internal error if group service return some error",
			setup: func(gs *mocks.GroupService) {
				gs.EXPECT().ListUserGroups(mock.AnythingOfType("context.todoCtx"), someUserID, someRoleID).Return([]group.Group{}, errors.New("some error"))
			},
			request: &shieldv1beta1.ListUserGroupsRequest{
				Id:   someUserID,
				Role: someRoleID,
			},
			want:    nil,
			wantErr: grpcInternalServerError,
		},
		{
			name: "should return empty list if user does not exist",
			setup: func(gs *mocks.GroupService) {
				gs.EXPECT().ListUserGroups(mock.AnythingOfType("context.todoCtx"), someUserID, someRoleID).Return([]group.Group{}, nil)
			},
			request: &shieldv1beta1.ListUserGroupsRequest{
				Id:   someUserID,
				Role: someRoleID,
			},
			want:    &shieldv1beta1.ListUserGroupsResponse{},
			wantErr: nil,
		},
		// if user id empty, it would not go to this handler
		{
			name: "should return groups if group service return not nil",
			setup: func(gs *mocks.GroupService) {
				var testGroupList []group.Group
				for _, g := range testGroupMap {
					testGroupList = append(testGroupList, g)
				}
				gs.EXPECT().ListUserGroups(mock.AnythingOfType("context.todoCtx"), someUserID, someRoleID).Return(testGroupList, nil)
			},
			request: &shieldv1beta1.ListUserGroupsRequest{
				Id:   someUserID,
				Role: someRoleID,
			},
			want: &shieldv1beta1.ListUserGroupsResponse{
				Groups: []*shieldv1beta1.Group{
					{
						Id:    "9f256f86-31a3-11ec-8d3d-0242ac130003",
						Name:  "Group 1",
						Slug:  "group-1",
						OrgId: "9f256f86-31a3-11ec-8d3d-0242ac130003",
						Metadata: &structpb.Struct{
							Fields: map[string]*structpb.Value{
								"foo": structpb.NewStringValue("bar"),
							},
						},
						CreatedAt: timestamppb.New(time.Time{}),
						UpdatedAt: timestamppb.New(time.Time{}),
					},
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGroupSvc := new(mocks.GroupService)
			if tt.setup != nil {
				tt.setup(mockGroupSvc)
			}
			h := Handler{
				groupService: mockGroupSvc,
			}
			got, err := h.ListUserGroups(context.TODO(), tt.request)
			assert.EqualValues(t, got, tt.want)
			assert.EqualValues(t, err, tt.wantErr)
		})
	}
}

func TestCreateMetadataKey(t *testing.T) {
	email := "user@gotocompany.com"
	table := []struct {
		title string
		setup func(ctx context.Context, us *mocks.UserService) context.Context
		req   *shieldv1beta1.CreateMetadataKeyRequest
		want  *shieldv1beta1.CreateMetadataKeyResponse
		err   error
	}{
		{
			title: "should return error if body is empty",
			setup: func(ctx context.Context, us *mocks.UserService) context.Context {
				us.EXPECT().CreateMetadataKey(mock.AnythingOfType("*context.valueCtx"), shieldv1beta1.CreateMetadataKeyRequest{Body: nil}).Return(user.UserMetadataKey{}, grpcBadBodyError)
				return user.SetContextWithEmail(ctx, email)
			},
			req:  &shieldv1beta1.CreateMetadataKeyRequest{Body: nil},
			want: nil,
			err:  grpcBadBodyError,
		},
		{
			title: "should return error conflict if key already exists",
			setup: func(ctx context.Context, us *mocks.UserService) context.Context {
				us.EXPECT().CreateMetadataKey(mock.AnythingOfType("*context.valueCtx"), user.UserMetadataKey{
					Key:         "k1",
					Description: "key one",
				}).Return(user.UserMetadataKey{}, user.ErrConflict)
				return user.SetContextWithEmail(ctx, email)
			},
			req: &shieldv1beta1.CreateMetadataKeyRequest{Body: &shieldv1beta1.MetadataKeyRequestBody{
				Key:         "k1",
				Description: "key one",
			}},
			want: nil,
			err:  grpcConflictError,
		},
	}

	for _, tt := range table {
		t.Run(tt.title, func(t *testing.T) {
			var resp *shieldv1beta1.CreateMetadataKeyResponse
			var err error

			ctx := context.TODO()
			mockUserSrv := new(mocks.UserService)
			if tt.setup != nil {
				ctx = tt.setup(ctx, mockUserSrv)
			}
			mockDep := Handler{userService: mockUserSrv}
			resp, err = mockDep.CreateMetadataKey(ctx, tt.req)
			assert.EqualValues(t, tt.want, resp)
			assert.EqualValues(t, tt.err, err)
		})
	}
}
