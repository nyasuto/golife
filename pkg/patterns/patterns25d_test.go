package patterns

import (
	"golife/pkg/rules"
	"testing"
)

func TestVerticalGlider(t *testing.T) {
	p := VerticalGlider()

	if p.Name != "Vertical Glider" {
		t.Error("Pattern name should be 'Vertical Glider'")
	}

	if p.Width != 12 || p.Height != 12 || p.Depth != 5 {
		t.Errorf("Expected dimensions (12,12,5), got (%d,%d,%d)", p.Width, p.Height, p.Depth)
	}

	if len(p.Cells) == 0 {
		t.Error("Pattern should have cells")
	}

	// Verify cells exist in multiple layers
	layersWithCells := make(map[int]bool)
	for coord := range p.Cells {
		layersWithCells[coord.Z] = true
	}

	if len(layersWithCells) < 3 {
		t.Errorf("Vertical glider should span at least 3 layers, got %d", len(layersWithCells))
	}
}

func TestLayerOscillator(t *testing.T) {
	p := LayerOscillator()

	if p.Name != "Layer Oscillator" {
		t.Error("Pattern name should be 'Layer Oscillator'")
	}

	if p.Depth != 3 {
		t.Errorf("Expected depth 3, got %d", p.Depth)
	}

	if len(p.Cells) == 0 {
		t.Error("Pattern should have cells")
	}

	// Verify cells in different layers
	layerCounts := make(map[int]int)
	for coord := range p.Cells {
		layerCounts[coord.Z]++
	}

	if len(layerCounts) < 2 {
		t.Error("Oscillator should have cells in at least 2 layers")
	}
}

func TestLayerStack(t *testing.T) {
	p := LayerStack()

	if p.Name != "Layer Stack" {
		t.Error("Pattern name should be 'Layer Stack'")
	}

	if p.Depth != 3 {
		t.Errorf("Expected depth 3, got %d", p.Depth)
	}

	// Each layer should have 4 cells (2x2 blocks)
	layerCounts := make(map[int]int)
	for coord := range p.Cells {
		layerCounts[coord.Z]++
	}

	for z := 0; z < 3; z++ {
		if layerCounts[z] != 4 {
			t.Errorf("Layer %d should have 4 cells, got %d", z, layerCounts[z])
		}
	}
}

func TestVerticalBlinker(t *testing.T) {
	p := VerticalBlinker()

	if p.Name != "Vertical Blinker" {
		t.Error("Pattern name should be 'Vertical Blinker'")
	}

	// Should have cells in all 3 layers
	layerCounts := make(map[int]int)
	for coord := range p.Cells {
		layerCounts[coord.Z]++
	}

	if len(layerCounts) != 3 {
		t.Errorf("Vertical blinker should span 3 layers, got %d", len(layerCounts))
	}

	// Each layer should have 3 cells
	for z := 0; z < 3; z++ {
		if layerCounts[z] != 3 {
			t.Errorf("Layer %d should have 3 cells, got %d", z, layerCounts[z])
		}
	}
}

func TestLayerSandwich(t *testing.T) {
	p := LayerSandwich()

	if p.Name != "Layer Sandwich" {
		t.Error("Pattern name should be 'Layer Sandwich'")
	}

	// Should have cells in layer 0 and 2, but possibly not layer 1
	layerCounts := make(map[int]int)
	for coord := range p.Cells {
		layerCounts[coord.Z]++
	}

	if layerCounts[0] == 0 {
		t.Error("Layer 0 should have cells")
	}
	if layerCounts[2] == 0 {
		t.Error("Layer 2 should have cells")
	}
}

func TestEnergyWave(t *testing.T) {
	p := EnergyWave()

	if p.Name != "Energy Wave" {
		t.Error("Pattern name should be 'Energy Wave'")
	}

	// Should have different energy levels
	hasHighEnergy := false
	hasMediumEnergy := false
	hasLowEnergy := false

	for _, state := range p.Cells {
		if state >= 150 {
			hasHighEnergy = true
		} else if state >= 80 {
			hasMediumEnergy = true
		} else if state > 0 {
			hasLowEnergy = true
		}
	}

	if !hasHighEnergy {
		t.Error("Energy wave should have high energy cells")
	}
	if !hasMediumEnergy {
		t.Error("Energy wave should have medium energy cells")
	}
	if !hasLowEnergy {
		t.Error("Energy wave should have low energy cells")
	}
}

func TestLoadIntoUniverse25D(t *testing.T) {
	p := VerticalGlider()
	u := p.CreateUniverse(rules.ConwayRule{})

	if u == nil {
		t.Fatal("CreateUniverse should not return nil")
	}

	size := u.Size()
	if size.X != p.Width || size.Y != p.Height || size.Z != p.Depth {
		t.Errorf("Universe size doesn't match pattern dimensions")
	}

	// Verify cells were loaded
	livingCount := u.CountLiving()
	if livingCount != len(p.Cells) {
		t.Errorf("Expected %d living cells, got %d", len(p.Cells), livingCount)
	}
}

func TestLoadIntoUniverse25D_WithOffset(t *testing.T) {
	p := LayerStack()
	u := p.CreateUniverse(rules.ConwayRule{})

	initialCount := u.CountLiving()

	// Clear and reload with offset
	u.Clear()
	p.LoadIntoUniverse25D(u, 2, 2, 0)

	offsetCount := u.CountLiving()
	if offsetCount != initialCount {
		t.Errorf("Cell count should be the same after offset load, got %d vs %d", offsetCount, initialCount)
	}
}

func TestGetPatterns25D(t *testing.T) {
	patterns := GetPatterns25D()

	expectedPatterns := []string{
		"vertical-glider",
		"layer-oscillator",
		"layer-stack",
		"vertical-blinker",
		"layer-sandwich",
		"energy-wave",
	}

	for _, name := range expectedPatterns {
		if _, exists := patterns[name]; !exists {
			t.Errorf("Pattern '%s' should be in GetPatterns25D()", name)
		}
	}

	if len(patterns) != len(expectedPatterns) {
		t.Errorf("Expected %d patterns, got %d", len(expectedPatterns), len(patterns))
	}
}

func TestListPatterns25D(t *testing.T) {
	list := ListPatterns25D()

	expectedCount := 6
	if len(list) != expectedCount {
		t.Errorf("Expected %d patterns in list, got %d", expectedCount, len(list))
	}

	// Verify all listed patterns exist in GetPatterns25D
	allPatterns := GetPatterns25D()
	for _, name := range list {
		if _, exists := allPatterns[name]; !exists {
			t.Errorf("Listed pattern '%s' not found in GetPatterns25D()", name)
		}
	}
}

func TestPattern25D_CreateUniverse(t *testing.T) {
	patterns := GetPatterns25D()

	for name, pattern := range patterns {
		u := pattern.CreateUniverse(rules.ConwayRule{})

		if u == nil {
			t.Errorf("Pattern '%s' CreateUniverse returned nil", name)
			continue
		}

		size := u.Size()
		if size.X != pattern.Width || size.Y != pattern.Height || size.Z != pattern.Depth {
			t.Errorf("Pattern '%s': universe size mismatch", name)
		}

		if u.CountLiving() == 0 {
			t.Errorf("Pattern '%s' has no living cells after loading", name)
		}
	}
}

func TestPattern25D_Evolution(t *testing.T) {
	// Test that patterns evolve without crashing
	patterns := GetPatterns25D()

	for name, pattern := range patterns {
		u := pattern.CreateUniverse(rules.ConwayRule{})
		u.SetLayerInteraction(true)

		initialCount := u.CountLiving()

		// Run 10 generations
		for i := 0; i < 10; i++ {
			u.Step()
		}

		// Just verify it doesn't crash and produces some activity
		// (patterns may die out or explode, both are valid behaviors)
		t.Logf("Pattern '%s': initial=%d, after 10 steps=%d", name, initialCount, u.CountLiving())
	}
}
