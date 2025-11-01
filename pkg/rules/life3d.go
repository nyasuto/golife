package rules

import "golife/pkg/core"

// Life3D_B6S567 implements the 3D Life rule: B6/S567 (also known as Life 5766)
// This rule is well-researched and produces stable 3D patterns including gliders.
//
// Birth: exactly 6 neighbors
// Survival: 5, 6, or 7 neighbors
// Neighborhood: Moore (26 neighbors in 3D)
//
// References:
// - Carter Bays (1987): "Candidates for the Game of Life in Three Dimensions"
// - This rule supports stable structures and 3D gliders
type Life3D_B6S567 struct{}

// Name returns the name of this rule
func (r Life3D_B6S567) Name() string {
	return "B6/S567 (3D Life)"
}

// ShouldBirth determines if a dead cell should become alive
// In B6/S567, birth occurs only when there are exactly 6 neighbors
func (r Life3D_B6S567) ShouldBirth(neighborCount int) bool {
	return neighborCount == 6
}

// ShouldSurvive determines if a live cell should stay alive
// In B6/S567, survival occurs when there are 5, 6, or 7 neighbors
func (r Life3D_B6S567) ShouldSurvive(neighborCount int, currentState core.CellState) bool {
	if currentState == core.Dead {
		return false
	}
	return neighborCount >= 5 && neighborCount <= 7
}

// NeighborWeight returns 1.0 for all neighbors (uniform weight)
func (r Life3D_B6S567) NeighborWeight(distance float64) float64 {
	return 1.0
}
