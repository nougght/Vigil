package event

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	metrics_model "github.com/nougght/monitoring-system/server/internal/model/metrics"
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

type AgentMetricsEvent struct {
	AgentID uuid.UUID
	Metric  metrics_model.MetricSample
}

func (e *AgentMetricsEvent) Name() string {
	return "AgentMetricsEvent"
}

func (e *AgentMetricsEvent) Subject() string {
	return fmt.Sprintf("agent.detailed.%s", e.AgentID)
}
