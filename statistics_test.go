package main

import (
	"testing"
	"time"
)

// TestCountLivingCells tests the countLivingCells function
func TestCountLivingCells(t *testing.T) {
	testCases := []struct {
		name     string
		grid     [][]int
		expected int
	}{
		{
			name: "Empty grid",
			grid: [][]int{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			},
			expected: 0,
		},
		{
			name: "Full grid",
			grid: [][]int{
				{1, 1, 1},
				{1, 1, 1},
				{1, 1, 1},
			},
			expected: 9,
		},
		{
			name: "Partial grid",
			grid: [][]int{
				{1, 0, 1},
				{0, 1, 0},
				{1, 0, 1},
			},
			expected: 5,
		},
		{
			name: "Glider pattern",
			grid: [][]int{
				{0, 1, 0},
				{0, 0, 1},
				{1, 1, 1},
			},
			expected: 5,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := countLivingCells(tc.grid)
			if result != tc.expected {
				t.Errorf("Expected %d living cells, got %d", tc.expected, result)
			}
		})
	}
}

// TestUpdateStatistics tests the updateStatistics function
func TestUpdateStatistics(t *testing.T) {
	testCases := []struct {
		name            string
		grid            [][]int
		prevLivingCells int
		expectedGen     int
		expectedLiving  int
		expectedBirths  int
		expectedDeaths  int
	}{
		{
			name: "No change",
			grid: [][]int{
				{1, 1, 0},
				{1, 1, 0},
				{0, 0, 0},
			},
			prevLivingCells: 4,
			expectedGen:     1,
			expectedLiving:  4,
			expectedBirths:  0,
			expectedDeaths:  0,
		},
		{
			name: "Births only",
			grid: [][]int{
				{1, 1, 1},
				{1, 1, 1},
				{0, 0, 0},
			},
			prevLivingCells: 4,
			expectedGen:     1,
			expectedLiving:  6,
			expectedBirths:  2,
			expectedDeaths:  0,
		},
		{
			name: "Deaths only",
			grid: [][]int{
				{1, 0, 0},
				{0, 1, 0},
				{0, 0, 0},
			},
			prevLivingCells: 5,
			expectedGen:     1,
			expectedLiving:  2,
			expectedBirths:  0,
			expectedDeaths:  3,
		},
		{
			name: "First generation from empty",
			grid: [][]int{
				{0, 1, 0},
				{0, 0, 1},
				{1, 1, 1},
			},
			prevLivingCells: 0,
			expectedGen:     1,
			expectedLiving:  5,
			expectedBirths:  5,
			expectedDeaths:  0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stats := Statistics{
				Generation:  0,
				LivingCells: tc.prevLivingCells,
				StartTime:   time.Now(),
			}

			updateStatistics(&stats, tc.grid, tc.prevLivingCells)

			if stats.Generation != tc.expectedGen {
				t.Errorf("Expected generation %d, got %d", tc.expectedGen, stats.Generation)
			}
			if stats.LivingCells != tc.expectedLiving {
				t.Errorf("Expected %d living cells, got %d", tc.expectedLiving, stats.LivingCells)
			}
			if stats.Births != tc.expectedBirths {
				t.Errorf("Expected %d births, got %d", tc.expectedBirths, stats.Births)
			}
			if stats.Deaths != tc.expectedDeaths {
				t.Errorf("Expected %d deaths, got %d", tc.expectedDeaths, stats.Deaths)
			}
		})
	}
}

// TestUpdateStatisticsFPS tests FPS calculation
func TestUpdateStatisticsFPS(t *testing.T) {
	stats := Statistics{
		Generation:  0,
		LivingCells: 0,
		StartTime:   time.Now(),
	}

	grid := [][]int{
		{1, 0, 0},
		{0, 1, 0},
		{0, 0, 1},
	}

	// First update - no FPS yet
	updateStatistics(&stats, grid, 0)
	if stats.FPS != 0 {
		t.Errorf("Expected FPS 0 on first update, got %.2f", stats.FPS)
	}

	// Wait a bit and update again
	time.Sleep(100 * time.Millisecond)
	updateStatistics(&stats, grid, stats.LivingCells)

	// FPS should be calculated now (approximately 10 FPS for 100ms)
	if stats.FPS <= 0 {
		t.Errorf("Expected positive FPS, got %.2f", stats.FPS)
	}
	if stats.FPS < 5 || stats.FPS > 15 {
		t.Errorf("Expected FPS around 10, got %.2f", stats.FPS)
	}
}

// TestUpdateStatisticsProgression tests multiple generations
func TestUpdateStatisticsProgression(t *testing.T) {
	stats := Statistics{
		Generation:  0,
		LivingCells: 0,
		StartTime:   time.Now(),
	}

	grids := [][][]int{
		{
			{0, 1, 0},
			{0, 0, 1},
			{1, 1, 1},
		},
		{
			{0, 0, 0},
			{1, 0, 1},
			{0, 1, 1},
		},
		{
			{0, 0, 0},
			{0, 0, 1},
			{1, 0, 1},
		},
	}

	expectedGenerations := []int{1, 2, 3}
	expectedLiving := []int{5, 4, 3}

	for i, grid := range grids {
		prevLiving := stats.LivingCells
		updateStatistics(&stats, grid, prevLiving)

		if stats.Generation != expectedGenerations[i] {
			t.Errorf("Step %d: Expected generation %d, got %d", i, expectedGenerations[i], stats.Generation)
		}
		if stats.LivingCells != expectedLiving[i] {
			t.Errorf("Step %d: Expected %d living cells, got %d", i, expectedLiving[i], stats.LivingCells)
		}
	}
}
