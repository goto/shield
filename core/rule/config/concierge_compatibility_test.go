package config_test

import (
	"github.com/goto/shield/core/rule"
	"github.com/goto/shield/core/rule/config"
	"os"
	"path/filepath"
	"testing"
)

func TestConciergeConfigurationsRemainLoadable(t *testing.T) {
	root := os.Getenv("SHIELD_TEST_CONCIERGE_ROOT")
	if root == "" {
		t.Skip("optional local Concierge compatibility corpus")
	}
	files, basicAuth := 0, 0
	for _, env := range []string{"integration", "production"} {
		paths, err := filepath.Glob(filepath.Join(root, env, "rules", "proxy", "*"))
		if err != nil {
			t.Fatal(err)
		}
		for _, p := range paths {
			if filepath.Ext(p) != ".yaml" && filepath.Ext(p) != ".yml" {
				continue
			}
			data, err := os.ReadFile(p)
			if err != nil {
				t.Fatal(err)
			}
			set, err := config.ParseRulesetYaml(data)
			if err != nil {
				t.Fatalf("%s: %v", p, err)
			}
			normalized := rule.YamlRulesetToRuleset(set)
			if err := rule.ValidateRuleset(normalized); err != nil {
				t.Fatalf("%s: %v", p, err)
			}
			for _, r := range normalized.Rules {
				for _, m := range r.Middlewares {
					if m.Name == "basic_auth" {
						basicAuth++
					}
				}
			}
			files++
		}
	}
	if files == 0 {
		t.Fatal("no compatibility fixtures found")
	}
	t.Logf("Accepted %d rule files; Basic Auth routes: %d", files, basicAuth)
}
