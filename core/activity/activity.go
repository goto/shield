package activity

import (
	"context"

	"github.com/goto/salt/audit"
)

type Repository interface {
	Insert(ctx context.Context, log *audit.Log) error
}
