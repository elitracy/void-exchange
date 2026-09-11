package entity_test

import (
	"testing"

	"github.com/elitracy/space-war-sim/pkg/entity"
	"github.com/stretchr/testify/assert"
)

func TestResourceDeposity_Registration(t *testing.T) {
	tests := []struct {
		resourceType entity.ResourceType
		resource     *entity.Resource
		total        int
		id           entity.EntityId
	}{
		{entity.ResourceMineral, entity.NewResource(entity.ResourceMineral), 100, 0},
		{entity.ResourceEnergy, entity.NewResource(entity.ResourceEnergy), 100, 1},
	}

	em := entity.NewEntityManager()

	for _, tt := range tests {

		e := entity.NewResourceDeposit(tt.resourceType, 100)
		entity.Register(em, e)

		tt.resource.Id = e.Resource.Id // might need a manager in the future. Should this be an entity?

		assert.Equal(t, tt.id, e.Id())
		assert.Equal(t, tt.total, e.Total)
		assert.Equal(t, tt.resource, e.Resource)
	}

}
