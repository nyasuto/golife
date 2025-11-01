package core

// Dimension represents the dimensionality of the universe
type Dimension int

const (
	Dim2D  Dimension = 2
	Dim25D Dimension = 25 // Multiple 2D layers
	Dim3D  Dimension = 3
	Dim4D  Dimension = 4
)

// CellState represents the state of a cell
// 0 = dead, 1-255 = alive with gradual states for Generations rules
type CellState uint8

const (
	Dead  CellState = 0
	Alive CellState = 255
)

// Coord represents a dimension-agnostic coordinate
// For 2D, only X and Y are used
// For 3D, X, Y, Z are used
// For 4D, all fields are used
type Coord struct {
	X, Y, Z, W int
}

// NewCoord2D creates a 2D coordinate
func NewCoord2D(x, y int) Coord {
	return Coord{X: x, Y: y, Z: 0, W: 0}
}

// NewCoord3D creates a 3D coordinate
func NewCoord3D(x, y, z int) Coord {
	return Coord{X: x, Y: y, Z: z, W: 0}
}

// NewCoord4D creates a 4D coordinate
func NewCoord4D(x, y, z, w int) Coord {
	return Coord{X: x, Y: y, Z: z, W: w}
}

// Universe interface supports all dimensions
type Universe interface {
	// Dimension returns the dimensionality of this universe
	Dimension() Dimension

	// Get returns the state of a cell at the given coordinate
	Get(coord Coord) CellState

	// Set sets the state of a cell at the given coordinate
	Set(coord Coord, state CellState)

	// Step executes one generation
	Step()

	// Size returns the dimensions of the universe
	Size() Coord

	// Clone creates a deep copy of the universe
	Clone() Universe

	// Clear resets all cells to dead state
	Clear()

	// CountLiving returns the number of living cells
	CountLiving() int
}

// Rule defines birth/survival conditions
type Rule interface {
	// Name returns the name of this rule
	Name() string

	// ShouldBirth determines if a dead cell should become alive
	ShouldBirth(neighborCount int) bool

	// ShouldSurvive determines if a live cell should stay alive
	ShouldSurvive(neighborCount int, currentState CellState) bool

	// NeighborWeight returns the weight for a neighbor at given distance
	// Used for distance-decay rules. Default implementation returns 1.0 for all distances.
	NeighborWeight(distance float64) float64
}

// NeighborhoodType defines the type of neighborhood calculation
type NeighborhoodType int

const (
	// Moore neighborhood includes all adjacent cells (8 in 2D, 26 in 3D, 80 in 4D)
	Moore NeighborhoodType = iota

	// VonNeumann neighborhood includes only face-adjacent cells (4 in 2D, 6 in 3D, 8 in 4D)
	VonNeumann

	// Custom allows user-defined neighborhood offsets
	Custom
)

// String returns the string representation of the neighborhood type
func (n NeighborhoodType) String() string {
	switch n {
	case Moore:
		return "Moore"
	case VonNeumann:
		return "VonNeumann"
	case Custom:
		return "Custom"
	default:
		return "Unknown"
	}
}
