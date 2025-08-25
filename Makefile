.PHONY: build test fmt lint clean all help

# Default target
all: fmt lint build test

# Build the binary
build:
	@echo "Building git-get..."
	go build -o git-get ./cmd/

# Run tests with race detection
test:
	@echo "Running tests..."
	go test -race ./...

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run linter with auto-fix
lint:
	@echo "Running linter..."
	golangci-lint run --fix

# Clean built binaries
clean:
	@echo "Cleaning..."
	rm -f git-get

# Show help
help:
	@echo "Available targets:"
	@echo "  build  - Build the git-get binary"
	@echo "  test   - Run tests with race detection"
	@echo "  fmt    - Format Go code"
	@echo "  lint   - Run golangci-lint with auto-fix"
	@echo "  clean  - Remove built binaries"
	@echo "  all    - Run fmt, lint, build, and test (default)"
	@echo "  help   - Show this help message"