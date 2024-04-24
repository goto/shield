package activity

import (
	"context"
	"time"

	"github.com/goto/salt/audit"
)

const (
	SinkTypeDB     = "db"
	SinkTypeStdout = "stdout"
	SinkTypeNone   = "none"
)

type Repository interface {
	Insert(ctx context.Context, log *audit.Log) error
	List(ctx context.Context, filter Filter) ([]audit.Log, error)
}

type AppConfig struct {
	Version string
}

type Filter struct {
	Actor     string
	Action    string
	Data      map[string]string
	Metadata  map[string]string
	StartTime time.Time
	EndTime   time.Time
	Limit     int32
	Page      int32
}

type PagedActivity struct {
	Count      int32
	Activities []audit.Log
}

type Actor struct {
	ID    string
	Email string
}
