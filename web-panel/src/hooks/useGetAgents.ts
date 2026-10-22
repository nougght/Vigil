
import { useQuery } from '@tanstack/react-query';
import { getAllAgents as getAllAgents } from '../api/client/monitoringServerAPI';
import { convertAgentFromDTO } from '../api/mapper';

interface GetAgentsParams {
  groupIds: string[];
}
export function useAgents(params?:GetAgentsParams, enabled: boolean = true) {
    return useQuery({
        queryKey: ['agents', params],
        queryFn: async () => (await getAllAgents(params)).map(convertAgentFromDTO),
        refetchOnWindowFocus: true,
        enabled: enabled,
        staleTime: 60000,
    });
}