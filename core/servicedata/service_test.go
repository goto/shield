package servicedata_test

import (
	"context"
	"errors"
	"testing"

	"github.com/goto/shield/core/project"
	"github.com/goto/shield/core/relation"
	"github.com/goto/shield/core/resource"
	"github.com/goto/shield/core/servicedata"
	"github.com/goto/shield/core/servicedata/mocks"
	"github.com/goto/shield/core/user"
	"github.com/goto/shield/internal/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_CreateKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		email   string
		key     servicedata.Key
		setup   func(t *testing.T) *servicedata.Service
		want    servicedata.Key
		wantErr error
	}{
		{
			name:  "CreateKey",
			email: "john.doe@gotocompany.com",
			key: servicedata.Key{
				ProjectID:   "test-project-slug",
				Key:         "test-key",
				Description: "test key no 01",
			},
			setup: func(t *testing.T) *servicedata.Service {
				t.Helper()
				repository := &mocks.Repository{}
				resourceService := &mocks.ResourceService{}
				relationService := &mocks.RelationService{}
				projectService := &mocks.ProjectService{}
				userService := &mocks.UserService{}
				repository.On("WithTransaction", mock.Anything).Return(context.TODO())
				repository.On("Commit", mock.Anything).Return(nil)
				userService.EXPECT().FetchCurrentUser(mock.Anything).
					Return(user.User{
						ID:    "test-user-id",
						Email: "john.doe@gotocompany.com",
					}, nil)
				projectService.EXPECT().Get(mock.Anything, "test-project-slug").
					Return(project.Project{
						ID: "test-project-id",
					}, nil)
				resourceService.EXPECT().Create(mock.Anything, resource.Resource{
					Name:        "test-project-id:servicedata_key:test-key",
					ProjectID:   "test-project-id",
					NamespaceID: "shield/servicedata_key",
					UserID:      "test-user-id",
				}).Return(resource.Resource{
					Idxa: "test-resource-id",
				}, nil)
				repository.EXPECT().CreateKey(mock.Anything, servicedata.Key{
					URN:         "test-project-id:servicedata_key:test-key",
					ProjectID:   "test-project-id",
					Key:         "test-key",
					Description: "test key no 01",
					ResourceID:  "test-resource-id",
				}).Return(servicedata.Key{
					URN:         "test-project-id:servicedata_key:test-key",
					ProjectID:   "test-project-id",
					Key:         "test-key",
					Description: "test key no 01",
					ResourceID:  "test-resource-id",
				}, nil)
				relationService.EXPECT().Create(mock.Anything, relation.RelationV2{
					Object: relation.Object{
						ID:          "test-resource-id",
						NamespaceID: schema.ServiceDataKeyNamespace,
					},
					Subject: relation.Subject{
						ID:        "test-user-id",
						RoleID:    schema.OwnerRole,
						Namespace: schema.UserPrincipal,
					},
				}).Return(relation.RelationV2{}, nil)
				return servicedata.NewService(repository, resourceService, relationService, projectService, userService)
			},
			want: servicedata.Key{
				URN:         "test-project-id:servicedata_key:test-key",
				ProjectID:   "test-project-id",
				Key:         "test-key",
				Description: "test key no 01",
				ResourceID:  "test-resource-id",
			},
		},
		{
			name: "CreateKeyEmpty",
			key: servicedata.Key{
				ProjectID:   "test-project-slug",
				Key:         "",
				Description: "test key no 01",
			},
			setup: func(t *testing.T) *servicedata.Service {
				t.Helper()
				repository := &mocks.Repository{}
				resourceService := &mocks.ResourceService{}
				relationService := &mocks.RelationService{}
				projectService := &mocks.ProjectService{}
				userService := &mocks.UserService{}
				return servicedata.NewService(repository, resourceService, relationService, projectService, userService)
			},
			wantErr: servicedata.ErrInvalidDetail,
		},
		{
			name: "CreateKeyMissingEmail",
			key: servicedata.Key{
				ProjectID:   "test-project-slug",
				Key:         "test-key",
				Description: "test key no 01",
			},
			setup: func(t *testing.T) *servicedata.Service {
				t.Helper()
				repository := &mocks.Repository{}
				resourceService := &mocks.ResourceService{}
				relationService := &mocks.RelationService{}
				projectService := &mocks.ProjectService{}
				userService := &mocks.UserService{}
				userService.EXPECT().FetchCurrentUser(mock.Anything).Return(user.User{}, user.ErrMissingEmail)
				return servicedata.NewService(repository, resourceService, relationService, projectService, userService)
			},
			wantErr: user.ErrMissingEmail,
		},
		{
			name: "CreateKeyInvalidEmail",
			key: servicedata.Key{
				ProjectID:   "test-project-slug",
				Key:         "test-key",
				Description: "test key no 01",
			},
			email: "jane.doe@gotocompany.com",
			setup: func(t *testing.T) *servicedata.Service {
				t.Helper()
				repository := &mocks.Repository{}
				resourceService := &mocks.ResourceService{}
				relationService := &mocks.RelationService{}
				projectService := &mocks.ProjectService{}
				userService := &mocks.UserService{}
				userService.EXPECT().FetchCurrentUser(mock.Anything).Return(user.User{}, user.ErrInvalidEmail)
				return servicedata.NewService(repository, resourceService, relationService, projectService, userService)
			},
			wantErr: user.ErrInvalidEmail,
		},
		{
			name: "CreateKeyInvalidProjectID",
			key: servicedata.Key{
				ProjectID:   "invalid-test-project-slug",
				Key:         "test-key",
				Description: "test key no 01",
			},
			email: "jane.doe@gotocompany.com",
			setup: func(t *testing.T) *servicedata.Service {
				t.Helper()
				repository := &mocks.Repository{}
				resourceService := &mocks.ResourceService{}
				relationService := &mocks.RelationService{}
				projectService := &mocks.ProjectService{}
				userService := &mocks.UserService{}
				userService.EXPECT().FetchCurrentUser(mock.Anything).Return(user.User{Email: "jane.doe@gotocompany.com"}, nil)
				projectService.EXPECT().Get(mock.Anything, "invalid-test-project-slug").Return(project.Project{}, project.ErrNotExist)
				return servicedata.NewService(repository, resourceService, relationService, projectService, userService)
			},
			wantErr: project.ErrNotExist,
		},
		{
			name: "CreateKeyErrCreateResource",
			key: servicedata.Key{
				ProjectID:   "test-project-slug",
				Key:         "test-key",
				Description: "test key no 01",
			},
			email: "john.doe@gotocompany.com",
			setup: func(t *testing.T) *servicedata.Service {
				t.Helper()
				repository := &mocks.Repository{}
				resourceService := &mocks.ResourceService{}
				relationService := &mocks.RelationService{}
				projectService := &mocks.ProjectService{}
				userService := &mocks.UserService{}
				repository.On("WithTransaction", mock.Anything).Return(context.TODO())
				repository.On("Rollback", mock.Anything, mock.Anything).Return(nil)
				userService.EXPECT().FetchCurrentUser(mock.Anything).
					Return(user.User{
						ID:    "test-user-id",
						Email: "john.doe@gotocompany.com",
					}, nil)
				projectService.EXPECT().Get(mock.Anything, "test-project-slug").
					Return(project.Project{
						ID: "test-project-id",
					}, nil)
				resourceService.EXPECT().Create(mock.Anything, resource.Resource{
					Name:        "test-project-id:servicedata_key:test-key",
					ProjectID:   "test-project-id",
					NamespaceID: "shield/servicedata_key",
					UserID:      "test-user-id",
				}).Return(resource.Resource{}, resource.ErrConflict)
				return servicedata.NewService(repository, resourceService, relationService, projectService, userService)
			},
			wantErr: resource.ErrConflict,
		},
		{
			name: "CreateKeyErrCreateKey",
			key: servicedata.Key{
				ProjectID:   "test-project-slug",
				Key:         "test-key",
				Description: "test key no 01",
			},
			email: "john.doe@gotocompany.com",
			setup: func(t *testing.T) *servicedata.Service {
				t.Helper()
				repository := &mocks.Repository{}
				resourceService := &mocks.ResourceService{}
				relationService := &mocks.RelationService{}
				projectService := &mocks.ProjectService{}
				userService := &mocks.UserService{}
				repository.On("WithTransaction", mock.Anything).Return(context.TODO())
				repository.On("Rollback", mock.Anything, mock.Anything).Return(nil)
				userService.EXPECT().FetchCurrentUser(mock.Anything).
					Return(user.User{
						ID:    "test-user-id",
						Email: "john.doe@gotocompany.com",
					}, nil)
				projectService.EXPECT().Get(mock.Anything, "test-project-slug").
					Return(project.Project{
						ID: "test-project-id",
					}, nil)
				resourceService.EXPECT().Create(mock.Anything, resource.Resource{
					Name:        "test-project-id:servicedata_key:test-key",
					ProjectID:   "test-project-id",
					NamespaceID: "shield/servicedata_key",
					UserID:      "test-user-id",
				}).Return(resource.Resource{
					Idxa: "test-resource-id",
				}, nil)
				repository.EXPECT().CreateKey(mock.Anything, servicedata.Key{
					URN:         "test-project-id:servicedata_key:test-key",
					ProjectID:   "test-project-id",
					Key:         "test-key",
					Description: "test key no 01",
					ResourceID:  "test-resource-id",
				}).Return(servicedata.Key{}, servicedata.ErrConflict)
				return servicedata.NewService(repository, resourceService, relationService, projectService, userService)
			},
			wantErr: servicedata.ErrConflict,
		},
		{
			name: "CreateKeyErrCreateRelation",
			key: servicedata.Key{
				ProjectID:   "test-project-slug",
				Key:         "test-key",
				Description: "test key no 01",
			},
			email: "john.doe@gotocompany.com",
			setup: func(t *testing.T) *servicedata.Service {
				t.Helper()
				repository := &mocks.Repository{}
				resourceService := &mocks.ResourceService{}
				relationService := &mocks.RelationService{}
				projectService := &mocks.ProjectService{}
				userService := &mocks.UserService{}
				repository.On("WithTransaction", mock.Anything).Return(context.TODO())
				repository.On("Rollback", mock.Anything, mock.Anything).Return(nil)
				userService.EXPECT().FetchCurrentUser(mock.Anything).
					Return(user.User{
						ID:    "test-user-id",
						Email: "john.doe@gotocompany.com",
					}, nil)
				projectService.EXPECT().Get(mock.Anything, "test-project-slug").
					Return(project.Project{
						ID: "test-project-id",
					}, nil)
				resourceService.EXPECT().Create(mock.Anything, resource.Resource{
					Name:        "test-project-id:servicedata_key:test-key",
					ProjectID:   "test-project-id",
					NamespaceID: "shield/servicedata_key",
					UserID:      "test-user-id",
				}).Return(resource.Resource{
					Idxa: "test-resource-id",
				}, nil)
				repository.EXPECT().CreateKey(mock.Anything, servicedata.Key{
					URN:         "test-project-id:servicedata_key:test-key",
					ProjectID:   "test-project-id",
					Key:         "test-key",
					Description: "test key no 01",
					ResourceID:  "test-resource-id",
				}).Return(servicedata.Key{
					URN:         "test-project-id:servicedata_key:test-key",
					ProjectID:   "test-project-id",
					Key:         "test-key",
					Description: "test key no 01",
					ResourceID:  "test-resource-id",
				}, nil)
				relationService.EXPECT().Create(mock.Anything, relation.RelationV2{
					Object: relation.Object{
						ID:          "test-resource-id",
						NamespaceID: schema.ServiceDataKeyNamespace,
					},
					Subject: relation.Subject{
						ID:        "test-user-id",
						RoleID:    schema.OwnerRole,
						Namespace: schema.UserPrincipal,
					},
				}).Return(relation.RelationV2{}, relation.ErrCreatingRelationInAuthzEngine)
				return servicedata.NewService(repository, resourceService, relationService, projectService, userService)
			},
			wantErr: relation.ErrCreatingRelationInAuthzEngine,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := tt.setup(t)

			assert.NotNil(t, svc)

			ctx := user.SetContextWithEmail(context.TODO(), tt.email)
			got, err := svc.CreateKey(ctx, tt.key)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.wantErr))
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
