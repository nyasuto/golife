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

.PHONY: all build clean test coverage run run-auto run-pattern demo-multi demo-25d web-viewer help fmt vet quality deps tidy

# Default target
all: help

## help: Display this help message
help:
	@echo "Available targets:"
	@echo "  make build        - Build the binary"
	@echo "  make run          - Run with interactive mode and all features"
	@echo "  make run-auto     - Run automatic mode (100 generations)"
	@echo "  make run-pattern  - Run with Gosper's Glider Gun pattern"
	@echo "  make demo-multi   - Run multi-layer 2.5D visualization demo"
	@echo "  make demo-25d     - Run 2.5D patterns catalog demo"
	@echo "  make web-viewer   - Run WebGL 3D viewer (http://localhost:8080)"
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
