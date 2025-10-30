package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	termbox "github.com/nsf/termbox-go"
)

// Default values
const (
	defaultWidth       = 100
	defaultHeight      = 40
	defaultSpeed       = 200
	defaultGenerations = 300
)

// Configuration variables
var (
	width       int
	height      int
	speed       int
	generations int
	pattern     string
	showStats   bool
)

// Statistics holds simulation statistics
type Statistics struct {
	Generation    int
	LivingCells   int
	Births        int
	Deaths        int
	StartTime     time.Time
	LastFrameTime time.Time
	FPS           float64
}

// DX is width
var DX = defaultWidth

// DY is height
var DY = defaultHeight

// Pattern represents a predefined pattern
type Pattern struct {
	Name        string
	Description string
	Width       int
	Height      int
	Cells       [][]int
}

// availablePatterns returns a map of all available patterns
func availablePatterns() map[string]Pattern {
	return map[string]Pattern{
		"glider": {
			Name:        "Glider",
			Description: "A small pattern that moves diagonally",
			Width:       3,
			Height:      3,
			Cells: [][]int{
				{0, 1, 0},
				{0, 0, 1},
				{1, 1, 1},
			},
		},
		"blinker": {
			Name:        "Blinker",
			Description: "A period-2 oscillator",
			Width:       3,
			Height:      1,
			Cells: [][]int{
				{1, 1, 1},
			},
		},
		"toad": {
			Name:        "Toad",
			Description: "A period-2 oscillator",
			Width:       4,
			Height:      2,
			Cells: [][]int{
				{0, 1, 1, 1},
				{1, 1, 1, 0},
			},
		},
		"beacon": {
			Name:        "Beacon",
			Description: "A period-2 oscillator",
			Width:       4,
			Height:      4,
			Cells: [][]int{
				{1, 1, 0, 0},
				{1, 1, 0, 0},
				{0, 0, 1, 1},
				{0, 0, 1, 1},
			},
		},
		"pulsar": {
			Name:        "Pulsar",
			Description: "A period-3 oscillator",
			Width:       13,
			Height:      13,
			Cells: [][]int{
				{0, 0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				{1, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1},
				{1, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1},
				{1, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1},
				{0, 0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0, 0},
				{1, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1},
				{1, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1},
				{1, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1},
				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0, 0},
			},
		},
		"glider-gun": {
			Name:        "Gosper's Glider Gun",
			Description: "A pattern that continuously generates gliders",
			Width:       36,
			Height:      9,
			Cells: [][]int{
				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1},
				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1},
				{1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				{1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 1, 1, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
		},
	}
}

// loadPattern loads a predefined pattern into the center of the grid
func loadPattern(patternName string) ([][]int, error) {
	patterns := availablePatterns()
	p, exists := patterns[patternName]
	if !exists {
		return nil, fmt.Errorf("pattern '%s' not found", patternName)
	}

	// Create empty grid
	result := make([][]int, DY)
	for y := 0; y < DY; y++ {
		result[y] = make([]int, DX)
	}

	// Calculate center position
	startX := (DX - p.Width) / 2
	startY := (DY - p.Height) / 2

	// Place pattern in the center
	for y := 0; y < p.Height && startY+y < DY; y++ {
		for x := 0; x < p.Width && startX+x < DX; x++ {
			if startY+y >= 0 && startX+x >= 0 {
				result[startY+y][startX+x] = p.Cells[y][x]
			}
		}
	}

	return result, nil
}

// listPatterns returns a formatted string of all available patterns
func listPatterns() string {
	patterns := availablePatterns()
	result := "Available patterns:\n"
	for name, p := range patterns {
		result += fmt.Sprintf("  %s: %s\n", name, p.Description)
	}
	return result
}

// countLivingCells counts the number of living cells in the grid
func countLivingCells(data [][]int) int {
	count := 0
	for y := 0; y < len(data); y++ {
		for x := 0; x < len(data[y]); x++ {
			if data[y][x] == 1 {
				count++
			}
		}
	}
	return count
}

// updateStatistics updates the statistics for the current generation
func updateStatistics(stats *Statistics, data [][]int, prevLivingCells int) {
	stats.Generation++
	stats.LivingCells = countLivingCells(data)

	// Calculate births and deaths
	diff := stats.LivingCells - prevLivingCells
	if diff > 0 {
		stats.Births = diff
		stats.Deaths = 0
	} else if diff < 0 {
		stats.Births = 0
		stats.Deaths = -diff
	} else {
		stats.Births = 0
		stats.Deaths = 0
	}

	// Calculate FPS
	now := time.Now()
	if !stats.LastFrameTime.IsZero() {
		frameDuration := now.Sub(stats.LastFrameTime).Seconds()
		if frameDuration > 0 {
			stats.FPS = 1.0 / frameDuration
		}
	}
	stats.LastFrameTime = now
}

func randomize() [][]int {
	result := make([][]int, DY)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for y := 0; y < DY; y++ {
		result[y] = make([]int, DX)
		for x := 0; x < DX; x++ {

			result[y][x] = r.Intn(2)

		}
	}
	return result
}

// countNeighbors counts the number of alive neighbors around a cell
func countNeighbors(data [][]int, x, y int) int {
	count := 0

	// Check all 8 directions around the cell
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			// Skip the center cell itself
			if dx == 0 && dy == 0 {
				continue
			}

			// Calculate neighbor coordinates
			nx := x + dx
			ny := y + dy

			// Check boundaries
			if nx >= 0 && nx < DX && ny >= 0 && ny < DY {
				count += data[ny][nx]
			}
		}
	}

	return count
}

func step(data [][]int) [][]int {
	result := make([][]int, DY)

	for y := 0; y < DY; y++ {
		result[y] = make([]int, DX)
		for x := 0; x < DX; x++ {
			neighbors := countNeighbors(data, x, y)
			isAlive := data[y][x] == 1

			// Conway's Game of Life rules:
			// 1. Any live cell with 2 or 3 live neighbors survives
			// 2. Any dead cell with exactly 3 live neighbors becomes alive
			// 3. All other cells die or stay dead
			if isAlive && (neighbors == 2 || neighbors == 3) {
				result[y][x] = 1
			} else if !isAlive && neighbors == 3 {
				result[y][x] = 1
			} else {
				result[y][x] = 0
			}
		}
	}

	return result
}

func flush(data [][]int) error {
	for y := 0; y < DY; y++ {
		for x := 0; x < DX; x++ {
			var dot = ' '
			if data[y][x] == 1 {
				dot = '*'
			}
			termbox.SetCell(x, y, dot, termbox.ColorDefault, termbox.ColorDefault)

		}
	}

	return termbox.Flush()

}

// displayStatistics displays statistics on the screen
func displayStatistics(stats Statistics) {
	if !showStats {
		return
	}

	// Calculate position (top-right corner)
	startX := DX - 35
	startY := 0

	// Only display if there's enough space
	if startX < 0 {
		return
	}

	lines := []string{
		"╔═══════════════════════════════╗",
		fmt.Sprintf("║ Generation: %-17d ║", stats.Generation),
		fmt.Sprintf("║ Living cells: %-15d ║", stats.LivingCells),
		fmt.Sprintf("║ Births: +%-20d ║", stats.Births),
		fmt.Sprintf("║ Deaths: -%-20d ║", stats.Deaths),
		fmt.Sprintf("║ FPS: %-24.1f ║", stats.FPS),
		"╚═══════════════════════════════╝",
	}

	for i, line := range lines {
		for j, ch := range line {
			termbox.SetCell(startX+j, startY+i, ch, termbox.ColorCyan, termbox.ColorDefault)
		}
	}
}

func init() {
	flag.IntVar(&width, "width", defaultWidth, "Grid width")
	flag.IntVar(&height, "height", defaultHeight, "Grid height")
	flag.IntVar(&speed, "speed", defaultSpeed, "Animation speed in milliseconds")
	flag.IntVar(&generations, "generations", defaultGenerations, "Number of generations to simulate")
	flag.StringVar(&pattern, "pattern", "", "Pattern to load (use 'list' to see available patterns)")
	flag.BoolVar(&showStats, "stats", false, "Show statistics during simulation")
}

func main() {
	flag.Parse()

	// Handle pattern list request
	if pattern == "list" {
		fmt.Print(listPatterns())
		return
	}

	// Validate parameters
	if width <= 0 || height <= 0 {
		fmt.Println("Error: width and height must be positive integers")
		flag.Usage()
		return
	}
	if speed <= 0 {
		fmt.Println("Error: speed must be a positive integer")
		flag.Usage()
		return
	}
	if generations <= 0 {
		fmt.Println("Error: generations must be a positive integer")
		flag.Usage()
		return
	}

	// Update global dimensions
	DX = width
	DY = height

	// Initialize matrix
	var matrix [][]int
	var err error

	if pattern != "" {
		// Load predefined pattern
		matrix, err = loadPattern(pattern)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			fmt.Print(listPatterns())
			return
		}
	} else {
		// Use random initialization
		matrix = randomize()
	}

	termboxErr := termbox.Init()
	if termboxErr != nil {
		panic(termboxErr)
	}
	defer termbox.Close()

	termboxErr = termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	if termboxErr != nil {
		panic(termboxErr)
	}

	// Initialize statistics
	stats := Statistics{
		Generation:    0,
		LivingCells:   countLivingCells(matrix),
		StartTime:     time.Now(),
		LastFrameTime: time.Now(),
	}

	for i := 0; i < generations; i++ {
		prevLivingCells := stats.LivingCells
		matrix = step(matrix)
		updateStatistics(&stats, matrix, prevLivingCells)

		termboxErr = flush(matrix)
		if termboxErr != nil {
			panic(termboxErr)
		}

		displayStatistics(stats)
		termboxErr = termbox.Flush()
		if termboxErr != nil {
			panic(termboxErr)
		}

		time.Sleep(time.Duration(speed) * time.Millisecond)
	}
}
