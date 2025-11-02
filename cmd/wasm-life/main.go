//go:build js && wasm
// +build js,wasm

package main

import (
	"fmt"
	"golife/pkg/wasm"
	"syscall/js"
)

func main() {
	fmt.Println("Go WASM initialized!")
	js.Global().Get("console").Call("log", "🚀 Go Life 3D WASM loaded successfully")

	// Register Go functions as JavaScript callbacks
	wasm.RegisterCallbacks()
	js.Global().Get("console").Call("log", "✅ Go functions registered:")
	js.Global().Get("console").Call("log", "  - goInitUniverse(width, height, depth)")
	js.Global().Get("console").Call("log", "  - goLoadPattern(name, x, y, z)")
	js.Global().Get("console").Call("log", "  - goStep()")
	js.Global().Get("console").Call("log", "  - goGetLivingCells()")
	js.Global().Get("console").Call("log", "  - goGetUniverseInfo()")
	js.Global().Get("console").Call("log", "  - goClearUniverse()")
	js.Global().Get("console").Call("log", "  - goSetCell(x, y, z, alive)")

	// Keep the program running
	<-make(chan struct{})
}
