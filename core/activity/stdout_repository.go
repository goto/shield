package activity

import (
	"context"
	"encoding/json"
	"io"

	"github.com/goto/salt/audit"
)

type StdoutRepository struct {
	writer io.Writer
}

func NewStdoutRepository(writer io.Writer) *StdoutRepository {
	return &StdoutRepository{
		writer: writer,
	}
}

func (r StdoutRepository) Insert(ctx context.Context, log *audit.Log) error {
	return json.NewEncoder(r.writer).Encode(log)
}

// TODO : Fix this
func (r StdoutRepository) List(ctx context.Context, filter Filter) ([]audit.Log, error) {
	return []audit.Log{}, nil
}
