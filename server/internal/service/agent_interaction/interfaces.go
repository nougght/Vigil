package agent

import (
	"context"

	"github.com/google/uuid"
	agent_model "github.com/nougght/monitoring-system/server/internal/model/agent"
	metrics_model "github.com/nougght/monitoring-system/server/internal/model/metrics"
)

type MetricsService interface {
	HandleActivityUpdate(ctx context.Context, update *metrics_model.ActivityUpdate) error
	HandleMetrics(ctx context.Context, agentID uuid.UUID, metrics metrics_model.MetricsBatch) error
}

type AgentRegistryService interface {

	// CreateAgent(ctx context.Context, name string, description *string) (*agent_model.CreateAgentResult, error)
	CreateSession(agentID uuid.UUID) *agent_model.AgentSession
	Enroll(ctx context.Context, params *agent_model.EnrollParams) (*agent_model.EnrollResult, error)
	// GenerateAgentSetupConfig(ctx context.Context, agentID uuid.UUID, enrollmentKey string) *agent_model.AgentSetupConfig
	// GetAgentByID(ctx context.Context, id uuid.UUID) (*agent_model.Agent, error)
	// GetAgentNamesByIDs(ctx context.Context, agentIDs []uuid.UUID) (map[uuid.UUID]string, error)
	// GetAllAgents(ctx context.Context) ([]*agent_model.Agent, error)
	// GetNewAgentFiles(ctx context.Context, agentID uuid.UUID, enrollmentKey string) ([]byte, error)
	GetSession(agentID uuid.UUID) (*agent_model.AgentSession, bool)
	// GetSpecifications(ctx context.Context, agentID uuid.UUID) (*agent_model.Specs, error)
	// GetSpecsTotalList(ctx context.Context, agentIDs []uuid.UUID) (map[uuid.UUID]agent_model.SpecsTotal, error)
	// GetTotalAgents(ctx context.Context) (int, error)
	// IsOnline(agentID uuid.UUID) bool
	// MoveAgentsToGroup(ctx context.Context, agentIDs []uuid.UUID, groupID uuid.UUID) (updated int, err error)
	// OnlineList() []uuid.UUID
	RemoveSession(agentID uuid.UUID)
	UpdateSpecifications(ctx context.Context, agentID uuid.UUID, specifications *agent_model.Specs) error
}
