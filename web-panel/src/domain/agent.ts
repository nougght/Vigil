import {  validate } from 'uuid';

export interface Agent {
    id: string
    name: string
    description?: string
    groupID: string
    createdAt: Date
    status?: string
    lastSeenAt?: Date
    isOnline: boolean
}

export const IsDefaultGroup = (agent: Agent)=> {
  return !validate(agent.groupID)
}