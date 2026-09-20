import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { renderHook, waitFor } from '@testing-library/react'
import type { ReactNode } from 'react'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { useFactionsQuery, useTerritoriesQuery } from './useGameData'

vi.mock('../../wailsjs/go/main/App', () => ({
    GetFactions: vi.fn(),
    GetTerritories: vi.fn(),
}))

import { GetFactions, GetTerritories } from '../../wailsjs/go/main/App'

function wrapper({ children }: { children: ReactNode }) {
    const queryClient = new QueryClient({
        defaultOptions: { queries: { retry: false } },
    })
    return (
        <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
    )
}

beforeEach(() => {
    vi.mocked(GetFactions).mockReset()
    vi.mocked(GetTerritories).mockReset()
})

describe('useFactionsQuery', () => {
    it('does not call GetFactions when no scenario is selected', () => {
        renderHook(() => useFactionsQuery(0, ''), { wrapper })
        expect(GetFactions).not.toHaveBeenCalled()
    })

    it('calls GetFactions once a scenario is selected and returns its data', async () => {
        vi.mocked(GetFactions).mockResolvedValue([
            { id: 0, name: 'faction_a', territories: [], drones: [], resources: {} },
        ])

        const { result } = renderHook(() => useFactionsQuery(0, 'simple'), { wrapper })

        await waitFor(() => expect(result.current.data).toBeDefined())
        expect(GetFactions).toHaveBeenCalledTimes(1)
        expect(result.current.data?.[0].name).toBe('faction_a')
    })

    it('re-fetches when tick changes for the same scenario', async () => {
        vi.mocked(GetFactions).mockResolvedValue([])
        const { rerender } = renderHook(
            ({ tick }) => useFactionsQuery(tick, 'simple'),
            { wrapper, initialProps: { tick: 0 } },
        )

        await waitFor(() => expect(GetFactions).toHaveBeenCalledTimes(1))
        rerender({ tick: 1 })
        await waitFor(() => expect(GetFactions).toHaveBeenCalledTimes(2))
    })
})

describe('useTerritoriesQuery', () => {
    it('does not call GetTerritories when no scenario is selected', () => {
        renderHook(() => useTerritoriesQuery(0, ''), { wrapper })
        expect(GetTerritories).not.toHaveBeenCalled()
    })

    it('calls GetTerritories once a scenario is selected', async () => {
        vi.mocked(GetTerritories).mockResolvedValue([])
        renderHook(() => useTerritoriesQuery(0, 'simple'), { wrapper })
        await waitFor(() => expect(GetTerritories).toHaveBeenCalledTimes(1))
    })
})
