

export interface CountDistribution {
    high: CountByPercent;
    medium: CountByPercent;
    low: CountByPercent;
}

export const CountDistToPieList = (dist: CountDistribution) => {
    return [
        {
            Name: `>${dist.high.percent}`,
            Value: dist.high.count,
            fill: "#c94b4b",
        },
        {
            Name: `>${dist.medium.percent}`,
            Value: dist.medium.count,
            fill: "#c9ac4b",
        },
        {
            Name: `>${dist.low.percent}`,
            Value: dist.low.count,
            fill: "#4bc95c",
        }
    ]
}

export interface CountByPercent {
    percent: number;
    count: number;
}

export interface AgentsSummary {
    averageCPUUsage: number;
    averageMemoryUsage: number;
    cpuUsageDistribution: CountDistribution;
    memoryUsageDistribution: CountDistribution;
    diskUsageDistribution: CountDistribution;
    onlineAgents: number;
    totalAgents: number;
}

export interface TopNAgents {
    cpuUsage: AgentWithCPUUsage[];
    memoryUsage: AgentWithMemoryUsage[];
}

export interface AgentWithCPUUsage {
    id: string;
    name: string;
    cpuUsage: number;
}

export interface AgentWithMemoryUsage {
    id: string;
    name: string;
    memoryUsage: number;
}

export interface Overview {
    summary: AgentsSummary;
    topN: TopNAgents;
}