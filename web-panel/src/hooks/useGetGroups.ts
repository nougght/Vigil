import {useQuery, useQueryClient} from "@tanstack/react-query";
import {getAllAgentGroups, getAgentGroupByID} from "../api/client/monitoringServerAPI.ts";
import {convertGroupFromDTO} from "../api/mapper.ts";
import type {Group} from "../domain/group.ts";


export function useGroup(id: string) {
  const queryClient = useQueryClient();
  return useQuery({
    queryKey: [`group-${id}`],
    queryFn: async () => convertGroupFromDTO(await getAgentGroupByID(id)),
    // try get data from cache
    initialData: () => {
      const itemsList = queryClient.getQueryData<Group[]>(['groups']);
      return itemsList?.find((item) => item.id === id);
    },
    staleTime: 60000,
  });
}

export function useGroups() {
    return useQuery({
        queryKey: ['groups'],
        queryFn: async () => (await getAllAgentGroups()).map(convertGroupFromDTO),
        refetchOnWindowFocus: true,
        staleTime: 60000,
    });
}