package storage

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log"

	"github.com/golang-migrate/migrate/v4"
	migrate_pgx "github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

type DB interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	// begin real or pseudo transaction
	Begin(ctx context.Context) (pgx.Tx, error)
}

func RunMigrations(migrationFS fs.FS, pool *pgxpool.Pool) error {
	sourceDriver, err := iofs.New(migrationFS, "migrations")
	if err != nil {
		return fmt.Errorf("failed to create iofs driver: %w", err)
	}
	db := stdlib.OpenDBFromPool(pool)
	dbDriver, err := migrate_pgx.WithInstance(db, &migrate_pgx.Config{})
	if err != nil {
		return fmt.Errorf("failed to create pgx driver: %w", err)
	}
	m, err := migrate.NewWithInstance("iofs", sourceDriver, "postgres", dbDriver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("no changes to migrate")
			return nil
		}
		return fmt.Errorf("migrate up failed: %w", err)
	}

	return nil
}

func connect(ctx context.Context, config *PostgresConfig) (*pgxpool.Pool, error) {
	// TODO: add retrier
	pool, err := pgxpool.New(ctx, config.ConnString())
	if err != nil {
		return nil, fmt.Errorf("postgres connection failed: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}
	return pool, nil
}

func ConnectAndMigrate(ctx context.Context, config *PostgresConfig, migrationFS fs.FS) (*pgxpool.Pool, error) {
	pool, err := connect(ctx, config)
	if err != nil {
		return nil, err
	}

	if err := RunMigrations(migrationFS, pool); err != nil {
		return nil, fmt.Errorf("migrations failed: %w", err)
	}
	log.Println("successful migrations")
	return pool, nil
}
