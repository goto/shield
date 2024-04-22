package v1beta1

import (
	"context"
	"errors"
	"testing"

	"github.com/goto/shield/core/activity"
	"github.com/goto/shield/internal/api/v1beta1/mocks"
	shieldv1beta1 "github.com/goto/shield/proto/v1beta1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_ListActivity(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(gs *mocks.ActivityService, us *mocks.UserService)
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
			got, err := h.ListActivities(context.TODO(), tt.request)
			assert.EqualValues(t, got, tt.want)
			assert.EqualValues(t, err, tt.wantErr)
		})
	}
}
