package rules

import "golife/pkg/core"

// ConwayRule implements the classic Conway's Game of Life rule: B3/S23
// Birth: exactly 3 neighbors
// Survival: 2 or 3 neighbors
type ConwayRule struct{}

// Name returns the name of this rule
func (r ConwayRule) Name() string {
	return "Conway B3/S23"
}

// ShouldBirth determines if a dead cell should become alive
func (r ConwayRule) ShouldBirth(neighborCount int) bool {
	return neighborCount == 3
}

// ShouldSurvive determines if a live cell should stay alive
func (r ConwayRule) ShouldSurvive(neighborCount int, currentState core.CellState) bool {
	if currentState == core.Dead {
		return false
	}
	return neighborCount == 2 || neighborCount == 3
}

// NeighborWeight returns 1.0 for all neighbors (uniform weight)
func (r ConwayRule) NeighborWeight(distance float64) float64 {
	return 1.0
}
