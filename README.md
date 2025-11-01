# golife

Conway's Game of Life implementation in Go with terminal UI using termbox-go.

[![CI](https://github.com/nyasuto/golife/actions/workflows/ci.yml/badge.svg)](https://github.com/nyasuto/golife/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/nyasuto/golife)](https://goreportcard.com/report/github.com/nyasuto/golife)

## Features

- 🎮 **Configurable Parameters**: Customize grid size, animation speed, and generation count via command-line flags
- 🎨 **Pattern Presets**: Built-in famous Life patterns (Glider, Pulsar, Gosper's Glider Gun, etc.)
- 🌐 **WebGL 3D Viewer**: Real-time 3D visualization using Three.js with Bays's Glider
- 📊 **Multi-dimensional Support**: 2D, 2.5D, and 3D cellular automata
- 🧪 **Comprehensive Tests**: Extensive unit tests with high coverage
- ⚡ **Performance**: Optimized algorithms with parallel processing (4.5-6.6x speedup)
- 🛠️ **Modern Development**: CI/CD pipeline with automated testing, linting, and quality checks

## Installation

```bash
# Clone the repository
git clone https://github.com/nyasuto/golife.git
cd golife

# Build the application
make build

# Or build directly with go
go build -o bin/golife
```

## Usage

### Basic Usage

```bash
# Run with default settings (100x40 grid, 200ms speed, 300 generations)
./bin/golife

# Or use make
make run
```

### Command-Line Options

```bash
# Custom grid size
./bin/golife --width=150 --height=50

# Adjust animation speed (milliseconds)
./bin/golife --speed=100

# Set number of generations
./bin/golife --generations=500

# Combine multiple options
./bin/golife --width=120 --height=45 --speed=150 --generations=1000
```

### Pattern Presets

golife includes famous Game of Life patterns as presets:

```bash
# List available patterns
./bin/golife --pattern=list

# Run specific patterns
./bin/golife --pattern=glider
./bin/golife --pattern=pulsar
./bin/golife --pattern=glider-gun --width=150 --height=60 --generations=500

# Available patterns:
# - glider: A small pattern that moves diagonally
# - blinker: A period-2 oscillator
# - toad: A period-2 oscillator
# - beacon: A period-2 oscillator
# - pulsar: A period-3 oscillator
# - glider-gun: Gosper's Glider Gun (continuously generates gliders)
```

### 3D WebGL Viewer

Experience Game of Life in 3D with real-time WebGL visualization:

```bash
# Start the WebGL viewer
make web-viewer

# Or run directly
./bin/web-viewer
```

Then open http://localhost:8080 in your browser.

**Features:**
- 🎬 Real-time 3D voxel rendering with Three.js
- 🔄 WebSocket streaming for live updates
- 🎮 Interactive camera controls (orbit, zoom, pan)
- 🎨 Gradient coloring based on Z-depth
- 📊 Live statistics (generation, population, FPS)
- ⚡ Instanced rendering for performance
- 🧬 Simulates Bays's Glider (10-cell, period-4) in B6/S567 rule

**Controls:**
- **Mouse drag**: Rotate camera
- **Mouse wheel**: Zoom in/out
- **Right click drag**: Pan view

## Development

### Prerequisites

- Go 1.25 or later
- Make (optional, for using Makefile commands)

### Building from Source

```bash
# Install dependencies
go mod download

# Build the application
make build

# Run tests
make test

# Run all quality checks (format, vet, test)
make quality

# Generate coverage report
make coverage

# Clean build artifacts
make clean
```

### Available Make Commands

```bash
make help          # Show all available commands
make build         # Build the application
make test          # Run tests
make coverage      # Generate coverage report
make quality       # Run all quality checks
make run           # Build and run the application
make clean         # Clean build artifacts
```

## Game of Life Rules

Conway's Game of Life follows these simple rules:

1. **Survival**: Any live cell with 2 or 3 live neighbors survives
2. **Birth**: Any dead cell with exactly 3 live neighbors becomes alive
3. **Death**: All other cells die or stay dead

Despite these simple rules, complex and fascinating patterns emerge!

## Project Structure

```
golife/
├── main.go              # Main application logic
├── main_test.go         # Tests for core Game of Life logic
├── patterns_test.go     # Tests for pattern presets
├── Makefile            # Build and development commands
├── go.mod              # Go module definition
├── CLAUDE.md           # Development guidelines
└── README.md           # This file
```

## CI/CD

This project uses GitHub Actions for continuous integration:

- **Lint**: Code quality checks with golangci-lint
- **Test**: Automated testing with race detection and coverage reporting
- **Build**: Compilation verification across the codebase
- **Quality Checks**: go vet, go fmt, and go mod tidy validation

All checks must pass before merging pull requests.

## Contributing

1. Fork the repository
2. Create a feature branch (`feat/your-feature` or `fix/your-fix`)
3. Make your changes
4. Run `make quality` to ensure all checks pass
5. Commit your changes with conventional commit messages
6. Push to your fork and submit a pull request

See [CLAUDE.md](./CLAUDE.md) for detailed development guidelines.

## License

This project is open source. Feel free to use and modify as needed.

## Acknowledgments

- [John Conway](https://en.wikipedia.org/wiki/John_Horton_Conway) for creating the Game of Life
- [termbox-go](https://github.com/nsf/termbox-go) for terminal UI rendering
- The Go community for excellent tools and libraries
