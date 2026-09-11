package api_test

import (
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/elitracy/space-war-sim/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestListScenarios(t *testing.T) {
	tests := []struct {
		names         []string
		expectedNames []string
	}{
		{[]string{"scenario_a.json", "scenario_b.json"}, []string{"scenario_a", "scenario_b"}},
		{[]string{"scenario_a.json", "bad_scenario.txt"}, []string{"scenario_a"}},
	}

	for _, tt := range tests {
		dir := t.TempDir()
		for _, fileName := range tt.names {
			os.WriteFile(filepath.Join(dir, fileName), []byte("{}"), 0644)
		}

		names, err := api.ListScenarios(dir)

		assert.Equal(t, tt.expectedNames, names)
		assert.Nil(t, err)
	}
}

func TestLoadScenario(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test_scenario.json")
	os.WriteFile(path, []byte("{}"), 0644)

	seed := 123
	gs, err := api.LoadScenario(path, seed)

	assert.Equal(t, gs.Rng, rand.New(rand.NewSource(int64(seed))))
	assert.NotNil(t, gs)
	assert.Nil(t, err)
}
