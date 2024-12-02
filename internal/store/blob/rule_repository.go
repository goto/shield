package blob

import (
	"context"
	"io"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/goto/salt/log"

	"github.com/robfig/cron/v3"

	"github.com/goto/shield/core/rule"
	"github.com/goto/shield/core/rule/config"
	"github.com/pkg/errors"
	"gocloud.dev/blob"
)

type RuleRepository struct {
	log log.Logger
	mu  *sync.Mutex

	cron      *cron.Cron
	bucket    Bucket
	cached    []rule.Ruleset
	updatedAt time.Time
}

func (repo *RuleRepository) GetAll(ctx context.Context) ([]rule.Ruleset, error) {
	repo.mu.Lock()
	currentCache := repo.cached
	repo.mu.Unlock()
	if repo.cron != nil {
		// cache must have been refreshed automatically, just return
		return currentCache, nil
	}

	err := repo.refresh(ctx)
	return repo.cached, err
}

func (repo *RuleRepository) Fetch(ctx context.Context) ([]rule.Ruleset, error) {
	return repo.GetAll(ctx)
}

func (repo *RuleRepository) refresh(ctx context.Context) error {
	var rulesets []rule.Ruleset

	// get all items
	it := repo.bucket.List(&blob.ListOptions{})
	for {
		obj, err := it.Next(ctx)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if obj.IsDir {
			continue
		}
		if !(strings.HasSuffix(obj.Key, ".yaml") || strings.HasSuffix(obj.Key, ".yml")) {
			continue
		}
		fileBytes, err := repo.bucket.ReadAll(ctx, obj.Key)
		if err != nil {
			return errors.Wrap(err, "bucket.ReadAll: "+obj.Key)
		}

		s, err := config.ParseRulesetYaml(fileBytes)
		if err != nil {
			return errors.Wrap(err, "yaml.Unmarshal: "+obj.Key)
		}
		if len(s.Rules) == 0 {
			continue
		}

		targetRuleSet := rule.YamlRulesetToRuleset(s)

		// parse all urls at this time only to avoid doing it usage
		rxParsingSuccess := true
		for ruleIdx, rule := range targetRuleSet.Rules {
			// TODO: only compile between delimiter, maybe angular brackets
			targetRuleSet.Rules[ruleIdx].Frontend.URLRx, err = regexp.Compile(rule.Frontend.URL)
			if err != nil {
				rxParsingSuccess = false
				repo.log.Error("failed to parse rule frontend as a valid regular expression",
					"url", rule.Frontend.URL, "err", err)
			}
		}

		if rxParsingSuccess {
			rulesets = append(rulesets, targetRuleSet)
		} else {
			repo.log.Warn("skipping rule set due to parsing errors", "content", string(fileBytes))
		}
	}

	repo.mu.Lock()
	repo.cached = rulesets
	repo.updatedAt = time.Now()
	repo.mu.Unlock()
	repo.log.Debug("rule cache refreshed", "ruleset_count", len(repo.cached))
	return nil
}

func (repo *RuleRepository) InitCache(ctx context.Context, refreshDelay time.Duration) error {
	repo.cron = cron.New(cron.WithChain(
		cron.SkipIfStillRunning(cron.DefaultLogger),
	))
	if _, err := repo.cron.AddFunc("@every "+refreshDelay.String(), func() {
		if err := repo.refresh(ctx); err != nil {
			repo.log.Warn("failed to refresh rule repository", "err", err)
		}
	}); err != nil {
		return err
	}
	repo.cron.Start()

	// do it once right now
	return repo.refresh(ctx)
}

func (repo *RuleRepository) Close() error {
	<-repo.cron.Stop().Done()
	return repo.bucket.Close()
}

func (repo *RuleRepository) Upsert(ctx context.Context, name string, config rule.Ruleset) (rule.Config, error) {
	// upsert is currently not supported for BLOB rule config storage type
	return rule.Config{}, rule.ErrUpsertConfigNotSupported
}

func (repo *RuleRepository) IsUpdated(ctx context.Context, lastUpdated time.Time) bool {
	return repo.updatedAt.After(lastUpdated)
}

func NewRuleRepository(logger log.Logger, b Bucket) *RuleRepository {
	return &RuleRepository{
		log:    logger,
		bucket: b,
		mu:     new(sync.Mutex),
	}
}
