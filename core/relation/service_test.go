package relation_test

import (
	"context"
	"errors"
	"testing"

	"github.com/goto/shield/core/action"
	"github.com/goto/shield/core/activity"
	"github.com/goto/shield/core/namespace"
	"github.com/goto/shield/core/relation"
	"github.com/goto/shield/core/relation/mocks"
	"github.com/goto/shield/core/user"
	"github.com/goto/shield/internal/schema"
	"github.com/goto/shield/pkg/logger"
	"github.com/goto/shield/pkg/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	testLogger = logger.InitLogger(logger.Config{
		Level:  "info",
		Format: "json",
	})
	testRelation = relation.Relation{
		ObjectID:         uuid.NewString(),
		ObjectNamespace:  namespace.Namespace{ID: schema.ServiceDataKeyNamespace},
		SubjectID:        uuid.NewString(),
		SubjectNamespace: namespace.DefinitionUser,
	}
	testAction     = action.Action{ID: schema.EditPermission}
	testRelationV2 = relation.RelationV2{
		ID: uuid.NewString(),
		Object: relation.Object{
			ID:          testRelation.ObjectID,
			NamespaceID: testRelation.ObjectNamespace.ID,
		},
		Subject: relation.Subject{
			ID:        testRelation.SubjectID,
			Namespace: testRelation.SubjectNamespace.ID,
		},
	}
	testUserID                    = uuid.NewString()
	auditKeyRelationCreate        = "relation.create"
	auditKeyRelationSubjectDelete = "relation_subject.delete"
	testResourceID                = uuid.NewString()
)

func TestService_Get(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		id      string
		setup   func(t *testing.T) *relation.Service
		want    relation.RelationV2
		wantErr error
	}{
		{
			name: "GetSuccess",
			id:   testRelationV2.ID,
			setup: func(t *testing.T) *relation.Service {
				t.Helper()
				repository := &mocks.Repository{}
				authzRepository := &mocks.AuthzRepository{}
				userService := &mocks.UserService{}
				activityService := &mocks.ActivityService{}
				repository.EXPECT().Get(mock.Anything, testRelationV2.ID).Return(testRelationV2, nil)
				return relation.NewService(testLogger, repository, authzRepository, userService, activityService)
			},
			want:    testRelationV2,
			wantErr: nil,
		},
		{
			name: "GetErr",
			id:   testRelationV2.ID,
			setup: func(t *testing.T) *relation.Service {
				t.Helper()
				repository := &mocks.Repository{}
				authzRepository := &mocks.AuthzRepository{}
				userService := &mocks.UserService{}
				activityService := &mocks.ActivityService{}
				repository.EXPECT().Get(mock.Anything, testRelationV2.ID).Return(relation.RelationV2{}, relation.ErrInvalidID)
				return relation.NewService(testLogger, repository, authzRepository, userService, activityService)
			},
			wantErr: relation.ErrInvalidID,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := tt.setup(t)

			assert.NotNil(t, svc)

			got, err := svc.Get(context.TODO(), tt.id)
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

func TestService_Create(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		rel     relation.RelationV2
		setup   func(t *testing.T) *relation.Service
		want    relation.RelationV2
		wantErr error
	}{
		{
			name: "CreateSuccess",
			rel:  testRelationV2,
			setup: func(t *testing.T) *relation.Service {
				t.Helper()
				repository := &mocks.Repository{}
				authzRepository := &mocks.AuthzRepository{}
				userService := &mocks.UserService{}
				activityService := &mocks.ActivityService{}
				userService.EXPECT().FetchCurrentUser(mock.Anything).Return(user.User{
					ID: testUserID, Email: "john.doe@gotocompany.com",
				}, nil)
				repository.EXPECT().Create(mock.Anything, testRelationV2).Return(testRelationV2, nil)
				authzRepository.EXPECT().AddV2(mock.Anything, testRelationV2).Return(nil)
				activityService.EXPECT().Log(mock.Anything, auditKeyRelationCreate,
					activity.Actor{Email: "john.doe@gotocompany.com", ID: testUserID}, testRelationV2.ToLogData()).Return(nil)
				return relation.NewService(testLogger, repository, authzRepository, userService, activityService)
			},
			want:    testRelationV2,
			wantErr: nil,
		},
		{
			name: "CreateFetchUserErr",
			rel:  testRelationV2,
			setup: func(t *testing.T) *relation.Service {
				t.Helper()
				repository := &mocks.Repository{}
				authzRepository := &mocks.AuthzRepository{}
				userService := &mocks.UserService{}
				activityService := &mocks.ActivityService{}
				userService.EXPECT().FetchCurrentUser(mock.Anything).Return(user.User{}, user.ErrNotExist)
				return relation.NewService(testLogger, repository, authzRepository, userService, activityService)
			},
			wantErr: user.ErrNotExist,
		},
		{
			name: "CreateErr",
			rel:  testRelationV2,
			setup: func(t *testing.T) *relation.Service {
				t.Helper()
				repository := &mocks.Repository{}
				authzRepository := &mocks.AuthzRepository{}
				userService := &mocks.UserService{}
				activityService := &mocks.ActivityService{}
				userService.EXPECT().FetchCurrentUser(mock.Anything).Return(user.User{
					ID: testUserID, Email: "john.doe@gotocompany.com",
				}, nil)
				repository.EXPECT().Create(mock.Anything, testRelationV2).Return(relation.RelationV2{}, relation.ErrCreatingRelationInStore)
				return relation.NewService(testLogger, repository, authzRepository, userService, activityService)
			},
			wantErr: relation.ErrCreatingRelationInStore,
		},
		{
			name: "CreateAddV2Err",
			rel:  testRelationV2,
			setup: func(t *testing.T) *relation.Service {
				t.Helper()
				repository := &mocks.Repository{}
				authzRepository := &mocks.AuthzRepository{}
				userService := &mocks.UserService{}
				activityService := &mocks.ActivityService{}
				userService.EXPECT().FetchCurrentUser(mock.Anything).Return(user.User{
					ID: testUserID, Email: "john.doe@gotocompany.com",
				}, nil)
				repository.EXPECT().Create(mock.Anything, testRelationV2).Return(testRelationV2, nil)
				authzRepository.EXPECT().AddV2(mock.Anything, testRelationV2).Return(relation.ErrCreatingRelationInAuthzEngine)
				return relation.NewService(testLogger, repository, authzRepository, userService, activityService)
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

			got, err := svc.Create(context.TODO(), tt.rel)
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

func TestService_List(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		setup   func(t *testing.T) *relation.Service
		want    []relation.RelationV2
		wantErr error
	}{
		{
			name: "ListSuccess",
			setup: func(t *testing.T) *relation.Service {
				t.Helper()
				repository := &mocks.Repository{}
				authzRepository := &mocks.AuthzRepository{}
				userService := &mocks.UserService{}
				activityService := &mocks.ActivityService{}
				repository.EXPECT().List(mock.Anything).Return([]relation.RelationV2{testRelationV2}, nil)
				return relation.NewService(testLogger, repository, authzRepository, userService, activityService)
			},
			want: []relation.RelationV2{testRelationV2},
		},
		{
			name: "ListErr",
			setup: func(t *testing.T) *relation.Service {
				t.Helper()
				repository := &mocks.Repository{}
				authzRepository := &mocks.AuthzRepository{}
				userService := &mocks.UserService{}
				activityService := &mocks.ActivityService{}
				repository.EXPECT().List(mock.Anything).Return([]relation.RelationV2{}, relation.ErrFetchingUser)
				return relation.NewService(testLogger, repository, authzRepository, userService, activityService)
			},
			want:    []relation.RelationV2{},
			wantErr: relation.ErrFetchingUser,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := tt.setup(t)

			assert.NotNil(t, svc)

			got, err := svc.List(context.TODO())
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

func TestService_GetRelationByField(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		rel     relation.RelationV2
		setup   func(t *testing.T) *relation.Service
		want    relation.RelationV2
		wantErr error
	}{
		{
			name: "GetRelationByFieldSuccess",
			rel:  testRelationV2,
			setup: func(t *testing.T) *relation.Service {
				t.Helper()
				repository := &mocks.Repository{}
				authzRepository := &mocks.AuthzRepository{}
				userService := &mocks.UserService{}
				activityService := &mocks.ActivityService{}
				repository.EXPECT().GetByFields(mock.Anything, testRelationV2).Return(testRelationV2, nil)
				return relation.NewService(testLogger, repository, authzRepository, userService, activityService)
			},
			want:    testRelationV2,
			wantErr: nil,
		},
		{
			name: "GetRelationByFieldErr",
			rel:  testRelationV2,
			setup: func(t *testing.T) *relation.Service {
				t.Helper()
				repository := &mocks.Repository{}
				authzRepository := &mocks.AuthzRepository{}
				userService := &mocks.UserService{}
				activityService := &mocks.ActivityService{}
				repository.EXPECT().GetByFields(mock.Anything, testRelationV2).Return(relation.RelationV2{}, relation.ErrInvalidDetail)
				return relation.NewService(testLogger, repository, authzRepository, userService, activityService)
			},
			wantErr: relation.ErrInvalidDetail,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := tt.setup(t)

			assert.NotNil(t, svc)

			got, err := svc.GetRelationByFields(context.TODO(), tt.rel)
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

func TestService_DeleteV2(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		rel     relation.RelationV2
		setup   func(t *testing.T) *relation.Service
		wantErr error
	}{
		{
			name: "DeleteV2Success",
			rel:  testRelationV2,
			setup: func(t *testing.T) *relation.Service {
				t.Helper()
				repository := &mocks.Repository{}
				authzRepository := &mocks.AuthzRepository{}
				userService := &mocks.UserService{}
				activityService := &mocks.ActivityService{}
				repository.EXPECT().GetByFields(mock.Anything, testRelationV2).Return(testRelationV2, nil)
				authzRepository.EXPECT().DeleteV2(mock.Anything, testRelationV2).Return(nil)
				repository.EXPECT().DeleteByID(mock.Anything, testRelationV2.ID).Return(nil)
				return relation.NewService(testLogger, repository, authzRepository, userService, activityService)
			},
			wantErr: nil,
		},
		{
			name: "DeleteV2GetByFieldsErr",
			rel:  testRelationV2,
			setup: func(t *testing.T) *relation.Service {
				t.Helper()
				repository := &mocks.Repository{}
				authzRepository := &mocks.AuthzRepository{}
				userService := &mocks.UserService{}
				activityService := &mocks.ActivityService{}
				repository.EXPECT().GetByFields(mock.Anything, testRelationV2).Return(relation.RelationV2{}, relation.ErrInvalidDetail)
				return relation.NewService(testLogger, repository, authzRepository, userService, activityService)
			},
			wantErr: relation.ErrInvalidDetail,
		},
		{
			name: "DeleteV2Err",
			rel:  testRelationV2,
			setup: func(t *testing.T) *relation.Service {
				t.Helper()
				repository := &mocks.Repository{}
				authzRepository := &mocks.AuthzRepository{}
				userService := &mocks.UserService{}
				activityService := &mocks.ActivityService{}
				repository.EXPECT().GetByFields(mock.Anything, testRelationV2).Return(testRelationV2, nil)
				authzRepository.EXPECT().DeleteV2(mock.Anything, testRelationV2).Return(relation.ErrInvalidDetail)
				return relation.NewService(testLogger, repository, authzRepository, userService, activityService)
			},
			wantErr: relation.ErrInvalidDetail,
		},
		{
			name: "DeleteV2DeleteByIDErr",
			rel:  testRelationV2,
			setup: func(t *testing.T) *relation.Service {
				t.Helper()
				repository := &mocks.Repository{}
				authzRepository := &mocks.AuthzRepository{}
				userService := &mocks.UserService{}
				activityService := &mocks.ActivityService{}
				repository.EXPECT().GetByFields(mock.Anything, testRelationV2).Return(testRelationV2, nil)
				authzRepository.EXPECT().DeleteV2(mock.Anything, testRelationV2).Return(nil)
				repository.EXPECT().DeleteByID(mock.Anything, testRelationV2.ID).Return(relation.ErrInvalidDetail)
				return relation.NewService(testLogger, repository, authzRepository, userService, activityService)
			},
			wantErr: relation.ErrInvalidDetail,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := tt.setup(t)

			assert.NotNil(t, svc)

			err := svc.DeleteV2(context.TODO(), tt.rel)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.wantErr))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_CheckPermission(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		usr          user.User
		resNamespace namespace.Namespace
		resIdxa      string
		act          action.Action
		setup        func(t *testing.T) *relation.Service
		want         bool
		wantErr      error
	}{
		{
			name:         "CheckPermissionSuccess",
			usr:          user.User{ID: testRelation.SubjectID},
			resNamespace: testRelation.ObjectNamespace,
			resIdxa:      testRelation.ObjectID,
			act:          testAction,
			setup: func(t *testing.T) *relation.Service {
				t.Helper()
				repository := &mocks.Repository{}
				authzRepository := &mocks.AuthzRepository{}
				userService := &mocks.UserService{}
				activityService := &mocks.ActivityService{}
				authzRepository.EXPECT().Check(mock.Anything, testRelation, testAction).Return(true, nil)
				return relation.NewService(testLogger, repository, authzRepository, userService, activityService)
			},
			want:    true,
			wantErr: nil,
		},
		{
			name:         "CheckPermissionErr",
			usr:          user.User{ID: testRelation.SubjectID},
			resNamespace: testRelation.ObjectNamespace,
			resIdxa:      testRelation.ObjectID,
			act:          testAction,
			setup: func(t *testing.T) *relation.Service {
				t.Helper()
				repository := &mocks.Repository{}
				authzRepository := &mocks.AuthzRepository{}
				userService := &mocks.UserService{}
				activityService := &mocks.ActivityService{}
				authzRepository.EXPECT().Check(mock.Anything, testRelation, testAction).Return(false, relation.ErrInvalidID)
				return relation.NewService(testLogger, repository, authzRepository, userService, activityService)
			},
			want:    false,
			wantErr: relation.ErrInvalidID,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := tt.setup(t)

			assert.NotNil(t, svc)

			got, err := svc.CheckPermission(context.TODO(), tt.usr, tt.resNamespace, tt.resIdxa, tt.act)
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

func TestService_BulkCheckPermission(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		relations []relation.Relation
		actions   []action.Action
		setup     func(t *testing.T) *relation.Service
		want      []relation.Permission
		wantErr   error
	}{
		{
			name:      "BulkCheckPermissionSuccess",
			relations: []relation.Relation{testRelation},
			actions:   []action.Action{testAction},
			setup: func(t *testing.T) *relation.Service {
				t.Helper()
				repository := &mocks.Repository{}
				authzRepository := &mocks.AuthzRepository{}
				userService := &mocks.UserService{}
				activityService := &mocks.ActivityService{}
				authzRepository.EXPECT().BulkCheck(mock.Anything, []relation.Relation{testRelation}, []action.Action{testAction}).
					Return([]relation.Permission{{
						ObjectID:        testRelation.ObjectID,
						ObjectNamespace: testRelation.ObjectNamespace.ID,
						Permission:      testAction.ID,
						Allowed:         true,
					}}, nil)
				return relation.NewService(testLogger, repository, authzRepository, userService, activityService)
			},
			want: []relation.Permission{{
				ObjectID:        testRelation.ObjectID,
				ObjectNamespace: testRelation.ObjectNamespace.ID,
				Permission:      testAction.ID,
				Allowed:         true,
			}},
			wantErr: nil,
		},
		{
			name:      "BulkCheckPermissionErr",
			relations: []relation.Relation{testRelation},
			actions:   []action.Action{testAction},
			setup: func(t *testing.T) *relation.Service {
				t.Helper()
				repository := &mocks.Repository{}
				authzRepository := &mocks.AuthzRepository{}
				userService := &mocks.UserService{}
				activityService := &mocks.ActivityService{}
				authzRepository.EXPECT().BulkCheck(mock.Anything, []relation.Relation{testRelation}, []action.Action{testAction}).
					Return([]relation.Permission{}, relation.ErrInvalidDetail)
				return relation.NewService(testLogger, repository, authzRepository, userService, activityService)
			},
			want:    []relation.Permission{},
			wantErr: relation.ErrInvalidDetail,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := tt.setup(t)

			assert.NotNil(t, svc)

			got, err := svc.BulkCheckPermission(context.TODO(), tt.relations, tt.actions)
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

func TestService_DeleteSubjectRelation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		resourceType       string
		optionalResourceID string
		setup              func(t *testing.T) *relation.Service
		wantErr            error
	}{
		{
			name:               "DeleteSubjectRelationSuccess",
			resourceType:       testRelation.ObjectNamespace.ID,
			optionalResourceID: testRelation.ObjectID,
			setup: func(t *testing.T) *relation.Service {
				t.Helper()
				repository := &mocks.Repository{}
				authzRepository := &mocks.AuthzRepository{}
				userService := &mocks.UserService{}
				activityService := &mocks.ActivityService{}
				userService.EXPECT().FetchCurrentUser(mock.Anything).Return(user.User{ID: testUserID, Email: "john.doe@gotocompany.com"}, nil)
				authzRepository.EXPECT().DeleteSubjectRelations(mock.Anything, testRelation.ObjectNamespace.ID, testRelation.ObjectID).Return(nil)
				activityService.EXPECT().Log(mock.Anything, auditKeyRelationSubjectDelete, activity.Actor{ID: testUserID, Email: "john.doe@gotocompany.com"},
					relation.ToSubjectLogData(testRelation.ObjectNamespace.ID, testRelation.ObjectID)).Return(nil)
				return relation.NewService(testLogger, repository, authzRepository, userService, activityService)
			},
			wantErr: nil,
		},
		{
			name:               "DeleteSubjectRelationFetchUserErr",
			resourceType:       testRelation.ObjectNamespace.ID,
			optionalResourceID: testRelation.ObjectID,
			setup: func(t *testing.T) *relation.Service {
				t.Helper()
				repository := &mocks.Repository{}
				authzRepository := &mocks.AuthzRepository{}
				userService := &mocks.UserService{}
				activityService := &mocks.ActivityService{}
				userService.EXPECT().FetchCurrentUser(mock.Anything).Return(user.User{}, user.ErrInvalidEmail)
				return relation.NewService(testLogger, repository, authzRepository, userService, activityService)
			},
			wantErr: user.ErrInvalidEmail,
		},
		{
			name:               "DeleteSubjectRelationErr",
			resourceType:       testRelation.ObjectNamespace.ID,
			optionalResourceID: testRelation.ObjectID,
			setup: func(t *testing.T) *relation.Service {
				t.Helper()
				repository := &mocks.Repository{}
				authzRepository := &mocks.AuthzRepository{}
				userService := &mocks.UserService{}
				activityService := &mocks.ActivityService{}
				userService.EXPECT().FetchCurrentUser(mock.Anything).Return(user.User{ID: testUserID, Email: "john.doe@gotocompany.com"}, nil)
				authzRepository.EXPECT().DeleteSubjectRelations(mock.Anything, testRelation.ObjectNamespace.ID, testRelation.ObjectID).Return(relation.ErrInvalidDetail)
				return relation.NewService(testLogger, repository, authzRepository, userService, activityService)
			},
			wantErr: relation.ErrInvalidDetail,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := tt.setup(t)

			assert.NotNil(t, svc)

			err := svc.DeleteSubjectRelations(context.TODO(), tt.resourceType, tt.optionalResourceID)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.wantErr))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_LookupResources(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		resourceType string
		permission   string
		subjectType  string
		subjectID    string
		setup        func(t *testing.T) *relation.Service
		want         []string
		wantErr      error
	}{
		{
			name:         "LookupResourceSuccess",
			resourceType: testRelation.ObjectNamespace.ID,
			permission:   testAction.ID,
			subjectType:  testRelation.SubjectNamespace.ID,
			subjectID:    testRelation.SubjectID,
			setup: func(t *testing.T) *relation.Service {
				t.Helper()
				repository := &mocks.Repository{}
				authzRepository := &mocks.AuthzRepository{}
				userService := &mocks.UserService{}
				activityService := &mocks.ActivityService{}
				authzRepository.EXPECT().LookupResources(mock.Anything, testRelation.ObjectNamespace.ID, testAction.ID, testRelation.SubjectNamespace.ID, testRelation.SubjectID).
					Return([]string{testResourceID}, nil)
				return relation.NewService(testLogger, repository, authzRepository, userService, activityService)
			},
			want:    []string{testResourceID},
			wantErr: nil,
		},
		{
			name:         "LookupResourceErr",
			resourceType: testRelation.ObjectNamespace.ID,
			permission:   testAction.ID,
			subjectType:  testRelation.SubjectNamespace.ID,
			subjectID:    testRelation.SubjectID,
			setup: func(t *testing.T) *relation.Service {
				t.Helper()
				repository := &mocks.Repository{}
				authzRepository := &mocks.AuthzRepository{}
				userService := &mocks.UserService{}
				activityService := &mocks.ActivityService{}
				authzRepository.EXPECT().LookupResources(mock.Anything, testRelation.ObjectNamespace.ID, testAction.ID, testRelation.SubjectNamespace.ID, testRelation.SubjectID).
					Return([]string{}, relation.ErrInvalidDetail)
				return relation.NewService(testLogger, repository, authzRepository, userService, activityService)
			},
			want:    []string{},
			wantErr: relation.ErrInvalidDetail,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := tt.setup(t)

			assert.NotNil(t, svc)

			got, err := svc.LookupResources(context.TODO(), tt.resourceType, tt.permission, tt.subjectType, tt.subjectID)
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
