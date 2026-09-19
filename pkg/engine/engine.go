package engine

import (
	"context"
	"time"

	"github.com/elitracy/space-war-sim/pkg/logging"
)

type GameState interface {
	Tick() error
	CurrentTick() int
	IsPaused() bool
}

// onTick, if non-nil, is invoked synchronously after every successful tick
// so callers (e.g. the UI) can be notified of the exact tick as it happens,
// instead of polling CurrentTick() on a separate timer and risking drift.
func RunGame(ctx context.Context, tickInterval time.Duration, logPath string, gs GameState, onTick func(tick int)) error {

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

		if onTick != nil {
			onTick(gs.CurrentTick())
		}

		select {
		case <-ctx.Done():
			logging.Info("Exiting...")
			return ctx.Err()
		case <-time.After(tickInterval):
		}
	}
}
