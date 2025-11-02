//go:build js && wasm
// +build js,wasm

package wasm

import (
	"encoding/json"
	"golife/pkg/core"
	"golife/pkg/patterns"
	"golife/pkg/rules"
	"golife/pkg/universe"
	"syscall/js"
)

// Global universe instance
var (
	currentUniverse *universe.Universe3D
	generation      int
)

// CellData represents a single living cell for JSON serialization
type CellData struct {
	X int `json:"x"`
	Y int `json:"y"`
	Z int `json:"z"`
}

// UniverseState represents the current state of the universe
type UniverseState struct {
	Cells      []CellData `json:"cells"`
	Generation int        `json:"generation"`
	Population int        `json:"population"`
	Width      int        `json:"width"`
	Height     int        `json:"height"`
	Depth      int        `json:"depth"`
}

// InitUniverse creates a new 3D universe with the given dimensions
// JavaScript call: initUniverse(width, height, depth)
func InitUniverse(this js.Value, args []js.Value) interface{} {
	if len(args) != 3 {
		return map[string]interface{}{
			"error": "initUniverse requires 3 arguments: width, height, depth",
		}
	}

	width := args[0].Int()
	height := args[1].Int()
	depth := args[2].Int()

	// Validate dimensions
	if width <= 0 || height <= 0 || depth <= 0 {
		return map[string]interface{}{
			"error": "dimensions must be positive",
		}
	}

	// Create universe with B6/S567 rule (3D Life)
	rule := rules.Life3D_B6S567{}
	currentUniverse = universe.New3D(width, height, depth, rule)
	generation = 0

	return map[string]interface{}{
		"success": true,
		"width":   width,
		"height":  height,
		"depth":   depth,
	}
}

// LoadPattern loads a predefined pattern into the universe
// JavaScript call: loadPattern(patternName, x, y, z)
func LoadPattern(this js.Value, args []js.Value) interface{} {
	if currentUniverse == nil {
		return map[string]interface{}{
			"error": "universe not initialized",
		}
	}

	if len(args) != 4 {
		return map[string]interface{}{
			"error": "loadPattern requires 4 arguments: patternName, x, y, z",
		}
	}

	patternName := args[0].String()
	x := args[1].Int()
	y := args[2].Int()
	z := args[3].Int()

	// Get pattern
	var pattern *patterns.Pattern3D
	switch patternName {
	case "glider":
		pattern = patterns.BaysGlider()
	case "block":
		pattern = patterns.Block3D()
	case "beehive":
		pattern = patterns.Beehive3D()
	case "blinker":
		pattern = patterns.Blinker3D()
	case "flashlight":
		pattern = patterns.Flashlight3D()
	case "wheel":
		pattern = patterns.Wheel3D()
	case "bucket":
		pattern = patterns.Bucket3D()
	default:
		return map[string]interface{}{
			"error": "unknown pattern: " + patternName,
		}
	}

	// Load pattern
	pattern.LoadIntoUniverse3D(currentUniverse, x, y, z)

	return map[string]interface{}{
		"success": true,
		"pattern": patternName,
		"x":       x,
		"y":       y,
		"z":       z,
	}
}

// Step advances the simulation by one generation
// JavaScript call: step()
func Step(this js.Value, args []js.Value) interface{} {
	if currentUniverse == nil {
		return map[string]interface{}{
			"error": "universe not initialized",
		}
	}

	// Use parallel step for better performance
	currentUniverse.StepParallel()
	generation++

	return map[string]interface{}{
		"success":    true,
		"generation": generation,
	}
}

// GetLivingCells returns all living cells in the universe
// JavaScript call: getLivingCells()
func GetLivingCells(this js.Value, args []js.Value) interface{} {
	if currentUniverse == nil {
		return map[string]interface{}{
			"error": "universe not initialized",
		}
	}

	state := extractUniverseState()

	// Convert to JSON string for JavaScript
	jsonBytes, err := json.Marshal(state)
	if err != nil {
		return map[string]interface{}{
			"error": "failed to marshal state: " + err.Error(),
		}
	}

	return string(jsonBytes)
}

// GetUniverseInfo returns basic information about the universe
// JavaScript call: getUniverseInfo()
func GetUniverseInfo(this js.Value, args []js.Value) interface{} {
	if currentUniverse == nil {
		return map[string]interface{}{
			"error": "universe not initialized",
		}
	}

	size := currentUniverse.Size()
	return map[string]interface{}{
		"width":      size.X,
		"height":     size.Y,
		"depth":      size.Z,
		"generation": generation,
		"population": currentUniverse.CountLiving(),
	}
}

// ClearUniverse resets all cells to dead state
// JavaScript call: clearUniverse()
func ClearUniverse(this js.Value, args []js.Value) interface{} {
	if currentUniverse == nil {
		return map[string]interface{}{
			"error": "universe not initialized",
		}
	}

	currentUniverse.Clear()
	generation = 0

	return map[string]interface{}{
		"success": true,
	}
}

// SetCell sets the state of a specific cell
// JavaScript call: setCell(x, y, z, alive)
func SetCell(this js.Value, args []js.Value) interface{} {
	if currentUniverse == nil {
		return map[string]interface{}{
			"error": "universe not initialized",
		}
	}

	if len(args) != 4 {
		return map[string]interface{}{
			"error": "setCell requires 4 arguments: x, y, z, alive",
		}
	}

	x := args[0].Int()
	y := args[1].Int()
	z := args[2].Int()
	alive := args[3].Bool()

	coord := core.NewCoord3D(x, y, z)
	if alive {
		currentUniverse.Set(coord, core.Alive)
	} else {
		currentUniverse.Set(coord, core.Dead)
	}

	return map[string]interface{}{
		"success": true,
	}
}

// extractUniverseState extracts the current state of the universe
func extractUniverseState() UniverseState {
	size := currentUniverse.Size()
	cells := make([]CellData, 0, currentUniverse.CountLiving())

	// Iterate through all cells and collect living ones
	for z := 0; z < size.Z; z++ {
		for y := 0; y < size.Y; y++ {
			for x := 0; x < size.X; x++ {
				coord := core.NewCoord3D(x, y, z)
				if currentUniverse.Get(coord) != core.Dead {
					cells = append(cells, CellData{X: x, Y: y, Z: z})
				}
			}
		}
	}

	return UniverseState{
		Cells:      cells,
		Generation: generation,
		Population: currentUniverse.CountLiving(),
		Width:      size.X,
		Height:     size.Y,
		Depth:      size.Z,
	}
}

// RegisterCallbacks registers all Go functions as JavaScript callbacks
func RegisterCallbacks() {
	js.Global().Set("goInitUniverse", js.FuncOf(InitUniverse))
	js.Global().Set("goLoadPattern", js.FuncOf(LoadPattern))
	js.Global().Set("goStep", js.FuncOf(Step))
	js.Global().Set("goGetLivingCells", js.FuncOf(GetLivingCells))
	js.Global().Set("goGetUniverseInfo", js.FuncOf(GetUniverseInfo))
	js.Global().Set("goClearUniverse", js.FuncOf(ClearUniverse))
	js.Global().Set("goSetCell", js.FuncOf(SetCell))
}
