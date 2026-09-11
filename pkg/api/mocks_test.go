package api_test

import (
	"sync"
	"testing"
	"time"
)

type mockGameState struct {
	mu     sync.Mutex
	tick   int
	paused bool
	ticked chan struct{}
}

func newMockGameState() *mockGameState {
	return &mockGameState{
		ticked: make(chan struct{}, 1),
	}
}

func (gs *mockGameState) CurrentTick() int { gs.mu.Lock(); defer gs.mu.Unlock(); return gs.tick }
func (gs *mockGameState) IsPaused() bool   { gs.mu.Lock(); defer gs.mu.Unlock(); return gs.paused }
func (gs *mockGameState) Pause()           { gs.mu.Lock(); defer gs.mu.Unlock(); gs.paused = true }
func (gs *mockGameState) Resume()          { gs.mu.Lock(); defer gs.mu.Unlock(); gs.paused = false }

func (gs *mockGameState) Tick() error {
	gs.mu.Lock()
	gs.tick++
	gs.mu.Unlock()

	select {
	case gs.ticked <- struct{}{}:
	default:
	}

	return nil
}

func waitForTick(t *testing.T, gs *mockGameState) {
	t.Helper()

	select {
	case <-gs.ticked:
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for tick")
	}
}
