

export const MessageTypeSeries = "series"


export interface Message {
    type: string
    agentID: string
    payload: SeriesDTO
}


export const ClientMessageTypeAgentDetailed = "agent.detailed"


export interface AgentDetailedMessage {
    agents: string[]
}

export interface ClientMessage {
    type: string
    payload: AgentDetailedMessage
}

export interface SeriesDTO {
    key: number
    label: string
    unit: string
    ts: Date
    values: number[]
}

// const METRIC_KIND_UNSPECIFIED = 0;
const CPU_PERCENT = 1;
const MEM_USED = 30;
const DISK_USED = 50;   // label = точка монтирования
const NET_UPLOAD = 60;   // label = интерфейс
const NET_DOWNLOAD = 70;

export interface Metrics {
    agentID?: string
    focusedWindow?: string
    cpuPercent?: number
    memoryUsed?: number
    diskUsage?: Map<string, number> 
    uploadMbps?: number
    downloadMbps?: number
    // processList?: Process[]
    timestamp?: Date
}

export const FillMetricsFromSeries = (s: SeriesDTO, metrics: Metrics, agentID: string):Metrics => {
    let m: Metrics = {
        cpuPercent: metrics.cpuPercent,
        memoryUsed: metrics.memoryUsed,
        diskUsage: metrics.diskUsage,
        uploadMbps: metrics.uploadMbps,
        downloadMbps: metrics.downloadMbps,
        agentID: agentID
    }
    const last = s.values[s.values.length - 1]
    switch (s.key) {
        case CPU_PERCENT:
            m.cpuPercent = last
            break
        case MEM_USED:
            m.memoryUsed = last
            break
        case DISK_USED:
            m.diskUsage = new Map(m.diskUsage)
            m.diskUsage.set(s.label, last)
            break
        case NET_UPLOAD:
            m.uploadMbps = last
            break
        case NET_DOWNLOAD:
            m.downloadMbps = last
            break
        default:
            console.log("unknown metric key: ", s.key)
            break
    }
    return m
}