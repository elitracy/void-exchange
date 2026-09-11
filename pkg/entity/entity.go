package entity

import (
	"errors"
	"fmt"
)

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
	Entities        []EntityId
	entityLookup    map[EntityId]Entity
}

func NewEntityManager() *EntityManager {
	return &EntityManager{
		currentEntityId: 0,
		entityLookup:    make(map[EntityId]Entity),
	}
}

func Register[T Entity](em *EntityManager, e T) T {

	e.setId(em.currentEntityId)
	em.entityLookup[e.Id()] = e
	em.Entities = append(em.Entities, e.Id())

	em.currentEntityId++

	return e
}

func GetAs[T Entity](em *EntityManager, id EntityId) (T, error) {
	e, ok := em.GetEntity(id)

	if !ok {
		return *new(T), ErrEntityNotFound
	}

	typed, ok := e.(T)

	if !ok { // NOTE: shouldn't ever get hit because a non-Entity can't be passed at compile time to func
		return *new(T), ErrInvalidEntity
	}

	return typed, nil
}

func (em EntityManager) GetEntity(id EntityId) (Entity, bool) {
	e, ok := em.entityLookup[id]

	if !ok {
		return nil, ok
	}

	return e, ok

}

func (em *EntityManager) Tick() error {
	var errs []error
	for _, id := range em.Entities {
		e, ok := em.GetEntity(id)
		if !ok {
			return fmt.Errorf("invalid entity id %d", id)
		}

		err := e.Tick()
		if err != nil {
			errs = append(errs, fmt.Errorf("ticking entity (%d): %w", id, err))
		}
	}

	return errors.Join(errs...)
}
