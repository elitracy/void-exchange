package models_test

import (
	"testing"

	"github.com/elitracy/space-war-sim/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTerritory_Deposits(t *testing.T) {
	tests := []struct {
		deposits      []models.ResourceType
		depositAmount int
	}{
		{[]models.ResourceType{models.ResourceMineral}, 100},
		{[]models.ResourceType{models.ResourceMineral, models.ResourceEnergy}, 100},
	}

	em := models.NewEntityManager()

	for _, tt := range tests {

		territory := models.NewTerritory()
		models.Register(em, territory)

		for _, d := range tt.deposits {
			deposit := models.NewResourceDeposit(d, tt.depositAmount)
			models.Register(em, deposit)

			if territory.AddDeposit(deposit.Id()) != nil || deposit.Id() == -1 {
				t.Errorf("resource deposit not registered: %d", deposit.Id())
			}
		}

		assert.Equal(t, len(tt.deposits), len(territory.ResourceDeposits))

		for i, d := range territory.ResourceDeposits {
			e, ok := em.GetEntity(d)
			require.True(t, ok, "resource deposit not registered: %d", d)

			deposit, ok := e.(*models.ResourceDeposit)
			require.True(t, ok, "entity is not resource deposit: %d", d)

			assert.Equal(t, tt.deposits[i], deposit.Resource.Type)
			assert.Equal(t, tt.depositAmount, deposit.Total)
			assert.Equal(t, tt.depositAmount, deposit.Remaining)
		}

	}
}

func TestTerritory_UpdateOwner(t *testing.T) {
	tests := []struct {
		owner         models.EntityId
		expectedOwner models.EntityId
	}{
		{20, 20},
		{-1, -1},
	}

	for _, tt := range tests {
		em := models.NewEntityManager()

		territory := models.NewTerritory()
		models.Register(em, territory)

		territory.UpdatedOwner(tt.owner)

		assert.Equal(t, tt.expectedOwner, territory.Owner)
	}
}
