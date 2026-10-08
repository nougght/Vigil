package repository

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/nougght/monitoring-system/server/internal/model"
)

type Transactor struct {
	db DB
}

func NewTransactor(db DB) *Transactor {
	return &Transactor{
		db: db,
	}
}

// do func in transaction with nesting check
func (t *Transactor) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	// if nested - use provided tx
	if tx := getTx(ctx); tx != nil {
		return fn(ctx)
	}
	// otherwise - begin new tx

	tx, err := t.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		err := tx.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			log.Printf("transaction rollback failed: %s", err.Error())
		}
	}()

	ctx = context.WithValue(ctx, model.ContextKeyTx, tx)
	if err = fn(ctx); err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commit transaction failed: %w", err)
	}
	return nil
}

// returns tx from ctx or nil
func getTx(ctx context.Context) (tx *pgx.Tx) {
	tx, ok := ctx.Value(model.ContextKeyTx).(*pgx.Tx)
	if !ok {
		return
	}
	return tx
}

// try get tx from ctx or use provided db
func conn(ctx context.Context, db DB) DB {
	if tx := getTx(ctx); tx != nil {
		return *tx
	}
	return db
}
