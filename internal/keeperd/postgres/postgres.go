package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/derpartizanen/gophkeeper/internal/logger"
)

const (
	ConnAttempts = 10
	ConnTimeout  = time.Second
)

var DefaultTxOptions = pgx.TxOptions{
	IsoLevel: pgx.ReadCommitted,
}

type PgxIface interface {
	Acquire(ctx context.Context) (*pgxpool.Conn, error)

	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row

	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	Close()
}

// Postgres represents abstraction over database connection pool.
type Postgres struct {
	Pool PgxIface
}

// New initializes connection to database and creates new wrapper object.
func New(url string, log *logger.Logger) (*Postgres, error) {
	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("postgres - New - pgxpool.ParseConfig: %w", err)
	}

	pg := new(Postgres)
	attempts := ConnAttempts

	for attempts > 0 {
		pg.Pool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err == nil {
			break
		}

		log.Info().Msgf("Postgres is trying to connect, attempts left: %d", attempts)
		time.Sleep(ConnTimeout)

		attempts--
	}

	if err != nil {
		return nil, fmt.Errorf("postgres - New - attempts == 0: %w", err)
	}

	return pg, nil
}

// Close gracefully closes connection to database.
func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}

// BeginTx starts new database transaction.
func (p *Postgres) BeginTx(ctx context.Context) (pgx.Tx, error) {
	tx, err := p.Pool.BeginTx(ctx, DefaultTxOptions)
	if err != nil {
		return nil, fmt.Errorf("postgres - BeginTx - conn.BeginTx: %w", err)
	}

	return tx, nil
}

// Select returns rows matching query.
func (p *Postgres) Select(
	ctx context.Context,
	dst interface{},
	query string,
	args ...interface{},
) error {
	return pgxscan.Select(ctx, p.Pool, dst, query, args...)
}

// RunAtomic executes provided function in Postgres transaction.
func (p *Postgres) RunAtomic(ctx context.Context, operation func(tx Transaction) error) error {
	tx, err := p.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("postgres - RunAtomic - p.BeginTx: %w", err)
	}

	if err := operation(tx); err != nil {
		if rErr := tx.Rollback(ctx); rErr != nil {
			logger.FromContext(ctx).Error().Err(rErr).Msg("postgres - RunAtomic - tx.Rollback")
		}

		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("postgres - RunAtomic - tx.Commit: %w", err)
	}

	return nil
}
