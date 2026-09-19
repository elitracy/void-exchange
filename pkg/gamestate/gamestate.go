package gamestate

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"

	"github.com/elitracy/void-exchange/pkg/entity"
)

const MiningRatePerDrone = 5

type GameState struct {
	EM          *entity.EntityManager
	Territories []entity.EntityId
	Factions    []entity.EntityId
	currentTick int
	Rng         *rand.Rand

	paused   bool
	pausedMu sync.Mutex
}

func (gs *GameState) Territory(id entity.EntityId) (*entity.Territory, error) {
	t, err := entity.GetAs[*entity.Territory](gs.EM, id)
	return t, err
}

func (gs *GameState) Faction(id entity.EntityId) (*entity.Faction, error) {
	f, err := entity.GetAs[*entity.Faction](gs.EM, id)
	return f, err
}

func (gs *GameState) Deposit(id entity.EntityId) (*entity.ResourceDeposit, error) {
	d, err := entity.GetAs[*entity.ResourceDeposit](gs.EM, id)
	return d, err
}

func (gs *GameState) Drone(id entity.EntityId) (*entity.Drone, error) {
	d, err := entity.GetAs[*entity.Drone](gs.EM, id)
	return d, err
}

func (gs *GameState) IsPaused() bool {
	gs.pausedMu.Lock()
	defer gs.pausedMu.Unlock()
	return gs.paused
}

func (gs *GameState) Pause() {
	gs.pausedMu.Lock()
	defer gs.pausedMu.Unlock()
	gs.paused = true
}

func (gs *GameState) Resume() {
	gs.pausedMu.Lock()
	defer gs.pausedMu.Unlock()
	gs.paused = false
}

func NewGameState(seed int64) *GameState {
	return &GameState{
		EM:  entity.NewEntityManager(),
		Rng: rand.New(rand.NewSource(seed)),
	}
}

func PopulateTerritory(
	em *entity.EntityManager,
	t *entity.Territory,
	resourceTypes []entity.ResourceType,
	min, max int64,
	rng *rand.Rand,
) error {
	if min > max || min < 0 {
		return fmt.Errorf("invalid min,max (%d,%d)", min, max)
	}

	for _, rt := range resourceTypes {
		amount := rng.Int63n(max-min+1) + min
		deposit := entity.NewResourceDeposit(rt, int(amount))
		entity.Register(em, deposit)

		if err := t.AddDeposit(deposit.Id()); err != nil {
			return err
		}
	}

	return nil
}

func (gs *GameState) Tick() error {
	gs.currentTick += 1
	err := gs.EM.Tick()
	conflictErr := gs.ResolveTerritoryConflicts()
	miningErr := gs.ResolveMining()
	return errors.Join(err, conflictErr, miningErr)
}

func (gs *GameState) CurrentTick() int { return gs.currentTick }
func (gs *GameState) SetSeed(seed int) { gs.Rng = rand.New(rand.NewSource(int64(seed))) }

func (gs *GameState) DispatchDrones(factionId, territoryId entity.EntityId, droneIds []entity.EntityId, activity entity.DroneActivity) error {
	if len(droneIds) == 0 {
		return ErrNoDronesSpecified
	}

	faction, err := gs.Faction(factionId)
	if err != nil {
		return fmt.Errorf("faction (%d): %w", factionId, err)
	}

	terr, err := gs.Territory(territoryId)
	if err != nil {
		return fmt.Errorf("territory (%d): %w", territoryId, err)
	}

	var wantType entity.DroneType
	switch activity {
	case entity.DroneFighting:
		wantType = entity.DroneFighter
	case entity.DroneMining:
		if terr.Owner != factionId {
			return fmt.Errorf("territory (%d): %w", territoryId, ErrTerritoryNotOwned)
		}
		wantType = entity.DroneMiner
	default:
		return fmt.Errorf("%q: %w", activity, ErrInvalidActivity)
	}

	drones := make([]*entity.Drone, 0, len(droneIds))
	for _, droneId := range droneIds {
		if !faction.HasDrone(droneId) {
			return fmt.Errorf("drone (%d): %w", droneId, ErrDroneNotOwned)
		}

		drone, err := gs.Drone(droneId)
		if err != nil {
			return fmt.Errorf("drone (%d): %w", droneId, err)
		}

		if drone.Activity != entity.DroneIdle {
			return fmt.Errorf("drone (%d): %w", droneId, ErrDroneNotIdle)
		}

		if drone.Type != wantType {
			return fmt.Errorf("drone (%d): %w", droneId, ErrWrongDroneType)
		}

		drones = append(drones, drone)
	}

	for _, drone := range drones {
		drone.Activity = activity
		drone.Target = territoryId
	}

	if activity == entity.DroneFighting {
		present := false
		for _, fid := range terr.Factions {
			if fid == factionId {
				present = true
				break
			}
		}
		if !present {
			terr.AddFaction(factionId)
		}
	}

	return nil
}

func (gs *GameState) RecallDrones(factionId entity.EntityId, droneIds []entity.EntityId) error {
	if len(droneIds) == 0 {
		return ErrNoDronesSpecified
	}

	faction, err := gs.Faction(factionId)
	if err != nil {
		return fmt.Errorf("faction (%d): %w", factionId, err)
	}

	drones := make([]*entity.Drone, 0, len(droneIds))
	for _, droneId := range droneIds {
		if !faction.HasDrone(droneId) {
			return fmt.Errorf("drone (%d): %w", droneId, ErrDroneNotOwned)
		}

		drone, err := gs.Drone(droneId)
		if err != nil {
			return fmt.Errorf("drone (%d): %w", droneId, err)
		}

		if drone.Activity == entity.DroneIdle {
			return fmt.Errorf("drone (%d): %w", droneId, ErrDroneNotCommitted)
		}

		drones = append(drones, drone)
	}

	touchedTerritories := map[entity.EntityId]struct{}{}
	for _, drone := range drones {
		if drone.Activity == entity.DroneFighting {
			touchedTerritories[drone.Target] = struct{}{}
		}
		drone.Activity = entity.DroneIdle
		drone.Target = -1
	}

	for territoryId := range touchedTerritories {
		gs.pruneFactionPresence(factionId, territoryId)
	}

	return nil
}

func (gs *GameState) pruneFactionPresence(factionId, territoryId entity.EntityId) {
	terr, err := gs.Territory(territoryId)
	if err != nil {
		return
	}

	faction, err := gs.Faction(factionId)
	if err != nil {
		return
	}

	for _, droneId := range faction.Fleet {
		drone, err := gs.Drone(droneId)
		if err != nil {
			continue
		}
		if drone.Activity == entity.DroneFighting && drone.Target == territoryId {
			return // still fighting here, stay present
		}
	}

	terr.RemoveFaction(factionId)
}

type combatant struct {
	factionId   entity.EntityId
	drones      []*entity.Drone
	totalAttack int
}

// Each faction's total attack is dealt as damage split evenly across the
// opposing side's drones.
func (gs *GameState) ResolveTerritoryConflicts() error {
	errs := []error{}

	for _, id := range gs.Territories {
		terr, err := gs.Territory(id)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		if len(terr.Factions) == 0 {
			continue
		}

		if len(terr.Factions) == 1 {
			terr.Owner = terr.Factions[0]
			continue
		}

		combatants := make([]*combatant, 0, len(terr.Factions))
		for _, factionId := range terr.Factions {
			faction, err := gs.Faction(factionId)
			if err != nil {
				errs = append(errs, fmt.Errorf("invalid faction (%d): %w", factionId, err))
				continue
			}

			c := &combatant{factionId: factionId}
			for _, droneId := range faction.Fleet {
				drone, err := gs.Drone(droneId)
				if err != nil {
					errs = append(errs, fmt.Errorf("invalid drone (%d): %w", droneId, err))
					continue
				}

				if drone.Activity != entity.DroneFighting || drone.Target != id {
					continue
				}

				c.drones = append(c.drones, drone)
				c.totalAttack += drone.Attack
			}

			combatants = append(combatants, c)
		}

		totalAttackAll := 0
		for _, c := range combatants {
			totalAttackAll += c.totalAttack
		}

		for _, c := range combatants {
			if len(c.drones) == 0 {
				continue
			}

			incoming := totalAttackAll - c.totalAttack
			if incoming <= 0 {
				continue
			}

			perDrone := incoming / len(c.drones)
			remainder := incoming % len(c.drones)

			for i, drone := range c.drones {
				dmg := perDrone
				if i < remainder {
					dmg++ // spread the remainder so total damage dealt isn't lost to integer division
				}
				drone.HP -= dmg
			}
		}

		remainingFactions := []entity.EntityId{}
		for _, c := range combatants {
			faction, err := gs.Faction(c.factionId)
			if err != nil {
				continue
			}

			stillFighting := false
			for _, drone := range c.drones {
				if drone.HP <= 0 {
					if _, err := faction.RemoveDrone(drone.Id()); err != nil {
						errs = append(errs, err)
					}
					continue
				}
				stillFighting = true
			}

			if stillFighting {
				remainingFactions = append(remainingFactions, c.factionId)
			}
		}

		terr.Factions = remainingFactions

		if len(terr.Factions) == 1 {
			terr.Owner = terr.Factions[0]
		}
	}

	return errors.Join(errs...)
}

func (gs *GameState) ResolveMining() error {
	errs := []error{}

	for _, id := range gs.Territories {
		terr, err := gs.Territory(id)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		if terr.Owner == -1 {
			continue
		}

		faction, err := gs.Faction(terr.Owner)
		if err != nil {
			errs = append(errs, fmt.Errorf("invalid owner faction (%d): %w", terr.Owner, err))
			continue
		}

		minerCount := 0
		for _, droneId := range faction.Fleet {
			drone, err := gs.Drone(droneId)
			if err != nil {
				errs = append(errs, fmt.Errorf("invalid drone (%d): %w", droneId, err))
				continue
			}
			if drone.Activity == entity.DroneMining && drone.Target == id {
				minerCount++
			}
		}

		if minerCount == 0 {
			continue
		}

		for _, depositId := range terr.ResourceDeposits {
			deposit, err := gs.Deposit(depositId)
			if err != nil {
				errs = append(errs, fmt.Errorf("invalid deposit (%d): %w", depositId, err))
				continue
			}

			extracted := min(deposit.Remaining, minerCount*MiningRatePerDrone)
			if extracted <= 0 {
				continue
			}

			deposit.Remaining -= extracted
			faction.UpdateResource(deposit.Resource.Type, extracted)
		}
	}

	return errors.Join(errs...)
}
