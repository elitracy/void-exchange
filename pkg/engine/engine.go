package engine

import (
	"context"
	"time"

	"github.com/elitracy/void-exchange/pkg/logging"
)

type GameState interface {
	Tick() error
	CurrentTick() int
	IsPaused() bool
}

func RunGame(ctx context.Context, tickInterval time.Duration, logPath string, gs GameState) error {

	logging.Init(logPath, gs.CurrentTick())

	for {
		select {
		case <-ctx.Done():
			logging.Info("Exiting...")
			return ctx.Err()
		default:
		}

		if gs.IsPaused() {
			time.Sleep(tickInterval)
			continue
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
