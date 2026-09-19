import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import App from './App'
import {
    CurrentTick,
    GetFactions,
    GetTerritories,
    ListScenarios,
    LoadScenario,
    Pause,
    Resume,
    StartRun,
} from '../wailsjs/go/main/App'

vi.mock('../wailsjs/go/main/App', () => ({
    CurrentTick: vi.fn(),
    ListScenarios: vi.fn(),
    LoadScenario: vi.fn(),
    Pause: vi.fn(),
    Resume: vi.fn(),
    StartRun: vi.fn(),
    GetFactions: vi.fn(),
    GetTerritories: vi.fn(),
}))

const mocked = {
    CurrentTick: vi.mocked(CurrentTick),
    ListScenarios: vi.mocked(ListScenarios),
    LoadScenario: vi.mocked(LoadScenario),
    Pause: vi.mocked(Pause),
    Resume: vi.mocked(Resume),
    StartRun: vi.mocked(StartRun),
    GetFactions: vi.mocked(GetFactions),
    GetTerritories: vi.mocked(GetTerritories),
}

function renderApp() {
    const queryClient = new QueryClient({
        defaultOptions: { queries: { retry: false } },
    })
    return render(
        <QueryClientProvider client={queryClient}>
            <App />
        </QueryClientProvider>,
    )
}

beforeEach(() => {
    Object.values(mocked).forEach(fn => fn.mockReset())
    mocked.ListScenarios.mockResolvedValue(['simple', 'medium'])
    mocked.LoadScenario.mockResolvedValue(undefined)
    mocked.StartRun.mockResolvedValue(undefined)
    mocked.Pause.mockResolvedValue(undefined)
    mocked.Resume.mockResolvedValue(undefined)
    mocked.CurrentTick.mockResolvedValue(0)
    mocked.GetFactions.mockResolvedValue([])
    mocked.GetTerritories.mockResolvedValue([])
})

describe('App', () => {
    it('lists scenarios on load and shows idle status with Start disabled', async () => {
        renderApp()

        expect(await screen.findByText('simple')).toBeInTheDocument()
        expect(screen.getByText('medium')).toBeInTheDocument()
        expect(screen.getByText('idle')).toBeInTheDocument()
        expect(screen.getByRole('button', { name: 'Start' })).toBeDisabled()
    })

    it('loading a scenario enables Start and marks it active', async () => {
        const user = userEvent.setup()
        renderApp()

        const scenarioButton = await screen.findByText('simple')
        await user.click(scenarioButton)

        await waitFor(() => expect(mocked.LoadScenario).toHaveBeenCalledWith('simple', 0))
        expect(screen.getByRole('button', { name: 'Start' })).toBeEnabled()
        expect(screen.getByText('paused')).toBeInTheDocument()
    })

    it('starting a run flips status to running and disables Start while enabling Pause', async () => {
        const user = userEvent.setup()
        renderApp()

        await user.click(await screen.findByText('simple'))
        await user.click(screen.getByRole('button', { name: 'Start' }))

        await waitFor(() => expect(mocked.StartRun).toHaveBeenCalledWith(100))
        expect(screen.getByText('running')).toBeInTheDocument()
        expect(screen.getByRole('button', { name: 'Start' })).toBeDisabled()
        expect(screen.getByRole('button', { name: 'Pause' })).toBeEnabled()
    })

    it('pausing a running scenario flips status back to paused', async () => {
        const user = userEvent.setup()
        renderApp()

        await user.click(await screen.findByText('simple'))
        await user.click(screen.getByRole('button', { name: 'Start' }))
        await waitFor(() => expect(screen.getByText('running')).toBeInTheDocument())

        await user.click(screen.getByRole('button', { name: 'Pause' }))

        await waitFor(() => expect(mocked.Pause).toHaveBeenCalledTimes(1))
        expect(screen.getByText('paused')).toBeInTheDocument()
    })

    it('shows "No scenarios found." when the backend returns an empty list', async () => {
        mocked.ListScenarios.mockResolvedValue([])
        renderApp()
        expect(await screen.findByText('No scenarios found.')).toBeInTheDocument()
    })
})
