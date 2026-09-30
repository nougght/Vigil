package rest

import (
	"context"

	"github.com/google/uuid"
	agent_model "github.com/nougght/monitoring-system/server/internal/model/agent"
	overview_model "github.com/nougght/monitoring-system/server/internal/model/overview"
)

type AgentRegistryService interface {
	CreateAgent(ctx context.Context, name string, description *string) (*agent_model.CreateAgentResult, error)
	// CreateSession(agentID uuid.UUID) *agent_model.AgentSession
	// Enroll(ctx context.Context, params *agent_model.EnrollParams) (*agent_model.EnrollResult, error)
	// GenerateAgentSetupConfig(ctx context.Context, agentID uuid.UUID, enrollmentKey string) *agent_model.AgentSetupConfig
	GetAgentByID(ctx context.Context, id uuid.UUID) (*agent_model.Agent, error)
	// GetAgentNamesByIDs(ctx context.Context, agentIDs []uuid.UUID) (map[uuid.UUID]string, error)
	GetAllAgents(ctx context.Context) ([]*agent_model.Agent, error)
	GetNewAgentFiles(ctx context.Context, agentID uuid.UUID, enrollmentKey string) ([]byte, error)
	// GetSession(agentID uuid.UUID) (*agent_model.AgentSession, bool)
	GetSpecifications(ctx context.Context, agentID uuid.UUID) (*agent_model.Specs, error)
	// GetSpecsTotalList(ctx context.Context, agentIDs []uuid.UUID) (map[uuid.UUID]agent_model.SpecsTotal, error)
	// GetTotalAgents(ctx context.Context) (int, error)
	// IsOnline(agentID uuid.UUID) bool
	// OnlineList() []uuid.UUID
	// RemoveSession(agentID uuid.UUID)
	// UpdateSpecifications(ctx context.Context, agentID uuid.UUID, specifications *agent_model.Specs) error
}

type AgentInteractionService interface {
	// Enroll(ctx context.Context, params *agent_model.EnrollParams) (*agent_model.EnrollResult, error)
	// HandleActivityUpdate(ctx context.Context, update *metrics_model.ActivityUpdate) error
	// HandleConnection(agentID uuid.UUID)
	// HandleDisconnection(agentID uuid.UUID)
	// HandleFrame(frame []byte, agentID uuid.UUID)
	// HandleMetricsBatch(ctx context.Context, batch *metrics_model.MetricsBatch) error
	// HandleSpecifications(ctx context.Context, agentID uuid.UUID, specifications *agent_model.Specs) error
	// RequestSpecifications(ctx context.Context, agentID uuid.UUID) (*agent_model.Specs, error)
	// SetRequester(requester agent.Requester)
	SubStreaming(agentID uuid.UUID, viewerID uuid.UUID) (<-chan []byte, error)
	UnsubAllStreaming(viewerID uuid.UUID)
}

type OverviewService interface {
	GetOverview(ctx context.Context) (*overview_model.AgentsOverview, error)
	// RunAggregator(ctx context.Context)
}
