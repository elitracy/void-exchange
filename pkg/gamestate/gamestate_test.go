package gamestate_test

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"testing"

	"github.com/elitracy/void-exchange/pkg/entity"
	"github.com/elitracy/void-exchange/pkg/gamestate"
	"github.com/stretchr/testify/assert"
)

var (
	testErr = errors.New("big ol bad error")
)

type mockEntity struct {
	*entity.CoreEntity
	ticks int
}

type errEntity struct {
	*entity.CoreEntity
	err error
}

func newMockEntity() *mockEntity {
	return &mockEntity{
		CoreEntity: entity.NewCoreEntity(),
	}
}

func (e *mockEntity) Tick() error {
	e.ticks += 1
	return nil
}

func newErrEntity() *errEntity {
	return &errEntity{
		CoreEntity: entity.NewCoreEntity(),
		err:        testErr,
	}
}

func (e *errEntity) Tick() error {
	return e.err
}

func TestGameState_New(t *testing.T) {
	gs := gamestate.NewGameState(0)

	assert.NotNil(t, gs.EM)
	assert.NotNil(t, gs.Rng)
	assert.Equal(t, 0, gs.CurrentTick())
}

func TestGameState_Tick(t *testing.T) {
	tests := []struct {
		numTicks      int
		expectedTicks int
		numEntities   int
	}{
		{10, 10, 3},
		{10, 10, 0},
	}

	for _, tt := range tests {
		gs := gamestate.NewGameState(0)

		for range tt.numEntities {
			e := newMockEntity()
			entity.Register(gs.EM, e)
		}

		for range tt.numTicks {
			err := gs.Tick()
			assert.Nil(t, err)
		}

		assert.Equal(t, tt.expectedTicks, gs.CurrentTick())

		for _, id := range gs.EM.Entities {
			e, ok := gs.EM.GetEntity(id)
			assert.True(t, ok)

			mockE, ok := e.(*mockEntity)
			assert.True(t, ok)

			assert.Equal(t, tt.expectedTicks, mockE.ticks)
		}

	}
}

func TestGameState_TickErrorsPropogate(t *testing.T) {
	tests := []struct {
		expectedError error
		numEntities   int
	}{
		{testErr, 10},
	}

	for _, tt := range tests {
		gs := gamestate.NewGameState(0)

		for range tt.numEntities {
			e := newErrEntity()
			entity.Register(gs.EM, e)
		}

		err := gs.Tick()

		for _, id := range gs.EM.Entities {
			assert.ErrorContains(t, err, strconv.Itoa(int(id)))
			assert.ErrorContains(t, err, tt.expectedError.Error())
		}
	}

}

func TestPopulateTerritory(t *testing.T) {
	tests := []struct {
		resourceTypes []entity.ResourceType
		min, max      int64
		rng           *rand.Rand
	}{
		{
			[]entity.ResourceType{entity.ResourceEnergy},
			0, 10, rand.New(rand.NewSource(0)),
		},
		{
			[]entity.ResourceType{entity.ResourceMineral},
			-10, 10, rand.New(rand.NewSource(0)),
		},
		{
			[]entity.ResourceType{entity.ResourceMineral, entity.ResourceEnergy},
			0, 10, rand.New(rand.NewSource(0)),
		},
	}

	for _, tt := range tests {

		em := entity.NewEntityManager()
		terr := entity.NewTerritory()

		entity.Register(em, terr)

		err := gamestate.PopulateTerritory(em, terr, tt.resourceTypes, tt.min, tt.max, tt.rng)

		if tt.min > tt.max || tt.min < 0 {
			assert.ErrorContains(t, err, "invalid min,max")
		}

		for _, id := range terr.ResourceDeposits {
			e, ok := em.GetEntity(id)
			assert.True(t, ok)

			deposit, ok := e.(*entity.ResourceDeposit)
			assert.True(t, ok)

			assert.Contains(t, tt.resourceTypes, deposit.Resource.Type)

			assert.GreaterOrEqual(t, int(tt.max), deposit.Total)
			assert.LessOrEqual(t, int(tt.min), deposit.Total)
		}

		if tt.min <= tt.max && tt.min > 0 {
			assert.Equal(t, len(tt.resourceTypes), len(terr.ResourceDeposits))
		}

	}
}

type droneConf struct {
	factionIdx entity.EntityId
	hp, attack int
}

// TestResolveTerritoryConflicts exercises the even-split combat model:
// each tick, every present faction's total attack (summed across its
// drones actually dispatched to fight at this territory) is split evenly
// as damage across the opposing side's drones. Drones here are wired up
// directly (HP/Attack/Activity/Target set by hand) rather than through
// DispatchDrones, so each case can pin down exact combat stats regardless
// of the level-derived defaults NewDrone would otherwise apply.
func TestResolveTerritoryConflicts(t *testing.T) {
	tests := []struct {
		name            string
		drones          []droneConf
		ticks           int
		expectedOwnerId entity.EntityId
	}{
		{"uncontested single faction does not change owner", []droneConf{{0, 100, 10}}, 10, 0},
		{"2 factions - even split, sturdier side wins", []droneConf{{0, 100, 20}, {1, 50, 10}}, 5, 0},
		{"2 factions - mutual elimination (no owner resolution)", []droneConf{{0, 100, 50}, {1, 100, 50}}, 10, -1},
		{"2 factions - both survive (no owner resolution)", []droneConf{{0, 100, 5}, {1, 100, 5}}, 10, -1},
		{"5 factions - even split free-for-all, highest hp wins", []droneConf{
			{0, 20, 10},
			{1, 30, 10},
			{2, 40, 10},
			{3, 50, 10},
			{4, 60, 10},
		}, 10, 4},
	}

	for _, tt := range tests {
		gs := gamestate.NewGameState(0)
		terr := entity.Register(gs.EM, entity.NewTerritory())
		gs.Territories = append(gs.Territories, terr.Id())

		factions := []*entity.Faction{}

		for range tt.drones {
			faction := entity.Register(gs.EM, entity.NewFaction("test_faction"))
			factions = append(factions, faction)
			gs.Factions = append(gs.Factions, faction.Id())
			terr.AddFaction(faction.Id())
		}

		for _, conf := range tt.drones {
			drone := entity.NewDrone(fmt.Sprintf("test_drone_%d", conf.factionIdx), entity.DroneFighter, 1)
			drone.HP = conf.hp
			drone.Attack = conf.attack
			drone.Activity = entity.DroneFighting
			drone.Target = terr.Id()
			entity.Register(gs.EM, drone)
			assert.Nil(t, factions[conf.factionIdx].AddDrone(drone.Id()))
		}

		var err error
		for range tt.ticks {
			err = gs.ResolveTerritoryConflicts()
		}

		assert.Nil(t, err)
		if tt.expectedOwnerId != -1 {
			assert.Equal(t, factions[tt.expectedOwnerId].Id(), terr.Owner)
		} else {
			assert.Equal(t, tt.expectedOwnerId, terr.Owner)

		}
	}
}

func setupDispatchGame(t *testing.T) (gs *gamestate.GameState, faction *entity.Faction, enemyFaction *entity.Faction, terr *entity.Territory, fighter, miner *entity.Drone) {
	t.Helper()

	gs = gamestate.NewGameState(0)

	faction = entity.Register(gs.EM, entity.NewFaction("faction_a"))
	gs.Factions = append(gs.Factions, faction.Id())

	enemyFaction = entity.Register(gs.EM, entity.NewFaction("faction_b"))
	gs.Factions = append(gs.Factions, enemyFaction.Id())

	terr = entity.Register(gs.EM, entity.NewTerritory())
	gs.Territories = append(gs.Territories, terr.Id())

	fighter = entity.Register(gs.EM, entity.NewDrone("fighter", entity.DroneFighter, 1))
	assert.Nil(t, faction.AddDrone(fighter.Id()))

	miner = entity.Register(gs.EM, entity.NewDrone("miner", entity.DroneMiner, 1))
	assert.Nil(t, faction.AddDrone(miner.Id()))

	return gs, faction, enemyFaction, terr, fighter, miner
}

func TestDispatchDrones_Fighting(t *testing.T) {
	gs, faction, _, terr, fighter, _ := setupDispatchGame(t)

	err := gs.DispatchDrones(faction.Id(), terr.Id(), []entity.EntityId{fighter.Id()}, entity.DroneFighting)
	assert.Nil(t, err)

	got, err := gs.Drone(fighter.Id())
	assert.Nil(t, err)
	assert.Equal(t, entity.DroneFighting, got.Activity)
	assert.Equal(t, terr.Id(), got.Target)
	assert.Contains(t, terr.Factions, faction.Id())
}

func TestDispatchDrones_MiningRequiresOwnership(t *testing.T) {
	gs, faction, _, terr, _, miner := setupDispatchGame(t)

	err := gs.DispatchDrones(faction.Id(), terr.Id(), []entity.EntityId{miner.Id()}, entity.DroneMining)
	assert.ErrorIs(t, err, gamestate.ErrTerritoryNotOwned)

	terr.Owner = faction.Id()
	err = gs.DispatchDrones(faction.Id(), terr.Id(), []entity.EntityId{miner.Id()}, entity.DroneMining)
	assert.Nil(t, err)

	got, err := gs.Drone(miner.Id())
	assert.Nil(t, err)
	assert.Equal(t, entity.DroneMining, got.Activity)
	assert.Equal(t, terr.Id(), got.Target)
}

func TestDispatchDrones_WrongTypeRejected(t *testing.T) {
	gs, faction, _, terr, fighter, _ := setupDispatchGame(t)
	terr.Owner = faction.Id()

	err := gs.DispatchDrones(faction.Id(), terr.Id(), []entity.EntityId{fighter.Id()}, entity.DroneMining)
	assert.ErrorIs(t, err, gamestate.ErrWrongDroneType)
}

func TestDispatchDrones_NotIdleRejected(t *testing.T) {
	gs, faction, _, terr, fighter, _ := setupDispatchGame(t)

	assert.Nil(t, gs.DispatchDrones(faction.Id(), terr.Id(), []entity.EntityId{fighter.Id()}, entity.DroneFighting))
	err := gs.DispatchDrones(faction.Id(), terr.Id(), []entity.EntityId{fighter.Id()}, entity.DroneFighting)
	assert.ErrorIs(t, err, gamestate.ErrDroneNotIdle)
}

func TestDispatchDrones_NotOwnedRejected(t *testing.T) {
	gs, _, enemyFaction, terr, fighter, _ := setupDispatchGame(t)

	err := gs.DispatchDrones(enemyFaction.Id(), terr.Id(), []entity.EntityId{fighter.Id()}, entity.DroneFighting)
	assert.ErrorIs(t, err, gamestate.ErrDroneNotOwned)
}

func TestRecallDrones(t *testing.T) {
	gs, faction, _, terr, fighter, _ := setupDispatchGame(t)

	assert.Nil(t, gs.DispatchDrones(faction.Id(), terr.Id(), []entity.EntityId{fighter.Id()}, entity.DroneFighting))
	assert.Contains(t, terr.Factions, faction.Id())

	assert.Nil(t, gs.RecallDrones(faction.Id(), []entity.EntityId{fighter.Id()}))

	got, err := gs.Drone(fighter.Id())
	assert.Nil(t, err)
	assert.Equal(t, entity.DroneIdle, got.Activity)
	assert.Equal(t, entity.EntityId(-1), got.Target)
	assert.NotContains(t, terr.Factions, faction.Id())
}

func TestRecallDrones_NotCommittedRejected(t *testing.T) {
	gs, faction, _, _, fighter, _ := setupDispatchGame(t)

	err := gs.RecallDrones(faction.Id(), []entity.EntityId{fighter.Id()})
	assert.ErrorIs(t, err, gamestate.ErrDroneNotCommitted)
}

func TestResolveMining(t *testing.T) {
	gs, faction, _, terr, _, miner := setupDispatchGame(t)
	terr.Owner = faction.Id()

	// Total is deliberately not a multiple of MiningRatePerDrone so the
	// second tick has to clamp extraction to whatever remains.
	total := gamestate.MiningRatePerDrone + 3
	deposit := entity.Register(gs.EM, entity.NewResourceDeposit(entity.ResourceMineral, total))
	assert.Nil(t, terr.AddDeposit(deposit.Id()))

	assert.Nil(t, gs.DispatchDrones(faction.Id(), terr.Id(), []entity.EntityId{miner.Id()}, entity.DroneMining))

	assert.Nil(t, gs.ResolveMining())
	assert.Equal(t, 3, deposit.Remaining)
	assert.Equal(t, gamestate.MiningRatePerDrone, faction.OwnedResources[entity.ResourceMineral])

	// second tick should clamp extraction to whatever remains, not go negative
	assert.Nil(t, gs.ResolveMining())
	assert.Equal(t, 0, deposit.Remaining)
	assert.Equal(t, total, faction.OwnedResources[entity.ResourceMineral])
}
