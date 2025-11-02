//go:build js && wasm
// +build js,wasm

package main

import (
	"fmt"
	"syscall/js"
)

func main() {
	fmt.Println("Go WASM initialized!")
	js.Global().Get("console").Call("log", "🚀 Go Life 3D WASM loaded successfully")

	// Keep the program running
	<-make(chan struct{})
}
