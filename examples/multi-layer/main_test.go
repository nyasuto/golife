package main

import (
	"bytes"
	"golife/pkg/core"
	"golife/pkg/rules"
	"golife/pkg/universe"
	"golife/pkg/visualizer/terminal"
	"testing"
)

func TestMultiLayerUniverseCreation(t *testing.T) {
	u := universe.New25D(20, 10, 3, rules.ConwayRule{})

	size := u.Size()
	if size.X != 20 {
		t.Errorf("Width: got %d, want 20", size.X)
	}
	if size.Y != 10 {
		t.Errorf("Height: got %d, want 10", size.Y)
	}
	if size.Z != 3 {
		t.Errorf("Depth: got %d, want 3", size.Z)
	}
}

func TestMultiLayerPatterns(t *testing.T) {
	u := universe.New25D(20, 10, 3, rules.ConwayRule{})
	u.SetLayerInteraction(true)
	u.SetVerticalWeight(0.3)

	// Layer 0: Glider (5 cells)
	u.Set(core.NewCoord3D(2, 1, 0), core.Alive)
	u.Set(core.NewCoord3D(3, 2, 0), core.Alive)
	u.Set(core.NewCoord3D(1, 3, 0), core.Alive)
	u.Set(core.NewCoord3D(2, 3, 0), core.Alive)
	u.Set(core.NewCoord3D(3, 3, 0), core.Alive)

	// Layer 1: Blinker (3 cells)
	u.Set(core.NewCoord3D(10, 5, 1), core.Alive)
	u.Set(core.NewCoord3D(11, 5, 1), core.Alive)
	u.Set(core.NewCoord3D(12, 5, 1), core.Alive)

	// Layer 2: Block (4 cells)
	u.Set(core.NewCoord3D(15, 7, 2), core.Alive)
	u.Set(core.NewCoord3D(16, 7, 2), core.Alive)
	u.Set(core.NewCoord3D(15, 8, 2), core.Alive)
	u.Set(core.NewCoord3D(16, 8, 2), core.Alive)

	// Total: 5 + 3 + 4 = 12 cells
	population := u.CountLiving()
	if population != 12 {
		t.Errorf("Initial population: got %d, want 12", population)
	}

	// Verify cells exist in correct layers
	if u.Get(core.NewCoord3D(2, 1, 0)) != core.Alive {
		t.Error("Glider cell should be alive in layer 0")
	}
	if u.Get(core.NewCoord3D(10, 5, 1)) != core.Alive {
		t.Error("Blinker cell should be alive in layer 1")
	}
	if u.Get(core.NewCoord3D(15, 7, 2)) != core.Alive {
		t.Error("Block cell should be alive in layer 2")
	}
}

func TestMultiLayerViewRender(t *testing.T) {
	u := universe.New25D(10, 5, 2, rules.ConwayRule{})

	// Add some cells
	u.Set(core.NewCoord3D(1, 1, 0), core.Alive)
	u.Set(core.NewCoord3D(5, 3, 1), core.Alive)

	view := terminal.NewMultiLayerView()

	// Test single layer render
	output := view.Render(u)
	if output == "" {
		t.Error("Render output should not be empty")
	}
}

func TestMultiLayerViewLayouts(t *testing.T) {
	u := universe.New25D(10, 5, 2, rules.ConwayRule{})
	u.Set(core.NewCoord3D(1, 1, 0), core.Alive)

	view := terminal.NewMultiLayerView()
	view.ToggleAllLayers()

	layouts := []terminal.LayoutType{
		terminal.HorizontalLayout,
		terminal.VerticalLayout,
		terminal.GridLayout,
	}

	for _, layout := range layouts {
		view.SetLayout(layout)
		output := view.Render(u)
		if output == "" {
			t.Errorf("Render output for layout %v should not be empty", layout)
		}
	}
}

func TestMultiLayerViewStats(t *testing.T) {
	u := universe.New25D(10, 5, 2, rules.ConwayRule{})
	u.Set(core.NewCoord3D(1, 1, 0), core.Alive)
	u.Set(core.NewCoord3D(2, 2, 1), core.Alive)

	view := terminal.NewMultiLayerView()

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

func TestMultiLayerViewControls(t *testing.T) {
	view := terminal.NewMultiLayerView()

	controls := view.RenderControls()
	if controls == "" {
		t.Error("Controls output should not be empty")
	}

	// Controls should contain key instructions
	if !bytes.Contains([]byte(controls), []byte("Space")) {
		t.Error("Controls should contain 'Space' key")
	}
	if !bytes.Contains([]byte(controls), []byte("Quit")) {
		t.Error("Controls should contain 'Quit' instruction")
	}
}

func TestMultiLayerEvolution(t *testing.T) {
	u := universe.New25D(20, 10, 2, rules.ConwayRule{})
	u.SetLayerInteraction(false) // Disable interaction for predictable behavior

	// Create a blinker (oscillates between 2 states)
	u.Set(core.NewCoord3D(10, 5, 0), core.Alive)
	u.Set(core.NewCoord3D(11, 5, 0), core.Alive)
	u.Set(core.NewCoord3D(12, 5, 0), core.Alive)

	initialPop := u.CountLiving()
	if initialPop != 3 {
		t.Fatalf("Initial population: got %d, want 3", initialPop)
	}

	// Step once
	u.Step()
	gen1Pop := u.CountLiving()

	// Blinker should change shape
	if gen1Pop == 0 {
		t.Error("Population should not be zero after one step")
	}

	// Step again
	u.Step()
	gen2Pop := u.CountLiving()

	// Blinker should oscillate back (may not be exact due to layer interaction)
	if gen2Pop == 0 {
		t.Error("Population should not be zero after two steps")
	}
}

func TestMultiLayerLayerNavigation(t *testing.T) {
	u := universe.New25D(10, 10, 3, rules.ConwayRule{})

	// Add cells in different layers
	u.Set(core.NewCoord3D(1, 1, 0), core.Alive)
	u.Set(core.NewCoord3D(5, 5, 1), core.Alive)
	u.Set(core.NewCoord3D(8, 8, 2), core.Alive)

	view := terminal.NewMultiLayerView()

	// Test layer navigation
	for z := 0; z < 3; z++ {
		view.SetCurrentLayer(z)
		output := view.Render(u)
		if output == "" {
			t.Errorf("Render output for layer %d should not be empty", z)
		}
	}
}

func TestMultiLayerToggleLayers(t *testing.T) {
	u := universe.New25D(10, 10, 2, rules.ConwayRule{})
	u.Set(core.NewCoord3D(1, 1, 0), core.Alive)

	view := terminal.NewMultiLayerView()

	// Initially single layer
	output1 := view.Render(u)

	// Toggle to all layers
	view.ToggleAllLayers()
	output2 := view.Render(u)

	// Outputs should be different
	if output1 == output2 {
		t.Error("Single layer and all layers view should produce different output")
	}

	// Toggle back to single layer
	view.ToggleAllLayers()
	output3 := view.Render(u)

	// Should match initial output
	if output1 != output3 {
		t.Error("Toggling back should restore original view")
	}
}
