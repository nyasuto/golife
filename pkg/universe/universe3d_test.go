package universe

import (
	"golife/pkg/core"
	"golife/pkg/rules"
	"testing"
)

func TestNew3D(t *testing.T) {
	u := New3D(10, 10, 10, rules.ConwayRule{})

	if u == nil {
		t.Fatal("New3D should not return nil")
	}

	if u.width != 10 || u.height != 10 || u.depth != 10 {
		t.Errorf("Expected dimensions (10,10,10), got (%d,%d,%d)", u.width, u.height, u.depth)
	}

	size := u.Size()
	if size.X != 10 || size.Y != 10 || size.Z != 10 {
		t.Errorf("Size() returned incorrect dimensions: %v", size)
	}

	// Should have 26 neighbor offsets
	if len(u.neighborOffsets) != 26 {
		t.Errorf("Expected 26 neighbor offsets, got %d", len(u.neighborOffsets))
	}
}

func TestUniverse3D_Dimension(t *testing.T) {
	u := New3D(10, 10, 10, rules.ConwayRule{})

	if u.Dimension() != core.Dim3D {
		t.Errorf("Expected Dim3D, got %v", u.Dimension())
	}
}

func TestUniverse3D_GetSet(t *testing.T) {
	u := New3D(10, 10, 10, rules.ConwayRule{})

	// Test setting and getting a cell
	coord := core.NewCoord3D(5, 5, 5)
	u.Set(coord, core.Alive)

	state := u.Get(coord)
	if state != core.Alive {
		t.Errorf("Expected Alive, got %v", state)
	}

	// Test multiple cells
	coords := []core.Coord{
		core.NewCoord3D(0, 0, 0),
		core.NewCoord3D(9, 9, 9),
		core.NewCoord3D(3, 7, 2),
	}

	for _, c := range coords {
		u.Set(c, core.Alive)
		if u.Get(c) != core.Alive {
			t.Errorf("Cell at %v should be alive", c)
		}
	}
}

func TestUniverse3D_BoundsChecking(t *testing.T) {
	u := New3D(10, 10, 10, rules.ConwayRule{})

	// Out of bounds should return Dead
	testCases := []core.Coord{
		core.NewCoord3D(-1, 5, 5),
		core.NewCoord3D(5, -1, 5),
		core.NewCoord3D(5, 5, -1),
		core.NewCoord3D(10, 5, 5),
		core.NewCoord3D(5, 10, 5),
		core.NewCoord3D(5, 5, 10),
		core.NewCoord3D(100, 100, 100),
	}

	for _, coord := range testCases {
		state := u.Get(coord)
		if state != core.Dead {
			t.Errorf("Out of bounds coord %v should return Dead, got %v", coord, state)
		}

		// Setting out of bounds should not crash
		u.Set(coord, core.Alive)
	}
}

func TestUniverse3D_CountNeighbors(t *testing.T) {
	u := New3D(10, 10, 10, rules.ConwayRule{})

	// Test center cell with all 26 neighbors alive
	// Set all 26 neighbors
	for dz := -1; dz <= 1; dz++ {
		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				if dx == 0 && dy == 0 && dz == 0 {
					continue
				}
				coord := core.NewCoord3D(5+dx, 5+dy, 5+dz)
				u.Set(coord, core.Alive)
			}
		}
	}

	count := u.countNeighbors(5, 5, 5)
	if count != 26 {
		t.Errorf("Expected 26 neighbors, got %d", count)
	}

	// Test corner cell (should have fewer neighbors due to bounds)
	u.Clear()
	u.Set(core.NewCoord3D(1, 0, 0), core.Alive)
	u.Set(core.NewCoord3D(0, 1, 0), core.Alive)
	u.Set(core.NewCoord3D(0, 0, 1), core.Alive)

	count = u.countNeighbors(0, 0, 0)
	if count != 3 {
		t.Errorf("Corner cell should have 3 neighbors, got %d", count)
	}
}

func TestUniverse3D_CountNeighbors_Partial(t *testing.T) {
	u := New3D(10, 10, 10, rules.ConwayRule{})

	// Test with specific neighbor pattern
	// Set 6 face neighbors
	u.Set(core.NewCoord3D(6, 5, 5), core.Alive) // +X
	u.Set(core.NewCoord3D(4, 5, 5), core.Alive) // -X
	u.Set(core.NewCoord3D(5, 6, 5), core.Alive) // +Y
	u.Set(core.NewCoord3D(5, 4, 5), core.Alive) // -Y
	u.Set(core.NewCoord3D(5, 5, 6), core.Alive) // +Z
	u.Set(core.NewCoord3D(5, 5, 4), core.Alive) // -Z

	count := u.countNeighbors(5, 5, 5)
	if count != 6 {
		t.Errorf("Expected 6 face neighbors, got %d", count)
	}

	// Add 12 edge neighbors
	u.Set(core.NewCoord3D(6, 6, 5), core.Alive)
	u.Set(core.NewCoord3D(6, 4, 5), core.Alive)
	u.Set(core.NewCoord3D(4, 6, 5), core.Alive)
	u.Set(core.NewCoord3D(4, 4, 5), core.Alive)
	u.Set(core.NewCoord3D(6, 5, 6), core.Alive)
	u.Set(core.NewCoord3D(6, 5, 4), core.Alive)
	u.Set(core.NewCoord3D(4, 5, 6), core.Alive)
	u.Set(core.NewCoord3D(4, 5, 4), core.Alive)
	u.Set(core.NewCoord3D(5, 6, 6), core.Alive)
	u.Set(core.NewCoord3D(5, 6, 4), core.Alive)
	u.Set(core.NewCoord3D(5, 4, 6), core.Alive)
	u.Set(core.NewCoord3D(5, 4, 4), core.Alive)

	count = u.countNeighbors(5, 5, 5)
	if count != 18 {
		t.Errorf("Expected 18 neighbors (6 face + 12 edge), got %d", count)
	}
}

func TestUniverse3D_Step(t *testing.T) {
	u := New3D(10, 10, 10, rules.ConwayRule{})

	// Create a simple pattern
	u.Set(core.NewCoord3D(5, 5, 5), core.Alive)
	u.Set(core.NewCoord3D(6, 5, 5), core.Alive)
	u.Set(core.NewCoord3D(5, 6, 5), core.Alive)

	initialCount := u.CountLiving()
	if initialCount != 3 {
		t.Errorf("Expected 3 initial living cells, got %d", initialCount)
	}

	// Execute one step
	u.Step()

	// Count should change based on rule
	afterCount := u.CountLiving()
	t.Logf("Living cells after step: %d", afterCount)

	// The exact count depends on the rule, but it should be deterministic
	// For Conway's rule in 3D, these 3 cells would likely die (too few neighbors)
}

func TestUniverse3D_Clear(t *testing.T) {
	u := New3D(10, 10, 10, rules.ConwayRule{})

	// Set some cells
	u.Set(core.NewCoord3D(5, 5, 5), core.Alive)
	u.Set(core.NewCoord3D(6, 5, 5), core.Alive)
	u.Set(core.NewCoord3D(5, 6, 5), core.Alive)

	if u.CountLiving() == 0 {
		t.Error("Should have living cells before clear")
	}

	u.Clear()

	if u.CountLiving() != 0 {
		t.Errorf("Expected 0 living cells after clear, got %d", u.CountLiving())
	}
}

func TestUniverse3D_CountLiving(t *testing.T) {
	u := New3D(10, 10, 10, rules.ConwayRule{})

	if u.CountLiving() != 0 {
		t.Error("New universe should have 0 living cells")
	}

	// Add cells
	cells := []core.Coord{
		core.NewCoord3D(0, 0, 0),
		core.NewCoord3D(5, 5, 5),
		core.NewCoord3D(9, 9, 9),
		core.NewCoord3D(2, 3, 4),
		core.NewCoord3D(7, 1, 8),
	}

	for _, c := range cells {
		u.Set(c, core.Alive)
	}

	count := u.CountLiving()
	if count != len(cells) {
		t.Errorf("Expected %d living cells, got %d", len(cells), count)
	}
}

func TestUniverse3D_GetSlice(t *testing.T) {
	u := New3D(5, 5, 5, rules.ConwayRule{})

	// Set some cells in layer z=2
	u.Set(core.NewCoord3D(1, 1, 2), core.Alive)
	u.Set(core.NewCoord3D(2, 2, 2), core.Alive)
	u.Set(core.NewCoord3D(3, 3, 2), core.Alive)

	// Get slice at z=2
	slice := u.GetSlice(2)

	if slice == nil {
		t.Fatal("GetSlice should not return nil for valid z")
	}

	if len(slice) != 5 {
		t.Errorf("Expected 5 rows, got %d", len(slice))
	}

	if len(slice[0]) != 5 {
		t.Errorf("Expected 5 columns, got %d", len(slice[0]))
	}

	// Check specific cells
	if slice[1][1] != core.Alive {
		t.Error("Cell [1][1] should be alive")
	}
	if slice[2][2] != core.Alive {
		t.Error("Cell [2][2] should be alive")
	}
	if slice[3][3] != core.Alive {
		t.Error("Cell [3][3] should be alive")
	}
	if slice[0][0] != core.Dead {
		t.Error("Cell [0][0] should be dead")
	}

	// Test invalid z
	invalidSlice := u.GetSlice(-1)
	if invalidSlice != nil {
		t.Error("GetSlice with invalid z should return nil")
	}

	invalidSlice = u.GetSlice(10)
	if invalidSlice != nil {
		t.Error("GetSlice with z >= depth should return nil")
	}
}

func TestUniverse3D_MemoryUsage(t *testing.T) {
	// Test 64^3 grid (mentioned in issue)
	u := New3D(64, 64, 64, rules.ConwayRule{})

	totalCells := 64 * 64 * 64
	expectedSize := totalCells * 2 // cells + nextCells

	if len(u.cells) != totalCells {
		t.Errorf("Expected %d cells, got %d", totalCells, len(u.cells))
	}

	if len(u.nextCells) != totalCells {
		t.Errorf("Expected %d nextCells, got %d", totalCells, len(u.nextCells))
	}

	t.Logf("Memory usage for 64³ grid: ~%d bytes (cells array)", totalCells)
	t.Logf("Total with double buffering: ~%d bytes", expectedSize)
}

func BenchmarkUniverse3D_Step_32x32x32(b *testing.B) {
	u := New3D(32, 32, 32, rules.ConwayRule{})

	// Create random pattern
	for i := 0; i < 1000; i++ {
		u.Set(core.NewCoord3D(i%32, (i/32)%32, (i/1024)%32), core.Alive)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u.Step()
	}
}

func BenchmarkUniverse3D_Step_64x64x64(b *testing.B) {
	u := New3D(64, 64, 64, rules.ConwayRule{})

	// Create random pattern
	for i := 0; i < 5000; i++ {
		u.Set(core.NewCoord3D(i%64, (i/64)%64, (i/4096)%64), core.Alive)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u.Step()
	}
}

func BenchmarkUniverse3D_CountNeighbors(b *testing.B) {
	u := New3D(64, 64, 64, rules.ConwayRule{})

	// Set up a dense region
	for z := 30; z < 34; z++ {
		for y := 30; y < 34; y++ {
			for x := 30; x < 34; x++ {
				u.Set(core.NewCoord3D(x, y, z), core.Alive)
			}
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u.countNeighbors(32, 32, 32)
	}
}
