import { render, screen } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import TerritoryTable from './TerritoryTable'
import { useFactionsQuery, useTerritoriesQuery } from '../hooks/useGameData'

vi.mock('../hooks/useGameData', () => ({
    useFactionsQuery: vi.fn(),
    useTerritoriesQuery: vi.fn(),
}))

const mockedUseFactionsQuery = vi.mocked(useFactionsQuery)
const mockedUseTerritoriesQuery = vi.mocked(useTerritoriesQuery)

beforeEach(() => {
    mockedUseFactionsQuery.mockReset()
    mockedUseTerritoriesQuery.mockReset()
    mockedUseFactionsQuery.mockReturnValue({ data: [], isLoading: false } as never)
})

describe('TerritoryTable', () => {
    it('prompts for a scenario when none is selected', () => {
        mockedUseTerritoriesQuery.mockReturnValue({ data: undefined, isLoading: false } as never)
        render(<TerritoryTable tick={0} scenario="" />)
        expect(screen.getByText('Select a scenario to view territories.')).toBeInTheDocument()
    })

    it('shows a loading state while the first fetch is pending', () => {
        mockedUseTerritoriesQuery.mockReturnValue({ data: undefined, isLoading: true } as never)
        render(<TerritoryTable tick={0} scenario="simple" />)
        expect(screen.getByText('Loading territories…')).toBeInTheDocument()
    })

    it('shows an empty-state row when the scenario has no territories', () => {
        mockedUseTerritoriesQuery.mockReturnValue({ data: [], isLoading: false } as never)
        render(<TerritoryTable tick={0} scenario="simple" />)
        expect(screen.getByText('No territories yet.')).toBeInTheDocument()
    })

    it('labels an unclaimed territory (owner < 0) as "Unclaimed"', () => {
        mockedUseTerritoriesQuery.mockReturnValue({
            data: [{ id: 1, owner: -1, factions: [], deposit_types: [], deposits: [] }],
            isLoading: false,
        } as never)
        render(<TerritoryTable tick={0} scenario="simple" />)
        expect(screen.getByText('Unclaimed')).toBeInTheDocument()
        expect(screen.getByText('#001')).toBeInTheDocument()
    })

    it('resolves an owner id to the matching faction name', () => {
        mockedUseFactionsQuery.mockReturnValue({
            data: [{ id: 0, name: 'faction_a', territories: [], drones: [], resources: {} }],
            isLoading: false,
        } as never)
        mockedUseTerritoriesQuery.mockReturnValue({
            data: [{ id: 1, owner: 0, factions: [0], deposit_types: ['resource_mineral'], deposits: [5] }],
            isLoading: false,
        } as never)

        render(<TerritoryTable tick={0} scenario="simple" />)

        expect(screen.getByText('faction_a')).toBeInTheDocument()
        expect(screen.getByText('MINERAL')).toBeInTheDocument()
    })

    it('falls back to "Faction #<id>" when the owner id has no matching faction', () => {
        mockedUseFactionsQuery.mockReturnValue({ data: [], isLoading: false } as never)
        mockedUseTerritoriesQuery.mockReturnValue({
            data: [{ id: 2, owner: 7, factions: [], deposit_types: [], deposits: [] }],
            isLoading: false,
        } as never)

        render(<TerritoryTable tick={0} scenario="simple" />)

        expect(screen.getByText('Faction #7')).toBeInTheDocument()
    })
})
