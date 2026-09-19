# Void Exchange — Game Design Document

Source sketch: [`../assets/vision-whiteboard.png`](../assets/vision-whiteboard.png)

This document transcribes the whiteboard vision into a written reference, maps it
onto the code that already exists, and scopes the work that is actually in
flight right now (**Phase 1: Claim/Fight**). Sections marked **Status** compare
the vision against `pkg/entity`, `pkg/gamestate`, `pkg/api`, and the frontend as
they stand today, so this document stays useful as an implementation checklist
rather than only a historical pitch.

## 1. Elevator pitch

An idle/incremental space strategy sim. Factions compete over territories for
their resource deposits, spend those resources to produce drones and
infrastructure, trade on an open market, and reinvest the proceeds into
expansion — closing the loop back into more territory claims.

## 2. Core game loop

The loop is a four-phase cycle:

```
Claim/Fight → Produce → Buy/Sell → Expand → (back to Claim/Fight)
```

### 2.1 Claim / Fight

- Dispatch fighting drones (or ships) to unclaimed territories.
- **Combat model (decided):** fighter drones will come in variations later;
  for now each drone has HP and attack stats derived from its level. Per
  tick, each side sums total attack from its committed drones; that total is
  split **evenly** across the opposing side's committed drones (simplest
  option — favors whichever side has more drones committed, not just more
  total HP). Free-for-all with 3+ factions: each faction's incoming damage is
  the sum of every other present faction's attack, split evenly across its
  own drones.
- Dispatch resolves immediately (no travel-time/en-route state for Phase 1 —
  a deliberate scope cut to keep this milestone simple; revisit if the game
  needs travel time later).
- Reinforcements can be dispatched to a territory a faction already owns, not
  just contested/unclaimed ones.
- A dispatch can be recalled before it resolves.
- Once a territory is claimed, its owning faction can start mining its
  resource deposits — but mining drones must be explicitly dispatched to
  that territory to do so; mining is not automatic just from ownership.

**Status:** partially built. `entity.Territory` already tracks `Owner` and the
contesting `Factions`, and `gamestate.ResolveTerritoryConflicts` runs every
tick to resolve contested territories. Today's resolution model is a
placeholder, not the sum-damage rule above — every contesting drone loses a
flat 10 HP per tick regardless of the size of the opposing fleet, until only
one faction has surviving drones. `entity.Drone` also has no level/attack
stat yet, only `HP`. There is no "dispatch drones to this territory" command
yet either; drones exist on a faction's `Fleet` but nothing assigns them to
contest, reinforce, or mine a specific territory, and there's no recall path.
This is the gap Phase 1 work should close first — see §5.

### 2.2 Produce

- Build drones, power generators, and currency miners.
- Currency miners occasionally yield a by-product worth more than standard
  currency.
- **Drone funding (decided):** for now, factions/players start with a fixed
  pool of resources to spend building drones — no separate "drone factory
  queue economy" yet, just spend-from-starting-stock.

**Status:** mostly built. `entity.Factory`, `entity.Generator`, and
`entity.CreditMiner` all exist and tick independently, gated on having a
power source assigned (`SetGenerator`/`SetOutputTarget`). `Factory` production
is recipe-driven via `ProductRecipes` (currently `product_drone` and
`product_generator`, costed in `resource_mineral`). Not yet built: the
occasional higher-value by-product from currency miners, and any UI to
build/manage these from the frontend.

### 2.3 Buy / Sell

- An open market where factions purchase goods from each other.
- Simple table-style interface.
- Should feel closer to a crypto exchange than a traditional storefront.
- Open question from the sketch: sell directly to players too?

**Status:** not started. This is the `Exchange` / `Transaction` / `Item` /
`Credit` model in §3.3 below — none of those types exist in `pkg/entity` yet.
This phase depends on Produce being in place (something to sell) and is
explicitly out of scope for the current milestone.

### 2.4 Expand

- Build more drones/ships.
- Add more factories to planets/territories.
- Add miners to territory a faction already owns.

**Status:** not started as a distinct phase, though the primitives it needs
(`Faction.AddFactory`, `AddDrone`, `AddGenerator`, `AddCreditMiner`) already
exist — expand is largely "call the Produce primitives against territory you
already hold," so it may end up more like a UI/flow concern than new entity
types.

## 3. Entity model (from the sketch)

The whiteboard sketches out four clusters of data. Each subsection lists the
sketch's fields, then how (or whether) that maps onto `pkg/entity` today.

### 3.1 Claim/Fight cluster

**Territory** — `id`, `name`, `resource_deposits[]`, `owner`, `conflicts`

- Matches `entity.Territory` (`ID` via `CoreEntity`, `ResourceDeposits`,
  `Owner`, `Factions`) with two gaps: no `Name` field, and no first-class
  `Conflicts` list — conflict state today is implicit in `Territory.Factions`
  having more than one entry, not an explicit record of a conflict.

**Faction** — `id`, `name`, `territories[]`, `fleet`, `factories`,
`generators`, `currency_miners`, `resources[]`

- Matches `entity.Faction` almost exactly (`Name`, `Territories`, `Fleet`,
  `Factories`, `Generators`, `CreditMiners`, `OwnedResources`). This one's
  essentially built as sketched.

**Conflict** — `id`, `factions[]`, `fleets[]`

- Not implemented as its own entity. `ResolveTerritoryConflicts` resolves
  conflicts procedurally against `Territory.Factions` instead of creating a
  persistent `Conflict` record. Worth deciding (see open questions) whether
  Phase 1 needs a real `Conflict` entity — e.g. for GUI display of an
  in-progress fight, or a combat log — or whether the implicit model is
  sufficient for now.

### 3.2 Produce cluster

**Resource Deposit** — `id`, `resource`, `total`, `remaining`
**Resource** — `type` (e.g. `mineral`)
**Generators (solar)** — `status (on/off)`
**Factory** — `id`, `product`, `rate`, `generator`
**Currency Miner** — `id`, `generator`, `rate`

- `entity.ResourceDeposit`, `entity.Resource`, `entity.Factory`,
  `entity.Generator`, and `entity.CreditMiner` all exist and align closely.
  `Factory` is actually richer than the sketch (tracks `CurrentResources` and
  `CurrentUnits` toward its recipe). `Generator` is simpler than the sketch —
  it has an `OutputTarget` but no explicit on/off `status`; today "off" is
  implied by nothing pointing a `PowerSource` at it, not a state on the
  generator itself.

### 3.3 Buy/Sell cluster

**Exchange** — `id`, `name`, `factions[]`, `transactions[]`, `items[]`
**Transaction** — `id`, `hash`, `purchaser`, `purchasee`, `items[]`, `value`,
`credit_hashes[]`
**Item** — `id`, `name`, `owner`, `value`, `subitems[]`
**Credit** — `id`, `hash`, `prev_hash`, `owner`, `value`

- None of this exists yet. Note the sketch's `Transaction`/`Credit` shape
  (hash + prev_hash) implies a hash-chained ledger — worth flagging explicitly
  to the user (see open questions) since that's a meaningfully bigger design
  commitment than a plain balance field, and it's not needed for Phase 1.

### 3.4 Fleet

**Fleet** — `miners`, `transport`, `fighters`

- `entity.Drone` exists with a `Type` enum (`drone_miner`, `drone_transport`,
  `drone_fighter`) plus `Name` and `HP`, referenced by id from
  `Faction.Fleet`. This matches the sketch's intent (three drone roles) even
  though the sketch groups them as three separate lists rather than one typed
  list.

## 4. GUI vision (from the sketch)

The right half of the whiteboard sketches the dashboard layout:

- A left-hand **Territories** list/rail showing each territory's id and
  current owner id.
- A main pane with tabs for **Drones**, **Factions**, **Resources**, each
  showing a table. The Drones tab's columns are sketched explicitly: **Name,
  Health, Type, Activity**.

**Status:** the Territories and Factions tables exist today
(`frontend/src/components/TerritoryTable.tsx`, `FactionTable.tsx`) using the
"Void Exchange" dark-purple palette (see the `gui-palette-void-exchange`
project memory). A dedicated **Drones** table matching the sketch's
Name/Health/Type/Activity columns does not exist yet, nor does a
**Resources** tab. `Activity` isn't a field the backend tracks on `Drone` at
all right now (idle vs. traveling vs. fighting vs. mining) — that's new state
to add if the Drones tab is built next, and it's also exactly the kind of
state Claim/Fight dispatch needs (see §5).

## 5. Phase 1 scope: Claim/Fight

This is the milestone actually being built right now. Based on decisions
above, Phase 1 needs:

1. **Drone level/attack stat** — add a `Level` (or direct `Attack`) field to
   `entity.Drone` alongside existing `HP`, as the basis for combat damage.
   Variation by drone type comes later; level is enough for now.
2. **A dispatch command** — send N drones from a faction's fleet to a
   specific territory, covering all three cases: contest (unclaimed/enemy),
   reinforce (already owned), and mine (owned territory, miner-type drones).
   Needs a **recall** path before the dispatch resolves.
3. **Drone `Activity`/status field** — idle / fighting / mining, plus a
   `Target` territory id, both to know which drones are actually committed
   on a given tick and to back the sketched Drones GUI table. No separate
   "recalled" state needed since dispatch/recall are instant — recall just
   resets a drone straight back to idle.
4. **Even-split conflict resolution** — replace the flat HP-drain placeholder
   in `ResolveTerritoryConflicts` with the even-split sum-damage rule from
   §2.1, scoped to only the drones actually dispatched to fight at that
   territory (not a faction's whole fleet, which was the old code's bug —
   see implementation notes below).
5. **Territory → mining** — once a mining drone is dispatched to an owned
   territory, it should consume `ResourceDeposit.Remaining` there over time
   and feed the result into the faction's `OwnedResources`.
6. **Starting resource pool** — factions begin with a fixed resource stock
   spendable on drones via the existing `Factory`/`ProductRecipes` path.
7. **GUI**: a Drones table (Name/Health/Type/Activity as sketched) and a way
   to issue dispatch/recall commands from the frontend.

## 6. Decisions log (formerly open questions)

- **Damage split — decided:** even split, both for the 2-faction case and the
  N-faction free-for-all case (see §2.1). No focus-fire, no pooled HP model.
  Simplest option, revisit only if playtesting shows it feels wrong.
- **Conflict as an entity — decided against, for now:** no persistent
  `Conflict` entity. "Which drones are currently committed to what" is
  tracked directly on `entity.Drone` via its `Activity`/`Target` fields
  instead (see §5) — querying "drones with `Target == territoryId`" gives the
  same GUI-displayable in-progress-fight info without a new entity type.
  Revisit if a persisted combat log/history becomes a real requirement.
- **Ledger design for Buy/Sell — approach agreed, implementation deferred:**
  confirmed a plain **append-only transaction log** (`Transaction{id, from,
  to, item(s), value, tick}`) plus ordinary mutable balance fields, not a
  hash-chained ledger — chosen specifically so a future hacking/stealing
  mechanic can just make an illegitimate balance/log write rather than
  needing to break cryptographic chain integrity. Not implemented yet; this
  is Buy/Sell-phase work, out of scope for Phase 1.
