package universe

import (
	"golife/pkg/core"
	"golife/pkg/rules"
	"testing"
)

func TestNew2D(t *testing.T) {
	rule := rules.ConwayRule{}
	u := New2D(10, 10, rule)

	if u.width != 10 || u.height != 10 {
		t.Errorf("Expected dimensions (10,10), got (%d,%d)", u.width, u.height)
	}

	if u.Dimension() != core.Dim2D {
		t.Errorf("Expected Dim2D, got %d", u.Dimension())
	}

	size := u.Size()
	if size.X != 10 || size.Y != 10 {
		t.Errorf("Expected size (10,10), got (%d,%d)", size.X, size.Y)
	}
}

func TestUniverse2D_GetSet(t *testing.T) {
	rule := rules.ConwayRule{}
	u := New2D(5, 5, rule)

	// Test Set and Get
	coord := core.NewCoord2D(2, 3)
	u.Set(coord, core.Alive)

	if u.Get(coord) != core.Alive {
		t.Errorf("Expected Alive, got %d", u.Get(coord))
	}

	// Test boundary
	outOfBounds := core.NewCoord2D(10, 10)
	if u.Get(outOfBounds) != core.Dead {
		t.Error("Out of bounds should return Dead")
	}
}

func TestUniverse2D_Step_Glider(t *testing.T) {
	rule := rules.ConwayRule{}
	u := New2D(10, 10, rule)

	// Set up a glider pattern
	// Gen 0:
	//  . ● .
	//  . . ●
	//  ● ● ●
	u.Set(core.NewCoord2D(1, 0), core.Alive)
	u.Set(core.NewCoord2D(2, 1), core.Alive)
	u.Set(core.NewCoord2D(0, 2), core.Alive)
	u.Set(core.NewCoord2D(1, 2), core.Alive)
	u.Set(core.NewCoord2D(2, 2), core.Alive)

	initialCount := u.CountLiving()
	if initialCount != 5 {
		t.Errorf("Expected 5 living cells, got %d", initialCount)
	}

	// Step once
	u.Step()

	// After 1 step, glider should still have 5 cells
	afterStep := u.CountLiving()
	if afterStep != 5 {
		t.Errorf("Expected 5 living cells after step, got %d", afterStep)
	}
}

func TestUniverse2D_Step_Block(t *testing.T) {
	rule := rules.ConwayRule{}
	u := New2D(10, 10, rule)

	// Set up a block (2x2 still life)
	// ● ●
	// ● ●
	u.Set(core.NewCoord2D(1, 1), core.Alive)
	u.Set(core.NewCoord2D(2, 1), core.Alive)
	u.Set(core.NewCoord2D(1, 2), core.Alive)
	u.Set(core.NewCoord2D(2, 2), core.Alive)

	// Step multiple times - block should remain stable
	for i := 0; i < 10; i++ {
		u.Step()
		if u.CountLiving() != 4 {
			t.Errorf("Generation %d: Expected 4 living cells, got %d", i+1, u.CountLiving())
		}
	}

	// Verify block cells are still alive
	if u.Get(core.NewCoord2D(1, 1)) != core.Alive {
		t.Error("Block cell (1,1) should be alive")
	}
	if u.Get(core.NewCoord2D(2, 2)) != core.Alive {
		t.Error("Block cell (2,2) should be alive")
	}
}

func TestUniverse2D_Step_Blinker(t *testing.T) {
	rule := rules.ConwayRule{}
	u := New2D(10, 10, rule)

	// Set up a blinker (period-2 oscillator)
	// Horizontal: ● ● ●
	u.Set(core.NewCoord2D(1, 1), core.Alive)
	u.Set(core.NewCoord2D(2, 1), core.Alive)
	u.Set(core.NewCoord2D(3, 1), core.Alive)

	// Should have 3 cells
	if u.CountLiving() != 3 {
		t.Errorf("Expected 3 living cells, got %d", u.CountLiving())
	}

	// After 1 step, should rotate to vertical
	u.Step()
	if u.CountLiving() != 3 {
		t.Errorf("Expected 3 living cells after step 1, got %d", u.CountLiving())
	}

	// Vertical blinker should have cells at (2,0), (2,1), (2,2)
	if u.Get(core.NewCoord2D(2, 0)) != core.Alive {
		t.Error("Blinker should be vertical at (2,0)")
	}
	if u.Get(core.NewCoord2D(2, 1)) != core.Alive {
		t.Error("Blinker should be vertical at (2,1)")
	}
	if u.Get(core.NewCoord2D(2, 2)) != core.Alive {
		t.Error("Blinker should be vertical at (2,2)")
	}

	// After 2 steps total, should return to horizontal
	u.Step()
	if u.Get(core.NewCoord2D(1, 1)) != core.Alive {
		t.Error("Blinker should return to horizontal at (1,1)")
	}
	if u.Get(core.NewCoord2D(3, 1)) != core.Alive {
		t.Error("Blinker should return to horizontal at (3,1)")
	}
}

func TestUniverse2D_Clear(t *testing.T) {
	rule := rules.ConwayRule{}
	u := New2D(5, 5, rule)

	// Set some cells
	u.Set(core.NewCoord2D(0, 0), core.Alive)
	u.Set(core.NewCoord2D(1, 1), core.Alive)

	if u.CountLiving() != 2 {
		t.Errorf("Expected 2 living cells, got %d", u.CountLiving())
	}

	// Clear
	u.Clear()

	if u.CountLiving() != 0 {
		t.Errorf("Expected 0 living cells after clear, got %d", u.CountLiving())
	}
}

func TestUniverse2D_Clone(t *testing.T) {
	rule := rules.ConwayRule{}
	u := New2D(5, 5, rule)

	// Set some cells
	u.Set(core.NewCoord2D(2, 2), core.Alive)
	u.Set(core.NewCoord2D(3, 3), core.Alive)

	// Clone
	clone := u.Clone()

	// Verify clone has same state
	if clone.Get(core.NewCoord2D(2, 2)) != core.Alive {
		t.Error("Clone should have alive cell at (2,2)")
	}

	// Modify original
	u.Set(core.NewCoord2D(0, 0), core.Alive)

	// Clone should not be affected
	if clone.Get(core.NewCoord2D(0, 0)) != core.Dead {
		t.Error("Clone should not be affected by original modification")
	}
}

func TestUniverse2D_CountLiving(t *testing.T) {
	rule := rules.ConwayRule{}
	u := New2D(10, 10, rule)

	if u.CountLiving() != 0 {
		t.Errorf("Empty universe should have 0 living cells, got %d", u.CountLiving())
	}

	u.Set(core.NewCoord2D(0, 0), core.Alive)
	u.Set(core.NewCoord2D(5, 5), core.Alive)
	u.Set(core.NewCoord2D(9, 9), core.Alive)

	if u.CountLiving() != 3 {
		t.Errorf("Expected 3 living cells, got %d", u.CountLiving())
	}
}
