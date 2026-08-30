package models_test

import (
	"testing"

	"github.com/elitracy/space-war-sim/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestDrone_RegisterTypes(t *testing.T) {
	tests := []struct {
		t    models.DroneType
		name string
		id   int
	}{
		{models.DroneMiner, "minerDrone", 0},
		{models.DroneTransport, "transportDrone", 1},
		{models.DroneFighter, "fighterDrone", 2},
	}

	em := models.NewEntityManager()

	for _, tt := range tests {

		drone := models.NewDrone(tt.name, tt.t)
		got := models.Register(em, drone)

		assert.Equal(t, got.Name, tt.name)
		assert.Equal(t, got.Type, tt.t)
		assert.Equal(t, got.Id(), models.EntityId(tt.id))
	}
}
