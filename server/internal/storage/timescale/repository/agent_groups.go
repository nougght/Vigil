package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/nougght/monitoring-system/server/internal/model/agent_groups"
)

type AgentGroupsRepository struct {
	BaseRepository
}

func NewAgentGroupsRepository(db DB) *AgentGroupsRepository {
	return &AgentGroupsRepository{
		BaseRepository{db: db},
	}
}

func (r *AgentGroupsRepository) CreateGroup(ctx context.Context, group *agent_groups.AgentGroupInfo) (*agent_groups.AgentGroup, error) {
	conn := r.conn(ctx)

	query := `
		INSERT INTO agent_groups(name, description)
		VALUES ($1, $2)
		RETURNING *;
	`

	rows, err := conn.Query(ctx, query, group.Name, group.Description)
	if IsConflict(err) {
		return nil, ErrConflict
	}
	if err != nil {
		return nil, fmt.Errorf("insert failed: %w", err)
	}
	defer rows.Close()

	res, err := CollectOnePtr[agent_groups.AgentGroup](rows)
	if err != nil {
		return nil, fmt.Errorf("collect row failed: %w", err)
	}
	return res, nil
}

func (r *AgentGroupsRepository) GetGroupByID(ctx context.Context, id uuid.UUID) (*agent_groups.AgentGroup, error) {
	conn := r.conn(ctx)

	query := `
		SELECT * 
		FROM agent_groups ag
		WHERE ag.id = $1
	`

	rows, err := conn.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("select failed: %w", err)
	}

	res, err := CollectOnePtr[agent_groups.AgentGroup](rows)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("collect failed: %w", err)
	}

	return res, nil
}

// all groups without deleted
func (r *AgentGroupsRepository) GetAllGroups(ctx context.Context) ([]*agent_groups.AgentGroup, error) {
	conn := r.conn(ctx)

	query := `
		SELECT * 
		FROM agent_groups ag
		WHERE deleted_at IS NULL;
	`

	rows, err := conn.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("select failed: %w", err)
	}

	res, err := CollectRowsPtr[agent_groups.AgentGroup](rows)
	if err != nil {
		return nil, fmt.Errorf("collect failed: %w", err)
	}

	return res, nil
}

// update group
func (r *AgentGroupsRepository) UpdateGroup(ctx context.Context, group *agent_groups.UpdateAgentGroupInput) error {
	conn := r.conn(ctx)

	query := `
		UPDATE agent_groups ag
		SET name = COALESCE($2, ag.name),
			description = CASE WHEN $3::bool THEN $4::text ELSE ag.description END
		WHERE ag.id = $1
			AND deleted_at IS NULL;
	`

	set, desc := group.Description.GetPtrIfSet()
	res, err := conn.Exec(ctx, query, group.ID, group.Name, set, desc)
	if IsConflict(err) {
		return ErrConflict
	}
	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	if res.RowsAffected() == 0 {
		return ErrNoAffectedRows
	}

	return nil
}

// mark group as deleted
func (r *AgentGroupsRepository) DeleteGroupByID(ctx context.Context, id uuid.UUID) error {
	conn := r.conn(ctx)

	query := `
		UPDATE agent_groups ag
		SET deleted_at = now()
		WHERE id = $1 
			AND deleted_at IS NULL
	`

	res, err := conn.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete failed: %w", err)
	}

	if res.RowsAffected() == 0 {
		return ErrNoAffectedRows
	}

	return nil
}
