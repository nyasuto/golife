package rules

import (
	"golife/pkg/core"
	"testing"
)

func TestLife3D_B6S567_Name(t *testing.T) {
	rule := Life3D_B6S567{}
	name := rule.Name()

	if name != "B6/S567 (3D Life)" {
		t.Errorf("Expected name 'B6/S567 (3D Life)', got '%s'", name)
	}
}

func TestLife3D_B6S567_ShouldBirth(t *testing.T) {
	rule := Life3D_B6S567{}

	tests := []struct {
		name          string
		neighborCount int
		shouldBirth   bool
	}{
		{"0 neighbors", 0, false},
		{"1 neighbor", 1, false},
		{"2 neighbors", 2, false},
		{"3 neighbors", 3, false},
		{"4 neighbors", 4, false},
		{"5 neighbors", 5, false},
		{"6 neighbors", 6, true}, // Birth only at 6
		{"7 neighbors", 7, false},
		{"8 neighbors", 8, false},
		{"26 neighbors (max)", 26, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := rule.ShouldBirth(tt.neighborCount)
			if result != tt.shouldBirth {
				t.Errorf("ShouldBirth(%d) = %v, want %v",
					tt.neighborCount, result, tt.shouldBirth)
			}
		})
	}
}

func TestLife3D_B6S567_ShouldSurvive(t *testing.T) {
	rule := Life3D_B6S567{}

	tests := []struct {
		name          string
		neighborCount int
		shouldSurvive bool
	}{
		{"0 neighbors", 0, false},
		{"1 neighbor", 1, false},
		{"2 neighbors", 2, false},
		{"3 neighbors", 3, false},
		{"4 neighbors", 4, false},
		{"5 neighbors", 5, true},  // Survival starts at 5
		{"6 neighbors", 6, true},  // Optimal
		{"7 neighbors", 7, true},  // Survival ends at 7
		{"8 neighbors", 8, false}, // Too many
		{"9 neighbors", 9, false},
		{"26 neighbors (max)", 26, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := rule.ShouldSurvive(tt.neighborCount, core.Alive)
			if result != tt.shouldSurvive {
				t.Errorf("ShouldSurvive(%d) = %v, want %v",
					tt.neighborCount, result, tt.shouldSurvive)
			}
		})
	}
}

func TestLife3D_B6S567_ShouldSurvive_DeadCell(t *testing.T) {
	rule := Life3D_B6S567{}

	// Dead cells should never survive
	for neighbors := 0; neighbors <= 26; neighbors++ {
		result := rule.ShouldSurvive(neighbors, core.Dead)
		if result {
			t.Errorf("Dead cell should not survive with %d neighbors", neighbors)
		}
	}
}

func TestLife3D_B6S567_NeighborWeight(t *testing.T) {
	rule := Life3D_B6S567{}

	// Test various distances
	distances := []float64{0.0, 1.0, 1.414, 2.0, 3.0}
	for _, dist := range distances {
		weight := rule.NeighborWeight(dist)
		if weight != 1.0 {
			t.Errorf("NeighborWeight(%f) = %f, want 1.0", dist, weight)
		}
	}
}

func TestLife3D_B6S567_RuleSymmetry(t *testing.T) {
	rule := Life3D_B6S567{}

	// B6/S567 characteristics
	// Birth at 6 is in the middle of survival range [5,7]
	// This creates stability

	// Test that birth neighbor count is within survival range
	birthCount := 6
	if !rule.ShouldSurvive(birthCount, core.Alive) {
		t.Error("Birth count should be within survival range for stability")
	}

	// Test boundary conditions
	if rule.ShouldSurvive(4, core.Alive) {
		t.Error("Should not survive with 4 neighbors (underpopulation)")
	}

	if rule.ShouldSurvive(8, core.Alive) {
		t.Error("Should not survive with 8 neighbors (overpopulation)")
	}
}

func TestLife3D_B6S567_CompareWithConway(t *testing.T) {
	life3D := Life3D_B6S567{}
	conway := ConwayRule{}

	t.Run("Different birth thresholds", func(t *testing.T) {
		// Conway: B3, Life3D: B6
		if life3D.ShouldBirth(3) {
			t.Error("Life3D should not birth at 3 (Conway does)")
		}
		if !conway.ShouldBirth(3) {
			t.Error("Conway should birth at 3")
		}

		if !life3D.ShouldBirth(6) {
			t.Error("Life3D should birth at 6")
		}
		if conway.ShouldBirth(6) {
			t.Error("Conway should not birth at 6")
		}
	})

	t.Run("Different survival ranges", func(t *testing.T) {
		// Conway: S23, Life3D: S567
		if life3D.ShouldSurvive(2, core.Alive) {
			t.Error("Life3D should not survive at 2 (Conway does)")
		}
		if !conway.ShouldSurvive(2, core.Alive) {
			t.Error("Conway should survive at 2")
		}

		if !life3D.ShouldSurvive(5, core.Alive) {
			t.Error("Life3D should survive at 5")
		}
		if conway.ShouldSurvive(5, core.Alive) {
			t.Error("Conway should not survive at 5")
		}
	})
}

func TestLife3D_B6S567_StabilityProperties(t *testing.T) {
	rule := Life3D_B6S567{}

	t.Run("Survival range contains birth point", func(t *testing.T) {
		// This property ensures newly born cells can survive
		birthCount := 6
		if !rule.ShouldSurvive(birthCount, core.Alive) {
			t.Error("Newly born cells should be able to survive")
		}
	})

	t.Run("Survival range is continuous", func(t *testing.T) {
		// Check that survival range [5,7] is continuous
		survivalStart := 5
		survivalEnd := 7

		for n := survivalStart; n <= survivalEnd; n++ {
			if !rule.ShouldSurvive(n, core.Alive) {
				t.Errorf("Survival should be continuous in range [5,7], failed at %d", n)
			}
		}
	})

	t.Run("No survival outside range", func(t *testing.T) {
		// Check boundaries
		if rule.ShouldSurvive(4, core.Alive) {
			t.Error("Should not survive below range")
		}
		if rule.ShouldSurvive(8, core.Alive) {
			t.Error("Should not survive above range")
		}
	})
}

// BenchmarkLife3D_B6S567_ShouldBirth benchmarks birth decision
func BenchmarkLife3D_B6S567_ShouldBirth(b *testing.B) {
	rule := Life3D_B6S567{}
	for i := 0; i < b.N; i++ {
		rule.ShouldBirth(6)
	}
}

// BenchmarkLife3D_B6S567_ShouldSurvive benchmarks survival decision
func BenchmarkLife3D_B6S567_ShouldSurvive(b *testing.B) {
	rule := Life3D_B6S567{}
	for i := 0; i < b.N; i++ {
		rule.ShouldSurvive(6, core.Alive)
	}
}
