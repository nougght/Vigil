package overview

import (
	"context"
	"log"
	"maps"
	"slices"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/nougght/monitoring-system/server/internal/config"
	"github.com/nougght/monitoring-system/server/internal/model"
	agent_model "github.com/nougght/monitoring-system/server/internal/model/agent"
	metrics_model "github.com/nougght/monitoring-system/server/internal/model/metrics"
	overview_model "github.com/nougght/monitoring-system/server/internal/model/overview"
	agentregistry "github.com/nougght/monitoring-system/server/internal/service/agent_registry"
	"github.com/nougght/monitoring-system/server/internal/service/metrics"
	"github.com/nougght/monitoring-system/shared/go/util"
)

type OverviewService struct {
	cfg            *config.Config
	transactor     model.Transactor
	metricsService *metrics.MetricsService
	agentRegistry  *agentregistry.AgentRegistryService
	overviewCache  *overview_model.AgentsOverview
	mu             sync.RWMutex
}

func NewOverviewService(cfg *config.Config,
	transactor model.Transactor,
	metricsService *metrics.MetricsService,
	agentRegistry *agentregistry.AgentRegistryService,
) (*OverviewService, error) {
	s := &OverviewService{
		cfg:            cfg,
		transactor:     transactor,
		metricsService: metricsService,
		agentRegistry:  agentRegistry,
		overviewCache:  &overview_model.AgentsOverview{},
		mu:             sync.RWMutex{},
	}
	return s, nil
}

func (s *OverviewService) RunAggregator(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:

				agentsCount, err := s.agentRegistry.GetTotalAgents(ctx)
				if err != nil {
					log.Printf("failed to get total agents: %s", err.Error())
					continue
				}
				onlineAgents := 0

				snapshots := s.metricsService.GetSnapshots()
				agentIDs := slices.Collect(maps.Keys(snapshots))
				totals, err := s.agentRegistry.GetSpecsTotalList(ctx, agentIDs)
				if err != nil {
					log.Printf("failed to get specs total list: %s", err.Error())
					continue
				}

				var cpuSum float64
				var memorySum float64
				var diskSum float64

				var cpuDistribution overview_model.CPUUsageDistribution
				var memoryDistribution overview_model.MemoryUsageDistribution
				var diskDistribution overview_model.DiskUsageDistribution

				// TODO: dynamic topN count
				n := 5
				cpuSorted := util.Map(slices.Collect(maps.Values(snapshots)),
					func(snapshot *metrics_model.Snapshot) overview_model.AgentWithCPUUsage {
						return overview_model.AgentWithCPUUsage{
							AgentShort: agent_model.AgentShort{
								ID: snapshot.AgentID,
							},
							CPUUsage: snapshot.GetCPUUsage(),
						}
					})
				for i := len(cpuSorted); i <= n; i++ {
					cpuSorted = append(cpuSorted, overview_model.AgentWithCPUUsage{
						AgentShort: agent_model.AgentShort{
							ID: uuid.Nil,
						},
						CPUUsage: 0,
					})
				}
				topNCpuUsage := cpuSorted[:n]
				memorySorted := util.Map(slices.Collect(maps.Values(snapshots)),
					func(snapshot *metrics_model.Snapshot) overview_model.AgentWithMemoryUsage {
						return overview_model.AgentWithMemoryUsage{
							AgentShort: agent_model.AgentShort{
								ID: snapshot.AgentID,
							},
							MemoryUsage: snapshot.GetMemoryUsage(),
						}
					})
				for i := len(memorySorted); i <= n; i++ {
					memorySorted = append(memorySorted, overview_model.AgentWithMemoryUsage{
						AgentShort: agent_model.AgentShort{
							ID: uuid.Nil,
						},
						MemoryUsage: 0,
					})
				}
				topNMemoryUsage := memorySorted[:n]

				for agentID, snapshot := range snapshots {
					if s.agentRegistry.IsOnline(agentID) {
						onlineAgents++
					}
					cpuUsage := snapshot.GetCPUUsage()
					memoryUsage := snapshot.GetMemoryUsage()
					diskUsage := snapshot.GetDiskUsage()

					cpuSum += cpuUsage
					memorySum += memoryUsage
					diskSum += diskUsage

					cpuDistribution.Add(cpuUsage, 1)
					if total := totals[agentID].TotalMemory; total > 0 {
						memoryDistribution.Add(memoryUsage/float64(total)*100, 1)
					}
					if total := totals[agentID].TotalDisk; total > 0 {
						diskDistribution.Add(diskUsage/float64(total)*100, 1)
					}
				}

				s.mu.Lock()
				s.overviewCache.Summary.TotalAgents = agentsCount
				s.overviewCache.Summary.OnlineAgents = onlineAgents
				if agentsCount > 0 {
					s.overviewCache.Summary.AverageCPUUsage = cpuSum / float64(agentsCount)
					s.overviewCache.Summary.AverageMemoryUsage = memorySum / float64(agentsCount)
				} else {
					s.overviewCache.Summary.AverageCPUUsage = 0
					s.overviewCache.Summary.AverageMemoryUsage = 0
				}

				s.overviewCache.Summary.CPUUsageDistribution = cpuDistribution
				s.overviewCache.Summary.MemoryUsageDistribution = memoryDistribution
				s.overviewCache.Summary.DiskUsageDistribution = diskDistribution

				s.overviewCache.TopN.CPUUsage = topNCpuUsage
				s.overviewCache.TopN.MemoryUsage = topNMemoryUsage
				s.mu.Unlock()
			}
		}
	}()
}

func (s *OverviewService) GetOverview(ctx context.Context) (*overview_model.AgentsOverview, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.overviewCache, nil
}
