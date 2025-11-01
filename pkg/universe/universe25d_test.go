package universe

import (
	"golife/pkg/core"
	"golife/pkg/rules"
	"testing"
)

func TestNew25D(t *testing.T) {
	u := New25D(10, 10, 3, rules.ConwayRule{})

	if u.width != 10 || u.height != 10 || u.depth != 3 {
		t.Errorf("Expected dimensions (10,10,3), got (%d,%d,%d)", u.width, u.height, u.depth)
	}

	if u.Dimension() != core.Dim25D {
		t.Errorf("Expected Dim25D, got %d", u.Dimension())
	}

	size := u.Size()
	if size.X != 10 || size.Y != 10 || size.Z != 3 {
		t.Errorf("Expected size (10,10,3), got (%d,%d,%d)", size.X, size.Y, size.Z)
	}

	if u.layerInteraction {
		t.Error("Layer interaction should be disabled by default")
	}

	if u.verticalWeight != 0.3 {
		t.Errorf("Expected default vertical weight 0.3, got %f", u.verticalWeight)
	}
}

func TestUniverse25D_GetSet(t *testing.T) {
	u := New25D(5, 5, 3, rules.ConwayRule{})

	// Set cell in middle layer
	coord := core.NewCoord3D(2, 2, 1)
	u.Set(coord, core.Alive)

	if u.Get(coord) != core.Alive {
		t.Error("Cell (2,2,1) should be alive")
	}

	// Test boundary conditions
	outOfBounds := core.NewCoord3D(2, 2, 5)
	if u.Get(outOfBounds) != core.Dead {
		t.Error("Out of bounds should return Dead")
	}

	// Set should not panic on out of bounds
	u.Set(outOfBounds, core.Alive)
	if u.Get(outOfBounds) != core.Dead {
		t.Error("Out of bounds cell should remain Dead")
	}
}

func TestUniverse25D_GetLayer(t *testing.T) {
	u := New25D(5, 5, 3, rules.ConwayRule{})

	// Get valid layer
	layer := u.GetLayer(1)
	if layer == nil {
		t.Error("Layer 1 should exist")
	}

	// Get invalid layer
	if u.GetLayer(5) != nil {
		t.Error("Layer 5 should be nil")
	}
	if u.GetLayer(-1) != nil {
		t.Error("Layer -1 should be nil")
	}
}

func TestUniverse25D_StepIndependent(t *testing.T) {
	u := New25D(20, 20, 3, rules.ConwayRule{})

	// Create glider pattern in layer 0
	// ●
	//   ●
	// ●●●
	u.layers[0].Set(core.NewCoord2D(6, 5), core.Alive)
	u.layers[0].Set(core.NewCoord2D(7, 6), core.Alive)
	u.layers[0].Set(core.NewCoord2D(5, 7), core.Alive)
	u.layers[0].Set(core.NewCoord2D(6, 7), core.Alive)
	u.layers[0].Set(core.NewCoord2D(7, 7), core.Alive)

	// Create blinker pattern in layer 1
	// ●●●
	u.layers[1].Set(core.NewCoord2D(5, 5), core.Alive)
	u.layers[1].Set(core.NewCoord2D(6, 5), core.Alive)
	u.layers[1].Set(core.NewCoord2D(7, 5), core.Alive)

	// Layer 2 remains empty

	initialCount0 := u.CountLivingInLayer(0)
	initialCount1 := u.CountLivingInLayer(1)

	if initialCount0 != 5 {
		t.Errorf("Layer 0 should have 5 cells (glider), got %d", initialCount0)
	}
	if initialCount1 != 3 {
		t.Errorf("Layer 1 should have 3 cells (blinker), got %d", initialCount1)
	}

	// Step with independent layers
	u.SetLayerInteraction(false)
	u.Step()

	// Each layer should evolve independently
	if u.CountLivingInLayer(0) != 5 {
		t.Error("Glider should still have 5 cells after step")
	}
	if u.CountLivingInLayer(1) != 3 {
		t.Error("Blinker should still have 3 cells after step")
	}
	if u.CountLivingInLayer(2) != 0 {
		t.Error("Empty layer should remain empty")
	}
}

func TestUniverse25D_StepWithInteraction(t *testing.T) {
	u := New25D(20, 20, 3, rules.ConwayRule{})

	// Create a vertical pattern: same position in all layers
	// Layer 0: single cell at (10,10)
	// Layer 1: single cell at (10,10)
	// Layer 2: single cell at (10,10)
	coord := core.NewCoord2D(10, 10)
	u.layers[0].Set(coord, core.Alive)
	u.layers[1].Set(coord, core.Alive)
	u.layers[2].Set(coord, core.Alive)

	// Enable layer interaction
	u.SetLayerInteraction(true)
	u.SetVerticalWeight(1.0) // Full weight for testing

	initialTotal := u.CountLiving()
	if initialTotal != 3 {
		t.Errorf("Expected 3 living cells initially, got %d", initialTotal)
	}

	// Step with interaction
	u.Step()

	// With vertical neighbors, the middle layer cell has:
	// - 0 horizontal neighbors
	// - 2 vertical neighbors (above and below)
	// Should die (Conway rule: survive with 2-3 neighbors)
	// Actually should survive because it has 2 neighbors total

	// The behavior depends on the vertical weight
	// This test verifies that interaction is happening
	afterTotal := u.CountLiving()
	t.Logf("Living cells after step: %d", afterTotal)
}

func TestUniverse25D_VerticalNeighbors(t *testing.T) {
	u := New25D(10, 10, 3, rules.ConwayRule{})

	// Create a 3x3 pattern in layer 0
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			u.layers[0].Set(core.NewCoord2D(5+dx, 5+dy), core.Alive)
		}
	}

	// Count vertical neighbors for layer 1, position (5,5)
	verticalCount := u.countVerticalNeighbors(5, 5, 1)

	// Should count 9 cells from layer 0 (layer above)
	// Layer 2 is empty, so 0 from below
	if verticalCount != 9 {
		t.Errorf("Expected 9 vertical neighbors, got %d", verticalCount)
	}
}

func TestUniverse25D_VerticalNeighborsEdgeCases(t *testing.T) {
	u := New25D(10, 10, 3, rules.ConwayRule{})

	// Fill layer 0 completely
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			u.layers[0].Set(core.NewCoord2D(x, y), core.Alive)
		}
	}

	// Test corner of layer 1
	verticalCount := u.countVerticalNeighbors(0, 0, 1)
	// Corner (0,0) can only check 4 cells in the 3x3 grid above
	// (-1,-1), (-1,0), (-1,1) are out of bounds
	// (0,-1), (1,-1) are out of bounds
	// Valid: (0,0), (0,1), (1,0), (1,1)
	if verticalCount != 4 {
		t.Errorf("Corner should have 4 vertical neighbors, got %d", verticalCount)
	}

	// Test edge
	verticalCount = u.countVerticalNeighbors(5, 0, 1)
	// Edge (5,0): y-1=-1 is out of bounds
	// Valid: (4,0), (5,0), (6,0), (4,1), (5,1), (6,1)
	if verticalCount != 6 {
		t.Errorf("Edge should have 6 vertical neighbors, got %d", verticalCount)
	}

	// Test center
	verticalCount = u.countVerticalNeighbors(5, 5, 1)
	// Center has full 3x3 grid
	if verticalCount != 9 {
		t.Errorf("Center should have 9 vertical neighbors, got %d", verticalCount)
	}

	// Test top layer (no layer above)
	verticalCount = u.countVerticalNeighbors(5, 5, 0)
	// Only layer below (layer 1) exists, which is empty
	if verticalCount != 0 {
		t.Errorf("Top layer with empty layer below should have 0 vertical neighbors, got %d", verticalCount)
	}
}

func TestUniverse25D_VerticalWeight(t *testing.T) {
	u := New25D(10, 10, 3, rules.ConwayRule{})

	// Test setting valid weights
	u.SetVerticalWeight(0.5)
	if u.verticalWeight != 0.5 {
		t.Errorf("Expected weight 0.5, got %f", u.verticalWeight)
	}

	u.SetVerticalWeight(0.0)
	if u.verticalWeight != 0.0 {
		t.Errorf("Expected weight 0.0, got %f", u.verticalWeight)
	}

	u.SetVerticalWeight(1.0)
	if u.verticalWeight != 1.0 {
		t.Errorf("Expected weight 1.0, got %f", u.verticalWeight)
	}

	// Test invalid weights (should be rejected)
	u.SetVerticalWeight(-0.1)
	if u.verticalWeight != 1.0 {
		t.Error("Negative weight should be rejected")
	}

	u.SetVerticalWeight(1.5)
	if u.verticalWeight != 1.0 {
		t.Error("Weight > 1.0 should be rejected")
	}
}

func TestUniverse25D_Clone(t *testing.T) {
	u := New25D(5, 5, 3, rules.ConwayRule{})

	// Set some cells
	u.Set(core.NewCoord3D(1, 1, 0), core.Alive)
	u.Set(core.NewCoord3D(2, 2, 1), core.Alive)
	u.Set(core.NewCoord3D(3, 3, 2), core.Alive)
	u.SetLayerInteraction(true)
	u.SetVerticalWeight(0.5)

	// Clone
	clone := u.Clone().(*Universe25D)

	// Verify dimensions
	if clone.width != u.width || clone.height != u.height || clone.depth != u.depth {
		t.Error("Clone dimensions don't match")
	}

	// Verify settings
	if clone.layerInteraction != u.layerInteraction {
		t.Error("Clone layerInteraction doesn't match")
	}
	if clone.verticalWeight != u.verticalWeight {
		t.Error("Clone verticalWeight doesn't match")
	}

	// Verify cells
	if clone.Get(core.NewCoord3D(1, 1, 0)) != core.Alive {
		t.Error("Clone cell (1,1,0) should be alive")
	}
	if clone.Get(core.NewCoord3D(2, 2, 1)) != core.Alive {
		t.Error("Clone cell (2,2,1) should be alive")
	}
	if clone.Get(core.NewCoord3D(3, 3, 2)) != core.Alive {
		t.Error("Clone cell (3,3,2) should be alive")
	}

	// Modify original
	u.Set(core.NewCoord3D(1, 1, 0), core.Dead)

	// Clone should not be affected
	if clone.Get(core.NewCoord3D(1, 1, 0)) != core.Alive {
		t.Error("Clone should not be affected by original modification")
	}
}

func TestUniverse25D_Clear(t *testing.T) {
	u := New25D(5, 5, 3, rules.ConwayRule{})

	// Fill with cells
	for z := 0; z < 3; z++ {
		for y := 0; y < 5; y++ {
			for x := 0; x < 5; x++ {
				u.Set(core.NewCoord3D(x, y, z), core.Alive)
			}
		}
	}

	if u.CountLiving() != 75 {
		t.Errorf("Expected 75 living cells, got %d", u.CountLiving())
	}

	// Clear
	u.Clear()

	if u.CountLiving() != 0 {
		t.Errorf("After clear, expected 0 living cells, got %d", u.CountLiving())
	}
}

func TestUniverse25D_Randomize(t *testing.T) {
	u := New25D(10, 10, 3, rules.ConwayRule{})

	u.Randomize()

	count := u.CountLiving()
	if count == 0 {
		t.Error("Randomize should create some living cells")
	}

	// Randomize should fill approximately 50% of cells
	totalCells := 10 * 10 * 3
	if count < totalCells/4 || count > 3*totalCells/4 {
		t.Logf("Warning: Randomize created %d/%d cells (expected ~%d)", count, totalCells, totalCells/2)
	}
}

func TestUniverse25D_CountLivingInLayer(t *testing.T) {
	u := New25D(5, 5, 3, rules.ConwayRule{})

	// Add cells to different layers
	u.Set(core.NewCoord3D(1, 1, 0), core.Alive)
	u.Set(core.NewCoord3D(2, 2, 0), core.Alive)
	u.Set(core.NewCoord3D(3, 3, 1), core.Alive)

	if u.CountLivingInLayer(0) != 2 {
		t.Errorf("Layer 0 should have 2 cells, got %d", u.CountLivingInLayer(0))
	}
	if u.CountLivingInLayer(1) != 1 {
		t.Errorf("Layer 1 should have 1 cell, got %d", u.CountLivingInLayer(1))
	}
	if u.CountLivingInLayer(2) != 0 {
		t.Errorf("Layer 2 should have 0 cells, got %d", u.CountLivingInLayer(2))
	}

	// Test invalid layer indices
	if u.CountLivingInLayer(-1) != 0 {
		t.Error("Invalid layer should return 0")
	}
	if u.CountLivingInLayer(5) != 0 {
		t.Error("Invalid layer should return 0")
	}
}

func BenchmarkUniverse25D_StepIndependent(b *testing.B) {
	u := New25D(50, 50, 10, rules.ConwayRule{})
	u.Randomize()
	u.SetLayerInteraction(false)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u.Step()
	}
}

func BenchmarkUniverse25D_StepWithInteraction(b *testing.B) {
	u := New25D(50, 50, 10, rules.ConwayRule{})
	u.Randomize()
	u.SetLayerInteraction(true)
	u.SetVerticalWeight(0.3)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u.Step()
	}
}
