package terminal

import (
	"golife/pkg/core"
	"golife/pkg/rules"
	"golife/pkg/universe"
	"strings"
	"testing"
)

func TestNewMultiLayerView(t *testing.T) {
	v := NewMultiLayerView()

	if v.currentLayer != 0 {
		t.Errorf("Expected currentLayer 0, got %d", v.currentLayer)
	}

	if v.showAllLayers {
		t.Error("showAllLayers should be false by default")
	}

	if v.layout != SingleLayer {
		t.Errorf("Expected SingleLayer layout, got %d", v.layout)
	}
}

func TestMultiLayerView_SetLayout(t *testing.T) {
	v := NewMultiLayerView()

	v.SetLayout(HorizontalLayout)
	if v.layout != HorizontalLayout {
		t.Error("Layout should be HorizontalLayout")
	}

	v.SetLayout(VerticalLayout)
	if v.layout != VerticalLayout {
		t.Error("Layout should be VerticalLayout")
	}

	v.SetLayout(GridLayout)
	if v.layout != GridLayout {
		t.Error("Layout should be GridLayout")
	}
}

func TestMultiLayerView_SetCurrentLayer(t *testing.T) {
	v := NewMultiLayerView()

	v.SetCurrentLayer(2)
	if v.GetCurrentLayer() != 2 {
		t.Errorf("Expected layer 2, got %d", v.GetCurrentLayer())
	}
}

func TestMultiLayerView_NextLayer(t *testing.T) {
	v := NewMultiLayerView()
	maxLayers := 3

	// Start at 0, next should be 1
	v.NextLayer(maxLayers)
	if v.GetCurrentLayer() != 1 {
		t.Errorf("Expected layer 1, got %d", v.GetCurrentLayer())
	}

	// Next should be 2
	v.NextLayer(maxLayers)
	if v.GetCurrentLayer() != 2 {
		t.Errorf("Expected layer 2, got %d", v.GetCurrentLayer())
	}

	// Next should wrap to 0
	v.NextLayer(maxLayers)
	if v.GetCurrentLayer() != 0 {
		t.Errorf("Expected layer 0 (wrap), got %d", v.GetCurrentLayer())
	}
}

func TestMultiLayerView_PrevLayer(t *testing.T) {
	v := NewMultiLayerView()
	maxLayers := 3

	// Start at 0, prev should wrap to 2
	v.PrevLayer(maxLayers)
	if v.GetCurrentLayer() != 2 {
		t.Errorf("Expected layer 2 (wrap), got %d", v.GetCurrentLayer())
	}

	// Prev should be 1
	v.PrevLayer(maxLayers)
	if v.GetCurrentLayer() != 1 {
		t.Errorf("Expected layer 1, got %d", v.GetCurrentLayer())
	}

	// Prev should be 0
	v.PrevLayer(maxLayers)
	if v.GetCurrentLayer() != 0 {
		t.Errorf("Expected layer 0, got %d", v.GetCurrentLayer())
	}
}

func TestMultiLayerView_ToggleAllLayers(t *testing.T) {
	v := NewMultiLayerView()

	if v.IsShowingAllLayers() {
		t.Error("Should not be showing all layers initially")
	}

	v.ToggleAllLayers()
	if !v.IsShowingAllLayers() {
		t.Error("Should be showing all layers after toggle")
	}

	v.ToggleAllLayers()
	if v.IsShowingAllLayers() {
		t.Error("Should not be showing all layers after second toggle")
	}
}

func TestMultiLayerView_RenderSingleLayer(t *testing.T) {
	u := universe.New25D(5, 3, 3, rules.ConwayRule{})

	// Set some cells in layer 1
	u.Set(core.NewCoord3D(1, 1, 1), core.Alive)
	u.Set(core.NewCoord3D(2, 1, 1), core.Alive)
	u.Set(core.NewCoord3D(3, 1, 1), core.Alive)

	v := NewMultiLayerView()
	v.SetCurrentLayer(1)

	output := v.Render(u)

	// Check output contains layer header
	if !strings.Contains(output, "Layer 1") {
		t.Error("Output should contain 'Layer 1'")
	}

	// Check output contains borders
	if !strings.Contains(output, "╔") || !strings.Contains(output, "╗") {
		t.Error("Output should contain top borders")
	}

	if !strings.Contains(output, "╚") || !strings.Contains(output, "╝") {
		t.Error("Output should contain bottom borders")
	}

	// Check output contains cells
	if !strings.Contains(output, "●") {
		t.Error("Output should contain alive cells")
	}

	// Count alive cells in output
	aliveCount := strings.Count(output, "●")
	if aliveCount != 3 {
		t.Errorf("Expected 3 alive cells in output, got %d", aliveCount)
	}
}

func TestMultiLayerView_RenderHorizontal(t *testing.T) {
	u := universe.New25D(5, 3, 2, rules.ConwayRule{})

	// Set some cells
	u.Set(core.NewCoord3D(1, 1, 0), core.Alive)
	u.Set(core.NewCoord3D(2, 1, 1), core.Alive)

	v := NewMultiLayerView()
	v.SetLayout(HorizontalLayout)
	v.ToggleAllLayers()

	output := v.Render(u)

	// Check both layers are present
	if !strings.Contains(output, "Layer 0") {
		t.Error("Output should contain 'Layer 0'")
	}
	if !strings.Contains(output, "Layer 1") {
		t.Error("Output should contain 'Layer 1'")
	}

	// Check cells are present
	aliveCount := strings.Count(output, "●")
	if aliveCount != 2 {
		t.Errorf("Expected 2 alive cells in output, got %d", aliveCount)
	}
}

func TestMultiLayerView_RenderVertical(t *testing.T) {
	u := universe.New25D(5, 3, 2, rules.ConwayRule{})

	// Set some cells
	u.Set(core.NewCoord3D(1, 1, 0), core.Alive)
	u.Set(core.NewCoord3D(2, 1, 1), core.Alive)

	v := NewMultiLayerView()
	v.SetLayout(VerticalLayout)
	v.ToggleAllLayers()

	output := v.Render(u)

	// Check both layers are present
	if !strings.Contains(output, "Layer 0") {
		t.Error("Output should contain 'Layer 0'")
	}
	if !strings.Contains(output, "Layer 1") {
		t.Error("Output should contain 'Layer 1'")
	}

	// Vertical layout should have layers stacked
	lines := strings.Split(output, "\n")
	if len(lines) < 10 {
		t.Error("Vertical layout should have multiple lines")
	}
}

func TestMultiLayerView_RenderGrid(t *testing.T) {
	u := universe.New25D(5, 3, 4, rules.ConwayRule{})

	// Set some cells
	for z := 0; z < 4; z++ {
		u.Set(core.NewCoord3D(z, 1, z), core.Alive)
	}

	v := NewMultiLayerView()
	v.SetLayout(GridLayout)
	v.ToggleAllLayers()

	output := v.Render(u)

	// Check all layers are present
	for z := 0; z < 4; z++ {
		layerStr := "Layer " + string(rune('0'+z))
		if !strings.Contains(output, layerStr) {
			t.Errorf("Output should contain '%s'", layerStr)
		}
	}

	// Check cells are present
	aliveCount := strings.Count(output, "●")
	if aliveCount != 4 {
		t.Errorf("Expected 4 alive cells in output, got %d", aliveCount)
	}
}

func TestMultiLayerView_RenderStats(t *testing.T) {
	v := NewMultiLayerView()

	stats := v.RenderStats(42, 487, 30.2)

	if !strings.Contains(stats, "Gen: 42") {
		t.Error("Stats should contain generation")
	}
	if !strings.Contains(stats, "Living: 487") {
		t.Error("Stats should contain living count")
	}
	if !strings.Contains(stats, "FPS: 30.2") {
		t.Error("Stats should contain FPS")
	}
}

func TestMultiLayerView_RenderControls(t *testing.T) {
	v := NewMultiLayerView()

	// Single layer controls
	controls := v.RenderControls()
	if !strings.Contains(controls, "↑↓") {
		t.Error("Single layer controls should mention layer navigation")
	}
	if !strings.Contains(controls, "All Layers") {
		t.Error("Single layer controls should mention 'All Layers'")
	}

	// All layers controls
	v.ToggleAllLayers()
	controls = v.RenderControls()
	if !strings.Contains(controls, "Toggle View") {
		t.Error("All layers controls should mention 'Toggle View'")
	}
	if !strings.Contains(controls, "Layout") {
		t.Error("All layers controls should mention 'Layout'")
	}
}

func TestMultiLayerView_InvalidLayer(t *testing.T) {
	u := universe.New25D(5, 3, 3, rules.ConwayRule{})

	v := NewMultiLayerView()
	v.SetCurrentLayer(10) // Invalid layer

	output := v.Render(u)

	if !strings.Contains(output, "Invalid layer") {
		t.Error("Should show error for invalid layer")
	}
}

func TestMultiLayerView_EmptyUniverse(t *testing.T) {
	u := universe.New25D(5, 3, 2, rules.ConwayRule{})

	v := NewMultiLayerView()
	output := v.Render(u)

	// Should still render with no alive cells
	if !strings.Contains(output, "Layer 0") {
		t.Error("Should render even with no alive cells")
	}

	// Should have no alive cells
	aliveCount := strings.Count(output, "●")
	if aliveCount != 0 {
		t.Errorf("Expected 0 alive cells, got %d", aliveCount)
	}
}

func TestMultiLayerView_CellChar(t *testing.T) {
	v := NewMultiLayerView()

	// Dead cell
	if v.cellChar(core.Dead) != " " {
		t.Error("Dead cell should be rendered as space")
	}

	// Alive cell
	if v.cellChar(core.Alive) != "●" {
		t.Error("Alive cell should be rendered as ●")
	}

	// High energy cell (should still render as alive)
	if v.cellChar(100) != "●" {
		t.Error("High energy cell should be rendered as ●")
	}
}
