package dto

type AgentsOverview struct {
	Summary *AgentsSummary `json:"summary"`
	TopN    *TopNAgents    `json:"topN"`
} //@name AgentsOverview

type AgentsSummary struct {
	TotalAgents             int               `json:"totalAgents"`
	OnlineAgents            int               `json:"onlineAgents"`
	AverageCPUUsage         float64           `json:"averageCPUUsage"`
	AverageMemoryUsage      float64           `json:"averageMemoryUsage"`
	CPUUsageDistribution    CountDistribution `json:"cpuUsageDistribution"`
	MemoryUsageDistribution CountDistribution `json:"memoryUsageDistribution"`
	DiskUsageDistribution   CountDistribution `json:"diskUsageDistribution"`
} //@name AgentsSummary

type CountDistribution struct {
	High   CountByPercent `json:"high"`
	Medium CountByPercent `json:"medium"`
	Low    CountByPercent `json:"low"`
} //@name CountDistribution

type CountByPercent struct {
	Percent float64 `json:"percent"`
	Count   int     `json:"count"`
} //@name CountByPercent

type AgentWithCPUUsageDTO struct {
	AgentShortDTO
	CPUUsage float64 `json:"cpuUsage"`
} //@name AgentWithCPUUsageDTO

type AgentWithMemoryUsageDTO struct {
	AgentShortDTO
	MemoryUsage float64 `json:"memoryUsage"`
} //@name AgentWithMemoryUsageDTO
type TopNAgents struct {
	CPUUsage    []AgentWithCPUUsageDTO    `json:"cpuUsage"`
	MemoryUsage []AgentWithMemoryUsageDTO `json:"memoryUsage"`
} //@name TopNAgents
