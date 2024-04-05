package postgres_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/goto/salt/audit"
	"github.com/goto/salt/log"
	"github.com/goto/shield/internal/store/postgres"
	"github.com/goto/shield/pkg/db"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/suite"
)

type ActivityRepositoryTestSuite struct {
	suite.Suite
	ctx        context.Context
	client     *db.Client
	pool       *dockertest.Pool
	resource   *dockertest.Resource
	repository *postgres.ActivityRepository
}

func (s *ActivityRepositoryTestSuite) SetupSuite() {
	var err error

	logger := log.NewZap()
	s.client, s.pool, s.resource, err = newTestClient(logger)
	if err != nil {
		s.T().Fatal(err)
	}

	s.ctx = context.TODO()
	s.repository = postgres.NewActivityRepository(s.client)
}

func (s *ActivityRepositoryTestSuite) TearDownSuite() {
	if err := purgeDocker(s.pool, s.resource); err != nil {
		s.T().Fatal(err)
	}
}

func (s *ActivityRepositoryTestSuite) TearDownTest() {
	if err := s.cleanup(); err != nil {
		s.T().Fatal(err)
	}
}

func (s *ActivityRepositoryTestSuite) cleanup() error {
	queries := []string{
		fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", postgres.TABLE_ACTIVITY),
	}
	return execQueries(context.TODO(), s.client, queries)
}

func (s *ActivityRepositoryTestSuite) TestInsert() {
	type testCase struct {
		Description string
		LogToCreate *audit.Log
		ErrString   string
	}

	var testCases = []testCase{
		{
			Description: "should insert a log",
			LogToCreate: &audit.Log{
				Actor:  "0000000a-0000-00aa-0aa0-aaa00a000aa0",
				Action: "group.update",
				Data: map[string]string{
					"entity": "group",
					"id":     "1234-5678-1234",
				},
				Metadata: map[string]string{
					"app_name":    "shield",
					"app_version": "v1.0",
				},
				Timestamp: time.Now(),
			},
		},
		{
			Description: "should return error if metadata is invalid",
			LogToCreate: &audit.Log{
				Actor:  "0000000a-0000-00aa-0aa0-aaa00a000aa0",
				Action: "group.update",
				Data: map[string]string{
					"entity": "group",
					"id":     "1234-5678-1234",
				},
				Metadata:  make(chan int),
				Timestamp: time.Now(),
			},
			ErrString: "parsing error: json: unsupported type: chan int",
		},
		{
			Description: "should return error if data is invalid",
			LogToCreate: &audit.Log{
				Actor:     "0000000a-0000-00aa-0aa0-aaa00a000aa0",
				Action:    "group.update",
				Data:      make(chan int),
				Metadata:  map[string]string{},
				Timestamp: time.Now(),
			},
			ErrString: "parsing error: json: unsupported type: chan int",
		},
		{
			Description: "should return error if actor uuid is invalid",
			LogToCreate: &audit.Log{
				Actor:  "invalid-uuid",
				Action: "group.update",
				Data: map[string]string{
					"entity": "group",
					"id":     "1234-5678-1234",
				},
				Metadata: map[string]string{
					"app_name":    "shield",
					"app_version": "v1.0",
				},
				Timestamp: time.Now(),
			},
			ErrString: "invalid syntax of uuid",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			err := s.repository.Insert(s.ctx, tc.LogToCreate)
			if err != nil && tc.ErrString != "" {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
		})
	}
}

func TestActivityRepository(t *testing.T) {
	suite.Run(t, new(ActivityRepositoryTestSuite))
}
