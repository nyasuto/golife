package main

import (
	"golife/pkg/core"
	"golife/pkg/patterns"
	"golife/pkg/rules"
	"golife/pkg/universe"
	"testing"
)

func TestCellData(t *testing.T) {
	tests := []struct {
		name    string
		cell    CellData
		x, y, z int
	}{
		{
			name: "origin cell",
			cell: CellData{X: 0, Y: 0, Z: 0},
			x:    0, y: 0, z: 0,
		},
		{
			name: "positive coordinates",
			cell: CellData{X: 10, Y: 20, Z: 30},
			x:    10, y: 20, z: 30,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.cell.X != tt.x || tt.cell.Y != tt.y || tt.cell.Z != tt.z {
				t.Errorf("CellData mismatch: got (%d, %d, %d), want (%d, %d, %d)",
					tt.cell.X, tt.cell.Y, tt.cell.Z, tt.x, tt.y, tt.z)
			}
		})
	}
}

func TestUniverseState(t *testing.T) {
	cells := []CellData{
		{X: 0, Y: 0, Z: 0},
		{X: 1, Y: 1, Z: 1},
	}

	state := UniverseState{
		Cells:      cells,
		Generation: 10,
		Population: 2,
		Width:      32,
		Height:     32,
		Depth:      32,
	}

	if len(state.Cells) != 2 {
		t.Errorf("Cells length: got %d, want 2", len(state.Cells))
	}
	if state.Generation != 10 {
		t.Errorf("Generation: got %d, want 10", state.Generation)
	}
	if state.Population != 2 {
		t.Errorf("Population: got %d, want 2", state.Population)
	}
	if state.Width != 32 || state.Height != 32 || state.Depth != 32 {
		t.Errorf("Dimensions: got (%d, %d, %d), want (32, 32, 32)",
			state.Width, state.Height, state.Depth)
	}
}

func TestExtractUniverseState_EmptyUniverse(t *testing.T) {
	rule := rules.Life3D_B6S567{}
	u := universe.New3D(8, 8, 8, rule)

	state := extractUniverseState(u, 0)

	if len(state.Cells) != 0 {
		t.Errorf("Empty universe should have 0 cells, got %d", len(state.Cells))
	}
	if state.Generation != 0 {
		t.Errorf("Generation: got %d, want 0", state.Generation)
	}
	if state.Population != 0 {
		t.Errorf("Population: got %d, want 0", state.Population)
	}
	if state.Width != 8 || state.Height != 8 || state.Depth != 8 {
		t.Errorf("Dimensions: got (%d, %d, %d), want (8, 8, 8)",
			state.Width, state.Height, state.Depth)
	}
}

func TestExtractUniverseState_SingleCell(t *testing.T) {
	rule := rules.Life3D_B6S567{}
	u := universe.New3D(8, 8, 8, rule)

	// Set a single cell alive
	coord := core.NewCoord3D(3, 4, 5)
	u.Set(coord, core.Alive)

	state := extractUniverseState(u, 1)

	if len(state.Cells) != 1 {
		t.Fatalf("Expected 1 cell, got %d", len(state.Cells))
	}
	if state.Cells[0].X != 3 || state.Cells[0].Y != 4 || state.Cells[0].Z != 5 {
		t.Errorf("Cell position: got (%d, %d, %d), want (3, 4, 5)",
			state.Cells[0].X, state.Cells[0].Y, state.Cells[0].Z)
	}
	if state.Generation != 1 {
		t.Errorf("Generation: got %d, want 1", state.Generation)
	}
	if state.Population != 1 {
		t.Errorf("Population: got %d, want 1", state.Population)
	}
}

func TestExtractUniverseState_MultipleCells(t *testing.T) {
	rule := rules.Life3D_B6S567{}
	u := universe.New3D(16, 16, 16, rule)

	// Set multiple cells alive
	u.Set(core.NewCoord3D(0, 0, 0), core.Alive)
	u.Set(core.NewCoord3D(1, 1, 1), core.Alive)
	u.Set(core.NewCoord3D(5, 7, 9), core.Alive)

	state := extractUniverseState(u, 5)

	if len(state.Cells) != 3 {
		t.Fatalf("Expected 3 cells, got %d", len(state.Cells))
	}
	if state.Generation != 5 {
		t.Errorf("Generation: got %d, want 5", state.Generation)
	}
	if state.Population != 3 {
		t.Errorf("Population: got %d, want 3", state.Population)
	}

	// Verify cell positions (cells are collected in z-y-x order)
	expectedCells := []CellData{
		{X: 0, Y: 0, Z: 0},
		{X: 1, Y: 1, Z: 1},
		{X: 5, Y: 7, Z: 9},
	}

	for i, expected := range expectedCells {
		if state.Cells[i] != expected {
			t.Errorf("Cell %d: got (%d, %d, %d), want (%d, %d, %d)",
				i, state.Cells[i].X, state.Cells[i].Y, state.Cells[i].Z,
				expected.X, expected.Y, expected.Z)
		}
	}
}

func TestExtractUniverseState_WithPattern(t *testing.T) {
	rule := rules.Life3D_B6S567{}
	size := 32
	u := universe.New3D(size, size, size, rule)

	// Load Bays's Glider pattern
	glider := patterns.BaysGlider()
	glider.LoadIntoUniverse3D(u, size/2-2, size/2-2, size/2-2)

	state := extractUniverseState(u, 0)

	// Bays's Glider has 7 cells
	expectedCells := len(glider.Cells)
	if len(state.Cells) != expectedCells {
		t.Errorf("Expected %d cells, got %d", expectedCells, len(state.Cells))
	}
	if state.Population != expectedCells {
		t.Errorf("Population: got %d, want %d", state.Population, expectedCells)
	}
	if state.Width != size || state.Height != size || state.Depth != size {
		t.Errorf("Dimensions: got (%d, %d, %d), want (%d, %d, %d)",
			state.Width, state.Height, state.Depth, size, size, size)
	}
}

func TestExtractUniverseState_AfterStep(t *testing.T) {
	rule := rules.Life3D_B6S567{}
	size := 32
	u := universe.New3D(size, size, size, rule)

	// Load Bays's Glider pattern
	glider := patterns.BaysGlider()
	glider.LoadIntoUniverse3D(u, size/2-2, size/2-2, size/2-2)

	// Initial state
	initialState := extractUniverseState(u, 0)
	initialPop := initialState.Population

	// Step the simulation
	u.StepParallel()

	// State after one step
	nextState := extractUniverseState(u, 1)

	// Check generation incremented
	if nextState.Generation != 1 {
		t.Errorf("Generation: got %d, want 1", nextState.Generation)
	}

	// Population should change after step
	if nextState.Population == initialPop {
		t.Logf("Warning: Population unchanged after step (both %d)", initialPop)
	}

	// Cells count should match population
	if len(nextState.Cells) != nextState.Population {
		t.Errorf("Cells length (%d) != Population (%d)",
			len(nextState.Cells), nextState.Population)
	}
}
