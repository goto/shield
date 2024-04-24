package activity

import (
	"context"

	"github.com/goto/salt/audit"
	"github.com/mitchellh/mapstructure"
)

type Service struct {
	appConfig    AppConfig
	auditService *audit.Service
	repository   Repository
}

func NewService(appConfig AppConfig, repository Repository) *Service {
	return &Service{
		appConfig:    appConfig,
		repository:   repository,
		auditService: audit.New(audit.WithRepository(repository)),
	}
}

func (s Service) Log(ctx context.Context, action string, actor string, data any) error {
	if data == nil {
		return ErrInvalidData
	}

	var logDataMap map[string]interface{}
	if err := mapstructure.Decode(data, &logDataMap); err != nil {
		return err
	}

	metadata := map[string]any{
		"app_name":    "shield",
		"app_version": s.appConfig.Version,
	}

	ctx = audit.WithActor(ctx, actor)
	ctx, err := audit.WithMetadata(ctx, metadata)
	if err != nil {
		return err
	}

	return s.auditService.Log(ctx, action, logDataMap)
}

func (s Service) List(ctx context.Context, flt Filter) (PagedActivity, error) {
	pagedLogs, err := s.auditService.List(ctx, audit.Filter{
		Actor:     flt.Actor,
		Action:    flt.Action,
		Data:      flt.Data,
		Metadata:  flt.Metadata,
		StartTime: flt.StartTime,
		EndTime:   flt.EndTime,
		Limit:     flt.Limit,
		Page:      flt.Page,
	})
	if err != nil {
		return PagedActivity{}, nil
	}

	return PagedActivity{
		Count:      pagedLogs.Count,
		Activities: pagedLogs.Logs,
	}, nil
}
