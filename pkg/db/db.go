package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var (
	transactionContextKey = struct{}{}
)

type Client struct {
	db           *sqlx.DB
	queryTimeOut time.Duration
}

func New(cfg Config) (*Client, error) {
	d, err := sqlx.Open(cfg.Driver, cfg.URL)
	if err != nil {
		return nil, err
	}

	if err = d.Ping(); err != nil {
		return nil, err
	}

	d.SetMaxIdleConns(cfg.MaxIdleConns)
	d.SetMaxOpenConns(cfg.MaxOpenConns)
	d.SetConnMaxLifetime(cfg.ConnMaxLifeTime)

	return &Client{db: d, queryTimeOut: cfg.MaxQueryTimeoutInMS}, err
}

func (c Client) WithTimeout(ctx context.Context, op func(ctx context.Context) error) (err error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, c.queryTimeOut)
	defer cancel()

	return op(ctxWithTimeout)
}

// Handling transactions: https://stackoverflow.com/a/23502629/8244298
func (c Client) WithTxn(ctx context.Context, txnOptions sql.TxOptions, txFunc func(*sqlx.Tx) error) (err error) {
	txn, err := c.db.BeginTxx(ctx, &txnOptions)
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

func (c Client) WithTransaction(ctx context.Context, txnOptions sql.TxOptions) context.Context {
	if tx := extractTransaction(ctx); tx != nil {
		return ctx
	}
	tx, err := c.db.BeginTxx(ctx, &txnOptions)
	if err != nil {
		return ctx
	}
	return context.WithValue(ctx, transactionContextKey, tx)
}

func (c Client) Commit(ctx context.Context) error {
	if tx := extractTransaction(ctx); tx != nil {
		if err := tx.Commit(); err != nil {
			return err
		}
		return nil
	}
	return errors.New("no transaction")
}

func (c Client) Rollback(ctx context.Context) error {
	if tx := extractTransaction(ctx); tx != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return nil
	}
	return errors.New("no transaction")
}

func extractTransaction(ctx context.Context) *sqlx.Tx {
	if tx, ok := ctx.Value(transactionContextKey).(*sqlx.Tx); !ok {
		return nil
	} else {
		return tx
	}
}

func (c Client) GetDB(ctx context.Context) sqlx.ExtContext {
	if tx := extractTransaction(ctx); tx != nil {
		return tx
	}
	return c.db
}

func (c Client) QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	return c.GetDB(ctx).QueryRowxContext(ctx, query, args...)
}

func (c Client) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return c.GetDB(ctx).QueryContext(ctx, query, args...)
}

func (c Client) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return c.GetDB(ctx).ExecContext(ctx, query, args...)
}

func (c Client) Query(query string, args ...any) (*sql.Rows, error) {
	return c.db.Query(query, args...)
}

func (c Client) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return c.db.BeginTx(ctx, opts)
}

func (c Client) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return c.db.GetContext(ctx, dest, query, args...)
}

func (c Client) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return c.db.SelectContext(ctx, dest, query, args...)
}

func (c Client) Close() error {
	return c.db.Close()
}
