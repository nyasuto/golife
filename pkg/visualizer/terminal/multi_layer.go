package terminal

import (
	"fmt"
	"golife/pkg/core"
	"golife/pkg/universe"
	"strings"
)

// LayoutType defines how multiple layers are displayed
type LayoutType int

const (
	// SingleLayer shows one layer at a time
	SingleLayer LayoutType = iota
	// HorizontalLayout displays layers side by side
	HorizontalLayout
	// VerticalLayout stacks layers vertically
	VerticalLayout
	// GridLayout arranges layers in a grid
	GridLayout
)

// MultiLayerView manages visualization of multiple 2D layers
type MultiLayerView struct {
	currentLayer  int
	showAllLayers bool
	layout        LayoutType
	cellWidth     int // Width of each cell in characters
	cellHeight    int // Height of each cell in lines
}

// NewMultiLayerView creates a new multi-layer visualizer
func NewMultiLayerView() *MultiLayerView {
	return &MultiLayerView{
		currentLayer:  0,
		showAllLayers: false,
		layout:        SingleLayer,
		cellWidth:     1,
		cellHeight:    1,
	}
}

// SetLayout changes the display layout
func (v *MultiLayerView) SetLayout(layout LayoutType) {
	v.layout = layout
}

// SetCurrentLayer sets which layer to display in single-layer mode
func (v *MultiLayerView) SetCurrentLayer(layer int) {
	v.currentLayer = layer
}

// NextLayer moves to the next layer
func (v *MultiLayerView) NextLayer(maxLayers int) {
	v.currentLayer = (v.currentLayer + 1) % maxLayers
}

// PrevLayer moves to the previous layer
func (v *MultiLayerView) PrevLayer(maxLayers int) {
	v.currentLayer--
	if v.currentLayer < 0 {
		v.currentLayer = maxLayers - 1
	}
}

// ToggleAllLayers switches between single and all-layers view
func (v *MultiLayerView) ToggleAllLayers() {
	v.showAllLayers = !v.showAllLayers
}

// GetCurrentLayer returns the current layer index
func (v *MultiLayerView) GetCurrentLayer() int {
	return v.currentLayer
}

// IsShowingAllLayers returns whether all layers are displayed
func (v *MultiLayerView) IsShowingAllLayers() bool {
	return v.showAllLayers
}

// Render renders the universe to a string
func (v *MultiLayerView) Render(u *universe.Universe25D) string {
	if v.showAllLayers {
		switch v.layout {
		case HorizontalLayout:
			return v.renderHorizontal(u)
		case VerticalLayout:
			return v.renderVertical(u)
		case GridLayout:
			return v.renderGrid(u)
		default:
			return v.renderGrid(u)
		}
	}
	return v.renderSingleLayer(u, v.currentLayer)
}

// renderSingleLayer renders a single layer
func (v *MultiLayerView) renderSingleLayer(u *universe.Universe25D, layerIndex int) string {
	size := u.Size()
	if layerIndex < 0 || layerIndex >= size.Z {
		return fmt.Sprintf("Invalid layer: %d (max: %d)", layerIndex, size.Z-1)
	}

	var sb strings.Builder

	// Header
	sb.WriteString(fmt.Sprintf("Layer %d (Z=%d)\n", layerIndex, layerIndex))
	sb.WriteString(v.renderBorder(size.X, "top"))

	// Cells
	for y := 0; y < size.Y; y++ {
		sb.WriteString("║")
		for x := 0; x < size.X; x++ {
			coord := core.NewCoord3D(x, y, layerIndex)
			state := u.Get(coord)
			sb.WriteString(v.cellChar(state))
		}
		sb.WriteString("║\n")
	}

	// Footer
	sb.WriteString(v.renderBorder(size.X, "bottom"))

	return sb.String()
}

// renderHorizontal renders all layers horizontally
func (v *MultiLayerView) renderHorizontal(u *universe.Universe25D) string {
	size := u.Size()
	var sb strings.Builder

	// Headers for each layer
	for z := 0; z < size.Z; z++ {
		if z > 0 {
			sb.WriteString("  ")
		}
		sb.WriteString(fmt.Sprintf("Layer %d (Z=%d)", z, z))
		// Pad to match layer width
		padding := size.X + 2 - len(fmt.Sprintf("Layer %d (Z=%d)", z, z))
		if padding > 0 {
			sb.WriteString(strings.Repeat(" ", padding))
		}
	}
	sb.WriteString("\n")

	// Top borders
	for z := 0; z < size.Z; z++ {
		if z > 0 {
			sb.WriteString("  ")
		}
		sb.WriteString(v.renderBorder(size.X, "top"))
	}
	sb.WriteString("\n")

	// Cells row by row
	for y := 0; y < size.Y; y++ {
		for z := 0; z < size.Z; z++ {
			if z > 0 {
				sb.WriteString("  ")
			}
			sb.WriteString("║")
			for x := 0; x < size.X; x++ {
				coord := core.NewCoord3D(x, y, z)
				state := u.Get(coord)
				sb.WriteString(v.cellChar(state))
			}
			sb.WriteString("║")
		}
		sb.WriteString("\n")
	}

	// Bottom borders
	for z := 0; z < size.Z; z++ {
		if z > 0 {
			sb.WriteString("  ")
		}
		sb.WriteString(v.renderBorder(size.X, "bottom"))
	}
	sb.WriteString("\n")

	return sb.String()
}

// renderVertical renders all layers vertically
func (v *MultiLayerView) renderVertical(u *universe.Universe25D) string {
	size := u.Size()
	var sb strings.Builder

	for z := 0; z < size.Z; z++ {
		if z > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(v.renderSingleLayer(u, z))
	}

	return sb.String()
}

// renderGrid renders layers in a grid layout
func (v *MultiLayerView) renderGrid(u *universe.Universe25D) string {
	size := u.Size()

	// Calculate grid dimensions (try to make it roughly square)
	cols := 3
	if size.Z <= 2 {
		cols = size.Z
	} else if size.Z <= 4 {
		cols = 2
	}
	rows := (size.Z + cols - 1) / cols

	var sb strings.Builder

	for row := 0; row < rows; row++ {
		// Headers
		for col := 0; col < cols; col++ {
			z := row*cols + col
			if z >= size.Z {
				break
			}
			if col > 0 {
				sb.WriteString("  ")
			}
			header := fmt.Sprintf("Layer %d (Z=%d)", z, z)
			sb.WriteString(header)
			// Pad to match layer width
			padding := size.X + 2 - len(header)
			if padding > 0 {
				sb.WriteString(strings.Repeat(" ", padding))
			}
		}
		sb.WriteString("\n")

		// Top borders
		for col := 0; col < cols; col++ {
			z := row*cols + col
			if z >= size.Z {
				break
			}
			if col > 0 {
				sb.WriteString("  ")
			}
			sb.WriteString(v.renderBorder(size.X, "top"))
		}
		sb.WriteString("\n")

		// Cells
		for y := 0; y < size.Y; y++ {
			for col := 0; col < cols; col++ {
				z := row*cols + col
				if z >= size.Z {
					break
				}
				if col > 0 {
					sb.WriteString("  ")
				}
				sb.WriteString("║")
				for x := 0; x < size.X; x++ {
					coord := core.NewCoord3D(x, y, z)
					state := u.Get(coord)
					sb.WriteString(v.cellChar(state))
				}
				sb.WriteString("║")
			}
			sb.WriteString("\n")
		}

		// Bottom borders
		for col := 0; col < cols; col++ {
			z := row*cols + col
			if z >= size.Z {
				break
			}
			if col > 0 {
				sb.WriteString("  ")
			}
			sb.WriteString(v.renderBorder(size.X, "bottom"))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// renderBorder renders a border line
func (v *MultiLayerView) renderBorder(width int, position string) string {
	if position == "top" {
		return "╔" + strings.Repeat("═", width) + "╗"
	}
	return "╚" + strings.Repeat("═", width) + "╝"
}

// cellChar returns the character representation of a cell state
func (v *MultiLayerView) cellChar(state core.CellState) string {
	if state == core.Dead {
		return " "
	}
	return "●"
}

// RenderStats renders statistics information
func (v *MultiLayerView) RenderStats(generation int, living int, fps float64) string {
	return fmt.Sprintf("Gen: %d  Living: %d  FPS: %.1f", generation, living, fps)
}

// RenderControls renders control hints
func (v *MultiLayerView) RenderControls() string {
	if v.showAllLayers {
		return "[Space] Pause  [a] Toggle View  [1-3] Layout  [r] Restart  [q] Quit"
	}
	return "[Space] Pause  [↑↓] Layer  [a] All Layers  [r] Restart  [q] Quit"
}
