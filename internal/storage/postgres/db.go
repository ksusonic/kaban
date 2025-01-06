package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
)

type txCtx struct{}
type connCtx struct{}

type DB struct {
	pool *pgxpool.Pool
}

func NewDB(ctx context.Context, log *slog.Logger) (*DB, func(), error) {
	// Using environment variables instead of a connection string.
	conf, err := pgxpool.ParseConfig("")
	if err != nil {
		return nil, nil, err
	}

	conf.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   &pgxLogWrapper{Logger: log},
		LogLevel: logLevelFromEnv(),
	}

	pool, err := pgxpool.NewWithConfig(ctx, conf)
	if err != nil {
		return nil, nil, err
	}

	log.DebugContext(ctx, "pinging postgres", "host", conf.ConnConfig.Host, "port", conf.ConnConfig.Port)
	if err = pool.Ping(ctx); err != nil {
		return nil, nil, fmt.Errorf("ping db: %w", err)
	}
	log.DebugContext(ctx, "pg ping OK")

	return &DB{pool}, pool.Close, err
}

// TransactionContext returns a copy of the parent context which begins a transaction
// to PostgreSQL.
//
// Once the transaction is over, you must call db.Commit(ctx) to make the changes effective.
// This might live in the go-pkg/postgres package later for the sake of code reuse.
func (db DB) TransactionContext(ctx context.Context) (context.Context, error) {
	tx, err := db.Conn(ctx).Begin(ctx)
	if err != nil {
		return nil, err
	}
	return context.WithValue(ctx, txCtx{}, tx), nil
}

// Commit transaction from context.
func (db DB) Commit(ctx context.Context) error {
	if tx, ok := ctx.Value(txCtx{}).(pgx.Tx); ok && tx != nil {
		return tx.Commit(ctx)
	}
	return errors.New("context has no transaction")
}

// Rollback transaction from context.
func (db DB) Rollback(ctx context.Context) error {
	if tx, ok := ctx.Value(txCtx{}).(pgx.Tx); ok && tx != nil {
		return tx.Rollback(ctx)
	}
	return errors.New("context has no transaction")
}

// WithAcquire defer postgres.Release(dbCtx).
func (db DB) WithAcquire(ctx context.Context) (context.Context, error) {
	if _, ok := ctx.Value(connCtx{}).(*pgxpool.Conn); ok {
		panic("context already has a connection acquired")
	}
	res, err := db.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	return context.WithValue(ctx, connCtx{}, res), nil
}

// Release PostgreSQL connection acquired by context back to the pool.
func (db DB) Release(ctx context.Context) {
	if res, ok := ctx.Value(connCtx{}).(*pgxpool.Conn); ok && res != nil {
		res.Release()
	}
}

// Conn returns a PostgreSQL transaction if one exists.
// If not, returns a connection if a connection has been acquired by calling WithAcquire.
// Otherwise, it returns *pgxpool.Pool which acquires the connection and closes after a SQL command is executed.
func (db DB) Conn(ctx context.Context) PGXQuerier {
	if tx, ok := ctx.Value(txCtx{}).(pgx.Tx); ok && tx != nil {
		return tx
	}
	if res, ok := ctx.Value(connCtx{}).(*pgxpool.Conn); ok && res != nil {
		return res
	}
	return db.pool
}
