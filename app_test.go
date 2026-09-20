package main

import (
	"context"
	"testing"
	"time"

	"github.com/elitracy/void-exchange/pkg/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestApp(t *testing.T) *App {
	t.Helper()
	a := NewApp()
	a.startup(context.Background())
	a.emitEvent = func(context.Context, string, ...interface{}) {}
	return a
}

func TestApp_ListScenarios(t *testing.T) {
	a := newTestApp(t)

	names, err := a.ListScenarios()

	require.NoError(t, err)
	assert.Contains(t, names, "simple")
	assert.Contains(t, names, "medium")
}

func TestApp_QueriesBeforeScenarioLoaded_ReturnErrors(t *testing.T) {
	a := newTestApp(t)

	_, err := a.GetTerritories()
	assert.Error(t, err, "GetTerritories should refuse to run without a loaded scenario")

	_, err = a.GetFactions()
	assert.Error(t, err, "GetFactions should refuse to run without a loaded scenario")

	err = a.StartRun(10)
	assert.Error(t, err, "StartRun should refuse to run without a loaded scenario")
}

func TestApp_ControlsBeforeScenarioLoaded_DoNotPanic(t *testing.T) {
	a := newTestApp(t)

	assert.NotPanics(t, func() { a.Pause() })
	assert.NotPanics(t, func() { a.Resume() })
	assert.NotPanics(t, func() { a.Stop() })

	var paused bool
	assert.NotPanics(t, func() { paused = a.IsPaused() })
	assert.False(t, paused)

	assert.Equal(t, -1, a.CurrentTick())
}

func TestApp_LoadScenario_PopulatesFactionsAndTerritories(t *testing.T) {
	a := newTestApp(t)

	err := a.LoadScenario("simple", 1)
	require.NoError(t, err)
	assert.Equal(t, 0, a.CurrentTick())

	factions, err := a.GetFactions()
	require.NoError(t, err)
	require.Len(t, factions, 1)
	assert.Equal(t, "faction_a", factions[0].Name)

	territories, err := a.GetTerritories()
	require.NoError(t, err)
	require.Len(t, territories, 1)
	assert.Len(t, territories[0].DepositTypes, 1)
}

func TestApp_GetTerritory_GetDeposits_GetDeposit(t *testing.T) {
	a := newTestApp(t)
	require.NoError(t, a.LoadScenario("simple", 1))

	territories, err := a.GetTerritories()
	require.NoError(t, err)
	require.Len(t, territories, 1)
	territoryId := territories[0].ID

	territory, err := a.GetTerritory(territoryId)
	require.NoError(t, err)
	assert.Equal(t, territoryId, territory.ID)

	deposits, err := a.GetDeposits(territoryId)
	require.NoError(t, err)
	require.Len(t, deposits, 1)
	assert.Equal(t, entity.ResourceType("resource_mineral"), deposits[0].Type)

	deposit, err := a.GetDeposit(deposits[0].ID)
	require.NoError(t, err)
	assert.Equal(t, deposits[0].ID, deposit.ID)
}

func TestApp_GetDeposits_UnknownTerritory_Errors(t *testing.T) {
	a := newTestApp(t)
	require.NoError(t, a.LoadScenario("simple", 1))

	_, err := a.GetDeposits(entity.EntityId(9999))
	assert.Error(t, err)
}

func TestApp_StartPauseResumeStop_Lifecycle(t *testing.T) {
	a := newTestApp(t)
	require.NoError(t, a.LoadScenario("simple", 1))

	require.NoError(t, a.StartRun(5))

	require.Eventually(t, func() bool { return a.CurrentTick() > 0 }, time.Second, time.Millisecond,
		"tick should advance once the run starts")

	err := a.StartRun(5)
	assert.Error(t, err, "starting a run that is already active should not silently spawn a second ticker")

	a.Pause()
	require.Eventually(t, func() bool { return a.IsPaused() }, time.Second, time.Millisecond)
	frozen := a.CurrentTick()
	time.Sleep(20 * time.Millisecond)
	assert.Equal(t, frozen, a.CurrentTick(), "tick must not advance while paused")

	a.Resume()
	require.Eventually(t, func() bool { return a.CurrentTick() > frozen }, time.Second, time.Millisecond)

	a.Stop()
	time.Sleep(20 * time.Millisecond)
}
