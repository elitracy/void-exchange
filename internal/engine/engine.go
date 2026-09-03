package engine

import (
	"time"

	"github.com/elitracy/space-war-sim/internal/gamestate"
	"github.com/elitracy/space-war-sim/internal/logging"
)

func RunGame(gs *gamestate.GameState) error {
	logging.Init("./logs/debug.log", gs.CurrentTick)

	for {

		// update
		if err := gs.Tick(); err != nil {
			return err
		}

		time.Sleep(time.Second)
	}
}
