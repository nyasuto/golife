package terminal

import (
	"golife/pkg/core"
	"golife/pkg/patterns"
	"golife/pkg/rules"
	"golife/pkg/universe"
	"testing"
)

func TestNewSlice3DView(t *testing.T) {
	rule := rules.Life3D_B6S567{}
	u := universe.New3D(20, 20, 20, rule)

	view := NewSlice3DView(u)

	if view == nil {
		t.Fatal("NewSlice3DView should not return nil")
	}

	if view.planeType != PlaneXY {
		t.Errorf("Default plane should be XY, got %d", view.planeType)
	}

	if view.slicePosition != 10 {
		t.Errorf("Default slice position should be 10 (middle), got %d", view.slicePosition)
	}

	if view.showMulti {
		t.Error("showMulti should default to false")
	}
}

func TestSlice3DView_SetPlaneType(t *testing.T) {
	rule := rules.Life3D_B6S567{}
	u := universe.New3D(30, 20, 10, rule)
	view := NewSlice3DView(u)

	tests := []struct {
		plane       PlaneType
		expectedPos int
		name        string
	}{
		{PlaneXY, 5, "XY plane (Z middle)"},
		{PlaneXZ, 10, "XZ plane (Y middle)"},
		{PlaneYZ, 15, "YZ plane (X middle)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			view.SetPlaneType(tt.plane)
			if view.planeType != tt.plane {
				t.Errorf("Plane type not set correctly")
			}
			if view.slicePosition != tt.expectedPos {
				t.Errorf("Expected position %d, got %d", tt.expectedPos, view.slicePosition)
			}
		})
	}
}

func TestSlice3DView_NextPrevSlice(t *testing.T) {
	rule := rules.Life3D_B6S567{}
	u := universe.New3D(10, 10, 10, rule)
	view := NewSlice3DView(u)

	// Start at middle (5)
	if view.slicePosition != 5 {
		t.Fatalf("Expected initial position 5, got %d", view.slicePosition)
	}

	// Move forward
	view.NextSlice()
	if view.slicePosition != 6 {
		t.Errorf("NextSlice: expected 6, got %d", view.slicePosition)
	}

	// Move backward
	view.PrevSlice()
	if view.slicePosition != 5 {
		t.Errorf("PrevSlice: expected 5, got %d", view.slicePosition)
	}

	// Test boundaries - move to start
	for i := 0; i < 10; i++ {
		view.PrevSlice()
	}
	if view.slicePosition != 0 {
		t.Errorf("Should stop at 0, got %d", view.slicePosition)
	}

	// Test boundaries - move to end
	for i := 0; i < 20; i++ {
		view.NextSlice()
	}
	if view.slicePosition != 9 {
		t.Errorf("Should stop at 9 (max), got %d", view.slicePosition)
	}
}

func TestSlice3DView_ToggleMultiView(t *testing.T) {
	rule := rules.Life3D_B6S567{}
	u := universe.New3D(10, 10, 10, rule)
	view := NewSlice3DView(u)

	if view.showMulti {
		t.Error("Should start with showMulti=false")
	}

	view.ToggleMultiView()
	if !view.showMulti {
		t.Error("ToggleMultiView should set to true")
	}

	view.ToggleMultiView()
	if view.showMulti {
		t.Error("ToggleMultiView should set back to false")
	}
}

func TestSlice3DView_GetMaxSlicePosition(t *testing.T) {
	rule := rules.Life3D_B6S567{}
	u := universe.New3D(30, 20, 10, rule)
	view := NewSlice3DView(u)

	tests := []struct {
		plane    PlaneType
		expected int
		name     string
	}{
		{PlaneXY, 10, "XY plane (Z depth)"},
		{PlaneXZ, 20, "XZ plane (Y height)"},
		{PlaneYZ, 30, "YZ plane (X width)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			view.SetPlaneType(tt.plane)
			max := view.getMaxSlicePosition()
			if max != tt.expected {
				t.Errorf("Expected max %d, got %d", tt.expected, max)
			}
		})
	}
}

func TestSlice3DView_GetSliceDimensions(t *testing.T) {
	rule := rules.Life3D_B6S567{}
	u := universe.New3D(30, 20, 10, rule)
	view := NewSlice3DView(u)

	tests := []struct {
		plane     PlaneType
		expectedW int
		expectedH int
		name      string
	}{
		{PlaneXY, 30, 20, "XY plane (width x height)"},
		{PlaneXZ, 30, 10, "XZ plane (width x depth)"},
		{PlaneYZ, 20, 10, "YZ plane (height x depth)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			view.SetPlaneType(tt.plane)
			w, h := view.getSliceDimensions()
			if w != tt.expectedW || h != tt.expectedH {
				t.Errorf("Expected (%d,%d), got (%d,%d)", tt.expectedW, tt.expectedH, w, h)
			}
		})
	}
}

func TestSlice3DView_GetCellAt(t *testing.T) {
	rule := rules.Life3D_B6S567{}
	u := universe.New3D(10, 10, 10, rule)
	view := NewSlice3DView(u)

	// Set a specific cell
	u.Set(core.NewCoord3D(3, 4, 5), core.Alive)

	// Test XY plane at Z=5
	view.SetPlaneType(PlaneXY)
	view.slicePosition = 5
	state := view.getCellAt(3, 4)
	if state != core.Alive {
		t.Error("XY plane: Cell at (3,4) with Z=5 should be alive")
	}

	// Test XZ plane at Y=4
	view.SetPlaneType(PlaneXZ)
	view.slicePosition = 4
	state = view.getCellAt(3, 5)
	if state != core.Alive {
		t.Error("XZ plane: Cell at (3,5) with Y=4 should be alive")
	}

	// Test YZ plane at X=3
	view.SetPlaneType(PlaneYZ)
	view.slicePosition = 3
	state = view.getCellAt(4, 5)
	if state != core.Alive {
		t.Error("YZ plane: Cell at (4,5) with X=3 should be alive")
	}
}

func TestSlice3DView_GetCellAt_DeadCells(t *testing.T) {
	rule := rules.Life3D_B6S567{}
	u := universe.New3D(10, 10, 10, rule)
	view := NewSlice3DView(u)

	// All cells should be dead initially
	view.SetPlaneType(PlaneXY)
	view.slicePosition = 5

	for j := 0; j < 10; j++ {
		for i := 0; i < 10; i++ {
			state := view.getCellAt(i, j)
			if state != core.Dead {
				t.Errorf("Cell at (%d,%d) should be dead", i, j)
			}
		}
	}
}

func TestSlice3DView_GetPlaneInfo(t *testing.T) {
	rule := rules.Life3D_B6S567{}
	u := universe.New3D(20, 20, 20, rule)
	view := NewSlice3DView(u)

	tests := []struct {
		plane    PlaneType
		position int
		expected string
	}{
		{PlaneXY, 5, "XY Plane (Z=5/19)"},
		{PlaneXZ, 10, "XZ Plane (Y=10/19)"},
		{PlaneYZ, 15, "YZ Plane (X=15/19)"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			view.SetPlaneType(tt.plane)
			view.slicePosition = tt.position
			info := view.getPlaneInfo()
			if info != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, info)
			}
		})
	}
}

func TestSlice3DView_GetSliceLabel(t *testing.T) {
	rule := rules.Life3D_B6S567{}
	u := universe.New3D(10, 10, 10, rule)
	view := NewSlice3DView(u)

	tests := []struct {
		plane    PlaneType
		expected string
	}{
		{PlaneXY, "Z=5"},
		{PlaneXZ, "Y=5"},
		{PlaneYZ, "X=5"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			view.SetPlaneType(tt.plane)
			label := view.getSliceLabel(5)
			if label != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, label)
			}
		})
	}
}

func TestSlice3DView_WithPattern(t *testing.T) {
	rule := rules.Life3D_B6S567{}
	u := universe.New3D(30, 30, 30, rule)
	view := NewSlice3DView(u)

	// Load Bays's Glider (10 cells)
	glider := patterns.Glider3D()
	glider.LoadIntoUniverse3D(u, 10, 10, 10)

	// Verify total cell count
	totalCells := u.CountLiving()
	if totalCells != 10 {
		t.Errorf("Expected 10 cells in Bays's Glider, got %d", totalCells)
	}

	// Verify pattern is visible in appropriate slices
	view.SetPlaneType(PlaneXY)
	view.slicePosition = 10 // Z=10 layer

	// Bays's Glider has cells at z=0 and z=1 layers
	// At offset (10,10,10), cells at z=10 should include:
	// (11,10,10), (12,10,10), (10,11,10), (12,11,10), (11,12,10)
	state1 := view.getCellAt(11, 10) // (1,0,0) + offset
	state2 := view.getCellAt(12, 10) // (2,0,0) + offset

	if state1 != core.Alive || state2 != core.Alive {
		t.Error("Bays's Glider pattern should be visible at Z=10")
	}

	// At Z=11: should have the second layer of cells
	view.slicePosition = 11
	state3 := view.getCellAt(11, 10)
	if state3 != core.Alive {
		t.Error("Bays's Glider should have cells at Z=11")
	}

	// Count living cells in Z=10 slice (should be 5 cells in layer 0)
	view.slicePosition = 10
	count := 0
	for j := 0; j < 30; j++ {
		for i := 0; i < 30; i++ {
			if view.getCellAt(i, j) != core.Dead {
				count++
			}
		}
	}
	if count != 5 {
		t.Errorf("Expected 5 living cells in Z=10 slice, got %d", count)
	}
}

func TestSlice3DView_AllPlanesShowSamePattern(t *testing.T) {
	rule := rules.Life3D_B6S567{}
	u := universe.New3D(20, 20, 20, rule)

	// Create a simple 2x2x2 block at (5,6,7)
	block := patterns.Block3D()
	block.LoadIntoUniverse3D(u, 5, 6, 7)

	totalCells := u.CountLiving()
	if totalCells != 8 {
		t.Fatalf("Expected 8 living cells in block, got %d", totalCells)
	}

	view := NewSlice3DView(u)

	// Count cells visible in each plane
	countPlane := func(plane PlaneType, pos int) int {
		view.SetPlaneType(plane)
		view.slicePosition = pos
		w, h := view.getSliceDimensions()
		count := 0
		for j := 0; j < h; j++ {
			for i := 0; i < w; i++ {
				if view.getCellAt(i, j) != core.Dead {
					count++
				}
			}
		}
		return count
	}

	// XY plane at Z=7 should show 4 cells (2x2 square)
	xyCount := countPlane(PlaneXY, 7)
	if xyCount != 4 {
		t.Errorf("XY plane Z=7: expected 4 cells, got %d", xyCount)
	}

	// XZ plane at Y=6 should show 4 cells
	xzCount := countPlane(PlaneXZ, 6)
	if xzCount != 4 {
		t.Errorf("XZ plane Y=6: expected 4 cells, got %d", xzCount)
	}

	// YZ plane at X=5 should show 4 cells
	yzCount := countPlane(PlaneYZ, 5)
	if yzCount != 4 {
		t.Errorf("YZ plane X=5: expected 4 cells, got %d", yzCount)
	}
}

func BenchmarkSlice3DView_GetCellAt(b *testing.B) {
	rule := rules.Life3D_B6S567{}
	u := universe.New3D(100, 100, 100, rule)
	view := NewSlice3DView(u)

	view.SetPlaneType(PlaneXY)
	view.slicePosition = 50

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		view.getCellAt(i%100, (i/100)%100)
	}
}

func BenchmarkSlice3DView_GetSliceDimensions(b *testing.B) {
	rule := rules.Life3D_B6S567{}
	u := universe.New3D(100, 100, 100, rule)
	view := NewSlice3DView(u)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		view.getSliceDimensions()
	}
}
