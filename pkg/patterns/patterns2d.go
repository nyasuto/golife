package patterns

import (
	"golife/pkg/core"
	"golife/pkg/universe"
)

// Pattern2D represents a 2D pattern definition
type Pattern2D struct {
	Name        string
	Description string
	Width       int
	Height      int
	Cells       [][]core.CellState
}

// LoadIntoUniverse loads this pattern into a universe at the given offset
func (p *Pattern2D) LoadIntoUniverse(u *universe.Universe2D, offsetX, offsetY int) {
	for y := 0; y < len(p.Cells) && y < p.Height; y++ {
		for x := 0; x < len(p.Cells[y]) && x < p.Width; x++ {
			if p.Cells[y][x] != core.Dead {
				coord := core.NewCoord2D(offsetX+x, offsetY+y)
				u.Set(coord, p.Cells[y][x])
			}
		}
	}
}

// Glider returns the classic glider pattern
func Glider() Pattern2D {
	const (
		O = core.Dead
		X = core.Alive
	)
	return Pattern2D{
		Name:        "Glider",
		Description: "A small pattern that moves diagonally",
		Width:       3,
		Height:      3,
		Cells: [][]core.CellState{
			{O, X, O},
			{O, O, X},
			{X, X, X},
		},
	}
}

// Blinker returns a period-2 oscillator
func Blinker() Pattern2D {
	const X = core.Alive
	return Pattern2D{
		Name:        "Blinker",
		Description: "A period-2 oscillator",
		Width:       3,
		Height:      1,
		Cells: [][]core.CellState{
			{X, X, X},
		},
	}
}

// Toad returns a period-2 oscillator
func Toad() Pattern2D {
	const (
		O = core.Dead
		X = core.Alive
	)
	return Pattern2D{
		Name:        "Toad",
		Description: "A period-2 oscillator",
		Width:       4,
		Height:      2,
		Cells: [][]core.CellState{
			{O, X, X, X},
			{X, X, X, O},
		},
	}
}

// Beacon returns a period-2 oscillator
func Beacon() Pattern2D {
	const (
		O = core.Dead
		X = core.Alive
	)
	return Pattern2D{
		Name:        "Beacon",
		Description: "A period-2 oscillator",
		Width:       4,
		Height:      4,
		Cells: [][]core.CellState{
			{X, X, O, O},
			{X, X, O, O},
			{O, O, X, X},
			{O, O, X, X},
		},
	}
}

// Pulsar returns a period-3 oscillator
func Pulsar() Pattern2D {
	const (
		O = core.Dead
		X = core.Alive
	)
	return Pattern2D{
		Name:        "Pulsar",
		Description: "A period-3 oscillator",
		Width:       13,
		Height:      13,
		Cells: [][]core.CellState{
			{O, O, X, X, X, O, O, O, X, X, X, O, O},
			{O, O, O, O, O, O, O, O, O, O, O, O, O},
			{X, O, O, O, O, X, O, X, O, O, O, O, X},
			{X, O, O, O, O, X, O, X, O, O, O, O, X},
			{X, O, O, O, O, X, O, X, O, O, O, O, X},
			{O, O, X, X, X, O, O, O, X, X, X, O, O},
			{O, O, O, O, O, O, O, O, O, O, O, O, O},
			{O, O, X, X, X, O, O, O, X, X, X, O, O},
			{X, O, O, O, O, X, O, X, O, O, O, O, X},
			{X, O, O, O, O, X, O, X, O, O, O, O, X},
			{X, O, O, O, O, X, O, X, O, O, O, O, X},
			{O, O, O, O, O, O, O, O, O, O, O, O, O},
			{O, O, X, X, X, O, O, O, X, X, X, O, O},
		},
	}
}

// GliderGun returns the Gosper Glider Gun
func GliderGun() Pattern2D {
	const (
		O = core.Dead
		X = core.Alive
	)
	return Pattern2D{
		Name:        "Gosper Glider Gun",
		Description: "A pattern that generates gliders",
		Width:       36,
		Height:      9,
		Cells: [][]core.CellState{
			{O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, X, O, O, O, O, O, O, O, O, O, O, O},
			{O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, X, O, X, O, O, O, O, O, O, O, O, O, O, O},
			{O, O, O, O, O, O, O, O, O, O, O, O, X, X, O, O, O, O, O, O, X, X, O, O, O, O, O, O, O, O, O, O, O, O, X, X},
			{O, O, O, O, O, O, O, O, O, O, O, X, O, O, O, X, O, O, O, O, X, X, O, O, O, O, O, O, O, O, O, O, O, O, X, X},
			{X, X, O, O, O, O, O, O, O, O, X, O, O, O, O, O, X, O, O, O, X, X, O, O, O, O, O, O, O, O, O, O, O, O, O, O},
			{X, X, O, O, O, O, O, O, O, O, X, O, O, O, X, O, X, X, O, O, O, O, X, O, X, O, O, O, O, O, O, O, O, O, O, O},
			{O, O, O, O, O, O, O, O, O, O, X, O, O, O, O, O, X, O, O, O, O, O, O, O, X, O, O, O, O, O, O, O, O, O, O, O},
			{O, O, O, O, O, O, O, O, O, O, O, X, O, O, O, X, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O},
			{O, O, O, O, O, O, O, O, O, O, O, O, X, X, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O, O},
		},
	}
}

// Block returns a still life (2x2 square)
func Block() Pattern2D {
	const X = core.Alive
	return Pattern2D{
		Name:        "Block",
		Description: "A still life (stable pattern)",
		Width:       2,
		Height:      2,
		Cells: [][]core.CellState{
			{X, X},
			{X, X},
		},
	}
}

// AllPatterns returns a map of all available 2D patterns
func AllPatterns() map[string]Pattern2D {
	return map[string]Pattern2D{
		"glider":     Glider(),
		"blinker":    Blinker(),
		"toad":       Toad(),
		"beacon":     Beacon(),
		"pulsar":     Pulsar(),
		"glider-gun": GliderGun(),
		"block":      Block(),
	}
}
