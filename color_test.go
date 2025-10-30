package main

import (
	"testing"

	termbox "github.com/nsf/termbox-go"
)

// TestGetColorByAge tests the color assignment based on cell age
func TestGetColorByAge(t *testing.T) {
	testCases := []struct {
		name          string
		age           int
		expectedColor termbox.Attribute
	}{
		{
			name:          "Dead cell (age 0)",
			age:           0,
			expectedColor: termbox.ColorDefault,
		},
		{
			name:          "Newborn cell (age 1)",
			age:           1,
			expectedColor: termbox.ColorGreen,
		},
		{
			name:          "Young cell (age 2)",
			age:           2,
			expectedColor: termbox.ColorYellow,
		},
		{
			name:          "Young cell (age 3)",
			age:           3,
			expectedColor: termbox.ColorYellow,
		},
		{
			name:          "Old cell (age 4)",
			age:           4,
			expectedColor: termbox.ColorRed,
		},
		{
			name:          "Old cell (age 10)",
			age:           10,
			expectedColor: termbox.ColorRed,
		},
		{
			name:          "Very old cell (age 11)",
			age:           11,
			expectedColor: termbox.ColorBlue,
		},
		{
			name:          "Very old cell (age 100)",
			age:           100,
			expectedColor: termbox.ColorBlue,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := getColorByAge(tc.age)
			if result != tc.expectedColor {
				t.Errorf("Expected color %v for age %d, got %v", tc.expectedColor, tc.age, result)
			}
		})
	}
}

// TestStepWithAge tests the age tracking functionality
func TestStepWithAge(t *testing.T) {
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
		name           string
		grid           [][]int
		ageMap         [][]int
		expectedGrid   [][]int
		expectedAgeMap [][]int
		checkPosition  [2]int
		expectedAge    int
		description    string
	}{
		{
			name: "Newborn cell gets age 1",
			grid: [][]int{
				{0, 0, 0, 0, 0},
				{0, 1, 1, 0, 0},
				{0, 1, 0, 0, 0},
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
			},
			ageMap: [][]int{
				{0, 0, 0, 0, 0},
				{0, 5, 5, 0, 0},
				{0, 5, 0, 0, 0},
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
			},
			expectedGrid: [][]int{
				{0, 0, 0, 0, 0},
				{0, 1, 1, 0, 0},
				{0, 1, 1, 0, 0},
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
			},
			checkPosition: [2]int{2, 2}, // y=2, x=2
			expectedAge:   1,
			description:   "Newly born cell at (2,2) should have age 1",
		},
		{
			name: "Surviving cell increments age",
			grid: [][]int{
				{0, 0, 0, 0, 0},
				{0, 1, 1, 0, 0},
				{0, 1, 1, 0, 0},
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
			},
			ageMap: [][]int{
				{0, 0, 0, 0, 0},
				{0, 5, 5, 0, 0},
				{0, 5, 5, 0, 0},
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
			},
			expectedGrid: [][]int{
				{0, 0, 0, 0, 0},
				{0, 1, 1, 0, 0},
				{0, 1, 1, 0, 0},
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
			},
			checkPosition: [2]int{1, 1}, // y=1, x=1
			expectedAge:   6,
			description:   "Surviving cell at (1,1) should increment from age 5 to 6",
		},
		{
			name: "Dead cell has age 0",
			grid: [][]int{
				{0, 0, 0, 0, 0},
				{0, 1, 0, 0, 0},
				{0, 1, 0, 0, 0},
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
			},
			ageMap: [][]int{
				{0, 0, 0, 0, 0},
				{0, 10, 0, 0, 0},
				{0, 10, 0, 0, 0},
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
			},
			expectedGrid: [][]int{
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0},
			},
			checkPosition: [2]int{1, 1}, // y=1, x=1
			expectedAge:   0,
			description:   "Dead cell should have age 0",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultGrid, resultAgeMap := stepWithAge(tc.grid, tc.ageMap)

			// Verify the grid state
			for y := 0; y < DY; y++ {
				for x := 0; x < DX; x++ {
					if resultGrid[y][x] != tc.expectedGrid[y][x] {
						t.Errorf("Grid mismatch at [%d][%d]: expected %d, got %d",
							y, x, tc.expectedGrid[y][x], resultGrid[y][x])
					}
				}
			}

			// Verify the specific age
			y, x := tc.checkPosition[0], tc.checkPosition[1]
			if resultAgeMap[y][x] != tc.expectedAge {
				t.Errorf("%s: expected age %d at [%d][%d], got %d",
					tc.description, tc.expectedAge, y, x, resultAgeMap[y][x])
			}
		})
	}
}

// TestStepWithAgeProgression tests age progression over multiple generations
func TestStepWithAgeProgression(t *testing.T) {
	// Save original DX and DY
	originalDX := DX
	originalDY := DY
	DX = 5
	DY = 5
	defer func() {
		DX = originalDX
		DY = originalDY
	}()

	// Start with a stable block pattern (2x2)
	grid := [][]int{
		{0, 0, 0, 0, 0},
		{0, 1, 1, 0, 0},
		{0, 1, 1, 0, 0},
		{0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0},
	}

	ageMap := [][]int{
		{0, 0, 0, 0, 0},
		{0, 1, 1, 0, 0},
		{0, 1, 1, 0, 0},
		{0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0},
	}

	// Simulate 10 generations
	for gen := 1; gen <= 10; gen++ {
		grid, ageMap = stepWithAge(grid, ageMap)

		// Block pattern should remain stable
		// All cells in the block should have age = gen + 1
		expectedAge := gen + 1
		for y := 1; y <= 2; y++ {
			for x := 1; x <= 2; x++ {
				if grid[y][x] != 1 {
					t.Errorf("Generation %d: Block cell at [%d][%d] should be alive", gen, y, x)
				}
				if ageMap[y][x] != expectedAge {
					t.Errorf("Generation %d: Expected age %d at [%d][%d], got %d",
						gen, expectedAge, y, x, ageMap[y][x])
				}
			}
		}
	}
}
