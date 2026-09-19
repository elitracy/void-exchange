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

// RemoveDrone drops a single drone from the faction's fleet, e.g. when it is
// destroyed in combat. Unlike the old conflict resolution code, this never
// touches the rest of the fleet.
func (f *Faction) RemoveDrone(id EntityId) (removed bool, err error) {
	if id == -1 {
		return false, ErrInvalidEntityId
	}

	fleet := []EntityId{}
	for _, droneId := range f.Fleet {
		if droneId == id {
			removed = true
			continue
		}

		fleet = append(fleet, droneId)
	}

	f.Fleet = fleet
	return removed, nil
}

// HasDrone reports whether the given drone id belongs to this faction's
// fleet, used to validate dispatch/recall commands.
func (f *Faction) HasDrone(id EntityId) bool {
	for _, droneId := range f.Fleet {
		if droneId == id {
			return true
		}
	}
	return false
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

func (f *Faction) Tick() error { return nil }
