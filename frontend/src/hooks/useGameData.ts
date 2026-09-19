import { useQuery } from "@tanstack/react-query"
import { GetFactions, GetTerritories } from "../../wailsjs/go/main/App"

// Shared query keys so TerritoryTable and FactionTable dedupe their
// network calls through react-query's cache instead of double-fetching.
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
