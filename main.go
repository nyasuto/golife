package main

import (
	"flag"
	"fmt"
	"time"

	"golife/pkg/core"
	"golife/pkg/patterns"
	"golife/pkg/rules"
	"golife/pkg/universe"

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
	width        int
	height       int
	speed        int
	generations  int
	pattern      string
	showStats    bool
	colorMode    string
	interactive  bool
	currentSpeed int // Current speed in interactive mode
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

// listPatterns returns a formatted string of all available patterns
func listPatterns() string {
	patterns := patterns.AllPatterns()
	result := "Available patterns:\n"
	for name, p := range patterns {
		result += fmt.Sprintf("  %s: %s\n", name, p.Description)
	}
	return result
}

// loadPattern loads a predefined pattern into the universe
func loadPattern(u *universe.Universe2D, patternName string) error {
	allPatterns := patterns.AllPatterns()
	p, exists := allPatterns[patternName]
	if !exists {
		return fmt.Errorf("pattern '%s' not found", patternName)
	}

	// Calculate center position
	startX := (u.Width() - p.Width) / 2
	startY := (u.Height() - p.Height) / 2

	// Load pattern into universe
	p.LoadIntoUniverse(u, startX, startY)
	return nil
}

// updateStatistics updates the statistics for the current generation
func updateStatistics(stats *Statistics, u *universe.Universe2D, prevLivingCells int) {
	stats.Generation++
	stats.LivingCells = u.CountLiving()

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

// getColorByAge returns the color based on cell age
func getColorByAge(age int) termbox.Attribute {
	if age == 0 {
		return termbox.ColorDefault
	} else if age == 1 {
		return termbox.ColorGreen // Newborn
	} else if age <= 3 {
		return termbox.ColorYellow // Young
	} else if age <= 10 {
		return termbox.ColorRed // Old
	} else {
		return termbox.ColorBlue // Very old
	}
}

// flush renders the universe to the terminal
func flush(u *universe.Universe2D) error {
	for y := 0; y < u.Height(); y++ {
		for x := 0; x < u.Width(); x++ {
			var dot = ' '
			coord := core.NewCoord2D(x, y)
			if u.Get(coord) != core.Dead {
				dot = '*'
			}
			termbox.SetCell(x, y, dot, termbox.ColorDefault, termbox.ColorDefault)
		}
	}
	return termbox.Flush()
}

// flushWithColor renders the universe with age-based colors
func flushWithColor(u *universe.Universe2D) error {
	for y := 0; y < u.Height(); y++ {
		for x := 0; x < u.Width(); x++ {
			var dot = ' '
			color := termbox.ColorDefault
			coord := core.NewCoord2D(x, y)

			if u.Get(coord) != core.Dead {
				dot = '*'
				if colorMode == "age" {
					age := u.GetAge(x, y)
					color = getColorByAge(age)
				} else {
					color = termbox.ColorDefault
				}
			}
			termbox.SetCell(x, y, dot, color, termbox.ColorDefault)
		}
	}
	return termbox.Flush()
}

// displayStatistics displays statistics on the screen
func displayStatistics(stats Statistics, width, height int) {
	if !showStats {
		return
	}

	// Calculate position (top-right corner)
	startX := width - 35
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

// displayHelp displays keyboard controls
func displayHelp(height int) {
	startX := 2
	startY := height - 10

	if startY < 0 {
		return
	}

	lines := []string{
		"╔═══════════════════════════════╗",
		"║ Controls:                     ║",
		"║ Space - Pause/Resume          ║",
		"║ n     - Next step             ║",
		"║ +/-   - Speed up/down         ║",
		"║ r     - Restart (random)      ║",
		"║ q/Esc - Quit                  ║",
		"╚═══════════════════════════════╝",
	}

	for i, line := range lines {
		for j, ch := range line {
			termbox.SetCell(startX+j, startY+i, ch, termbox.ColorYellow, termbox.ColorDefault)
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
	flag.StringVar(&colorMode, "color", "", "Color mode: 'age' for age-based coloring")
	flag.BoolVar(&interactive, "interactive", false, "Enable interactive mode (keyboard controls)")
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

	// Create universe with Conway's rule
	rule := rules.ConwayRule{}
	u := universe.New2D(width, height, rule)

	// Initialize universe
	if pattern != "" {
		// Load predefined pattern
		err := loadPattern(u, pattern)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			fmt.Print(listPatterns())
			return
		}
	} else {
		// Use random initialization
		u.Randomize()
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
		LivingCells:   u.CountLiving(),
		StartTime:     time.Now(),
		LastFrameTime: time.Now(),
	}

	currentSpeed = speed

	if interactive {
		runInteractive(u, stats)
	} else {
		runAutomatic(u, stats, generations)
	}
}

// runAutomatic runs the simulation automatically for a fixed number of generations
func runAutomatic(u *universe.Universe2D, stats Statistics, gens int) {
	useColorMode := colorMode == "age"

	for i := 0; i < gens; i++ {
		prevLivingCells := stats.LivingCells
		u.Step()
		updateStatistics(&stats, u, prevLivingCells)

		if useColorMode {
			_ = flushWithColor(u)
		} else {
			_ = flush(u)
		}

		displayStatistics(stats, u.Width(), u.Height())
		_ = termbox.Flush()

		time.Sleep(time.Duration(speed) * time.Millisecond)
	}
}

// runInteractive runs the simulation in interactive mode with keyboard controls
func runInteractive(u *universe.Universe2D, stats Statistics) {
	paused := false
	running := true
	eventQueue := make(chan termbox.Event)

	// Start event polling goroutine
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	// Initial render
	render(u, stats)

	for running {
		select {
		case ev := <-eventQueue:
			if ev.Type == termbox.EventKey {
				switch ev.Key {
				case termbox.KeyEsc:
					running = false
				case termbox.KeySpace:
					paused = !paused
				default:
					switch ev.Ch {
					case 'q':
						running = false
					case ' ':
						paused = !paused
					case 'n':
						// Step once
						if paused {
							doStep(u, &stats)
							render(u, stats)
						}
					case '+', '=':
						// Speed up (decrease delay)
						if currentSpeed > 10 {
							currentSpeed -= 10
						}
					case '-', '_':
						// Slow down (increase delay)
						if currentSpeed < 1000 {
							currentSpeed += 10
						}
					case 'r':
						// Restart with random
						u.Clear()
						u.Randomize()
						stats = Statistics{
							Generation:    0,
							LivingCells:   u.CountLiving(),
							StartTime:     time.Now(),
							LastFrameTime: time.Now(),
						}
						render(u, stats)
					}
				}
			}

		default:
			if !paused {
				doStep(u, &stats)
				render(u, stats)
				time.Sleep(time.Duration(currentSpeed) * time.Millisecond)
			} else {
				time.Sleep(50 * time.Millisecond)
			}
		}
	}
}

// doStep performs one simulation step and updates statistics
func doStep(u *universe.Universe2D, stats *Statistics) {
	prevLivingCells := stats.LivingCells
	u.Step()
	updateStatistics(stats, u, prevLivingCells)
}

// render draws the current state to the screen
func render(u *universe.Universe2D, stats Statistics) {
	useColorMode := colorMode == "age"

	if useColorMode {
		_ = flushWithColor(u)
	} else {
		_ = flush(u)
	}

	displayStatistics(stats, u.Width(), u.Height())
	if interactive {
		displayHelp(u.Height())
	}
	_ = termbox.Flush()
}
