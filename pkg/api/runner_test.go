package api_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/elitracy/space-war-sim/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestStartStop(t *testing.T) {
	logPath := filepath.Join(t.TempDir(), "test.log")
	gs := newMockGameState()
	runner := api.NewRunner(gs, logPath)

	done, err := runner.Start(context.Background(), time.Millisecond, nil)
	assert.Nil(t, err)

	for range 10 {
		waitForTick(t, gs)
	}

	runner.Stop()
	err = <-done
	assert.Equal(t, context.Canceled, err)
}

func TestPauseResume(t *testing.T) {
	logPath := filepath.Join(t.TempDir(), "test.log")
	gs := newMockGameState()
	runner := api.NewRunner(gs, logPath)

	done, err := runner.Start(context.Background(), time.Millisecond, nil)
	assert.Nil(t, err)

	waitForTick(t, gs)
	gs.Pause()
	assert.True(t, gs.IsPaused())
	frozenTick := gs.CurrentTick()

	time.Sleep(20 * time.Millisecond)
	assert.Equal(t, frozenTick, gs.CurrentTick())

	select {
	case err = <-done:
		t.Fatal("run declard done while paused")
	default:
	}

	gs.Resume()

	waitForTick(t, gs)
	runner.Stop()
	err = <-done
	assert.ErrorIs(t, context.Canceled, err)
}
