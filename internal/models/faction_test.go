package models_test

import (
	"testing"

	"github.com/elitracy/space-war-sim/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestFaction_New(t *testing.T) {
	tests := []struct {
		expectedName  string
		resources     []models.ResourceType
		resourceDelta int
		expectedId    models.EntityId
	}{
		{"faction_a", []models.ResourceType{models.ResourceEnergy, models.ResourceMineral}, 10, 0},
	}

	em := models.NewEntityManager()
	for _, tt := range tests {
		faction := models.NewFaction(tt.expectedName)
		models.Register(em, faction)

		assert.Equal(t, tt.expectedName, faction.Name)
		assert.Equal(t, tt.expectedId, faction.Id())
	}
}
func TestFaction_UpdatedResources(t *testing.T) {
	tests := []struct {
		resource       models.ResourceType
		initialAmount  int
		resourceDelta  int
		ticks          int
		expectedAmount int
	}{
		{models.ResourceMineral, 0, 10, 1, 10},
		{models.ResourceMineral, 10, -10, 1, 0},
		{models.ResourceMineral, 10, -10, 5, 0},
	}

	em := models.NewEntityManager()
	for _, tt := range tests {
		faction := models.NewFaction("")
		models.Register(em, faction)

		faction.Resources[tt.resource] = tt.initialAmount

		for range tt.ticks {
			faction.UpdateResource(tt.resource, tt.resourceDelta)
		}

		assert.Equal(t, tt.expectedAmount, faction.Resources[tt.resource])
	}
}

func TestFaction_AddTerritory(t *testing.T) {
	tests := []struct {
		name      string
		territory *models.Territory
		register  bool
	}{
		{"faction_a", models.NewTerritory(), true},
		{"faction_a", models.NewTerritory(), false},
	}

	for _, tt := range tests {
		em := models.NewEntityManager()
		faction := models.NewFaction(tt.name)

		models.Register(em, faction)

		if tt.register {
			models.Register(em, tt.territory)
		}

		err := faction.AddTerritory(tt.territory.Id())

		if tt.register {
			assert.Equal(t, tt.territory.Id(), faction.Territories[0])
			assert.NoError(t, err)
		} else {
			assert.Equal(t, models.ErrInvalidEntityId, err)
		}
	}
}
func TestFaction_AddDrone(t *testing.T) {
	tests := []struct {
		name     string
		Drone    *models.Drone
		register bool
	}{
		{"faction_a", models.NewDrone("drone_a", models.DroneMiner), true},
		{"faction_a", models.NewDrone("drone_b", models.DroneMiner), false},
	}

	for _, tt := range tests {
		em := models.NewEntityManager()
		faction := models.NewFaction(tt.name)

		models.Register(em, faction)

		if tt.register {
			models.Register(em, tt.Drone)
		}

		err := faction.AddDrone(tt.Drone.Id())

		if tt.register {
			assert.Equal(t, tt.Drone.Id(), faction.Fleet[0])
			assert.NoError(t, err)
		} else {
			assert.Equal(t, models.ErrInvalidEntityId, err)
		}
	}
}

func TestFaction_AddFactory(t *testing.T) {
	tests := []struct {
		name     string
		Factory  *models.Factory
		register bool
	}{
		{"faction_a", models.NewFactory(models.ProductDrone), true},
		{"faction_a", models.NewFactory(models.ProductDrone), false},
	}

	for _, tt := range tests {
		em := models.NewEntityManager()
		faction := models.NewFaction(tt.name)

		models.Register(em, faction)

		if tt.register {
			models.Register(em, tt.Factory)
		}

		err := faction.AddFactory(tt.Factory.Id())

		if tt.register {
			assert.Equal(t, tt.Factory.Id(), faction.Factories[0])
			assert.NoError(t, err)
		} else {
			assert.Equal(t, models.ErrInvalidEntityId, err)
		}
	}
}

func TestFaction_AddGenerator(t *testing.T) {
	tests := []struct {
		name      string
		Generator *models.Generator
		register  bool
	}{
		{"faction_a", models.NewGenerator(), true},
		{"faction_a", models.NewGenerator(), false},
	}

	for _, tt := range tests {
		em := models.NewEntityManager()
		faction := models.NewFaction(tt.name)

		models.Register(em, faction)

		if tt.register {
			models.Register(em, tt.Generator)
		}

		err := faction.AddGenerator(tt.Generator.Id())

		if tt.register {
			assert.Equal(t, tt.Generator.Id(), faction.Generators[0])
			assert.NoError(t, err)
		} else {
			assert.Equal(t, models.ErrInvalidEntityId, err)
		}
	}
}

func TestFaction_AddCreditMiner(t *testing.T) {
	tests := []struct {
		name        string
		CreditMiner *models.CreditMiner
		register    bool
	}{
		{"faction_a", models.NewCreditMiner(), true},
		{"faction_a", models.NewCreditMiner(), false},
	}

	for _, tt := range tests {
		em := models.NewEntityManager()
		faction := models.NewFaction(tt.name)

		models.Register(em, faction)

		if tt.register {
			models.Register(em, tt.CreditMiner)
		}

		err := faction.AddCreditMiner(tt.CreditMiner.Id())

		if tt.register {
			assert.Equal(t, tt.CreditMiner.Id(), faction.CreditMiners[0])
			assert.NoError(t, err)
		} else {
			assert.Equal(t, models.ErrInvalidEntityId, err)
		}
	}
}
