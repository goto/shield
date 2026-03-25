package cmd

import (
	"context"
	"errors"
	"net/url"

	"github.com/goto/salt/log"
	"github.com/goto/shield/core/rule"
	"github.com/goto/shield/internal/proxy"
	"github.com/goto/shield/internal/proxy/envoy/xds"
	"github.com/goto/shield/internal/proxy/envoy/xds/ads"
	"github.com/goto/shield/internal/store/blob"
	"github.com/goto/shield/internal/store/postgres"
)

func serveXDS(ctx context.Context, logger *log.Zap, cfg proxy.ServicesConfig, pgRuleRepository *postgres.RuleRepository) ([]func() error, error) {
	cleanUpBlobs, repositories, err := buildXDSDependencies(ctx, logger, cfg, pgRuleRepository)
	if err != nil {
		return nil, err
	}

	errChan := make(chan error)
	go func() {
		err := xds.Serve(ctx, logger, cfg, repositories)
		if err != nil {
			errChan <- err
			logger.Error("error while running envoy xds server", "error", err)
		}
	}()

	return cleanUpBlobs, nil
}

func buildXDSDependencies(ctx context.Context, logger *log.Zap, cfg proxy.ServicesConfig, pgRuleRepository *postgres.RuleRepository) ([]func() error, map[string]ads.Repository, error) {
	var cleanUpBlobs []func() error
	repositories := make(map[string]ads.Repository)

	for _, svcConfig := range cfg.Services {
		parsedRuleConfigURL, err := url.Parse(svcConfig.RulesPath)
		if err != nil {
			return nil, nil, err
		}

		var repository ads.Repository
		switch parsedRuleConfigURL.Scheme {
		case rule.RULES_CONFIG_STORAGE_PG:
			repository = pgRuleRepository
		case rule.RULES_CONFIG_STORAGE_GS,
			rule.RULES_CONFIG_STORAGE_FILE,
			rule.RULES_CONFIG_STORAGE_MEM:
			ruleBlobFS, err := blob.NewStore(ctx, svcConfig.RulesPath, svcConfig.RulesPathSecret)
			if err != nil {
				return nil, nil, err
			}

			blobRuleRepository := blob.NewRuleRepository(logger, ruleBlobFS)
			if err := blobRuleRepository.InitCache(ctx, ruleCacheRefreshDelay); err != nil {
				return nil, nil, err
			}
			cleanUpBlobs = append(cleanUpBlobs, blobRuleRepository.Close)
			repository = blobRuleRepository
		default:
			return nil, nil, errors.New("invalid rule config storage")
		}
		repositories[svcConfig.Name] = repository
	}

	return cleanUpBlobs, repositories, nil
}
