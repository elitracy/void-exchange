package logging_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/elitracy/space-war-sim/pkg/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInfo_WritesToLogFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.log")

	logging.Init(path, 3)
	logging.Info("hello %s", "world")
	logging.Flush()

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Contains(t, string(data), "hello world")
	assert.Contains(t, string(data), "[TELEMETRY]")
}

func TestLogFunctions_BeforeInit_DoNotPanic(t *testing.T) {
	assert.NotPanics(t, func() {
		logging.Info("info")
		logging.Debug("debug")
		logging.Warn("warn")
		logging.Error("error")
		logging.Ok("ok")
	})
}

func TestFlush_WithoutInit_DoesNotPanic(t *testing.T) {
	assert.NotPanics(t, func() { logging.Flush() })
}

func TestFlush_IsIdempotent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.log")
	logging.Init(path, 0)
	logging.Info("first")

	assert.NotPanics(t, func() {
		logging.Flush()
		logging.Flush()
	})
}

func TestLogCallsAfterFlush_AreDroppedNotPanicked(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.log")
	logging.Init(path, 0)
	logging.Flush()

	assert.NotPanics(t, func() { logging.Info("after flush") })
}

func TestReinit_StartsAFreshLogFile(t *testing.T) {
	dir := t.TempDir()
	pathA := filepath.Join(dir, "a.log")
	pathB := filepath.Join(dir, "b.log")

	logging.Init(pathA, 0)
	logging.Info("in a")
	logging.Flush()

	logging.Init(pathB, 0)
	logging.Info("in b")
	logging.Flush()

	dataA, err := os.ReadFile(pathA)
	require.NoError(t, err)
	dataB, err := os.ReadFile(pathB)
	require.NoError(t, err)

	assert.Contains(t, string(dataA), "in a")
	assert.NotContains(t, string(dataA), "in b")
	assert.Contains(t, string(dataB), "in b")
	assert.NotContains(t, string(dataB), "in a")
}

func TestConcurrentInitAndLogCalls_DoNotRace(t *testing.T) {
	dir := t.TempDir()

	done := make(chan struct{})
	go func() {
		defer close(done)
		for i := range 20 {
			logging.Init(filepath.Join(dir, "concurrent.log"), i)
			logging.Info("tick %d", i)
			time.Sleep(time.Millisecond)
			logging.Flush()
		}
	}()

	for i := range 20 {
		logging.Warn("observer %d", i)
		time.Sleep(time.Millisecond)
	}

	<-done
}

func TestInfo_WithNoFormatArgs_DoesNotTreatMessageAsFormatString(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.log")
	logging.Init(path, 0)
	msg := "100" + "%" + " complete"
	infoFn := logging.Info
	infoFn(msg)
	logging.Flush()

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.True(t, strings.Contains(string(data), "100% complete"))
}
