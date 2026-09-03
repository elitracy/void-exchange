package entity

import "fmt"

type EntityId int

type Entity interface {
	Id() EntityId
	setId(id EntityId)
	Tick() error
}

type CoreEntity struct {
	id EntityId
}

func NewCoreEntity() *CoreEntity {
	return &CoreEntity{id: -1}
}

func (ce CoreEntity) Id() EntityId       { return ce.id }
func (ce *CoreEntity) setId(id EntityId) { ce.id = id }

type EntityManager struct {
	currentEntityId EntityId
	Entities        map[EntityId]Entity
}

func NewEntityManager() *EntityManager {
	return &EntityManager{
		currentEntityId: 0,
		Entities:        make(map[EntityId]Entity),
	}
}

func Register[T Entity](em *EntityManager, e T) T {

	e.setId(em.currentEntityId)
	em.Entities[e.Id()] = e
	em.currentEntityId++

	return e
}

func (em EntityManager) GetEntity(id EntityId) (Entity, bool) {
	e, ok := em.Entities[id]

	if !ok {
		return nil, ok
	}

	return e, ok

}

func (em EntityManager) Tick() error {
	for id, entity := range em.Entities {
		err := entity.Tick()
		if err != nil {
			return fmt.Errorf("ticking entity (%d): %e", id, err)
		}
	}

	return nil
}
