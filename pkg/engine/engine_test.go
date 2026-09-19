package engine_test

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/elitracy/void-exchange/pkg/engine"
	"github.com/stretchr/testify/assert"
)

type mockGameState struct {
	currentTick int
	cancel      context.CancelFunc
}

func (gs mockGameState) CurrentTick() int {
	return gs.currentTick
}

func (gs mockGameState) IsPaused() bool {
	return false
}

func (gs *mockGameState) Tick() error {
	gs.currentTick++
	if gs.currentTick >= 1 {
		gs.cancel()
	}
	return nil
}

type errGameState struct{}

func (gs errGameState) CurrentTick() int {
	return 0
}

func (gs errGameState) IsPaused() bool {
	return false
}

func (gs *errGameState) Tick() error {
	return errors.New("big ahhhh error")
}

func TestRunGame_StopsOnCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	gs := &mockGameState{currentTick: 0, cancel: cancel}
	logPath := filepath.Join(t.TempDir(), "test.log")

	err := engine.RunGame(ctx, time.Second, logPath, gs)
	assert.Equal(t, err, context.Canceled)
}

func TestRunGame_PropogatesError(t *testing.T) {
	gs := &errGameState{}
	ctx := context.Background()
	logPath := filepath.Join(t.TempDir(), "test.log")
	err := engine.RunGame(ctx, 0, logPath, gs)

	assert.EqualError(t, err, "big ahhhh error")
}
