# Testing gaps: audit, fixes, and new coverage

This documents the state of testing across the app (frontend UI, the Go
`App`/Wails glue layer, and the Go domain packages), the concrete bugs found
while adding tests, which ones were fixed, and which are documented but left
alone (with a recommended fix) because they're bigger than "add a test."

Layering used throughout:

- **UI (frontend, `frontend/src`)** — needs *unit* tests. Covered with Vitest.
- **"In between" (`app.go` / `main.go`, the Wails-bound glue between the UI
  and the simulation)** — needs *integration* tests, since its only job is
  wiring other already-unit-tested packages together correctly. Covered with
  Go tests in `app_test.go` that exercise the real `pkg/api` /
  `pkg/gamestate` / `pkg/engine` stack (no mocks).
- **Domain logic (`pkg/entity`, `pkg/gamestate`, `pkg/engine`, `pkg/scenario`,
  `pkg/api`)** — already had unit tests before this change; gaps found there
  are noted below and a few easy ones were filled in.

## How to run everything

```sh
go test ./...                 # Go unit + integration tests
go test ./... -race           # same, with the race detector (recommended)
cd frontend && pnpm test      # frontend unit tests (Vitest)
```

---

## Bugs found and fixed

### 1. `App.Pause/Resume/Stop/IsPaused` crashed if called before a scenario loaded

**What was wrong:** `app.go` guarded every other method (`GetTerritories`,
`GetFactions`, `StartRun`, ...) with a `nil` check on `a.runner`/`a.gs`, but
`Pause`, `Resume`, `Stop`, and `IsPaused` dereferenced `a.runner`/`a.gs`
directly:

```go
func (a *App) Pause()  { a.runner.Pause() }        // nil pointer panic
func (a *App) IsPaused() bool { return a.gs.IsPaused() } // nil pointer panic
```

Since these are bound directly to the frontend via Wails, any call from the
UI before `LoadScenario` succeeds (a stray click, a race between button
disabled-state and an in-flight click, or simply a future UI bug) would
panic and crash the whole desktop process, not just show an error.

**Fix:** added the same `nil` guard used elsewhere; `Pause`/`Resume`/`Stop`
become no-ops and `IsPaused` returns `false` if nothing is loaded yet.

**Test:** `TestApp_ControlsBeforeScenarioLoaded_DoNotPanic` in `app_test.go`.

### 2. `pkg/logging`'s global logger raced itself and could panic

**What was wrong:** `pkg/logging` keeps one package-level `*logger`
singleton (`_logger`). `engine.RunGame` calls `logging.Init(...)` at the
start of *every single run* and (after this change) `logging.Flush()` when
it ends. Two problems compounded:

1. `Flush()` read/closed whatever `_logger` **currently** was, not
   necessarily the instance its own `Init()` call had created. If a second
   run started (a new `Init()`) before the first run's deferred `Flush()`
   executed, `Flush()` would close the *second* run's channel while that
   run's own goroutine was still sending to it — a data race, and, once hit,
   a `panic: send on closed channel` (this is what `go test -race ./...`
   caught in `pkg/api`, where `TestStartStop` and `TestPauseResume` run
   back-to-back and share this same process-wide singleton).
2. `Runner.Start()` logged `"Started run"` **before** the goroutine that
   calls `engine.RunGame` (and therefore `logging.Init`) had necessarily
   run yet — so that message was always at risk of being written through
   whatever the *previous* run's (possibly already-closed) logger was,
   rather than the new one.

**Fix (`pkg/logging/logging.go`):**
- Added a `globalMu sync.Mutex` guarding reads/writes of `_logger`, and a
  per-logger `sendMu`/`closed` flag so `Flush()` marks the logger closed and
  closes its channel under the same lock that `log()` uses to send — so a
  send can never race a close.
- `Flush()` now also clears the global pointer, so any log call that lands
  after a flush (e.g. from a slow, still-unwinding goroutine) is a silent
  no-op instead of touching a dead channel.
- Moved the `"Started run"` log line from `Runner.Start()` into
  `engine.RunGame`, right after `logging.Init(...)`, so it's always written
  through the logger that call just created.

**Known remaining limitation (not fixed, documented only):** `log()` sends
on `l.queue` *while holding* `l.sendMu`, and the queue is a small buffered
channel (capacity 10). If producers ever outpace the single consumer
goroutine for more than 10 messages, `log()` blocks while holding `sendMu`,
and a concurrent `Flush()` would then block forever waiting for the same
lock — a deadlock. This is unlikely at today's log volume (a handful of
lifecycle messages per run) but would resurface if per-tick logging is ever
added. **Suggested fix:** send with `select { case l.queue <- msg: default:
<drop or grow> }`, or restructure shutdown to signal via a separate `done`
channel instead of closing the message channel.

**Tests:** `pkg/logging/logging_test.go` (new file) — `Init`/`Flush`
lifecycle, no-panic-before-`Init`, idempotent `Flush`, logging after
`Flush` is dropped not panicked, re-`Init` starts a clean file, and a
concurrent `Init`/`Flush`-vs-log-calls test intended to be run with `-race`.

### 3. `GameState.currentTick` was an unsynchronized `int` read from two goroutines

**What was wrong:** `pkg/gamestate/gamestate.go` protected `paused` with a
mutex but not `currentTick`:

```go
currentTick int
...
func (gs *GameState) Tick() error { gs.currentTick += 1; ... }
func (gs *GameState) CurrentTick() int { return gs.currentTick }
```

`Tick()` runs on the engine's background goroutine; `CurrentTick()` is
called from `App.CurrentTick()`, which the UI polls once a second on a
different goroutine. That's a textbook data race under the Go memory model
— confirmed by `go test -race .` failing as soon as an integration test
polled `CurrentTick()` while a run was ticking.

**Fix:** changed `currentTick` to `atomic.Int64` (`Add(1)` / `Load()`).
Minimal, no behavior change, and it's the only field that needed it —
`paused` already had its own mutex.

**Test:** `TestApp_StartPauseResumeStop_Lifecycle` in `app_test.go` polls
`CurrentTick()` concurrently with the engine ticking; run with `go test
-race` to verify.

### 4. `App` was untestable outside a real Wails runtime (silent `os.Exit(1)`)

**What was wrong:** `App.StartRun`'s completion goroutine calls
`runtime.EventsEmit(a.ctx, "run:finished", ...)` (from
`wailsapp/wails/v2/pkg/runtime`). That package's internals do this when the
context isn't the one `wails.Run` builds:

```go
func getEvents(ctx context.Context) frontend.Events {
    if ctx == nil { ... }
    result := ctx.Value("events")
    if result != nil { return result.(frontend.Events) }
    log.Fatalf("cannot call '%s': %s", funcName, contextError) // os.Exit(1)
}
```

`log.Fatalf` calls `os.Exit(1)` — it is **not** a panic and cannot be
recovered. Worse, `pkg/logging` calls `log.SetOutput(...)` to redirect the
*standard library's* global `log` package at a rotating file, so this fatal
message doesn't even reach the terminal — the whole test process (or, in
principle, the whole desktop app if this path were ever hit with a bad
context) just vanishes with exit code 1 and no visible diagnostic. This is
exactly why an early version of `app_test.go`'s lifecycle test would kill
the entire `go test` binary with no stack trace: `context.Background()` (a
plain test context) has no `"events"` value.

In production this is low-risk — `a.ctx` is always the real context Wails
passes to `startup()` — but it made `App`'s public API fundamentally
impossible to integration-test as-is.

**Fix:** added a small dependency-injection seam:

```go
type App struct {
    ...
    emitEvent func(ctx context.Context, eventName string, optionalData ...interface{})
}
func NewApp() *App { return &App{emitEvent: runtime.EventsEmit} }
```

Production code is unchanged (`NewApp()` still wires up the real
`runtime.EventsEmit`); tests construct `App` via a `newTestApp` helper that
swaps in a no-op. This is what unblocked writing `app_test.go` at all.

### 5. Root Go package couldn't be built or tested at all

**What was wrong:** `main.go` has `//go:embed all:frontend/dist`, and
`frontend/dist` is generated by `vite build` / `wails build` and is
git-ignored. On a clean checkout (or CI, before the frontend build step),
`go build .` / `go test .` fails outright:

```
main.go:11:12: pattern all:frontend/dist: no matching files found
```

Since `app.go` (the file this whole "in between" layer lives in) is
`package main`, this meant **the integration layer had zero tests and could
not even compile for testing**, independent of any assertions — that's a
harder blocker than "gaps in coverage."

**Fix:** added a tracked placeholder `frontend/dist/.gitkeep` (with a
`.gitignore` exception, `!frontend/dist/.gitkeep`) so the embed pattern
always matches at least one file. A real `pnpm build` / `wails build`
still populates `frontend/dist` with the actual app for packaging; this
only unblocks compiling the Go side without requiring a frontend build
first.

**Test:** `app_test.go` (new file) exercises `App` end-to-end against the
real `simple` scenario: `ListScenarios`, error paths before a scenario is
loaded, `LoadScenario` populating factions/territories, `GetTerritory` /
`GetDeposits` / `GetDeposit`, and the full start/pause/resume/stop
lifecycle (plus the double-`StartRun` "already running" error).

---

## Gaps documented, not fixed (bigger than a test)

These are real issues found while writing tests, left alone because fixing
them changes behavior/scope beyond "add tests and fix what's easy."

### Hardcoded, cwd-relative paths

`app.go`'s `LoadScenario` builds `fmt.Sprintf("./scenarios/%s.json", name)`,
and `ListScenarios("./scenarios")` is called the same way; `LoadScenario`
also hardcodes `api.NewRunner(gs, "./logs/run.log")`. These all resolve
relative to the process's **current working directory**, not the
executable's location. That's fine under `go run` / `wails dev` / `go
test` (cwd happens to be the repo root), but a packaged app (e.g. a macOS
`.app` bundle) can be launched with an unrelated cwd, silently breaking
`ListScenarios`/`LoadScenario`, or scattering `run.log` wherever the OS
happened to start the process (this is also why running the Go tests in
this repo now produces a `logs/run.log` in the repo root — harmless but
untidy).

**Suggested fix:** resolve these relative to `os.Executable()` (or a
configured data directory), not the cwd; or make the scenarios directory
and log path constructor arguments to `App` (they're already parameters to
the lower-level `api.ListScenarios`/`api.LoadScenario`/`api.NewRunner`
functions — `app.go` just doesn't forward that flexibility up).

### `App.Stop()` / `Runner.Stop()` don't block until the run has actually stopped

`Stop()` only cancels the run's `context.Context`; the engine goroutine
notices on its next loop iteration (up to one tick interval later) and
*then* runs its cleanup (`logging.Flush()`, etc.). There's no way for a
caller — the UI or a test — to know the run has fully unwound. In practice
this means, e.g., loading a new scenario immediately after `Stop()` starts
a second concurrent engine goroutine that briefly overlaps with the one
still finishing. This is what made `app_test.go`'s lifecycle test need a
manual `time.Sleep` at the end (documented inline in that file) instead of
a deterministic wait.

**Suggested fix:** have `Runner.Stop()` return (or `Start()`'s `done`
channel double as) a way to block until the goroutine has exited, e.g.
`func (r *Runner) Stop() { r.cancel(); <-r.doneInternal }`.

### Frontend: `useMemo` used for an async side effect

`App.tsx`'s scenario list is fetched like this:

```tsx
useMemo(() => {
    (async () => { let names = await ListScenarios(); setScenarios(names ?? []) })();
}, [])
```

`useMemo` is a *pure* memoization hint — React does not guarantee it won't
be called more than once, or that its result will be used, and using it to
kick off a side effect (a network/IPC call plus a `setState`) is
unsupported by React's own documentation even though it happens to work
today. **Fix:** use `useEffect` instead; behavior is identical for this
one-time-on-mount case, but it's the correct primitive.

### Frontend: no error handling on any Wails call

None of `App.tsx`'s `startRun`/`pauseRun`/`resumeRun`/`loadScenario`
handle a rejected promise. Since the Go side does return errors for
several of these (e.g. `StartRun` returns `"no scenario loaded"`, `"already
running"`), a real failure currently surfaces only as an unhandled promise
rejection in the browser console — the user sees no feedback and the UI
state (`running`, `currentScenario`) can end up out of sync with the
backend (e.g. `setRunning(true)` runs even though `StartRun` never actually
succeeded, since it's called before checking for an error in some paths).
**Fix:** wrap these calls in try/catch and surface failures in the UI
(toast, inline status, etc.).

### Frontend: derived state kept in `useState` + `useEffect` instead of computed directly

`FactionTable`/`TerritoryTable` copy query data into local state on every
change:

```tsx
const [factions, setFactions] = useState<Array<Faction>>([])
useEffect(() => { setFactions(transformFactions(data ?? [])) }, [scenario, tick, data])
```

`transformFactions`/`transformTerritories` are pure and cheap; this could
just be `const factions = useMemo(() => transformFactions(data ?? []),
[data])` (or even inline, uncached) with no behavior loss, one fewer
render per update, and no one-tick-stale window between `data` changing and
the effect running. Not fixed since it's a pure refactor with no bug behind
it today, but worth doing alongside any future touch of these files.

### Frontend: "paused" status label shown before anything ever started

`App.tsx`'s status derivation — `!currentScenario ? "idle" : running ?
"running" : "paused"` — labels the state right after `LoadScenario` (before
`StartRun` has ever been called) as "paused," which reads as "there's a run
you can resume" when there never was one. Minor UX/naming issue, not a
functional bug; a three-way `idle | running | paused` state doesn't
distinguish "never started" from "started then paused." Covered by
`App.test.tsx`'s `'loading a scenario enables Start and marks it active'`
test, which asserts the current (labeled) behavior rather than papering
over it.

### Coverage gaps in already-tested domain packages

Not new bugs, but worth tracking: after this change, `go test ./... -cover`
reports:

| package | coverage |
|---|---|
| `.` (App/main) | 81.4% |
| `pkg/api` | 91.8% |
| `pkg/engine` | 73.3% |
| `pkg/entity` | 89.1% |
| `pkg/gamestate` | 67.6% |
| `pkg/logging` | 100% |
| `pkg/scenario` | 100% |

The lowest, `pkg/gamestate` (67.6%), is mostly simple accessors
(`Faction`, `Deposit`, `SetSeed`) that are exercised indirectly through
other packages' tests but never directly in `gamestate`'s own test file,
plus a couple of branches in `ResolveTerritoryConflicts` (multi-faction
contested-territory combat resolution) that aren't hit by the existing
fixtures. `pkg/engine`'s uncovered 26.7% is almost entirely the `Tick()`
error-return path (`RunGame` returning early because `gs.Tick()` failed),
which none of the current fixtures trigger. These are reasonable next
targets but are pre-existing gaps in domain logic, not part of the UI/"in
between" scope this pass focused on.

`main.go`'s `main()` (Wails bootstrap: `wails.Run(...)`) and `App.Greet`
(a template leftover, unused by the UI) are untestable/uninteresting in the
normal sense — `main()` starts a real native window and event loop, which
isn't something a unit/integration test should do.

---

## New test inventory

**Go — integration ("in between"):**
- `app_test.go` (new) — `App` end-to-end against the real scenario/engine
  stack: scenario listing, pre-load error paths, nil-guard regression,
  load → factions/territories/deposits, and the full
  start/pause/resume/stop lifecycle including the double-start error.

**Go — unit (domain, gaps filled):**
- `pkg/logging/logging_test.go` (new) — see bug #2 above.
- `pkg/api/views_test.go` (new) — `NewTerritoryView`/`NewFactionView`/
  `NewDepositView` were previously untested (0% coverage): dedup/sort of
  deposit types, and the partial-view-plus-error behavior when a territory
  references a deposit id that no longer resolves.
- `pkg/api/runner_test.go` — added `TestRunner_PauseResume_DelegateToGameState`
  covering the `Runner.Pause`/`Resume` wrapper methods directly (previously
  only the underlying mock's `Pause`/`Resume` were called in tests, not the
  `Runner` methods themselves).

**Frontend — unit (UI, previously nonexistent):**
There was no test runner, config, or single test file in `frontend/`
before this change. Added Vitest + Testing Library
(`frontend/vitest.config.ts`, `frontend/src/test/setup.ts`,
`pnpm test` / `pnpm test:watch` scripts) and:
- `src/lib/theme.test.ts` — pure functions: faction accent
  determinism/wraparound, resource label formatting, chip class hashing.
- `src/components/ui/StatTile.test.tsx`, `ResourceChips.test.tsx` —
  presentational component behavior (empty states, accent classes, amount
  rendering).
- `src/hooks/useGameData.test.tsx` — the `enabled: !!scenario` gating and
  tick-based re-fetching, with the Wails bindings mocked.
- `src/components/FactionTable.test.tsx`, `TerritoryTable.test.tsx` — empty
  states, loading states, computed columns (stockpile sum, deposit types),
  and owner-name resolution (including the "Faction #N" fallback for an
  owner id with no matching faction), with `useGameData` mocked so these
  are true unit tests independent of Wails/react-query.
- `src/App.test.tsx` — scenario list rendering, load/start/pause button
  enablement and status label transitions, with all Wails bindings mocked.

No end-to-end or visual regression tests exist for the frontend (e.g. a
real Wails webview driven by a browser automation tool) — out of scope for
unit testing but worth naming as a gap if a release-gating check is ever
wanted.
