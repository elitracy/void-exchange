package entity

type Faction struct {
	*CoreEntity
	Name           string
	Territories    []EntityId
	Fleet          []EntityId
	Factories      []EntityId
	Generators     []EntityId
	CreditMiners   []EntityId
	OwnedResources map[ResourceType]int
}

func NewFaction(name string) *Faction {
	return &Faction{
		CoreEntity:     NewCoreEntity(),
		Name:           name,
		OwnedResources: make(map[ResourceType]int),
	}
}

func (f *Faction) AddTerritory(id EntityId) error {
	if id == -1 {
		return ErrInvalidEntityId
	}

	f.Territories = append(f.Territories, id)

	return nil
}

func (f *Faction) AddDrone(id EntityId) error {
	if id == -1 {
		return ErrInvalidEntityId
	}

	f.Fleet = append(f.Fleet, id)

	return nil
}

func (f *Faction) AddFactory(id EntityId) error {
	if id == -1 {
		return ErrInvalidEntityId
	}

	f.Factories = append(f.Factories, id)

	return nil
}

func (f *Faction) AddGenerator(id EntityId) error {
	if id == -1 {
		return ErrInvalidEntityId
	}

	f.Generators = append(f.Generators, id)

	return nil
}

func (f *Faction) AddCreditMiner(id EntityId) error {
	if id == -1 {
		return ErrInvalidEntityId
	}

	f.CreditMiners = append(f.CreditMiners, id)

	return nil
}

func (f *Faction) UpdateResource(r ResourceType, quantity int) error {
	f.OwnedResources[r] += quantity
	f.OwnedResources[r] = max(f.OwnedResources[r], 0)

	return nil
}
