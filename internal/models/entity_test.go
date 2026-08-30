package models_test

import (
	"testing"

	"github.com/elitracy/space-war-sim/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestEntity_Register(t *testing.T) {
	e0 := models.NewCoreEntity()
	e1 := models.NewCoreEntity()
	em := models.NewEntityManager()

	models.Register(em, e0)
	assert.Equal(t, e0.Id(), models.EntityId(0))

	models.Register(em, e1)
	assert.Equal(t, e1.Id(), models.EntityId(1))
}
