package v1beta1

import (
	"context"
	"errors"
	"testing"

	"github.com/goto/shield/core/action"
	"github.com/goto/shield/core/relation"
	"github.com/goto/shield/core/resource"
	"github.com/goto/shield/core/user"
	"github.com/goto/shield/internal/adapter"
	"github.com/goto/shield/internal/api/v1beta1/mocks"
	"github.com/goto/shield/internal/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	shieldv1beta1 "github.com/goto/shield/proto/v1beta1"
)

var (
	testRelationV2 = relation.RelationV2{
		ID: "relation-id-1",
		Subject: relation.Subject{
			ID:        "subject-id",
			Namespace: "ns1",
			RoleID:    "role1",
		},
		Object: relation.Object{
			ID:          "5e70ba45-dc63-4152-9d72-27cbc34d9d13",
			NamespaceID: "ns2",
		},
	}

	testRelationV2Invalid = relation.RelationV2{
		ID: "relation-id-1",
		Subject: relation.Subject{
			ID:        "subject-id",
			Namespace: "ns1",
			RoleID:    "role1",
		},
		Object: relation.Object{
			ID:          "non-uuid",
			NamespaceID: "ns2",
		},
	}

	testRelationPB = &shieldv1beta1.Relation{
		Id:              "relation-id-1",
		ObjectId:        "5e70ba45-dc63-4152-9d72-27cbc34d9d13",
		ObjectNamespace: "ns2",
		Subject:         "ns1:subject-id",
		RoleName:        "role1",
	}
)

func TestHandler_ListRelations(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(rs *mocks.RelationService)
		want    *shieldv1beta1.ListRelationsResponse
		wantErr error
	}{
		{
			name: "should return internal error if relation service return some error",
			setup: func(rs *mocks.RelationService) {
				rs.EXPECT().List(mock.AnythingOfType("context.todoCtx")).Return([]relation.RelationV2{}, errors.New("some error"))
			},
			want:    nil,
			wantErr: grpcInternalServerError,
		},
		{
			name: "should return relations if relation service return nil error",
			setup: func(rs *mocks.RelationService) {
				rs.EXPECT().List(mock.AnythingOfType("context.todoCtx")).Return([]relation.RelationV2{
					testRelationV2,
				}, nil)
			},
			want: &shieldv1beta1.ListRelationsResponse{
				Relations: []*shieldv1beta1.Relation{
					testRelationPB,
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRelationSrv := new(mocks.RelationService)
			if tt.setup != nil {
				tt.setup(mockRelationSrv)
			}
			mockDep := Handler{relationService: mockRelationSrv}
			resp, err := mockDep.ListRelations(context.TODO(), &shieldv1beta1.ListRelationsRequest{})
			assert.EqualValues(t, tt.want, resp)
			assert.EqualValues(t, tt.wantErr, err)
		})
	}
}

func TestHandler_CreateRelation(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(rs *mocks.RelationService, res *mocks.ResourceService)
		request *shieldv1beta1.CreateRelationRequest
		want    *shieldv1beta1.CreateRelationResponse
		wantErr error
	}{
		{
			name: "should return resource does not exist error when object id does not belong to a resource",
			setup: func(rs *mocks.RelationService, res *mocks.ResourceService) {
				res.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), testRelationV2.Object.ID).Return(resource.Resource{}, errors.New("resource does not exist"))
			},
			request: &shieldv1beta1.CreateRelationRequest{
				Body: &shieldv1beta1.RelationRequestBody{
					ObjectId:        testRelationV2.Object.ID,
					ObjectNamespace: testRelationV2.Object.NamespaceID,
					Subject:         generateSubject(testRelationV2.Subject.ID, testRelationV2.Subject.Namespace),
					RoleName:        testRelationV2.Subject.RoleID,
				},
			},
			want:    nil,
			wantErr: status.Errorf(codes.NotFound, "resource does not exist"),
		},
		{
			name: "should return internal server error if resource service's CheckAuthz function returns some error",
			setup: func(rs *mocks.RelationService, res *mocks.ResourceService) {
				res.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), testRelationV2.Object.ID).Return(resource.Resource{}, nil)

				res.EXPECT().CheckAuthz(mock.AnythingOfType("context.todoCtx"), resource.Resource{
					Name:        testRelationV2.Object.ID,
					NamespaceID: testRelationV2.Object.NamespaceID,
				}, action.Action{ID: schema.EditPermission}).Return(false, errors.New("some error"))
			},
			request: &shieldv1beta1.CreateRelationRequest{
				Body: &shieldv1beta1.RelationRequestBody{
					ObjectId:        testRelationV2.Object.ID,
					ObjectNamespace: testRelationV2.Object.NamespaceID,
					Subject:         generateSubject(testRelationV2.Subject.ID, testRelationV2.Subject.Namespace),
					RoleName:        testRelationV2.Subject.RoleID,
				},
			},
			want:    nil,
			wantErr: grpcInternalServerError,
		},
		{
			name:  "should return bad body error if object id is not uuid",
			setup: func(rs *mocks.RelationService, res *mocks.ResourceService) {},
			request: &shieldv1beta1.CreateRelationRequest{
				Body: &shieldv1beta1.RelationRequestBody{
					ObjectId:        testRelationV2Invalid.Object.ID,
					ObjectNamespace: testRelationV2Invalid.Object.NamespaceID,
					Subject:         generateSubject(testRelationV2Invalid.Subject.ID, testRelationV2Invalid.Subject.Namespace),
					RoleName:        testRelationV2Invalid.Subject.RoleID,
				},
			},
			want:    nil,
			wantErr: grpcBadBodyError,
		},
		{
			name: "should return permision denied error if resource service's CheckAuthz function returns false",
			setup: func(rs *mocks.RelationService, res *mocks.ResourceService) {
				res.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), testRelationV2.Object.ID).Return(resource.Resource{}, nil)

				res.EXPECT().CheckAuthz(mock.AnythingOfType("context.todoCtx"), resource.Resource{
					Name:        testRelationV2.Object.ID,
					NamespaceID: testRelationV2.Object.NamespaceID,
				}, action.Action{ID: schema.EditPermission}).Return(false, nil)
			},
			request: &shieldv1beta1.CreateRelationRequest{
				Body: &shieldv1beta1.RelationRequestBody{
					ObjectId:        testRelationV2.Object.ID,
					ObjectNamespace: testRelationV2.Object.NamespaceID,
					Subject:         generateSubject(testRelationV2.Subject.ID, testRelationV2.Subject.Namespace),
					RoleName:        testRelationV2.Subject.RoleID,
				},
			},
			want:    nil,
			wantErr: grpcPermissionDenied,
		},
		{
			name: "should return unauthenticated error if resource service's CheckAuthz function returns invalid email error",
			setup: func(rs *mocks.RelationService, res *mocks.ResourceService) {
				res.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), testRelationV2.Object.ID).Return(resource.Resource{}, nil)

				res.EXPECT().CheckAuthz(mock.AnythingOfType("context.todoCtx"), resource.Resource{
					Name:        testRelationV2.Object.ID,
					NamespaceID: testRelationV2.Object.NamespaceID,
				}, action.Action{ID: schema.EditPermission}).Return(false, user.ErrInvalidEmail)
			},
			request: &shieldv1beta1.CreateRelationRequest{
				Body: &shieldv1beta1.RelationRequestBody{
					ObjectId:        testRelationV2.Object.ID,
					ObjectNamespace: testRelationV2.Object.NamespaceID,
					Subject:         generateSubject(testRelationV2.Subject.ID, testRelationV2.Subject.Namespace),
					RoleName:        testRelationV2.Subject.RoleID,
				},
			},
			want:    nil,
			wantErr: grpcUnauthenticated,
		},
		{
			name: "should return internal error if relation service return some error",
			setup: func(rs *mocks.RelationService, res *mocks.ResourceService) {
				res.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), testRelationV2.Object.ID).Return(resource.Resource{}, nil)

				res.EXPECT().CheckAuthz(mock.AnythingOfType("context.todoCtx"), resource.Resource{
					Name:        testRelationV2.Object.ID,
					NamespaceID: testRelationV2.Object.NamespaceID,
				}, action.Action{ID: schema.EditPermission}).Return(true, nil)

				rs.EXPECT().Create(mock.AnythingOfType("context.todoCtx"), relation.RelationV2{
					Subject: relation.Subject{
						ID:        testRelationV2.Subject.ID,
						Namespace: testRelationV2.Subject.Namespace,
						RoleID:    testRelationV2.Subject.RoleID,
					},
					Object: relation.Object{
						ID:          testRelationV2.Object.ID,
						NamespaceID: testRelationV2.Object.NamespaceID,
					},
				}).Return(relation.RelationV2{}, errors.New("some error"))
			},
			request: &shieldv1beta1.CreateRelationRequest{
				Body: &shieldv1beta1.RelationRequestBody{
					ObjectId:        testRelationV2.Object.ID,
					ObjectNamespace: testRelationV2.Object.NamespaceID,
					Subject:         generateSubject(testRelationV2.Subject.ID, testRelationV2.Subject.Namespace),
					RoleName:        testRelationV2.Subject.RoleID,
				},
			},
			want:    nil,
			wantErr: grpcInternalServerError,
		},
		{
			name: "should return bad request error if field value not exist in foreign reference",
			setup: func(rs *mocks.RelationService, res *mocks.ResourceService) {
				res.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), testRelationV2.Object.ID).Return(resource.Resource{}, nil)

				res.EXPECT().CheckAuthz(mock.AnythingOfType("context.todoCtx"), resource.Resource{
					Name:        testRelationV2.Object.ID,
					NamespaceID: testRelationV2.Object.NamespaceID,
				}, action.Action{ID: schema.EditPermission}).Return(true, nil)

				rs.EXPECT().Create(mock.AnythingOfType("context.todoCtx"), relation.RelationV2{
					Subject: relation.Subject{
						ID:        testRelationV2.Subject.ID,
						Namespace: testRelationV2.Subject.Namespace,
						RoleID:    testRelationV2.Subject.RoleID,
					},
					Object: relation.Object{
						ID:          testRelationV2.Object.ID,
						NamespaceID: testRelationV2.Object.NamespaceID,
					},
				}).Return(relation.RelationV2{}, relation.ErrInvalidDetail)
			},
			request: &shieldv1beta1.CreateRelationRequest{
				Body: &shieldv1beta1.RelationRequestBody{
					ObjectId:        testRelationV2.Object.ID,
					ObjectNamespace: testRelationV2.Object.NamespaceID,
					Subject:         generateSubject(testRelationV2.Subject.ID, testRelationV2.Subject.Namespace),
					RoleName:        testRelationV2.Subject.RoleID,
				},
			},
			want:    nil,
			wantErr: grpcBadBodyError,
		},
		{
			name: "should return success if relation service return nil",
			setup: func(rs *mocks.RelationService, res *mocks.ResourceService) {
				res.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), testRelationV2.Object.ID).Return(resource.Resource{}, nil)

				res.EXPECT().CheckAuthz(mock.AnythingOfType("context.todoCtx"), resource.Resource{
					Name:        testRelationV2.Object.ID,
					NamespaceID: testRelationV2.Object.NamespaceID,
				}, action.Action{ID: schema.EditPermission}).Return(true, nil)

				rs.EXPECT().Create(mock.AnythingOfType("context.todoCtx"), relation.RelationV2{
					Subject: relation.Subject{
						ID:        testRelationV2.Subject.ID,
						Namespace: testRelationV2.Subject.Namespace,
						RoleID:    testRelationV2.Subject.RoleID,
					},
					Object: relation.Object{
						ID:          testRelationV2.Object.ID,
						NamespaceID: testRelationV2.Object.NamespaceID,
					},
				}).Return(testRelationV2, nil)
			},
			request: &shieldv1beta1.CreateRelationRequest{
				Body: &shieldv1beta1.RelationRequestBody{
					ObjectId:        testRelationV2.Object.ID,
					ObjectNamespace: testRelationV2.Object.NamespaceID,
					Subject:         generateSubject(testRelationV2.Subject.ID, testRelationV2.Subject.Namespace),
					RoleName:        testRelationV2.Subject.RoleID,
				},
			},
			want: &shieldv1beta1.CreateRelationResponse{
				Relation: testRelationPB,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRelationSrv := new(mocks.RelationService)
			mockResourceSrv := new(mocks.ResourceService)
			relationAdapter := new(adapter.Relation)
			if tt.setup != nil {
				tt.setup(mockRelationSrv, mockResourceSrv)
			}

			mockDep := Handler{relationService: mockRelationSrv, resourceService: mockResourceSrv, relationAdapter: relationAdapter}
			resp, err := mockDep.CreateRelation(context.TODO(), tt.request)
			assert.EqualValues(t, tt.want, resp)
			assert.EqualValues(t, tt.wantErr, err)
		})
	}
}

func TestHandler_GetRelation(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(rs *mocks.RelationService)
		request *shieldv1beta1.GetRelationRequest
		want    *shieldv1beta1.GetRelationResponse
		wantErr error
	}{
		{
			name: "should return internal error if relation service return some error",
			setup: func(rs *mocks.RelationService) {
				rs.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), testRelationV2.ID).Return(relation.RelationV2{}, errors.New("some error"))
			},
			request: &shieldv1beta1.GetRelationRequest{
				Id: testRelationV2.ID,
			},
			want:    nil,
			wantErr: grpcInternalServerError,
		},
		{
			name: "should return not found error if id is empty",
			setup: func(rs *mocks.RelationService) {
				rs.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), "").Return(relation.RelationV2{}, relation.ErrInvalidID)
			},
			request: &shieldv1beta1.GetRelationRequest{},
			want:    nil,
			wantErr: grpcRelationNotFoundErr,
		},
		{
			name: "should return not found error if id is not uuid",
			setup: func(rs *mocks.RelationService) {
				rs.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), "some-id").Return(relation.RelationV2{}, relation.ErrInvalidUUID)
			},
			request: &shieldv1beta1.GetRelationRequest{
				Id: "some-id",
			},
			want:    nil,
			wantErr: grpcRelationNotFoundErr,
		},
		{
			name: "should return not found error if id not exist",
			setup: func(rs *mocks.RelationService) {
				rs.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), testRelationV2.ID).Return(relation.RelationV2{}, relation.ErrNotExist)
			},
			request: &shieldv1beta1.GetRelationRequest{
				Id: testRelationV2.ID,
			},
			want:    nil,
			wantErr: grpcRelationNotFoundErr,
		},
		{
			name: "should return success if relation service return nil error",
			setup: func(rs *mocks.RelationService) {
				rs.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), testRelationV2.ID).Return(testRelationV2, nil)
			},
			request: &shieldv1beta1.GetRelationRequest{
				Id: testRelationV2.ID,
			},
			want: &shieldv1beta1.GetRelationResponse{
				Relation: testRelationPB,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRelationSrv := new(mocks.RelationService)
			if tt.setup != nil {
				tt.setup(mockRelationSrv)
			}
			mockDep := Handler{relationService: mockRelationSrv}
			resp, err := mockDep.GetRelation(context.TODO(), tt.request)
			assert.EqualValues(t, tt.want, resp)
			assert.EqualValues(t, tt.wantErr, err)
		})
	}
}

func TestHandler_DeleteRelation(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(rs *mocks.RelationService, res *mocks.ResourceService)
		request *shieldv1beta1.DeleteRelationRequest
		want    *shieldv1beta1.DeleteRelationResponse
		wantErr error
	}{
		{
			name: "should return internal error if relation service return some error",
			setup: func(rs *mocks.RelationService, res *mocks.ResourceService) {
				rs.EXPECT().GetRelationByFields(mock.AnythingOfType("context.todoCtx"), relation.RelationV2{
					Subject: relation.Subject{
						ID:     testRelationV2.Subject.ID,
						RoleID: testRelationV2.Subject.RoleID,
					},
					Object: relation.Object{
						ID: testRelationV2.Object.ID,
					},
				}).Return(relation.RelationV2{}, errors.New("some error"))
			},
			request: &shieldv1beta1.DeleteRelationRequest{
				ObjectId:  testRelationV2.Object.ID,
				SubjectId: testRelationV2.Subject.ID,
				Role:      testRelationV2.Subject.RoleID,
			},
			want:    nil,
			wantErr: grpcInternalServerError,
		},
		{
			name: "should return not found error if relation does not exist",
			setup: func(rs *mocks.RelationService, res *mocks.ResourceService) {
				rs.EXPECT().GetRelationByFields(mock.AnythingOfType("context.todoCtx"), relation.RelationV2{
					Subject: relation.Subject{
						ID:     testRelationV2.Subject.ID,
						RoleID: testRelationV2.Subject.RoleID,
					},
					Object: relation.Object{
						ID: testRelationV2.Object.ID,
					},
				}).Return(relation.RelationV2{}, relation.ErrNotExist)
			},
			request: &shieldv1beta1.DeleteRelationRequest{
				ObjectId:  testRelationV2.Object.ID,
				SubjectId: testRelationV2.Subject.ID,
				Role:      testRelationV2.Subject.RoleID,
			},
			want:    nil,
			wantErr: grpcRelationNotFoundErr,
		},
		{
			name: "should return internal server error if resource service's CheckAuthz returns some error",
			setup: func(rs *mocks.RelationService, res *mocks.ResourceService) {
				rs.EXPECT().GetRelationByFields(mock.AnythingOfType("context.todoCtx"), relation.RelationV2{
					Subject: relation.Subject{
						ID:     testRelationV2.Subject.ID,
						RoleID: testRelationV2.Subject.RoleID,
					},
					Object: relation.Object{
						ID: testRelationV2.Object.ID,
					},
				}).Return(relation.RelationV2{
					Object: relation.Object{
						ID:          testRelationV2.Object.ID,
						NamespaceID: testRelationV2.Object.NamespaceID,
					},
					Subject: relation.Subject{
						ID:        testRelationV2.Subject.ID,
						Namespace: testRelationV2.Subject.Namespace,
						RoleID:    testRelationV2.Subject.RoleID,
					},
				}, nil)

				res.EXPECT().CheckAuthz(mock.AnythingOfType("context.todoCtx"), resource.Resource{
					Name:        testRelationV2.Object.ID,
					NamespaceID: testRelationV2.Object.NamespaceID,
				}, action.Action{ID: schema.EditPermission}).Return(false, errors.New("some error"))
			},
			request: &shieldv1beta1.DeleteRelationRequest{
				ObjectId:  testRelationV2.Object.ID,
				SubjectId: testRelationV2.Subject.ID,
				Role:      testRelationV2.Subject.RoleID,
			},
			want:    nil,
			wantErr: grpcInternalServerError,
		},
		{
			name: "should return permission denied error if resource service returns false response while checking permission",
			setup: func(rs *mocks.RelationService, res *mocks.ResourceService) {
				rs.EXPECT().GetRelationByFields(mock.AnythingOfType("context.todoCtx"), relation.RelationV2{
					Subject: relation.Subject{
						ID:     testRelationV2.Subject.ID,
						RoleID: testRelationV2.Subject.RoleID,
					},
					Object: relation.Object{
						ID: testRelationV2.Object.ID,
					},
				}).Return(relation.RelationV2{
					Object: relation.Object{
						ID:          testRelationV2.Object.ID,
						NamespaceID: testRelationV2.Object.NamespaceID,
					},
					Subject: relation.Subject{
						ID:        testRelationV2.Subject.ID,
						Namespace: testRelationV2.Subject.Namespace,
						RoleID:    testRelationV2.Subject.RoleID,
					},
				}, nil)

				res.EXPECT().CheckAuthz(mock.AnythingOfType("context.todoCtx"), resource.Resource{
					Name:        testRelationV2.Object.ID,
					NamespaceID: testRelationV2.Object.NamespaceID,
				}, action.Action{ID: schema.EditPermission}).Return(false, nil)
			},
			request: &shieldv1beta1.DeleteRelationRequest{
				ObjectId:  testRelationV2.Object.ID,
				SubjectId: testRelationV2.Subject.ID,
				Role:      testRelationV2.Subject.RoleID,
			},
			want:    nil,
			wantErr: grpcPermissionDenied,
		},
		{
			name: "should return internal server error when relation service returns some error while deletion",
			setup: func(rs *mocks.RelationService, res *mocks.ResourceService) {
				rs.EXPECT().GetRelationByFields(mock.AnythingOfType("context.todoCtx"), relation.RelationV2{
					Subject: relation.Subject{
						ID:     testRelationV2.Subject.ID,
						RoleID: testRelationV2.Subject.RoleID,
					},
					Object: relation.Object{
						ID: testRelationV2.Object.ID,
					},
				}).Return(relation.RelationV2{
					Object: relation.Object{
						ID:          testRelationV2.Object.ID,
						NamespaceID: testRelationV2.Object.NamespaceID,
					},
					Subject: relation.Subject{
						ID:        testRelationV2.Subject.ID,
						Namespace: testRelationV2.Subject.Namespace,
						RoleID:    testRelationV2.Subject.RoleID,
					},
				}, nil)

				res.EXPECT().CheckAuthz(mock.AnythingOfType("context.todoCtx"), resource.Resource{
					Name:        testRelationV2.Object.ID,
					NamespaceID: testRelationV2.Object.NamespaceID,
				}, action.Action{ID: schema.EditPermission}).Return(true, nil)

				rs.EXPECT().DeleteV2(mock.AnythingOfType("context.todoCtx"), relation.RelationV2{
					Subject: relation.Subject{
						ID:     testRelationV2.Subject.ID,
						RoleID: testRelationV2.Subject.RoleID,
					},
					Object: relation.Object{
						ID: testRelationV2.Object.ID,
					},
				}).Return(nil)
			},
			request: &shieldv1beta1.DeleteRelationRequest{
				ObjectId:  testRelationV2.Object.ID,
				SubjectId: testRelationV2.Subject.ID,
				Role:      testRelationV2.Subject.RoleID,
			},
			want: &shieldv1beta1.DeleteRelationResponse{
				Message: "Relation deleted",
			},
			wantErr: nil,
		},
		{
			name: "should successfully delete when relation exist and user has permission to edit it",
			setup: func(rs *mocks.RelationService, res *mocks.ResourceService) {
				rs.EXPECT().GetRelationByFields(mock.AnythingOfType("context.todoCtx"), relation.RelationV2{
					Subject: relation.Subject{
						ID:     testRelationV2.Subject.ID,
						RoleID: testRelationV2.Subject.RoleID,
					},
					Object: relation.Object{
						ID: testRelationV2.Object.ID,
					},
				}).Return(relation.RelationV2{
					Object: relation.Object{
						ID:          testRelationV2.Object.ID,
						NamespaceID: testRelationV2.Object.NamespaceID,
					},
					Subject: relation.Subject{
						ID:        testRelationV2.Subject.ID,
						Namespace: testRelationV2.Subject.Namespace,
						RoleID:    testRelationV2.Subject.RoleID,
					},
				}, nil)

				res.EXPECT().CheckAuthz(mock.AnythingOfType("context.todoCtx"), resource.Resource{
					Name:        testRelationV2.Object.ID,
					NamespaceID: testRelationV2.Object.NamespaceID,
				}, action.Action{ID: schema.EditPermission}).Return(true, nil)

				rs.EXPECT().DeleteV2(mock.AnythingOfType("context.todoCtx"), relation.RelationV2{
					Subject: relation.Subject{
						ID:     testRelationV2.Subject.ID,
						RoleID: testRelationV2.Subject.RoleID,
					},
					Object: relation.Object{
						ID: testRelationV2.Object.ID,
					},
				}).Return(nil)
			},
			request: &shieldv1beta1.DeleteRelationRequest{
				ObjectId:  testRelationV2.Object.ID,
				SubjectId: testRelationV2.Subject.ID,
				Role:      testRelationV2.Subject.RoleID,
			},
			want: &shieldv1beta1.DeleteRelationResponse{
				Message: "Relation deleted",
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRelationSrv := new(mocks.RelationService)
			mockResourceSrv := new(mocks.ResourceService)
			if tt.setup != nil {
				tt.setup(mockRelationSrv, mockResourceSrv)
			}
			mockDep := Handler{relationService: mockRelationSrv, resourceService: mockResourceSrv}
			resp, err := mockDep.DeleteRelation(context.TODO(), tt.request)
			assert.EqualValues(t, tt.want, resp)
			assert.EqualValues(t, tt.wantErr, err)
		})
	}
}
