package terminal

import (
	"fmt"

	"golife/pkg/core"
	"golife/pkg/engine"
	"golife/pkg/universe"

	termbox "github.com/nsf/termbox-go"
)

// Renderer2D renders a 2D universe to the terminal
type Renderer2D struct {
	showStats bool
	colorMode string
}

// NewRenderer2D creates a new 2D renderer
func NewRenderer2D(showStats bool, colorMode string) *Renderer2D {
	return &Renderer2D{
		showStats: showStats,
		colorMode: colorMode,
	}
}

// Render renders the universe to the terminal
func (r *Renderer2D) Render(u *universe.Universe2D, stats *engine.Statistics, showHelp bool) error {
	if r.colorMode == "age" {
		if err := r.renderWithColor(u); err != nil {
			return err
		}
	} else {
		if err := r.renderPlain(u); err != nil {
			return err
		}
	}

	if r.showStats {
		r.displayStatistics(stats, u.Width(), u.Height())
	}

	if showHelp {
		r.displayHelp(u.Height())
	}

	return termbox.Flush()
}

// renderPlain renders the universe without colors
func (r *Renderer2D) renderPlain(u *universe.Universe2D) error {
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
	return nil
}

// renderWithColor renders the universe with age-based colors
func (r *Renderer2D) renderWithColor(u *universe.Universe2D) error {
	for y := 0; y < u.Height(); y++ {
		for x := 0; x < u.Width(); x++ {
			var dot = ' '
			color := termbox.ColorDefault
			coord := core.NewCoord2D(x, y)

			if u.Get(coord) != core.Dead {
				dot = '*'
				age := u.GetAge(x, y)
				color = getColorByAge(age)
			}
			termbox.SetCell(x, y, dot, color, termbox.ColorDefault)
		}
	}
	return nil
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

// displayStatistics displays statistics on the screen
func (r *Renderer2D) displayStatistics(stats *engine.Statistics, width, height int) {
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
func (r *Renderer2D) displayHelp(height int) {
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
