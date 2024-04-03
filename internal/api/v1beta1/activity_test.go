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
	testActorID  = uuid.NewString()
	testActivity = audit.Log{
		Actor:     testActorID,
		Action:    "user.create",
		Timestamp: time.Time{},
	}
	testActivityPB = &shieldv1beta1.Activity{}
)

func TestHandler_ListActivity(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(gs *mocks.ActivityService, us *mocks.UserService)
		request *shieldv1beta1.ListActivitiesRequest
		want    *shieldv1beta1.ListActivitiesResponse
		wantErr error
	}{
		// return error if activity service return error
		// return activities if activity service return nil error
		// return error if request start time format return parse error
		// return error if request end time format return parse error
		// return error if actor uuid is invalid
		// return error if actor email is invalid
		// return error if actor email not found
		{
			name: "should return internal error if activity service return error",
			setup: func(as *mocks.ActivityService, _ *mocks.UserService) {
				as.EXPECT().List(mock.AnythingOfType("*context.emptyCtx"), activity.Filter{}).Return(activity.PagedActivity{}, errors.New("some error"))
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
			request: &shieldv1beta1.ListActivitiesRequest{
				Actor: "invalid-uuid",
			},
			want:    nil,
			wantErr: grpcBadBodyError,
		},
		{
			name: "should return bad request error if email is not found",
			setup: func(_ *mocks.ActivityService, us *mocks.UserService) {
				us.EXPECT().GetByEmail(mock.AnythingOfType("*context.emptyCtx"), testActorID).Return(user.User{}, user.ErrNotExist)
			},
			request: &shieldv1beta1.ListActivitiesRequest{
				Actor: testActorID,
			},
			want:    nil,
			wantErr: grpcBadBodyError,
		},
		{
			name: "should return activities if activity service return none error",
			setup: func(gs *mocks.ActivityService, _ *mocks.UserService) {
				testActivityList := []audit.Log{testActivity}
				gs.EXPECT().List(mock.AnythingOfType("*context.emptyCtx"), activity.Filter{}).Return(
					activity.PagedActivity{
						Count:      int32(len(testActivityList)),
						Activities: testActivityList,
					}, nil)
			},
			request: &shieldv1beta1.ListActivitiesRequest{},
			want: &shieldv1beta1.ListActivitiesResponse{
				Count: int32(len([]*shieldv1beta1.Activity{
					testActivityPB,
				})),
				Activities: []*shieldv1beta1.Activity{testActivityPB},
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
			}
			got, err := h.ListActivities(context.Background(), tt.request)
			assert.EqualValues(t, got, tt.want)
			assert.EqualValues(t, err, tt.wantErr)
		})
	}
}
