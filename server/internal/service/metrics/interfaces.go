package metrics

import (
	"context"

	metrics_model "github.com/nougght/monitoring-system/server/internal/model/metrics"
	metrickinds "github.com/nougght/monitoring-system/shared/go/metric_kinds"
)

type MetricsRepository interface {
	SaveRows(ctx context.Context, rows []metrics_model.MetricRow) error
	UpsertMetricKinds(ctx context.Context, kinds []metrickinds.MetricKindInfo) error
}

type SeriesProvider interface {
	LoadAllSeries(ctx context.Context) error
	ResolveAllSeries(ctx context.Context, keys []metrics_model.MetricSeriesKey) (resultIDs map[metrics_model.MetricSeriesKey]metrics_model.MetricSeriesID, err error)
}
