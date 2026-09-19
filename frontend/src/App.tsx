import { useState, useMemo, useEffect } from 'react';
import { ListScenarios, LoadScenario, Pause, Resume, StartRun } from "../wailsjs/go/main/App";
import { EventsOn } from "../wailsjs/runtime/runtime";
import TerritoryTable from './components/TerritoryTable';
import FactionTable from './components/FactionTable';
import DroneTable from './components/DroneTable';
import StatTile from './components/ui/StatTile';
import { useFactionsQuery, useTerritoriesQuery } from './hooks/useGameData';
import { formatMissionClock } from './lib/time';

const SIM_TICK_MS = 100
const UI_REFRESH_MS = 1000
// The backend now pushes a "tick" event on every simulation tick instead of
// the UI polling CurrentTick() on its own timer, so the displayed tick can
// no longer drift or skip. Querying territories/factions on every single
// simulation tick would be excessive, so those refresh once every N ticks —
// still driven off the real tick count rather than a separate wall-clock
// timer, so the refresh cadence can't drift out of sync either.
const TICKS_PER_REFRESH = Math.max(1, Math.round(UI_REFRESH_MS / SIM_TICK_MS))

type View = "territories" | "factions" | "drones"

function App() {
    const [scenarios, setScenarios] = useState<string[]>()
    const [currentScenario, setCurrentScenario] = useState<string>()
    const [tick, setTick] = useState(0)
    const [running, setRunning] = useState(false)
    // Tracks whether a run has been started at all, separate from `running`.
    // `running` flips false on pause too, so gating the Resume button on
    // `running` alone made it disable itself the moment you paused. `started`
    // is what actually distinguishes "never started" (idle) from "paused".
    const [started, setStarted] = useState(false)
    const [view, setView] = useState<View>("territories")
    // Stand-in for a future dev tools panel: toggled with the backtick key.
    const [showTickDebug, setShowTickDebug] = useState(false)

    const queryTick = Math.floor(tick / TICKS_PER_REFRESH)

    const { data: territories } = useTerritoriesQuery(queryTick, currentScenario ?? "")
    const { data: factions } = useFactionsQuery(queryTick, currentScenario ?? "")

    const claimedTerritories = useMemo(
        () => (territories ?? []).filter(t => t.owner >= 0).length,
        [territories],
    )

    async function startRun() {
        await StartRun(100)
        setRunning(true)
        setStarted(true)
    }

    async function pauseRun() {
        await Pause()
        setRunning(false)
    }

    async function resumeRun() {
        await Resume()
        setRunning(true)
    }

    async function loadScenario(scenario: string) {
        await LoadScenario(scenario, 0)
        setCurrentScenario(scenario)
        setRunning(false)
        setStarted(false)
        setTick(0)
    }

    useMemo(() => {
        (async () => {
            let names = await ListScenarios()
            setScenarios(names ?? [])
        })();
    }, [])

    useEffect(() => {
        if (!running) return

        return EventsOn("tick", (t: number) => {
            setTick(t)
        })
    }, [running])

    useEffect(() => {
        function onKeyDown(e: KeyboardEvent) {
            if (e.key === "`") setShowTickDebug(v => !v)
        }
        window.addEventListener("keydown", onKeyDown)
        return () => window.removeEventListener("keydown", onKeyDown)
    }, [])

    const status = !currentScenario ? "idle" : running ? "running" : started ? "paused" : "idle"
    const statusStyles: Record<string, string> = {
        idle: "bg-space-600 text-ink-500",
        running: "bg-signal-up-dim text-signal-up",
        paused: "bg-signal-warn/15 text-signal-warn",
    }
    const statusDot: Record<string, string> = {
        idle: "bg-ink-600",
        running: "bg-signal-up animate-pulse",
        paused: "bg-signal-warn",
    }

    return (
        <div id="app" className="flex h-screen flex-col overflow-hidden bg-starfield text-ink-200">
            <header className="flex items-center justify-between border-b border-space-600 bg-space-900/80 px-6 py-3 backdrop-blur">
                <div className="flex items-center gap-3">
                    <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-violet-core/20 shadow-glow-sm">
                        <span className="text-sm font-black text-violet-glow">VE</span>
                    </div>
                    <div>
                        <h1 className="text-sm font-black uppercase tracking-[0.2em] text-ink-100">
                            Void Exchange
                        </h1>
                        <p className="text-[0.65rem] uppercase tracking-widest text-ink-500">
                            Territory Control Terminal
                        </p>
                    </div>
                </div>

                <span className={`inline-flex items-center gap-2 rounded-full px-3 py-1 text-[0.7rem] font-bold uppercase tracking-widest ${statusStyles[status]}`}>
                    <span className={`h-1.5 w-1.5 rounded-full ${statusDot[status]}`} />
                    {status}
                </span>
            </header>

            <div className="flex flex-1 overflow-hidden">
                <aside className="flex w-64 flex-shrink-0 flex-col gap-4 overflow-y-auto border-r border-space-600 bg-space-900/60 p-4">
                    <div>
                        <h2 className="mb-2 text-[0.7rem] font-bold uppercase tracking-widest text-ink-500">
                            Scenarios
                        </h2>
                        <div className="flex flex-col gap-2">
                            {scenarios?.map(s => {
                                const active = s === currentScenario
                                return (
                                    <button
                                        key={s}
                                        onClick={() => loadScenario(s)}
                                        className={`rounded-lg border px-3 py-2 text-left text-sm font-semibold transition-colors ${active
                                                ? "border-violet-core bg-violet-dim/50 text-violet-glow shadow-glow-sm"
                                                : "border-space-600 bg-space-800/60 text-ink-300 hover:border-space-500 hover:text-ink-100"
                                            }`}
                                    >
                                        {s}
                                    </button>
                                )
                            })}
                            {scenarios?.length === 0 && (
                                <p className="text-xs text-ink-500">No scenarios found.</p>
                            )}
                        </div>
                    </div>

                    <div className="mt-auto rounded-lg border border-space-600 bg-space-800/60 p-3">
                        <p className="text-[0.65rem] font-bold uppercase tracking-widest text-ink-500">
                            Active Scenario
                        </p>
                        <p className="mt-1 text-sm font-semibold text-ink-200">
                            {currentScenario ?? "None selected"}
                        </p>
                    </div>
                </aside>

                <main className="flex flex-1 flex-col overflow-hidden p-6">
                    <div className="mb-4 flex flex-wrap items-center justify-between gap-4">
                        <div className="flex flex-wrap gap-3">
                            <StatTile label="Mission Time" value={formatMissionClock(tick)} accent="violet" />
                            <StatTile label="Factions" value={factions?.length ?? 0} accent="neutral" />
                            <StatTile
                                label="Territories Claimed"
                                value={`${claimedTerritories}/${territories?.length ?? 0}`}
                                accent="up"
                            />
                        </div>

                        <div className="flex gap-2">
                            <button
                                disabled={!currentScenario || started}
                                onClick={() => startRun()}
                                className="rounded-lg bg-violet-core px-4 py-2 text-sm font-bold uppercase tracking-wide text-white shadow-glow-sm transition-opacity enabled:hover:bg-violet-bright disabled:cursor-not-allowed disabled:opacity-30"
                            >
                                Start
                            </button>
                            <button
                                disabled={!started || running}
                                onClick={() => resumeRun()}
                                className="rounded-lg border border-space-600 bg-space-800 px-4 py-2 text-sm font-bold uppercase tracking-wide text-ink-200 transition-colors enabled:hover:border-signal-up enabled:hover:text-signal-up disabled:cursor-not-allowed disabled:opacity-30"
                            >
                                Resume
                            </button>
                            <button
                                disabled={!started || !running}
                                onClick={() => pauseRun()}
                                className="rounded-lg border border-space-600 bg-space-800 px-4 py-2 text-sm font-bold uppercase tracking-wide text-ink-200 transition-colors enabled:hover:border-signal-warn enabled:hover:text-signal-warn disabled:cursor-not-allowed disabled:opacity-30"
                            >
                                Pause
                            </button>
                        </div>
                    </div>

                    <div className="mb-4 flex gap-1 rounded-lg border border-space-600 bg-space-850/60 p-1 self-start">
                        {(["territories", "factions", "drones"] as const).map(v => (
                            <button
                                key={v}
                                onClick={() => setView(v)}
                                className={`rounded-md px-4 py-1.5 text-xs font-bold uppercase tracking-widest transition-colors ${view === v
                                        ? "bg-violet-core text-white shadow-glow-sm"
                                        : "text-ink-500 hover:text-ink-200"
                                    }`}
                            >
                                {v === "territories" ? "Territories" : v === "factions" ? "Factions" : "Drones"}
                            </button>
                        ))}
                    </div>

                    <div className="flex flex-1 flex-col overflow-hidden">
                        {view === "territories" && <TerritoryTable tick={queryTick} scenario={currentScenario ?? ""} />}
                        {view === "factions" && <FactionTable tick={queryTick} scenario={currentScenario ?? ""} />}
                        {view === "drones" && <DroneTable tick={queryTick} scenario={currentScenario ?? ""} />}
                    </div>
                </main>
            </div>

            {showTickDebug && (
                <div className="fixed bottom-3 right-3 rounded-md border border-space-600 bg-space-900/90 px-2 py-1 font-mono text-[0.65rem] text-ink-500 shadow-glow-sm backdrop-blur">
                    tick {tick}
                </div>
            )}
        </div>
    )
}

export default App
