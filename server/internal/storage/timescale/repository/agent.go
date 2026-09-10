package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/nougght/monitoring-system/server/internal/model"
	agent_model "github.com/nougght/monitoring-system/server/internal/model/agent"
)

type AgentRepository struct {
	pool DB
}

func NewAgentRepository(db DB) *AgentRepository {
	return &AgentRepository{
		pool: db,
	}

}

func (r *AgentRepository) db(ctx context.Context) DB {
	res := r.pool
	if tx := ctx.Value(model.ContextKeyTx); tx != nil {
		res = tx.(DB)
	}
	return res
}

func (r *AgentRepository) CreateAgent(ctx context.Context, agent *agent_model.Agent) (*agent_model.Agent, error) {
	query := `
	INSERT INTO agents (name, description) 
	VALUES($1, $2)
	RETURNING id, created_at, last_seen_at
	`
	err := r.db(ctx).QueryRow(ctx, query, agent.Name, agent.Description).Scan(&agent.ID, &agent.CreatedAt, &agent.LastSeenAt)
	if err != nil {
		return nil, fmt.Errorf("insert failed: %w", err)
	}

	return agent, nil
}

func (r *AgentRepository) GetAllAgents(ctx context.Context) (res []*agent_model.Agent, err error) {
	query := `
		SELECT * FROM agents;
		`
	// SELECT *,  EXISTS (
	//     SELECT 1
	//     FROM enrollment_keys keys
	//     WHERE keys.agent_id = agent.id AND keys.used_at != NULL
	// ) AS is_enrolled FROM agents
	// `
	rows, err := r.db(ctx).Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("select failed: %w", err)
	}
	defer rows.Close()

	res, err = pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[agent_model.Agent])
	if err != nil {
		return nil, fmt.Errorf("collect rows failed: %w", err)
	}

	return res, nil
}

func (r *AgentRepository) GetAgentByID(ctx context.Context, id uuid.UUID) (res *agent_model.Agent, err error) {
	query := `
		SELECT * FROM agents WHERE agents.id = $1;
		`
	rows, err := r.db(ctx).Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("select failed: %w", err)
	}
	defer rows.Close()

	res, err = pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[agent_model.Agent])
	if errors.Is(err, pgx.ErrNoRows) {
		err = ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("collect row failed: %w", err)
	}

	return res, nil
}

func (r *AgentRepository) GetTotalAgents(ctx context.Context) (int, error) {
	query := `
	SELECT COUNT(*) FROM agents;
	`
	var total int
	err := r.db(ctx).QueryRow(ctx, query).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("select failed: %w", err)
	}
	return total, nil
}

func (r *AgentRepository) GetAgentNamesByIDs(ctx context.Context, agentIDs []uuid.UUID) (map[uuid.UUID]string, error) {
	query := `
		SELECT id, name FROM agents WHERE id = ANY($1);
	`
	rows, err := r.db(ctx).Query(ctx, query, agentIDs)
	if err != nil {
		return nil, fmt.Errorf("select failed: %w", err)
	}
	defer rows.Close()
	names := make(map[uuid.UUID]string)
	for rows.Next() {
		var id uuid.UUID
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		names[id] = name
	}
	return names, nil
}

func (r *AgentRepository) UpdateStatus(ctx context.Context, agentID uuid.UUID, status agent_model.AgentStatus) error {
	query := `
	UPDATE agents SET status = $1 WHERE ID = $2
	`
	_, err := r.db(ctx).Exec(ctx, query, status, agentID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("insert failed: %w", ErrNotFound)
		}
		return fmt.Errorf("insert failed: %w", err)
	}

	return nil
}
