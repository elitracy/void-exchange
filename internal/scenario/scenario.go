package scenario

import (
	"encoding/json"
	"io"

	"github.com/elitracy/space-war-sim/internal/entity"
	"github.com/elitracy/space-war-sim/internal/gamestate"
)

type Config struct {
	Seed        int               `json:"seed"`
	Factions    []FactionConfig   `json:"factions"`
	Territories []TerritoryConfig `json:"territories"`
}

type FactionConfig struct {
	Name string `json:"name"`
}

type TerritoryConfig struct {
	Deposits   []DepositConfig `json:"deposits"`
	DepositMin int             `json:"depositMin"`
	DepositMax int             `json:"depositMax"`
}

type DepositConfig struct {
	Type     entity.ResourceType
	Quantity int
}

type CoreScenario struct {
	cfg Config
}

func Load(r io.Reader) (Config, error) {
	var cfg Config
	err := json.NewDecoder(r).Decode(&cfg)
	return cfg, err
}

func Build(cfg Config, gs *gamestate.GameState) {
	for _, f := range cfg.Factions {
		faction := entity.NewFaction(f.Name)
		gs.Factions = append(gs.Factions, faction)
	}

	for _, t := range cfg.Territories {
		territory := entity.NewTerritory()
		entity.Register(gs.EM, territory)

		for _, d := range t.Deposits {
			for range d.Quantity {
				deposit := entity.NewResourceDeposit(d.Type, t.DepositMax)
				entity.Register(gs.EM, deposit)
				territory.AddDeposit(deposit.Id())
			}
		}
	}
}
