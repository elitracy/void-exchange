package models

type TerritoryId int

type Territory struct {
	*CoreEntity
	ResourceDeposits []EntityId
}

func NewTerritory() *Territory {
	return &Territory{
		CoreEntity: NewCoreEntity(),
	}
}

func (t *Territory) AddDeposit(id EntityId) error {
	if id == -1 {
		return ErrInvalidEntityId
	}

	t.ResourceDeposits = append(t.ResourceDeposits, id)

	return nil
}
