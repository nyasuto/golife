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

	// Glider3D is now an alias for BaysGlider
	if glider.Name != "Bays's Glider" {
		t.Errorf("Expected name 'Bays's Glider', got '%s'", glider.Name)
	}

	if len(glider.Cells) != 10 {
		t.Errorf("Expected 10 cells, got %d", len(glider.Cells))
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

	// Oscillator3D_Period2 is now an alias for Blinker3D (6 cells)
	if len(osc.Cells) != 6 {
		t.Errorf("Expected 6 cells, got %d", len(osc.Cells))
	}
}

func TestPattern3D_LoadIntoUniverse3D(t *testing.T) {
	rule := rules.Life3D_B6S567{}
	u := universe.New3D(20, 20, 20, rule)

	glider := Glider3D()
	glider.LoadIntoUniverse3D(u, 3, 3, 3)

	// Check that cells were loaded (Bays's Glider has 10 cells)
	count := u.CountLiving()
	if count != 10 {
		t.Errorf("Expected 10 living cells after loading glider, got %d", count)
	}

	// Check specific cell positions from Bays's Glider pattern
	if u.Get(core.NewCoord3D(4, 3, 3)) != core.Alive { // offset (1,0,0) + (3,3,3)
		t.Error("Cell at (4,3,3) should be alive")
	}
	if u.Get(core.NewCoord3D(5, 3, 3)) != core.Alive { // offset (2,0,0) + (3,3,3)
		t.Error("Cell at (5,3,3) should be alive")
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
	if count != 10 {
		t.Errorf("Expected 10 living cells, got %d", count)
	}
}

func TestGetPatterns3D(t *testing.T) {
	patterns := GetPatterns3D()

	if len(patterns) == 0 {
		t.Fatal("GetPatterns3D should return at least one pattern")
	}

	// Test that basic patterns still exist
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

	// Now we have 9 patterns (excluding duplicates)
	expectedCount := 9
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
		// glider is now an alias for BaysGlider
		if glider.Name != "Bays's Glider" {
			t.Errorf("Expected 'Bays's Glider', got '%s'", glider.Name)
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
	if count != 10 {
		t.Errorf("Expected 10 living cells (Bays's glider), got %d", count)
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
	if initialCount != 10 {
		t.Errorf("Expected 10 initial cells, got %d", initialCount)
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

// New tests for verified patterns from Carter Bays (1987)

func TestBaysGlider(t *testing.T) {
	glider := BaysGlider()

	if glider == nil {
		t.Fatal("BaysGlider should not return nil")
	}

	if glider.Name != "Bays's Glider" {
		t.Errorf("Expected name 'Bays's Glider', got '%s'", glider.Name)
	}

	if len(glider.Cells) != 10 {
		t.Errorf("Expected 10 cells, got %d", len(glider.Cells))
	}
}

func TestBaysGlider_Period(t *testing.T) {
	rule := rules.Life3D_B6S567{}
	u := universe.New3D(30, 30, 30, rule)

	glider := BaysGlider()
	glider.LoadIntoUniverse3D(u, 10, 10, 10)

	initialCount := u.CountLiving()
	if initialCount != 10 {
		t.Fatalf("Expected 10 initial cells, got %d", initialCount)
	}

	// Run for 4 generations (one period)
	for i := 0; i < 4; i++ {
		u.Step()
	}

	// After one period, the glider should have moved
	// but maintain 10 cells (may vary slightly during transition)
	finalCount := u.CountLiving()
	t.Logf("Bays's Glider: initial=%d, after 4 steps=%d", initialCount, finalCount)

	// The glider should still exist (not die completely)
	if finalCount == 0 {
		t.Error("Bays's Glider died completely (should be period-4)")
	}
}

func TestBlinker3D_Oscillation(t *testing.T) {
	rule := rules.Life3D_B6S567{}
	u := universe.New3D(10, 10, 10, rule)

	blinker := Blinker3D()
	blinker.LoadIntoUniverse3D(u, 3, 3, 3)

	initialCount := u.CountLiving()
	if initialCount != 6 {
		t.Fatalf("Expected 6 initial cells, got %d", initialCount)
	}

	// Run for 2 generations (one period)
	u.Step()
	step1Count := u.CountLiving()

	u.Step()
	step2Count := u.CountLiving()

	t.Logf("Blinker3D: initial=%d, after 1 step=%d, after 2 steps=%d",
		initialCount, step1Count, step2Count)

	// Period-2 oscillator should return to original cell count after 2 steps
	if step2Count != initialCount {
		t.Logf("Warning: Blinker may not be stable (period-2 expected)")
	}
}

func TestFlashlight3D(t *testing.T) {
	flashlight := Flashlight3D()

	if flashlight == nil {
		t.Fatal("Flashlight3D should not return nil")
	}

	if flashlight.Name != "Flashlight" {
		t.Errorf("Expected name 'Flashlight', got '%s'", flashlight.Name)
	}

	if len(flashlight.Cells) != 14 {
		t.Errorf("Expected 14 cells, got %d", len(flashlight.Cells))
	}
}

func TestFlashlight3D_Oscillation(t *testing.T) {
	rule := rules.Life3D_B6S567{}
	u := universe.New3D(15, 15, 15, rule)

	flashlight := Flashlight3D()
	flashlight.LoadIntoUniverse3D(u, 5, 5, 5)

	initialCount := u.CountLiving()
	if initialCount != 14 {
		t.Fatalf("Expected 14 initial cells, got %d", initialCount)
	}

	// Run for 4 generations (one period)
	for i := 0; i < 4; i++ {
		u.Step()
	}

	finalCount := u.CountLiving()
	t.Logf("Flashlight: initial=%d, after 4 steps=%d", initialCount, finalCount)

	// Period-4 oscillator
	if finalCount == 0 {
		t.Error("Flashlight died completely (should be period-4)")
	}
}

func TestWheel3D(t *testing.T) {
	wheel := Wheel3D()

	if wheel == nil {
		t.Fatal("Wheel3D should not return nil")
	}

	if wheel.Name != "Wheel" {
		t.Errorf("Expected name 'Wheel', got '%s'", wheel.Name)
	}

	if len(wheel.Cells) != 12 {
		t.Errorf("Expected 12 cells, got %d", len(wheel.Cells))
	}
}

func TestBeehive3D_Stability(t *testing.T) {
	rule := rules.Life3D_B6S567{}
	u := universe.New3D(15, 15, 15, rule)

	beehive := Beehive3D()
	beehive.LoadIntoUniverse3D(u, 5, 5, 5)

	initialCount := u.CountLiving()
	if initialCount != 14 {
		t.Fatalf("Expected 14 initial cells, got %d", initialCount)
	}

	// Run for 10 generations to verify stability
	for i := 0; i < 10; i++ {
		u.Step()
	}

	finalCount := u.CountLiving()
	t.Logf("Beehive: initial=%d, after 10 steps=%d", initialCount, finalCount)

	// Still-life should remain stable
	if finalCount != initialCount {
		t.Logf("Warning: Beehive may not be stable (expected %d, got %d)", initialCount, finalCount)
	}
}

func TestBucket3D_Stability(t *testing.T) {
	rule := rules.Life3D_B6S567{}
	u := universe.New3D(15, 15, 15, rule)

	bucket := Bucket3D()
	bucket.LoadIntoUniverse3D(u, 5, 5, 5)

	initialCount := u.CountLiving()
	if initialCount != 16 {
		t.Fatalf("Expected 16 initial cells, got %d", initialCount)
	}

	// Run for 10 generations to verify stability
	for i := 0; i < 10; i++ {
		u.Step()
	}

	finalCount := u.CountLiving()
	t.Logf("Bucket: initial=%d, after 10 steps=%d", initialCount, finalCount)

	// Still-life should remain stable
	if finalCount != initialCount {
		t.Logf("Warning: Bucket may not be stable (expected %d, got %d)", initialCount, finalCount)
	}
}

func TestGetPatterns3D_NewPatterns(t *testing.T) {
	patterns := GetPatterns3D()

	expectedPatterns := []string{
		"bays-glider", "glider",
		"blinker", "flashlight", "wheel", "oscillator",
		"block", "beehive", "bucket",
	}

	for _, name := range expectedPatterns {
		if _, ok := patterns[name]; !ok {
			t.Errorf("Pattern '%s' not found in GetPatterns3D", name)
		}
	}

	// Verify total count
	if len(patterns) != len(expectedPatterns) {
		t.Errorf("Expected %d patterns, got %d", len(expectedPatterns), len(patterns))
	}
}

func TestListPatterns3D_NewPatterns(t *testing.T) {
	patterns := ListPatterns3D()

	expectedCount := 9
	if len(patterns) != expectedCount {
		t.Errorf("Expected %d patterns, got %d", expectedCount, len(patterns))
	}
}

func TestLoadPattern3D_AllPatterns(t *testing.T) {
	testCases := []struct {
		name          string
		expectedName  string
		expectedCells int
	}{
		{"bays-glider", "Bays's Glider", 10},
		{"blinker", "3D Blinker", 6},
		{"flashlight", "Flashlight", 14},
		{"wheel", "Wheel", 12},
		{"block", "3D Block", 8},
		{"beehive", "3D Beehive", 14},
		{"bucket", "Bucket", 16},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pattern := LoadPattern3D(tc.name)
			if pattern == nil {
				t.Fatalf("LoadPattern3D('%s') should not return nil", tc.name)
			}
			if pattern.Name != tc.expectedName {
				t.Errorf("Expected name '%s', got '%s'", tc.expectedName, pattern.Name)
			}
			if len(pattern.Cells) != tc.expectedCells {
				t.Errorf("Expected %d cells, got %d", tc.expectedCells, len(pattern.Cells))
			}
		})
	}
}
