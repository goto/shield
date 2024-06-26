package v1beta1

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/goto/salt/audit"
	"github.com/goto/shield/core/activity"
	"github.com/goto/shield/core/user"
	"github.com/goto/shield/internal/api/v1beta1/mocks"
	"github.com/goto/shield/pkg/uuid"
	shieldv1beta1 "github.com/goto/shield/proto/v1beta1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	testActorID    = uuid.NewString()
	testActorEmail = "test@gotocompany.com"
	testActivity   = audit.Log{
		Actor:     testActorID,
		Action:    "user.create",
		Timestamp: time.Time{},
	}
	testActivityPB, _ = transformActivityToPB(testActivity)
)

func TestHandler_ListActivity(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(as *mocks.ActivityService, us *mocks.UserService)
		request *shieldv1beta1.ListActivitiesRequest
		want    *shieldv1beta1.ListActivitiesResponse
		wantErr error
	}{
		{
			name: "should return internal error if activity service return error",
			setup: func(as *mocks.ActivityService, _ *mocks.UserService) {
				as.EXPECT().List(mock.AnythingOfType("context.todoCtx"), activity.Filter{}).Return(activity.PagedActivity{}, errors.New("some error"))
			},
			request: &shieldv1beta1.ListActivitiesRequest{},
			want:    nil,
			wantErr: grpcInternalServerError,
		},
		{
			name: "should return bad request error if start time parsing return error",
			request: &shieldv1beta1.ListActivitiesRequest{
				StartTime: "invalid-start-time",
			},
			want:    nil,
			wantErr: grpcBadBodyError,
		},
		{
			name: "should return bad request error if end time parsing return error",
			request: &shieldv1beta1.ListActivitiesRequest{
				EndTime: "invalid-end-time",
			},
			want:    nil,
			wantErr: grpcBadBodyError,
		},
		{
			name: "should return bad request error if uuid is invalid",
			setup: func(_ *mocks.ActivityService, us *mocks.UserService) {
				us.EXPECT().GetByEmail(mock.AnythingOfType("context.todoCtx"), "invalid-uuid").Return(user.User{}, user.ErrInvalidEmail)
			},
			request: &shieldv1beta1.ListActivitiesRequest{
				Actor: "invalid-uuid",
			},
			want:    nil,
			wantErr: grpcBadBodyError,
		},
		{
			name: "should return bad request error if email is not found",
			setup: func(_ *mocks.ActivityService, us *mocks.UserService) {
				us.EXPECT().GetByEmail(mock.AnythingOfType("context.todoCtx"), testActorEmail).Return(user.User{}, user.ErrNotExist)
			},
			request: &shieldv1beta1.ListActivitiesRequest{
				Actor: testActorEmail,
			},
			want:    nil,
			wantErr: grpcBadBodyError,
		},
		{
			name: "should return activities if activity service return none error",
			setup: func(as *mocks.ActivityService, _ *mocks.UserService) {
				testActivityList := []audit.Log{testActivity}
				as.EXPECT().List(mock.AnythingOfType("context.todoCtx"), activity.Filter{}).Return(
					activity.PagedActivity{
						Count:      int32(len(testActivityList)),
						Activities: testActivityList,
					}, nil)
			},
			request: &shieldv1beta1.ListActivitiesRequest{},
			want: &shieldv1beta1.ListActivitiesResponse{
				Count: int32(len([]*shieldv1beta1.Activity{
					&testActivityPB,
				})),
				Activities: []*shieldv1beta1.Activity{&testActivityPB},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockActivitySvc := new(mocks.ActivityService)
			mockUserSvc := new(mocks.UserService)
			if tt.setup != nil {
				tt.setup(mockActivitySvc, mockUserSvc)
			}
			h := Handler{
				activityService: mockActivitySvc,
				userService:     mockUserSvc,
			}
			got, err := h.ListActivities(context.TODO(), tt.request)
			assert.EqualValues(t, got, tt.want)
			assert.EqualValues(t, err, tt.wantErr)
		})
	}
}
