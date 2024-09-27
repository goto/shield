package v1beta1

import (
	"context"
	"errors"

	"github.com/goto/shield/core/rule"
	shieldv1beta1 "github.com/goto/shield/proto/v1beta1"
	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type RuleService interface {
	GetAllConfigs(ctx context.Context) ([]rule.Ruleset, error)
	UpsertRulesConfigs(ctx context.Context, name string, config string) (rule.RuleConfig, error)
}

func (h Handler) UpsertRulesConfig(ctx context.Context, request *shieldv1beta1.UpsertRulesConfigRequest) (*shieldv1beta1.UpsertRulesConfigResponse, error) {
	logger := grpczap.Extract(ctx)

	rc, err := h.ruleService.UpsertRulesConfigs(ctx, request.Name, request.Config)
	if err != nil {
		logger.Error(err.Error())
		switch {
		case errors.Is(err, rule.ErrUpsertConfigNotSupported):
			return nil, grpcUnsupportedError
		case errors.Is(err, rule.ErrInvalidRuleConfig):
			return nil, grpcBadBodyError
		default:
			return nil, grpcInternalServerError
		}
	}

	rc.Config = request.Config
	return ruleConfigToPB(rc), nil
}

func ruleConfigToPB(from rule.RuleConfig) *shieldv1beta1.UpsertRulesConfigResponse {
	return &shieldv1beta1.UpsertRulesConfigResponse{
		Id:        from.ID,
		Name:      from.Name,
		Config:    from.Config,
		CreatedAt: timestamppb.New(from.CreatedAt),
		UpdatedAt: timestamppb.New(from.UpdatedAt),
	}
}
