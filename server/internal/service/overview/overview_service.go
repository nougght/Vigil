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

				cpuDistribution := overview_model.CPUUsageDistribution{
					CountDistribution: overview_model.CountDistribution{
						High:   overview_model.CountByPercent{Percent: overview_model.HighCPUUsagePercent},
						Medium: overview_model.CountByPercent{Percent: overview_model.MediumCPUUsagePercent},
						Low:    overview_model.CountByPercent{Percent: overview_model.LowCPUUsagePercent},
					},
				}
				memoryDistribution := overview_model.MemoryUsageDistribution{
					CountDistribution: overview_model.CountDistribution{
						High:   overview_model.CountByPercent{Percent: overview_model.HighMemoryUsagePercent},
						Medium: overview_model.CountByPercent{Percent: overview_model.MediumMemoryUsagePercent},
						Low:    overview_model.CountByPercent{Percent: overview_model.LowMemoryUsagePercent},
					},
				}
				diskDistribution := overview_model.DiskUsageDistribution{
					CountDistribution: overview_model.CountDistribution{
						High:   overview_model.CountByPercent{Percent: overview_model.HighDiskUsagePercent},
						Medium: overview_model.CountByPercent{Percent: overview_model.MediumDiskUsagePercent},
						Low:    overview_model.CountByPercent{Percent: overview_model.LowDiskUsagePercent},
					},
				}

				// TODO: dynamic topN count
				n := 5
				cpuSorted := slices.SortedFunc(slices.Values(util.Map(slices.Collect(maps.Values(snapshots)),
					func(snapshot *metrics_model.Snapshot) overview_model.AgentWithCPUUsage {
						return overview_model.AgentWithCPUUsage{
							AgentShort: agent_model.AgentShort{
								ID: snapshot.AgentID,
							},
							CPUUsage: snapshot.GetCPUUsage(),
						}
					})),
					func(a overview_model.AgentWithCPUUsage, b overview_model.AgentWithCPUUsage) int {
						return int(b.CPUUsage - a.CPUUsage)
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
				memorySorted := slices.SortedFunc(slices.Values(util.Map(slices.Collect(maps.Values(snapshots)),
					func(snapshot *metrics_model.Snapshot) overview_model.AgentWithMemoryUsage {
						return overview_model.AgentWithMemoryUsage{
							AgentShort: agent_model.AgentShort{
								ID: snapshot.AgentID,
							},
							MemoryUsage: snapshot.GetMemoryUsage(),
						}
					})),
					func(a overview_model.AgentWithMemoryUsage, b overview_model.AgentWithMemoryUsage) int {
						return int(b.MemoryUsage - a.MemoryUsage)
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

				topAgentIDs := make([]uuid.UUID, 0, n*2)
				for _, ag := range topNCpuUsage {
					topAgentIDs = append(topAgentIDs, ag.ID)
				}
				for _, ag := range topNMemoryUsage {
					topAgentIDs = append(topAgentIDs, ag.ID)
				}
				names, err := s.agentRegistry.GetAgentNamesByIDs(ctx, topAgentIDs)
				if err != nil {
					log.Printf("failed get agent names: %s", err.Error())
				}
				for i, ag := range topNCpuUsage {
					name, ok := names[ag.ID]
					if !ok {
						log.Println("name not found in list")
						continue
					}
					topNCpuUsage[i].Name = name
				}
				for i, ag := range topNMemoryUsage {
					name, ok := names[ag.ID]
					if !ok {
						log.Println("name not found in list")
						continue
					}
					topNMemoryUsage[i].Name = name
				}
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
					} else {
						log.Println("memory total is 0")
					}
					if total := totals[agentID].TotalDisk; total > 0 {
						diskDistribution.Add(diskUsage/float64(total)*100, 1)
					} else {
						log.Println("disk total is 0")
					}
				}

				s.mu.Lock()
				s.overviewCache.Summary.TotalAgents = agentsCount
				s.overviewCache.Summary.OnlineAgents = onlineAgents
				if onlineAgents > 0 {
					s.overviewCache.Summary.AverageCPUUsage = cpuSum / float64(onlineAgents)
					s.overviewCache.Summary.AverageMemoryUsage = memorySum / float64(onlineAgents)
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
