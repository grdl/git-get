.PHONY: build test fmt lint man clean all help

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

# Generate man page from README
man:
	@echo "Generating man page from README..."
	@mkdir -p docs
	go tool go-md2man -in README.md -out docs/git-get.1
	@echo "Man page generated: docs/git-get.1"
	@echo "To install: sudo cp docs/git-get.1 /usr/share/man/man1/"

# Clean built binaries and generated files
clean:
	@echo "Cleaning..."
	rm -f git-get
	rm -f docs/git-get.1 docs/git-get.1.gz

# Show help
help:
	@echo "Available targets:"
	@echo "  build  - Build the git-get binary"
	@echo "  test   - Run tests with race detection"
	@echo "  fmt    - Format Go code"
	@echo "  lint   - Run golangci-lint with auto-fix"
	@echo "  man    - Generate man page from README"
	@echo "  clean  - Remove built binaries and generated files"
	@echo "  all    - Run fmt, lint, build, and test (default)"
	@echo "  help   - Show this help message"