package gamestate

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/elitracy/void-exchange/pkg/entity"
)

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
	gs.ResolveTerritoryConflicts()
	return err
}

func (gs *GameState) CurrentTick() int { return gs.currentTick }
func (gs *GameState) SetSeed(seed int) { gs.Rng = rand.New(rand.NewSource(int64(seed))) }

func (gs *GameState) ResolveTerritoryConflicts() error {
	errs := []error{}

	for _, id := range gs.Territories {
		terr, err := gs.Territory(id)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		if len(terr.Factions) == 1 {
			terr.Owner = terr.Factions[0]
			continue
		}

		if len(terr.Factions) <= 0 {
			continue
		}

		for _, faction_id := range terr.Factions {

			faction, err := entity.GetAs[*entity.Faction](gs.EM, faction_id)
			if err != nil {
				errs = append(errs, fmt.Errorf("invalid faction (%d): %w", faction_id, err))
				continue
			}

			var drones []*entity.Drone
			for _, drone_id := range faction.Fleet {
				drone, err := entity.GetAs[*entity.Drone](gs.EM, drone_id)
				if err != nil {
					errs = append(errs, fmt.Errorf("invalid drone (%d): %w", drone_id, err))
					continue
				}

				drones = append(drones, drone)

				if drone.HP > 0 {
					drone.HP -= 10
				}

			}

			survivors := []entity.EntityId{}
			for _, drone := range drones {
				if drone.HP > 0 {
					survivors = append(survivors, drone.Id())
				}
			}

			faction.Fleet = survivors
			if len(faction.Fleet) == 0 {
				_, err := terr.RemoveFaction(faction.Id())
				if err != nil {
					errs = append(errs, err)
				}
			}
		}

		if len(terr.Factions) == 1 {
			terr.Owner = terr.Factions[0]
		}
	}

	return nil
}
