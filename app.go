package main

import (
	"context"
	"fmt"
	"time"

	"github.com/elitracy/void-exchange/pkg/api"
	"github.com/elitracy/void-exchange/pkg/entity"
	"github.com/elitracy/void-exchange/pkg/gamestate"
	"github.com/elitracy/void-exchange/pkg/logging"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx    context.Context
	gs     *gamestate.GameState
	runner *api.Runner
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) ListScenarios() ([]string, error) {
	names, err := api.ListScenarios("./scenarios")
	return names, err
}

func (a *App) LoadScenario(name string, seed int) error {
	path := fmt.Sprintf("./scenarios/%s.json", name)
	gs, err := api.LoadScenario(path, seed)

	a.gs = gs
	a.runner = api.NewRunner(gs, "./logs/run.log")

	logging.Info("Loaded scenario: %s", name)
	return err
}

func (a *App) StartRun(tickMs int) error {
	if a.runner == nil {
		return fmt.Errorf("no scenario loaded")
	}

	done, err := a.runner.Start(a.ctx, time.Duration(tickMs)*time.Millisecond)

	if err != nil {
		return err
	}

	go func() {
		err := <-done
		runtime.EventsEmit(a.ctx, "run:finished", fmt.Sprint(err))
	}()

	return err
}

func (a *App) Pause()  { a.runner.Pause() }
func (a *App) Resume() { a.runner.Resume() }
func (a *App) Stop()   { a.runner.Stop() }
func (a *App) CurrentTick() int {
	if a.gs == nil {
		return -1
	}

	return a.gs.CurrentTick()
}
func (a *App) IsPaused() bool { return a.gs.IsPaused() }

func (a *App) GetTerritories() ([]*api.TerritoryView, error) {
	if a.gs == nil || a.runner == nil {
		return nil, fmt.Errorf("no scenario loaded")
	}

	tViews := []*api.TerritoryView{}
	for _, id := range a.gs.Territories {
		t, err := a.gs.Territory(id)
		if err != nil {
			return nil, err
		}

		view, err := api.NewTerritoryView(a.gs.EM, t)
		if err != nil {
			return nil, fmt.Errorf("couldn't create territory view (%d): %w", t.Id(), err)
		}

		tViews = append(tViews, view)
	}

	return tViews, nil
}

func (a *App) GetTerritory(id entity.EntityId) (*api.TerritoryView, error) {

	t, err := entity.GetAs[*entity.Territory](a.gs.EM, id)
	if err != nil {
		return nil, err
	}

	view, err := api.NewTerritoryView(a.gs.EM, t)
	if err != nil {
		return nil, fmt.Errorf("couldn't create territory view (%d): %w", t.Id(), err)
	}

	return view, nil

}

func (a *App) GetDeposits(territoryId entity.EntityId) ([]*api.DepositView, error) {
	if a.gs == nil || a.runner == nil {
		return nil, fmt.Errorf("no scenario loaded")
	}

	territory, err := entity.GetAs[*entity.Territory](a.gs.EM, territoryId)
	if err != nil {
		return nil, err
	}

	dViews := []*api.DepositView{}
	for _, id := range territory.ResourceDeposits {
		d, err := entity.GetAs[*entity.ResourceDeposit](a.gs.EM, id)
		if err != nil {
			return nil, fmt.Errorf("invalid deposit: %w", err)
		}

		view, err := api.NewDepositView(d)
		if err != nil {
			return nil, fmt.Errorf("couldn't create deposit view (%d): %w", d.Id(), err)
		}

		dViews = append(dViews, view)
	}

	return dViews, nil
}

func (a *App) GetDeposit(depositId entity.EntityId) (*api.DepositView, error) {
	if a.gs == nil || a.runner == nil {
		return nil, fmt.Errorf("no scenario loaded")
	}

	deposit, err := a.gs.Deposit(depositId)
	if err != nil {
		return nil, err
	}

	view, err := api.NewDepositView(deposit)
	if err != nil {
		return nil, fmt.Errorf("couldn't create deposit view (%d): %w", deposit.Id(), err)
	}

	return view, nil
}

func (a *App) GetFactions() ([]*api.FactionView, error) {
	if a.gs == nil || a.runner == nil {
		return nil, fmt.Errorf("no scenario loaded")
	}

	fViews := []*api.FactionView{}
	for _, id := range a.gs.Factions {
		f, err := a.gs.Faction(id)
		if err != nil {
			return nil, err
		}

		view, err := api.NewFactionView(f)
		if err != nil {
			return nil, fmt.Errorf("couldn't create faction view (%d): %w", f.Id(), err)
		}

		fViews = append(fViews, view)
	}

	return fViews, nil
}

func (a *App) GetDrones() ([]*api.DroneView, error) {
	if a.gs == nil || a.runner == nil {
		return nil, fmt.Errorf("no scenario loaded")
	}

	dViews := []*api.DroneView{}
	for _, factionId := range a.gs.Factions {
		faction, err := a.gs.Faction(factionId)
		if err != nil {
			return nil, err
		}

		for _, droneId := range faction.Fleet {
			d, err := a.gs.Drone(droneId)
			if err != nil {
				return nil, fmt.Errorf("invalid drone (%d): %w", droneId, err)
			}

			view, err := api.NewDroneView(factionId, d)
			if err != nil {
				return nil, fmt.Errorf("couldn't create drone view (%d): %w", d.Id(), err)
			}

			dViews = append(dViews, view)
		}
	}

	return dViews, nil
}

// DispatchDrones: drone_mining requires the faction already own the territory.
func (a *App) DispatchDrones(factionId entity.EntityId, territoryId entity.EntityId, droneIds []entity.EntityId, activity entity.DroneActivity) error {
	if a.gs == nil || a.runner == nil {
		return fmt.Errorf("no scenario loaded")
	}

	return a.gs.DispatchDrones(factionId, territoryId, droneIds, activity)
}

func (a *App) RecallDrones(factionId entity.EntityId, droneIds []entity.EntityId) error {
	if a.gs == nil || a.runner == nil {
		return fmt.Errorf("no scenario loaded")
	}

	return a.gs.RecallDrones(factionId, droneIds)
}
