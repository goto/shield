package v1beta1

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/goto/salt/audit"
	"github.com/goto/shield/core/activity"
	"github.com/goto/shield/core/user"
	"github.com/goto/shield/pkg/uuid"
	shieldv1beta1 "github.com/goto/shield/proto/v1beta1"
	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ActivityService interface {
	List(ctx context.Context, filter activity.Filter) (activity.PagedActivity, error)
}

func (h Handler) ListActivities(ctx context.Context, request *shieldv1beta1.ListActivitiesRequest) (*shieldv1beta1.ListActivitiesResponse, error) {
	logger := grpczap.Extract(ctx)
	var activities []*shieldv1beta1.Activity

	startTime := time.Time{}
	endTime := time.Time{}

	if request.StartTime != "" {
		parseStartTime, err := strconv.ParseInt(request.StartTime, 10, 64)
		if err != nil {
			logger.Error(err.Error())
			return nil, grpcBadBodyError
		}
		startTime = time.Unix(parseStartTime, 0)
	}

	if request.EndTime != "" {
		parseEndTime, err := strconv.ParseInt(request.EndTime, 10, 64)
		if err != nil {
			logger.Error(err.Error())
			return nil, grpcBadBodyError
		}
		endTime = time.Unix(parseEndTime, 0)
	}

	if request.Actor != "" && !uuid.IsValid(request.Actor) {
		actor, err := h.userService.GetByEmail(ctx, request.Actor)
		if err != nil {
			logger.Error(err.Error())
			switch {
			case errors.Is(err, user.ErrInvalidEmail), errors.Is(err, user.ErrNotExist):
				return nil, grpcBadBodyError
			default:
				return nil, grpcInternalServerError
			}
		}
		request.Actor = actor.ID
	}

	filter := activity.Filter{
		Actor:     request.Actor,
		Action:    request.Action,
		Data:      request.Data,
		Metadata:  request.Metadata,
		StartTime: startTime,
		EndTime:   endTime,
		Limit:     request.PageSize,
		Page:      request.PageNum,
	}

	activityResp, err := h.activityService.List(ctx, filter)
	if err != nil {
		logger.Error(err.Error())
		return nil, grpcInternalServerError
	}

	for _, activity := range activityResp.Activities {
		activityPB, err := transformActivityToPB(activity)
		if err != nil {
			logger.Error(err.Error())
			return nil, grpcInternalServerError
		}
		activities = append(activities, &activityPB)
	}

	return &shieldv1beta1.ListActivitiesResponse{
		Count:      int32(len(activities)),
		Activities: activities,
	}, nil
}

func transformActivityToPB(from audit.Log) (shieldv1beta1.Activity, error) {
	var dataMapString map[string]string
	if from.Data != nil {
		var dataMap map[string]interface{}
		if err := json.Unmarshal(from.Data.([]uint8), &dataMap); err != nil {
			return shieldv1beta1.Activity{}, err
		}
		var err error
		dataMapString, err = mapInterfaceToMapString(dataMap)
		if err != nil {
			return shieldv1beta1.Activity{}, err
		}
	}

	var metadataMapString map[string]string
	if from.Metadata != nil {
		var metadataMap map[string]interface{}
		if err := json.Unmarshal(from.Metadata.([]uint8), &metadataMap); err != nil {
			return shieldv1beta1.Activity{}, err
		}
		var err error
		metadataMapString, err = mapInterfaceToMapString(metadataMap)
		if err != nil {
			return shieldv1beta1.Activity{}, err
		}
	}

	return shieldv1beta1.Activity{
		Actor:     from.Actor,
		Action:    from.Action,
		Data:      dataMapString,
		Metadata:  metadataMapString,
		Timestamp: timestamppb.New(from.Timestamp),
	}, nil
}

func mapInterfaceToMapString(from map[string]interface{}) (map[string]string, error) {
	to := make(map[string]string)

	for k, v := range from {
		switch v.(type) {
		case string:
			to[k] = v.(string)
		case []interface{}:
			to[k] = fmt.Sprintf("%v", v)
		default:
			return map[string]string{}, ErrInternalServer
		}
	}

	return to, nil
}
