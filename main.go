package main

import (
	"math/rand"
	"time"

	termbox "github.com/nsf/termbox-go"
)

// DX is width
var DX = 100

// DY is height
var DY = 40

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

func main() {
	var matrix = randomize()

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	err = termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 300; i++ {
		matrix = step(matrix)
		err = flush(matrix)
		if err != nil {
			panic(err)
		}

		time.Sleep(200 * time.Millisecond)
	}

	defer termbox.Close()
}
