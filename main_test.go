package main

import (
	"testing"
)

// TestRandomize tests the randomize function
func TestRandomize(t *testing.T) {
	// Save original DX and DY
	originalDX := DX
	originalDY := DY
	defer func() {
		DX = originalDX
		DY = originalDY
	}()

	// Test with default size
	result := randomize()

	// Check grid dimensions
	if len(result) != DY {
		t.Errorf("Expected height %d, got %d", DY, len(result))
	}

	for y := 0; y < len(result); y++ {
		if len(result[y]) != DX {
			t.Errorf("Expected width %d at row %d, got %d", DX, y, len(result[y]))
		}

		// Check all values are 0 or 1
		for x := 0; x < len(result[y]); x++ {
			if result[y][x] != 0 && result[y][x] != 1 {
				t.Errorf("Expected value 0 or 1 at [%d][%d], got %d", y, x, result[y][x])
			}
		}
	}

	// Test with different sizes
	testCases := []struct {
		width  int
		height int
	}{
		{10, 10},
		{5, 5},
		{1, 1},
	}

	for _, tc := range testCases {
		DX = tc.width
		DY = tc.height
		result = randomize()

		if len(result) != tc.height {
			t.Errorf("Expected height %d, got %d", tc.height, len(result))
		}

		for y := 0; y < len(result); y++ {
			if len(result[y]) != tc.width {
				t.Errorf("Expected width %d, got %d", tc.width, len(result[y]))
			}
		}
	}
}

// TestStepSurvivalRules tests the survival rules (2-3 neighbors)
func TestStepSurvivalRules(t *testing.T) {
	// Save original DX and DY
	originalDX := DX
	originalDY := DY
	DX = 5
	DY = 5
	defer func() {
		DX = originalDX
		DY = originalDY
	}()

	testCases := []struct {
		name     string
		input    [][]int
		expected [][]int
	}{
		{
			name: "cell survives with 2 neighbors",
			input: [][]int{
				{0, 0, 0, 0, 0},
				{0, 1, 1, 0, 0},
				{0, 1, 0, 0, 0},
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
			},
			expected: [][]int{
				{0, 0, 0, 0, 0},
				{0, 1, 1, 0, 0},
				{0, 1, 1, 0, 0},
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
			},
		},
		{
			name: "cell survives with 3 neighbors",
			input: [][]int{
				{0, 0, 0, 0, 0},
				{0, 1, 1, 0, 0},
				{0, 1, 1, 0, 0},
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
			},
			expected: [][]int{
				{0, 0, 0, 0, 0},
				{0, 1, 1, 0, 0},
				{0, 1, 1, 0, 0},
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
			},
		},
		{
			name: "cell dies with 1 neighbor (underpopulation)",
			input: [][]int{
				{0, 0, 0, 0, 0},
				{0, 1, 0, 0, 0},
				{0, 1, 0, 0, 0},
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
			},
			expected: [][]int{
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
			},
		},
		{
			name: "cell dies with 4+ neighbors (overpopulation)",
			input: [][]int{
				{0, 0, 0, 0, 0},
				{0, 1, 1, 1, 0},
				{0, 1, 1, 0, 0},
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
			},
			expected: [][]int{
				{0, 0, 1, 0, 0},
				{0, 1, 0, 1, 0},
				{0, 1, 0, 1, 0},
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := step(tc.input)
			for y := 0; y < DY; y++ {
				for x := 0; x < DX; x++ {
					if result[y][x] != tc.expected[y][x] {
						t.Errorf("At [%d][%d]: expected %d, got %d", y, x, tc.expected[y][x], result[y][x])
					}
				}
			}
		})
	}
}

// TestStepBirthRule tests the birth rule (exactly 3 neighbors)
func TestStepBirthRule(t *testing.T) {
	// Save original DX and DY
	originalDX := DX
	originalDY := DY
	DX = 5
	DY = 5
	defer func() {
		DX = originalDX
		DY = originalDY
	}()

	testCases := []struct {
		name     string
		input    [][]int
		expected [][]int
	}{
		{
			name: "cell is born with exactly 3 neighbors",
			input: [][]int{
				{0, 0, 0, 0, 0},
				{0, 1, 1, 0, 0},
				{0, 1, 0, 0, 0},
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
			},
			expected: [][]int{
				{0, 0, 0, 0, 0},
				{0, 1, 1, 0, 0},
				{0, 1, 1, 0, 0},
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
			},
		},
		{
			name: "cell is not born with 2 neighbors",
			input: [][]int{
				{0, 0, 0, 0, 0},
				{0, 1, 0, 0, 0},
				{0, 1, 0, 0, 0},
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
			},
			expected: [][]int{
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := step(tc.input)
			for y := 0; y < DY; y++ {
				for x := 0; x < DX; x++ {
					if result[y][x] != tc.expected[y][x] {
						t.Errorf("At [%d][%d]: expected %d, got %d", y, x, tc.expected[y][x], result[y][x])
					}
				}
			}
		})
	}
}

// TestStepEdgeCases tests edge cases (corners and edges)
func TestStepEdgeCases(t *testing.T) {
	// Save original DX and DY
	originalDX := DX
	originalDY := DY
	DX = 3
	DY = 3
	defer func() {
		DX = originalDX
		DY = originalDY
	}()

	testCases := []struct {
		name     string
		input    [][]int
		expected [][]int
	}{
		{
			name: "top-left corner with 2 neighbors survives",
			input: [][]int{
				{1, 1, 0},
				{1, 0, 0},
				{0, 0, 0},
			},
			expected: [][]int{
				{1, 1, 0},
				{1, 1, 0},
				{0, 0, 0},
			},
		},
		{
			name: "top-right corner with 2 neighbors survives",
			input: [][]int{
				{0, 1, 1},
				{0, 0, 1},
				{0, 0, 0},
			},
			expected: [][]int{
				{0, 1, 1},
				{0, 1, 1},
				{0, 0, 0},
			},
		},
		{
			name: "bottom-left corner with 2 neighbors survives",
			input: [][]int{
				{0, 0, 0},
				{1, 0, 0},
				{1, 1, 0},
			},
			expected: [][]int{
				{0, 0, 0},
				{1, 1, 0},
				{1, 1, 0},
			},
		},
		{
			name: "bottom-right corner with 2 neighbors survives",
			input: [][]int{
				{0, 0, 0},
				{0, 0, 1},
				{0, 1, 1},
			},
			expected: [][]int{
				{0, 0, 0},
				{0, 1, 1},
				{0, 1, 1},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := step(tc.input)
			for y := 0; y < DY; y++ {
				for x := 0; x < DX; x++ {
					if result[y][x] != tc.expected[y][x] {
						t.Errorf("At [%d][%d]: expected %d, got %d", y, x, tc.expected[y][x], result[y][x])
					}
				}
			}
		})
	}
}

// TestStepKnownPatterns tests well-known Game of Life patterns
func TestStepKnownPatterns(t *testing.T) {
	// Save original DX and DY
	originalDX := DX
	originalDY := DY
	DX = 5
	DY = 5
	defer func() {
		DX = originalDX
		DY = originalDY
	}()

	t.Run("Blinker (period 2 oscillator)", func(t *testing.T) {
		// Horizontal blinker
		input := [][]int{
			{0, 0, 0, 0, 0},
			{0, 0, 1, 0, 0},
			{0, 0, 1, 0, 0},
			{0, 0, 1, 0, 0},
			{0, 0, 0, 0, 0},
		}

		// Should become vertical
		expectedStep1 := [][]int{
			{0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0},
			{0, 1, 1, 1, 0},
			{0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0},
		}

		result := step(input)
		for y := 0; y < DY; y++ {
			for x := 0; x < DX; x++ {
				if result[y][x] != expectedStep1[y][x] {
					t.Errorf("Step 1: At [%d][%d]: expected %d, got %d", y, x, expectedStep1[y][x], result[y][x])
				}
			}
		}

		// Should oscillate back to horizontal
		result = step(result)
		for y := 0; y < DY; y++ {
			for x := 0; x < DX; x++ {
				if result[y][x] != input[y][x] {
					t.Errorf("Step 2: At [%d][%d]: expected %d, got %d", y, x, input[y][x], result[y][x])
				}
			}
		}
	})

	t.Run("Block (still life)", func(t *testing.T) {
		// 2x2 block should remain stable
		input := [][]int{
			{0, 0, 0, 0, 0},
			{0, 1, 1, 0, 0},
			{0, 1, 1, 0, 0},
			{0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0},
		}

		result := step(input)
		for y := 0; y < DY; y++ {
			for x := 0; x < DX; x++ {
				if result[y][x] != input[y][x] {
					t.Errorf("At [%d][%d]: expected %d, got %d", y, x, input[y][x], result[y][x])
				}
			}
		}
	})
}

// TestStepEmptyGrid tests that an empty grid stays empty
func TestStepEmptyGrid(t *testing.T) {
	// Save original DX and DY
	originalDX := DX
	originalDY := DY
	DX = 5
	DY = 5
	defer func() {
		DX = originalDX
		DY = originalDY
	}()

	input := [][]int{
		{0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0},
	}

	result := step(input)
	for y := 0; y < DY; y++ {
		for x := 0; x < DX; x++ {
			if result[y][x] != 0 {
				t.Errorf("At [%d][%d]: expected 0, got %d", y, x, result[y][x])
			}
		}
	}
}
