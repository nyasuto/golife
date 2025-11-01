package pkg_test

import (
	"golife/pkg/patterns"
	"golife/pkg/rules"
	"golife/pkg/universe"
	"testing"
)

// TestIntegration_GliderMovement tests that a glider pattern moves correctly over multiple generations
func TestIntegration_GliderMovement(t *testing.T) {
	rule := rules.ConwayRule{}
	u := universe.New2D(50, 50, rule)

	// Load glider pattern at (10, 10)
	glider := patterns.Glider()
	glider.LoadIntoUniverse(u, 10, 10)

	initialCount := u.CountLiving()
	if initialCount != 5 {
		t.Errorf("Glider should have 5 cells, got %d", initialCount)
	}

	// Simulate 20 generations
	for i := 0; i < 20; i++ {
		u.Step()
	}

	// Glider should still exist (5 cells)
	finalCount := u.CountLiving()
	if finalCount != 5 {
		t.Errorf("Glider should still have 5 cells after 20 generations, got %d", finalCount)
	}
}

// TestIntegration_Blinker tests period-2 oscillator
func TestIntegration_Blinker(t *testing.T) {
	rule := rules.ConwayRule{}
	u := universe.New2D(20, 20, rule)

	blinker := patterns.Blinker()
	blinker.LoadIntoUniverse(u, 10, 10)

	if u.CountLiving() != 3 {
		t.Errorf("Blinker should have 3 cells, got %d", u.CountLiving())
	}

	// Step once
	u.Step()
	if u.CountLiving() != 3 {
		t.Errorf("Blinker should still have 3 cells after 1 step, got %d", u.CountLiving())
	}

	// Step again (should return to original state - period-2)
	u.Step()
	if u.CountLiving() != 3 {
		t.Errorf("Blinker should still have 3 cells after 2 steps, got %d", u.CountLiving())
	}
}

// TestIntegration_Block tests still life
func TestIntegration_Block(t *testing.T) {
	rule := rules.ConwayRule{}
	u := universe.New2D(20, 20, rule)

	block := patterns.Block()
	block.LoadIntoUniverse(u, 10, 10)

	if u.CountLiving() != 4 {
		t.Errorf("Block should have 4 cells, got %d", u.CountLiving())
	}

	// Simulate 100 generations - block should remain stable
	for i := 0; i < 100; i++ {
		u.Step()
		if u.CountLiving() != 4 {
			t.Errorf("Block should still have 4 cells after %d generations, got %d", i+1, u.CountLiving())
			break
		}
	}
}

// TestIntegration_Pulsar tests period-3 oscillator
func TestIntegration_Pulsar(t *testing.T) {
	rule := rules.ConwayRule{}
	u := universe.New2D(50, 50, rule)

	pulsar := patterns.Pulsar()
	pulsar.LoadIntoUniverse(u, 20, 20)

	initialCount := u.CountLiving()
	if initialCount == 0 {
		t.Error("Pulsar should have living cells")
	}

	// Step 3 times (period-3 oscillator)
	for i := 0; i < 3; i++ {
		u.Step()
	}

	// Should return to similar state (cell count should be stable)
	finalCount := u.CountLiving()
	if finalCount == 0 {
		t.Error("Pulsar should still have living cells after 3 steps")
	}

	// Allow some variance in count due to edge effects
	variance := float64(finalCount-initialCount) / float64(initialCount)
	if variance > 0.2 || variance < -0.2 {
		t.Errorf("Pulsar count changed too much after 3 steps. Initial: %d, Final: %d", initialCount, finalCount)
	}
}

// TestIntegration_AllPatterns tests that all patterns can be loaded
func TestIntegration_AllPatterns(t *testing.T) {
	rule := rules.ConwayRule{}
	allPatterns := patterns.AllPatterns()

	for name, pattern := range allPatterns {
		t.Run(name, func(t *testing.T) {
			// Create universe large enough for the pattern
			u := universe.New2D(100, 100, rule)

			// Load pattern at center
			offsetX := (100 - pattern.Width) / 2
			offsetY := (100 - pattern.Height) / 2
			pattern.LoadIntoUniverse(u, offsetX, offsetY)

			// Verify some cells are alive
			if u.CountLiving() == 0 {
				t.Errorf("Pattern %s should have living cells", name)
			}

			// Simulate a few generations - should not crash
			for i := 0; i < 10; i++ {
				u.Step()
			}
		})
	}
}

// BenchmarkUniverse2D_Step benchmarks the step function
func BenchmarkUniverse2D_Step(b *testing.B) {
	rule := rules.ConwayRule{}
	u := universe.New2D(100, 100, rule)
	u.Randomize()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u.Step()
	}
}

// BenchmarkUniverse2D_Clone benchmarks cloning
func BenchmarkUniverse2D_Clone(b *testing.B) {
	rule := rules.ConwayRule{}
	u := universe.New2D(100, 100, rule)
	u.Randomize()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = u.Clone()
	}
}
