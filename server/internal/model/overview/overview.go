package overview_model

import (
	agent_model "github.com/nougght/monitoring-system/server/internal/model/agent"
)

type AgentsOverview struct {
	Summary AgentsSummary `json:"summary"`
	TopN    TopNAgents    `json:"top_n"`
}

type AgentsSummary struct {
	TotalAgents             int
	OnlineAgents            int
	AverageCPUUsage         float64
	AverageMemoryUsage      float64
	CPUUsageDistribution    CPUUsageDistribution
	MemoryUsageDistribution MemoryUsageDistribution
	DiskUsageDistribution   DiskUsageDistribution
}

type AgentWithCPUUsage struct {
	agent_model.AgentShort
	CPUUsage float64
}

type AgentWithMemoryUsage struct {
	agent_model.AgentShort
	MemoryUsage float64
}

type TopNAgents struct {
	CPUUsage    []AgentWithCPUUsage    `json:"cpu_usage"`
	MemoryUsage []AgentWithMemoryUsage `json:"memory_usage"`
}

type CPUUsageDistribution struct {
	CountDistribution
}
type MemoryUsageDistribution struct {
	CountDistribution
}
type DiskUsageDistribution struct {
	CountDistribution
}

const (
	HighCPUUsagePercent      = 80
	MediumCPUUsagePercent    = 50
	LowCPUUsagePercent       = 20
	HighMemoryUsagePercent   = 80
	MediumMemoryUsagePercent = 50
	LowMemoryUsagePercent    = 20
	HighDiskUsagePercent     = 80
	MediumDiskUsagePercent   = 50
	LowDiskUsagePercent      = 20
)

type CountDistribution struct {
	High   CountByPercent
	Medium CountByPercent
	Low    CountByPercent
}

func (d *CountDistribution) Add(percent float64, count int) {
	switch {
	case percent >= d.High.Percent:
		d.High.Count += count
	case percent >= d.Medium.Percent:
		d.Medium.Count += count
	default:
		d.Low.Count += count
	}
}

type CountByPercent struct {
	Percent float64
	Count   int
}
