#!/bin/bash
set -e

# GitAI Installer Script
# This script installs the latest version of GitAI for macOS and Linux

REPO="xyue92/gitai"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="gitai"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Print colored output
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# Detect OS and architecture
detect_platform() {
    local os=$(uname -s | tr '[:upper:]' '[:lower:]')
    local arch=$(uname -m)

    case "$os" in
        linux*)
            OS="linux"
            ;;
        darwin*)
            OS="darwin"
            ;;
        *)
            print_error "Unsupported operating system: $os"
            exit 1
            ;;
    esac

    case "$arch" in
        x86_64|amd64)
            ARCH="amd64"
            ;;
        aarch64|arm64)
            ARCH="arm64"
            ;;
        *)
            print_error "Unsupported architecture: $arch"
            exit 1
            ;;
    esac

    PLATFORM="${OS}-${ARCH}"
    print_info "Detected platform: $PLATFORM"
}

# Get latest release version
get_latest_version() {
    print_info "Fetching latest release..."
    VERSION=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

    if [ -z "$VERSION" ]; then
        print_error "Failed to get latest version"
        exit 1
    fi

    print_info "Latest version: $VERSION"
}

# Download binary
download_binary() {
    DOWNLOAD_FILENAME="${BINARY_NAME}-${PLATFORM}"
    DOWNLOAD_URL="https://github.com/$REPO/releases/download/$VERSION/${DOWNLOAD_FILENAME}"
    CHECKSUM_URL="${DOWNLOAD_URL}.sha256"

    TMP_DIR=$(mktemp -d)
    TMP_FILE="$TMP_DIR/$DOWNLOAD_FILENAME"
    TMP_CHECKSUM="$TMP_DIR/${DOWNLOAD_FILENAME}.sha256"

    print_info "Downloading from: $DOWNLOAD_URL"

    if ! curl -fsSL "$DOWNLOAD_URL" -o "$TMP_FILE"; then
        print_error "Failed to download binary"
        rm -rf "$TMP_DIR"
        exit 1
    fi

    # Download checksum
    if curl -fsSL "$CHECKSUM_URL" -o "$TMP_CHECKSUM" 2>/dev/null; then
        print_info "Verifying checksum..."
        cd "$TMP_DIR"
        if command -v sha256sum >/dev/null 2>&1; then
            sha256sum -c "${DOWNLOAD_FILENAME}.sha256" >/dev/null 2>&1 || {
                print_error "Checksum verification failed"
                rm -rf "$TMP_DIR"
                exit 1
            }
        elif command -v shasum >/dev/null 2>&1; then
            shasum -a 256 -c "${DOWNLOAD_FILENAME}.sha256" >/dev/null 2>&1 || {
                print_error "Checksum verification failed"
                rm -rf "$TMP_DIR"
                exit 1
            }
        else
            print_warning "No checksum tool found, skipping verification"
        fi
        print_info "Checksum verified successfully"
    else
        print_warning "Checksum file not found, skipping verification"
    fi
}

# Install binary
install_binary() {
    chmod +x "$TMP_FILE"

    # Check if we need sudo
    if [ -w "$INSTALL_DIR" ]; then
        mv "$TMP_FILE" "$INSTALL_DIR/$BINARY_NAME"
    else
        print_info "Installing to $INSTALL_DIR (requires sudo)..."
        sudo mv "$TMP_FILE" "$INSTALL_DIR/$BINARY_NAME"
    fi

    rm -rf "$TMP_DIR"
    print_info "Installation completed successfully!"
}

# Verify installation
verify_installation() {
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        INSTALLED_VERSION=$($BINARY_NAME --version 2>/dev/null || echo "unknown")
        print_info "GitAI installed: $INSTALLED_VERSION"
        echo ""
        echo -e "${GREEN}✓${NC} GitAI has been installed to $INSTALL_DIR/$BINARY_NAME"
        echo ""
        echo "Try it out:"
        echo "  gitai --help"
    else
        print_error "Installation verification failed"
        echo "Please make sure $INSTALL_DIR is in your PATH"
        exit 1
    fi
}

# Main installation flow
main() {
    echo "╔════════════════════════════════════════╗"
    echo "║   GitAI Installer                      ║"
    echo "╚════════════════════════════════════════╝"
    echo ""

    detect_platform
    get_latest_version
    download_binary
    install_binary
    verify_installation

    echo ""
    echo -e "${GREEN}Next steps:${NC}"
    echo "1. Make sure Ollama is installed: https://ollama.com"
    echo "2. Pull an AI model: ollama pull qwen2.5-coder:7b"
    echo "3. Navigate to a git repository and run: gitai commit"
    echo ""
}

main
