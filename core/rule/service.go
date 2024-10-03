package rule

import (
	"context"
	"strings"

	rulecfg "github.com/goto/shield/core/rule/config"
)

const (
	RULES_CONFIG_STORAGE_PG   = "postgres"
	RULES_CONFIG_STORAGE_GS   = "gs"
	RULES_CONFIG_STORAGE_FILE = "file"
	RULES_CONFIG_STORAGE_MEM  = "mem"
)

type Service struct {
	configRepository ConfigRepository
}

func NewService(configRepository ConfigRepository) *Service {
	return &Service{
		configRepository: configRepository,
	}
}

func (s Service) GetAllConfigs(ctx context.Context) ([]Ruleset, error) {
	return s.configRepository.GetAll(ctx)
}

func (s Service) UpsertRulesConfigs(ctx context.Context, name string, config string) (Config, error) {
	if strings.TrimSpace(name) == "" {
		return Config{}, ErrInvalidRuleConfig
	}

	if strings.TrimSpace(config) == "" {
		return Config{}, ErrInvalidRuleConfig
	}

	yamlRuleset, err := rulecfg.ParseRulesetYaml([]byte(config))
	if err != nil {
		return Config{}, ErrInvalidRuleConfig
	}

	targetRuleset := YamlRulesetToRuleset(yamlRuleset)
	return s.configRepository.Upsert(ctx, name, targetRuleset)
}
