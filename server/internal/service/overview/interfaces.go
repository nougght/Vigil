package overview

import (
	"context"

	"github.com/google/uuid"
	agent_model "github.com/nougght/monitoring-system/server/internal/model/agent"
	metrics_model "github.com/nougght/monitoring-system/server/internal/model/metrics"
)

type MetricsService interface {
	GetSnapshots() map[uuid.UUID]*metrics_model.Snapshot
}

type AgentRegistryService interface {
	GetAgentNamesByIDs(ctx context.Context, agentIDs []uuid.UUID) (map[uuid.UUID]string, error)
	GetSpecsTotalList(ctx context.Context, agentIDs []uuid.UUID) (map[uuid.UUID]agent_model.SpecsTotal, error)
	GetTotalAgents(ctx context.Context) (int, error)
	IsOnline(agentID uuid.UUID) bool
}
