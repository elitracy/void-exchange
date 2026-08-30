package models_test

import (
	"testing"

	"github.com/elitracy/space-war-sim/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestResourceDeposity_Registration(t *testing.T) {
	tests := []struct {
		resourceType models.ResourceType
		resource     *models.Resource
		total        int
		id           models.EntityId
	}{
		{models.ResourceMineral, models.NewResource(models.ResourceMineral), 100, 0},
		{models.ResourceEnergy, models.NewResource(models.ResourceEnergy), 100, 1},
	}

	em := models.NewEntityManager()

	for _, tt := range tests {

		e := models.NewResourceDeposit(tt.resourceType, 100)
		models.Register(em, e)

		tt.resource.Id = e.Resource.Id // might need a manager in the future. Should this be an entity?

		assert.Equal(t, tt.id, e.Id())
		assert.Equal(t, tt.total, e.Total)
		assert.Equal(t, tt.resource, e.Resource)
	}

}
