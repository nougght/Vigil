package metrics_model

import (
	"sync"
	"time"

	"github.com/google/uuid"
	metrickinds "github.com/nougght/monitoring-system/shared/go/metric_kinds"
)

type AgentSeriesKey struct {
	Kind  int32
	Label string
}

type MetricState struct {
	Value float64
	Ts    time.Time
	// Ring
}

type Snapshot struct {
	AgentID  uuid.UUID
	LastSeen time.Time

	mu     sync.RWMutex
	values map[AgentSeriesKey]MetricState
}

func NewSnapshot(agentID uuid.UUID) *Snapshot {
	return &Snapshot{
		AgentID: agentID,
		values:  make(map[AgentSeriesKey]MetricState, 10),
	}
}

func (s *Snapshot) UpdateMetrics(metrics []*MetricSample) {
	s.mu.Lock()
	for _, m := range metrics {
		key := AgentSeriesKey{Kind: m.Kind, Label: m.Label}
		s.values[key] = MetricState{
			Value: m.Value,
			Ts:    m.Timestamp,
		}
	}
	s.mu.Unlock()
}

func (s *Snapshot) GetCPUUsage() float64 {
	return s.GetValueByKey(metrickinds.GetKindByKey("cpu_usage"), "")
}

func (s *Snapshot) GetMemoryUsage() float64 {
	return s.GetValueByKey(metrickinds.GetKindByKey("memory_usage"), "")
}

func (s *Snapshot) GetDiskUsage() float64 {
	return s.GetValueByKey(metrickinds.GetKindByKey("disk_usage"), "")
}

func (s *Snapshot) GetValueByKey(kind int32, label string) float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	key := AgentSeriesKey{Kind: kind, Label: label}
	value, ok := s.values[key]
	if !ok {
		return 0
	}
	return value.Value
}
