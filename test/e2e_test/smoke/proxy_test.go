package e2e_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/goto/shield/config"
	"github.com/goto/shield/pkg/db"
	shieldv1beta1 "github.com/goto/shield/proto/v1beta1"
	"github.com/goto/shield/test/e2e_test/testbench"
)

type EndToEndProxySmokeTestSuite struct {
	suite.Suite
	userID       string
	orgID        string
	orgSlug      string
	projID       string
	projSlug     string
	groupID      string
	client       shieldv1beta1.ShieldServiceClient
	cancelClient func()
	testBench    *testbench.TestBench
	dbClient     *db.Client
	appConfig    *config.Shield
}

func (s *EndToEndProxySmokeTestSuite) SetupTest() {
	ctx := context.Background()
	s.client, s.appConfig, s.cancelClient, _ = testbench.SetupTests(s.T())

	dbClient, err := testbench.SetupDB(s.appConfig.DB)
	if err != nil {
		s.T().Fatal("failed to setup database")
		return
	}
	s.dbClient = dbClient

	// validate
	uRes, err := s.client.ListUsers(ctx, &shieldv1beta1.ListUsersRequest{})
	s.Require().NoError(err)
	s.Require().Equal(9, len(uRes.GetUsers()))
	s.userID = uRes.GetUsers()[0].GetId()

	oRes, err := s.client.ListOrganizations(ctx, &shieldv1beta1.ListOrganizationsRequest{})
	s.Require().NoError(err)
	s.Require().Equal(1, len(oRes.GetOrganizations()))
	s.orgID = oRes.GetOrganizations()[0].GetId()
	s.orgSlug = oRes.GetOrganizations()[0].GetSlug()

	pRes, err := s.client.ListProjects(ctx, &shieldv1beta1.ListProjectsRequest{})
	s.Require().NoError(err)
	s.Require().Equal(1, len(pRes.GetProjects()))
	s.projID = pRes.GetProjects()[0].GetId()
	s.projSlug = pRes.GetProjects()[0].GetSlug()

	gRes, err := s.client.ListGroups(ctx, &shieldv1beta1.ListGroupsRequest{})
	s.Require().NoError(err)
	s.Require().Equal(3, len(gRes.GetGroups()))
	s.groupID = gRes.GetGroups()[0].GetId()
}

func (s *EndToEndProxySmokeTestSuite) TearDownTest() {
	s.cancelClient()
	// Clean tests
	err := s.testBench.CleanUp()
	s.Require().NoError(err)
}

func (s *EndToEndProxySmokeTestSuite) TestProxyToEchoServer() {
	s.Run("should return unauthenticated error if user is not registered", func() {
		url := fmt.Sprintf("http://localhost:%d/api/ping", s.appConfig.Proxy.Services[0].Port)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		s.Require().NoError(err)

		req.Header.Set(testbench.IdentityHeader, "john.doe@gotocompany.com")

		res, err := http.DefaultClient.Do(req)
		s.Require().NoError(err)

		defer res.Body.Close()
		s.Assert().Equal(401, res.StatusCode)
	})
	s.Run("should be able to proxy to an echo server", func() {
		url := fmt.Sprintf("http://localhost:%d/api/ping", s.appConfig.Proxy.Services[0].Port)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		s.Require().NoError(err)

		req.Header.Set(testbench.IdentityHeader, "member2-group1@gotocompany.com")

		res, err := http.DefaultClient.Do(req)
		s.Require().NoError(err)

		defer res.Body.Close()
		s.Assert().Equal(200, res.StatusCode)
	})
	s.Run("resource created on echo server should persist in shieldDB", func() {
		url := fmt.Sprintf("http://localhost:%d/api/resource", s.appConfig.Proxy.Services[0].Port)
		reqBodyMap := map[string]string{
			"project": s.projID,
			"name":    "test-resource",
			"group":   s.groupID,
		}
		reqBodyBytes, err := json.Marshal(reqBodyMap)
		s.Require().NoError(err)

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBodyBytes))
		s.Require().NoError(err)

		req.Header.Set(testbench.IdentityHeader, "member2-group1@gotocompany.com")
		req.Header.Set("X-Shield-Org", s.orgID)

		res, err := http.DefaultClient.Do(req)
		s.Require().NoError(err)

		defer res.Body.Close()

		resourceSelectQuery := "SELECT name FROM resources"
		resources, err := s.dbClient.DB.Query(resourceSelectQuery)
		s.Require().NoError(err)
		defer resources.Close()

		var resourceName = ""
		for resources.Next() {
			if err := resources.Scan(&resourceName); err != nil {
				s.Require().NoError(err)
			}
		}
		s.Assert().Equal(200, res.StatusCode)
		s.Assert().Equal("test-resource", resourceName)
	})

	s.Run("user not part of group will not be authenticated by middleware auth", func() {
		groupDetail, err := s.client.GetGroup(context.Background(), &shieldv1beta1.GetGroupRequest{Id: s.groupID})
		s.Require().NoError(err)

		url := fmt.Sprintf("http://localhost:%d/api/resource_slug", s.appConfig.Proxy.Services[0].Port)
		reqBodyMap := map[string]string{
			"project":    s.projID,
			"name":       "test-resource-group-slug",
			"group_slug": groupDetail.GetGroup().GetSlug(),
		}
		reqBodyBytes, err := json.Marshal(reqBodyMap)
		s.Require().NoError(err)

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBodyBytes))
		s.Require().NoError(err)

		req.Header.Set(testbench.IdentityHeader, "member2-group1@gotocompany.com")
		req.Header.Set("X-Shield-Org", s.orgID)

		res, err := http.DefaultClient.Do(req)
		s.Require().NoError(err)

		defer res.Body.Close()

		s.Assert().Equal(401, res.StatusCode)
	})

	s.Run("permission expression: user not having permission at proj level will not be authenticated by middleware auth", func() {
		url := fmt.Sprintf("http://localhost:%d/api/create_firehose_based_on_sink", s.appConfig.Proxy.Services[0].Port)
		reqBodyMap := map[string]any{
			"organization": s.orgID,
			"project":      s.projID,
			"configs": map[string]any{
				"env_vars": map[string]any{
					"SINK_TYPE": "bigquery",
				},
			},
		}
		reqBodyBytes, err := json.Marshal(reqBodyMap)
		s.Require().NoError(err)

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBodyBytes))
		s.Require().NoError(err)

		req.Header.Set(testbench.IdentityHeader, "member2-group1@gotocompany.com")

		res, err := http.DefaultClient.Do(req)
		s.Require().NoError(err)

		defer res.Body.Close()
		s.Assert().Equal(401, res.StatusCode)
	})

	s.Run("permission expression: user not having permission at org level will not be authenticated by middleware auth", func() {
		url := fmt.Sprintf("http://localhost:%d/api/create_firehose_based_on_sink", s.appConfig.Proxy.Services[0].Port)
		reqBodyMap := map[string]any{
			"organization": s.orgID,
			"project":      s.projID,
			"configs": map[string]any{
				"env_vars": map[string]any{
					"SINK_TYPE": "blob",
				},
			},
		}
		reqBodyBytes, err := json.Marshal(reqBodyMap)
		s.Require().NoError(err)

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBodyBytes))
		s.Require().NoError(err)

		req.Header.Set(testbench.IdentityHeader, "member2-group1@gotocompany.com")

		res, err := http.DefaultClient.Do(req)
		s.Require().NoError(err)

		defer res.Body.Close()
		s.Assert().Equal(401, res.StatusCode)
	})

	s.Run("permission expression: user not having permission at org level will not be authenticated by middleware auth with org passed as slug", func() {
		url := fmt.Sprintf("http://localhost:%d/api/create_firehose_based_on_sink", s.appConfig.Proxy.Services[0].Port)
		reqBodyMap := map[string]any{
			"organization": s.orgSlug,
			"project":      s.projSlug,
			"configs": map[string]any{
				"env_vars": map[string]any{
					"SINK_TYPE": "blob",
				},
			},
		}
		reqBodyBytes, err := json.Marshal(reqBodyMap)
		s.Require().NoError(err)

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBodyBytes))
		s.Require().NoError(err)

		req.Header.Set(testbench.IdentityHeader, "member2-group1@gotocompany.com")

		res, err := http.DefaultClient.Do(req)
		s.Require().NoError(err)

		defer res.Body.Close()
		s.Assert().Equal(401, res.StatusCode)
	})

	s.Run("permission expression: user not having permission at proj level will not be authenticated by middleware auth with proj passed as slug", func() {
		url := fmt.Sprintf("http://localhost:%d/api/create_firehose_based_on_sink", s.appConfig.Proxy.Services[0].Port)
		reqBodyMap := map[string]any{
			"organization": s.orgSlug,
			"project":      s.projSlug,
			"configs": map[string]any{
				"env_vars": map[string]any{
					"SINK_TYPE": "bigquery",
				},
			},
		}
		reqBodyBytes, err := json.Marshal(reqBodyMap)
		s.Require().NoError(err)

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBodyBytes))
		s.Require().NoError(err)

		req.Header.Set(testbench.IdentityHeader, "member2-group1@gotocompany.com")

		res, err := http.DefaultClient.Do(req)
		s.Require().NoError(err)

		defer res.Body.Close()
		s.Assert().Equal(401, res.StatusCode)
	})

	s.Run("resource created on echo server should persist in shieldDB when using group slug", func() {
		groupDetail, err := s.client.GetGroup(context.Background(), &shieldv1beta1.GetGroupRequest{Id: s.groupID})
		s.Require().NoError(err)

		url := fmt.Sprintf("http://localhost:%d/api/resource_slug", s.appConfig.Proxy.Services[0].Port)
		reqBodyMap := map[string]string{
			"project":    s.projID,
			"name":       "test-resource-group-slug",
			"group_slug": groupDetail.GetGroup().GetSlug(),
		}
		reqBodyBytes, err := json.Marshal(reqBodyMap)
		s.Require().NoError(err)

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBodyBytes))
		s.Require().NoError(err)

		req.Header.Set(testbench.IdentityHeader, "admin1-group1-org1@gotocompany.com")
		req.Header.Set("X-Shield-Org", s.orgID)

		res, err := http.DefaultClient.Do(req)
		s.Require().NoError(err)

		defer res.Body.Close()

		s.Assert().Equal(200, res.StatusCode)

		resourceSelectQuery := "SELECT name FROM resources"
		resources, err := s.dbClient.DB.Query(resourceSelectQuery)
		s.Require().NoError(err)
		defer resources.Close()

		var resourceName = ""
		for resources.Next() {
			if err := resources.Scan(&resourceName); err != nil {
				s.Require().NoError(err)
			}
		}
		s.Assert().Equal("test-resource-group-slug", resourceName)

		relationSelectQuery := "SELECT subject_id FROM relations ORDER BY created_at DESC LIMIT 1"
		relations, err := s.dbClient.DB.Query(relationSelectQuery)
		s.Require().NoError(err)
		defer resources.Close()

		var subjectID = ""
		for relations.Next() {
			if err := relations.Scan(&subjectID); err != nil {
				s.Require().NoError(err)
			}
		}
		s.Assert().Equal(s.groupID, subjectID)
	})
	s.Run("resource created on echo server should persist in shieldDB when using user id", func() {
		url := fmt.Sprintf("http://localhost:%d/api/resource_user_id", s.appConfig.Proxy.Services[0].Port)
		reqBodyMap := map[string]string{
			"project": s.projID,
			"name":    "test-resource-user-id",
			"user_id": s.userID,
		}
		reqBodyBytes, err := json.Marshal(reqBodyMap)
		s.Require().NoError(err)

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBodyBytes))
		s.Require().NoError(err)

		req.Header.Set(testbench.IdentityHeader, "member2-group1@gotocompany.com")
		req.Header.Set("X-Shield-Org", s.orgID)

		res, err := http.DefaultClient.Do(req)
		s.Require().NoError(err)

		defer res.Body.Close()

		s.Assert().Equal(200, res.StatusCode)

		resourceSelectQuery := "SELECT name FROM resources"
		resources, err := s.dbClient.DB.Query(resourceSelectQuery)
		s.Require().NoError(err)
		defer resources.Close()

		var resourceName = ""
		for resources.Next() {
			if err := resources.Scan(&resourceName); err != nil {
				s.Require().NoError(err)
			}
		}
		s.Assert().Equal("test-resource-user-id", resourceName)

		relationSelectQuery := "SELECT subject_id FROM relations ORDER BY created_at DESC LIMIT 1"
		relations, err := s.dbClient.DB.Query(relationSelectQuery)
		s.Require().NoError(err)
		defer resources.Close()

		var subjectID = ""
		for relations.Next() {
			if err := relations.Scan(&subjectID); err != nil {
				s.Require().NoError(err)
			}
		}
		s.Assert().Equal(s.userID, subjectID)
	})
	s.Run("resource created on echo server should persist in shieldDB when using user e-mail", func() {
		userDetail, err := s.client.GetUser(context.Background(), &shieldv1beta1.GetUserRequest{Id: s.userID})
		s.Require().NoError(err)

		url := fmt.Sprintf("http://localhost:%d/api/resource_user_email", s.appConfig.Proxy.Services[0].Port)
		reqBodyMap := map[string]string{
			"project":    s.projID,
			"name":       "test-resource-user-email",
			"user_email": userDetail.GetUser().GetEmail(),
		}
		reqBodyBytes, err := json.Marshal(reqBodyMap)
		s.Require().NoError(err)

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBodyBytes))
		s.Require().NoError(err)

		req.Header.Set(testbench.IdentityHeader, "member2-group1@gotocompany.com")
		req.Header.Set("X-Shield-Org", s.orgID)

		res, err := http.DefaultClient.Do(req)
		s.Require().NoError(err)

		defer res.Body.Close()

		s.Assert().Equal(200, res.StatusCode)

		resourceSelectQuery := "SELECT name FROM resources"
		resources, err := s.dbClient.DB.Query(resourceSelectQuery)
		s.Require().NoError(err)
		defer resources.Close()

		var resourceName = ""
		for resources.Next() {
			if err := resources.Scan(&resourceName); err != nil {
				s.Require().NoError(err)
			}
		}
		s.Assert().Equal("test-resource-user-email", resourceName)

		relationSelectQuery := "SELECT subject_id FROM relations ORDER BY created_at DESC LIMIT 1"
		relations, err := s.dbClient.DB.Query(relationSelectQuery)
		s.Require().NoError(err)
		defer resources.Close()

		var subjectID = ""
		for relations.Next() {
			if err := relations.Scan(&subjectID); err != nil {
				s.Require().NoError(err)
			}
		}
		s.Assert().Equal(s.userID, subjectID)
	})
}

func TestEndToEndProxySmokeTestSuite(t *testing.T) {
	suite.Run(t, new(EndToEndProxySmokeTestSuite))
}
