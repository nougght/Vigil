import {useQuery, useQueryClient} from '@tanstack/react-query';
import {getAgentByID} from '../api/client/monitoringServerAPI';
import {convertAgentFromDTO} from '../api/mapper';
import type {Agent} from "../domain/agent";


export function useAgent(id: string) {
  const queryClient = useQueryClient();
  return useQuery({
    queryKey: [`agent-${id}`],
    queryFn: async () => convertAgentFromDTO(await getAgentByID(id)),
    initialData: () => {
      const itemsList = queryClient.getQueryData<Agent[]>(['agents']);
      return itemsList?.find((item) => item.id === id);
    },
    staleTime: 60000,
  });
}