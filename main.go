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

func step(data [][]int) [][]int {

	result := make([][]int, DY)

	for y := 0; y < DY; y++ {
		result[y] = make([]int, DX)
		for x := 0; x < DX; x++ {
			if data[y][x] == 1 {
				//生命がいるところは、周囲に２個、または３個の生命がいる場合に、そのまま生命が残ります。 そうでない場合には死んでしまいます。
				var check int
				if x == 0 && y == 0 {
					check = data[y][x+1] + data[y+1][x] + data[y+1][x+1]
				} else if x == DX-1 && y == DY-1 {
					check = data[y-1][x-1] + data[y-1][x] + data[y][x-1]
				} else if y == 0 && x == DX-1 {
					check = data[y][x-1] + data[y+1][x-1] + data[y+1][x]
				} else if x == 0 && y == DY-1 {
					check = data[y-1][x] + data[y-1][x+1] + data[y][x+1]
				} else if y == 0 {
					check = data[y][x-1] + data[y][x+1] + data[y+1][x-1] + data[y+1][x] + data[y+1][x+1]
				} else if y == DY-1 {
					check = data[y-1][x-1] + data[y-1][x] + data[y-1][x+1] + data[y][x-1] + data[y][x+1]
				} else if x == 0 {
					check = data[y-1][x] + data[y-1][x+1] + data[y][x+1] + data[y+1][x] + data[y+1][x+1]
				} else if x == DX-1 {
					check = data[y-1][x-1] + data[y-1][x] + data[y][x-1] + data[y+1][x-1] + data[y+1][x]
				} else {
					check = data[y-1][x-1] + data[y-1][x] + data[y-1][x+1] + data[y][x-1] + data[y][x+1] + data[y+1][x-1] + data[y+1][x] + data[y+1][x+1]
				}

				if check == 2 || check == 3 {
					result[y][x] = 1
				} else {
					result[y][x] = 0
				}

			} else {
				//生命のいないところには周囲にちょうど３個の生命がある場合に新しく生命が誕生します。
				var check int
				if x == 0 && y == 0 {
					check = data[y][x+1] + data[y+1][x] + data[y+1][x+1]
				} else if x == DX-1 && y == DY-1 {
					check = data[y-1][x-1] + data[y-1][x] + data[y][x-1]
				} else if y == 0 && x == DX-1 {
					check = data[y][x-1] + data[y+1][x-1] + data[y+1][x]
				} else if x == 0 && y == DY-1 {
					check = data[y-1][x] + data[y-1][x+1] + data[y][x+1]
				} else if y == 0 {
					check = data[y][x-1] + data[y][x+1] + data[y+1][x-1] + data[y+1][x] + data[y+1][x+1]
				} else if y == DY-1 {
					check = data[y-1][x-1] + data[y-1][x] + data[y-1][x+1] + data[y][x-1] + data[y][x+1]
				} else if x == 0 {
					check = data[y-1][x] + data[y-1][x+1] + data[y][x+1] + data[y+1][x] + data[y+1][x+1]
				} else if x == DX-1 {
					check = data[y-1][x-1] + data[y-1][x] + data[y][x-1] + data[y+1][x-1] + data[y+1][x]
				} else {
					check = data[y-1][x-1] + data[y-1][x] + data[y-1][x+1] + data[y][x-1] + data[y][x+1] + data[y+1][x-1] + data[y+1][x] + data[y+1][x+1]
				}

				if check == 3 {
					result[y][x] = 1
				} else {
					result[y][x] = 0
				}
			}
			//result[y][x] = r.Intn(2)

		}
	}
	return result
}

func flush(data [][]int) {
	for y := 0; y < DY; y++ {
		for x := 0; x < DX; x++ {
			var dot = ' '
			if data[y][x] == 1 {
				dot = '*'
			}
			termbox.SetCell(x, y, dot, termbox.ColorDefault, termbox.ColorDefault)

		}
	}

	termbox.Flush()

}

func main() {
	var matrix = randomize()

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	for i := 0; i < 300; i++ {
		matrix = step(matrix)
		flush(matrix)
		time.Sleep(200 * time.Millisecond)
	}

	defer termbox.Close()
}
