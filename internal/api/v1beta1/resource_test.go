package v1beta1

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/goto/shield/core/organization"
	"github.com/goto/shield/core/project"
	"github.com/goto/shield/core/relation"
	"github.com/goto/shield/core/resource"
	"github.com/goto/shield/core/user"
	"github.com/goto/shield/internal/api/v1beta1/mocks"
	"github.com/goto/shield/pkg/uuid"
	shieldv1beta1 "github.com/goto/shield/proto/v1beta1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	testResourceID = uuid.NewString()
	testResource   = resource.Resource{
		Idxa:           testResourceID,
		URN:            "res-urn",
		Name:           "a resource name",
		ProjectID:      testProjectID,
		OrganizationID: testOrgID,
		NamespaceID:    testNSID,
		UserID:         testUserID,
	}
	testResourcePB = &shieldv1beta1.Resource{
		Id:   testResource.Idxa,
		Name: testResource.Name,
		Urn:  testResource.URN,
		Project: &shieldv1beta1.Project{
			Id: testProjectID,
		},
		Organization: &shieldv1beta1.Organization{
			Id: testOrgID,
		},
		Namespace: &shieldv1beta1.Namespace{
			Id: testNSID,
		},
		User: &shieldv1beta1.User{
			Id: testUserID,
		},
		CreatedAt: timestamppb.New(time.Time{}),
		UpdatedAt: timestamppb.New(time.Time{}),
	}
)

func TestHandler_ListResources(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(rs *mocks.ResourceService)
		request *shieldv1beta1.ListResourcesRequest
		want    *shieldv1beta1.ListResourcesResponse
		wantErr error
	}{
		{
			name: "should return internal error if resource service return some error",
			setup: func(rs *mocks.ResourceService) {
				rs.EXPECT().List(mock.AnythingOfType("context.todoCtx"), resource.Filter{}).Return(resource.PagedResources{}, errors.New("some error"))
			},
			request: &shieldv1beta1.ListResourcesRequest{},
			want:    nil,
			wantErr: grpcInternalServerError,
		},
		{
			name: "should return resources if resource service return nil error",
			setup: func(rs *mocks.ResourceService) {
				testResourceList := []resource.Resource{testResource}
				rs.EXPECT().List(mock.AnythingOfType("context.todoCtx"), resource.Filter{}).Return(
					resource.PagedResources{
						Count:     int32(len(testResourceList)),
						Resources: testResourceList,
					}, nil)
			},
			request: &shieldv1beta1.ListResourcesRequest{},
			want: &shieldv1beta1.ListResourcesResponse{
				Count: int32(
					len([]*shieldv1beta1.Resource{
						testResourcePB,
					},
					)),
				Resources: []*shieldv1beta1.Resource{
					testResourcePB,
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockResourceSrv := new(mocks.ResourceService)
			if tt.setup != nil {
				tt.setup(mockResourceSrv)
			}
			mockDep := Handler{resourceService: mockResourceSrv}
			resp, err := mockDep.ListResources(context.TODO(), tt.request)
			assert.EqualValues(t, tt.want, resp)
			assert.EqualValues(t, tt.wantErr, err)
		})
	}
}

func TestHandler_CreateResource(t *testing.T) {
	email := "user@gotocompany.com"
	tests := []struct {
		name    string
		setup   func(ctx context.Context, rs *mocks.ResourceService, ps *mocks.ProjectService, rls *mocks.RelationService, rt *mocks.RelationTransformer) context.Context
		request *shieldv1beta1.CreateResourceRequest
		want    *shieldv1beta1.CreateResourceResponse
		wantErr error
	}{
		{
			name: "should return unauthenticated error if auth email in context is empty and org service return invalid user email",
			setup: func(ctx context.Context, rs *mocks.ResourceService, ps *mocks.ProjectService, rls *mocks.RelationService, _ *mocks.RelationTransformer) context.Context {
				ps.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), testResource.ProjectID).Return(project.Project{}, user.ErrInvalidEmail)

				rs.EXPECT().Upsert(mock.AnythingOfType("context.todoCtx"), resource.Resource{
					Name:           testResource.Name,
					ProjectID:      testResource.ProjectID,
					OrganizationID: testResource.OrganizationID,
					NamespaceID:    testResource.NamespaceID,
				}).Return(resource.Resource{}, user.ErrInvalidEmail)
				return ctx
			},
			request: &shieldv1beta1.CreateResourceRequest{
				Body: &shieldv1beta1.ResourceRequestBody{
					Name:        testResource.Name,
					ProjectId:   testResource.ProjectID,
					NamespaceId: testResource.NamespaceID,
					Relations: []*shieldv1beta1.Relation{
						{
							RoleName: "owner",
							Subject:  "user:" + testResource.UserID,
						},
					},
				},
			},
			want:    nil,
			wantErr: grpcUnauthenticated,
		},
		{
			name: "should return internal error if project service return some error",
			setup: func(ctx context.Context, rs *mocks.ResourceService, ps *mocks.ProjectService, rls *mocks.RelationService, _ *mocks.RelationTransformer) context.Context {
				ps.EXPECT().Get(mock.AnythingOfType("*context.valueCtx"), testResource.ProjectID).Return(project.Project{}, errors.New("some error"))
				return user.SetContextWithEmail(ctx, email)
			},
			request: &shieldv1beta1.CreateResourceRequest{
				Body: &shieldv1beta1.ResourceRequestBody{
					Name:        testResource.Name,
					ProjectId:   testResource.ProjectID,
					NamespaceId: testResource.NamespaceID,
					Relations: []*shieldv1beta1.Relation{
						{
							RoleName: "owner",
							Subject:  "user:" + testUserID,
						},
					},
				},
			},
			want:    nil,
			wantErr: grpcInternalServerError,
		},
		{
			name: "should return internal error if resource service return some error",
			setup: func(ctx context.Context, rs *mocks.ResourceService, ps *mocks.ProjectService, rls *mocks.RelationService, _ *mocks.RelationTransformer) context.Context {
				ps.EXPECT().Get(mock.AnythingOfType("*context.valueCtx"), testResource.ProjectID).Return(project.Project{
					ID: testResourceID,
					Organization: organization.Organization{
						ID: testResource.OrganizationID,
					},
				}, nil)

				rs.EXPECT().Upsert(mock.AnythingOfType("*context.valueCtx"), resource.Resource{
					Name:           testResource.Name,
					ProjectID:      testResource.ProjectID,
					NamespaceID:    testResource.NamespaceID,
					OrganizationID: testResource.OrganizationID,
				}).Return(resource.Resource{}, errors.New("some error"))
				return user.SetContextWithEmail(ctx, email)
			},
			request: &shieldv1beta1.CreateResourceRequest{
				Body: &shieldv1beta1.ResourceRequestBody{
					Name:        testResource.Name,
					ProjectId:   testResource.ProjectID,
					NamespaceId: testResource.NamespaceID,
					Relations: []*shieldv1beta1.Relation{
						{
							RoleName: "owner",
							Subject:  "user:" + testUserID,
						},
					},
				},
			},
			want:    nil,
			wantErr: grpcInternalServerError,
		},
		{
			name: "should return bad request error if field value not exist in foreign reference",
			setup: func(ctx context.Context, rs *mocks.ResourceService, ps *mocks.ProjectService, rls *mocks.RelationService, _ *mocks.RelationTransformer) context.Context {
				ps.EXPECT().Get(mock.AnythingOfType("*context.valueCtx"), testResource.ProjectID).Return(project.Project{
					ID: testResourceID,
					Organization: organization.Organization{
						ID: testResource.OrganizationID,
					},
				}, nil)

				rs.EXPECT().Upsert(mock.AnythingOfType("*context.valueCtx"), resource.Resource{
					Name:           testResource.Name,
					ProjectID:      testResource.ProjectID,
					OrganizationID: testResource.OrganizationID,
					NamespaceID:    testResource.NamespaceID,
				}).Return(resource.Resource{}, resource.ErrInvalidDetail)
				return user.SetContextWithEmail(ctx, email)
			},
			request: &shieldv1beta1.CreateResourceRequest{
				Body: &shieldv1beta1.ResourceRequestBody{
					Name:        testResource.Name,
					ProjectId:   testResource.ProjectID,
					NamespaceId: testResource.NamespaceID,
					Relations: []*shieldv1beta1.Relation{
						{
							RoleName: "owner",
							Subject:  "user:" + testUserID,
						},
					},
				},
			},
			want:    nil,
			wantErr: grpcBadBodyError,
		},
		{
			name: "should return success if resource service return nil",
			setup: func(ctx context.Context, rs *mocks.ResourceService, ps *mocks.ProjectService, rls *mocks.RelationService, rt *mocks.RelationTransformer) context.Context {
				ps.EXPECT().Get(mock.AnythingOfType("*context.valueCtx"), testResource.ProjectID).Return(project.Project{
					ID: testResourceID,
					Organization: organization.Organization{
						ID: testResource.OrganizationID,
					},
				}, nil)

				theRelation := relation.RelationV2{
					Object: relation.Object{
						ID:          testResource.Idxa,
						NamespaceID: testResource.NamespaceID,
					},
					Subject: relation.Subject{
						RoleID:    "owner",
						Namespace: "user",
						ID:        testUserID,
					},
				}
				rt.EXPECT().TransformRelation(mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("relation.RelationV2")).Return(theRelation, nil)

				rls.EXPECT().Create(mock.AnythingOfType("*context.valueCtx"), theRelation).Return(relation.RelationV2{}, nil)

				rs.EXPECT().Upsert(mock.AnythingOfType("*context.valueCtx"), resource.Resource{
					Name:           testResource.Name,
					ProjectID:      testResource.ProjectID,
					OrganizationID: testResource.OrganizationID,
					NamespaceID:    testResource.NamespaceID,
				}).Return(testResource, nil)
				return user.SetContextWithEmail(ctx, email)
			},
			request: &shieldv1beta1.CreateResourceRequest{
				Body: &shieldv1beta1.ResourceRequestBody{
					Name:        testResource.Name,
					ProjectId:   testResource.ProjectID,
					NamespaceId: testResource.NamespaceID,
					Relations: []*shieldv1beta1.Relation{
						{
							RoleName: "owner",
							Subject:  "user:" + testUserID,
						},
					},
				},
			},
			want: &shieldv1beta1.CreateResourceResponse{
				Resource: testResourcePB,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockResourceSrv := new(mocks.ResourceService)
			mockProjectSrv := new(mocks.ProjectService)
			mockRelationSrv := new(mocks.RelationService)
			relationAdapter := new(mocks.RelationTransformer)
			ctx := context.TODO()
			if tt.setup != nil {
				ctx = tt.setup(ctx, mockResourceSrv, mockProjectSrv, mockRelationSrv, relationAdapter)
			}
			mockDep := Handler{resourceService: mockResourceSrv, projectService: mockProjectSrv, relationService: mockRelationSrv, relationAdapter: relationAdapter}
			resp, err := mockDep.CreateResource(ctx, tt.request)
			assert.EqualValues(t, tt.want, resp)
			assert.EqualValues(t, tt.wantErr, err)
		})
	}
}

func TestHandler_GetResource(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(rs *mocks.ResourceService)
		request *shieldv1beta1.GetResourceRequest
		want    *shieldv1beta1.GetResourceResponse
		wantErr error
	}{
		{
			name: "should return internal error if resource service return some error",
			setup: func(rs *mocks.ResourceService) {
				rs.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), testResource.Idxa).Return(resource.Resource{}, errors.New("some error"))
			},
			request: &shieldv1beta1.GetResourceRequest{
				Id: testResource.Idxa,
			},
			want:    nil,
			wantErr: grpcInternalServerError,
		},
		{
			name: "should return not found error if id is empty",
			setup: func(rs *mocks.ResourceService) {
				rs.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), "").Return(resource.Resource{}, resource.ErrInvalidID)
			},
			request: &shieldv1beta1.GetResourceRequest{},
			want:    nil,
			wantErr: grpcResourceNotFoundErr,
		},
		{
			name: "should return not found error if id is not uuid",
			setup: func(rs *mocks.ResourceService) {
				rs.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), "some-id").Return(resource.Resource{}, resource.ErrInvalidUUID)
			},
			request: &shieldv1beta1.GetResourceRequest{
				Id: "some-id",
			},
			want:    nil,
			wantErr: grpcResourceNotFoundErr,
		},
		{
			name: "should return not found error if id not exist",
			setup: func(rs *mocks.ResourceService) {
				rs.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), testResource.Idxa).Return(resource.Resource{}, resource.ErrNotExist)
			},
			request: &shieldv1beta1.GetResourceRequest{
				Id: testResource.Idxa,
			},
			want:    nil,
			wantErr: grpcResourceNotFoundErr,
		},
		{
			name: "should return success if resource service return nil error",
			setup: func(rs *mocks.ResourceService) {
				rs.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), testResource.Idxa).Return(testResource, nil)
			},
			request: &shieldv1beta1.GetResourceRequest{
				Id: testResource.Idxa,
			},
			want: &shieldv1beta1.GetResourceResponse{
				Resource: testResourcePB,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockResourceSrv := new(mocks.ResourceService)
			if tt.setup != nil {
				tt.setup(mockResourceSrv)
			}
			mockDep := Handler{resourceService: mockResourceSrv}
			resp, err := mockDep.GetResource(context.TODO(), tt.request)
			assert.EqualValues(t, tt.want, resp)
			assert.EqualValues(t, tt.wantErr, err)
		})
	}
}

func TestHandler_UpdateResource(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(rs *mocks.ResourceService, ps *mocks.ProjectService)
		request *shieldv1beta1.UpdateResourceRequest
		want    *shieldv1beta1.UpdateResourceResponse
		wantErr error
	}{
		{
			name: "should return internal error if project service return some error",
			setup: func(rs *mocks.ResourceService, ps *mocks.ProjectService) {
				ps.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), testResource.ProjectID).Return(project.Project{}, errors.New("some error"))
			},
			request: &shieldv1beta1.UpdateResourceRequest{
				Id: testResourceID,
				Body: &shieldv1beta1.ResourceRequestBody{
					Name:        testResource.Name,
					ProjectId:   testResource.ProjectID,
					NamespaceId: testResource.NamespaceID,
				},
			},
			want:    nil,
			wantErr: grpcInternalServerError,
		},
		{
			name: "should return internal error if resource service return some error",
			setup: func(rs *mocks.ResourceService, ps *mocks.ProjectService) {
				ps.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), testResource.ProjectID).Return(project.Project{
					ID: testResourceID,
					Organization: organization.Organization{
						ID: testResource.OrganizationID,
					},
				}, nil)

				rs.EXPECT().Update(mock.AnythingOfType("context.todoCtx"), testResourceID, resource.Resource{
					Name:           testResource.Name,
					ProjectID:      testResource.ProjectID,
					OrganizationID: testResource.OrganizationID,
					NamespaceID:    testResource.NamespaceID,
				}).Return(resource.Resource{}, errors.New("some error"))
			},
			request: &shieldv1beta1.UpdateResourceRequest{
				Id: testResourceID,
				Body: &shieldv1beta1.ResourceRequestBody{
					Name:        testResource.Name,
					ProjectId:   testResource.ProjectID,
					NamespaceId: testResource.NamespaceID,
				},
			},
			want:    nil,
			wantErr: grpcInternalServerError,
		},
		{
			name: "should return not found error if id is empty",
			setup: func(rs *mocks.ResourceService, ps *mocks.ProjectService) {
				ps.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), testResource.ProjectID).Return(project.Project{
					ID: testResourceID,
					Organization: organization.Organization{
						ID: testResource.OrganizationID,
					},
				}, nil)

				rs.EXPECT().Update(mock.AnythingOfType("context.todoCtx"), "", resource.Resource{
					Name:           testResource.Name,
					ProjectID:      testResource.ProjectID,
					NamespaceID:    testResource.NamespaceID,
					OrganizationID: testResource.OrganizationID,
				}).Return(resource.Resource{}, resource.ErrInvalidID)
			},
			request: &shieldv1beta1.UpdateResourceRequest{
				Body: &shieldv1beta1.ResourceRequestBody{
					Name:        testResource.Name,
					ProjectId:   testResource.ProjectID,
					NamespaceId: testResource.NamespaceID,
				},
			},
			want:    nil,
			wantErr: grpcResourceNotFoundErr,
		},
		{
			name: "should return not found error if id is not exist",
			setup: func(rs *mocks.ResourceService, ps *mocks.ProjectService) {
				ps.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), testResource.ProjectID).Return(project.Project{
					ID: testResourceID,
					Organization: organization.Organization{
						ID: testResource.OrganizationID,
					},
				}, nil)

				rs.EXPECT().Update(mock.AnythingOfType("context.todoCtx"), testResourceID, resource.Resource{
					Name:           testResource.Name,
					ProjectID:      testResource.ProjectID,
					OrganizationID: testResource.OrganizationID,
					NamespaceID:    testResource.NamespaceID,
				}).Return(resource.Resource{}, resource.ErrNotExist)
			},
			request: &shieldv1beta1.UpdateResourceRequest{
				Id: testResourceID,
				Body: &shieldv1beta1.ResourceRequestBody{
					Name:        testResource.Name,
					ProjectId:   testResource.ProjectID,
					NamespaceId: testResource.NamespaceID,
				},
			},
			want:    nil,
			wantErr: grpcResourceNotFoundErr,
		},
		{
			name: "should return not found error if id is not uuid",
			setup: func(rs *mocks.ResourceService, ps *mocks.ProjectService) {
				ps.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), testResource.ProjectID).Return(project.Project{
					ID: testResourceID,
					Organization: organization.Organization{
						ID: testResource.OrganizationID,
					},
				}, nil)

				rs.EXPECT().Update(mock.AnythingOfType("context.todoCtx"), "some-id", resource.Resource{
					Name:           testResource.Name,
					ProjectID:      testResource.ProjectID,
					OrganizationID: testResource.OrganizationID,
					NamespaceID:    testResource.NamespaceID,
				}).Return(resource.Resource{}, resource.ErrInvalidUUID)
			},
			request: &shieldv1beta1.UpdateResourceRequest{
				Id: "some-id",
				Body: &shieldv1beta1.ResourceRequestBody{
					Name:        testResource.Name,
					ProjectId:   testResource.ProjectID,
					NamespaceId: testResource.NamespaceID,
				},
			},
			want:    nil,
			wantErr: grpcResourceNotFoundErr,
		},
		{
			name: "should return bad request error if field value not exist in foreign reference",
			setup: func(rs *mocks.ResourceService, ps *mocks.ProjectService) {
				ps.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), testResource.ProjectID).Return(project.Project{
					ID: testResourceID,
					Organization: organization.Organization{
						ID: testResource.OrganizationID,
					},
				}, nil)

				rs.EXPECT().Update(mock.AnythingOfType("context.todoCtx"), testResourceID, resource.Resource{
					Name:           testResource.Name,
					ProjectID:      testResource.ProjectID,
					OrganizationID: testResource.OrganizationID,
					NamespaceID:    testResource.NamespaceID,
				}).Return(resource.Resource{}, resource.ErrInvalidDetail)
			},
			request: &shieldv1beta1.UpdateResourceRequest{
				Id: testResourceID,
				Body: &shieldv1beta1.ResourceRequestBody{
					Name:        testResource.Name,
					ProjectId:   testResource.ProjectID,
					NamespaceId: testResource.NamespaceID,
				},
			},
			want:    nil,
			wantErr: grpcBadBodyError,
		},
		{
			name: "should return already exist error if resource service return err conflict",
			setup: func(rs *mocks.ResourceService, ps *mocks.ProjectService) {
				ps.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), testResource.ProjectID).Return(project.Project{
					ID: testResourceID,
					Organization: organization.Organization{
						ID: testResource.OrganizationID,
					},
				}, nil)

				rs.EXPECT().Update(mock.AnythingOfType("context.todoCtx"), testResourceID, resource.Resource{
					Name:           testResource.Name,
					ProjectID:      testResource.ProjectID,
					OrganizationID: testResource.OrganizationID,
					NamespaceID:    testResource.NamespaceID,
				}).Return(resource.Resource{}, resource.ErrConflict)
			},
			request: &shieldv1beta1.UpdateResourceRequest{
				Id: testResourceID,
				Body: &shieldv1beta1.ResourceRequestBody{
					Name:        testResource.Name,
					ProjectId:   testResource.ProjectID,
					NamespaceId: testResource.NamespaceID,
				},
			},
			want:    nil,
			wantErr: grpcConflictError,
		},
		{
			name: "should return success if resource service return nil",
			setup: func(rs *mocks.ResourceService, ps *mocks.ProjectService) {
				ps.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), testResource.ProjectID).Return(project.Project{
					ID: testResourceID,
					Organization: organization.Organization{
						ID: testResource.OrganizationID,
					},
				}, nil)

				rs.EXPECT().Update(mock.AnythingOfType("context.todoCtx"), testResourceID, resource.Resource{
					Name:           testResource.Name,
					ProjectID:      testResource.ProjectID,
					NamespaceID:    testResource.NamespaceID,
					OrganizationID: testResource.OrganizationID,
				}).Return(testResource, nil)
			},
			request: &shieldv1beta1.UpdateResourceRequest{
				Id: testResourceID,
				Body: &shieldv1beta1.ResourceRequestBody{
					Name:        testResource.Name,
					ProjectId:   testResource.ProjectID,
					NamespaceId: testResource.NamespaceID,
				},
			},
			want: &shieldv1beta1.UpdateResourceResponse{
				Resource: testResourcePB,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockResourceSrv := new(mocks.ResourceService)
			mockProjectSrv := new(mocks.ProjectService)
			if tt.setup != nil {
				tt.setup(mockResourceSrv, mockProjectSrv)
			}
			mockDep := Handler{resourceService: mockResourceSrv, projectService: mockProjectSrv}
			resp, err := mockDep.UpdateResource(context.TODO(), tt.request)
			assert.EqualValues(t, tt.want, resp)
			assert.EqualValues(t, tt.wantErr, err)
		})
	}
}
