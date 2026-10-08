package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	agent_model "github.com/nougght/monitoring-system/server/internal/model/agent"
)

type EnrollmentKeysRepository struct {
	BaseRepository
}

func NewEnrollmentKeysRepository(db DB) *EnrollmentKeysRepository {
	return &EnrollmentKeysRepository{
		BaseRepository{db: db},
	}
}

func (r *EnrollmentKeysRepository) CreateKey(ctx context.Context, key *agent_model.EnrollmentKey) (*agent_model.EnrollmentKey, error) {
	query := `
	INSERT INTO enrollment_keys (key_hash, agent_id, expires_at, selector) 
	VALUES($1, $2, $3, $4)
	`
	_, err := r.conn(ctx).Exec(ctx, query, key.HashString, key.AgentID, key.ExpiresAt, key.Selector)
	if err != nil {
		return nil, fmt.Errorf("insert failed: %w", err)
	}

	return key, nil
}

func (r *EnrollmentKeysRepository) GetKeyBySelector(ctx context.Context, selector string) (*agent_model.EnrollmentKey, error) {
	query := `
	SELECT * FROM enrollment_keys k WHERE k.selector = $1
	`
	rows, err := r.conn(ctx).Query(ctx, query, selector)
	if err != nil {
		return nil, fmt.Errorf("select failed: %w", err)
	}
	defer rows.Close()

	key, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[agent_model.EnrollmentKey])
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("row collect failed: %w", err)
	}
	return &key, nil
}

func (r *EnrollmentKeysRepository) GetKeyByAgentId(ctx context.Context, agentID uuid.UUID) (*agent_model.EnrollmentKey, error) {
	query := `
	SELECT * FROM enrollment_keys k WHERE k.agent_id = $1
	`
	rows, err := r.conn(ctx).Query(ctx, query, agentID)
	if err != nil {
		return nil, fmt.Errorf("select failed: %w", err)
	}
	defer rows.Close()

	key, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[agent_model.EnrollmentKey])
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("row collect failed: %w", err)
	}
	return &key, nil
}

func (r *EnrollmentKeysRepository) SetUsed(ctx context.Context, agentID uuid.UUID, usedAt time.Time) error {
	query := `
		UPDATE enrollment_keys keys 
		SET used_at = $1
		WHERE agent_id = $2 AND used_at IS NULL
	`
	res, err := r.conn(ctx).Exec(ctx, query, usedAt, agentID)
	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	if res.RowsAffected() == 0 {
		return ErrNoAffectedRows
	}

	return nil
}
