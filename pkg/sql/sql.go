package sql

import (
	"context"
	"database/sql"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Config struct {
	Driver              string
	URL                 string
	MaxIdleConns        int
	MaxOpenConns        int
	ConnMaxLifeTime     time.Duration
	MaxQueryTimeoutInMS time.Duration
}

type SQL struct {
	*sqlx.DB
	queryTimeOut time.Duration
}

func New(config Config) (*SQL, error) {
	d, err := sqlx.Open(config.Driver, config.URL)

	if err != nil {
		return nil, err
	}

	if err = d.Ping(); err != nil {
		return nil, err
	}

	d.SetMaxIdleConns(config.MaxIdleConns)
	d.SetMaxOpenConns(config.MaxOpenConns)
	d.SetConnMaxLifetime(config.ConnMaxLifeTime)
	fmt.Println("sql timeout config: ", config.MaxQueryTimeoutInMS)
	return &SQL{DB: d, queryTimeOut: config.MaxQueryTimeoutInMS}, err
}

func (s SQL) WithTimeout(ctx context.Context, op func(ctx context.Context) error) (err error) {
	if ctx.Err() != nil {
		debug.PrintStack()
		fmt.Println("⚠️ Parent context is already canceled! Reason:", ctx.Err())
	}
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), s.queryTimeOut)
	defer cancel()
	// Logging: show deadlines if available
	if parentDeadline, ok := ctx.Deadline(); ok {
		fmt.Printf("🧭 Parent context deadline: %v (in %d ms)\n", parentDeadline, time.Until(parentDeadline).Milliseconds())
	} else {
		fmt.Println("🧭 Parent context has no deadline")
	}

	if childDeadline, ok := ctxWithTimeout.Deadline(); ok {
		fmt.Printf("⏳ Child context deadline:  %v (in %d ms)\n", childDeadline, time.Until(childDeadline).Milliseconds())
	}
	err = op(ctxWithTimeout)

	// After the operation, check context errors
	select {
	case <-ctxWithTimeout.Done():
		fmt.Println("⚠️ CHILD context done. Error:", ctxWithTimeout.Err())
	case <-ctx.Done():
		fmt.Println("⚠️ PARENT context done. Error:", ctx.Err())
	default:
		fmt.Println("✅ No context cancellation detected")
	}

	return err
}

// Handling transactions: https://stackoverflow.com/a/23502629/8244298
func (s SQL) WithTxn(ctx context.Context, txnOptions sql.TxOptions, txFunc func(*sqlx.Tx) error) (err error) {
	txn, err := s.BeginTxx(ctx, &txnOptions)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			switch p := p.(type) {
			case error:
				err = p
			default:
				err = errors.Errorf("%s", p)
			}
			err = txn.Rollback()
			panic(p)
		} else if err != nil {
			if rlbErr := txn.Rollback(); err != nil {
				err = fmt.Errorf("rollback error: %s while executing: %w", rlbErr, err)
			} else {
				err = fmt.Errorf("rollback: %w", err)
			}
			err = fmt.Errorf("rollback: %w", err)
		} else {
			err = txn.Commit()
		}
	}()

	err = txFunc(txn)
	return err
}
