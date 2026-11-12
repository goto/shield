package rule

import rulecfg "github.com/goto/shield/core/rule/config"

// ValidateRuleset is shared by persisted JSON and YAML-backed configurations.
func ValidateRuleset(set Ruleset) error {
	for _, r := range set.Rules {
		for _, m := range r.Middlewares {
			if err := rulecfg.ValidateMiddlewareConfig(m.Name, m.Config); err != nil {
				return err
			}
		}
	}
	return nil
}
