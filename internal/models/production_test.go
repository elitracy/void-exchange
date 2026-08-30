package models_test

import (
	"testing"

	"github.com/elitracy/space-war-sim/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestFactory_Register(t *testing.T) {
	em := models.NewEntityManager()

	generator := models.NewGenerator()
	models.Register(em, generator)

	factory := models.NewFactory(models.ProductDrone)
	models.Register(em, factory)

	assert.Equal(t, generator.Id(), models.EntityId(0))
	assert.Equal(t, factory.Id(), models.EntityId(1))

	for r := range models.ProductRecipes[models.ProductDrone].RequiredResources {
		assert.Equal(t, factory.CurrentResources[r], 0)
	}
}

func TestFactory_Tick(t *testing.T) {
	em := models.NewEntityManager()

	tests := []struct {
		factory            *models.Factory
		generator          *models.Generator
		rate               int
		initialMaterials   int
		remainingMaterials int
		unitsProduced      int
		ticks              int
	}{
		{models.NewFactory(models.ProductDrone), models.NewGenerator(), 1, 5, 0, 1, 10},
		{models.NewFactory(models.ProductDrone), models.NewGenerator(), 1, 3, 3, 0, 10},
		{models.NewFactory(models.ProductGenerator), models.NewGenerator(), 1, 20, 0, 2, 10},
		{models.NewFactory(models.ProductGenerator), models.NewGenerator(), 1, 12, 2, 1, 10},
		{models.NewFactory(models.ProductGenerator), nil, 1, 12, 12, 0, 10},
		{models.NewFactory(models.ProductGenerator), models.NewGenerator(), 2, 40, 0, 4, 2},
	}

	for _, tt := range tests {
		models.Register(em, tt.factory)
		tt.factory.Rate = tt.rate

		if tt.generator != nil {
			models.Register(em, tt.generator)
			tt.factory.SetGenerator(tt.generator.Id())
		}

		assert.Equal(t, tt.factory.CurrentResources[models.ResourceMineral], 0)

		tt.factory.CurrentResources[models.ResourceMineral] = tt.initialMaterials
		for range tt.ticks {
			assert.Equal(t, tt.factory.Tick(), nil)
		}

		assert.Equal(t, tt.remainingMaterials, tt.factory.CurrentResources[models.ResourceMineral])
		assert.Equal(t, tt.unitsProduced, tt.factory.CurrentUnits)
	}

}

func TestCreditMiner_Tick(t *testing.T) {
	em := models.NewEntityManager()

	tests := []struct {
		creditMiner     *models.CreditMiner
		generator       *models.Generator
		rate            int
		creditsProduced int
		ticks           int
	}{
		{models.NewCreditMiner(), models.NewGenerator(), 100, 500, 5},
	}

	for _, tt := range tests {
		models.Register(em, tt.creditMiner)
		tt.creditMiner.Rate = tt.rate

		if tt.generator != nil {
			models.Register(em, tt.generator)
			tt.creditMiner.SetGenerator(tt.generator.Id())
		}

		assert.Equal(t, tt.creditMiner.CurrentCredits, 0)

		for range tt.ticks {
			assert.Equal(t, tt.creditMiner.Tick(), nil)
		}

		assert.Equal(t, tt.creditsProduced, tt.creditMiner.CurrentCredits)
	}

}
