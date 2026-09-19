package entity_test

import (
	"testing"

	"github.com/elitracy/void-exchange/pkg/entity"
	"github.com/stretchr/testify/assert"
)

func TestFaction_New(t *testing.T) {
	tests := []struct {
		expectedName  string
		resources     []entity.ResourceType
		resourceDelta int
		expectedId    entity.EntityId
	}{
		{"faction_a", []entity.ResourceType{entity.ResourceEnergy, entity.ResourceMineral}, 10, 0},
	}

	em := entity.NewEntityManager()
	for _, tt := range tests {
		faction := entity.NewFaction(tt.expectedName)
		entity.Register(em, faction)

		assert.Equal(t, tt.expectedName, faction.Name)
		assert.Equal(t, tt.expectedId, faction.Id())
	}
}
func TestFaction_UpdatedResources(t *testing.T) {
	tests := []struct {
		resource       entity.ResourceType
		initialAmount  int
		resourceDelta  int
		ticks          int
		expectedAmount int
	}{
		{entity.ResourceMineral, 0, 10, 1, 10},
		{entity.ResourceMineral, 10, -10, 1, 0},
		{entity.ResourceMineral, 10, -10, 5, 0},
	}

	em := entity.NewEntityManager()
	for _, tt := range tests {
		faction := entity.NewFaction("")
		entity.Register(em, faction)

		faction.OwnedResources[tt.resource] = tt.initialAmount

		for range tt.ticks {
			faction.UpdateResource(tt.resource, tt.resourceDelta)
		}

		assert.Equal(t, tt.expectedAmount, faction.OwnedResources[tt.resource])
	}
}

func TestFaction_AddTerritory(t *testing.T) {
	tests := []struct {
		name      string
		territory *entity.Territory
		register  bool
	}{
		{"faction_a", entity.NewTerritory(), true},
		{"faction_a", entity.NewTerritory(), false},
	}

	for _, tt := range tests {
		em := entity.NewEntityManager()
		faction := entity.NewFaction(tt.name)

		entity.Register(em, faction)

		if tt.register {
			entity.Register(em, tt.territory)
		}

		err := faction.AddTerritory(tt.territory.Id())

		if tt.register {
			assert.Equal(t, tt.territory.Id(), faction.Territories[0])
			assert.NoError(t, err)
		} else {
			assert.Equal(t, entity.ErrInvalidEntityId, err)
		}
	}
}
func TestFaction_AddDrone(t *testing.T) {
	tests := []struct {
		name     string
		Drone    *entity.Drone
		register bool
	}{
		{"faction_a", entity.NewDrone("drone_a", entity.DroneMiner, 1), true},
		{"faction_a", entity.NewDrone("drone_b", entity.DroneMiner, 1), false},
	}

	for _, tt := range tests {
		em := entity.NewEntityManager()
		faction := entity.NewFaction(tt.name)

		entity.Register(em, faction)

		if tt.register {
			entity.Register(em, tt.Drone)
		}

		err := faction.AddDrone(tt.Drone.Id())

		if tt.register {
			assert.Equal(t, tt.Drone.Id(), faction.Fleet[0])
			assert.NoError(t, err)
		} else {
			assert.Equal(t, entity.ErrInvalidEntityId, err)
		}
	}
}

func TestFaction_AddFactory(t *testing.T) {
	tests := []struct {
		name     string
		Factory  *entity.Factory
		register bool
	}{
		{"faction_a", entity.NewFactory(entity.ProductDrone), true},
		{"faction_a", entity.NewFactory(entity.ProductDrone), false},
	}

	for _, tt := range tests {
		em := entity.NewEntityManager()
		faction := entity.NewFaction(tt.name)

		entity.Register(em, faction)

		if tt.register {
			entity.Register(em, tt.Factory)
		}

		err := faction.AddFactory(tt.Factory.Id())

		if tt.register {
			assert.Equal(t, tt.Factory.Id(), faction.Factories[0])
			assert.NoError(t, err)
		} else {
			assert.Equal(t, entity.ErrInvalidEntityId, err)
		}
	}
}

func TestFaction_AddGenerator(t *testing.T) {
	tests := []struct {
		name      string
		Generator *entity.Generator
		register  bool
	}{
		{"faction_a", entity.NewGenerator(), true},
		{"faction_a", entity.NewGenerator(), false},
	}

	for _, tt := range tests {
		em := entity.NewEntityManager()
		faction := entity.NewFaction(tt.name)

		entity.Register(em, faction)

		if tt.register {
			entity.Register(em, tt.Generator)
		}

		err := faction.AddGenerator(tt.Generator.Id())

		if tt.register {
			assert.Equal(t, tt.Generator.Id(), faction.Generators[0])
			assert.NoError(t, err)
		} else {
			assert.Equal(t, entity.ErrInvalidEntityId, err)
		}
	}
}

func TestFaction_AddCreditMiner(t *testing.T) {
	tests := []struct {
		name        string
		CreditMiner *entity.CreditMiner
		register    bool
	}{
		{"faction_a", entity.NewCreditMiner(), true},
		{"faction_a", entity.NewCreditMiner(), false},
	}

	for _, tt := range tests {
		em := entity.NewEntityManager()
		faction := entity.NewFaction(tt.name)

		entity.Register(em, faction)

		if tt.register {
			entity.Register(em, tt.CreditMiner)
		}

		err := faction.AddCreditMiner(tt.CreditMiner.Id())

		if tt.register {
			assert.Equal(t, tt.CreditMiner.Id(), faction.CreditMiners[0])
			assert.NoError(t, err)
		} else {
			assert.Equal(t, entity.ErrInvalidEntityId, err)
		}
	}
}
