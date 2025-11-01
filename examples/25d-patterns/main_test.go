package main

import (
	"bytes"
	"golife/pkg/patterns"
	"golife/pkg/rules"
	"golife/pkg/visualizer/terminal"
	"testing"
)

func TestPatternCatalog(t *testing.T) {
	patternList := patterns.ListPatterns25D()
	allPatterns := patterns.GetPatterns25D()

	if len(patternList) == 0 {
		t.Fatal("Pattern list should not be empty")
	}

	if len(allPatterns) == 0 {
		t.Fatal("All patterns map should not be empty")
	}

	// Verify all patterns in list exist in map
	for _, name := range patternList {
		if _, exists := allPatterns[name]; !exists {
			t.Errorf("Pattern '%s' in list but not in map", name)
		}
	}
}

func TestPatternProperties(t *testing.T) {
	allPatterns := patterns.GetPatterns25D()

	for name, pattern := range allPatterns {
		t.Run(name, func(t *testing.T) {
			// Verify pattern has a name
			if pattern.Name == "" {
				t.Error("Pattern should have a name")
			}

			// Verify pattern has a description
			if pattern.Description == "" {
				t.Error("Pattern should have a description")
			}

			// Verify positive dimensions
			if pattern.Width <= 0 {
				t.Errorf("Width should be positive, got %d", pattern.Width)
			}
			if pattern.Height <= 0 {
				t.Errorf("Height should be positive, got %d", pattern.Height)
			}
			if pattern.Depth <= 0 {
				t.Errorf("Depth should be positive, got %d", pattern.Depth)
			}

			// Verify has cells
			if len(pattern.Cells) == 0 {
				t.Error("Pattern should have at least one cell")
			}
		})
	}
}

func TestPatternCreation(t *testing.T) {
	allPatterns := patterns.GetPatterns25D()
	rule := rules.ConwayRule{}

	for name, pattern := range allPatterns {
		t.Run(name, func(t *testing.T) {
			u := pattern.CreateUniverse(rule)

			size := u.Size()
			if size.X != pattern.Width {
				t.Errorf("Universe width: got %d, want %d", size.X, pattern.Width)
			}
			if size.Y != pattern.Height {
				t.Errorf("Universe height: got %d, want %d", size.Y, pattern.Height)
			}
			if size.Z != pattern.Depth {
				t.Errorf("Universe depth: got %d, want %d", size.Z, pattern.Depth)
			}

			// Verify cells were loaded
			livingCells := u.CountLiving()
			expectedCells := len(pattern.Cells)
			if livingCells != expectedCells {
				t.Errorf("Living cells: got %d, want %d", livingCells, expectedCells)
			}
		})
	}
}

func TestPatternEvolution(t *testing.T) {
	allPatterns := patterns.GetPatterns25D()
	rule := rules.ConwayRule{}

	for name, pattern := range allPatterns {
		t.Run(name, func(t *testing.T) {
			u := pattern.CreateUniverse(rule)
			u.SetLayerInteraction(true)
			u.SetVerticalWeight(0.3)

			initialPop := u.CountLiving()

			// Evolve for a few generations
			for i := 0; i < 3; i++ {
				u.Step()
			}

			finalPop := u.CountLiving()

			// Population should either stay positive or go to zero
			if finalPop < 0 {
				t.Errorf("Population should not be negative: %d", finalPop)
			}

			// Most patterns should have some activity
			if initialPop > 0 && finalPop == 0 {
				t.Logf("Pattern '%s' died out after 3 generations", name)
			}
		})
	}
}

func TestMultiLayerViewWithPatterns(t *testing.T) {
	allPatterns := patterns.GetPatterns25D()
	rule := rules.ConwayRule{}
	view := terminal.NewMultiLayerView()
	view.SetLayout(terminal.GridLayout)
	view.ToggleAllLayers()

	for name, pattern := range allPatterns {
		t.Run(name, func(t *testing.T) {
			u := pattern.CreateUniverse(rule)

			// Render should produce output
			output := view.Render(u)
			if output == "" {
				t.Error("Render output should not be empty")
			}

			// Stats should produce output
			stats := view.RenderStats(0, u.CountLiving(), 30.0)
			if stats == "" {
				t.Error("Stats output should not be empty")
			}
		})
	}
}

func TestSpecificPatterns(t *testing.T) {
	allPatterns := patterns.GetPatterns25D()

	tests := []struct {
		name       string
		patternKey string
		checkFunc  func(*testing.T, *patterns.Pattern25D)
	}{
		{
			name:       "layer-sandwich exists",
			patternKey: "layer-sandwich",
			checkFunc: func(t *testing.T, p *patterns.Pattern25D) {
				if p.Depth < 2 {
					t.Errorf("layer-sandwich should have at least 2 layers, got %d", p.Depth)
				}
			},
		},
		{
			name:       "layer-oscillator exists",
			patternKey: "layer-oscillator",
			checkFunc: func(t *testing.T, p *patterns.Pattern25D) {
				if len(p.Cells) == 0 {
					t.Error("layer-oscillator should have cells")
				}
			},
		},
		{
			name:       "energy-wave exists",
			patternKey: "energy-wave",
			checkFunc: func(t *testing.T, p *patterns.Pattern25D) {
				if len(p.Cells) == 0 {
					t.Error("energy-wave should have cells")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pattern, exists := allPatterns[tt.patternKey]
			if !exists {
				t.Fatalf("Pattern '%s' not found", tt.patternKey)
			}
			tt.checkFunc(t, pattern)
		})
	}
}

func TestLayerInteractionModes(t *testing.T) {
	allPatterns := patterns.GetPatterns25D()
	pattern := allPatterns["layer-oscillator"]
	rule := rules.ConwayRule{}

	tests := []struct {
		name        string
		interaction bool
		weight      float64
	}{
		{
			name:        "no interaction",
			interaction: false,
			weight:      0.0,
		},
		{
			name:        "low weight interaction",
			interaction: true,
			weight:      0.3,
		},
		{
			name:        "high weight interaction",
			interaction: true,
			weight:      0.8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := pattern.CreateUniverse(rule)
			u.SetLayerInteraction(tt.interaction)
			if tt.interaction {
				u.SetVerticalWeight(tt.weight)
			}

			initialPop := u.CountLiving()

			// Evolve for 3 generations
			for i := 0; i < 3; i++ {
				u.Step()
			}

			finalPop := u.CountLiving()

			// Population should be non-negative
			if finalPop < 0 {
				t.Errorf("Population should not be negative: %d", finalPop)
			}

			// Results should vary based on interaction mode
			t.Logf("Initial: %d, Final: %d (interaction=%v, weight=%.1f)",
				initialPop, finalPop, tt.interaction, tt.weight)
		})
	}
}

func TestEnergyDiffusionRule(t *testing.T) {
	allPatterns := patterns.GetPatterns25D()
	pattern := allPatterns["energy-wave"]
	baseRule := rules.ConwayRule{}

	u := pattern.CreateUniverse(baseRule)
	u.SetLayerInteraction(true)

	// Apply energy diffusion rule
	energyRule := rules.NewEnergyDiffusionRule(baseRule, 0.5, 10)
	u.SetInteractionRule(energyRule)

	initialPop := u.CountLiving()

	// Evolve
	for i := 0; i < 5; i++ {
		u.Step()
	}

	finalPop := u.CountLiving()

	// Population should change with energy diffusion
	if finalPop < 0 {
		t.Errorf("Population should not be negative: %d", finalPop)
	}

	t.Logf("Energy diffusion: Initial=%d, Final=%d", initialPop, finalPop)
}

func TestBirthBetweenLayersRule(t *testing.T) {
	allPatterns := patterns.GetPatterns25D()
	pattern := allPatterns["layer-oscillator"]
	baseRule := rules.ConwayRule{}

	u := pattern.CreateUniverse(baseRule)
	u.SetLayerInteraction(true)

	// Apply birth between layers rule
	birthRule := rules.NewBirthBetweenLayersRule(baseRule, false)
	u.SetInteractionRule(birthRule)

	initialPop := u.CountLiving()

	// Evolve
	for i := 0; i < 3; i++ {
		u.Step()
	}

	finalPop := u.CountLiving()

	// Population should be non-negative
	if finalPop < 0 {
		t.Errorf("Population should not be negative: %d", finalPop)
	}

	t.Logf("Birth between layers: Initial=%d, Final=%d", initialPop, finalPop)
}

func TestRenderOutput(t *testing.T) {
	allPatterns := patterns.GetPatterns25D()
	view := terminal.NewMultiLayerView()
	rule := rules.ConwayRule{}

	for name, pattern := range allPatterns {
		t.Run(name, func(t *testing.T) {
			u := pattern.CreateUniverse(rule)

			// Test different layouts
			layouts := []terminal.LayoutType{
				terminal.HorizontalLayout,
				terminal.VerticalLayout,
				terminal.GridLayout,
			}

			view.ToggleAllLayers()

			for _, layout := range layouts {
				view.SetLayout(layout)
				output := view.Render(u)

				if output == "" {
					t.Errorf("Render output for layout %v should not be empty", layout)
				}

				// Output should contain some visual representation
				if len(output) < 10 {
					t.Errorf("Render output seems too short: %d chars", len(output))
				}
			}
		})
	}
}

func TestStatsRendering(t *testing.T) {
	allPatterns := patterns.GetPatterns25D()
	view := terminal.NewMultiLayerView()
	rule := rules.ConwayRule{}

	pattern := allPatterns["layer-oscillator"]
	u := pattern.CreateUniverse(rule)

	stats := view.RenderStats(5, u.CountLiving(), 30.0)

	if stats == "" {
		t.Error("Stats output should not be empty")
	}

	// Stats should contain key information
	if !bytes.Contains([]byte(stats), []byte("Gen")) {
		t.Error("Stats should contain 'Gen'")
	}
	if !bytes.Contains([]byte(stats), []byte("Living")) {
		t.Error("Stats should contain 'Living'")
	}
	if !bytes.Contains([]byte(stats), []byte("FPS")) {
		t.Error("Stats should contain 'FPS'")
	}
}
