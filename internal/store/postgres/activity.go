package postgres

import "time"

type Activity struct {
	Actor     string            `db:"actor"`
	Action    string            `db:"action"`
	Data      map[string]string `db:"data"`
	Metadata  map[string]string `db:"metadata"`
	Timestamp time.Time         `db:"timestamp"`
}
