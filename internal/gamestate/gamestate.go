package gamestate

import (
	"math/rand"

	"github.com/elitracy/space-war-sim/internal/entity"
)

type GameState struct {
	EM          *entity.EntityManager
	Territories []*entity.Territory
	Factions    []*entity.Faction
	CurrentTick int
	Rng         *rand.Rand
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
	types []entity.ResourceType,
	min, max int64,
	rng *rand.Rand,
) error {

	for _, rt := range types {
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
	gs.CurrentTick += 1
	return gs.EM.Tick()
}
