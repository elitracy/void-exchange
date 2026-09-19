package entity_test

import (
	"testing"

	"github.com/elitracy/void-exchange/pkg/entity"
	"github.com/stretchr/testify/assert"
)

func TestDrone_RegisterTypes(t *testing.T) {
	tests := []struct {
		t    entity.DroneType
		name string
		id   int
	}{
		{entity.DroneMiner, "minerDrone", 0},
		{entity.DroneTransport, "transportDrone", 1},
		{entity.DroneFighter, "fighterDrone", 2},
	}

	em := entity.NewEntityManager()

	for _, tt := range tests {

		drone := entity.NewDrone(tt.name, tt.t, 1)
		got := entity.Register(em, drone)

		assert.Equal(t, got.Name, tt.name)
		assert.Equal(t, got.Type, tt.t)
		assert.Equal(t, got.Id(), entity.EntityId(tt.id))
	}
}
