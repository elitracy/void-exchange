package engine

import (
	"context"
	"fmt"
	"time"

	"github.com/elitracy/space-war-sim/pkg/logging"
)

type GameState interface {
	Tick() error
	CurrentTick() int
}

func RunGame(ctx context.Context, gs GameState, tickInterval time.Duration) error {

	if gs == nil {
		return fmt.Errorf("nil gamestate")
	}

	logging.Init("./logs/debug.log", gs.CurrentTick())

	for {
		select {
		case <-ctx.Done():
			logging.Info("Exiting...")
			return ctx.Err()
		default:
		}

		if err := gs.Tick(); err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			logging.Info("Exiting...")
			return ctx.Err()
		case <-time.After(tickInterval):
		}
	}
}
