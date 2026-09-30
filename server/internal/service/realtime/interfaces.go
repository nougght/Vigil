package realtime

import (
	"context"

	overview_model "github.com/nougght/monitoring-system/server/internal/model/overview"
)

type OverviewService interface {
	GetOverview(ctx context.Context) (*overview_model.AgentsOverview, error)
	// RunAggregator(ctx context.Context)
}
