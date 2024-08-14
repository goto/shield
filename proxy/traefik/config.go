package traefik

import (
	_ "embed"
)

var (
	//go:embed dynamic-config.yaml
	RuleYaml string
)
