package entity_test

import (
	"testing"

	"github.com/elitracy/void-exchange/pkg/entity"
	"github.com/stretchr/testify/assert"
)

func TestFactory_Register(t *testing.T) {
	em := entity.NewEntityManager()

	generator := entity.NewGenerator()
	entity.Register(em, generator)

	factory := entity.NewFactory(entity.ProductDrone)
	entity.Register(em, factory)

	assert.Equal(t, generator.Id(), entity.EntityId(0))
	assert.Equal(t, factory.Id(), entity.EntityId(1))

	for r := range entity.ProductRecipes[entity.ProductDrone].RequiredResources {
		assert.Equal(t, factory.CurrentResources[r], 0)
	}
}

func TestFactory_Tick(t *testing.T) {
	em := entity.NewEntityManager()

	tests := []struct {
		factory            *entity.Factory
		generator          *entity.Generator
		rate               int
		initialMaterials   int
		remainingMaterials int
		unitsProduced      int
		ticks              int
	}{
		{entity.NewFactory(entity.ProductDrone), entity.NewGenerator(), 1, 5, 0, 1, 10},
		{entity.NewFactory(entity.ProductDrone), entity.NewGenerator(), 1, 3, 3, 0, 10},
		{entity.NewFactory(entity.ProductGenerator), entity.NewGenerator(), 1, 20, 0, 2, 10},
		{entity.NewFactory(entity.ProductGenerator), entity.NewGenerator(), 1, 12, 2, 1, 10},
		{entity.NewFactory(entity.ProductGenerator), nil, 1, 12, 12, 0, 10},
		{entity.NewFactory(entity.ProductGenerator), entity.NewGenerator(), 2, 40, 0, 4, 2},
	}

	for _, tt := range tests {
		entity.Register(em, tt.factory)
		tt.factory.Rate = tt.rate

		if tt.generator != nil {
			entity.Register(em, tt.generator)
			tt.factory.SetGenerator(tt.generator.Id())
		}

		assert.Equal(t, tt.factory.CurrentResources[entity.ResourceMineral], 0)

		tt.factory.CurrentResources[entity.ResourceMineral] = tt.initialMaterials
		for range tt.ticks {
			assert.Equal(t, tt.factory.Tick(), nil)
		}

		assert.Equal(t, tt.remainingMaterials, tt.factory.CurrentResources[entity.ResourceMineral])
		assert.Equal(t, tt.unitsProduced, tt.factory.CurrentUnits)
	}

}

func TestCreditMiner_Tick(t *testing.T) {
	em := entity.NewEntityManager()

	tests := []struct {
		creditMiner     *entity.CreditMiner
		generator       *entity.Generator
		rate            int
		creditsProduced int
		ticks           int
	}{
		{entity.NewCreditMiner(), entity.NewGenerator(), 100, 500, 5},
	}

	for _, tt := range tests {
		entity.Register(em, tt.creditMiner)
		tt.creditMiner.Rate = tt.rate

		if tt.generator != nil {
			entity.Register(em, tt.generator)
			tt.creditMiner.SetGenerator(tt.generator.Id())
		}

		assert.Equal(t, tt.creditMiner.CurrentCredits, 0)

		for range tt.ticks {
			assert.Equal(t, tt.creditMiner.Tick(), nil)
		}

		assert.Equal(t, tt.creditsProduced, tt.creditMiner.CurrentCredits)
	}

}
