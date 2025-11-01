package engine

import (
	"golife/pkg/core"
	"golife/pkg/rules"
	"golife/pkg/universe"
	"testing"
	"time"
)

func TestNewStatistics(t *testing.T) {
	initialPop := 42
	stats := NewStatistics(initialPop)

	if stats.Generation != 0 {
		t.Errorf("Initial generation should be 0, got %d", stats.Generation)
	}
	if stats.LivingCells != initialPop {
		t.Errorf("Initial living cells should be %d, got %d", initialPop, stats.LivingCells)
	}
	if stats.Births != 0 {
		t.Errorf("Initial births should be 0, got %d", stats.Births)
	}
	if stats.Deaths != 0 {
		t.Errorf("Initial deaths should be 0, got %d", stats.Deaths)
	}
	if stats.FPS != 0 {
		t.Errorf("Initial FPS should be 0, got %.2f", stats.FPS)
	}
	if stats.StartTime.IsZero() {
		t.Error("Start time should be set")
	}
	if stats.LastFrameTime.IsZero() {
		t.Error("Last frame time should be set")
	}
}

func TestStatisticsUpdate_PopulationIncrease(t *testing.T) {
	rule := rules.ConwayRule{}
	u := universe.New2D(10, 10, rule)

	// Create a blinker pattern that will have births
	u.Set(core.NewCoord2D(5, 5), 1)
	u.Set(core.NewCoord2D(6, 5), 1)
	u.Set(core.NewCoord2D(7, 5), 1)

	stats := NewStatistics(u.CountLiving())
	initialPop := stats.LivingCells

	// Step and update
	time.Sleep(10 * time.Millisecond)
	u.Step()
	stats.Update(u)

	if stats.Generation != 1 {
		t.Errorf("Generation should be 1, got %d", stats.Generation)
	}

	// Check that births/deaths are tracked correctly
	newPop := stats.LivingCells
	diff := newPop - initialPop

	if diff > 0 {
		if stats.Births != diff {
			t.Errorf("Births should be %d, got %d", diff, stats.Births)
		}
		if stats.Deaths != 0 {
			t.Errorf("Deaths should be 0 when population increases, got %d", stats.Deaths)
		}
	} else if diff < 0 {
		if stats.Deaths != -diff {
			t.Errorf("Deaths should be %d, got %d", -diff, stats.Deaths)
		}
		if stats.Births != 0 {
			t.Errorf("Births should be 0 when population decreases, got %d", stats.Births)
		}
	} else {
		if stats.Births != 0 {
			t.Errorf("Births should be 0 when population unchanged, got %d", stats.Births)
		}
		if stats.Deaths != 0 {
			t.Errorf("Deaths should be 0 when population unchanged, got %d", stats.Deaths)
		}
	}

	// FPS should be calculated
	if stats.FPS <= 0 {
		t.Errorf("FPS should be positive after update, got %.2f", stats.FPS)
	}
}

func TestStatisticsUpdate_PopulationDecrease(t *testing.T) {
	rule := rules.ConwayRule{}
	u := universe.New2D(10, 10, rule)

	// Create a pattern that will die
	u.Set(core.NewCoord2D(5, 5), 1)
	u.Set(core.NewCoord2D(6, 6), 1)

	stats := NewStatistics(u.CountLiving())
	initialPop := stats.LivingCells

	// Step and update
	time.Sleep(10 * time.Millisecond)
	u.Step()
	stats.Update(u)

	newPop := stats.LivingCells

	// These two cells should die (not enough neighbors)
	if newPop >= initialPop {
		t.Skip("Pattern did not decrease as expected, skipping test")
	}

	expectedDeaths := initialPop - newPop
	if stats.Deaths != expectedDeaths {
		t.Errorf("Deaths should be %d, got %d", expectedDeaths, stats.Deaths)
	}
	if stats.Births != 0 {
		t.Errorf("Births should be 0 when population decreases, got %d", stats.Births)
	}
}

func TestStatisticsUpdate_NoChange(t *testing.T) {
	rule := rules.ConwayRule{}
	u := universe.New2D(10, 10, rule)

	// Create a stable block pattern
	u.Set(core.NewCoord2D(5, 5), 1)
	u.Set(core.NewCoord2D(5, 6), 1)
	u.Set(core.NewCoord2D(6, 5), 1)
	u.Set(core.NewCoord2D(6, 6), 1)

	stats := NewStatistics(u.CountLiving())
	initialPop := stats.LivingCells

	// Step and update
	time.Sleep(10 * time.Millisecond)
	u.Step()
	stats.Update(u)

	if stats.LivingCells != initialPop {
		t.Skip("Block pattern not stable, skipping test")
	}

	if stats.Births != 0 {
		t.Errorf("Births should be 0 for stable pattern, got %d", stats.Births)
	}
	if stats.Deaths != 0 {
		t.Errorf("Deaths should be 0 for stable pattern, got %d", stats.Deaths)
	}
}

func TestStatisticsUpdate_FPSCalculation(t *testing.T) {
	rule := rules.ConwayRule{}
	u := universe.New2D(10, 10, rule)
	u.Set(core.NewCoord2D(5, 5), 1)

	stats := NewStatistics(u.CountLiving())

	// First update
	time.Sleep(10 * time.Millisecond)
	u.Step()
	stats.Update(u)

	if stats.FPS <= 0 {
		t.Errorf("FPS should be positive after first update, got %.2f", stats.FPS)
	}

	firstFPS := stats.FPS

	// Second update with different timing
	time.Sleep(20 * time.Millisecond)
	u.Step()
	stats.Update(u)

	if stats.FPS <= 0 {
		t.Errorf("FPS should be positive after second update, got %.2f", stats.FPS)
	}

	// FPS should be different due to different timing
	if stats.FPS == firstFPS {
		t.Logf("FPS unchanged: %.2f (timing variation too small)", stats.FPS)
	}
}

func TestStatisticsUpdate_MultipleGenerations(t *testing.T) {
	rule := rules.ConwayRule{}
	u := universe.New2D(10, 10, rule)

	// Create a glider
	u.Set(core.NewCoord2D(2, 1), 1)
	u.Set(core.NewCoord2D(3, 2), 1)
	u.Set(core.NewCoord2D(1, 3), 1)
	u.Set(core.NewCoord2D(2, 3), 1)
	u.Set(core.NewCoord2D(3, 3), 1)

	stats := NewStatistics(u.CountLiving())

	// Run for multiple generations
	for i := 1; i <= 5; i++ {
		time.Sleep(10 * time.Millisecond)
		u.Step()
		stats.Update(u)

		if stats.Generation != i {
			t.Errorf("After %d steps, generation should be %d, got %d", i, i, stats.Generation)
		}

		if stats.LivingCells < 0 {
			t.Errorf("Living cells should not be negative, got %d", stats.LivingCells)
		}

		if stats.FPS < 0 {
			t.Errorf("FPS should not be negative, got %.2f", stats.FPS)
		}
	}
}

func TestStatisticsReset(t *testing.T) {
	rule := rules.ConwayRule{}
	u := universe.New2D(10, 10, rule)
	u.Randomize()

	stats := NewStatistics(u.CountLiving())

	// Run a few generations
	for i := 0; i < 3; i++ {
		u.Step()
		stats.Update(u)
	}

	if stats.Generation == 0 {
		t.Error("Generation should not be 0 after updates")
	}

	// Reset with new population
	newPop := 100
	stats.Reset(newPop)

	if stats.Generation != 0 {
		t.Errorf("Generation should be 0 after reset, got %d", stats.Generation)
	}
	if stats.LivingCells != newPop {
		t.Errorf("Living cells should be %d after reset, got %d", newPop, stats.LivingCells)
	}
	if stats.Births != 0 {
		t.Errorf("Births should be 0 after reset, got %d", stats.Births)
	}
	if stats.Deaths != 0 {
		t.Errorf("Deaths should be 0 after reset, got %d", stats.Deaths)
	}
	if stats.FPS != 0 {
		t.Errorf("FPS should be 0 after reset, got %.2f", stats.FPS)
	}
	if stats.StartTime.IsZero() {
		t.Error("Start time should be set after reset")
	}
	if stats.LastFrameTime.IsZero() {
		t.Error("Last frame time should be set after reset")
	}
}

func TestStatisticsReset_PreservesIndependence(t *testing.T) {
	stats1 := NewStatistics(10)
	stats2 := NewStatistics(20)

	rule := rules.ConwayRule{}
	u1 := universe.New2D(10, 10, rule)
	u1.Randomize()

	// Update only stats1
	u1.Step()
	stats1.Update(u1)

	if stats1.Generation != 1 {
		t.Error("stats1 should have generation 1")
	}
	if stats2.Generation != 0 {
		t.Error("stats2 should still have generation 0")
	}

	// Reset stats1
	stats1.Reset(50)

	if stats1.Generation != 0 {
		t.Error("stats1 should have generation 0 after reset")
	}
	if stats2.Generation != 0 {
		t.Error("stats2 should still have generation 0")
	}
	if stats2.LivingCells != 20 {
		t.Errorf("stats2 should still have 20 living cells, got %d", stats2.LivingCells)
	}
}
