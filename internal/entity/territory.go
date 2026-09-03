package entity

type TerritoryId int

type Territory struct {
	*CoreEntity
	ResourceDeposits []EntityId
	Owner            EntityId
}

func NewTerritory() *Territory {
	return &Territory{
		CoreEntity: NewCoreEntity(),
		Owner:      EntityId(-1),
	}
}

func (t *Territory) AddDeposit(id EntityId) error {
	if id == -1 {
		return ErrInvalidEntityId
	}

	t.ResourceDeposits = append(t.ResourceDeposits, id)

	return nil
}

func (t *Territory) UpdatedOwner(id EntityId) error {
	if id == -1 {
		return ErrInvalidEntityId
	}

	t.Owner = id

	return nil
}
func (t *Territory) Tick() error { return nil }
