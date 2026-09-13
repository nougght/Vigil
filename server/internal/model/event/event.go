package event

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	metrics_model "github.com/nougght/monitoring-system/server/internal/model/metrics"
	overview_model "github.com/nougght/monitoring-system/server/internal/model/overview"
)

type Event interface {
	Name() string
	Subject() string
}

type EventHandler func(ctx context.Context, event Event)

type EventBus interface {
	Start(ctx context.Context)
	Shutdown(ctx context.Context) error
	Subscribe(subject string, handler EventHandler) error
	Publish(ctx context.Context, event Event) error
}

const (
	AgentMetricsEventName  = "AgentMetricsEvent"
	FleetOverviewEventName = "FleetOverviewEvent"
)

const (
	AgentMetricsEventSubjectPrefix  = "agent.detailed"
	FleetOverviewEventSubjectPrefix = "fleet.overview"
)

type AgentMetricsEvent struct {
	AgentID uuid.UUID
	Metric  metrics_model.MetricSample
}

func (e *AgentMetricsEvent) Name() string {
	return AgentMetricsEventName
}

func (e *AgentMetricsEvent) Subject() string {
	return fmt.Sprintf("%s.%s", AgentMetricsEventSubjectPrefix, e.AgentID)
}

type FleetOverviewEvent struct {
	Overview overview_model.AgentsOverview
}

func (e *FleetOverviewEvent) Name() string {
	return FleetOverviewEventName
}

func (e *FleetOverviewEvent) Subject() string {
	return FleetOverviewEventSubjectPrefix
}
