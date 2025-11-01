package legacy

import (
	"math/rand"
	"time"
)

// ClassicGrid represents the classic [][]int grid structure
type ClassicGrid struct {
	Width  int
	Height int
	Cells  [][]int
}

// NewClassicGrid creates a new classic grid
func NewClassicGrid(width, height int) *ClassicGrid {
	cells := make([][]int, height)
	for y := 0; y < height; y++ {
		cells[y] = make([]int, width)
	}
	return &ClassicGrid{
		Width:  width,
		Height: height,
		Cells:  cells,
	}
}

// Randomize fills the grid with random cells
func (g *ClassicGrid) Randomize() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			g.Cells[y][x] = r.Intn(2)
		}
	}
}

// CountNeighbors counts the number of alive neighbors around a cell
func (g *ClassicGrid) CountNeighbors(x, y int) int {
	count := 0

	// Check all 8 directions (Moore neighborhood)
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			// Skip the center cell itself
			if dx == 0 && dy == 0 {
				continue
			}

			// Calculate neighbor coordinates
			nx := x + dx
			ny := y + dy

			// Check boundaries
			if nx >= 0 && nx < g.Width && ny >= 0 && ny < g.Height {
				count += g.Cells[ny][nx]
			}
		}
	}

	return count
}

// Step executes one generation and returns a new grid
func (g *ClassicGrid) Step() *ClassicGrid {
	result := NewClassicGrid(g.Width, g.Height)

	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			neighbors := g.CountNeighbors(x, y)
			isAlive := g.Cells[y][x] == 1

			// Conway's Game of Life rules:
			// 1. Any live cell with 2 or 3 live neighbors survives
			// 2. Any dead cell with exactly 3 live neighbors becomes alive
			// 3. All other cells die or stay dead
			if isAlive && (neighbors == 2 || neighbors == 3) {
				result.Cells[y][x] = 1
			} else if !isAlive && neighbors == 3 {
				result.Cells[y][x] = 1
			} else {
				result.Cells[y][x] = 0
			}
		}
	}

	return result
}

// Clone creates a deep copy of the grid
func (g *ClassicGrid) Clone() *ClassicGrid {
	clone := NewClassicGrid(g.Width, g.Height)
	for y := 0; y < g.Height; y++ {
		copy(clone.Cells[y], g.Cells[y])
	}
	return clone
}

// CountLiving returns the number of living cells
func (g *ClassicGrid) CountLiving() int {
	count := 0
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			count += g.Cells[y][x]
		}
	}
	return count
}

// Set sets the value of a cell
func (g *ClassicGrid) Set(x, y, value int) {
	if x >= 0 && x < g.Width && y >= 0 && y < g.Height {
		g.Cells[y][x] = value
	}
}

// Get gets the value of a cell
func (g *ClassicGrid) Get(x, y int) int {
	if x >= 0 && x < g.Width && y >= 0 && y < g.Height {
		return g.Cells[y][x]
	}
	return 0
}
