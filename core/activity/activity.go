package activity

import (
	"context"

	"github.com/goto/salt/audit"
)

const (
	SinkTypeDB     = "db"
	SinkTypeStdout = "stdout"
	SinkTypeNone   = "none"
)

type Repository interface {
	Insert(ctx context.Context, log *audit.Log) error
}
