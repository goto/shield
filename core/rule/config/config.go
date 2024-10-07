package config

import (
	"gopkg.in/yaml.v2"
)

type Ruleset struct {
	Rules []Rule `yaml:"rules"`
}

type Rule struct {
	Backends []Backend `yaml:"backends"`
}

type Backend struct {
	Name      string     `yaml:"name"`
	Target    string     `yaml:"target"`
	Methods   []string   `yaml:"methods"`
	Frontends []Frontend `yaml:"frontends"`
	Prefix    string     `yaml:"prefix"`
}

type Frontend struct {
	Action      string       `yaml:"action"`
	Path        string       `yaml:"path"`
	Method      string       `yaml:"method"`
	Middlewares []Middleware `yaml:"middlewares"`
	Hooks       []Hook       `yaml:"hooks"`
}

type Middleware struct {
	Name   string                 `yaml:"name"`
	Config map[string]interface{} `yaml:"config"`
}

type Hook struct {
	Name   string                 `yaml:"name"`
	Config map[string]interface{} `yaml:"config"`
}

func ParseRulesetYaml(fileBytes []byte) (Ruleset, error) {
	var s Ruleset
	if err := yaml.Unmarshal(fileBytes, &s); err != nil {
		return Ruleset{}, err
	}
	return s, nil
}
