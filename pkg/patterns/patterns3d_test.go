package patterns

import (
	"golife/pkg/core"
	"golife/pkg/rules"
	"golife/pkg/universe"
	"testing"
)

func TestGlider3D(t *testing.T) {
	glider := Glider3D()

	if glider == nil {
		t.Fatal("Glider3D should not return nil")
	}

	if glider.Name != "3D Glider" {
		t.Errorf("Expected name '3D Glider', got '%s'", glider.Name)
	}

	if len(glider.Cells) != 4 {
		t.Errorf("Expected 4 cells, got %d", len(glider.Cells))
	}

	// Check that all cells are alive
	for _, state := range glider.Cells {
		if state != core.Alive {
			t.Error("All glider cells should be alive")
		}
	}
}

func TestBlock3D(t *testing.T) {
	block := Block3D()

	if block == nil {
		t.Fatal("Block3D should not return nil")
	}

	if block.Name != "3D Block" {
		t.Errorf("Expected name '3D Block', got '%s'", block.Name)
	}

	// 2x2x2 cube = 8 cells
	if len(block.Cells) != 8 {
		t.Errorf("Expected 8 cells, got %d", len(block.Cells))
	}
}

func TestOscillator3D_Period2(t *testing.T) {
	osc := Oscillator3D_Period2()

	if osc == nil {
		t.Fatal("Oscillator3D_Period2 should not return nil")
	}

	if len(osc.Cells) != 3 {
		t.Errorf("Expected 3 cells, got %d", len(osc.Cells))
	}
}

func TestPattern3D_LoadIntoUniverse3D(t *testing.T) {
	rule := rules.Life3D_B6S567{}
	u := universe.New3D(10, 10, 10, rule)

	glider := Glider3D()
	glider.LoadIntoUniverse3D(u, 3, 3, 3)

	// Check that cells were loaded
	count := u.CountLiving()
	if count != 4 {
		t.Errorf("Expected 4 living cells after loading glider, got %d", count)
	}

	// Check specific cell positions
	if u.Get(core.NewCoord3D(3, 3, 3)) != core.Alive {
		t.Error("Cell at (3,3,3) should be alive")
	}
	if u.Get(core.NewCoord3D(4, 3, 3)) != core.Alive {
		t.Error("Cell at (4,3,3) should be alive")
	}
}

func TestPattern3D_CreateUniverse(t *testing.T) {
	rule := rules.Life3D_B6S567{}
	glider := Glider3D()

	u := glider.CreateUniverse(rule)

	if u == nil {
		t.Fatal("CreateUniverse should not return nil")
	}

	size := u.Size()
	if size.X != glider.Width || size.Y != glider.Height || size.Z != glider.Depth {
		t.Errorf("Universe size mismatch: got (%d,%d,%d), want (%d,%d,%d)",
			size.X, size.Y, size.Z, glider.Width, glider.Height, glider.Depth)
	}

	count := u.CountLiving()
	if count != 4 {
		t.Errorf("Expected 4 living cells, got %d", count)
	}
}

func TestGetPatterns3D(t *testing.T) {
	patterns := GetPatterns3D()

	if len(patterns) == 0 {
		t.Fatal("GetPatterns3D should return at least one pattern")
	}

	expectedPatterns := []string{"glider", "block", "oscillator"}
	for _, name := range expectedPatterns {
		if _, ok := patterns[name]; !ok {
			t.Errorf("Pattern '%s' not found in GetPatterns3D", name)
		}
	}
}

func TestListPatterns3D(t *testing.T) {
	patterns := ListPatterns3D()

	if len(patterns) == 0 {
		t.Fatal("ListPatterns3D should return at least one pattern name")
	}

	expectedCount := 3
	if len(patterns) != expectedCount {
		t.Errorf("Expected %d patterns, got %d", expectedCount, len(patterns))
	}
}

func TestLoadPattern3D(t *testing.T) {
	t.Run("Load existing pattern", func(t *testing.T) {
		glider := LoadPattern3D("glider")
		if glider == nil {
			t.Fatal("LoadPattern3D('glider') should not return nil")
		}
		if glider.Name != "3D Glider" {
			t.Errorf("Expected '3D Glider', got '%s'", glider.Name)
		}
	})

	t.Run("Load non-existent pattern", func(t *testing.T) {
		pattern := LoadPattern3D("nonexistent")
		if pattern != nil {
			t.Error("LoadPattern3D with invalid name should return nil")
		}
	})
}

func TestDemoUniverse3D(t *testing.T) {
	u := DemoUniverse3D()

	if u == nil {
		t.Fatal("DemoUniverse3D should not return nil")
	}

	size := u.Size()
	if size.X != 20 || size.Y != 20 || size.Z != 20 {
		t.Errorf("Expected 20x20x20 universe, got %dx%dx%d", size.X, size.Y, size.Z)
	}

	count := u.CountLiving()
	if count != 4 {
		t.Errorf("Expected 4 living cells (glider), got %d", count)
	}
}

func TestPattern3D_Evolution(t *testing.T) {
	rule := rules.Life3D_B6S567{}

	t.Run("Block stability", func(t *testing.T) {
		block := Block3D()
		u := block.CreateUniverse(rule)

		initialCount := u.CountLiving()

		// Run for 10 generations
		for i := 0; i < 10; i++ {
			u.Step()
		}

		finalCount := u.CountLiving()

		// Block should remain stable (though it might not in B6/S567)
		t.Logf("Block: initial=%d, after 10 steps=%d", initialCount, finalCount)
	})

	t.Run("Glider evolution", func(t *testing.T) {
		glider := Glider3D()
		u := glider.CreateUniverse(rule)

		initialCount := u.CountLiving()

		// Run for 20 generations
		for i := 0; i < 20; i++ {
			u.Step()
		}

		finalCount := u.CountLiving()

		t.Logf("Glider: initial=%d, after 20 steps=%d", initialCount, finalCount)

		// The simple 4-cell pattern might not be a true glider in B6/S567
		// This is exploratory - we're documenting the behavior
		// B6/S567 requires specific patterns that are different from 2D gliders
	})
}

func TestPattern3D_WithLargeUniverse(t *testing.T) {
	rule := rules.Life3D_B6S567{}
	u := universe.New3D(50, 50, 50, rule)

	glider := Glider3D()
	// Place glider away from boundaries
	glider.LoadIntoUniverse3D(u, 20, 20, 20)

	initialCount := u.CountLiving()
	if initialCount != 4 {
		t.Errorf("Expected 4 initial cells, got %d", initialCount)
	}

	// Run simulation
	for i := 0; i < 50; i++ {
		u.Step()
	}

	finalCount := u.CountLiving()
	t.Logf("Glider in large universe: initial=%d, after 50 steps=%d", initialCount, finalCount)
}

func TestTestPattern3D(t *testing.T) {
	pattern := TestPattern3D()

	if pattern == nil {
		t.Fatal("TestPattern3D should not return nil")
	}

	if pattern.Name != "Test Pattern" {
		t.Errorf("Expected name 'Test Pattern', got '%s'", pattern.Name)
	}

	if len(pattern.Cells) != 3 {
		t.Errorf("Expected 3 cells, got %d", len(pattern.Cells))
	}
}

func BenchmarkPattern3D_LoadIntoUniverse(b *testing.B) {
	rule := rules.Life3D_B6S567{}
	glider := Glider3D()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u := universe.New3D(20, 20, 20, rule)
		glider.LoadIntoUniverse3D(u, 8, 8, 8)
	}
}

func BenchmarkPattern3D_CreateUniverse(b *testing.B) {
	rule := rules.Life3D_B6S567{}
	glider := Glider3D()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		glider.CreateUniverse(rule)
	}
}
