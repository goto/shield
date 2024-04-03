package activity

import (
	"context"
	"time"

	"github.com/goto/salt/audit"
	"github.com/goto/shield/config"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s Service) Log(ctx context.Context, action string, actor string, data map[string]string) error {
	metadata := map[string]string{
		"app_name":    "shield",
		"app_version": config.Version,
	}

	log := &audit.Log{
		Timestamp: time.Now(),
		Action:    action,
		Data:      data,
		Actor:     actor,
		Metadata:  metadata,
	}

	return s.repository.Insert(ctx, log)
}
