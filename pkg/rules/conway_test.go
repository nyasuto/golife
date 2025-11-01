package rules

import (
	"golife/pkg/core"
	"testing"
)

func TestConwayRule_Name(t *testing.T) {
	rule := ConwayRule{}
	if rule.Name() != "Conway B3/S23" {
		t.Errorf("Expected 'Conway B3/S23', got '%s'", rule.Name())
	}
}

func TestConwayRule_ShouldBirth(t *testing.T) {
	rule := ConwayRule{}

	tests := []struct {
		neighbors int
		want      bool
	}{
		{0, false},
		{1, false},
		{2, false},
		{3, true}, // Birth at exactly 3 neighbors
		{4, false},
		{5, false},
		{8, false},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := rule.ShouldBirth(tt.neighbors)
			if got != tt.want {
				t.Errorf("ShouldBirth(%d) = %v, want %v", tt.neighbors, got, tt.want)
			}
		})
	}
}

func TestConwayRule_ShouldSurvive(t *testing.T) {
	rule := ConwayRule{}

	tests := []struct {
		neighbors int
		state     core.CellState
		want      bool
	}{
		{0, core.Alive, false},
		{1, core.Alive, false},
		{2, core.Alive, true}, // Survive with 2 neighbors
		{3, core.Alive, true}, // Survive with 3 neighbors
		{4, core.Alive, false},
		{5, core.Alive, false},
		{8, core.Alive, false},
		{3, core.Dead, false}, // Dead cells never survive
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := rule.ShouldSurvive(tt.neighbors, tt.state)
			if got != tt.want {
				t.Errorf("ShouldSurvive(%d, %d) = %v, want %v", tt.neighbors, tt.state, got, tt.want)
			}
		})
	}
}

func TestConwayRule_NeighborWeight(t *testing.T) {
	rule := ConwayRule{}

	// Conway's rule uses uniform weights
	distances := []float64{0.5, 1.0, 1.414, 2.0, 10.0}
	for _, dist := range distances {
		weight := rule.NeighborWeight(dist)
		if weight != 1.0 {
			t.Errorf("NeighborWeight(%f) = %f, want 1.0", dist, weight)
		}
	}
}
