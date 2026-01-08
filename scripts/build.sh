#!/bin/bash
set -e

# Local Build Script for GitAI
# Builds binaries for all supported platforms

VERSION=${1:-"dev"}
OUTPUT_DIR="dist"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║   GitAI Multi-Platform Build          ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════╝${NC}"
echo ""
echo "Version: $VERSION"
echo "Output directory: $OUTPUT_DIR"
echo ""

# Clean and create output directory
rm -rf "$OUTPUT_DIR"
mkdir -p "$OUTPUT_DIR"

# Build function
build() {
    local goos=$1
    local goarch=$2
    local output=$3

    echo -e "${GREEN}Building:${NC} $goos/$goarch -> $output"

    CGO_ENABLED=0 GOOS=$goos GOARCH=$goarch go build \
        -v \
        -trimpath \
        -ldflags="-s -w -X main.Version=$VERSION" \
        -o "$OUTPUT_DIR/$output"

    # Create checksum
    if command -v sha256sum >/dev/null 2>&1; then
        (cd "$OUTPUT_DIR" && sha256sum "$output" > "$output.sha256")
    elif command -v shasum >/dev/null 2>&1; then
        (cd "$OUTPUT_DIR" && shasum -a 256 "$output" > "$output.sha256")
    fi
}

# Build all platforms
echo "Building binaries..."
echo ""

build "linux" "amd64" "gitai-linux-amd64"
build "linux" "arm64" "gitai-linux-arm64"
build "darwin" "amd64" "gitai-darwin-amd64"
build "darwin" "arm64" "gitai-darwin-arm64"
build "windows" "amd64" "gitai-windows-amd64.exe"

echo ""
echo -e "${GREEN}✓ Build completed successfully!${NC}"
echo ""
echo "Output files:"
ls -lh "$OUTPUT_DIR"
echo ""
echo "Total size:"
du -sh "$OUTPUT_DIR"
