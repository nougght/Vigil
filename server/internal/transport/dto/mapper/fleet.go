package mapper

import (
	overview_model "github.com/nougght/monitoring-system/server/internal/model/overview"
	dto "github.com/nougght/monitoring-system/server/internal/transport/dto/types"
)

func AgentWithCPUUsageToDTO(domain *overview_model.AgentWithCPUUsage) (res *dto.AgentWithCPUUsageDTO) {
	if domain == nil {
		return
	}

	res = &dto.AgentWithCPUUsageDTO{
		AgentShortDTO: *AgentShortToDTO(&domain.AgentShort),
		CPUUsage:      domain.CPUUsage,
	}
	return
}

func AgentWithMemoryUsageToDTO(domain *overview_model.AgentWithMemoryUsage) (res *dto.AgentWithMemoryUsageDTO) {
	if domain == nil {
		return
	}

	res = &dto.AgentWithMemoryUsageDTO{
		AgentShortDTO: *AgentShortToDTO(&domain.AgentShort),
		MemoryUsage:   domain.MemoryUsage,
	}
	return
}

func TopNAgentsToDTO(domain *overview_model.TopNAgents) (res *dto.TopNAgents) {
	if domain == nil {
		return
	}
	res = &dto.TopNAgents{
		CPUUsage:    make([]dto.AgentWithCPUUsageDTO, len(domain.CPUUsage)),
		MemoryUsage: make([]dto.AgentWithMemoryUsageDTO, len(domain.MemoryUsage)),
	}
	for i := range domain.CPUUsage {
		res.CPUUsage[i] = *AgentWithCPUUsageToDTO(&domain.CPUUsage[i])
	}
	for i := range domain.MemoryUsage {
		res.MemoryUsage[i] = *AgentWithMemoryUsageToDTO(&domain.MemoryUsage[i])
	}
	return
}

func SummaryToDTO(domain *overview_model.AgentsSummary) (res *dto.AgentsSummary) {
	if domain == nil {
		return
	}
	res = &dto.AgentsSummary{
		TotalAgents:        domain.TotalAgents,
		OnlineAgents:       domain.OnlineAgents,
		AverageCPUUsage:    domain.AverageCPUUsage,
		AverageMemoryUsage: domain.AverageMemoryUsage,
		CPUUsageDistribution: dto.CountDistribution{
			High:   dto.CountByPercent{Percent: domain.CPUUsageDistribution.High.Percent, Count: domain.CPUUsageDistribution.High.Count},
			Medium: dto.CountByPercent{Percent: domain.CPUUsageDistribution.Medium.Percent, Count: domain.CPUUsageDistribution.Medium.Count},
			Low:    dto.CountByPercent{Percent: domain.CPUUsageDistribution.Low.Percent, Count: domain.CPUUsageDistribution.Low.Count},
		},
		MemoryUsageDistribution: dto.CountDistribution{
			High:   dto.CountByPercent{Percent: domain.MemoryUsageDistribution.High.Percent, Count: domain.MemoryUsageDistribution.High.Count},
			Medium: dto.CountByPercent{Percent: domain.MemoryUsageDistribution.Medium.Percent, Count: domain.MemoryUsageDistribution.Medium.Count},
			Low:    dto.CountByPercent{Percent: domain.MemoryUsageDistribution.Low.Percent, Count: domain.MemoryUsageDistribution.Low.Count},
		},
		DiskUsageDistribution: dto.CountDistribution{
			High:   dto.CountByPercent{Percent: domain.DiskUsageDistribution.High.Percent, Count: domain.DiskUsageDistribution.High.Count},
			Medium: dto.CountByPercent{Percent: domain.DiskUsageDistribution.Medium.Percent, Count: domain.DiskUsageDistribution.Medium.Count},
			Low:    dto.CountByPercent{Percent: domain.DiskUsageDistribution.Low.Percent, Count: domain.DiskUsageDistribution.Low.Count},
		},
	}
	return
}

func OverviewToDTO(domain *overview_model.AgentsOverview) (res *dto.AgentsOverview) {
	if domain == nil {
		return
	}
	res = &dto.AgentsOverview{
		Summary: SummaryToDTO(&domain.Summary),
		TopN:    TopNAgentsToDTO(&domain.TopN),
	}
	return
}
