package agentregistry

import (
	"context"
	"time"

	"github.com/google/uuid"
	agent_model "github.com/nougght/monitoring-system/server/internal/model/agent"
)

type AgentRepository interface {
	CreateAgent(ctx context.Context, agent *agent_model.Agent) (*agent_model.Agent, error)
	GetAgentByID(ctx context.Context, id uuid.UUID) (res *agent_model.Agent, err error)
	GetAgentNamesByIDs(ctx context.Context, agentIDs []uuid.UUID) (map[uuid.UUID]string, error)
	GetAllAgents(ctx context.Context) (res []*agent_model.Agent, err error)
	GetTotalAgents(ctx context.Context) (int, error)
	UpdateStatus(ctx context.Context, agentID uuid.UUID, status agent_model.AgentStatus) error
}

type EnrollmentKeysRepository interface {
	CreateKey(ctx context.Context, key *agent_model.EnrollmentKey) (*agent_model.EnrollmentKey, error)
	GetKeyByAgentId(ctx context.Context, agentID uuid.UUID) (*agent_model.EnrollmentKey, error)
	GetKeyBySelector(ctx context.Context, selector string) (*agent_model.EnrollmentKey, error)
	SetUsed(ctx context.Context, agentID uuid.UUID, usedAt time.Time) error
}

type SpecsRepository interface {
	CreateOrUpdateSpecs(ctx context.Context, specs *agent_model.Specs) (*agent_model.Specs, error)
	GetCurrentSpecs(ctx context.Context, agentID uuid.UUID) (specs *agent_model.Specs, err error)
	GetSpecsTotalList(ctx context.Context, agentIDs []uuid.UUID) (map[uuid.UUID]agent_model.SpecsTotal, error)
}
