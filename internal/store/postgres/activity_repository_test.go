package postgres_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/goto/salt/audit"
	"github.com/goto/salt/log"
	"github.com/goto/shield/core/activity"
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
	activities []audit.Log
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

func (s *ActivityRepositoryTestSuite) SetupTest() {
	var err error
	s.activities, err = bootstrapActivity(s.client)
	if err != nil {
		s.T().Fatal(err)
	}
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

func (s *ActivityRepositoryTestSuite) TestList() {
	type testCase struct {
		Description        string
		Filter             activity.Filter
		ExpectedActivities []audit.Log
		ErrString          string
	}

	var testCases = []testCase{
		{
			Description:        "should get all activity",
			ExpectedActivities: s.activities,
		},
		{
			Description: "should return maximum 50 activities if limit is set < 1",
			Filter: activity.Filter{
				Limit: 0,
			},
			ExpectedActivities: s.activities,
		},
		{
			Description: "should return maximum specified activities if limit is set >= 1",
			Filter: activity.Filter{
				Limit: 2,
			},
			ExpectedActivities: []audit.Log{
				s.activities[0], s.activities[1],
			},
		},
		{
			Description: "should return the first page if page is set < 1",
			Filter: activity.Filter{
				Page: 0,
			},
			ExpectedActivities: s.activities,
		},
		{
			Description: "should return the specified page if page is set >= 1",
			Filter: activity.Filter{
				Page:  2,
				Limit: 1,
			},
			ExpectedActivities: []audit.Log{s.activities[1]},
		},
		{
			Description: "should return activities between start time and end time",
			Filter: activity.Filter{
				StartTime: time.Date(2024, time.January, 10, 15, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2024, time.January, 10, 25, 0, 0, 0, time.UTC),
			},
			ExpectedActivities: []audit.Log{s.activities[1]},
		},
		{
			Description: "should return filtered activities by actor",
			Filter: activity.Filter{
				Actor: s.activities[0].Actor,
			},
			ExpectedActivities: []audit.Log{s.activities[0]},
		},
		{
			Description: "should return filtered activity by action",
			Filter: activity.Filter{
				Action: s.activities[1].Action,
			},
			ExpectedActivities: []audit.Log{s.activities[1]},
		},
		{
			Description: "should return filtered activity by data",
			Filter: activity.Filter{
				Data: map[string]string{
					"entity": "user",
				},
			},
			ExpectedActivities: []audit.Log{s.activities[0]},
		},
		{
			Description: "should return filtered activity by metadata",
			Filter: activity.Filter{
				Metadata: map[string]string{
					"version": "v0.1.3",
				},
			},
			ExpectedActivities: []audit.Log{s.activities[2]},
		},
		{
			Description: "should return empty activities if all filtered out",
			Filter: activity.Filter{
				Actor:     s.activities[0].Actor,
				Action:    s.activities[1].Action,
				StartTime: time.Date(2024, time.January, 20, 0, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2024, time.January, 30, 0, 0, 0, 0, time.UTC),
				Data:      map[string]string{"entity": "project"},
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			got, err := s.repository.List(s.ctx, tc.Filter)
			if tc.ErrString != "" {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
			if !cmp.Equal(got, tc.ExpectedActivities) {
				s.T().Fatalf("got result %+v, expected was %+v", got, tc.ExpectedActivities)
			}
		})
	}
}

func TestActivityRepository(t *testing.T) {
	suite.Run(t, new(ActivityRepositoryTestSuite))
}
