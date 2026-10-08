package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	agent_model "github.com/nougght/monitoring-system/server/internal/model/agent"
)

type AgentRepository struct {
	BaseRepository
}

var (
	ErrInvalidGroup = errors.New("invalid group")
)

func NewAgentRepository(db DB) *AgentRepository {
	return &AgentRepository{
		BaseRepository{db: db},
	}

}

func (r *AgentRepository) CreateAgent(ctx context.Context, agent *agent_model.Agent) (*agent_model.Agent, error) {
	query := `
	INSERT INTO agents (name, description) 
	VALUES($1, $2)
	RETURNING id, created_at, last_seen_at
	`
	err := r.conn(ctx).QueryRow(ctx, query, agent.Name, agent.Description).Scan(&agent.ID, &agent.CreatedAt, &agent.LastSeenAt)
	if err != nil {
		return nil, fmt.Errorf("insert failed: %w", err)
	}
	return agent, nil
}

func (r *AgentRepository) GetAllAgents(ctx context.Context) (res []*agent_model.Agent, err error) {
	query := `
		SELECT a.id,
			   a.name,
			   a.description,
			    CASE 
					WHEN ag.deleted_at != NULL THEN 
						a.group_id
					ELSE 
						NULL
				END AS group_id, 
			   a.created_at,
			   a.status,
			   a.last_seen_at
		FROM agents a
		LEFT JOIN agent_groups ag
			ON a.group_id = ag.id;
		`
	// SELECT *,  EXISTS (
	//     SELECT 1
	//     FROM enrollment_keys keys
	//     WHERE keys.agent_id = agent.id AND keys.used_at != NULL
	// ) AS is_enrolled FROM agents
	// `
	rows, err := r.conn(ctx).Query(ctx, query)
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
		SELECT id,
			   name,
			   description,
			    CASE 
					WHEN ag.deleted_at != NULL THEN 
						a.group_id
					ELSE 
						NULL
				END AS group_id, 
			   created_at,
			   last_seen_at
		FROM agents a
		LEFT JOIN agent_groups ag
			ON a.group_id = ag.id
		WHERE a.id = $1
		LIMIT 1;
		`
	rows, err := r.conn(ctx).Query(ctx, query, id)
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
	err := r.conn(ctx).QueryRow(ctx, query).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("select failed: %w", err)
	}
	return total, nil
}

func (r *AgentRepository) GetAgentNamesByIDs(ctx context.Context, agentIDs []uuid.UUID) (map[uuid.UUID]string, error) {
	query := `
		SELECT id, name FROM agents WHERE id = ANY($1);
	`
	rows, err := r.conn(ctx).Query(ctx, query, agentIDs)
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
	res, err := r.conn(ctx).Exec(ctx, query, status, agentID)
	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	if res.RowsAffected() == 0 {
		return ErrNoAffectedRows
	}

	return nil
}

func (r *AgentRepository) UpdateAgentsGroup(ctx context.Context, agentIDs []uuid.UUID, groupID uuid.UUID) (updated int, err error) {
	if len(agentIDs) == 0 {
		return 0, nil
	}
	query := `
	UPDATE agents
	SET group_id = $2
	WHERE id = ANY($1);
	`
	res, err := r.conn(ctx).Exec(ctx, query, agentIDs, groupID)
	if IsInvalidForeignKey(err) {
		return 0, ErrInvalidGroup
	}
	if err != nil {
		return 0, fmt.Errorf("update failed: %w", err)
	}

	return int(res.RowsAffected()), nil
}
