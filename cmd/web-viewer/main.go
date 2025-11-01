package main

import (
	"flag"
	"golife/pkg/core"
	"golife/pkg/patterns"
	"golife/pkg/rules"
	"golife/pkg/universe"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var (
	addr     = flag.String("addr", ":8080", "http service address")
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for development
		},
	}
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

func main() {
	flag.Parse()
	log.SetFlags(0)

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", handleWebSocket)

	log.Printf("Starting WebGL 3D Life viewer on %s", *addr)
	log.Printf("Open http://localhost%s in your browser", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "web/index.html")
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer func() {
		if err := ws.Close(); err != nil {
			log.Println("WebSocket close error:", err)
		}
	}()

	// Create 3D universe with Bays's Glider
	rule := rules.Life3D_B6S567{}
	size := 32
	u := universe.New3D(size, size, size, rule)

	// Load Bays's Glider in the center
	glider := patterns.BaysGlider()
	glider.LoadIntoUniverse3D(u, size/2-2, size/2-2, size/2-2)

	generation := 0
	ticker := time.NewTicker(100 * time.Millisecond) // 10 FPS
	defer ticker.Stop()

	log.Println("WebSocket client connected")

	for range ticker.C {
		// Extract living cells
		state := extractUniverseState(u, generation)

		// Send to client
		if err := ws.WriteJSON(state); err != nil {
			log.Println("WebSocket write error:", err)
			return
		}

		// Step simulation (use parallel for better performance)
		u.StepParallel()
		generation++

		// Check for websocket close
		if err := ws.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
			log.Println("SetReadDeadline error:", err)
			return
		}
		if _, _, err := ws.NextReader(); err != nil {
			// Client disconnected
			break
		}
	}

	log.Println("WebSocket client disconnected")
}

func extractUniverseState(u *universe.Universe3D, generation int) UniverseState {
	size := u.Size()
	cells := make([]CellData, 0, u.CountLiving())

	// Iterate through all cells and collect living ones
	for z := 0; z < size.Z; z++ {
		for y := 0; y < size.Y; y++ {
			for x := 0; x < size.X; x++ {
				coord := core.NewCoord3D(x, y, z)
				if u.Get(coord) != core.Dead {
					cells = append(cells, CellData{X: x, Y: y, Z: z})
				}
			}
		}
	}

	return UniverseState{
		Cells:      cells,
		Generation: generation,
		Population: u.CountLiving(),
		Width:      size.X,
		Height:     size.Y,
		Depth:      size.Z,
	}
}
