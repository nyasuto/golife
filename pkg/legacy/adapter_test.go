package legacy

import (
	"golife/pkg/core"
	"golife/pkg/patterns"
	"golife/pkg/rules"
	"golife/pkg/universe"
	"testing"
)

func TestClassicGrid_ToUniverse(t *testing.T) {
	// Create a classic grid with a glider pattern
	g := NewClassicGrid(10, 10)
	g.Set(1, 0, 1)
	g.Set(2, 1, 1)
	g.Set(0, 2, 1)
	g.Set(1, 2, 1)
	g.Set(2, 2, 1)

	// Convert to Universe
	u := g.ToUniverse()

	// Verify dimension
	if u.Dimension() != core.Dim2D {
		t.Errorf("Expected Dim2D, got %d", u.Dimension())
	}

	// Verify size
	size := u.Size()
	if size.X != 10 || size.Y != 10 {
		t.Errorf("Expected size (10,10), got (%d,%d)", size.X, size.Y)
	}

	// Verify living cells
	if u.CountLiving() != 5 {
		t.Errorf("Expected 5 living cells, got %d", u.CountLiving())
	}

	// Verify specific cells
	if u.Get(core.NewCoord2D(1, 0)) != core.Alive {
		t.Error("Cell (1,0) should be alive")
	}
	if u.Get(core.NewCoord2D(0, 0)) != core.Dead {
		t.Error("Cell (0,0) should be dead")
	}
}

func TestNewFromUniverse(t *testing.T) {
	// Create a Universe2D with a block pattern
	u := universe.New2D(10, 10, rules.ConwayRule{})
	u.Set(core.NewCoord2D(1, 1), core.Alive)
	u.Set(core.NewCoord2D(2, 1), core.Alive)
	u.Set(core.NewCoord2D(1, 2), core.Alive)
	u.Set(core.NewCoord2D(2, 2), core.Alive)

	// Convert to ClassicGrid
	g := NewFromUniverse(u)

	// Verify dimensions
	if g.Width != 10 || g.Height != 10 {
		t.Errorf("Expected dimensions (10,10), got (%d,%d)", g.Width, g.Height)
	}

	// Verify living cells
	if g.CountLiving() != 4 {
		t.Errorf("Expected 4 living cells, got %d", g.CountLiving())
	}

	// Verify specific cells
	if g.Get(1, 1) != 1 {
		t.Error("Cell (1,1) should be alive")
	}
	if g.Get(0, 0) != 0 {
		t.Error("Cell (0,0) should be dead")
	}
}

func TestRoundTripConversion(t *testing.T) {
	// Create a classic grid
	original := NewClassicGrid(5, 5)
	original.Set(1, 1, 1)
	original.Set(2, 2, 1)
	original.Set(3, 3, 1)

	// Convert to Universe and back
	u := original.ToUniverse()
	converted := NewFromUniverse(u)

	// Verify they match
	if converted.Width != original.Width || converted.Height != original.Height {
		t.Error("Dimensions don't match after round-trip conversion")
	}

	if converted.CountLiving() != original.CountLiving() {
		t.Error("Living cell count doesn't match after round-trip conversion")
	}

	for y := 0; y < original.Height; y++ {
		for x := 0; x < original.Width; x++ {
			if converted.Get(x, y) != original.Get(x, y) {
				t.Errorf("Cell (%d,%d) mismatch after round-trip", x, y)
			}
		}
	}
}

func TestFromSlice(t *testing.T) {
	cells := [][]int{
		{0, 1, 0},
		{0, 0, 1},
		{1, 1, 1},
	}

	g := FromSlice(cells)

	if g.Width != 3 || g.Height != 3 {
		t.Errorf("Expected dimensions (3,3), got (%d,%d)", g.Width, g.Height)
	}

	if g.Get(1, 0) != 1 {
		t.Error("Cell (1,0) should be alive")
	}

	if g.CountLiving() != 5 {
		t.Errorf("Expected 5 living cells, got %d", g.CountLiving())
	}
}

func TestToSlice(t *testing.T) {
	g := NewClassicGrid(3, 3)
	g.Set(1, 0, 1)
	g.Set(2, 1, 1)
	g.Set(0, 2, 1)
	g.Set(1, 2, 1)
	g.Set(2, 2, 1)

	slice := g.ToSlice()

	if len(slice) != 3 || len(slice[0]) != 3 {
		t.Error("Slice dimensions incorrect")
	}

	if slice[0][1] != 1 {
		t.Error("Cell [0][1] should be 1")
	}

	// Modify slice should not affect original
	slice[0][0] = 1
	if g.Get(0, 0) != 0 {
		t.Error("Original grid should not be affected by slice modification")
	}
}

func TestCompatibilityWithPatterns(t *testing.T) {
	// Load a pattern using the new system
	u := universe.New2D(50, 50, rules.ConwayRule{})
	glider := patterns.Glider()
	glider.LoadIntoUniverse(u, 10, 10)

	// Convert to classic grid
	g := NewFromUniverse(u)

	// Verify it works the same
	if g.CountLiving() != 5 {
		t.Errorf("Glider should have 5 cells, got %d", g.CountLiving())
	}

	// Step using classic method
	g2 := g.Step()

	// Step using new method
	u.Step()

	// Convert back and compare
	g3 := NewFromUniverse(u)

	if g2.CountLiving() != g3.CountLiving() {
		t.Errorf("Living counts don't match: classic=%d, new=%d", g2.CountLiving(), g3.CountLiving())
	}
}

// BenchmarkClassicVsNew compares performance
func BenchmarkClassicGrid_Step(b *testing.B) {
	g := NewClassicGrid(100, 100)
	g.Randomize()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g = g.Step()
	}
}

func BenchmarkUniverse2D_Step(b *testing.B) {
	u := universe.New2D(100, 100, rules.ConwayRule{})
	u.Randomize()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u.Step()
	}
}
