package activity

import (
	"context"
	"time"

	"github.com/goto/salt/audit"
	"github.com/mitchellh/mapstructure"
)

type Service struct {
	appConfig  AppConfig
	repository Repository
}

func NewService(appConfig AppConfig, repository Repository) *Service {
	return &Service{
		appConfig:  appConfig,
		repository: repository,
	}
}

func (s Service) Log(ctx context.Context, action string, actor Actor, data any) error {
	if data == nil {
		return ErrInvalidData
	}

	var logDataMap map[string]interface{}
	if err := mapstructure.Decode(data, &logDataMap); err != nil {
		return err
	}

	metadata := map[string]string{
		"app_name":    "shield",
		"app_version": s.appConfig.Version,
		"email":       actor.Email,
	}

	log := &audit.Log{
		Timestamp: time.Now(),
		Action:    action,
		Data:      logDataMap,
		Actor:     actor.ID,
		Metadata:  metadata,
	}

	return s.repository.Insert(ctx, log)
}

func (s Service) List(ctx context.Context, filter Filter) (PagedActivity, error) {
	if !filter.EndTime.IsZero() && !filter.StartTime.IsZero() && filter.EndTime.Before(filter.StartTime) {
		return PagedActivity{}, ErrInvalidFilter
	}

	activities, err := s.repository.List(ctx, filter)
	if err != nil {
		return PagedActivity{}, err
	}

	return PagedActivity{
		Count:      int32(len(activities)),
		Activities: activities,
	}, nil
}
