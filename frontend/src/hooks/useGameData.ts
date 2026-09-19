import { keepPreviousData, useQuery } from "@tanstack/react-query"
import { GetFactions, GetTerritories } from "../../wailsjs/go/main/App"

// Shared query keys so TerritoryTable and FactionTable dedupe their
// network calls through react-query's cache instead of double-fetching.
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
