import { render, screen } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import FactionTable from './FactionTable'
import { useFactionsQuery } from '../hooks/useGameData'

vi.mock('../hooks/useGameData', () => ({
    useFactionsQuery: vi.fn(),
}))

const mockedUseFactionsQuery = vi.mocked(useFactionsQuery)

beforeEach(() => {
    mockedUseFactionsQuery.mockReset()
})

describe('FactionTable', () => {
    it('prompts for a scenario when none is selected', () => {
        mockedUseFactionsQuery.mockReturnValue({ data: undefined, isLoading: false } as never)
        render(<FactionTable tick={0} scenario="" />)
        expect(screen.getByText('Select a scenario to view factions.')).toBeInTheDocument()
    })

    it('shows a loading state while the first fetch is pending', () => {
        mockedUseFactionsQuery.mockReturnValue({ data: undefined, isLoading: true } as never)
        render(<FactionTable tick={0} scenario="simple" />)
        expect(screen.getByText('Loading factions…')).toBeInTheDocument()
    })

    it('shows an empty-state row when the scenario has no factions', () => {
        mockedUseFactionsQuery.mockReturnValue({ data: [], isLoading: false } as never)
        render(<FactionTable tick={0} scenario="simple" />)
        expect(screen.getByText('No factions yet.')).toBeInTheDocument()
    })

    it('renders one row per faction with computed stockpile and resource chips', () => {
        mockedUseFactionsQuery.mockReturnValue({
            data: [
                {
                    id: 0,
                    name: 'faction_a',
                    territories: [1, 2],
                    drones: [10],
                    resources: { resource_mineral: 5, resource_gas: 3 },
                },
            ],
            isLoading: false,
        } as never)

        render(<FactionTable tick={0} scenario="simple" />)

        expect(screen.getByText('faction_a')).toBeInTheDocument()
        // territories column
        expect(screen.getByText('2')).toBeInTheDocument()
        // fleet column
        expect(screen.getByText('1')).toBeInTheDocument()
        // stockpile = sum of resource amounts (5 + 3)
        expect(screen.getByText('8')).toBeInTheDocument()
        expect(screen.getByText('MINERAL')).toBeInTheDocument()
        expect(screen.getByText('GAS')).toBeInTheDocument()
    })

    it('treats a faction with no resources as having "No holdings" and zero stockpile', () => {
        mockedUseFactionsQuery.mockReturnValue({
            data: [{ id: 0, name: 'faction_a', territories: [], drones: [], resources: {} }],
            isLoading: false,
        } as never)

        render(<FactionTable tick={0} scenario="simple" />)

        expect(screen.getByText('No holdings')).toBeInTheDocument()
        expect(screen.getAllByText('0')).toHaveLength(3)
    })
})
