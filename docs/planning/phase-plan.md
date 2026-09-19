# Void Exchange — Phase Plan & Agent Dispatch Backlog

Source: `docs/design/game-design.md` (currently uncommitted, in the main
checkout at `/Users/elitracy/Documents/Coding/Projects/void_exchange`), cross-
checked against the actual uncommitted code in that checkout (`git diff
--stat` there shows ~763 lines of Phase 1 work already written but not yet
committed: dispatch/recall, even-split combat, mining, starting resources,
and a Drones GUI tab).

This document exists so each numbered item below can be handed to an agent
as a self-contained task. Each step lists: what exists today, what to build,
the files it touches, and how to know it's done. Steps are grouped into
tracks; tracks are roughly priority-ordered but Track A should go first
regardless, since it's cleanup of work that already exists.

---

## Track A — Close out Phase 1 (verify + land the uncommitted work)

The design doc's `Status` paragraphs say most of Phase 1 "hasn't started."
That's stale — the code already exists uncommitted in the main checkout.
Before scoping new work, that work needs to be verified and landed, and the
doc corrected so it stops being misleading to the next agent that reads it.

### A1. Verify the uncommitted Phase 1 work builds and passes tests
- **What exists:** uncommitted changes to `pkg/entity/{fleet,faction}.go`,
  `pkg/gamestate/gamestate.go` (+errors.go), `pkg/api/{views,commands}.go`,
  `pkg/scenario/scenario.go`, `app.go`, and frontend files including a new
  `frontend/src/components/DroneTable.tsx`.
- **Do:** in the main checkout, run `go build ./...`, `go vet ./...`, and
  `go test ./...`; run the frontend build (`pnpm install && pnpm build` or
  equivalent under `frontend/`). Fix anything broken.
- **Done when:** clean build + all tests green on both sides.

### A2. Review and commit the Phase 1 diff
- **Do:** split the uncommitted diff into logical commits (suggested split:
  entity model changes → gamestate dispatch/combat/mining → API/views →
  frontend DroneTable + wiring → scenario starting-resources config). Note
  the diff also contains an unrelated import-path rename
  (`elitracy/space-war-sim` → `elitracy/void-exchange`) across nearly every
  file — call that out as its own commit rather than folding it silently
  into a feature commit.
- **Done when:** Phase 1 work is committed on `main` (or a feature branch,
  per the user's usual workflow) with reviewable, logically-scoped commits.

### A3. Update `docs/design/game-design.md` Status sections
- **What exists:** §2.1, §5 item list, and the Status paragraphs throughout
  the doc describe Phase 1 as not-yet-built.
- **Do:** rewrite each `Status:` paragraph in §2.1–§2.4 and §5 to reflect
  what A1/A2 actually land — e.g. §2.1's Status paragraph should say the
  even-split combat model and dispatch/recall are implemented, not "This is
  the gap Phase 1 work should close first."
- **Done when:** doc's Status text matches the committed code, and §5's
  seven-item checklist is checked off or annotated with what's actually
  still open (see Track B for the pieces that are genuinely still open).

---

## Track B — Phase 1 loose ends (flagged in the doc, not yet built)

These are gaps the design doc calls out explicitly but that Track A's diff
does **not** cover. Each is small and independent — good candidates for
parallel agent dispatch.

### B1. Territory `Name` field
- **Gap:** `entity.Territory` (`pkg/entity/territory.go`) has no `Name`
  field; the sketch's Territory cluster (§3.1) calls for one, and the GUI
  currently labels territories by raw id (`#001`) in `DroneTable.tsx`'s
  territory dropdown and presumably `TerritoryTable.tsx`.
- **Do:** add `Name string` to `Territory`, set it in `scenario.Build` from
  a new `TerritoryConfig.Name` field (`pkg/scenario/scenario.go`), thread it
  through `NewTerritoryView` (`pkg/api/views.go`), and use it in the
  frontend wherever territory ids are currently displayed raw.
- **Done when:** territories have human-readable names end-to-end from
  scenario JSON to GUI.

### B2. Generator explicit on/off `Status`
- **Gap:** `entity.Generator` (`pkg/entity/production.go`) only has
  `OutputTarget`; "off" is implied by nothing pointing at it, not a state
  the generator itself holds. The sketch (§3.2) wants an explicit
  `status (on/off)`.
- **Do:** decide whether this is purely cosmetic (derive `on` from
  `OutputTarget != -1` in the view layer, no entity change) or a real field
  a player can toggle independent of assignment (entity change). See
  **Question Q1** below — this is a design decision, not just an
  implementation detail, so don't start until Q1 is answered.

### B3. Currency-miner byproduct
- **Gap:** §2.2 says "Currency miners occasionally yield a by-product worth
  more than standard currency," not implemented. `entity.CreditMiner.Tick()`
  (`pkg/entity/production.go`) currently only accumulates `CurrentCredits`.
- **Blocked on:** **Question Q2** (what the byproduct actually is — a
  distinct resource type, a rare multiplier, tradeable item) needs an
  answer before this can be scoped as a coded task.

### B4. Resources GUI tab
- **Gap:** §4 sketches a Resources tab alongside Drones/Factions; only
  Territories/Factions/Drones exist in `frontend/src/App.tsx`'s `View`
  union today.
- **Do:** add a `ResourceTable.tsx` (mirror `DroneTable.tsx`'s pattern:
  `useResourcesQuery` hook, a `GetResources` App.go method, a
  `ResourceView`/`ResourceDepositView` in `pkg/api/views.go`), showing
  per-territory deposits (`resource`, `total`, `remaining`) and per-faction
  `OwnedResources` totals.
- **Done when:** a Resources tab exists with the same read-only polish as
  the Factions/Territories tables (no dispatch actions needed here — this
  is a viewer, not a Produce/Expand control surface; those come in Track C
  UI work).

### B5. Produce/Expand build UI
- **Gap:** `Faction.AddFactory/AddDrone/AddGenerator/AddCreditMiner` exist
  in `pkg/entity/faction.go` and are called from scenario setup, but there
  is no player-facing "build a factory/generator/miner on this territory"
  command or UI — everything currently gets created only at scenario load.
- **Do:** add `App.go` methods (`BuildFactory`, `BuildGenerator`,
  `BuildCreditMiner`, mirroring the shape of `DispatchDrones`) that spend
  from `Faction.OwnedResources` via a resource cost (reuse
  `ProductRecipes`-style costing, or define new costs for
  generators/miners if they don't have one yet — `ProductRecipes` currently
  only covers `product_drone`/`product_generator`), and a frontend control
  (could live in the new Resources tab or a dedicated Territory detail
  view) to trigger them.
- **Done when:** a player can spend resources to add a factory/generator/
  miner to a territory they own, from the GUI, without editing scenario
  JSON.

---

## Track C — Phase 2: Buy/Sell (Exchange)

Nothing in `pkg/entity` exists for this yet. The design doc's decisions log
already settled the ledger shape (§6: plain append-only `Transaction` log +
mutable balances, not hash-chained), so this track can start from a fairly
firm spec.

### C1. Entity model: `Exchange`, `Transaction`, `Item`
- **Do:** add to `pkg/entity` (new file, e.g. `exchange.go`):
  - `Item{Id, Name, Owner EntityId, Value int, Subitems []EntityId}`
  - `Transaction{Id, From, To EntityId, Items []EntityId, Value int, Tick int}`
    (per §6's decided shape — no `hash`/`prev_hash`/`credit_hashes`, that
    part of the original sketch was explicitly rejected)
  - `Exchange{Id, Name string, Factions []EntityId, Transactions []EntityId,
    Items []EntityId}` as the container/market entity
  - Follow the existing `CoreEntity`/`entity.Register`/`EntityManager`
    pattern used by `Territory`/`Faction`/`Drone`.
- **Done when:** entities compile, have constructors matching the
  `New*`-style used elsewhere, and have unit tests mirroring
  `territory_test.go`/`fleet_test.go`'s style.

### C2. `GameState` market operations
- **Do:** add to `pkg/gamestate/gamestate.go`: `ListItems`, `ListOffers` (or
  equivalent), and a `Trade(exchangeId, buyerId, sellerId, itemIds
  []EntityId, price int) error` that validates ownership/balance, moves
  `Item.Owner`, debits/credits both factions' balances, and appends a
  `Transaction` to the exchange's log. Validate everything before mutating
  anything, matching `DispatchDrones`' "validate the whole batch first"
  pattern.
- **Blocked on:** **Question Q3** (does a faction's spendable balance come
  from `OwnedResources`, i.e. can you sell/buy minerals directly, or is
  there a separate `Credit`/currency balance distinct from resources? The
  sketch's Credit Miner produces "currency," which reads as distinct from
  `ResourceType` — needs to be resolved before `Trade`'s balance-mutation
  logic can be written correctly).
- **Done when:** `Trade` and listing methods exist with tests covering the
  happy path, insufficient-balance, and not-owned-item rejection cases.

### C3. API + views
- **Do:** `pkg/api/views.go`: `ExchangeView`, `ItemView`, `TransactionView`.
  `app.go`: `GetExchange`, `ListItems`, `Trade` methods (mirror
  `GetDrones`/`DispatchDrones` shape).
- **Done when:** Wails bindings regenerate cleanly and the new methods are
  callable from the frontend.

### C4. Market GUI
- **Gap:** §2.3 wants "a simple table-style interface... closer to a crypto
  exchange than a traditional storefront."
- **Do:** new `ExchangeTable.tsx` (or `MarketView.tsx`) — an order-book /
  ticker-style table (item, price, seller, action) rather than a shopping-
  cart layout. Reuse the `useTable`/`tableFeatures` pattern from
  `DroneTable.tsx`.
- **Blocked on:** **Question Q4** (sell directly to players too, or
  faction-to-faction only? — this is the open question the sketch itself
  raised and the doc never answered). This changes whether the UI needs a
  "player" identity/session concept at all, so resolve before building C4.
- **Done when:** a market tab lists tradeable items and lets a faction buy/
  sell, wired to `Trade`.

---

## Track D — Phase 3: Expand

§2.4 status says the underlying primitives already exist and this "may end
up more like a UI/flow concern than new entity types" — confirmed true after
reading the code; `Faction.AddFactory/AddDrone/AddGenerator/AddCreditMiner`
are all there. Track D is really just: does B5 (Produce/Expand build UI)
plus a way to do it specifically on territory a faction already owns, cover
the full ask, or does "Expand" need something Produce doesn't (e.g. a
distinct territory-capacity limit, or multi-territory logistics)?

### D1. Scope "Expand" as a phase
- **Blocked on:** **Question Q5** — is Expand just Produce-on-owned-
  territory (B5 covers it, no extra work needed), or does the sketch imply
  a genuinely separate mechanic (territory capacity/slots, cross-territory
  logistics, something else)? Needs a decision before this track gets its
  own steps; right now there's a real chance Track D is zero-work once B5
  ships.

---

## Open Questions (need a user decision before the blocked steps can start)

| # | Question | Blocks | Context |
|---|----------|--------|---------|
| Q1 | Does `Generator.Status` need to be a real toggle-able field, or is on/off fully derivable from `OutputTarget != -1`? | B2 | Sketch lists status as a field; current code treats it as implicit. |
| Q2 | What is the currency-miner "byproduct" mechanically — a rare second resource type, a value multiplier on existing output, or a tradeable `Item`? How rare? | B3 | §2.2 names the feature but gives no shape. |
| Q3 | Is there a distinct currency/credit balance per faction (separate from `OwnedResources[ResourceType]`), or does Buy/Sell trade resources and items directly against each other? | C2 | Sketch's `CreditMiner` produces "currency," which reads as separate from mineral-type resources, but nothing in `pkg/entity` currently models a currency balance. |
| Q4 | Buy/Sell: faction-to-faction only, or can a human player transact directly too? | C4 | Explicit open question left in the original sketch (§2.3), never answered. |
| Q5 | Is "Expand" (§2.4) fully covered by Produce-on-owned-territory (Track B5), or does it need its own mechanic (territory capacity, logistics)? | D1 | Doc's own Status note says this might collapse into a UI concern, but doesn't commit to that. |
| Q6 | Should `entity.Conflict` be built now, or does querying `Drone{Activity: fighting, Target: territoryId}` (today's approach) stay sufficient? | — (revisit trigger) | §6 decided against it "for now," conditional on GUI needing an in-progress-fight display or persisted combat log. Worth a yes/no once the Drones tab (already built) has been used for a while — does it already answer "what's fighting where," or is a log still missing? |

---

## Needed Features — flat checklist (cross-referenced to steps above)

- [ ] Verified, committed Phase 1 backend (dispatch/recall/combat/mining) — A1, A2
- [ ] Design doc Status sections accurate — A3
- [ ] Territory names — B1
- [ ] Generator on/off status (pending Q1) — B2
- [ ] Currency-miner byproduct (pending Q2) — B3
- [ ] Resources GUI tab — B4
- [ ] Produce/Expand build UI (factory/generator/miner from the frontend) — B5
- [ ] Exchange/Transaction/Item entities — C1
- [ ] Trade/market GameState logic (pending Q3) — C2
- [ ] Market API + views — C3
- [ ] Market GUI (pending Q4) — C4
- [ ] Expand-phase scoping decision (pending Q5) — D1
- [ ] Conflict entity revisit checkpoint (pending Q6) — tracked, not scheduled

## Suggested dispatch order

1. A1 → A2 → A3 (sequential, one agent, must land first)
2. B1, B4 in parallel (independent, no open questions blocking them)
3. Resolve Q1, Q2, Q3, Q4, Q5 with the user (batchable as one round of
   questions) — unblocks B2, B3, C2, C4, D1
4. B5 (depends on nothing but Track A)
5. C1 → C2 → C3 → C4 (sequential within the track; C1 could start
   immediately, C2 needs Q3 answered first)
6. D1 once Q5 is answered — likely folds into B5 rather than becoming new work
