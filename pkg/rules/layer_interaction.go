package rules

import "golife/pkg/core"

// LayerInteractionType defines how layers interact in 2.5D universe
type LayerInteractionType int

const (
	// NoInteraction - layers evolve independently
	NoInteraction LayerInteractionType = iota
	// WeightedNeighbors - vertical neighbors contribute with weight
	WeightedNeighbors
	// BirthBetweenLayers - cells can only birth between existing layers
	BirthBetweenLayers
	// EnergyDiffusion - cell states diffuse between layers
	EnergyDiffusion
)

// LayerInteractionRule defines how cells interact across layers in 2.5D
type LayerInteractionRule interface {
	// Type returns the interaction type
	Type() LayerInteractionType

	// CalculateNeighborCount calculates effective neighbor count considering vertical layers
	// horizontalNeighbors: count from same layer (0-8 for Moore)
	// verticalNeighbors: count from upper and lower layers (0-18)
	// currentState: current cell state
	// upperState, lowerState: cell states in layers above/below
	CalculateNeighborCount(horizontalNeighbors, verticalNeighbors int, currentState, upperState, lowerState core.CellState) int

	// ShouldBirth determines if a dead cell should be born
	// neighborCount: effective neighbor count from CalculateNeighborCount
	// upperState, lowerState: cell states in layers above/below
	ShouldBirth(neighborCount int, upperState, lowerState core.CellState) bool

	// ShouldSurvive determines if a live cell should survive
	// neighborCount: effective neighbor count from CalculateNeighborCount
	// currentState: current cell state
	// upperState, lowerState: cell states in layers above/below
	ShouldSurvive(neighborCount int, currentState, upperState, lowerState core.CellState) bool
}

// WeightedNeighborsRule implements weighted neighbor counting
type WeightedNeighborsRule struct {
	baseRule       core.Rule
	verticalWeight float64 // 0.0-1.0
}

// NewWeightedNeighborsRule creates a new weighted neighbors rule
func NewWeightedNeighborsRule(baseRule core.Rule, verticalWeight float64) *WeightedNeighborsRule {
	if verticalWeight < 0.0 {
		verticalWeight = 0.0
	}
	if verticalWeight > 1.0 {
		verticalWeight = 1.0
	}
	return &WeightedNeighborsRule{
		baseRule:       baseRule,
		verticalWeight: verticalWeight,
	}
}

func (r *WeightedNeighborsRule) Type() LayerInteractionType {
	return WeightedNeighbors
}

func (r *WeightedNeighborsRule) CalculateNeighborCount(horizontalNeighbors, verticalNeighbors int, currentState, upperState, lowerState core.CellState) int {
	// Weighted total: horizontal + (vertical × weight)
	totalNeighbors := float64(horizontalNeighbors) + float64(verticalNeighbors)*r.verticalWeight
	return int(totalNeighbors + 0.5) // Round to nearest int
}

func (r *WeightedNeighborsRule) ShouldBirth(neighborCount int, upperState, lowerState core.CellState) bool {
	return r.baseRule.ShouldBirth(neighborCount)
}

func (r *WeightedNeighborsRule) ShouldSurvive(neighborCount int, currentState, upperState, lowerState core.CellState) bool {
	return r.baseRule.ShouldSurvive(neighborCount, currentState)
}

// BirthBetweenLayersRule allows birth only where adjacent layers have live cells
type BirthBetweenLayersRule struct {
	baseRule          core.Rule
	requireBothLayers bool // If true, both upper AND lower must have cells
}

// NewBirthBetweenLayersRule creates a new birth between layers rule
func NewBirthBetweenLayersRule(baseRule core.Rule, requireBothLayers bool) *BirthBetweenLayersRule {
	return &BirthBetweenLayersRule{
		baseRule:          baseRule,
		requireBothLayers: requireBothLayers,
	}
}

func (r *BirthBetweenLayersRule) Type() LayerInteractionType {
	return BirthBetweenLayers
}

func (r *BirthBetweenLayersRule) CalculateNeighborCount(horizontalNeighbors, verticalNeighbors int, currentState, upperState, lowerState core.CellState) int {
	// Only horizontal neighbors count
	return horizontalNeighbors
}

func (r *BirthBetweenLayersRule) ShouldBirth(neighborCount int, upperState, lowerState core.CellState) bool {
	// Check base rule first
	if !r.baseRule.ShouldBirth(neighborCount) {
		return false
	}

	// Check layer condition
	hasUpper := upperState > core.Dead
	hasLower := lowerState > core.Dead

	if r.requireBothLayers {
		// Requires cells in both adjacent layers
		return hasUpper && hasLower
	} else {
		// Requires cells in at least one adjacent layer
		return hasUpper || hasLower
	}
}

func (r *BirthBetweenLayersRule) ShouldSurvive(neighborCount int, currentState, upperState, lowerState core.CellState) bool {
	return r.baseRule.ShouldSurvive(neighborCount, currentState)
}

// EnergyDiffusionRule treats cell state as energy that diffuses between layers
type EnergyDiffusionRule struct {
	baseRule        core.Rule
	diffusionRate   float64 // 0.0-1.0, how much energy transfers
	energyThreshold uint8   // Minimum energy for cell to be "alive"
}

// NewEnergyDiffusionRule creates a new energy diffusion rule
func NewEnergyDiffusionRule(baseRule core.Rule, diffusionRate float64, energyThreshold uint8) *EnergyDiffusionRule {
	if diffusionRate < 0.0 {
		diffusionRate = 0.0
	}
	if diffusionRate > 1.0 {
		diffusionRate = 1.0
	}
	if energyThreshold == 0 {
		energyThreshold = 1
	}
	return &EnergyDiffusionRule{
		baseRule:        baseRule,
		diffusionRate:   diffusionRate,
		energyThreshold: energyThreshold,
	}
}

func (r *EnergyDiffusionRule) Type() LayerInteractionType {
	return EnergyDiffusion
}

func (r *EnergyDiffusionRule) CalculateNeighborCount(horizontalNeighbors, verticalNeighbors int, currentState, upperState, lowerState core.CellState) int {
	// Count neighbors based on energy threshold
	effectiveVertical := 0

	// Add diffused energy from vertical neighbors
	diffusedEnergy := (float64(upperState) + float64(lowerState)) * r.diffusionRate

	// If diffused energy exceeds threshold, count as neighbors
	if diffusedEnergy >= float64(r.energyThreshold) {
		effectiveVertical = int(diffusedEnergy / float64(r.energyThreshold))
	}

	return horizontalNeighbors + effectiveVertical
}

func (r *EnergyDiffusionRule) ShouldBirth(neighborCount int, upperState, lowerState core.CellState) bool {
	return r.baseRule.ShouldBirth(neighborCount)
}

func (r *EnergyDiffusionRule) ShouldSurvive(neighborCount int, currentState, upperState, lowerState core.CellState) bool {
	return r.baseRule.ShouldSurvive(neighborCount, currentState)
}

// GetDiffusedEnergy calculates the new energy level after diffusion
func (r *EnergyDiffusionRule) GetDiffusedEnergy(currentEnergy, upperEnergy, lowerEnergy core.CellState) core.CellState {
	// Energy diffuses from adjacent layers
	incomingEnergy := (float64(upperEnergy) + float64(lowerEnergy)) * r.diffusionRate / 2.0

	// Current energy decays slightly
	decayedEnergy := float64(currentEnergy) * (1.0 - r.diffusionRate*0.1)

	totalEnergy := decayedEnergy + incomingEnergy

	if totalEnergy > 255 {
		return 255
	}
	return core.CellState(totalEnergy)
}
