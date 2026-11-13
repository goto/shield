package rule_test

import (
	"context"
	"testing"

	"github.com/goto/shield/core/rule"
)

type retirementRepository struct {
	name string
	set  rule.Ruleset
}

func (r *retirementRepository) GetAll(context.Context) ([]rule.Ruleset, error) {
	return []rule.Ruleset{r.set}, nil
}
func (r *retirementRepository) Upsert(_ context.Context, name string, set rule.Ruleset) (rule.Config, error) {
	r.name = name
	r.set = set
	return rule.Config{Name: name}, nil
}
func TestNamedEmptyRulesetRetiresPreviouslyRegisteredRoute(t *testing.T) {
	repo := &retirementRepository{set: rule.Ruleset{Rules: []rule.Rule{{Frontend: rule.Frontend{URL: "/obsolete"}}}}}
	service := rule.NewService(repo)
	_, err := service.UpsertRulesConfigs(context.Background(), "shield_grpc", "rules: []\n")
	if err != nil {
		t.Fatal(err)
	}
	if repo.name != "shield_grpc" || len(repo.set.Rules) != 0 {
		t.Fatal("old named route not replaced")
	}
}
