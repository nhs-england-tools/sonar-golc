# GoLC CLI — Makefile
#
# Binary naming rationale:
#   golc   — matches the project name "GoLC" (Go Line Counter), short, memorable,
#            and the obvious choice for CLI invocation: `golc ./src`
#   sgolc  — adds a Sonar prefix but users would wonder what the 's' stands for
#   sloc   — conflicts with existing SLOC tools
#   loc    — too generic
#
# We go with `golc` as the default. Override with: make build BINARY_NAME=sgolc

BINARY_NAME ?= golc
BIN_DIR     := bin
CMD_DIR     := ./cmd/golc-cli
INSTALL_DIR ?= $(HOME)/.local/bin
VERSION     ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS     := -s -w -X main.version=$(VERSION)

# Platforms for cross-compilation (OS/ARCH)
PLATFORMS := \
	linux/amd64 \
	linux/arm64 \
	darwin/amd64 \
	darwin/arm64 \
	windows/amd64 \
	windows/arm64

# Default target
.PHONY: all
all: test build install

# ─── Test ────────────────────────────────────────────────────────────────────

.PHONY: test
test: test-pkg test-cli test-golc test-integration

# Package-level tests (no build tags required)
.PHONY: test-pkg
test-pkg:
	go test ./pkg/...

# CLI tests
.PHONY: test-cli
test-cli:
	go test ./cmd/...

# Root-level tests that require the 'golc' build tag
.PHONY: test-golc
test-golc:
	go test -tags golc -run . .

# Integration tests (no build tag, but separated for clarity)
.PHONY: test-integration
test-integration:
	go test -tags golc -run Integration .

# Run all tests with verbose output
.PHONY: test-v
test-v:
	go test -v ./pkg/...
	go test -v ./cmd/...
	go test -v -tags golc .

# ─── Build ───────────────────────────────────────────────────────────────────

.PHONY: build
build:
	@mkdir -p $(BIN_DIR)
	go build -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/$(BINARY_NAME) $(CMD_DIR)
	@echo "Built $(BIN_DIR)/$(BINARY_NAME)"

.PHONY: build-all
build-all:
	@mkdir -p $(BIN_DIR)
	@for platform in $(PLATFORMS); do \
		OS=$${platform%/*}; \
		ARCH=$${platform#*/}; \
		EXT=""; \
		if [ "$$OS" = "windows" ]; then EXT=".exe"; fi; \
		OUTPUT="$(BIN_DIR)/$(BINARY_NAME)-$$OS-$$ARCH$$EXT"; \
		echo "Building $$OUTPUT ..."; \
		GOOS=$$OS GOARCH=$$ARCH go build -ldflags "$(LDFLAGS)" -o $$OUTPUT $(CMD_DIR) || exit 1; \
	done
	@echo "All platform builds in $(BIN_DIR)/"

# ─── Utilities ───────────────────────────────────────────────────────────────

.PHONY: install
install: build
	install -d $(INSTALL_DIR)
	install -m 755 $(BIN_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Installed $(INSTALL_DIR)/$(BINARY_NAME)"

.PHONY: uninstall
uninstall:
	rm -f $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Removed $(INSTALL_DIR)/$(BINARY_NAME)"

.PHONY: clean
clean:
	rm -rf $(BIN_DIR)/golc*

.PHONY: fmt
fmt:
	go fmt ./cmd/...

.PHONY: vet
vet:
	go vet ./cmd/...

.PHONY: help
help:
	@echo "Usage:"
	@echo "  make              — run tests, build then install"
	@echo "  make test         — run all tests (pkg + cli + golc + integration)"
	@echo "  make test-pkg     — run package tests only"
	@echo "  make test-cli     — run CLI tests only"
	@echo "  make test-golc    — run golc-tagged root tests"
	@echo "  make test-v       — run all tests with verbose output"
	@echo "  make build        — build CLI binary to bin/"
	@echo "  make build-all    — cross-compile for all platforms"
	@echo "  make install      — build and install to INSTALL_DIR (default: $(HOME)/.local/bin)"
	@echo "  make uninstall    — remove from INSTALL_DIR"
	@echo "  make clean        — clean up build artifacts"
	@echo "  make fmt          — go fmt all source"
	@echo "  make vet          — go vet all source"
	@echo ""
	@echo "Variables:"
	@echo "  BINARY_NAME=golc  — override output binary name"
	@echo "  VERSION=v1.0      — override version stamp"
