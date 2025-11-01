package main

import (
	"flag"
	"fmt"
	"time"

	"golife/pkg/engine"
	"golife/pkg/patterns"
	"golife/pkg/rules"
	"golife/pkg/universe"
	"golife/pkg/visualizer/terminal"

	termbox "github.com/nsf/termbox-go"
)

// Default values
const (
	defaultWidth       = 100
	defaultHeight      = 40
	defaultSpeed       = 200
	defaultGenerations = 300
)

// Configuration holds command-line configuration
type Configuration struct {
	Width        int
	Height       int
	Speed        int
	Generations  int
	Pattern      string
	ShowStats    bool
	ColorMode    string
	Interactive  bool
	CurrentSpeed int
}

var config Configuration

func init() {
	flag.IntVar(&config.Width, "width", defaultWidth, "Grid width")
	flag.IntVar(&config.Height, "height", defaultHeight, "Grid height")
	flag.IntVar(&config.Speed, "speed", defaultSpeed, "Animation speed in milliseconds")
	flag.IntVar(&config.Generations, "generations", defaultGenerations, "Number of generations to simulate")
	flag.StringVar(&config.Pattern, "pattern", "", "Pattern to load (use 'list' to see available patterns)")
	flag.BoolVar(&config.ShowStats, "stats", false, "Show statistics during simulation")
	flag.StringVar(&config.ColorMode, "color", "", "Color mode: 'age' for age-based coloring")
	flag.BoolVar(&config.Interactive, "interactive", false, "Enable interactive mode (keyboard controls)")
}

func main() {
	flag.Parse()

	// Handle pattern list request
	if config.Pattern == "list" {
		fmt.Print(listPatterns())
		return
	}

	// Validate parameters
	if config.Width <= 0 || config.Height <= 0 {
		fmt.Println("Error: width and height must be positive integers")
		flag.Usage()
		return
	}
	if config.Speed <= 0 {
		fmt.Println("Error: speed must be a positive integer")
		flag.Usage()
		return
	}
	if config.Generations <= 0 {
		fmt.Println("Error: generations must be a positive integer")
		flag.Usage()
		return
	}

	// Create universe with Conway's rule
	rule := rules.ConwayRule{}
	u := universe.New2D(config.Width, config.Height, rule)

	// Initialize universe
	if config.Pattern != "" {
		if err := loadPattern(u, config.Pattern); err != nil {
			fmt.Printf("Error: %v\n", err)
			fmt.Print(listPatterns())
			return
		}
	} else {
		u.Randomize()
	}

	// Initialize termbox
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()

	if err := termbox.Clear(termbox.ColorDefault, termbox.ColorDefault); err != nil {
		panic(err)
	}

	// Initialize statistics and renderer
	stats := engine.NewStatistics(u.CountLiving())
	renderer := terminal.NewRenderer2D(config.ShowStats, config.ColorMode)
	config.CurrentSpeed = config.Speed

	// Run simulation
	if config.Interactive {
		runInteractive(u, stats, renderer)
	} else {
		runAutomatic(u, stats, renderer)
	}
}

func listPatterns() string {
	allPatterns := patterns.AllPatterns()
	result := "Available patterns:\n"
	for name, p := range allPatterns {
		result += fmt.Sprintf("  %s: %s\n", name, p.Description)
	}
	return result
}

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

func runAutomatic(u *universe.Universe2D, stats *engine.Statistics, renderer *terminal.Renderer2D) {
	for i := 0; i < config.Generations; i++ {
		u.Step()
		stats.Update(u)

		if err := renderer.Render(u, stats, false); err != nil {
			panic(err)
		}

		time.Sleep(time.Duration(config.Speed) * time.Millisecond)
	}
}

func runInteractive(u *universe.Universe2D, stats *engine.Statistics, renderer *terminal.Renderer2D) {
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
	_ = renderer.Render(u, stats, true)

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
						if paused {
							u.Step()
							stats.Update(u)
							_ = renderer.Render(u, stats, true)
						}
					case '+', '=':
						if config.CurrentSpeed > 10 {
							config.CurrentSpeed -= 10
						}
					case '-', '_':
						if config.CurrentSpeed < 1000 {
							config.CurrentSpeed += 10
						}
					case 'r':
						u.Clear()
						u.Randomize()
						stats.Reset(u.CountLiving())
						_ = renderer.Render(u, stats, true)
					}
				}
			}

		default:
			if !paused {
				u.Step()
				stats.Update(u)
				_ = renderer.Render(u, stats, true)
				time.Sleep(time.Duration(config.CurrentSpeed) * time.Millisecond)
			} else {
				time.Sleep(50 * time.Millisecond)
			}
		}
	}
}
