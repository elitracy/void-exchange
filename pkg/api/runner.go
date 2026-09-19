package api

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/elitracy/void-exchange/pkg/engine"
	"github.com/elitracy/void-exchange/pkg/logging"
)

type Runner struct {
	gs      GameState
	logPath string

	runningMu sync.Mutex
	running   bool
	cancel    context.CancelFunc
	done      chan error
}

func NewRunner(gs GameState, logPath string) *Runner {
	return &Runner{
		gs:      gs,
		logPath: logPath,
	}
}

func (r *Runner) Start(appCtx context.Context, tickInterval time.Duration, onTick func(tick int)) (chan error, error) {
	r.runningMu.Lock()
	if r.running {
		r.runningMu.Unlock()
		return nil, fmt.Errorf("already running")
	}

	simCtx, cancel := context.WithCancel(appCtx)
	r.cancel = cancel

	done := make(chan error, 1)
	r.running = true
	r.runningMu.Unlock()

	go func() {
		err := engine.RunGame(simCtx, tickInterval, r.logPath, r.gs, onTick)
		r.runningMu.Lock()
		r.running = false
		r.runningMu.Unlock()
		done <- err
	}()

	return done, nil
}

func (r *Runner) Stop() {
	r.runningMu.Lock()
	defer r.runningMu.Unlock()

	if r.cancel != nil {
		r.cancel()
	}
	logging.Info("Stopped run...")
}

func (r *Runner) Pause() {
	r.gs.Pause()
	logging.Info("Paused run")
}
func (r *Runner) Resume() {
	r.gs.Resume()
	logging.Info("Resuming run")
}
