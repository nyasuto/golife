package universe

import (
	"golife/pkg/core"
	"golife/pkg/rules"
	"testing"
)

func TestUniverse25D_WithWeightedNeighborsRule(t *testing.T) {
	u := New25D(20, 20, 3, rules.ConwayRule{})

	// Set custom weighted rule
	customRule := rules.NewWeightedNeighborsRule(rules.ConwayRule{}, 0.8)
	u.SetInteractionRule(customRule)
	u.SetLayerInteraction(true)

	// Verify rule was set
	if u.GetInteractionRule().Type() != rules.WeightedNeighbors {
		t.Error("Interaction rule type should be WeightedNeighbors")
	}

	// Create a simple pattern
	u.layers[1].Set(core.NewCoord2D(10, 10), core.Alive)
	u.layers[1].Set(core.NewCoord2D(11, 10), core.Alive)
	u.layers[1].Set(core.NewCoord2D(10, 11), core.Alive)

	initialCount := u.CountLiving()
	if initialCount != 3 {
		t.Errorf("Expected 3 initial cells, got %d", initialCount)
	}

	// Step with weighted rule
	u.Step()

	// Pattern should evolve
	afterCount := u.CountLiving()
	t.Logf("Cells after step with weighted rule: %d", afterCount)

	// Verify interaction is working
	if afterCount == 0 {
		t.Error("All cells died - interaction may not be working")
	}
}

func TestUniverse25D_WithBirthBetweenLayersRule(t *testing.T) {
	u := New25D(20, 20, 3, rules.ConwayRule{})

	// Set birth between layers rule (requires at least one adjacent layer)
	customRule := rules.NewBirthBetweenLayersRule(rules.ConwayRule{}, false)
	u.SetInteractionRule(customRule)
	u.SetLayerInteraction(true)

	// Create pattern with cells in layer 0 and 2, but empty layer 1
	// This should allow birth in layer 1 if horizontal neighbors are right
	u.layers[0].Set(core.NewCoord2D(10, 10), core.Alive)
	u.layers[2].Set(core.NewCoord2D(10, 10), core.Alive)

	// Add pattern in layer 1 that would normally birth (glider)
	u.layers[1].Set(core.NewCoord2D(10, 9), core.Alive)
	u.layers[1].Set(core.NewCoord2D(11, 10), core.Alive)
	u.layers[1].Set(core.NewCoord2D(9, 11), core.Alive)
	u.layers[1].Set(core.NewCoord2D(10, 11), core.Alive)
	u.layers[1].Set(core.NewCoord2D(11, 11), core.Alive)

	initialCount := u.CountLiving()
	t.Logf("Initial cells: %d", initialCount)

	u.Step()

	afterCount := u.CountLiving()
	t.Logf("Cells after step with birth-between-layers rule: %d", afterCount)

	// The pattern should evolve with layer constraints
	if afterCount == 0 {
		t.Error("All cells died")
	}
}

func TestUniverse25D_WithBirthBetweenLayersRule_BothRequired(t *testing.T) {
	u := New25D(20, 20, 3, rules.ConwayRule{})

	// Set birth between layers rule (requires both adjacent layers)
	customRule := rules.NewBirthBetweenLayersRule(rules.ConwayRule{}, true)
	u.SetInteractionRule(customRule)
	u.SetLayerInteraction(true)

	// Pattern in middle layer with cells in both adjacent layers
	u.layers[0].Set(core.NewCoord2D(10, 10), core.Alive)
	u.layers[1].Set(core.NewCoord2D(10, 9), core.Alive)
	u.layers[1].Set(core.NewCoord2D(11, 10), core.Alive)
	u.layers[1].Set(core.NewCoord2D(9, 11), core.Alive)
	u.layers[1].Set(core.NewCoord2D(10, 11), core.Alive)
	u.layers[1].Set(core.NewCoord2D(11, 11), core.Alive)
	u.layers[2].Set(core.NewCoord2D(10, 10), core.Alive)

	initialCount := u.CountLiving()
	t.Logf("Initial cells: %d", initialCount)

	u.Step()

	afterCount := u.CountLiving()
	t.Logf("Cells after step with both-layers-required rule: %d", afterCount)

	// Pattern should be more constrained
	if afterCount == 0 {
		t.Error("All cells died")
	}
}

func TestUniverse25D_WithEnergyDiffusionRule(t *testing.T) {
	u := New25D(20, 20, 3, rules.ConwayRule{})

	// Set energy diffusion rule
	customRule := rules.NewEnergyDiffusionRule(rules.ConwayRule{}, 0.5, 10)
	u.SetInteractionRule(customRule)
	u.SetLayerInteraction(true)

	// Create high-energy pattern in top layer
	for y := 9; y <= 11; y++ {
		for x := 9; x <= 11; x++ {
			u.layers[0].Set(core.NewCoord2D(x, y), 100) // High energy
		}
	}

	// Middle layer has some cells
	u.layers[1].Set(core.NewCoord2D(10, 10), core.Alive)

	initialEnergy := 0
	for z := 0; z < 3; z++ {
		for y := 0; y < 20; y++ {
			for x := 0; x < 20; x++ {
				state := u.layers[z].Get(core.NewCoord2D(x, y))
				initialEnergy += int(state)
			}
		}
	}

	t.Logf("Initial total energy: %d", initialEnergy)

	u.Step()

	afterEnergy := 0
	for z := 0; z < 3; z++ {
		for y := 0; y < 20; y++ {
			for x := 0; x < 20; x++ {
				state := u.layers[z].Get(core.NewCoord2D(x, y))
				afterEnergy += int(state)
			}
		}
	}

	t.Logf("Total energy after step: %d", afterEnergy)

	// Energy should diffuse between layers
	// Note: This test just verifies the rule runs without error
	// Actual energy behavior depends on complex interactions
}

func TestUniverse25D_SwitchBetweenRules(t *testing.T) {
	u := New25D(10, 10, 3, rules.ConwayRule{})

	// Start with weighted rule
	rule1 := rules.NewWeightedNeighborsRule(rules.ConwayRule{}, 0.3)
	u.SetInteractionRule(rule1)

	if u.GetInteractionRule().Type() != rules.WeightedNeighbors {
		t.Error("Should have WeightedNeighbors rule")
	}

	// Switch to birth-between-layers
	rule2 := rules.NewBirthBetweenLayersRule(rules.ConwayRule{}, false)
	u.SetInteractionRule(rule2)

	if u.GetInteractionRule().Type() != rules.BirthBetweenLayers {
		t.Error("Should have BirthBetweenLayers rule")
	}

	// Switch to energy diffusion
	rule3 := rules.NewEnergyDiffusionRule(rules.ConwayRule{}, 0.5, 10)
	u.SetInteractionRule(rule3)

	if u.GetInteractionRule().Type() != rules.EnergyDiffusion {
		t.Error("Should have EnergyDiffusion rule")
	}
}

func TestUniverse25D_BackwardCompatibility(t *testing.T) {
	// Test that old SetVerticalWeight still works
	u := New25D(10, 10, 3, rules.ConwayRule{})

	u.SetVerticalWeight(0.7)
	u.SetLayerInteraction(true)

	// Verify the interaction rule was updated
	if u.GetInteractionRule().Type() != rules.WeightedNeighbors {
		t.Error("SetVerticalWeight should update interaction rule to WeightedNeighbors")
	}

	// Add pattern and step
	u.layers[1].Set(core.NewCoord2D(5, 5), core.Alive)
	u.layers[1].Set(core.NewCoord2D(6, 5), core.Alive)
	u.layers[1].Set(core.NewCoord2D(5, 6), core.Alive)

	u.Step()

	// Should not panic or error
	afterCount := u.CountLiving()
	t.Logf("Cells after step with backward-compatible API: %d", afterCount)
}

func BenchmarkUniverse25D_WeightedNeighborsRule(b *testing.B) {
	u := New25D(50, 50, 10, rules.ConwayRule{})
	u.Randomize()
	rule := rules.NewWeightedNeighborsRule(rules.ConwayRule{}, 0.3)
	u.SetInteractionRule(rule)
	u.SetLayerInteraction(true)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u.Step()
	}
}

func BenchmarkUniverse25D_BirthBetweenLayersRule(b *testing.B) {
	u := New25D(50, 50, 10, rules.ConwayRule{})
	u.Randomize()
	rule := rules.NewBirthBetweenLayersRule(rules.ConwayRule{}, false)
	u.SetInteractionRule(rule)
	u.SetLayerInteraction(true)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u.Step()
	}
}

func BenchmarkUniverse25D_EnergyDiffusionRule(b *testing.B) {
	u := New25D(50, 50, 10, rules.ConwayRule{})
	u.Randomize()
	rule := rules.NewEnergyDiffusionRule(rules.ConwayRule{}, 0.5, 10)
	u.SetInteractionRule(rule)
	u.SetLayerInteraction(true)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u.Step()
	}
}
