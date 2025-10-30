package main

import (
	"testing"
)

// TestAvailablePatterns tests that all patterns are properly defined
func TestAvailablePatterns(t *testing.T) {
	patterns := availablePatterns()

	expectedPatterns := []string{"glider", "blinker", "toad", "beacon", "pulsar", "glider-gun"}

	if len(patterns) != len(expectedPatterns) {
		t.Errorf("Expected %d patterns, got %d", len(expectedPatterns), len(patterns))
	}

	for _, name := range expectedPatterns {
		p, exists := patterns[name]
		if !exists {
			t.Errorf("Pattern '%s' not found in available patterns", name)
			continue
		}

		// Verify pattern has required fields
		if p.Name == "" {
			t.Errorf("Pattern '%s' has empty Name field", name)
		}
		if p.Description == "" {
			t.Errorf("Pattern '%s' has empty Description field", name)
		}
		if p.Width <= 0 {
			t.Errorf("Pattern '%s' has invalid Width: %d", name, p.Width)
		}
		if p.Height <= 0 {
			t.Errorf("Pattern '%s' has invalid Height: %d", name, p.Height)
		}
		if len(p.Cells) != p.Height {
			t.Errorf("Pattern '%s' cells height mismatch: expected %d, got %d", name, p.Height, len(p.Cells))
		}
		for y, row := range p.Cells {
			if len(row) != p.Width {
				t.Errorf("Pattern '%s' cells width mismatch at row %d: expected %d, got %d", name, y, p.Width, len(row))
			}
			// Verify all cells are 0 or 1
			for x, cell := range row {
				if cell != 0 && cell != 1 {
					t.Errorf("Pattern '%s' has invalid cell value at [%d][%d]: %d", name, y, x, cell)
				}
			}
		}
	}
}

// TestLoadPattern tests the pattern loading functionality
func TestLoadPattern(t *testing.T) {
	// Save original DX and DY
	originalDX := DX
	originalDY := DY
	DX = 50
	DY = 50
	defer func() {
		DX = originalDX
		DY = originalDY
	}()

	testCases := []struct {
		name        string
		patternName string
		expectError bool
	}{
		{
			name:        "Load glider pattern",
			patternName: "glider",
			expectError: false,
		},
		{
			name:        "Load blinker pattern",
			patternName: "blinker",
			expectError: false,
		},
		{
			name:        "Load pulsar pattern",
			patternName: "pulsar",
			expectError: false,
		},
		{
			name:        "Load non-existent pattern",
			patternName: "nonexistent",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := loadPattern(tc.patternName)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error for pattern '%s', but got none", tc.patternName)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error loading pattern '%s': %v", tc.patternName, err)
				return
			}

			// Verify grid dimensions
			if len(result) != DY {
				t.Errorf("Expected grid height %d, got %d", DY, len(result))
			}

			for y := 0; y < len(result); y++ {
				if len(result[y]) != DX {
					t.Errorf("Expected grid width %d at row %d, got %d", DX, y, len(result[y]))
				}
			}

			// Verify pattern is centered (at least some cells are alive)
			aliveCells := 0
			for y := 0; y < DY; y++ {
				for x := 0; x < DX; x++ {
					if result[y][x] == 1 {
						aliveCells++
					}
				}
			}

			if aliveCells == 0 {
				t.Errorf("Pattern '%s' loaded but no alive cells found", tc.patternName)
			}
		})
	}
}

// TestLoadPatternCentering tests that patterns are correctly centered
func TestLoadPatternCentering(t *testing.T) {
	// Save original DX and DY
	originalDX := DX
	originalDY := DY
	DX = 20
	DY = 20
	defer func() {
		DX = originalDX
		DY = originalDY
	}()

	// Load glider pattern (3x3)
	result, err := loadPattern("glider")
	if err != nil {
		t.Fatalf("Failed to load glider pattern: %v", err)
	}

	// Calculate expected center position
	patterns := availablePatterns()
	glider := patterns["glider"]
	expectedStartX := (DX - glider.Width) / 2
	expectedStartY := (DY - glider.Height) / 2

	// Verify pattern is placed at center
	for y := 0; y < glider.Height; y++ {
		for x := 0; x < glider.Width; x++ {
			gridY := expectedStartY + y
			gridX := expectedStartX + x
			if result[gridY][gridX] != glider.Cells[y][x] {
				t.Errorf("Pattern cell mismatch at [%d][%d]: expected %d, got %d",
					gridY, gridX, glider.Cells[y][x], result[gridY][gridX])
			}
		}
	}
}

// TestListPatterns tests the pattern listing function
func TestListPatterns(t *testing.T) {
	result := listPatterns()

	if result == "" {
		t.Error("listPatterns returned empty string")
	}

	// Verify it contains expected pattern names
	expectedPatterns := []string{"glider", "blinker", "pulsar", "beacon", "toad", "glider-gun"}
	for _, name := range expectedPatterns {
		if !contains(result, name) {
			t.Errorf("listPatterns output does not contain pattern '%s'", name)
		}
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsAt(s, substr))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// TestGliderPattern tests the glider pattern's behavior
func TestGliderPattern(t *testing.T) {
	// Save original DX and DY
	originalDX := DX
	originalDY := DY
	DX = 10
	DY = 10
	defer func() {
		DX = originalDX
		DY = originalDY
	}()

	// Create a grid with glider in top-left corner
	grid := make([][]int, DY)
	for y := 0; y < DY; y++ {
		grid[y] = make([]int, DX)
	}

	// Place glider at position (2, 2)
	grid[2][3] = 1
	grid[3][4] = 1
	grid[4][2] = 1
	grid[4][3] = 1
	grid[4][4] = 1

	// Run one step
	result := step(grid)

	// Verify glider has moved/transformed correctly
	// After one step, glider should change shape
	expectedAliveCells := 5
	aliveCells := 0
	for y := 0; y < DY; y++ {
		for x := 0; x < DX; x++ {
			if result[y][x] == 1 {
				aliveCells++
			}
		}
	}

	if aliveCells != expectedAliveCells {
		t.Errorf("Expected %d alive cells after glider step, got %d", expectedAliveCells, aliveCells)
	}
}
