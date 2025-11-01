package main

import (
	"fmt"
	"golife/pkg/patterns"
	"golife/pkg/rules"
	"golife/pkg/visualizer/terminal"
	"time"
)

func main() {
	fmt.Println("=== 2.5D Pattern Catalog ===")
	fmt.Println()

	patternList := patterns.ListPatterns25D()
	allPatterns := patterns.GetPatterns25D()

	// Display catalog
	fmt.Println("Available 2.5D Patterns:")
	for i, name := range patternList {
		p := allPatterns[name]
		fmt.Printf("%d. %s\n", i+1, p.Name)
		fmt.Printf("   %s\n", p.Description)
		fmt.Printf("   Size: %dx%dx%d\n", p.Width, p.Height, p.Depth)
		fmt.Printf("   Cells: %d\n", len(p.Cells))
		fmt.Println()
	}

	// Demo each pattern
	view := terminal.NewMultiLayerView()
	view.SetLayout(terminal.GridLayout)
	view.ToggleAllLayers()

	for _, name := range patternList {
		pattern := allPatterns[name]
		fmt.Printf("=== Demonstrating: %s ===\n", pattern.Name)
		fmt.Println()

		// Create universe with pattern
		u := pattern.CreateUniverse(rules.ConwayRule{})

		// Enable layer interaction for most patterns
		if name != "layer-sandwich" {
			u.SetLayerInteraction(true)

			// Use different rules for different patterns
			switch name {
			case "energy-wave":
				// Use energy diffusion for energy wave
				energyRule := rules.NewEnergyDiffusionRule(rules.ConwayRule{}, 0.5, 10)
				u.SetInteractionRule(energyRule)
			case "vertical-glider":
				// Use higher vertical weight for vertical movement
				u.SetVerticalWeight(0.8)
			default:
				// Default weighted neighbors
				u.SetVerticalWeight(0.3)
			}
		}

		// Show initial state
		fmt.Println("Generation 0:")
		fmt.Println(view.Render(u))
		fmt.Println(view.RenderStats(0, u.CountLiving(), 0))
		fmt.Println()

		time.Sleep(1 * time.Second)

		// Evolve for several generations
		generations := 5
		if name == "layer-sandwich" || name == "layer-stack" {
			generations = 10 // These evolve slower
		}

		for gen := 1; gen <= generations; gen++ {
			u.Step()

			if gen%2 == 0 || gen == generations {
				fmt.Printf("Generation %d:\n", gen)
				fmt.Println(view.Render(u))
				fmt.Println(view.RenderStats(gen, u.CountLiving(), 30.0))
				fmt.Println()
				time.Sleep(500 * time.Millisecond)
			}
		}

		fmt.Println("---")
		fmt.Println()
		time.Sleep(1 * time.Second)
	}

	// Comparison demo: same pattern with different rules
	fmt.Println("=== Rule Comparison Demo ===")
	fmt.Println()

	pattern := allPatterns["layer-oscillator"]

	// Test 1: No interaction
	fmt.Println("1. No Layer Interaction:")
	u1 := pattern.CreateUniverse(rules.ConwayRule{})
	u1.SetLayerInteraction(false)

	for i := 0; i < 3; i++ {
		u1.Step()
	}
	fmt.Println(view.Render(u1))
	fmt.Printf("After 3 steps: %d cells\n", u1.CountLiving())
	fmt.Println()

	// Test 2: Weighted neighbors (0.3)
	fmt.Println("2. Weighted Neighbors (weight=0.3):")
	u2 := pattern.CreateUniverse(rules.ConwayRule{})
	u2.SetLayerInteraction(true)
	u2.SetVerticalWeight(0.3)

	for i := 0; i < 3; i++ {
		u2.Step()
	}
	fmt.Println(view.Render(u2))
	fmt.Printf("After 3 steps: %d cells\n", u2.CountLiving())
	fmt.Println()

	// Test 3: Weighted neighbors (0.8)
	fmt.Println("3. Weighted Neighbors (weight=0.8):")
	u3 := pattern.CreateUniverse(rules.ConwayRule{})
	u3.SetLayerInteraction(true)
	u3.SetVerticalWeight(0.8)

	for i := 0; i < 3; i++ {
		u3.Step()
	}
	fmt.Println(view.Render(u3))
	fmt.Printf("After 3 steps: %d cells\n", u3.CountLiving())
	fmt.Println()

	// Test 4: Birth between layers
	fmt.Println("4. Birth Between Layers (single layer):")
	u4 := pattern.CreateUniverse(rules.ConwayRule{})
	u4.SetLayerInteraction(true)
	birthRule := rules.NewBirthBetweenLayersRule(rules.ConwayRule{}, false)
	u4.SetInteractionRule(birthRule)

	for i := 0; i < 3; i++ {
		u4.Step()
	}
	fmt.Println(view.Render(u4))
	fmt.Printf("After 3 steps: %d cells\n", u4.CountLiving())
	fmt.Println()

	fmt.Println("=== Demo Complete ===")
	fmt.Println()
	fmt.Println("Summary:")
	fmt.Printf("- Total patterns: %d\n", len(patternList))
	fmt.Println("- Patterns demonstrated:")
	for _, name := range patternList {
		fmt.Printf("  ✓ %s\n", allPatterns[name].Name)
	}
	fmt.Println()
	fmt.Println("Features showcased:")
	fmt.Println("  ✓ Multi-layer patterns")
	fmt.Println("  ✓ Layer interaction effects")
	fmt.Println("  ✓ Different interaction rules")
	fmt.Println("  ✓ Energy diffusion")
	fmt.Println("  ✓ Vertical movement")
	fmt.Println("  ✓ Layer-specific oscillation")
}
