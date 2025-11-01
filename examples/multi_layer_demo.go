package main

import (
	"fmt"
	"golife/pkg/core"
	"golife/pkg/rules"
	"golife/pkg/universe"
	"golife/pkg/visualizer/terminal"
	"time"
)

func main() {
	// Create a 20x10x3 universe (3 layers)
	u := universe.New25D(20, 10, 3, rules.ConwayRule{})

	// Enable layer interaction
	u.SetLayerInteraction(true)
	u.SetVerticalWeight(0.3)

	// Create different patterns in each layer
	// Layer 0: Glider
	u.Set(core.NewCoord3D(2, 1, 0), core.Alive)
	u.Set(core.NewCoord3D(3, 2, 0), core.Alive)
	u.Set(core.NewCoord3D(1, 3, 0), core.Alive)
	u.Set(core.NewCoord3D(2, 3, 0), core.Alive)
	u.Set(core.NewCoord3D(3, 3, 0), core.Alive)

	// Layer 1: Blinker
	u.Set(core.NewCoord3D(10, 5, 1), core.Alive)
	u.Set(core.NewCoord3D(11, 5, 1), core.Alive)
	u.Set(core.NewCoord3D(12, 5, 1), core.Alive)

	// Layer 2: Block
	u.Set(core.NewCoord3D(15, 7, 2), core.Alive)
	u.Set(core.NewCoord3D(16, 7, 2), core.Alive)
	u.Set(core.NewCoord3D(15, 8, 2), core.Alive)
	u.Set(core.NewCoord3D(16, 8, 2), core.Alive)

	view := terminal.NewMultiLayerView()

	fmt.Println("=== Multi-Layer Game of Life Demo ===")
	fmt.Println()

	// Demo 1: Single Layer View
	fmt.Println("1. Single Layer View (Layer 0):")
	fmt.Println(view.Render(u))
	fmt.Println(view.RenderStats(0, u.CountLiving(), 30.0))
	fmt.Println(view.RenderControls())
	fmt.Println()

	// Demo 2: Horizontal Layout
	view.SetLayout(terminal.HorizontalLayout)
	view.ToggleAllLayers()
	fmt.Println("2. Horizontal Layout (All Layers):")
	fmt.Println(view.Render(u))
	fmt.Println()

	// Demo 3: Vertical Layout
	view.SetLayout(terminal.VerticalLayout)
	fmt.Println("3. Vertical Layout (All Layers):")
	fmt.Println(view.Render(u))
	fmt.Println()

	// Demo 4: Grid Layout
	view.SetLayout(terminal.GridLayout)
	fmt.Println("4. Grid Layout (All Layers):")
	fmt.Println(view.Render(u))
	fmt.Println()

	// Demo 5: Evolution over time
	fmt.Println("5. Evolution (5 generations with Grid Layout):")
	for gen := 1; gen <= 5; gen++ {
		u.Step()
		fmt.Printf("\n--- Generation %d ---\n", gen)
		fmt.Println(view.Render(u))
		fmt.Println(view.RenderStats(gen, u.CountLiving(), 30.0))
		time.Sleep(500 * time.Millisecond)
	}

	// Demo 6: Layer Navigation
	view.ToggleAllLayers() // Back to single layer
	fmt.Println("\n6. Layer Navigation:")
	for z := 0; z < 3; z++ {
		view.SetCurrentLayer(z)
		fmt.Printf("\n--- Viewing Layer %d ---\n", z)
		fmt.Println(view.Render(u))
	}

	fmt.Println("\n=== Demo Complete ===")
	fmt.Println("\nFeatures demonstrated:")
	fmt.Println("✓ Single layer view")
	fmt.Println("✓ Horizontal layout")
	fmt.Println("✓ Vertical layout")
	fmt.Println("✓ Grid layout")
	fmt.Println("✓ Layer interaction")
	fmt.Println("✓ Layer navigation")
}
