package legacy

import (
	"golife/pkg/core"
	"golife/pkg/rules"
	"golife/pkg/universe"
)

// ToUniverse converts a ClassicGrid to a Universe2D
func (g *ClassicGrid) ToUniverse() core.Universe {
	u := universe.New2D(g.Width, g.Height, rules.ConwayRule{})

	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			if g.Cells[y][x] == 1 {
				u.Set(core.NewCoord2D(x, y), core.Alive)
			}
		}
	}

	return u
}

// NewFromUniverse creates a ClassicGrid from a Universe
func NewFromUniverse(u core.Universe) *ClassicGrid {
	size := u.Size()
	grid := NewClassicGrid(size.X, size.Y)

	for y := 0; y < size.Y; y++ {
		for x := 0; x < size.X; x++ {
			state := u.Get(core.NewCoord2D(x, y))
			if state > core.Dead {
				grid.Cells[y][x] = 1
			} else {
				grid.Cells[y][x] = 0
			}
		}
	}

	return grid
}

// FromSlice creates a ClassicGrid from a [][]int slice
func FromSlice(cells [][]int) *ClassicGrid {
	if len(cells) == 0 {
		return NewClassicGrid(0, 0)
	}

	height := len(cells)
	width := len(cells[0])
	grid := NewClassicGrid(width, height)

	for y := 0; y < height; y++ {
		copy(grid.Cells[y], cells[y])
	}

	return grid
}

// ToSlice converts a ClassicGrid to [][]int slice
func (g *ClassicGrid) ToSlice() [][]int {
	result := make([][]int, g.Height)
	for y := 0; y < g.Height; y++ {
		result[y] = make([]int, g.Width)
		copy(result[y], g.Cells[y])
	}
	return result
}
