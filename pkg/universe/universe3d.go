package universe

import (
	"golife/pkg/core"
	"runtime"
	"sync"
)

// Universe3D represents a true 3D universe with full 26-neighbor interaction
type Universe3D struct {
	width, height, depth int
	cells                []core.CellState // Flat array: [z*height*width + y*width + x]
	nextCells            []core.CellState
	rule                 core.Rule
	neighborOffsets      []int // Pre-computed neighbor offsets for performance
}

// New3D creates a new 3D universe with the given dimensions and rule
func New3D(width, height, depth int, rule core.Rule) *Universe3D {
	size := width * height * depth
	u := &Universe3D{
		width:     width,
		height:    height,
		depth:     depth,
		cells:     make([]core.CellState, size),
		nextCells: make([]core.CellState, size),
		rule:      rule,
	}

	// Pre-compute 26-neighbor offsets for performance
	u.precomputeNeighborOffsets()

	return u
}

// precomputeNeighborOffsets pre-computes the flat array offsets for all 26 neighbors
func (u *Universe3D) precomputeNeighborOffsets() {
	u.neighborOffsets = make([]int, 0, 26)

	// All 26 neighbors (3x3x3 cube minus center)
	for dz := -1; dz <= 1; dz++ {
		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				if dx == 0 && dy == 0 && dz == 0 {
					continue // Skip center cell
				}
				offset := dz*u.height*u.width + dy*u.width + dx
				u.neighborOffsets = append(u.neighborOffsets, offset)
			}
		}
	}
}

// Dimension returns the dimensionality (3D)
func (u *Universe3D) Dimension() core.Dimension {
	return core.Dim3D
}

// coordToIndex converts 3D coordinates to flat array index
func (u *Universe3D) coordToIndex(coord core.Coord) int {
	return coord.Z*u.height*u.width + coord.Y*u.width + coord.X
}

// isValid checks if coordinates are within bounds
func (u *Universe3D) isValid(x, y, z int) bool {
	return x >= 0 && x < u.width &&
		y >= 0 && y < u.height &&
		z >= 0 && z < u.depth
}

// Get returns the state of a cell at the given coordinate
func (u *Universe3D) Get(coord core.Coord) core.CellState {
	if !u.isValid(coord.X, coord.Y, coord.Z) {
		return core.Dead
	}
	return u.cells[u.coordToIndex(coord)]
}

// Set sets the state of a cell at the given coordinate
func (u *Universe3D) Set(coord core.Coord, state core.CellState) {
	if u.isValid(coord.X, coord.Y, coord.Z) {
		u.cells[u.coordToIndex(coord)] = state
	}
}

// Size returns the dimensions of the universe
func (u *Universe3D) Size() core.Coord {
	return core.NewCoord3D(u.width, u.height, u.depth)
}

// countNeighbors counts living neighbors for a cell (all 26 neighbors in 3D)
// This is the boundary-safe version with explicit coordinate checks
func (u *Universe3D) countNeighbors(x, y, z int) int {
	count := 0

	// Check all 26 neighbors (3x3x3 cube minus center)
	for dz := -1; dz <= 1; dz++ {
		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				// Skip center cell
				if dx == 0 && dy == 0 && dz == 0 {
					continue
				}

				nx := x + dx
				ny := y + dy
				nz := z + dz

				// Check bounds
				if nx < 0 || nx >= u.width || ny < 0 || ny >= u.height || nz < 0 || nz >= u.depth {
					continue
				}

				// Count if neighbor is alive
				idx := nz*u.height*u.width + ny*u.width + nx
				if u.cells[idx] != core.Dead {
					count++
				}
			}
		}
	}

	return count
}

// countNeighborsInterior counts neighbors for interior cells (no boundary check needed)
func (u *Universe3D) countNeighborsInterior(idx int) int {
	count := 0
	for _, offset := range u.neighborOffsets {
		if u.cells[idx+offset] != core.Dead {
			count++
		}
	}
	return count
}

// Step executes one generation using the rule
func (u *Universe3D) Step() {
	// Optimize by separating interior cells (no boundary check) from boundary cells
	// Interior region: cells that have all 26 neighbors within bounds
	interiorStartX, interiorEndX := 1, u.width-1
	interiorStartY, interiorEndY := 1, u.height-1
	interiorStartZ, interiorEndZ := 1, u.depth-1

	// Check if we have an interior region
	hasInterior := u.width > 2 && u.height > 2 && u.depth > 2

	if hasInterior {
		// Process interior cells (fast path - no boundary checks)
		for z := interiorStartZ; z < interiorEndZ; z++ {
			for y := interiorStartY; y < interiorEndY; y++ {
				for x := interiorStartX; x < interiorEndX; x++ {
					idx := z*u.height*u.width + y*u.width + x
					neighbors := u.countNeighborsInterior(idx)
					currentState := u.cells[idx]

					// Apply rule
					if currentState == core.Dead {
						if u.rule.ShouldBirth(neighbors) {
							u.nextCells[idx] = core.Alive
						} else {
							u.nextCells[idx] = core.Dead
						}
					} else {
						if u.rule.ShouldSurvive(neighbors, currentState) {
							u.nextCells[idx] = core.Alive
						} else {
							u.nextCells[idx] = core.Dead
						}
					}
				}
			}
		}
	}

	// Process boundary cells (with bounds checking)
	for z := 0; z < u.depth; z++ {
		for y := 0; y < u.height; y++ {
			for x := 0; x < u.width; x++ {
				// Skip interior cells (already processed)
				if hasInterior && x >= interiorStartX && x < interiorEndX &&
					y >= interiorStartY && y < interiorEndY &&
					z >= interiorStartZ && z < interiorEndZ {
					continue
				}

				idx := z*u.height*u.width + y*u.width + x
				neighbors := u.countNeighbors(x, y, z)
				currentState := u.cells[idx]

				// Apply rule
				if currentState == core.Dead {
					if u.rule.ShouldBirth(neighbors) {
						u.nextCells[idx] = core.Alive
					} else {
						u.nextCells[idx] = core.Dead
					}
				} else {
					if u.rule.ShouldSurvive(neighbors, currentState) {
						u.nextCells[idx] = core.Alive
					} else {
						u.nextCells[idx] = core.Dead
					}
				}
			}
		}
	}

	// Swap buffers
	u.cells, u.nextCells = u.nextCells, u.cells
}

// StepParallel executes one generation using parallel processing
// The grid is divided into Z-axis slices and processed concurrently
func (u *Universe3D) StepParallel() {
	numWorkers := runtime.NumCPU()
	if numWorkers > u.depth {
		numWorkers = u.depth
	}

	// Calculate slice size for each worker
	sliceSize := u.depth / numWorkers
	remainder := u.depth % numWorkers

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	// Launch workers to process Z-axis slices in parallel
	for workerID := 0; workerID < numWorkers; workerID++ {
		zStart := workerID * sliceSize
		zEnd := zStart + sliceSize
		if workerID == numWorkers-1 {
			zEnd += remainder // Last worker handles remainder
		}

		go func(zStart, zEnd int) {
			defer wg.Done()
			u.processZSlice(zStart, zEnd)
		}(zStart, zEnd)
	}

	wg.Wait()

	// Swap buffers (single-threaded)
	u.cells, u.nextCells = u.nextCells, u.cells
}

// processZSlice processes a range of Z layers [zStart, zEnd)
func (u *Universe3D) processZSlice(zStart, zEnd int) {
	// Determine interior region boundaries
	interiorStartX, interiorEndX := 1, u.width-1
	interiorStartY, interiorEndY := 1, u.height-1
	hasInterior := u.width > 2 && u.height > 2

	for z := zStart; z < zEnd; z++ {
		// Process interior cells (fast path)
		if hasInterior && z > 0 && z < u.depth-1 {
			for y := interiorStartY; y < interiorEndY; y++ {
				for x := interiorStartX; x < interiorEndX; x++ {
					idx := z*u.height*u.width + y*u.width + x
					neighbors := u.countNeighborsInterior(idx)
					currentState := u.cells[idx]

					if currentState == core.Dead {
						if u.rule.ShouldBirth(neighbors) {
							u.nextCells[idx] = core.Alive
						} else {
							u.nextCells[idx] = core.Dead
						}
					} else {
						if u.rule.ShouldSurvive(neighbors, currentState) {
							u.nextCells[idx] = core.Alive
						} else {
							u.nextCells[idx] = core.Dead
						}
					}
				}
			}
		}

		// Process boundary cells (safe path)
		for y := 0; y < u.height; y++ {
			for x := 0; x < u.width; x++ {
				// Skip interior cells if already processed
				if hasInterior && z > 0 && z < u.depth-1 &&
					x >= interiorStartX && x < interiorEndX &&
					y >= interiorStartY && y < interiorEndY {
					continue
				}

				idx := z*u.height*u.width + y*u.width + x
				neighbors := u.countNeighbors(x, y, z)
				currentState := u.cells[idx]

				if currentState == core.Dead {
					if u.rule.ShouldBirth(neighbors) {
						u.nextCells[idx] = core.Alive
					} else {
						u.nextCells[idx] = core.Dead
					}
				} else {
					if u.rule.ShouldSurvive(neighbors, currentState) {
						u.nextCells[idx] = core.Alive
					} else {
						u.nextCells[idx] = core.Dead
					}
				}
			}
		}
	}
}

// Clear sets all cells to dead
func (u *Universe3D) Clear() {
	for i := range u.cells {
		u.cells[i] = core.Dead
	}
}

// CountLiving returns the number of living cells
func (u *Universe3D) CountLiving() int {
	count := 0
	for _, state := range u.cells {
		if state != core.Dead {
			count++
		}
	}
	return count
}

// GetSlice returns a 2D slice at the given Z level
func (u *Universe3D) GetSlice(z int) [][]core.CellState {
	if z < 0 || z >= u.depth {
		return nil
	}

	slice := make([][]core.CellState, u.height)
	for y := 0; y < u.height; y++ {
		slice[y] = make([]core.CellState, u.width)
		for x := 0; x < u.width; x++ {
			idx := z*u.height*u.width + y*u.width + x
			slice[y][x] = u.cells[idx]
		}
	}
	return slice
}

// Clone creates a deep copy of the universe
func (u *Universe3D) Clone() core.Universe {
	clone := &Universe3D{
		width:  u.width,
		height: u.height,
		depth:  u.depth,
		rule:   u.rule,
	}

	// Copy cells
	clone.cells = make([]core.CellState, len(u.cells))
	copy(clone.cells, u.cells)

	// Copy nextCells
	clone.nextCells = make([]core.CellState, len(u.nextCells))
	copy(clone.nextCells, u.nextCells)

	// Pre-compute neighbor offsets
	clone.precomputeNeighborOffsets()

	return clone
}
