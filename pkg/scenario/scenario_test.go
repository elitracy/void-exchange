package scenario_test

import (
	"math/rand"
	"strings"
	"testing"

	"github.com/elitracy/space-war-sim/pkg/entity"
	"github.com/elitracy/space-war-sim/pkg/gamestate"
	"github.com/elitracy/space-war-sim/pkg/scenario"
	"github.com/stretchr/testify/assert"
)

func TestLoad_JSON(t *testing.T) {
	input := `
	{
		"factions": [
			{
				"name": "test_faction"
			}
		],
		"territories": [
			{
				"deposits": [
					{
						"type": "resource_mineral",
						"quantity": 1
					}
				],
				"depositMin": 1,
				"depositMax": 1
			}
		]
	}
	`

	cfg, err := scenario.Load(strings.NewReader(input))

	assert.Nil(t, err)

	assert.Equal(t, scenario.FactionConfig{Name: "test_faction"}, cfg.Factions[0])
	assert.Equal(t, scenario.TerritoryConfig{
		Deposits: []scenario.DepositConfig{
			{Type: entity.ResourceMineral, Quantity: 1},
		},
		DepositMin: 1,
		DepositMax: 1,
	},
		cfg.Territories[0],
	)
}
func TestBuild(t *testing.T) {
	cfg := scenario.Config{
		Factions: []scenario.FactionConfig{
			{"test_faction_a"},
			{"test_faction_b"},
		},
		Territories: []scenario.TerritoryConfig{
			{
				Deposits: []scenario.DepositConfig{
					{Type: entity.ResourceMineral, Quantity: 1},
				},
				DepositMin: 1,
				DepositMax: 1,
			},
			{
				Deposits: []scenario.DepositConfig{
					{Type: entity.ResourceEnergy, Quantity: 2},
					{Type: entity.ResourceMineral, Quantity: 1},
				},
				DepositMin: 1,
				DepositMax: 1,
			},
		},
	}

	gs := gamestate.NewGameState(0)
	scenario.Build(cfg, gs)

	assert.Equal(t, 2, len(gs.Factions))
	assert.Equal(t, "test_faction_a", gs.Factions[0].Name)
	assert.Equal(t, "test_faction_b", gs.Factions[1].Name)
	assert.Equal(t, rand.New(rand.NewSource(0)), gs.Rng)
	assert.Equal(t, 1, len(gs.Territories[0].ResourceDeposits))
	assert.Equal(t, 3, len(gs.Territories[1].ResourceDeposits))
}
