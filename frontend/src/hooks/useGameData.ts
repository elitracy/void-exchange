import { keepPreviousData, useQuery, useQueryClient } from "@tanstack/react-query"
import { GetDrones, GetFactions, GetTerritories } from "../../wailsjs/go/main/App"

// Shared query keys so TerritoryTable, FactionTable, and DroneTable dedupe
// their network calls through react-query's cache instead of double-fetching.
//
// `tick` changes on every refresh, which makes each refresh a "new" query
// key. Without placeholderData, react-query would report `data: undefined`
// for that instant (even though the previous tick's rows are still valid),
// which is what caused the table to blank out and flicker between ticks.
// keepPreviousData keeps rendering the last successful result until the
// new one lands, so the UI only ever updates in place.
export function useFactionsQuery(tick: number, scenario: string) {
    return useQuery({
        queryKey: ["factions", scenario, tick],
        queryFn: GetFactions,
        enabled: !!scenario,
        placeholderData: keepPreviousData,
    })
}

export function useTerritoriesQuery(tick: number, scenario: string) {
    return useQuery({
        queryKey: ["territories", scenario, tick],
        queryFn: GetTerritories,
        enabled: !!scenario,
        placeholderData: keepPreviousData,
    })
}

export function useDronesQuery(tick: number, scenario: string) {
    return useQuery({
        queryKey: ["drones", scenario, tick],
        queryFn: GetDrones,
        enabled: !!scenario,
        placeholderData: keepPreviousData,
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
