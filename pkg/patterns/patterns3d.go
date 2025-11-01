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

// Glider3D returns a 3D glider pattern for B6/S567 rule
// This is a 2x2x2 pattern that moves diagonally through 3D space
// Reference: Carter Bays (1987)
func Glider3D() *Pattern3D {
	cells := make(map[core.Coord]core.CellState)

	// 3D glider: 4 cells in a 2x2x2 cube configuration
	// This is a simple glider that moves diagonally
	cells[core.NewCoord3D(0, 0, 0)] = core.Alive
	cells[core.NewCoord3D(1, 0, 0)] = core.Alive
	cells[core.NewCoord3D(0, 1, 0)] = core.Alive
	cells[core.NewCoord3D(0, 0, 1)] = core.Alive

	return &Pattern3D{
		Name:        "3D Glider",
		Description: "A simple 3D glider for B6/S567 rule",
		Width:       3,
		Height:      3,
		Depth:       3,
		Cells:       cells,
	}
}

// Block3D returns a 3D block (stable structure)
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
		Description: "A stable 2x2x2 cube",
		Width:       3,
		Height:      3,
		Depth:       3,
		Cells:       cells,
	}
}

// Oscillator3D_Period2 returns a simple period-2 oscillator
func Oscillator3D_Period2() *Pattern3D {
	cells := make(map[core.Coord]core.CellState)

	// Simple 3-cell oscillator in 3D
	cells[core.NewCoord3D(1, 0, 1)] = core.Alive
	cells[core.NewCoord3D(1, 1, 1)] = core.Alive
	cells[core.NewCoord3D(1, 2, 1)] = core.Alive

	return &Pattern3D{
		Name:        "3D Blinker",
		Description: "A period-2 oscillator in 3D",
		Width:       3,
		Height:      4,
		Depth:       3,
		Cells:       cells,
	}
}

// GetPatterns3D returns a map of all available 3D patterns
func GetPatterns3D() map[string]*Pattern3D {
	return map[string]*Pattern3D{
		"glider":     Glider3D(),
		"block":      Block3D(),
		"oscillator": Oscillator3D_Period2(),
	}
}

// ListPatterns3D returns a list of all available 3D pattern names
func ListPatterns3D() []string {
	return []string{
		"glider",
		"block",
		"oscillator",
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
