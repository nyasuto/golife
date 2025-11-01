package universe

import (
	"golife/pkg/core"
	"golife/pkg/rules"
)

// Universe25D represents a 2.5D universe with multiple 2D layers
type Universe25D struct {
	width, height, depth int
	layers               []*Universe2D
	layerInteraction     bool
	verticalWeight       float64 // Weight for vertical neighbors (0.0-1.0) - deprecated, use interactionRule
	rule                 core.Rule
	interactionRule      rules.LayerInteractionRule // Optional: layer interaction rule
}

// New25D creates a new 2.5D universe with the given dimensions and rule
func New25D(width, height, depth int, rule core.Rule) *Universe25D {
	layers := make([]*Universe2D, depth)
	for z := 0; z < depth; z++ {
		layers[z] = New2D(width, height, rule)
	}

	// Default to weighted neighbors rule for backward compatibility
	defaultInteractionRule := rules.NewWeightedNeighborsRule(rule, 0.3)

	return &Universe25D{
		width:            width,
		height:           height,
		depth:            depth,
		layers:           layers,
		layerInteraction: false,
		verticalWeight:   0.3, // Deprecated: kept for backward compatibility
		rule:             rule,
		interactionRule:  defaultInteractionRule,
	}
}

// Dimension returns the dimensionality (2.5D)
func (u *Universe25D) Dimension() core.Dimension {
	return core.Dim25D
}

// Get returns the state of a cell at the given coordinate
func (u *Universe25D) Get(coord core.Coord) core.CellState {
	if coord.Z < 0 || coord.Z >= u.depth {
		return core.Dead
	}
	return u.layers[coord.Z].Get(core.NewCoord2D(coord.X, coord.Y))
}

// Set sets the state of a cell at the given coordinate
func (u *Universe25D) Set(coord core.Coord, state core.CellState) {
	if coord.Z >= 0 && coord.Z < u.depth {
		u.layers[coord.Z].Set(core.NewCoord2D(coord.X, coord.Y), state)
	}
}

// Size returns the dimensions of the universe
func (u *Universe25D) Size() core.Coord {
	return core.NewCoord3D(u.width, u.height, u.depth)
}

// GetLayer returns a specific layer
func (u *Universe25D) GetLayer(z int) *Universe2D {
	if z >= 0 && z < u.depth {
		return u.layers[z]
	}
	return nil
}

// SetLayerInteraction enables or disables layer interaction
func (u *Universe25D) SetLayerInteraction(enabled bool) {
	u.layerInteraction = enabled
}

// SetVerticalWeight sets the weight for vertical neighbors (0.0-1.0)
// Deprecated: Use SetInteractionRule with WeightedNeighborsRule instead
func (u *Universe25D) SetVerticalWeight(weight float64) {
	if weight >= 0.0 && weight <= 1.0 {
		u.verticalWeight = weight
		// Update interaction rule to match
		u.interactionRule = rules.NewWeightedNeighborsRule(u.rule, weight)
	}
}

// SetInteractionRule sets the layer interaction rule
func (u *Universe25D) SetInteractionRule(rule rules.LayerInteractionRule) {
	u.interactionRule = rule
}

// GetInteractionRule returns the current layer interaction rule
func (u *Universe25D) GetInteractionRule() rules.LayerInteractionRule {
	return u.interactionRule
}

// Step executes one generation
func (u *Universe25D) Step() {
	if u.layerInteraction {
		u.stepWithInteraction()
	} else {
		u.stepIndependent()
	}
}

// stepIndependent processes each layer independently
func (u *Universe25D) stepIndependent() {
	// Each layer evolves independently
	for _, layer := range u.layers {
		layer.Step()
	}
}

// stepWithInteraction processes layers with vertical influence
func (u *Universe25D) stepWithInteraction() {
	// Create temporary storage for new states
	newLayers := make([]*Universe2D, u.depth)
	for z := 0; z < u.depth; z++ {
		newLayers[z] = New2D(u.width, u.height, u.rule)
	}

	// Process each cell considering vertical neighbors
	for z := 0; z < u.depth; z++ {
		for y := 0; y < u.height; y++ {
			for x := 0; x < u.width; x++ {
				coord2D := core.NewCoord2D(x, y)

				// Count horizontal neighbors (same layer)
				horizontalNeighbors := u.layers[z].countNeighbors(x, y)

				// Count vertical neighbors (upper and lower layers)
				verticalNeighbors := u.countVerticalNeighbors(x, y, z)

				// Get current and adjacent layer states
				currentState := u.layers[z].Get(coord2D)
				var upperState, lowerState core.CellState
				if z > 0 {
					upperState = u.layers[z-1].Get(coord2D)
				}
				if z < u.depth-1 {
					lowerState = u.layers[z+1].Get(coord2D)
				}

				// Use interaction rule to calculate effective neighbor count
				neighborCount := u.interactionRule.CalculateNeighborCount(
					horizontalNeighbors, verticalNeighbors,
					currentState, upperState, lowerState,
				)

				// Apply rule with layer interaction
				var newState core.CellState
				if currentState == core.Dead {
					if u.interactionRule.ShouldBirth(neighborCount, upperState, lowerState) {
						newState = core.Alive
					} else {
						newState = core.Dead
					}
				} else {
					if u.interactionRule.ShouldSurvive(neighborCount, currentState, upperState, lowerState) {
						newState = core.Alive
					} else {
						newState = core.Dead
					}
				}

				newLayers[z].Set(coord2D, newState)
			}
		}
	}

	// Copy new states to layers
	u.layers = newLayers
}

// countVerticalNeighbors counts alive cells in the 9 positions above and below
func (u *Universe25D) countVerticalNeighbors(x, y, z int) int {
	count := 0

	// Check layer above (z-1)
	if z > 0 {
		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				nx := x + dx
				ny := y + dy
				if nx >= 0 && nx < u.width && ny >= 0 && ny < u.height {
					if u.layers[z-1].Get(core.NewCoord2D(nx, ny)) > core.Dead {
						count++
					}
				}
			}
		}
	}

	// Check layer below (z+1)
	if z < u.depth-1 {
		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				nx := x + dx
				ny := y + dy
				if nx >= 0 && nx < u.width && ny >= 0 && ny < u.height {
					if u.layers[z+1].Get(core.NewCoord2D(nx, ny)) > core.Dead {
						count++
					}
				}
			}
		}
	}

	return count
}

// Clone creates a deep copy of the universe
func (u *Universe25D) Clone() core.Universe {
	clone := New25D(u.width, u.height, u.depth, u.rule)
	clone.layerInteraction = u.layerInteraction
	clone.verticalWeight = u.verticalWeight

	for z := 0; z < u.depth; z++ {
		for y := 0; y < u.height; y++ {
			for x := 0; x < u.width; x++ {
				coord := core.NewCoord2D(x, y)
				state := u.layers[z].Get(coord)
				clone.layers[z].Set(coord, state)
			}
		}
	}

	return clone
}

// Clear resets all cells to dead state
func (u *Universe25D) Clear() {
	for _, layer := range u.layers {
		layer.Clear()
	}
}

// CountLiving returns the number of living cells across all layers
func (u *Universe25D) CountLiving() int {
	count := 0
	for _, layer := range u.layers {
		count += layer.CountLiving()
	}
	return count
}

// Randomize fills all layers with random cells
func (u *Universe25D) Randomize() {
	for _, layer := range u.layers {
		layer.Randomize()
	}
}

// CountLivingInLayer returns the number of living cells in a specific layer
func (u *Universe25D) CountLivingInLayer(z int) int {
	if z >= 0 && z < u.depth {
		return u.layers[z].CountLiving()
	}
	return 0
}
