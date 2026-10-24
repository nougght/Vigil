import {DefaultId} from "./common.ts";

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

export const IsDefaultGroup = (agent?: Agent)=> {
  if (!agent) {
    return false;
  }
  return agent.groupID == DefaultId;
}