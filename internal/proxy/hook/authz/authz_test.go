package authz

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/goto/shield/core/namespace"
	"github.com/goto/shield/core/relation"
	"github.com/goto/shield/core/resource"
	"github.com/goto/shield/core/rule"
	"github.com/goto/shield/internal/proxy/attribute"
	"github.com/goto/shield/internal/proxy/hook"
	"github.com/goto/shield/internal/proxy/hook/authz/mocks"
	shieldlogger "github.com/goto/shield/pkg/logger"
	"github.com/goto/shield/pkg/uuid"
)

var testPermissionAttributesMap = map[string]any{
	"project":       "ab657ae7-8c9e-45eb-9862-dd9ceb6d5c71",
	"resource":      []string{"resc1", "resc2"},
	"organization":  "org1",
	"namespace":     "ns1",
	"resource_type": "kind",
	"group":         "group@gotocompany.com",
	"user":          "user1@gotocompany.com",
}

var expectedResources = []resource.Resource{
	{
		ProjectID:   "ab657ae7-8c9e-45eb-9862-dd9ceb6d5c71",
		Name:        "resc1",
		NamespaceID: "ns1/kind",
	}, {
		ProjectID:   "ab657ae7-8c9e-45eb-9862-dd9ceb6d5c71",
		Name:        "resc2",
		NamespaceID: "ns1/kind",
	},
}

func TestCreateResources(t *testing.T) {
	t.Parallel()
	table := []struct {
		title                string
		permissionAttributes map[string]any
		a                    Authz
		want                 []resource.Resource
		err                  error
	}{
		{
			title:                "success/should return multiple resources",
			permissionAttributes: testPermissionAttributesMap,
			a:                    Authz{},
			want:                 expectedResources,
			err:                  nil,
		}, {
			title: "should should throw error if project is missing",
			permissionAttributes: map[string]any{
				"resource":      []string{"resc1", "resc2"},
				"namespace":     "ns1",
				"resource_type": "kind",
			},
			a:    Authz{},
			want: nil,
			err:  fmt.Errorf("namespace, resource type, projects, resource, and team are required"),
		}, {
			title: "success/should return resource",
			permissionAttributes: map[string]any{
				"project":       "c7772c63-fca4-4c7c-bf93-c8f85115de4b",
				"resource":      "res1",
				"namespace":     "ns1",
				"resource_type": "type",
			},
			a: Authz{},
			want: []resource.Resource{
				{
					ProjectID:   "c7772c63-fca4-4c7c-bf93-c8f85115de4b",
					Name:        "res1",
					NamespaceID: "ns1/type",
				},
			},
			err: nil,
		},
	}

	for _, tt := range table {
		tt := tt
		t.Run(tt.title, func(t *testing.T) {
			t.Parallel()

			resp, err := tt.a.createResources(tt.permissionAttributes)
			assert.EqualValues(t, tt.want, resp)
			assert.EqualValues(t, tt.err, err)
		})
	}
}

func TestServeHook(t *testing.T) {
	var (
		mockRelationService     = new(mocks.RelationService)
		mockResourceService     = new(mocks.ResourceService)
		mockRelationTransformer = new(mocks.RelationTransformer)
	)

	logger := shieldlogger.InitLogger(shieldlogger.Config{Level: "debug"})

	rootHook := hook.New()
	a := New(logger, rootHook, rootHook, mockResourceService, mockRelationService, mockRelationTransformer, "X-Shield-Email")

	t.Run("should return InternalServerError when non-nil error is sent", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "http://localhost:8080", nil)

		response := &http.Response{
			Request: req,
			Header:  http.Header{},
		}

		resp, err := a.ServeHook(response, errors.New("some error"))

		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("should return StatusBadRequest when response has status code 400", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "http://localhost:8080", nil)

		response := &http.Response{
			Request: req,
			Header:  http.Header{},
		}
		response.StatusCode = 400

		res, err := a.ServeHook(response, nil)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	t.Run("should not change status code if rule is not set", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "http://localhost:8080", nil)

		response := &http.Response{
			Request: req,
			Header:  http.Header{},
		}
		response.StatusCode = 200

		resp, err := a.ServeHook(response, nil)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("should return StatusInternalServerError if rule config are not set", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "http://localhost:8080", nil)

		response := &http.Response{
			Request: req,
			Header:  http.Header{},
		}

		rl := &rule.Rule{
			Hooks: rule.HookSpecs{
				rule.HookSpec{
					Name: "authz",
				},
			},
		}
		*response.Request = *response.Request.WithContext(rule.WithContext(req.Context(), rl))
		resp, err := a.ServeHook(response, nil)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("should return InternalServerError if backend namespace is empty", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "http://localhost:8080", nil)

		response := &http.Response{
			Request: req,
			Header:  http.Header{},
		}

		rl := &rule.Rule{
			Hooks: rule.HookSpecs{
				rule.HookSpec{
					Name:   "authz",
					Config: map[string]interface{}{},
				},
			},
			Backend: rule.Backend{
				Namespace: "",
			},
		}

		*response.Request = *response.Request.WithContext(rule.WithContext(req.Context(), rl))

		resp, err := a.ServeHook(response, nil)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("should return InternalServerError if identityProxyHeaderKey not set", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "http://localhost:8080", nil)
		response := &http.Response{
			Request: req,
			Header:  http.Header{},
		}
		rl := &rule.Rule{
			Hooks: rule.HookSpecs{
				rule.HookSpec{
					Name:   "authz",
					Config: map[string]interface{}{},
				},
			},
			Backend: rule.Backend{
				Namespace: "ns1",
			},
		}
		*response.Request = *response.Request.WithContext(rule.WithContext(req.Context(), rl))
		resp, err := a.ServeHook(response, nil)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("should return InternalServerError if all attributes are not set", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "http://localhost:8080", nil)

		response := &http.Response{
			Request: req,
			Header:  http.Header{},
		}

		rl := &rule.Rule{
			Hooks: rule.HookSpecs{
				rule.HookSpec{
					Name:   "authz",
					Config: map[string]interface{}{},
				},
			},
			Backend: rule.Backend{
				Namespace: "ns1",
			},
		}

		*response.Request = *response.Request.WithContext(rule.WithContext(req.Context(), rl))

		response.Request.Header.Set("X-Shield-Email", "user@gotocompany.com")

		resp, err := a.ServeHook(response, nil)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("should not change status code if all attributes are set", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "http://localhost:8080", nil)

		response := &http.Response{
			Request: req,
			Header:  http.Header{},
		}
		response.StatusCode = 200

		rl := &rule.Rule{
			Hooks: rule.HookSpecs{
				rule.HookSpec{
					Name: "authz",
					Config: map[string]interface{}{
						"attributes": map[string]attribute.Attribute{
							"project": {
								Type:  "constant",
								Value: testPermissionAttributesMap["project"].(string),
							},
							"resource": {
								Type:  "constant",
								Value: testPermissionAttributesMap["resource"].([]string)[0],
							},
							"namespace": {
								Type:  "constant",
								Value: testPermissionAttributesMap["namespace"].(string),
							},
							"resource_type": {
								Type:  "constant",
								Value: testPermissionAttributesMap["resource_type"].(string),
							},
						},
					},
				},
			},
			Backend: rule.Backend{
				Namespace: "ns1",
			},
		}

		*response.Request = *response.Request.WithContext(rule.WithContext(req.Context(), rl))

		response.Request.Header.Set("X-Shield-Email", "user@gotocompany.com")
		response.Request.Header.Set("organization", "org1")

		rsc := resource.Resource{
			Name:        testPermissionAttributesMap["resource"].([]string)[0],
			ProjectID:   testPermissionAttributesMap["project"].(string),
			NamespaceID: namespace.CreateID(testPermissionAttributesMap["namespace"].(string), testPermissionAttributesMap["resource_type"].(string)),
		}

		mockResourceService.EXPECT().Upsert(mock.AnythingOfType("*context.valueCtx"), rsc).Return(resource.Resource{
			Idxa:           uuid.NewString(),
			URN:            "new-resource-urn",
			ProjectID:      rsc.ProjectID,
			OrganizationID: rsc.OrganizationID,
			NamespaceID:    rsc.NamespaceID,
			UserID:         "user@gotocompany.com",
			Name:           rsc.Name,
			CreatedAt:      time.Time{},
			UpdatedAt:      time.Time{},
		}, nil)

		resp, err := a.ServeHook(response, nil)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("should not change status code if relations are set", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "http://localhost:8080", nil)
		body := io.NopCloser(bytes.NewBuffer([]byte(`{"foo" : "bar"}`)))

		response := &http.Response{
			Request: req,
			Header:  http.Header{},
			Body:    body,
		}
		response.StatusCode = 200

		rl := &rule.Rule{
			Hooks: rule.HookSpecs{
				rule.HookSpec{
					Name: "authz",
					Config: map[string]interface{}{
						"attributes": map[string]attribute.Attribute{
							"project": {
								Type:  "constant",
								Value: testPermissionAttributesMap["project"].(string),
							},
							"resource": {
								Type: "json_payload",
								Key:  "foo",
							},
							"namespace": {
								Type:  "constant",
								Value: testPermissionAttributesMap["namespace"].(string),
							},
							"resource_type": {
								Type:  "constant",
								Value: testPermissionAttributesMap["resource_type"].(string),
							},
							"group": {
								Type:  "constant",
								Value: testPermissionAttributesMap["group"].(string),
							},
							"user": {
								Type:  "constant",
								Value: testPermissionAttributesMap["user"].(string),
							},
						},
						"relations": []Relation{
							{
								Role:               "owner",
								SubjectPrincipal:   "group",
								SubjectIDAttribute: "group",
							},
							{
								Role:               "owner",
								SubjectPrincipal:   "user",
								SubjectIDAttribute: "user",
							},
						},
					},
				},
			},
			Backend: rule.Backend{
				Namespace: "ns1",
			},
		}

		*response.Request = *response.Request.WithContext(rule.WithContext(req.Context(), rl))

		response.Request.Header.Set("X-Shield-Email", "user@gotocompany.com")
		response.Request.Header.Set("organization", "org1")

		rsc := resource.Resource{
			Name:        "bar",
			ProjectID:   testPermissionAttributesMap["project"].(string),
			NamespaceID: namespace.CreateID(testPermissionAttributesMap["namespace"].(string), testPermissionAttributesMap["resource_type"].(string)),
		}

		mockRelationTransformer.EXPECT().TransformRelation(mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("relation.RelationV2")).Return(relation.RelationV2{}, nil)

		mockResourceService.EXPECT().Upsert(mock.AnythingOfType("*context.valueCtx"), rsc).Return(resource.Resource{
			Idxa:           uuid.NewString(),
			URN:            "new-resource-urn",
			ProjectID:      rsc.ProjectID,
			OrganizationID: rsc.OrganizationID,
			NamespaceID:    rsc.NamespaceID,
			UserID:         "user@gotocompany.com",
			Name:           "bar",
			CreatedAt:      time.Time{},
			UpdatedAt:      time.Time{},
		}, nil)

		mockRelationService.EXPECT().Create(mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("relation.RelationV2")).Return(
			relation.RelationV2{}, nil)

		resp, err := a.ServeHook(response, nil)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("should throw internal server error when header type attributes is missing", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "http://localhost:8080", nil)
		body := io.NopCloser(bytes.NewBuffer([]byte(`{"foo": "bar"}`)))

		response := &http.Response{
			Request: req,
			Header:  http.Header{},
			Body:    body,
		}
		response.StatusCode = 200

		rl := &rule.Rule{
			Hooks: rule.HookSpecs{
				rule.HookSpec{
					Name: "authz",
					Config: map[string]interface{}{
						"attributes": map[string]attribute.Attribute{
							"project": {
								Type:  "constant",
								Value: testPermissionAttributesMap["project"].(string),
							},
							"organization": {
								Type:   "header",
								Key:    "organization",
								Source: "request",
							},
							"resource": {
								Type: "json_payload",
								Key:  "foo",
							},
							"namespace": {
								Type:  "constant",
								Value: testPermissionAttributesMap["namespace"].(string),
							},
							"resource_type": {
								Type:  "constant",
								Value: testPermissionAttributesMap["resource_type"].(string),
							},
							"group": {
								Type:  "constant",
								Value: testPermissionAttributesMap["group"].(string),
							},
							"user": {
								Type:  "constant",
								Value: testPermissionAttributesMap["user"].(string),
							},
						},
						"relations": []Relation{
							{
								Role:               "owner",
								SubjectPrincipal:   "group",
								SubjectIDAttribute: "group",
							},
							{
								Role:               "owner",
								SubjectPrincipal:   "user",
								SubjectIDAttribute: "user",
							},
						},
					},
				},
			},
			Backend: rule.Backend{
				Namespace: "ns1",
			},
		}

		*response.Request = *response.Request.WithContext(rule.WithContext(req.Context(), rl))

		response.Request.Header.Set("X-Shield-Email", "user@gotocompany.com")

		rsc := resource.Resource{
			Name:           "bar",
			OrganizationID: testPermissionAttributesMap["organization"].(string),
			ProjectID:      testPermissionAttributesMap["project"].(string),
			NamespaceID:    namespace.CreateID(testPermissionAttributesMap["namespace"].(string), testPermissionAttributesMap["resource_type"].(string)),
		}

		mockResourceService.EXPECT().Upsert(mock.AnythingOfType("*context.valueCtx"), rsc).Return(resource.Resource{
			Idxa:           uuid.NewString(),
			URN:            "new-resource-urn",
			ProjectID:      rsc.ProjectID,
			OrganizationID: rsc.OrganizationID,
			NamespaceID:    rsc.NamespaceID,
			UserID:         "user@gotocompany.com",
			Name:           "bar",
			CreatedAt:      time.Time{},
			UpdatedAt:      time.Time{},
		}, nil)

		mockRelationService.EXPECT().Create(mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("relation.RelationV2")).Return(
			relation.RelationV2{}, nil)

		resp, err := a.ServeHook(response, nil)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("should throw internal server error when json_payload type attributes is missing", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "http://localhost:8080", nil)
		body := io.NopCloser(bytes.NewBuffer([]byte(`{}`)))

		response := &http.Response{
			Request: req,
			Header:  http.Header{},
			Body:    body,
		}
		response.StatusCode = 200

		rl := &rule.Rule{
			Hooks: rule.HookSpecs{
				rule.HookSpec{
					Name: "authz",
					Config: map[string]interface{}{
						"attributes": map[string]attribute.Attribute{
							"project": {
								Type:  "constant",
								Value: testPermissionAttributesMap["project"].(string),
							},
							"organization": {
								Type:   "header",
								Key:    "organization",
								Source: "request",
							},
							"resource": {
								Type: "json_payload",
								Key:  "foo",
							},
							"namespace": {
								Type:  "constant",
								Value: testPermissionAttributesMap["namespace"].(string),
							},
							"resource_type": {
								Type:  "constant",
								Value: testPermissionAttributesMap["resource_type"].(string),
							},
							"group": {
								Type:  "constant",
								Value: testPermissionAttributesMap["group"].(string),
							},
							"user": {
								Type:  "constant",
								Value: testPermissionAttributesMap["user"].(string),
							},
						},
						"relations": []Relation{
							{
								Role:               "owner",
								SubjectPrincipal:   "group",
								SubjectIDAttribute: "group",
							},
							{
								Role:               "owner",
								SubjectPrincipal:   "user",
								SubjectIDAttribute: "user",
							},
						},
					},
				},
			},
			Backend: rule.Backend{
				Namespace: "ns1",
			},
		}

		*response.Request = *response.Request.WithContext(rule.WithContext(req.Context(), rl))

		response.Request.Header.Set("X-Shield-Email", "user@gotocompany.com")
		response.Request.Header.Set("organization", "org1")

		rsc := resource.Resource{
			Name:           "bar",
			OrganizationID: testPermissionAttributesMap["organization"].(string),
			ProjectID:      testPermissionAttributesMap["project"].(string),
			NamespaceID:    namespace.CreateID(testPermissionAttributesMap["namespace"].(string), testPermissionAttributesMap["resource_type"].(string)),
		}

		mockResourceService.EXPECT().Upsert(mock.AnythingOfType("*context.valueCtx"), rsc).Return(resource.Resource{
			Idxa:           uuid.NewString(),
			URN:            "new-resource-urn",
			ProjectID:      rsc.ProjectID,
			OrganizationID: rsc.OrganizationID,
			NamespaceID:    rsc.NamespaceID,
			UserID:         "user@gotocompany.com",
			Name:           "bar",
			CreatedAt:      time.Time{},
			UpdatedAt:      time.Time{},
		}, nil)

		mockRelationService.EXPECT().Create(mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("relation.RelationV2")).Return(
			relation.RelationV2{}, nil)

		resp, err := a.ServeHook(response, nil)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("should throw internal server error when constant type attributes is missing", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "http://localhost:8080", nil)
		body := io.NopCloser(bytes.NewBuffer([]byte(`{"foo": "bar"}`)))

		response := &http.Response{
			Request: req,
			Header:  http.Header{},
			Body:    body,
		}
		response.StatusCode = 200

		rl := &rule.Rule{
			Hooks: rule.HookSpecs{
				rule.HookSpec{
					Name: "authz",
					Config: map[string]interface{}{
						"attributes": map[string]attribute.Attribute{
							"project": {
								Type:  "constant",
								Value: testPermissionAttributesMap["project"].(string),
							},
							"organization": {
								Type:   "header",
								Key:    "organization",
								Source: "request",
							},
							"resource": {
								Type: "json_payload",
								Key:  "foo",
							},
							"namespace": {
								Type:  "constant",
								Value: testPermissionAttributesMap["namespace"].(string),
							},
							"resource_type": {
								Type:  "constant",
								Value: testPermissionAttributesMap["resource_type"].(string),
							},
							"group": {
								Type:  "constant",
								Value: testPermissionAttributesMap["group"].(string),
							},
							"user": {
								Type: "constant",
							},
						},
						"relations": []Relation{
							{
								Role:               "owner",
								SubjectPrincipal:   "group",
								SubjectIDAttribute: "group",
							},
							{
								Role:               "owner",
								SubjectPrincipal:   "user",
								SubjectIDAttribute: "user",
							},
						},
					},
				},
			},
			Backend: rule.Backend{
				Namespace: "ns1",
			},
		}

		*response.Request = *response.Request.WithContext(rule.WithContext(req.Context(), rl))

		response.Request.Header.Set("X-Shield-Email", "user@gotocompany.com")
		response.Request.Header.Set("organization", "org1")

		rsc := resource.Resource{
			Name:           "bar",
			OrganizationID: testPermissionAttributesMap["organization"].(string),
			ProjectID:      testPermissionAttributesMap["project"].(string),
			NamespaceID:    namespace.CreateID(testPermissionAttributesMap["namespace"].(string), testPermissionAttributesMap["resource_type"].(string)),
		}

		mockResourceService.EXPECT().Upsert(mock.AnythingOfType("*context.valueCtx"), rsc).Return(resource.Resource{
			Idxa:           uuid.NewString(),
			URN:            "new-resource-urn",
			ProjectID:      rsc.ProjectID,
			OrganizationID: rsc.OrganizationID,
			NamespaceID:    rsc.NamespaceID,
			UserID:         "user@gotocompany.com",
			Name:           "bar",
			CreatedAt:      time.Time{},
			UpdatedAt:      time.Time{},
		}, nil)

		mockRelationService.EXPECT().Create(mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("relation.RelationV2")).Return(
			relation.RelationV2{}, nil)

		resp, err := a.ServeHook(response, nil)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
