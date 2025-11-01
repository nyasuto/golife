package patterns

import (
	"golife/pkg/core"
	"golife/pkg/universe"
)

// Pattern25D represents a 2.5D pattern with multiple layers
type Pattern25D struct {
	Name        string
	Description string
	Width       int
	Height      int
	Depth       int
	Cells       map[core.Coord]core.CellState
}

// VerticalGlider creates a pattern that moves upward through layers
// Requires layer interaction to be enabled
func VerticalGlider() *Pattern25D {
	cells := make(map[core.Coord]core.CellState)

	// Create a column pattern that encourages vertical movement
	// Layer 0: Base pattern
	cells[core.NewCoord3D(5, 4, 0)] = core.Alive
	cells[core.NewCoord3D(5, 5, 0)] = core.Alive
	cells[core.NewCoord3D(5, 6, 0)] = core.Alive
	cells[core.NewCoord3D(6, 5, 0)] = core.Alive

	// Layer 1: Middle pattern
	cells[core.NewCoord3D(5, 5, 1)] = core.Alive
	cells[core.NewCoord3D(6, 4, 1)] = core.Alive
	cells[core.NewCoord3D(6, 5, 1)] = core.Alive
	cells[core.NewCoord3D(6, 6, 1)] = core.Alive

	// Layer 2: Top pattern (lighter)
	cells[core.NewCoord3D(5, 5, 2)] = core.Alive
	cells[core.NewCoord3D(6, 5, 2)] = core.Alive

	return &Pattern25D{
		Name:        "Vertical Glider",
		Description: "Pattern that propagates through layers vertically",
		Width:       12,
		Height:      12,
		Depth:       5,
		Cells:       cells,
	}
}

// LayerOscillator creates a pattern that oscillates between layers
func LayerOscillator() *Pattern25D {
	cells := make(map[core.Coord]core.CellState)

	// Layer 0: Active phase
	cells[core.NewCoord3D(5, 5, 0)] = core.Alive
	cells[core.NewCoord3D(6, 5, 0)] = core.Alive
	cells[core.NewCoord3D(7, 5, 0)] = core.Alive
	cells[core.NewCoord3D(5, 6, 0)] = core.Alive
	cells[core.NewCoord3D(7, 6, 0)] = core.Alive

	// Layer 1: Support structure
	cells[core.NewCoord3D(6, 5, 1)] = core.Alive
	cells[core.NewCoord3D(6, 6, 1)] = core.Alive

	// Layer 2: Inactive phase (will activate alternately)
	cells[core.NewCoord3D(6, 4, 2)] = core.Alive
	cells[core.NewCoord3D(6, 7, 2)] = core.Alive

	return &Pattern25D{
		Name:        "Layer Oscillator",
		Description: "Pattern that oscillates between layers",
		Width:       15,
		Height:      15,
		Depth:       3,
		Cells:       cells,
	}
}

// LayerStack creates a stable multi-layer structure
func LayerStack() *Pattern25D {
	cells := make(map[core.Coord]core.CellState)

	// Create a 2x2 block in each layer, offset slightly
	for z := 0; z < 3; z++ {
		offsetX := z
		offsetY := z

		cells[core.NewCoord3D(5+offsetX, 5+offsetY, z)] = core.Alive
		cells[core.NewCoord3D(6+offsetX, 5+offsetY, z)] = core.Alive
		cells[core.NewCoord3D(5+offsetX, 6+offsetY, z)] = core.Alive
		cells[core.NewCoord3D(6+offsetX, 6+offsetY, z)] = core.Alive
	}

	return &Pattern25D{
		Name:        "Layer Stack",
		Description: "Stable structure spanning multiple layers",
		Width:       15,
		Height:      15,
		Depth:       3,
		Cells:       cells,
	}
}

// VerticalBlinker creates a blinker that spans layers
func VerticalBlinker() *Pattern25D {
	cells := make(map[core.Coord]core.CellState)

	// Vertical line through all layers
	for z := 0; z < 3; z++ {
		cells[core.NewCoord3D(5, 5, z)] = core.Alive
		cells[core.NewCoord3D(5, 6, z)] = core.Alive
		cells[core.NewCoord3D(5, 7, z)] = core.Alive
	}

	return &Pattern25D{
		Name:        "Vertical Blinker",
		Description: "Blinker pattern spanning multiple layers",
		Width:       12,
		Height:      12,
		Depth:       3,
		Cells:       cells,
	}
}

// LayerSandwich creates a pattern with activity in outer layers
func LayerSandwich() *Pattern25D {
	cells := make(map[core.Coord]core.CellState)

	// Layer 0: Glider
	cells[core.NewCoord3D(2, 1, 0)] = core.Alive
	cells[core.NewCoord3D(3, 2, 0)] = core.Alive
	cells[core.NewCoord3D(1, 3, 0)] = core.Alive
	cells[core.NewCoord3D(2, 3, 0)] = core.Alive
	cells[core.NewCoord3D(3, 3, 0)] = core.Alive

	// Layer 1: Empty (middle layer)

	// Layer 2: Blinker
	cells[core.NewCoord3D(7, 5, 2)] = core.Alive
	cells[core.NewCoord3D(8, 5, 2)] = core.Alive
	cells[core.NewCoord3D(9, 5, 2)] = core.Alive

	return &Pattern25D{
		Name:        "Layer Sandwich",
		Description: "Active patterns in outer layers with empty middle",
		Width:       15,
		Height:      10,
		Depth:       3,
		Cells:       cells,
	}
}

// EnergyWave creates a pattern designed for energy diffusion
func EnergyWave() *Pattern25D {
	cells := make(map[core.Coord]core.CellState)

	// Layer 0: High energy core
	for y := 4; y <= 6; y++ {
		for x := 4; x <= 6; x++ {
			cells[core.NewCoord3D(x, y, 0)] = 200 // High energy
		}
	}

	// Layer 1: Medium energy ring
	for x := 3; x <= 7; x++ {
		cells[core.NewCoord3D(x, 3, 1)] = 100
		cells[core.NewCoord3D(x, 7, 1)] = 100
	}
	for y := 4; y <= 6; y++ {
		cells[core.NewCoord3D(3, y, 1)] = 100
		cells[core.NewCoord3D(7, y, 1)] = 100
	}

	// Layer 2: Low energy scattered
	cells[core.NewCoord3D(5, 5, 2)] = 50
	cells[core.NewCoord3D(2, 2, 2)] = 30
	cells[core.NewCoord3D(8, 8, 2)] = 30

	return &Pattern25D{
		Name:        "Energy Wave",
		Description: "Pattern designed for energy diffusion demonstration",
		Width:       12,
		Height:      12,
		Depth:       3,
		Cells:       cells,
	}
}

// LoadIntoUniverse25D loads a 2.5D pattern into a Universe25D
func (p *Pattern25D) LoadIntoUniverse25D(u *universe.Universe25D, offsetX, offsetY, offsetZ int) {
	for coord, state := range p.Cells {
		newCoord := core.NewCoord3D(
			coord.X+offsetX,
			coord.Y+offsetY,
			coord.Z+offsetZ,
		)
		u.Set(newCoord, state)
	}
}

// CreateUniverse creates a new Universe25D with this pattern
func (p *Pattern25D) CreateUniverse(rule core.Rule) *universe.Universe25D {
	u := universe.New25D(p.Width, p.Height, p.Depth, rule)
	p.LoadIntoUniverse25D(u, 0, 0, 0)
	return u
}

// GetPatterns25D returns a map of all available 2.5D patterns
func GetPatterns25D() map[string]*Pattern25D {
	return map[string]*Pattern25D{
		"vertical-glider":  VerticalGlider(),
		"layer-oscillator": LayerOscillator(),
		"layer-stack":      LayerStack(),
		"vertical-blinker": VerticalBlinker(),
		"layer-sandwich":   LayerSandwich(),
		"energy-wave":      EnergyWave(),
	}
}

// ListPatterns25D returns a list of available 2.5D pattern names
func ListPatterns25D() []string {
	return []string{
		"vertical-glider",
		"layer-oscillator",
		"layer-stack",
		"vertical-blinker",
		"layer-sandwich",
		"energy-wave",
	}
}
