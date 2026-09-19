package api

import (
	"errors"
	"slices"
	"strings"

	"github.com/elitracy/space-war-sim/pkg/entity"
)

type TerritoryView struct {
	ID           entity.EntityId       `json:"id"`
	Owner        entity.EntityId       `json:"owner"`
	Factions     []entity.EntityId     `json:"factions"`
	DepositTypes []entity.ResourceType `json:"deposit_types"`
	Deposits     []entity.EntityId     `json:"deposits"`
}

type DepositView struct {
	ID        entity.EntityId     `json:"id"`
	Type      entity.ResourceType `json:"type"`
	Total     int                 `json:"total"`
	Remaining int                 `json:"remaining"`
}

type FactionView struct {
	ID          entity.EntityId             `json:"id"`
	Name        string                      `json:"name"`
	Territories []entity.EntityId           `json:"territories"`
	Drones      []entity.EntityId           `json:"drones"`
	Resources   map[entity.ResourceType]int `json:"resources"`
}

func NewTerritoryView(em *entity.EntityManager, t *entity.Territory) (*TerritoryView, error) {

	errs := []error{}

	depositTypes := make(map[entity.ResourceType]struct{})
	for _, id := range t.ResourceDeposits {
		deposit, err := entity.GetAs[*entity.ResourceDeposit](em, id)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		depositTypes[deposit.Resource.Type] = struct{}{}
	}

	sortedDepositTypes := []entity.ResourceType{}
	for dType := range depositTypes {
		sortedDepositTypes = append(sortedDepositTypes, dType)
	}

	slices.SortFunc(sortedDepositTypes, func(a, b entity.ResourceType) int {
		return strings.Compare(string(a), string(b))
	})

	return &TerritoryView{
		ID:           t.Id(),
		Owner:        t.Owner,
		Factions:     t.Factions,
		DepositTypes: sortedDepositTypes,
		Deposits:     t.ResourceDeposits,
	}, errors.Join(errs...)
}

func NewDepositView(d *entity.ResourceDeposit) (*DepositView, error) {
	return &DepositView{
		ID:        d.Id(),
		Type:      d.Resource.Type,
		Total:     d.Total,
		Remaining: d.Remaining,
	}, nil
}

func NewFactionView(f *entity.Faction) (*FactionView, error) {
	return &FactionView{
		ID:          f.Id(),
		Name:        f.Name,
		Territories: f.Territories,
		Drones:      f.Fleet,
		Resources:   f.OwnedResources,
	}, nil

}
