package servicedata_test

import (
	"context"
	"errors"
	"testing"

	"github.com/goto/shield/core/action"
	"github.com/goto/shield/core/namespace"
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

var (
	testResourceID  = "test-resource-id"
	testUserID      = "test-user-id"
	testProjectID   = "test-project-id"
	testProjectSlug = "test-project-slug"
	testKey         = servicedata.Key{
		ProjectID:   "test-project-slug",
		Key:         "test-key",
		Description: "test key no 01",
	}
	testCreateKey = servicedata.Key{
		URN:         "test-project-slug:servicedata_key:test-key",
		ProjectID:   testProjectID,
		ProjectSlug: testProjectSlug,
		Key:         "test-key",
		Description: "test key no 01",
		ResourceID:  testResourceID,
	}
	testCreatedKey = servicedata.Key{
		URN:         "test-project-slug:servicedata_key:test-key",
		ProjectID:   testProjectID,
		Key:         "test-key",
		Description: "test key no 01",
		ResourceID:  testResourceID,
	}
	testResource = resource.Resource{
		Name:        "test-project-slug:servicedata_key:test-key",
		ProjectID:   testProjectID,
		NamespaceID: schema.ServiceDataKeyNamespace,
		UserID:      testUserID,
	}
	testRelation = relation.RelationV2{
		Object: relation.Object{
			ID:          testResourceID,
			NamespaceID: schema.ServiceDataKeyNamespace,
		},
		Subject: relation.Subject{
			ID:        testUserID,
			RoleID:    schema.OwnerRole,
			Namespace: schema.UserPrincipal,
		},
	}
	testEntityID    = "test-entity-id"
	testNamespaceID = "test-namespace-id"
	testValue       = "test-value"
	testServiceData = servicedata.ServiceData{
		EntityID:    testEntityID,
		NamespaceID: testNamespaceID,
		Key:         testCreateKey,
		Value:       testValue,
	}
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
			key:   testKey,
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
						ID:    testUserID,
						Email: "john.doe@gotocompany.com",
					}, nil)
				projectService.EXPECT().Get(mock.Anything, "test-project-slug").
					Return(project.Project{
						ID:   testProjectID,
						Slug: testProjectSlug,
					}, nil)
				resourceService.EXPECT().Create(mock.Anything, testResource).Return(resource.Resource{
					Idxa: testResourceID,
				}, nil)
				repository.EXPECT().CreateKey(mock.Anything, testCreateKey).Return(testCreatedKey, nil)
				relationService.EXPECT().Create(mock.Anything, testRelation).Return(relation.RelationV2{}, nil)
				return servicedata.NewService(repository, resourceService, relationService, projectService, userService)
			},
			want: testCreatedKey,
		},
		{
			name: "CreateKeyEmpty",
			key: servicedata.Key{
				ProjectID:   testKey.ProjectID,
				Key:         "",
				Description: testKey.Description,
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
			key:  testKey,
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
			name:  "CreateKeyInvalidEmail",
			key:   testKey,
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
				Key:         testKey.Key,
				Description: testKey.Description,
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
			name:  "CreateKeyErrCreateResource",
			key:   testKey,
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
						ID:    testUserID,
						Email: "john.doe@gotocompany.com",
					}, nil)
				projectService.EXPECT().Get(mock.Anything, "test-project-slug").
					Return(project.Project{
						ID:   testProjectID,
						Slug: testProjectSlug,
					}, nil)
				resourceService.EXPECT().Create(mock.Anything, testResource).Return(resource.Resource{}, resource.ErrConflict)
				return servicedata.NewService(repository, resourceService, relationService, projectService, userService)
			},
			wantErr: resource.ErrConflict,
		},
		{
			name:  "CreateKeyErrCreateKey",
			key:   testKey,
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
						ID:    testUserID,
						Email: "john.doe@gotocompany.com",
					}, nil)
				projectService.EXPECT().Get(mock.Anything, "test-project-slug").
					Return(project.Project{
						ID:   testProjectID,
						Slug: testProjectSlug,
					}, nil)
				resourceService.EXPECT().Create(mock.Anything, testResource).Return(resource.Resource{
					Idxa: testResourceID,
				}, nil)
				repository.EXPECT().CreateKey(mock.Anything, testCreateKey).Return(servicedata.Key{}, servicedata.ErrConflict)
				return servicedata.NewService(repository, resourceService, relationService, projectService, userService)
			},
			wantErr: servicedata.ErrConflict,
		},
		{
			name:  "CreateKeyErrCreateRelation",
			key:   testKey,
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
						ID:    testUserID,
						Email: "john.doe@gotocompany.com",
					}, nil)
				projectService.EXPECT().Get(mock.Anything, "test-project-slug").
					Return(project.Project{
						ID:   testProjectID,
						Slug: testProjectSlug,
					}, nil)
				resourceService.EXPECT().Create(mock.Anything, testResource).Return(resource.Resource{
					Idxa: testResourceID,
				}, nil)
				repository.EXPECT().CreateKey(mock.Anything, testCreateKey).Return(testCreatedKey, nil)
				relationService.EXPECT().Create(mock.Anything, testRelation).Return(relation.RelationV2{}, relation.ErrCreatingRelationInAuthzEngine)
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

func TestService_Upsert(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		email   string
		data    servicedata.ServiceData
		setup   func(t *testing.T) *servicedata.Service
		want    servicedata.ServiceData
		wantErr error
	}{
		{
			name:  "Upsert",
			email: "john.doe@gotocompany.com",
			data:  testServiceData,
			setup: func(t *testing.T) *servicedata.Service {
				t.Helper()
				repository := &mocks.Repository{}
				resourceService := &mocks.ResourceService{}
				relationService := &mocks.RelationService{}
				projectService := &mocks.ProjectService{}
				userService := &mocks.UserService{}
				userService.EXPECT().FetchCurrentUser(mock.Anything).
					Return(user.User{
						ID:    testUserID,
						Email: "john.doe@gotocompany.com",
					}, nil)
				projectService.EXPECT().Get(mock.Anything, testProjectID).
					Return(project.Project{
						ID:   testProjectID,
						Slug: testProjectSlug,
					}, nil)
				resourceService.EXPECT().GetByURN(mock.Anything, testCreateKey.URN).Return(resource.Resource{
					Idxa: testResourceID,
				}, nil)
				relationService.EXPECT().CheckPermission(mock.Anything, user.User{
					ID:    testUserID,
					Email: "john.doe@gotocompany.com",
				}, namespace.Namespace{ID: schema.ServiceDataKeyNamespace},
					testResourceID, action.Action{ID: "edit"}).Return(true, nil)
				repository.EXPECT().Upsert(mock.Anything, testServiceData).Return(testServiceData, nil)
				return servicedata.NewService(repository, resourceService, relationService, projectService, userService)
			},
			want: testServiceData,
		},
		{
			name: "UpsertKeyEmpty",
			data: servicedata.ServiceData{
				Key: servicedata.Key{
					Key: "",
				},
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
			name:  "UpsertInvalidEmail",
			data:  testServiceData,
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
			name: "UpsertInvalidProjectID",
			data: servicedata.ServiceData{
				Key: servicedata.Key{
					Key:       testKey.Key,
					ProjectID: "invalid-test-project-slug",
				},
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
			name:  "UpsertErrGetResource",
			data:  testServiceData,
			email: "john.doe@gotocompany.com",
			setup: func(t *testing.T) *servicedata.Service {
				t.Helper()
				repository := &mocks.Repository{}
				resourceService := &mocks.ResourceService{}
				relationService := &mocks.RelationService{}
				projectService := &mocks.ProjectService{}
				userService := &mocks.UserService{}
				userService.EXPECT().FetchCurrentUser(mock.Anything).
					Return(user.User{
						ID:    testUserID,
						Email: "john.doe@gotocompany.com",
					}, nil)
				projectService.EXPECT().Get(mock.Anything, testProjectID).
					Return(project.Project{
						ID:   testProjectID,
						Slug: testProjectSlug,
					}, nil)
				resourceService.EXPECT().GetByURN(mock.Anything, testServiceData.Key.URN).Return(resource.Resource{}, resource.ErrNotExist)
				return servicedata.NewService(repository, resourceService, relationService, projectService, userService)
			},
			wantErr: resource.ErrNotExist,
		},
		{
			name:  "UpsertErrUnauthenticated",
			email: "john.doe@gotocompany.com",
			data:  testServiceData,
			setup: func(t *testing.T) *servicedata.Service {
				t.Helper()
				repository := &mocks.Repository{}
				resourceService := &mocks.ResourceService{}
				relationService := &mocks.RelationService{}
				projectService := &mocks.ProjectService{}
				userService := &mocks.UserService{}
				userService.EXPECT().FetchCurrentUser(mock.Anything).
					Return(user.User{
						ID:    testUserID,
						Email: "john.doe@gotocompany.com",
					}, nil)
				projectService.EXPECT().Get(mock.Anything, testProjectID).
					Return(project.Project{
						ID:   testProjectID,
						Slug: testProjectSlug,
					}, nil)
				resourceService.EXPECT().GetByURN(mock.Anything, testCreateKey.URN).Return(resource.Resource{
					Idxa: testResourceID,
				}, nil)
				relationService.EXPECT().CheckPermission(mock.Anything, user.User{
					ID:    testUserID,
					Email: "john.doe@gotocompany.com",
				}, namespace.Namespace{ID: schema.ServiceDataKeyNamespace},
					testResourceID, action.Action{ID: "edit"}).Return(false, nil)
				return servicedata.NewService(repository, resourceService, relationService, projectService, userService)
			},
			wantErr: user.ErrInvalidEmail,
		},
		{
			name:  "UpsertErr",
			email: "john.doe@gotocompany.com",
			data:  testServiceData,
			setup: func(t *testing.T) *servicedata.Service {
				t.Helper()
				repository := &mocks.Repository{}
				resourceService := &mocks.ResourceService{}
				relationService := &mocks.RelationService{}
				projectService := &mocks.ProjectService{}
				userService := &mocks.UserService{}
				userService.EXPECT().FetchCurrentUser(mock.Anything).
					Return(user.User{
						ID:    testUserID,
						Email: "john.doe@gotocompany.com",
					}, nil)
				projectService.EXPECT().Get(mock.Anything, testProjectID).
					Return(project.Project{
						ID:   testProjectID,
						Slug: testProjectSlug,
					}, nil)
				resourceService.EXPECT().GetByURN(mock.Anything, testCreateKey.URN).Return(resource.Resource{
					Idxa: testResourceID,
				}, nil)
				relationService.EXPECT().CheckPermission(mock.Anything, user.User{
					ID:    testUserID,
					Email: "john.doe@gotocompany.com",
				}, namespace.Namespace{ID: schema.ServiceDataKeyNamespace},
					testResourceID, action.Action{ID: "edit"}).Return(true, nil)
				repository.EXPECT().Upsert(mock.Anything, testServiceData).Return(servicedata.ServiceData{}, servicedata.ErrInvalidDetail)
				return servicedata.NewService(repository, resourceService, relationService, projectService, userService)
			},
			wantErr: servicedata.ErrInvalidDetail,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := tt.setup(t)

			assert.NotNil(t, svc)

			ctx := user.SetContextWithEmail(context.TODO(), tt.email)
			got, err := svc.Upsert(ctx, tt.data)
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
