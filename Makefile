# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

# Binary name
BINARY_NAME=golife
BINARY_UNIX=$(BINARY_NAME)_unix
BINARY_WINDOWS=$(BINARY_NAME).exe

# Build directory
BUILD_DIR=bin

.PHONY: all build clean test coverage run run-auto run-pattern demo demo-multi demo-25d demo-3d web-viewer build-wasm wasm-test wasm-api-test wasm-3d-viewer help fmt vet quality deps tidy

# Default target
all: help

## help: Display this help message
help:
	@echo "Available targets:"
	@echo "  make build        - Build the binary"
	@echo "  make run          - Run with interactive mode and all features"
	@echo "  make run-auto     - Run automatic mode (100 generations)"
	@echo "  make run-pattern  - Run with Gosper's Glider Gun pattern"
	@echo "  make demo         - Show all available demos"
	@echo "  make demo-multi   - Run multi-layer 2.5D visualization demo"
	@echo "  make demo-25d     - Run 2.5D patterns catalog demo"
	@echo "  make demo-3d      - Run WebGL 3D viewer demo (http://localhost:8080)"
	@echo "  make web-viewer   - Same as demo-3d"
	@echo "  make test         - Run tests"
	@echo "  make coverage     - Run tests with coverage report"
	@echo "  make quality      - Run all quality checks (fmt, vet, test)"
	@echo "  make clean        - Remove build artifacts"
	@echo "  make fmt          - Format Go code"
	@echo "  make vet          - Run go vet"
	@echo "  make deps         - Download dependencies"
	@echo "  make tidy         - Tidy go.mod"
	@echo "  make build-linux  - Build for Linux"
	@echo "  make build-windows- Build for Windows"
	@echo "  make build-all    - Build for all platforms"
	@echo "  make build-wasm   - Build WASM binary"
	@echo "  make wasm-test    - Run basic WASM test server (http://localhost:8081)"
	@echo "  make wasm-api-test- Run WASM API test server (http://localhost:8082)"
	@echo "  make wasm-3d-viewer- Run WASM 3D viewer (http://localhost:8083)"

## build: Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	cd cmd/golife && $(GOBUILD) -o ../../$(BUILD_DIR)/$(BINARY_NAME) -v

## run: Run the program with interactive mode and all features
run: build
	@echo "Running $(BINARY_NAME) in interactive mode with all features..."
	@./$(BUILD_DIR)/$(BINARY_NAME) --interactive --stats --color age --width 80 --height 40 --speed 100

## run-auto: Run automatic mode for 100 generations
run-auto: build
	@echo "Running $(BINARY_NAME) in automatic mode..."
	@./$(BUILD_DIR)/$(BINARY_NAME) --stats --color age --width 80 --height 40 --speed 50 --generations 100

## run-pattern: Run with Gosper's Glider Gun pattern in interactive mode
run-pattern: build
	@echo "Running $(BINARY_NAME) with Gosper's Glider Gun pattern..."
	@./$(BUILD_DIR)/$(BINARY_NAME) --interactive --stats --color age --pattern glider-gun --speed 100

## demo-multi: Run multi-layer 2.5D visualization demo
demo-multi:
	@echo "Running multi-layer 2.5D visualization demo..."
	$(GOCMD) run examples/multi-layer/main.go

## demo-25d: Run 2.5D patterns catalog demo
demo-25d:
	@echo "Running 2.5D patterns catalog demo..."
	$(GOCMD) run examples/25d-patterns/main.go

## demo-3d: Run WebGL 3D viewer demo (opens at http://localhost:8080)
demo-3d:
	@echo "=========================================="
	@echo "  3D Game of Life - WebGL Demo"
	@echo "=========================================="
	@echo ""
	@echo "Building 3D viewer..."
	@mkdir -p $(BUILD_DIR)
	@cd cmd/web-viewer && $(GOBUILD) -o ../../$(BUILD_DIR)/web-viewer
	@echo ""
	@echo "Starting 3D WebGL server..."
	@echo "  URL: http://localhost:8080"
	@echo ""
	@echo "Features:"
	@echo "  - 32x32x32 3D universe"
	@echo "  - Bays's Glider pattern (3D glider)"
	@echo "  - B6/S567 rule (3D Life)"
	@echo "  - Real-time WebGL visualization"
	@echo ""
	@echo "Controls:"
	@echo "  - Mouse drag: Rotate camera"
	@echo "  - Mouse wheel: Zoom in/out"
	@echo "  - Right-click drag: Pan view"
	@echo ""
	@echo "Press Ctrl+C to stop the server"
	@echo "=========================================="
	@echo ""
	@./$(BUILD_DIR)/web-viewer

## web-viewer: Run WebGL 3D viewer (opens at http://localhost:8080)
web-viewer:
	@echo "Building and starting WebGL 3D viewer..."
	@mkdir -p $(BUILD_DIR)
	@cd cmd/web-viewer && $(GOBUILD) -o ../../$(BUILD_DIR)/web-viewer
	@echo "Starting server at http://localhost:8080"
	@echo "Press Ctrl+C to stop"
	@./$(BUILD_DIR)/web-viewer

## test: Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

## coverage: Run tests with coverage report
coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## fmt: Format Go code
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...

## vet: Run go vet
vet:
	@echo "Running go vet..."
	$(GOVET) ./...

## quality: Run all quality checks
quality: fmt vet test
	@echo "All quality checks passed!"

## clean: Remove build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOGET) -v ./...

## tidy: Tidy go.mod
tidy:
	@echo "Tidying go.mod..."
	$(GOMOD) tidy

## build-linux: Build for Linux
build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	cd cmd/golife && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o ../../$(BUILD_DIR)/$(BINARY_UNIX) -v

## build-windows: Build for Windows
build-windows:
	@echo "Building for Windows..."
	@mkdir -p $(BUILD_DIR)
	cd cmd/golife && CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o ../../$(BUILD_DIR)/$(BINARY_WINDOWS) -v

## build-all: Build for all platforms
build-all: build build-linux build-windows
	@echo "Built binaries for all platforms in $(BUILD_DIR)/"

## build-wasm: Build WASM binary
build-wasm:
	@echo "Building WASM binary..."
	@GOOS=js GOARCH=wasm $(GOBUILD) -o web/life3d.wasm cmd/wasm-life/main.go
	@cp "$$(go env GOROOT)/lib/wasm/wasm_exec.js" web/ 2>/dev/null || \
		echo "Note: wasm_exec.js already present"
	@ls -lh web/life3d.wasm
	@echo "WASM binary built successfully"

## wasm-test: Run basic WASM test server
wasm-test: build-wasm
	@echo "=========================================="
	@echo "  WASM Basic Test Server"
	@echo "=========================================="
	@echo ""
	@echo "Starting HTTP server on :8081"
	@echo "  URL: http://localhost:8081/wasm-test.html"
	@echo ""
	@echo "Press Ctrl+C to stop"
	@echo "=========================================="
	@echo ""
	@cd web && python3 -m http.server 8081

## wasm-api-test: Run WASM API test server
wasm-api-test: build-wasm
	@echo "=========================================="
	@echo "  WASM API Test Server"
	@echo "=========================================="
	@echo ""
	@echo "Starting HTTP server on :8082"
	@echo "  URL: http://localhost:8082/wasm-api-test.html"
	@echo ""
	@echo "Features:"
	@echo "  - Interactive API testing"
	@echo "  - Pattern loading (7 patterns)"
	@echo "  - Step-by-step simulation"
	@echo "  - Animation control"
	@echo "  - Cell inspection"
	@echo ""
	@echo "API Functions:"
	@echo "  - initUniverse(width, height, depth)"
	@echo "  - loadPattern(name, x, y, z)"
	@echo "  - step()"
	@echo "  - getLivingCells()"
	@echo "  - getUniverseInfo()"
	@echo "  - clearUniverse()"
	@echo "  - setCell(x, y, z, alive)"
	@echo ""
	@echo "Press Ctrl+C to stop"
	@echo "=========================================="
	@echo ""
	@cd web && python3 -m http.server 8082

## wasm-3d-viewer: Run WASM 3D viewer
wasm-3d-viewer: build-wasm
	@echo "=========================================="
	@echo "  WASM 3D Viewer (Three.js)"
	@echo "=========================================="
	@echo ""
	@echo "Starting HTTP server on :8083"
	@echo "  URL: http://localhost:8083/wasm-3d-viewer.html"
	@echo ""
	@echo "Features:"
	@echo "  - Three.js 3D visualization"
	@echo "  - 10×10×10 default universe"
	@echo "  - Scalable to 15×15×15, 20×20×20"
	@echo "  - 7 pattern presets (Glider, Block, etc.)"
	@echo "  - Real-time animation"
	@echo "  - Interactive camera controls"
	@echo ""
	@echo "Controls:"
	@echo "  - Mouse drag: Rotate camera"
	@echo "  - Mouse wheel: Zoom in/out"
	@echo "  - Right-click drag: Pan view"
	@echo ""
	@echo "Press Ctrl+C to stop"
	@echo "=========================================="
	@echo ""
	@cd web && python3 -m http.server 8083
