package patterns

import (
	"golife/pkg/core"
	"golife/pkg/rules"
	"golife/pkg/universe"
)

// Pattern3D represents a 3D pattern with cells in 3D space
type Pattern3D struct {
	Name        string
	Description string
	Width       int
	Height      int
	Depth       int
	Cells       map[core.Coord]core.CellState
}

// LoadIntoUniverse3D loads the pattern into a 3D universe at the specified offset
func (p *Pattern3D) LoadIntoUniverse3D(u *universe.Universe3D, offsetX, offsetY, offsetZ int) {
	for coord, state := range p.Cells {
		newCoord := core.NewCoord3D(coord.X+offsetX, coord.Y+offsetY, coord.Z+offsetZ)
		u.Set(newCoord, state)
	}
}

// CreateUniverse creates a new 3D universe with this pattern loaded at the origin
func (p *Pattern3D) CreateUniverse(rule core.Rule) *universe.Universe3D {
	u := universe.New3D(p.Width, p.Height, p.Depth, rule)
	p.LoadIntoUniverse3D(u, 0, 0, 0)
	return u
}

// BaysGlider returns Bays's Glider (the first discovered 3D glider)
// 10 cells, period 4, moves √2 cells diagonally every 4 generations
// Reference: Carter Bays (1987) "Candidates for the Game of Life in Three Dimensions"
func BaysGlider() *Pattern3D {
	cells := make(map[core.Coord]core.CellState)

	// Bays's Glider Phase 0 (10 cells)
	// This is the verified glider from Bays's original paper
	// It moves at 45° to two axes and orthogonal to the third
	cells[core.NewCoord3D(1, 0, 0)] = core.Alive
	cells[core.NewCoord3D(2, 0, 0)] = core.Alive
	cells[core.NewCoord3D(0, 1, 0)] = core.Alive
	cells[core.NewCoord3D(2, 1, 0)] = core.Alive
	cells[core.NewCoord3D(1, 2, 0)] = core.Alive
	cells[core.NewCoord3D(1, 0, 1)] = core.Alive
	cells[core.NewCoord3D(2, 0, 1)] = core.Alive
	cells[core.NewCoord3D(0, 1, 1)] = core.Alive
	cells[core.NewCoord3D(2, 1, 1)] = core.Alive
	cells[core.NewCoord3D(1, 2, 1)] = core.Alive

	return &Pattern3D{
		Name:        "Bays's Glider",
		Description: "Period-4 glider, moves √2 cells/4 gens (Carter Bays 1987)",
		Width:       5,
		Height:      5,
		Depth:       4,
		Cells:       cells,
	}
}

// Glider3D is an alias for BaysGlider for backward compatibility
func Glider3D() *Pattern3D {
	return BaysGlider()
}

// Block3D returns a 3D block (stable structure)
// 8 cells (2×2×2 cube), the smallest still-life in B6/S567
// Reference: Carter Bays (1987)
func Block3D() *Pattern3D {
	cells := make(map[core.Coord]core.CellState)

	// 2x2x2 cube - stable in B6/S567
	for z := 0; z < 2; z++ {
		for y := 0; y < 2; y++ {
			for x := 0; x < 2; x++ {
				cells[core.NewCoord3D(x, y, z)] = core.Alive
			}
		}
	}

	return &Pattern3D{
		Name:        "3D Block",
		Description: "Stable 2×2×2 cube, 8 cells (Bays 1987)",
		Width:       3,
		Height:      3,
		Depth:       3,
		Cells:       cells,
	}
}

// Beehive3D returns the 3D Beehive still-life
// 14 cells, stable structure
// Reference: Carter Bays (1987)
func Beehive3D() *Pattern3D {
	cells := make(map[core.Coord]core.CellState)

	// Z=0 layer: 2D beehive shape
	cells[core.NewCoord3D(1, 0, 0)] = core.Alive
	cells[core.NewCoord3D(2, 0, 0)] = core.Alive
	cells[core.NewCoord3D(0, 1, 0)] = core.Alive
	cells[core.NewCoord3D(3, 1, 0)] = core.Alive
	cells[core.NewCoord3D(1, 2, 0)] = core.Alive
	cells[core.NewCoord3D(2, 2, 0)] = core.Alive

	// Z=1 layer: same pattern
	cells[core.NewCoord3D(1, 0, 1)] = core.Alive
	cells[core.NewCoord3D(2, 0, 1)] = core.Alive
	cells[core.NewCoord3D(0, 1, 1)] = core.Alive
	cells[core.NewCoord3D(3, 1, 1)] = core.Alive
	cells[core.NewCoord3D(1, 2, 1)] = core.Alive
	cells[core.NewCoord3D(2, 2, 1)] = core.Alive

	// Z=2 layer: one more layer to stabilize in 3D
	cells[core.NewCoord3D(1, 1, 2)] = core.Alive
	cells[core.NewCoord3D(2, 1, 2)] = core.Alive

	return &Pattern3D{
		Name:        "3D Beehive",
		Description: "Stable beehive structure, 14 cells (Bays 1987)",
		Width:       5,
		Height:      4,
		Depth:       4,
		Cells:       cells,
	}
}

// Bucket3D returns the Bucket still-life
// 16 cells, stable container-like structure
// Reference: Carter Bays (1987)
func Bucket3D() *Pattern3D {
	cells := make(map[core.Coord]core.CellState)

	// Bottom layer (z=0): square base
	cells[core.NewCoord3D(1, 1, 0)] = core.Alive
	cells[core.NewCoord3D(2, 1, 0)] = core.Alive
	cells[core.NewCoord3D(1, 2, 0)] = core.Alive
	cells[core.NewCoord3D(2, 2, 0)] = core.Alive

	// Middle layer (z=1): hollow square
	cells[core.NewCoord3D(0, 0, 1)] = core.Alive
	cells[core.NewCoord3D(1, 0, 1)] = core.Alive
	cells[core.NewCoord3D(2, 0, 1)] = core.Alive
	cells[core.NewCoord3D(3, 0, 1)] = core.Alive
	cells[core.NewCoord3D(0, 3, 1)] = core.Alive
	cells[core.NewCoord3D(1, 3, 1)] = core.Alive
	cells[core.NewCoord3D(2, 3, 1)] = core.Alive
	cells[core.NewCoord3D(3, 3, 1)] = core.Alive

	// Sides
	cells[core.NewCoord3D(0, 1, 1)] = core.Alive
	cells[core.NewCoord3D(0, 2, 1)] = core.Alive
	cells[core.NewCoord3D(3, 1, 1)] = core.Alive
	cells[core.NewCoord3D(3, 2, 1)] = core.Alive

	return &Pattern3D{
		Name:        "Bucket",
		Description: "Stable container structure, 16 cells (Bays 1987)",
		Width:       5,
		Height:      5,
		Depth:       3,
		Cells:       cells,
	}
}

// Blinker3D returns the 3D Blinker oscillator
// Period 2, 6 cells (three-cell bar in two adjacent planes)
// Reference: Carter Bays (1987)
func Blinker3D() *Pattern3D {
	cells := make(map[core.Coord]core.CellState)

	// 3-cell bar in z=0 plane
	cells[core.NewCoord3D(0, 1, 0)] = core.Alive
	cells[core.NewCoord3D(1, 1, 0)] = core.Alive
	cells[core.NewCoord3D(2, 1, 0)] = core.Alive

	// 3-cell bar in z=1 plane (same position)
	cells[core.NewCoord3D(0, 1, 1)] = core.Alive
	cells[core.NewCoord3D(1, 1, 1)] = core.Alive
	cells[core.NewCoord3D(2, 1, 1)] = core.Alive

	return &Pattern3D{
		Name:        "3D Blinker",
		Description: "Period-2 oscillator, 6 cells (Bays 1987)",
		Width:       4,
		Height:      3,
		Depth:       3,
		Cells:       cells,
	}
}

// Flashlight3D returns the Flashlight oscillator
// Period 4, 14 cells, flips between two mirror-image states
// Reference: Carter Bays (1987)
func Flashlight3D() *Pattern3D {
	cells := make(map[core.Coord]core.CellState)

	// Phase 0 of the Flashlight pattern
	// Z=0 layer
	cells[core.NewCoord3D(1, 1, 0)] = core.Alive
	cells[core.NewCoord3D(2, 1, 0)] = core.Alive
	cells[core.NewCoord3D(1, 2, 0)] = core.Alive
	cells[core.NewCoord3D(2, 2, 0)] = core.Alive

	// Z=1 layer (middle)
	cells[core.NewCoord3D(0, 1, 1)] = core.Alive
	cells[core.NewCoord3D(3, 1, 1)] = core.Alive
	cells[core.NewCoord3D(1, 0, 1)] = core.Alive
	cells[core.NewCoord3D(2, 0, 1)] = core.Alive
	cells[core.NewCoord3D(1, 3, 1)] = core.Alive
	cells[core.NewCoord3D(2, 3, 1)] = core.Alive

	// Z=2 layer
	cells[core.NewCoord3D(1, 1, 2)] = core.Alive
	cells[core.NewCoord3D(2, 1, 2)] = core.Alive
	cells[core.NewCoord3D(1, 2, 2)] = core.Alive
	cells[core.NewCoord3D(2, 2, 2)] = core.Alive

	return &Pattern3D{
		Name:        "Flashlight",
		Description: "Period-4 oscillator, 14 cells (Bays 1987)",
		Width:       5,
		Height:      5,
		Depth:       4,
		Cells:       cells,
	}
}

// Wheel3D returns the Wheel oscillator
// Period 2, 12 cells, ring that rotates 180° each tick
// Reference: Carter Bays (1987)
func Wheel3D() *Pattern3D {
	cells := make(map[core.Coord]core.CellState)

	// Ring pattern in middle layer (z=1)
	// Outer ring of 8 cells
	cells[core.NewCoord3D(0, 1, 1)] = core.Alive
	cells[core.NewCoord3D(1, 0, 1)] = core.Alive
	cells[core.NewCoord3D(2, 0, 1)] = core.Alive
	cells[core.NewCoord3D(3, 1, 1)] = core.Alive
	cells[core.NewCoord3D(3, 2, 1)] = core.Alive
	cells[core.NewCoord3D(2, 3, 1)] = core.Alive
	cells[core.NewCoord3D(1, 3, 1)] = core.Alive
	cells[core.NewCoord3D(0, 2, 1)] = core.Alive

	// Top and bottom cells
	cells[core.NewCoord3D(1, 1, 0)] = core.Alive
	cells[core.NewCoord3D(2, 2, 0)] = core.Alive
	cells[core.NewCoord3D(1, 2, 2)] = core.Alive
	cells[core.NewCoord3D(2, 1, 2)] = core.Alive

	return &Pattern3D{
		Name:        "Wheel",
		Description: "Period-2 oscillator, 12 cells (Bays 1987)",
		Width:       5,
		Height:      5,
		Depth:       4,
		Cells:       cells,
	}
}

// Oscillator3D_Period2 is an alias for Blinker3D for backward compatibility
func Oscillator3D_Period2() *Pattern3D {
	return Blinker3D()
}

// GetPatterns3D returns a map of all available 3D patterns
func GetPatterns3D() map[string]*Pattern3D {
	return map[string]*Pattern3D{
		// Spaceships
		"bays-glider": BaysGlider(),
		"glider":      Glider3D(), // Alias for backward compatibility

		// Oscillators
		"blinker":    Blinker3D(),
		"flashlight": Flashlight3D(),
		"wheel":      Wheel3D(),
		"oscillator": Oscillator3D_Period2(), // Alias

		// Still Lifes
		"block":   Block3D(),
		"beehive": Beehive3D(),
		"bucket":  Bucket3D(),
	}
}

// ListPatterns3D returns a list of all available 3D pattern names
func ListPatterns3D() []string {
	return []string{
		// Spaceships
		"bays-glider",
		"glider",

		// Oscillators
		"blinker",
		"flashlight",
		"wheel",
		"oscillator",

		// Still Lifes
		"block",
		"beehive",
		"bucket",
	}
}

// LoadPattern3D loads a 3D pattern by name
func LoadPattern3D(name string) *Pattern3D {
	patterns := GetPatterns3D()
	if pattern, ok := patterns[name]; ok {
		return pattern
	}
	return nil
}

// TestPattern3D creates a simple test pattern for development
func TestPattern3D() *Pattern3D {
	cells := make(map[core.Coord]core.CellState)

	// Simple L-shape in 3D
	cells[core.NewCoord3D(1, 1, 1)] = core.Alive
	cells[core.NewCoord3D(2, 1, 1)] = core.Alive
	cells[core.NewCoord3D(1, 2, 1)] = core.Alive

	return &Pattern3D{
		Name:        "Test Pattern",
		Description: "Simple test pattern for debugging",
		Width:       4,
		Height:      4,
		Depth:       3,
		Cells:       cells,
	}
}

// DemoUniverse3D creates a demo 3D universe with B6/S567 rule and a glider
func DemoUniverse3D() *universe.Universe3D {
	rule := rules.Life3D_B6S567{}
	u := universe.New3D(20, 20, 20, rule)

	// Load glider in the center
	glider := Glider3D()
	glider.LoadIntoUniverse3D(u, 8, 8, 8)

	return u
}
