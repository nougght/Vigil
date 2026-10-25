package event

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nougght/monitoring-system/pkg/eventbus"
	metrics_model "github.com/nougght/monitoring-system/server/internal/model/metrics"
	overview_model "github.com/nougght/monitoring-system/server/internal/model/overview"
)

type EventBus interface {
	Start(ctx context.Context)
	Shutdown(ctx context.Context) error
	Subscribe(subject string, handler eventbus.EventHandler) error
	Publish(ctx context.Context, event eventbus.Event) error
}

const (
	AgentMetricsEventName  = "AgentMetricsEvent"
	AgentActivityEventName = "AgentActivityEvent"
	FleetOverviewEventName = "FleetOverviewEvent"
)

const (
	AgentMetricsEventSubjectPrefix  = "agent.detailed"
	AgentActivityEventSubjectPrefix = "agent.activity"
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

type AgentActivityEvent struct {
	AgentID  uuid.UUID
	Activity metrics_model.ActivityUpdate
}

func (e *AgentActivityEvent) Name() string {
	return AgentActivityEventName
}

func (e *AgentActivityEvent) Subject() string {
	return fmt.Sprintf("%s.%s", AgentActivityEventSubjectPrefix, e.AgentID)
}
