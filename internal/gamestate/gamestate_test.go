package gamestate_test

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"testing"

	"github.com/elitracy/space-war-sim/internal/entity"
	"github.com/elitracy/space-war-sim/internal/gamestate"
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
	hp         int
}

func TestResolveTerritoryConflicts(t *testing.T) {
	tests := []struct {
		name            string
		drones          []droneConf
		ticks           int
		expectedOwnerId entity.EntityId
	}{
		{"uncontested single faction does not change owner", []droneConf{{0, 100}}, 10, 0},
		{"2 factions - higher hp wins", []droneConf{{0, 100}, {1, 50}}, 5, 0},
		{"2 factions - mutual elimination (no owner resolution)", []droneConf{{0, 100}, {1, 100}}, 10, -1},
		{"2 factions - both survive (no owner resolution)", []droneConf{{0, 20}, {1, 20}}, 10, -1},
		{"5 factions - highest hp wins", []droneConf{
			{0, 20},
			{1, 30},
			{2, 40},
			{3, 50},
			{4, 60},
		}, 10, 4},
	}

	for _, tt := range tests {
		gs := gamestate.NewGameState(0)
		terr := entity.Register(gs.EM, entity.NewTerritory())
		gs.Territories = append(gs.Territories, terr)

		factions := []*entity.Faction{}

		for range tt.drones {
			faction := entity.Register(gs.EM, entity.NewFaction("test_faction"))
			factions = append(factions, faction)
			gs.Factions = append(gs.Factions, faction)
			terr.AddFaction(faction.Id())
		}

		for _, conf := range tt.drones {
			drone := entity.NewDrone(fmt.Sprintf("test_drone_%d", conf.factionIdx), entity.DroneFighter, conf.hp)
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
