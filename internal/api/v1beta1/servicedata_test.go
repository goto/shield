package v1beta1

import (
	"context"
	"testing"

	"github.com/goto/shield/core/group"
	"github.com/goto/shield/core/project"
	"github.com/goto/shield/core/relation"
	"github.com/goto/shield/core/servicedata"
	"github.com/goto/shield/core/user"
	"github.com/goto/shield/internal/api/v1beta1/mocks"
	"github.com/goto/shield/pkg/uuid"
	shieldv1beta1 "github.com/goto/shield/proto/v1beta1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/structpb"
)

var (
	testKeyProjectID  = uuid.NewString()
	testKeyID         = uuid.NewString()
	testKeyResourceID = uuid.NewString()
	testKeyName       = "test-key"
	testValue         = "test-value"
	testEntityID      = uuid.NewString()
	testKey           = servicedata.Key{
		ID:          testKeyID,
		URN:         "key-urn",
		ProjectID:   testKeyProjectID,
		Name:        testKeyName,
		Description: "test description",
		ResourceID:  testKeyResourceID,
	}
	testKeyPB = &shieldv1beta1.ServiceDataKey{
		Id:  testKey.ID,
		Urn: testKey.URN,
	}
	testUserServiceDataCreate = servicedata.ServiceData{
		EntityID:    testEntityID,
		NamespaceID: userNamespaceID,
		Key: servicedata.Key{
			Name:      testKeyName,
			ProjectID: testKeyProjectID,
		},
		Value: testValue,
	}
	testGroupServiceDataCreate = servicedata.ServiceData{
		EntityID:    testEntityID,
		NamespaceID: groupNamespaceID,
		Key: servicedata.Key{
			Name:      testKeyName,
			ProjectID: testKeyProjectID,
		},
		Value: testValue,
	}
	email = "user@gotocompany.com"
)

func TestHandler_CreateKey(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(ctx context.Context, ss *mocks.ServiceDataService) context.Context
		request *shieldv1beta1.CreateServiceDataKeyRequest
		want    *shieldv1beta1.CreateServiceDataKeyResponse
		wantErr error
	}{
		{
			name: "should return unauthenticated error if auth email in context is empty",
			setup: func(ctx context.Context, ss *mocks.ServiceDataService) context.Context {
				ss.EXPECT().CreateKey(mock.AnythingOfType("context.todoCtx"), servicedata.Key{}).Return(servicedata.Key{}, user.ErrInvalidEmail)
				return ctx
			},
			request: &shieldv1beta1.CreateServiceDataKeyRequest{
				Body: &shieldv1beta1.ServiceDataKeyRequestBody{},
			},
			want:    nil,
			wantErr: grpcUnauthenticated,
		},
		{
			name:    "should return bad body error if no request body",
			request: &shieldv1beta1.CreateServiceDataKeyRequest{},
			want:    nil,
			wantErr: grpcBadBodyError,
		},
		{
			name: "should return bad body error if project not exist",
			setup: func(ctx context.Context, ss *mocks.ServiceDataService) context.Context {
				ss.EXPECT().CreateKey(mock.AnythingOfType("*context.valueCtx"), servicedata.Key{
					ProjectID:   "non-existing-project",
					Name:        testKey.Name,
					Description: testKey.Description,
				}).Return(servicedata.Key{}, project.ErrNotExist)
				return user.SetContextWithEmail(ctx, email)
			},
			request: &shieldv1beta1.CreateServiceDataKeyRequest{
				Body: &shieldv1beta1.ServiceDataKeyRequestBody{
					Project:     "non-existing-project",
					Key:         testKey.Name,
					Description: testKey.Description,
				},
			},
			want:    nil,
			wantErr: grpcBadBodyError,
		},
		{
			name: "should return conflict error if key already created",
			setup: func(ctx context.Context, ss *mocks.ServiceDataService) context.Context {
				ss.EXPECT().CreateKey(mock.AnythingOfType("*context.valueCtx"), servicedata.Key{
					ProjectID:   testKey.ProjectID,
					Name:        testKey.Name,
					Description: testKey.Description,
				}).Return(servicedata.Key{}, servicedata.ErrConflict)
				return user.SetContextWithEmail(ctx, email)
			},
			request: &shieldv1beta1.CreateServiceDataKeyRequest{
				Body: &shieldv1beta1.ServiceDataKeyRequestBody{
					Project:     testKey.ProjectID,
					Key:         testKey.Name,
					Description: testKey.Description,
				},
			},
			want:    nil,
			wantErr: grpcConflictError,
		},
		{
			name: "should return bad body error if error invalid detail in relation",
			setup: func(ctx context.Context, ss *mocks.ServiceDataService) context.Context {
				ss.EXPECT().CreateKey(mock.AnythingOfType("*context.valueCtx"), servicedata.Key{
					ProjectID:   testKey.ProjectID,
					Name:        testKey.Name,
					Description: testKey.Description,
				}).Return(servicedata.Key{}, relation.ErrInvalidDetail)
				return user.SetContextWithEmail(ctx, email)
			},
			request: &shieldv1beta1.CreateServiceDataKeyRequest{
				Body: &shieldv1beta1.ServiceDataKeyRequestBody{
					Project:     testKey.ProjectID,
					Key:         testKey.Name,
					Description: testKey.Description,
				},
			},
			want:    nil,
			wantErr: grpcBadBodyError,
		},
		{
			name: "should return bad body error if error invalid detail in service data",
			setup: func(ctx context.Context, ss *mocks.ServiceDataService) context.Context {
				ss.EXPECT().CreateKey(mock.AnythingOfType("*context.valueCtx"), servicedata.Key{
					ProjectID:   testKey.ProjectID,
					Name:        testKey.Name,
					Description: testKey.Description,
				}).Return(servicedata.Key{}, servicedata.ErrInvalidDetail)
				return user.SetContextWithEmail(ctx, email)
			},
			request: &shieldv1beta1.CreateServiceDataKeyRequest{
				Body: &shieldv1beta1.ServiceDataKeyRequestBody{
					Project:     testKey.ProjectID,
					Key:         testKey.Name,
					Description: testKey.Description,
				},
			},
			want:    nil,
			wantErr: grpcBadBodyError,
		},
		{
			name: "should return created key if no error",
			setup: func(ctx context.Context, ss *mocks.ServiceDataService) context.Context {
				ss.EXPECT().CreateKey(mock.AnythingOfType("*context.valueCtx"), servicedata.Key{
					ProjectID:   testKey.ProjectID,
					Name:        testKey.Name,
					Description: testKey.Description,
				}).Return(servicedata.Key{
					ID:  testKey.ID,
					URN: testKey.URN,
				}, nil)
				return user.SetContextWithEmail(ctx, email)
			},
			request: &shieldv1beta1.CreateServiceDataKeyRequest{
				Body: &shieldv1beta1.ServiceDataKeyRequestBody{
					Project:     testKey.ProjectID,
					Key:         testKey.Name,
					Description: testKey.Description,
				},
			},
			want: &shieldv1beta1.CreateServiceDataKeyResponse{
				ServiceDataKey: testKeyPB,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServiceDataService := new(mocks.ServiceDataService)
			ctx := context.TODO()
			if tt.setup != nil {
				ctx = tt.setup(ctx, mockServiceDataService)
			}
			mockDep := Handler{serviceDataService: mockServiceDataService}
			resp, err := mockDep.CreateServiceDataKey(ctx, tt.request)
			assert.EqualValues(t, tt.want, resp)
			assert.EqualValues(t, tt.wantErr, err)
		})
	}
}

func TestHandler_UpdateUserServiceData(t *testing.T) {
	email := "user@gotocompany.com"
	tests := []struct {
		name    string
		setup   func(ctx context.Context, ss *mocks.ServiceDataService, us *mocks.UserService) context.Context
		request *shieldv1beta1.UpsertUserServiceDataRequest
		want    *shieldv1beta1.UpsertUserServiceDataResponse
		wantErr error
	}{
		{
			name:    "should return bad body error if no request body",
			request: &shieldv1beta1.UpsertUserServiceDataRequest{},
			want:    nil,
			wantErr: grpcBadBodyError,
		},
		{
			name: "should return bad body error if no id in param",
			request: &shieldv1beta1.UpsertUserServiceDataRequest{
				UserId: "",
				Body:   &shieldv1beta1.UpsertServiceDataRequestBody{},
			},
			want:    nil,
			wantErr: grpcBadBodyError,
		},
		{
			name: "should return bad body error if request body data is empty",
			request: &shieldv1beta1.UpsertUserServiceDataRequest{
				UserId: "",
				Body: &shieldv1beta1.UpsertServiceDataRequestBody{
					Data: &structpb.Struct{},
				},
			},
			want:    nil,
			wantErr: grpcBadBodyError,
		},
		{
			name: "should return bad body error if request body data is more than max upsert (1)",
			request: &shieldv1beta1.UpsertUserServiceDataRequest{
				UserId: "",
				Body: &shieldv1beta1.UpsertServiceDataRequestBody{
					Data: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							"test-key-1": structpb.NewStringValue("test-value-1"),
							"test-key-2": structpb.NewStringValue("test-value-2"),
						},
					},
				},
			},
			want:    nil,
			wantErr: grpcBadBodyError,
		},
		{
			name: "should return bad body error if user id or email in param does not exist",
			request: &shieldv1beta1.UpsertUserServiceDataRequest{
				UserId: "",
				Body: &shieldv1beta1.UpsertServiceDataRequestBody{
					Data: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							testKeyName: structpb.NewStringValue(testValue),
						},
					},
				},
			},
			setup: func(ctx context.Context, ss *mocks.ServiceDataService, us *mocks.UserService) context.Context {
				us.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), "").Return(user.User{}, user.ErrInvalidEmail)
				return ctx
			},
			want:    nil,
			wantErr: grpcBadBodyError,
		},
		{
			name: "should return unauthenticated error if email in header invalid",
			request: &shieldv1beta1.UpsertUserServiceDataRequest{
				UserId: testEntityID,
				Body: &shieldv1beta1.UpsertServiceDataRequestBody{
					Project: testKeyProjectID,
					Data: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							testKeyName: structpb.NewStringValue(testValue),
						},
					},
				},
			},
			setup: func(ctx context.Context, ss *mocks.ServiceDataService, us *mocks.UserService) context.Context {
				us.EXPECT().Get(mock.AnythingOfType("*context.valueCtx"), testEntityID).Return(user.User{ID: testEntityID}, nil)
				ss.EXPECT().Upsert(mock.AnythingOfType("*context.valueCtx"), testUserServiceDataCreate).Return(servicedata.ServiceData{}, user.ErrInvalidEmail)
				return user.SetContextWithEmail(ctx, email)
			},
			want:    nil,
			wantErr: grpcUnauthenticated,
		},
		{
			name: "should return bad body error if project id or slug is invalid",
			request: &shieldv1beta1.UpsertUserServiceDataRequest{
				UserId: testEntityID,
				Body: &shieldv1beta1.UpsertServiceDataRequestBody{
					Project: testKeyProjectID,
					Data: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							testKeyName: structpb.NewStringValue(testValue),
						},
					},
				},
			},
			setup: func(ctx context.Context, ss *mocks.ServiceDataService, us *mocks.UserService) context.Context {
				us.EXPECT().Get(mock.AnythingOfType("*context.valueCtx"), testEntityID).Return(user.User{ID: testEntityID}, nil)
				ss.EXPECT().Upsert(mock.AnythingOfType("*context.valueCtx"), testUserServiceDataCreate).Return(servicedata.ServiceData{}, project.ErrNotExist)
				return user.SetContextWithEmail(ctx, email)
			},
			want:    nil,
			wantErr: grpcBadBodyError,
		},
		{
			name: "should return created service data urn",
			request: &shieldv1beta1.UpsertUserServiceDataRequest{
				UserId: testEntityID,
				Body: &shieldv1beta1.UpsertServiceDataRequestBody{
					Project: testKeyProjectID,
					Data: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							testKeyName: structpb.NewStringValue(testValue),
						},
					},
				},
			},
			setup: func(ctx context.Context, ss *mocks.ServiceDataService, us *mocks.UserService) context.Context {
				us.EXPECT().Get(mock.AnythingOfType("*context.valueCtx"), testEntityID).Return(user.User{ID: testEntityID}, nil)
				ss.EXPECT().Upsert(mock.AnythingOfType("*context.valueCtx"), testUserServiceDataCreate).Return(servicedata.ServiceData{
					Key:   servicedata.Key{Name: testKeyName},
					Value: testValue,
				}, nil)
				return user.SetContextWithEmail(ctx, email)
			},
			want: &shieldv1beta1.UpsertUserServiceDataResponse{
				Data: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						testKeyName: structpb.NewStringValue(testValue),
					},
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServiceDataService := new(mocks.ServiceDataService)
			mockUserService := new(mocks.UserService)
			ctx := context.TODO()
			if tt.setup != nil {
				ctx = tt.setup(ctx, mockServiceDataService, mockUserService)
			}
			mockDep := Handler{serviceDataService: mockServiceDataService, userService: mockUserService, serviceDataConfig: ServiceDataConfig{MaxUpsert: 1}}
			resp, err := mockDep.UpsertUserServiceData(ctx, tt.request)
			assert.EqualValues(t, tt.want, resp)
			assert.EqualValues(t, tt.wantErr, err)
		})
	}
}

func TestHandler_UpdateGroupServiceData(t *testing.T) {
	email := "user@gotocompany.com"
	tests := []struct {
		name    string
		setup   func(ctx context.Context, ss *mocks.ServiceDataService, gs *mocks.GroupService) context.Context
		request *shieldv1beta1.UpsertGroupServiceDataRequest
		want    *shieldv1beta1.UpsertGroupServiceDataResponse
		wantErr error
	}{
		{
			name:    "should return bad body error if no request body",
			request: &shieldv1beta1.UpsertGroupServiceDataRequest{},
			want:    nil,
			wantErr: grpcBadBodyError,
		},
		{
			name: "should return bad body error if no id in param",
			request: &shieldv1beta1.UpsertGroupServiceDataRequest{
				GroupId: "",
				Body:    &shieldv1beta1.UpsertServiceDataRequestBody{},
			},
			want:    nil,
			wantErr: grpcBadBodyError,
		},
		{
			name: "should return bad body error if request body data is empty",
			request: &shieldv1beta1.UpsertGroupServiceDataRequest{
				GroupId: "",
				Body: &shieldv1beta1.UpsertServiceDataRequestBody{
					Data: &structpb.Struct{},
				},
			},
			want:    nil,
			wantErr: grpcBadBodyError,
		},
		{
			name: "should return bad body error if request body data is more than max upsert (1)",
			request: &shieldv1beta1.UpsertGroupServiceDataRequest{
				GroupId: "",
				Body: &shieldv1beta1.UpsertServiceDataRequestBody{
					Data: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							"test-key-1": structpb.NewStringValue("test-value-1"),
							"test-key-2": structpb.NewStringValue("test-value-2"),
						},
					},
				},
			},
			want:    nil,
			wantErr: grpcBadBodyError,
		},
		{
			name: "should return bad body error if group id or slug in param does not exist",
			request: &shieldv1beta1.UpsertGroupServiceDataRequest{
				GroupId: testEntityID,
				Body: &shieldv1beta1.UpsertServiceDataRequestBody{
					Data: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							testKeyName: structpb.NewStringValue(testValue),
						},
					},
				},
			},
			setup: func(ctx context.Context, ss *mocks.ServiceDataService, gs *mocks.GroupService) context.Context {
				gs.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), testEntityID).Return(group.Group{}, group.ErrNotExist)
				return ctx
			},
			want:    nil,
			wantErr: grpcBadBodyError,
		},
		{
			name: "should return unauthenticated error if email in header invalid",
			request: &shieldv1beta1.UpsertGroupServiceDataRequest{
				GroupId: testEntityID,
				Body: &shieldv1beta1.UpsertServiceDataRequestBody{
					Project: testKeyProjectID,
					Data: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							testKeyName: structpb.NewStringValue(testValue),
						},
					},
				},
			},
			setup: func(ctx context.Context, ss *mocks.ServiceDataService, gs *mocks.GroupService) context.Context {
				gs.EXPECT().Get(mock.AnythingOfType("*context.valueCtx"), testEntityID).Return(group.Group{ID: testEntityID}, nil)
				ss.EXPECT().Upsert(mock.AnythingOfType("*context.valueCtx"), testGroupServiceDataCreate).Return(servicedata.ServiceData{}, user.ErrInvalidEmail)
				return user.SetContextWithEmail(ctx, email)
			},
			want:    nil,
			wantErr: grpcUnauthenticated,
		},
		{
			name: "should return bad body error if project id or slug is invalid",
			request: &shieldv1beta1.UpsertGroupServiceDataRequest{
				GroupId: testEntityID,
				Body: &shieldv1beta1.UpsertServiceDataRequestBody{
					Project: testKeyProjectID,
					Data: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							testKeyName: structpb.NewStringValue(testValue),
						},
					},
				},
			},
			setup: func(ctx context.Context, ss *mocks.ServiceDataService, gs *mocks.GroupService) context.Context {
				gs.EXPECT().Get(mock.AnythingOfType("*context.valueCtx"), testEntityID).Return(group.Group{ID: testEntityID}, nil)
				ss.EXPECT().Upsert(mock.AnythingOfType("*context.valueCtx"), testGroupServiceDataCreate).Return(servicedata.ServiceData{}, project.ErrNotExist)
				return user.SetContextWithEmail(ctx, email)
			},
			want:    nil,
			wantErr: grpcBadBodyError,
		},
		{
			name: "should return created service data urn",
			request: &shieldv1beta1.UpsertGroupServiceDataRequest{
				GroupId: testEntityID,
				Body: &shieldv1beta1.UpsertServiceDataRequestBody{
					Project: testKeyProjectID,
					Data: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							testKeyName: structpb.NewStringValue(testValue),
						},
					},
				},
			},
			setup: func(ctx context.Context, ss *mocks.ServiceDataService, gs *mocks.GroupService) context.Context {
				gs.EXPECT().Get(mock.AnythingOfType("*context.valueCtx"), testEntityID).Return(group.Group{ID: testEntityID}, nil)
				ss.EXPECT().Upsert(mock.AnythingOfType("*context.valueCtx"), testGroupServiceDataCreate).Return(servicedata.ServiceData{
					Key:   servicedata.Key{Name: testKeyName},
					Value: testValue,
				}, nil)
				return user.SetContextWithEmail(ctx, email)
			},
			want: &shieldv1beta1.UpsertGroupServiceDataResponse{
				Data: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						testKeyName: structpb.NewStringValue(testValue),
					},
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServiceDataService := new(mocks.ServiceDataService)
			mockGroupService := new(mocks.GroupService)
			ctx := context.TODO()
			if tt.setup != nil {
				ctx = tt.setup(ctx, mockServiceDataService, mockGroupService)
			}
			mockDep := Handler{serviceDataService: mockServiceDataService, groupService: mockGroupService, serviceDataConfig: ServiceDataConfig{MaxUpsert: 1}}
			resp, err := mockDep.UpsertGroupServiceData(ctx, tt.request)
			assert.EqualValues(t, tt.want, resp)
			assert.EqualValues(t, tt.wantErr, err)
		})
	}
}

func TestHandler_GetUserServiceData(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(ctx context.Context, ss *mocks.ServiceDataService, us *mocks.UserService) context.Context
		request *shieldv1beta1.GetUserServiceDataRequest
		want    *shieldv1beta1.GetUserServiceDataResponse
		wantErr error
	}{
		{
			name: "should return unauthenticated error if auth email in context is empty",
			setup: func(ctx context.Context, ss *mocks.ServiceDataService, us *mocks.UserService) context.Context {
				us.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), testEntityID).Return(user.User{ID: testEntityID}, nil)
				ss.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), mock.AnythingOfType("servicedata.Filter")).Return([]servicedata.ServiceData{}, user.ErrMissingEmail)
				return ctx
			},
			request: &shieldv1beta1.GetUserServiceDataRequest{
				UserId: testEntityID,
			},
			want:    nil,
			wantErr: grpcUnauthenticated,
		},
		{
			name: "should return bad body error if user entity not exist",
			setup: func(ctx context.Context, ss *mocks.ServiceDataService, us *mocks.UserService) context.Context {
				us.EXPECT().Get(mock.AnythingOfType("*context.valueCtx"), testEntityID).Return(user.User{}, user.ErrNotExist)
				return user.SetContextWithEmail(ctx, email)
			},
			request: &shieldv1beta1.GetUserServiceDataRequest{
				UserId: testEntityID,
			},
			want:    nil,
			wantErr: grpcBadBodyError,
		},
		{
			name: "should return bad body error if project not exist",
			setup: func(ctx context.Context, ss *mocks.ServiceDataService, us *mocks.UserService) context.Context {
				us.EXPECT().Get(mock.AnythingOfType("*context.valueCtx"), testEntityID).Return(user.User{ID: testEntityID}, nil)
				ss.EXPECT().Get(mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("servicedata.Filter")).Return([]servicedata.ServiceData{}, project.ErrNotExist)
				return user.SetContextWithEmail(ctx, email)
			},
			request: &shieldv1beta1.GetUserServiceDataRequest{
				UserId:  testEntityID,
				Project: "invalid-project-id",
			},
			want:    nil,
			wantErr: grpcBadBodyError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServiceDataService := new(mocks.ServiceDataService)
			mockUserService := new(mocks.UserService)

			ctx := context.TODO()
			if tt.setup != nil {
				ctx = tt.setup(ctx, mockServiceDataService, mockUserService)
			}
			mockDep := Handler{serviceDataService: mockServiceDataService, userService: mockUserService}
			resp, err := mockDep.GetUserServiceData(ctx, tt.request)
			assert.EqualValues(t, tt.want, resp)
			assert.EqualValues(t, tt.wantErr, err)
		})
	}
}

func TestHandler_GetGroupServiceData(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(ctx context.Context, ss *mocks.ServiceDataService, gs *mocks.GroupService) context.Context
		request *shieldv1beta1.GetGroupServiceDataRequest
		want    *shieldv1beta1.GetGroupServiceDataResponse
		wantErr error
	}{
		{
			name: "should return unauthenticated error if auth email in context is empty",
			setup: func(ctx context.Context, ss *mocks.ServiceDataService, us *mocks.GroupService) context.Context {
				us.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), testEntityID).Return(group.Group{ID: testEntityID}, nil)
				ss.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), mock.AnythingOfType("servicedata.Filter")).Return([]servicedata.ServiceData{}, user.ErrMissingEmail)
				return ctx
			},
			request: &shieldv1beta1.GetGroupServiceDataRequest{
				GroupId: testEntityID,
			},
			want:    nil,
			wantErr: grpcUnauthenticated,
		},
		{
			name: "should return bad body error if user entity not exist",
			setup: func(ctx context.Context, ss *mocks.ServiceDataService, gs *mocks.GroupService) context.Context {
				gs.EXPECT().Get(mock.AnythingOfType("*context.valueCtx"), testEntityID).Return(group.Group{}, group.ErrNotExist)
				return user.SetContextWithEmail(ctx, email)
			},
			request: &shieldv1beta1.GetGroupServiceDataRequest{
				GroupId: testEntityID,
			},
			want:    nil,
			wantErr: grpcBadBodyError,
		},
		{
			name: "should return bad body error if project not exist",
			setup: func(ctx context.Context, ss *mocks.ServiceDataService, gs *mocks.GroupService) context.Context {
				gs.EXPECT().Get(mock.AnythingOfType("*context.valueCtx"), testEntityID).Return(group.Group{ID: testEntityID}, nil)
				ss.EXPECT().Get(mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("servicedata.Filter")).Return([]servicedata.ServiceData{}, project.ErrNotExist)
				return user.SetContextWithEmail(ctx, email)
			},
			request: &shieldv1beta1.GetGroupServiceDataRequest{
				GroupId: testEntityID,
				Project: "invalid-project-id",
			},
			want:    nil,
			wantErr: grpcBadBodyError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServiceDataService := new(mocks.ServiceDataService)
			mockGroupService := new(mocks.GroupService)

			ctx := context.TODO()
			if tt.setup != nil {
				ctx = tt.setup(ctx, mockServiceDataService, mockGroupService)
			}
			mockDep := Handler{serviceDataService: mockServiceDataService, groupService: mockGroupService}
			resp, err := mockDep.GetGroupServiceData(ctx, tt.request)
			assert.EqualValues(t, tt.want, resp)
			assert.EqualValues(t, tt.wantErr, err)
		})
	}
}
