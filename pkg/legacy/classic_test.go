package legacy

import "testing"

func TestNewClassicGrid(t *testing.T) {
	g := NewClassicGrid(10, 10)

	if g.Width != 10 || g.Height != 10 {
		t.Errorf("Expected dimensions (10,10), got (%d,%d)", g.Width, g.Height)
	}

	if len(g.Cells) != 10 || len(g.Cells[0]) != 10 {
		t.Error("Cells array not initialized correctly")
	}
}

func TestClassicGrid_GetSet(t *testing.T) {
	g := NewClassicGrid(5, 5)

	// Test Set and Get
	g.Set(2, 3, 1)
	if g.Get(2, 3) != 1 {
		t.Errorf("Expected 1, got %d", g.Get(2, 3))
	}

	// Test boundary
	if g.Get(10, 10) != 0 {
		t.Error("Out of bounds should return 0")
	}
}

func TestClassicGrid_Step_Glider(t *testing.T) {
	g := NewClassicGrid(10, 10)

	// Set up a glider pattern
	g.Set(1, 0, 1)
	g.Set(2, 1, 1)
	g.Set(0, 2, 1)
	g.Set(1, 2, 1)
	g.Set(2, 2, 1)

	initialCount := g.CountLiving()
	if initialCount != 5 {
		t.Errorf("Expected 5 living cells, got %d", initialCount)
	}

	// Step once
	g = g.Step()

	// After 1 step, glider should still have 5 cells
	afterStep := g.CountLiving()
	if afterStep != 5 {
		t.Errorf("Expected 5 living cells after step, got %d", afterStep)
	}
}

func TestClassicGrid_Step_Block(t *testing.T) {
	g := NewClassicGrid(10, 10)

	// Set up a block (2x2 still life)
	g.Set(1, 1, 1)
	g.Set(2, 1, 1)
	g.Set(1, 2, 1)
	g.Set(2, 2, 1)

	// Step multiple times - block should remain stable
	for i := 0; i < 10; i++ {
		g = g.Step()
		if g.CountLiving() != 4 {
			t.Errorf("Generation %d: Expected 4 living cells, got %d", i+1, g.CountLiving())
		}
	}

	// Verify block cells are still alive
	if g.Get(1, 1) != 1 {
		t.Error("Block cell (1,1) should be alive")
	}
	if g.Get(2, 2) != 1 {
		t.Error("Block cell (2,2) should be alive")
	}
}

func TestClassicGrid_Step_Blinker(t *testing.T) {
	g := NewClassicGrid(10, 10)

	// Set up a blinker (period-2 oscillator)
	// Horizontal: ● ● ●
	g.Set(1, 1, 1)
	g.Set(2, 1, 1)
	g.Set(3, 1, 1)

	// Should have 3 cells
	if g.CountLiving() != 3 {
		t.Errorf("Expected 3 living cells, got %d", g.CountLiving())
	}

	// After 1 step, should rotate to vertical
	g = g.Step()
	if g.CountLiving() != 3 {
		t.Errorf("Expected 3 living cells after step 1, got %d", g.CountLiving())
	}

	// Vertical blinker should have cells at (2,0), (2,1), (2,2)
	if g.Get(2, 0) != 1 {
		t.Error("Blinker should be vertical at (2,0)")
	}
	if g.Get(2, 1) != 1 {
		t.Error("Blinker should be vertical at (2,1)")
	}
	if g.Get(2, 2) != 1 {
		t.Error("Blinker should be vertical at (2,2)")
	}

	// After 2 steps total, should return to horizontal
	g = g.Step()
	if g.Get(1, 1) != 1 {
		t.Error("Blinker should return to horizontal at (1,1)")
	}
	if g.Get(3, 1) != 1 {
		t.Error("Blinker should return to horizontal at (3,1)")
	}
}

func TestClassicGrid_Clone(t *testing.T) {
	g := NewClassicGrid(5, 5)
	g.Set(2, 2, 1)
	g.Set(3, 3, 1)

	// Clone
	clone := g.Clone()

	// Verify clone has same state
	if clone.Get(2, 2) != 1 {
		t.Error("Clone should have alive cell at (2,2)")
	}

	// Modify original
	g.Set(0, 0, 1)

	// Clone should not be affected
	if clone.Get(0, 0) != 0 {
		t.Error("Clone should not be affected by original modification")
	}
}

func TestClassicGrid_CountLiving(t *testing.T) {
	g := NewClassicGrid(10, 10)

	if g.CountLiving() != 0 {
		t.Errorf("Empty grid should have 0 living cells, got %d", g.CountLiving())
	}

	g.Set(0, 0, 1)
	g.Set(5, 5, 1)
	g.Set(9, 9, 1)

	if g.CountLiving() != 3 {
		t.Errorf("Expected 3 living cells, got %d", g.CountLiving())
	}
}

func TestClassicGrid_CountNeighbors(t *testing.T) {
	g := NewClassicGrid(5, 5)

	// Set up a 3x3 block
	for y := 1; y <= 3; y++ {
		for x := 1; x <= 3; x++ {
			g.Set(x, y, 1)
		}
	}

	// Center cell should have 8 neighbors
	if g.CountNeighbors(2, 2) != 8 {
		t.Errorf("Center cell should have 8 neighbors, got %d", g.CountNeighbors(2, 2))
	}

	// Corner cell should have 3 neighbors
	if g.CountNeighbors(1, 1) != 3 {
		t.Errorf("Corner cell should have 3 neighbors, got %d", g.CountNeighbors(1, 1))
	}

	// Edge cell should have 5 neighbors
	if g.CountNeighbors(2, 1) != 5 {
		t.Errorf("Edge cell should have 5 neighbors, got %d", g.CountNeighbors(2, 1))
	}
}
