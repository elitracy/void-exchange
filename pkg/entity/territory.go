package entity

type TerritoryId int

type Territory struct {
	*CoreEntity
	ResourceDeposits []EntityId
	Owner            EntityId
	Factions         []EntityId
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

func (t *Territory) AddFaction(id EntityId) error {
	if id == -1 {
		return ErrInvalidEntityId
	}

	t.Factions = append(t.Factions, id)

	return nil
}

func (t *Territory) RemoveFaction(id EntityId) (removed bool, err error) {
	if id == -1 {
		return false, ErrInvalidEntityId
	}

	factions := []EntityId{}
	removed = false
	for _, faction_id := range t.Factions {
		if faction_id == id {
			removed = true
			continue
		}

		factions = append(factions, faction_id)
	}

	t.Factions = factions
	return removed, nil
}

func (t *Territory) UpdatedOwner(id EntityId) error {
	if id == -1 {
		return ErrInvalidEntityId
	}

	t.Owner = id

	return nil
}
func (t *Territory) Tick() error { return nil }
