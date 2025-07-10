# BrewPy Makefile

BINARY_NAME=brewpy
VERSION?=1.0.0
BUILD_DIR=build
DIST_DIR=dist

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

LDFLAGS=-ldflags "-s -w"

.PHONY: all build clean deps help install uninstall release

all: clean deps build

help:
	@echo "BrewPy Build System"
	@echo ""
	@echo "Available targets:"
	@echo "  build     - Build the binary"
	@echo "  clean     - Clean build artifacts"
	@echo "  deps      - Download dependencies"
	@echo "  install   - Install locally"
	@echo "  uninstall - Remove local installation"
	@echo "  release   - Create release artifacts"
	@echo "  help      - Show this help message"

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR) $(DIST_DIR)
	@echo "Clean complete"

deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "Dependencies ready"

install: build
	@echo "Installing $(BINARY_NAME) locally..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "Installed to /usr/local/bin/$(BINARY_NAME)"
	@echo ""
	@echo "To complete setup, add this to your shell profile:"
	@echo "  eval \"\$$(brewpy init)\""

uninstall:
	@echo "Removing $(BINARY_NAME)..."
	@sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "Uninstalled"

release: clean deps
	@echo "Creating release artifacts for v$(VERSION)..."
	@mkdir -p $(DIST_DIR)
	
	@echo "Building for macOS ARM64..."
	@GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 .
	
	@echo "Building for macOS AMD64..."
	@GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 .
	
	@echo "Creating tarballs..."
	@cd $(DIST_DIR) && tar -czf $(BINARY_NAME)-v$(VERSION)-darwin-arm64.tar.gz $(BINARY_NAME)-darwin-arm64
	@cd $(DIST_DIR) && tar -czf $(BINARY_NAME)-v$(VERSION)-darwin-amd64.tar.gz $(BINARY_NAME)-darwin-amd64
	
	@echo "Generating checksums..."
	@cd $(DIST_DIR) && shasum -a 256 *.tar.gz > checksums.txt
	
	@echo "Release artifacts created in $(DIST_DIR)/"
	@echo ""
	@echo "Next steps:"
	@echo "1. Create a GitHub release with tag v$(VERSION)"
	@echo "2. Upload the tarball files"
	@echo "3. Update the Homebrew formula with new URL and SHA256"

dev:
	@echo "Building for development..."
	$(GOBUILD) -o $(BINARY_NAME) .
	@echo "Development build complete" 