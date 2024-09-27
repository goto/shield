package postgres

import (
	"time"

	"github.com/goto/shield/core/rule"
)

type RuleConfig struct {
	ID        uint32    `db:"id"`
	Name      string    `db:"name"`
	Config    string    `db:"config"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (from RuleConfig) transformToRuleConfig() rule.RuleConfig {
	return rule.RuleConfig{
		ID:        from.ID,
		Name:      from.Name,
		Config:    from.Config,
		CreatedAt: from.CreatedAt,
		UpdatedAt: from.UpdatedAt,
	}
}
