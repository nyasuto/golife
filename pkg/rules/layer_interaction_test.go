package rules

import (
	"golife/pkg/core"
	"testing"
)

// Mock rule for testing
type mockRule struct{}

func (r mockRule) Name() string { return "Mock" }
func (r mockRule) ShouldBirth(neighborCount int) bool {
	return neighborCount == 3
}
func (r mockRule) ShouldSurvive(neighborCount int, currentState core.CellState) bool {
	if currentState == core.Dead {
		return false
	}
	return neighborCount == 2 || neighborCount == 3
}
func (r mockRule) NeighborWeight(distance float64) float64 { return 1.0 }

func TestWeightedNeighborsRule(t *testing.T) {
	baseRule := mockRule{}
	rule := NewWeightedNeighborsRule(baseRule, 0.5)

	if rule.Type() != WeightedNeighbors {
		t.Errorf("Expected WeightedNeighbors type, got %d", rule.Type())
	}

	// Test neighbor calculation: 3 horizontal + 10 vertical × 0.5 = 3 + 5 = 8
	neighborCount := rule.CalculateNeighborCount(3, 10, core.Dead, core.Alive, core.Alive)
	if neighborCount != 8 {
		t.Errorf("Expected 8 neighbors, got %d", neighborCount)
	}

	// Test birth rule (should follow base rule)
	if !rule.ShouldBirth(3, core.Dead, core.Dead) {
		t.Error("Should birth with 3 neighbors")
	}
	if rule.ShouldBirth(2, core.Dead, core.Dead) {
		t.Error("Should not birth with 2 neighbors")
	}

	// Test survival rule
	if !rule.ShouldSurvive(2, core.Alive, core.Dead, core.Dead) {
		t.Error("Should survive with 2 neighbors")
	}
	if rule.ShouldSurvive(1, core.Alive, core.Dead, core.Dead) {
		t.Error("Should not survive with 1 neighbor")
	}
}

func TestWeightedNeighborsRule_WeightBounds(t *testing.T) {
	baseRule := mockRule{}

	// Test weight < 0
	rule1 := NewWeightedNeighborsRule(baseRule, -0.5)
	count1 := rule1.CalculateNeighborCount(5, 10, core.Dead, core.Dead, core.Dead)
	if count1 != 5 {
		t.Errorf("Negative weight should be clamped to 0, expected 5, got %d", count1)
	}

	// Test weight > 1.0
	rule2 := NewWeightedNeighborsRule(baseRule, 1.5)
	count2 := rule2.CalculateNeighborCount(5, 10, core.Dead, core.Dead, core.Dead)
	if count2 != 15 { // 5 + 10*1.0
		t.Errorf("Weight > 1.0 should be clamped to 1.0, expected 15, got %d", count2)
	}
}

func TestBirthBetweenLayersRule_SingleLayer(t *testing.T) {
	baseRule := mockRule{}
	rule := NewBirthBetweenLayersRule(baseRule, false) // Requires at least one layer

	if rule.Type() != BirthBetweenLayers {
		t.Errorf("Expected BirthBetweenLayers type, got %d", rule.Type())
	}

	// Test neighbor calculation (only horizontal counts)
	neighborCount := rule.CalculateNeighborCount(3, 10, core.Dead, core.Alive, core.Alive)
	if neighborCount != 3 {
		t.Errorf("Expected 3 neighbors, got %d", neighborCount)
	}

	// Test birth with upper layer having cells
	if !rule.ShouldBirth(3, core.Alive, core.Dead) {
		t.Error("Should birth when upper layer has cells and neighbor count is 3")
	}

	// Test birth with lower layer having cells
	if !rule.ShouldBirth(3, core.Dead, core.Alive) {
		t.Error("Should birth when lower layer has cells and neighbor count is 3")
	}

	// Test birth with no adjacent layer cells
	if rule.ShouldBirth(3, core.Dead, core.Dead) {
		t.Error("Should not birth when no adjacent layer has cells")
	}

	// Test birth with wrong neighbor count
	if rule.ShouldBirth(2, core.Alive, core.Alive) {
		t.Error("Should not birth with 2 neighbors even if layers have cells")
	}
}

func TestBirthBetweenLayersRule_BothLayers(t *testing.T) {
	baseRule := mockRule{}
	rule := NewBirthBetweenLayersRule(baseRule, true) // Requires both layers

	// Test birth with both layers having cells
	if !rule.ShouldBirth(3, core.Alive, core.Alive) {
		t.Error("Should birth when both layers have cells and neighbor count is 3")
	}

	// Test birth with only upper layer
	if rule.ShouldBirth(3, core.Alive, core.Dead) {
		t.Error("Should not birth when only upper layer has cells (requires both)")
	}

	// Test birth with only lower layer
	if rule.ShouldBirth(3, core.Dead, core.Alive) {
		t.Error("Should not birth when only lower layer has cells (requires both)")
	}

	// Test survival is not affected
	if !rule.ShouldSurvive(2, core.Alive, core.Dead, core.Dead) {
		t.Error("Survival should not require layer conditions")
	}
}

func TestEnergyDiffusionRule(t *testing.T) {
	baseRule := mockRule{}
	rule := NewEnergyDiffusionRule(baseRule, 0.5, 10)

	if rule.Type() != EnergyDiffusion {
		t.Errorf("Expected EnergyDiffusion type, got %d", rule.Type())
	}

	// Test with high energy in adjacent layers
	// upper = 100, lower = 100, diffusionRate = 0.5
	// diffused = (100 + 100) * 0.5 = 100
	// effective vertical = 100 / 10 = 10
	neighborCount := rule.CalculateNeighborCount(3, 0, core.Dead, 100, 100)
	expected := 3 + 10 // horizontal + diffused
	if neighborCount != expected {
		t.Errorf("Expected %d neighbors, got %d", expected, neighborCount)
	}

	// Test with low energy
	neighborCount2 := rule.CalculateNeighborCount(3, 0, core.Dead, 5, 5)
	// diffused = (5 + 5) * 0.5 = 5, which is below threshold 10
	if neighborCount2 != 3 {
		t.Errorf("Expected 3 neighbors (no diffusion below threshold), got %d", neighborCount2)
	}
}

func TestEnergyDiffusionRule_Bounds(t *testing.T) {
	baseRule := mockRule{}

	// Test negative diffusion rate
	rule1 := NewEnergyDiffusionRule(baseRule, -0.5, 10)
	count1 := rule1.CalculateNeighborCount(5, 0, core.Dead, 100, 100)
	if count1 != 5 {
		t.Errorf("Negative diffusion rate should be clamped to 0, expected 5, got %d", count1)
	}

	// Test diffusion rate > 1.0
	rule2 := NewEnergyDiffusionRule(baseRule, 1.5, 10)
	count2 := rule2.CalculateNeighborCount(5, 0, core.Dead, 100, 100)
	expected := 5 + ((100 + 100) * 1.0 / 10) // clamped to 1.0
	if count2 != int(expected) {
		t.Errorf("Diffusion rate > 1.0 should be clamped to 1.0, expected %d, got %d", int(expected), count2)
	}

	// Test zero threshold (should default to 1)
	rule3 := NewEnergyDiffusionRule(baseRule, 0.5, 0)
	count3 := rule3.CalculateNeighborCount(5, 0, core.Dead, 10, 10)
	// diffused = (10 + 10) * 0.5 = 10, threshold should be 1, so 10/1 = 10
	if count3 != 15 {
		t.Errorf("Zero threshold should default to 1, expected 15, got %d", count3)
	}
}

func TestEnergyDiffusionRule_GetDiffusedEnergy(t *testing.T) {
	baseRule := mockRule{}
	rule := NewEnergyDiffusionRule(baseRule, 0.5, 10)

	// Test energy diffusion
	newEnergy := rule.GetDiffusedEnergy(50, 100, 100)
	// current decays: 50 * (1 - 0.5*0.1) = 50 * 0.95 = 47.5
	// incoming: (100 + 100) * 0.5 / 2 = 50
	// total: 47.5 + 50 = 97.5 ≈ 97
	if newEnergy < 95 || newEnergy > 100 {
		t.Errorf("Expected energy around 97, got %d", newEnergy)
	}

	// Test energy cap at 255
	newEnergy2 := rule.GetDiffusedEnergy(200, 200, 200)
	if newEnergy2 != 255 {
		t.Errorf("Energy should be capped at 255, got %d", newEnergy2)
	}

	// Test with zero energy
	newEnergy3 := rule.GetDiffusedEnergy(0, 20, 20)
	// incoming: (20 + 20) * 0.5 / 2 = 10
	if newEnergy3 != 10 {
		t.Errorf("Expected energy 10 from diffusion, got %d", newEnergy3)
	}
}

func BenchmarkWeightedNeighborsRule(b *testing.B) {
	baseRule := mockRule{}
	rule := NewWeightedNeighborsRule(baseRule, 0.3)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rule.CalculateNeighborCount(5, 10, core.Alive, core.Alive, core.Alive)
		rule.ShouldBirth(3, core.Alive, core.Alive)
		rule.ShouldSurvive(2, core.Alive, core.Alive, core.Alive)
	}
}

func BenchmarkBirthBetweenLayersRule(b *testing.B) {
	baseRule := mockRule{}
	rule := NewBirthBetweenLayersRule(baseRule, false)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rule.CalculateNeighborCount(5, 10, core.Alive, core.Alive, core.Alive)
		rule.ShouldBirth(3, core.Alive, core.Alive)
		rule.ShouldSurvive(2, core.Alive, core.Alive, core.Alive)
	}
}

func BenchmarkEnergyDiffusionRule(b *testing.B) {
	baseRule := mockRule{}
	rule := NewEnergyDiffusionRule(baseRule, 0.5, 10)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rule.CalculateNeighborCount(5, 10, 50, 100, 100)
		rule.ShouldBirth(3, 100, 100)
		rule.ShouldSurvive(2, 50, 100, 100)
		rule.GetDiffusedEnergy(50, 100, 100)
	}
}
