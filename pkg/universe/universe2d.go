package universe

import (
	"golife/pkg/core"
	"math/rand"
	"time"
)

// Universe2D implements the Universe interface for 2D cellular automaton
type Universe2D struct {
	width, height int
	cells         []core.CellState // Flat array for cache locality: [y*width + x]
	nextCells     []core.CellState // Double buffering
	rule          core.Rule
	neighborhood  core.NeighborhoodType
}

// New2D creates a new 2D universe with the given dimensions and rule
func New2D(width, height int, rule core.Rule) *Universe2D {
	size := width * height
	return &Universe2D{
		width:        width,
		height:       height,
		cells:        make([]core.CellState, size),
		nextCells:    make([]core.CellState, size),
		rule:         rule,
		neighborhood: core.Moore, // Default to Moore neighborhood (8 neighbors)
	}
}

// Dimension returns the dimensionality (2D)
func (u *Universe2D) Dimension() core.Dimension {
	return core.Dim2D
}

// Get returns the state of a cell at the given coordinate
func (u *Universe2D) Get(coord core.Coord) core.CellState {
	if coord.X < 0 || coord.X >= u.width || coord.Y < 0 || coord.Y >= u.height {
		return core.Dead
	}
	return u.cells[coord.Y*u.width+coord.X]
}

// Set sets the state of a cell at the given coordinate
func (u *Universe2D) Set(coord core.Coord, state core.CellState) {
	if coord.X < 0 || coord.X >= u.width || coord.Y < 0 || coord.Y >= u.height {
		return
	}
	u.cells[coord.Y*u.width+coord.X] = state
}

// Size returns the dimensions of the universe
func (u *Universe2D) Size() core.Coord {
	return core.NewCoord2D(u.width, u.height)
}

// countNeighbors counts the number of alive neighbors around a cell
func (u *Universe2D) countNeighbors(x, y int) int {
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
			if nx >= 0 && nx < u.width && ny >= 0 && ny < u.height {
				if u.cells[ny*u.width+nx] > core.Dead {
					count++
				}
			}
		}
	}

	return count
}

// Step executes one generation
func (u *Universe2D) Step() {
	for y := 0; y < u.height; y++ {
		for x := 0; x < u.width; x++ {
			idx := y*u.width + x
			neighbors := u.countNeighbors(x, y)
			currentState := u.cells[idx]

			// Apply Conway's Game of Life rules
			if currentState == core.Dead {
				// Dead cell: check birth condition
				if u.rule.ShouldBirth(neighbors) {
					u.nextCells[idx] = core.Alive
				} else {
					u.nextCells[idx] = core.Dead
				}
			} else {
				// Living cell: check survival condition
				if u.rule.ShouldSurvive(neighbors, currentState) {
					u.nextCells[idx] = core.Alive
				} else {
					u.nextCells[idx] = core.Dead
				}
			}
		}
	}

	// Swap buffers
	u.cells, u.nextCells = u.nextCells, u.cells
}

// Clone creates a deep copy of the universe
func (u *Universe2D) Clone() core.Universe {
	clone := New2D(u.width, u.height, u.rule)
	copy(clone.cells, u.cells)
	return clone
}

// Clear resets all cells to dead state
func (u *Universe2D) Clear() {
	for i := range u.cells {
		u.cells[i] = core.Dead
	}
}

// CountLiving returns the number of living cells
func (u *Universe2D) CountLiving() int {
	count := 0
	for _, cell := range u.cells {
		if cell > core.Dead {
			count++
		}
	}
	return count
}

// Randomize fills the universe with random cells
func (u *Universe2D) Randomize() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range u.cells {
		if r.Intn(2) == 1 {
			u.cells[i] = core.Alive
		} else {
			u.cells[i] = core.Dead
		}
	}
}

// GetCells returns the internal cell array (for compatibility with legacy code)
func (u *Universe2D) GetCells() [][]int {
	result := make([][]int, u.height)
	for y := 0; y < u.height; y++ {
		result[y] = make([]int, u.width)
		for x := 0; x < u.width; x++ {
			if u.cells[y*u.width+x] > core.Dead {
				result[y][x] = 1
			} else {
				result[y][x] = 0
			}
		}
	}
	return result
}

// SetCells sets cells from a 2D int array (for compatibility with legacy code)
func (u *Universe2D) SetCells(cells [][]int) {
	for y := 0; y < u.height && y < len(cells); y++ {
		for x := 0; x < u.width && x < len(cells[y]); x++ {
			if cells[y][x] == 1 {
				u.cells[y*u.width+x] = core.Alive
			} else {
				u.cells[y*u.width+x] = core.Dead
			}
		}
	}
}
