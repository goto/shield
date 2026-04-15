package rule_test

import (
	"testing"

	"github.com/goto/shield/core/rule"
	"github.com/goto/shield/core/rule/config"
	"github.com/stretchr/testify/assert"
)

func TestYamlRulesetToRuleset_SkipReadBody(t *testing.T) {
	t.Run("should propagate SkipReadBody=true from config to rule", func(t *testing.T) {
		yamlRuleset := config.Ruleset{
			Rules: []config.Rule{
				{
					Backends: []config.Backend{
						{
							Name:   "dex",
							Target: "http://dex:80",
							Prefix: "/dex",
							Frontends: []config.Frontend{
								{
									Path:         "/dex/data_upload/upload",
									Method:       "POST",
									SkipReadBody: true,
								},
							},
						},
					},
				},
			},
		}

		result := rule.YamlRulesetToRuleset(yamlRuleset)

		assert.Len(t, result.Rules, 1)
		assert.True(t, result.Rules[0].SkipReadBody, "SkipReadBody should be true")
		assert.Equal(t, "/dex/data_upload/upload", result.Rules[0].Frontend.URL)
		assert.Equal(t, "POST", result.Rules[0].Frontend.Method)
	})

	t.Run("should default SkipReadBody to false", func(t *testing.T) {
		yamlRuleset := config.Ruleset{
			Rules: []config.Rule{
				{
					Backends: []config.Backend{
						{
							Name:   "dex",
							Target: "http://dex:80",
							Frontends: []config.Frontend{
								{
									Path:   "/dex/api/resource",
									Method: "GET",
								},
							},
						},
					},
				},
			},
		}

		result := rule.YamlRulesetToRuleset(yamlRuleset)

		assert.Len(t, result.Rules, 1)
		assert.False(t, result.Rules[0].SkipReadBody, "SkipReadBody should default to false")
	})
}
