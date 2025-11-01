package terminal

import (
	"fmt"
	"golife/pkg/core"
	"golife/pkg/universe"

	"github.com/nsf/termbox-go"
)

// PlaneType represents the orientation of a 2D slice through 3D space
type PlaneType int

const (
	PlaneXY PlaneType = iota // Z-axis fixed (top view)
	PlaneXZ                  // Y-axis fixed (front view)
	PlaneYZ                  // X-axis fixed (side view)
)

// Slice3DView handles visualization of 3D universes by showing 2D slices
type Slice3DView struct {
	universe      *universe.Universe3D
	planeType     PlaneType
	slicePosition int  // Current position along the fixed axis
	showMulti     bool // Show multiple slices in a grid
	cellWidth     int
	cellHeight    int
}

// NewSlice3DView creates a new 3D slice visualizer
func NewSlice3DView(u *universe.Universe3D) *Slice3DView {
	return &Slice3DView{
		universe:      u,
		planeType:     PlaneXY,
		slicePosition: u.Size().Z / 2, // Start at middle slice
		showMulti:     false,
		cellWidth:     2,
		cellHeight:    1,
	}
}

// SetPlaneType changes the viewing plane
func (v *Slice3DView) SetPlaneType(plane PlaneType) {
	v.planeType = plane
	size := v.universe.Size()

	// Reset slice position to middle of new axis
	switch plane {
	case PlaneXY:
		v.slicePosition = size.Z / 2
	case PlaneXZ:
		v.slicePosition = size.Y / 2
	case PlaneYZ:
		v.slicePosition = size.X / 2
	}
}

// NextSlice moves to the next slice along the fixed axis
func (v *Slice3DView) NextSlice() {
	maxPos := v.getMaxSlicePosition()

	if v.slicePosition < maxPos-1 {
		v.slicePosition++
	}
}

// PrevSlice moves to the previous slice along the fixed axis
func (v *Slice3DView) PrevSlice() {
	if v.slicePosition > 0 {
		v.slicePosition--
	}
}

// ToggleMultiView toggles between single and multi-slice view
func (v *Slice3DView) ToggleMultiView() {
	v.showMulti = !v.showMulti
}

// getMaxSlicePosition returns the maximum valid slice position for current plane
func (v *Slice3DView) getMaxSlicePosition() int {
	size := v.universe.Size()
	switch v.planeType {
	case PlaneXY:
		return size.Z
	case PlaneXZ:
		return size.Y
	case PlaneYZ:
		return size.X
	default:
		return 0
	}
}

// getSliceDimensions returns (width, height) of the current slice plane
func (v *Slice3DView) getSliceDimensions() (int, int) {
	size := v.universe.Size()
	switch v.planeType {
	case PlaneXY:
		return size.X, size.Y
	case PlaneXZ:
		return size.X, size.Z
	case PlaneYZ:
		return size.Y, size.Z
	default:
		return 0, 0
	}
}

// getCellAt returns the cell state at (i, j) in the current slice
func (v *Slice3DView) getCellAt(i, j int) core.CellState {
	switch v.planeType {
	case PlaneXY:
		// i=x, j=y, fixed z
		return v.universe.Get(core.NewCoord3D(i, j, v.slicePosition))
	case PlaneXZ:
		// i=x, j=z, fixed y
		return v.universe.Get(core.NewCoord3D(i, v.slicePosition, j))
	case PlaneYZ:
		// i=y, j=z, fixed x
		return v.universe.Get(core.NewCoord3D(v.slicePosition, i, j))
	default:
		return core.Dead
	}
}

// Render draws the 3D slice view to the terminal
func (v *Slice3DView) Render(startX, startY, width, height int, stats string) {
	if v.showMulti {
		v.renderMultiView(startX, startY, width, height, stats)
	} else {
		v.renderSingleSlice(startX, startY, width, height, stats)
	}
}

// renderSingleSlice renders a single 2D slice
func (v *Slice3DView) renderSingleSlice(startX, startY, width, height int, stats string) {
	sliceWidth, sliceHeight := v.getSliceDimensions()

	// Calculate centering offset
	gridWidth := sliceWidth * v.cellWidth
	gridHeight := sliceHeight * v.cellHeight
	offsetX := startX + (width-gridWidth)/2
	offsetY := startY + (height-gridHeight-3)/2 // Reserve space for info

	// Render cells
	for j := 0; j < sliceHeight; j++ {
		for i := 0; i < sliceWidth; i++ {
			state := v.getCellAt(i, j)
			x := offsetX + i*v.cellWidth
			y := offsetY + j*v.cellHeight

			char := ' '
			fg := termbox.ColorDefault
			bg := termbox.ColorDefault

			if state != core.Dead {
				char = '█'
				fg = termbox.ColorGreen
			}

			for dx := 0; dx < v.cellWidth; dx++ {
				termbox.SetCell(x+dx, y, char, fg, bg)
			}
		}
	}

	// Render slice info
	v.renderSliceInfo(offsetX, offsetY+gridHeight+1, stats)
}

// renderMultiView renders multiple slices in a grid layout
func (v *Slice3DView) renderMultiView(startX, startY, width, height int, stats string) {
	maxSlices := v.getMaxSlicePosition()
	sliceWidth, sliceHeight := v.getSliceDimensions()

	// Determine grid layout (e.g., 2x2, 3x3)
	cols := 2
	rows := 2
	if maxSlices > 4 {
		cols = 3
		rows = 3
	}

	// Calculate cell size to fit multiple slices
	availableWidth := width / cols
	availableHeight := (height - 3) / rows
	cellW := availableWidth / sliceWidth
	cellH := availableHeight / sliceHeight
	if cellW < 1 {
		cellW = 1
	}
	if cellH < 1 {
		cellH = 1
	}

	// Render grid of slices
	sliceIndex := 0
	for row := 0; row < rows && sliceIndex < maxSlices; row++ {
		for col := 0; col < cols && sliceIndex < maxSlices; col++ {
			panelX := startX + col*availableWidth
			panelY := startY + row*availableHeight

			v.renderSlicePanel(panelX, panelY, availableWidth, availableHeight,
				sliceIndex, cellW, cellH)
			sliceIndex++
		}
	}

	// Render stats at bottom
	v.renderSliceInfo(startX, startY+height-2, stats)
}

// renderSlicePanel renders a single slice panel in multi-view
func (v *Slice3DView) renderSlicePanel(x, y, w, h, slicePos, cellW, cellH int) {
	sliceWidth, sliceHeight := v.getSliceDimensions()

	// Save original position
	origPos := v.slicePosition
	v.slicePosition = slicePos

	// Center the slice in the panel
	gridW := sliceWidth * cellW
	gridH := sliceHeight * cellH
	offsetX := x + (w-gridW)/2
	offsetY := y + (h-gridH-1)/2

	// Render cells
	for j := 0; j < sliceHeight; j++ {
		for i := 0; i < sliceWidth; i++ {
			state := v.getCellAt(i, j)
			cellX := offsetX + i*cellW
			cellY := offsetY + j*cellH

			char := '·'
			fg := termbox.ColorDarkGray
			bg := termbox.ColorDefault

			if state != core.Dead {
				char = '█'
				fg = termbox.ColorGreen
			}

			for dx := 0; dx < cellW; dx++ {
				if cellX+dx < x+w && cellY < y+h {
					termbox.SetCell(cellX+dx, cellY, char, fg, bg)
				}
			}
		}
	}

	// Label the slice
	label := v.getSliceLabel(slicePos)
	for i, ch := range label {
		if x+i < x+w {
			termbox.SetCell(x+i, offsetY+gridH, ch, termbox.ColorWhite, termbox.ColorDefault)
		}
	}

	// Restore original position
	v.slicePosition = origPos
}

// renderSliceInfo renders slice information and stats
func (v *Slice3DView) renderSliceInfo(x, y int, stats string) {
	// Plane and position info
	planeInfo := v.getPlaneInfo()
	for i, ch := range planeInfo {
		termbox.SetCell(x+i, y, ch, termbox.ColorYellow, termbox.ColorDefault)
	}

	// Stats
	for i, ch := range stats {
		termbox.SetCell(x+i, y+1, ch, termbox.ColorWhite, termbox.ColorDefault)
	}

	// Controls help
	help := "[PgUp/PgDn] Slice | [1/2/3] Plane | [m] Multi-view"
	for i, ch := range help {
		termbox.SetCell(x+i, y+2, ch, termbox.ColorCyan, termbox.ColorDefault)
	}
}

// getPlaneInfo returns a string describing the current plane and position
func (v *Slice3DView) getPlaneInfo() string {
	size := v.universe.Size()
	switch v.planeType {
	case PlaneXY:
		return fmt.Sprintf("XY Plane (Z=%d/%d)", v.slicePosition, size.Z-1)
	case PlaneXZ:
		return fmt.Sprintf("XZ Plane (Y=%d/%d)", v.slicePosition, size.Y-1)
	case PlaneYZ:
		return fmt.Sprintf("YZ Plane (X=%d/%d)", v.slicePosition, size.X-1)
	default:
		return "Unknown Plane"
	}
}

// getSliceLabel returns a short label for a slice at given position
func (v *Slice3DView) getSliceLabel(pos int) string {
	switch v.planeType {
	case PlaneXY:
		return fmt.Sprintf("Z=%d", pos)
	case PlaneXZ:
		return fmt.Sprintf("Y=%d", pos)
	case PlaneYZ:
		return fmt.Sprintf("X=%d", pos)
	default:
		return ""
	}
}
