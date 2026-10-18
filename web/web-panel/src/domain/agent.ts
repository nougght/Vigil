

export interface Agent {
    id: string
    name: string
    description?: string
    createdAt: Date
    status?: string
    lastSeenAt?: Date
    isOnline: boolean
}