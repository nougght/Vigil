import { useQuery } from "@tanstack/react-query";
import type { Overview } from "../domain/overview";
import { convertOverviewFromDTO } from "../api/mapper";
import { getFleetOverview } from "../api/client/monitoringServerAPI";
import type {Error} from "./common"



interface UseOverviewResult {
    overview: Overview | undefined
    error: Error | null
}

const getOverview = async (): Promise<UseOverviewResult> => {
    const resp = await getFleetOverview();
    return resp.status == 200 ?
    {
        overview: convertOverviewFromDTO(resp.data),
        error: null
    } : 
    {
        overview: undefined,
        error: {
            status: resp.status,
            message: JSON.stringify(resp.data)
        }
    }
}


export function useOverview() {
    return useQuery({
        queryKey: ['overview'],
        queryFn: () => getOverview(),
        
        refetchOnWindowFocus: true,
        staleTime: 60000,
        retry: 2,
    });
}