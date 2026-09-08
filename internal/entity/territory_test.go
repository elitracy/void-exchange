package entity_test

import (
	"testing"

	"github.com/elitracy/space-war-sim/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTerritory_Deposits(t *testing.T) {
	tests := []struct {
		deposits      []entity.ResourceType
		depositAmount int
	}{
		{[]entity.ResourceType{entity.ResourceMineral}, 100},
		{[]entity.ResourceType{entity.ResourceMineral, entity.ResourceEnergy}, 100},
	}

	em := entity.NewEntityManager()

	for _, tt := range tests {

		territory := entity.NewTerritory()
		entity.Register(em, territory)

		for _, d := range tt.deposits {
			deposit := entity.NewResourceDeposit(d, tt.depositAmount)
			entity.Register(em, deposit)

			if territory.AddDeposit(deposit.Id()) != nil || deposit.Id() == -1 {
				t.Errorf("resource deposit not registered: %d", deposit.Id())
			}
		}

		assert.Equal(t, len(tt.deposits), len(territory.ResourceDeposits))

		for i, d := range territory.ResourceDeposits {
			e, ok := em.GetEntity(d)
			require.True(t, ok, "resource deposit not registered: %d", d)

			deposit, ok := e.(*entity.ResourceDeposit)
			require.True(t, ok, "entity is not resource deposit: %d", d)

			assert.Equal(t, tt.deposits[i], deposit.Resource.Type)
			assert.Equal(t, tt.depositAmount, deposit.Total)
			assert.Equal(t, tt.depositAmount, deposit.Remaining)
		}

	}
}

func TestTerritory_UpdateOwner(t *testing.T) {
	tests := []struct {
		owner         entity.EntityId
		expectedOwner entity.EntityId
	}{
		{20, 20},
		{-1, -1},
	}

	for _, tt := range tests {
		em := entity.NewEntityManager()

		territory := entity.NewTerritory()
		entity.Register(em, territory)

		territory.UpdatedOwner(tt.owner)

		assert.Equal(t, tt.expectedOwner, territory.Owner)
	}
}

func TestTerritory_RemoveFaction(t *testing.T) {
	tests := []struct {
		name               string
		factionIds         []entity.EntityId
		idsToRemove        []entity.EntityId
		expectedFactionIds []entity.EntityId
		expectRemoved      bool
	}{
		{"remove single id", []entity.EntityId{0, 1, 2}, []entity.EntityId{0}, []entity.EntityId{1, 2}, true},
		{"remove id with empty list", []entity.EntityId{}, []entity.EntityId{0, 1, 2}, []entity.EntityId{}, false},
		{"remove multiple ids", []entity.EntityId{0, 1, 2}, []entity.EntityId{0, 1, 2}, []entity.EntityId{}, true},
		{"remove nonexistant id", []entity.EntityId{0, 1, 2}, []entity.EntityId{100}, []entity.EntityId{0, 1, 2}, false},
	}

	for _, tt := range tests {
		territory := entity.NewTerritory()

		for _, id := range tt.factionIds {
			territory.AddFaction(id)
		}

		for _, id := range tt.idsToRemove {
			removed, err := territory.RemoveFaction(id)
			assert.Nil(t, err)
			assert.Equal(t, tt.expectRemoved, removed)
		}

		assert.Equal(t, tt.expectedFactionIds, territory.Factions)
	}
}
