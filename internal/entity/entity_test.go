package entity_test

import (
	"testing"

	"github.com/elitracy/space-war-sim/internal/entity"
	"github.com/stretchr/testify/assert"
)

type mockEntity struct {
	*entity.CoreEntity
	tickedCount int
}

func (e *mockEntity) Tick() error {
	e.tickedCount += 1
	return nil
}

func newMockEntity() *mockEntity {
	return &mockEntity{
		CoreEntity: entity.NewCoreEntity(),
	}
}

func TestEntityManager_Register(t *testing.T) {
	e0 := newMockEntity()
	e1 := newMockEntity()
	em := entity.NewEntityManager()

	entity.Register(em, e0)
	assert.Equal(t, e0.Id(), entity.EntityId(0))

	entity.Register(em, e1)
	assert.Equal(t, e1.Id(), entity.EntityId(1))
}

func TestEntityManager_GetEntity(t *testing.T) {
	e0 := newMockEntity()
	em := entity.NewEntityManager()

	entity.Register(em, e0)

	e, ok := em.GetEntity(e0.Id())

	assert.True(t, ok)
	assert.Equal(t, e, e0)
}

func TestEntityManager_Tick(t *testing.T) {
	e0 := newMockEntity()
	em := entity.NewEntityManager()

	entity.Register(em, e0)

	err := em.Tick()

	assert.Nil(t, err)
	assert.Equal(t, 1, e0.tickedCount)

}
