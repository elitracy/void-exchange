import { useQuery, useQueryClient } from "@tanstack/react-query"
import { GetDrones, GetFactions, GetTerritories } from "../../wailsjs/go/main/App"

// Shared query keys so TerritoryTable, FactionTable, and DroneTable dedupe
// their network calls through react-query's cache instead of double-fetching.
export function useFactionsQuery(tick: number, scenario: string) {
    return useQuery({
        queryKey: ["factions", scenario, tick],
        queryFn: GetFactions,
        enabled: !!scenario,
    })
}

export function useTerritoriesQuery(tick: number, scenario: string) {
    return useQuery({
        queryKey: ["territories", scenario, tick],
        queryFn: GetTerritories,
        enabled: !!scenario,
    })
}

export function useDronesQuery(tick: number, scenario: string) {
    return useQuery({
        queryKey: ["drones", scenario, tick],
        queryFn: GetDrones,
        enabled: !!scenario,
    })
}

// Dispatch/recall mutate backend state instantly (no travel time), so the
// UI needs to refetch drones/territories right away rather than waiting
// for the next tick's query key to change.
export function useInvalidateGameData(scenario: string) {
    const queryClient = useQueryClient()
    return () => {
        queryClient.invalidateQueries({ queryKey: ["drones", scenario] })
        queryClient.invalidateQueries({ queryKey: ["territories", scenario] })
    }
}
