package postgres

import (
	"time"

	"github.com/goto/salt/audit"
	"github.com/jmoiron/sqlx/types"
)

type AuditActivity struct {
	Timestamp time.Time          `db:"timestamp"`
	Action    string             `db:"action"`
	Actor     string             `db:"actor"`
	Data      types.NullJSONText `db:"data"`
	Metadata  types.NullJSONText `db:"metadata"`
}

func (from AuditActivity) transformToAuditLog() audit.Log {
	return audit.Log{
		Timestamp: from.Timestamp,
		Action:    from.Action,
		Actor:     from.Actor,
		Data:      from.Data,
		Metadata:  from.Metadata,
	}
}
