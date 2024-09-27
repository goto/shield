package rule

import (
	"context"
	"regexp"
	"time"

	"github.com/goto/shield/core/rule/config"
)

type ConfigRepository interface {
	GetAll(ctx context.Context) ([]Ruleset, error)
	Upsert(ctx context.Context, name string, config Ruleset) (RuleConfig, error)
}

type Ruleset struct {
	Rules []Rule `yaml:"rules"`
}

type RuleConfig struct {
	ID        uint32
	Name      string
	Config    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Rule struct {
	Frontend    Frontend        `yaml:"frontend"`
	Backend     Backend         `yaml:"backend"`
	Middlewares MiddlewareSpecs `yaml:"middlewares"`
	Hooks       HookSpecs       `yaml:"hooks"`
}

type MiddlewareSpec struct {
	Name   string                 `yaml:"name"`
	Config map[string]interface{} `yaml:"config"`
}

type MiddlewareSpecs []MiddlewareSpec

func (m MiddlewareSpecs) Get(name string) (MiddlewareSpec, bool) {
	for _, n := range m {
		if n.Name == name {
			return n, true
		}
	}
	return MiddlewareSpec{}, false
}

type HookSpec struct {
	Name   string                 `yaml:"name"`
	Config map[string]interface{} `yaml:"config"`
}

type HookSpecs []HookSpec

func (m HookSpecs) Get(name string) (HookSpec, bool) {
	for _, n := range m {
		if n.Name == name {
			return n, true
		}
	}
	return HookSpec{}, false
}

type Frontend struct {
	URL   string         `yaml:"url"`
	URLRx *regexp.Regexp `yaml:"-"`

	Method string `yaml:"method"`
}

type Backend struct {
	URL       string `yaml:"url"`
	Namespace string `yaml:"namespace"`
	Prefix    string `yaml:"prefix"`
}

func YamlRulesetToRuleset(YamlRuleset config.Ruleset) Ruleset {
	targetRuleSet := Ruleset{}
	for _, theRule := range YamlRuleset.Rules {
		for _, backend := range theRule.Backends {
			for _, frontend := range backend.Frontends {
				middlewares := MiddlewareSpecs{}
				for _, middleware := range frontend.Middlewares {
					middlewares = append(middlewares, MiddlewareSpec{
						Name:   middleware.Name,
						Config: middleware.Config,
					})
				}

				hooks := HookSpecs{}
				for _, hook := range frontend.Hooks {
					hooks = append(hooks, HookSpec{
						Name:   hook.Name,
						Config: hook.Config,
					})
				}

				targetRuleSet.Rules = append(targetRuleSet.Rules, Rule{
					Frontend: Frontend{
						URL:    frontend.Path,
						Method: frontend.Method,
					},
					Backend:     Backend{URL: backend.Target, Namespace: backend.Name, Prefix: backend.Prefix},
					Middlewares: middlewares,
					Hooks:       hooks,
				})
			}
		}
	}
	return targetRuleSet
}
