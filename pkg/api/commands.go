package api

import (
	"github.com/elitracy/space-war-sim/pkg/gamestate"
	"github.com/elitracy/space-war-sim/pkg/scenario"
	"os"
	"strings"
)

type GameState interface {
	Tick() error
	CurrentTick() int
	IsPaused() bool
	Pause()
	Resume()
}

func ListScenarios(dir string) ([]string, error) {
	var scenarioNames []string
	files, err := os.ReadDir(dir)

	if err != nil {
		return scenarioNames, err
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		name := strings.Split(file.Name(), ".")
		scenarioNames = append(scenarioNames, name[0])
	}

	return scenarioNames, nil
}

func LoadScenario(path string, seed int) (*gamestate.GameState, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	cfg, err := scenario.Load(file)

	if err != nil {
		return nil, err
	}

	gs := gamestate.NewGameState(int64(seed))
	scenario.Build(cfg, gs)

	return gs, nil
}
