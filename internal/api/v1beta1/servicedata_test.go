package v1beta1

import (
	"context"
	"testing"

	"github.com/goto/shield/core/project"
	"github.com/goto/shield/core/relation"
	"github.com/goto/shield/core/servicedata"
	"github.com/goto/shield/core/user"
	"github.com/goto/shield/internal/api/v1beta1/mocks"
	"github.com/goto/shield/pkg/uuid"
	shieldv1beta1 "github.com/goto/shield/proto/v1beta1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	testKeyProjectID  = uuid.NewString()
	testKeyID         = uuid.NewString()
	testKeyResourceID = uuid.NewString()
	testKey           = servicedata.Key{
		ID:          testKeyID,
		URN:         "key-urn",
		ProjectID:   testKeyProjectID,
		Key:         "test-key",
		Description: "test description",
		ResourceID:  testKeyResourceID,
	}
	testKeyPB = &shieldv1beta1.ServiceDataKey{
		Id:  testKey.ID,
		Urn: testKey.URN,
	}
)

func TestHandler_CreateKey(t *testing.T) {
	email := "user@gotocompany.com"
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
					Key:         testKey.Key,
					Description: testKey.Description,
				}).Return(servicedata.Key{}, project.ErrNotExist)
				return user.SetContextWithEmail(ctx, email)
			},
			request: &shieldv1beta1.CreateServiceDataKeyRequest{
				Body: &shieldv1beta1.ServiceDataKeyRequestBody{
					Project:     "non-existing-project",
					Key:         testKey.Key,
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
					Key:         testKey.Key,
					Description: testKey.Description,
				}).Return(servicedata.Key{}, servicedata.ErrConflict)
				return user.SetContextWithEmail(ctx, email)
			},
			request: &shieldv1beta1.CreateServiceDataKeyRequest{
				Body: &shieldv1beta1.ServiceDataKeyRequestBody{
					Project:     testKey.ProjectID,
					Key:         testKey.Key,
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
					Key:         testKey.Key,
					Description: testKey.Description,
				}).Return(servicedata.Key{}, relation.ErrInvalidDetail)
				return user.SetContextWithEmail(ctx, email)
			},
			request: &shieldv1beta1.CreateServiceDataKeyRequest{
				Body: &shieldv1beta1.ServiceDataKeyRequestBody{
					Project:     testKey.ProjectID,
					Key:         testKey.Key,
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
					Key:         testKey.Key,
					Description: testKey.Description,
				}).Return(servicedata.Key{}, servicedata.ErrInvalidDetail)
				return user.SetContextWithEmail(ctx, email)
			},
			request: &shieldv1beta1.CreateServiceDataKeyRequest{
				Body: &shieldv1beta1.ServiceDataKeyRequestBody{
					Project:     testKey.ProjectID,
					Key:         testKey.Key,
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
					Key:         testKey.Key,
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
					Key:         testKey.Key,
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
