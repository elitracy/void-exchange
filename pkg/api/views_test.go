package api_test

import (
	"testing"

	"github.com/elitracy/void-exchange/pkg/api"
	"github.com/elitracy/void-exchange/pkg/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDepositView(t *testing.T) {
	d := entity.NewResourceDeposit(entity.ResourceMineral, 100)
	d.Remaining = 40

	view, err := api.NewDepositView(d)

	require.NoError(t, err)
	assert.Equal(t, d.Id(), view.ID)
	assert.Equal(t, entity.ResourceMineral, view.Type)
	assert.Equal(t, 100, view.Total)
	assert.Equal(t, 40, view.Remaining)
}

func TestNewFactionView(t *testing.T) {
	f := entity.NewFaction("faction_a")
	em := entity.NewEntityManager()
	entity.Register(em, f)

	require.NoError(t, f.AddTerritory(entity.EntityId(1)))
	require.NoError(t, f.AddDrone(entity.EntityId(2)))
	require.NoError(t, f.UpdateResource(entity.ResourceMineral, 5))

	view, err := api.NewFactionView(f)

	require.NoError(t, err)
	assert.Equal(t, f.Id(), view.ID)
	assert.Equal(t, "faction_a", view.Name)
	assert.Equal(t, []entity.EntityId{1}, view.Territories)
	assert.Equal(t, []entity.EntityId{2}, view.Drones)
	assert.Equal(t, 5, view.Resources[entity.ResourceMineral])
}

func TestNewTerritoryView(t *testing.T) {
	em := entity.NewEntityManager()

	terr := entity.NewTerritory()
	entity.Register(em, terr)

	mineral := entity.Register(em, entity.NewResourceDeposit(entity.ResourceMineral, 10))
	energy := entity.Register(em, entity.NewResourceDeposit(entity.ResourceEnergy, 10))
	require.NoError(t, terr.AddDeposit(mineral.Id()))
	require.NoError(t, terr.AddDeposit(energy.Id()))
	require.NoError(t, terr.AddFaction(entity.EntityId(7)))
	require.NoError(t, terr.UpdatedOwner(entity.EntityId(7)))

	view, err := api.NewTerritoryView(em, terr)

	require.NoError(t, err)
	assert.Equal(t, terr.Id(), view.ID)
	assert.Equal(t, entity.EntityId(7), view.Owner)
	assert.Equal(t, []entity.EntityId{7}, view.Factions)
	assert.ElementsMatch(t, []entity.EntityId{mineral.Id(), energy.Id()}, view.Deposits)
	// Deposit types are deduplicated by type and sorted for stable UI output.
	assert.Equal(t, []entity.ResourceType{entity.ResourceEnergy, entity.ResourceMineral}, view.DepositTypes)
}

func TestNewTerritoryView_DedupesRepeatedDepositTypes(t *testing.T) {
	em := entity.NewEntityManager()
	terr := entity.NewTerritory()
	entity.Register(em, terr)

	a := entity.Register(em, entity.NewResourceDeposit(entity.ResourceMineral, 10))
	b := entity.Register(em, entity.NewResourceDeposit(entity.ResourceMineral, 20))
	require.NoError(t, terr.AddDeposit(a.Id()))
	require.NoError(t, terr.AddDeposit(b.Id()))

	view, err := api.NewTerritoryView(em, terr)

	require.NoError(t, err)
	assert.Equal(t, []entity.ResourceType{entity.ResourceMineral}, view.DepositTypes,
		"two deposits of the same type must collapse to one entry in DepositTypes")
}

func TestNewTerritoryView_UnknownDepositId_ReturnsPartialViewAndError(t *testing.T) {
	em := entity.NewEntityManager()
	terr := entity.NewTerritory()
	entity.Register(em, terr)
	require.NoError(t, terr.AddDeposit(entity.EntityId(999)))

	view, err := api.NewTerritoryView(em, terr)

	assert.Error(t, err, "a dangling deposit id should surface as an error, not be silently dropped")
	require.NotNil(t, view, "the view should still be usable for the deposits that did resolve")
	assert.Empty(t, view.DepositTypes)
}
