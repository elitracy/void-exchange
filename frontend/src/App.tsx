import { useState, useMemo, useEffect } from 'react';
import './App.css';
import { CurrentTick, GetTerritories, ListScenarios, LoadScenario, Pause, Resume, StartRun } from "../wailsjs/go/main/App";

function App() {
    const [scenarios, setScenarios] = useState<string[]>()
    const [currentScenario, setCurrentScenario] = useState<string>()
    const [tick, setTick] = useState(0)
    const [running, setRunning] = useState(false)
    const [territories, setTerritories] = useState<string[]>()


    async function startRun() {
        StartRun(100)
        setRunning(true)
    }

    async function pauseRun() {
        await Pause()
    }

    async function resumeRun() {
        await Resume()
    }

    async function loadScenario(scenario: string) {
        await LoadScenario(scenario, 0)
        setCurrentScenario(scenario)
        setTick(0)
    }

    useEffect(() => {
        (async () => {
            let ts = (await GetTerritories()).map(t => String(t.id))
            setTerritories(ts)
        })()
    }, [currentScenario])

    useMemo(() => {
        (async () => {
            let names = await ListScenarios()
            setScenarios(names ?? [])
        })();
    }, [])

    useEffect(() => {
        if (!running) return

        const id = setInterval(async () => {
            setTick(await CurrentTick())
        }, 100)

        return () => clearInterval(id)
    }, [running])


    return (
        <div id="app">
            <div id="header">
                <h1>Sim</h1>
            </div>
            <div id="content">
                <div id="options-pane">
                    <h2>Options</h2>

                    <h3>Scenarios</h3>
                    <div id="scenario-list">
                        {scenarios?.map(s => {
                            return (
                                <button onClick={() => loadScenario(s)}>{s}</button>
                            )
                        })}
                    </div>
                    <h3>Current Scenario: {currentScenario ?? "No scenario selected"}</h3>


                </div>
                <div id="content-pane">
                    <div id="content-pane-header">
                        <p>Tick: {tick}</p>
                        <div id="tick-buttons">
                            <>
                                <button className={!currentScenario ? "button-inactive" : ""} disabled={!currentScenario} onClick={() => startRun()}>Start</button>
                                <button className={!running ? "button-inactive" : ""} disabled={!running} onClick={() => resumeRun()}>Resume</button>
                                <button className={!running ? "button-inactive" : ""} disabled={!running} onClick={() => pauseRun()}>Pause</button>
                            </>
                        </div>
                    </div>
                    <div>
                        {territories?.map(t => {
                            return (
                                <p>{t}</p>
                            )
                        })}
                    </div>
                </div>
            </div>
        </div >
    )
}

export default App
