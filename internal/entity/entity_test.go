package entity_test

import (
	"testing"

	"github.com/elitracy/space-war-sim/internal/entity"
	"github.com/stretchr/testify/assert"
)

type mockEntity struct {
	*entity.CoreEntity
}

func (e mockEntity) Tick() error {
	return nil
}

func newMockEntity() *mockEntity {
	return &mockEntity{
		CoreEntity: entity.NewCoreEntity(),
	}
}

func TestEntity_Register(t *testing.T) {
	e0 := newMockEntity()
	e1 := newMockEntity()
	em := entity.NewEntityManager()

	entity.Register(em, e0)
	assert.Equal(t, e0.Id(), entity.EntityId(0))

	entity.Register(em, e1)
	assert.Equal(t, e1.Id(), entity.EntityId(1))
}
