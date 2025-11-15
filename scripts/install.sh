#!/bin/bash

# install.sh - Installation script for cntm
# Usage: curl -fsSL https://raw.githubusercontent.com/USER/REPO/main/scripts/install.sh | bash

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
VERSION="${CNTM_VERSION:-1.0.0}"
REPO="${CNTM_REPO:-yourusername/claude-nia-tool-management-cli}"
INSTALL_DIR="${CNTM_INSTALL_DIR:-/usr/local/bin}"
TMP_DIR=$(mktemp -d)

# Cleanup on exit
cleanup() {
    rm -rf "$TMP_DIR"
}
trap cleanup EXIT

echo -e "${BLUE}================================================${NC}"
echo -e "${BLUE}  Installing cntm v${VERSION}${NC}"
echo -e "${BLUE}================================================${NC}"
echo ""

# Detect OS and architecture
detect_platform() {
    local os arch

    # Detect OS
    case "$(uname -s)" in
        Linux*)     os="linux" ;;
        Darwin*)    os="darwin" ;;
        *)
            echo -e "${RED}Error: Unsupported operating system$(NC)"
            exit 1
            ;;
    esac

    # Detect architecture
    case "$(uname -m)" in
        x86_64)     arch="amd64" ;;
        amd64)      arch="amd64" ;;
        arm64)      arch="arm64" ;;
        aarch64)    arch="arm64" ;;
        *)
            echo -e "${RED}Error: Unsupported architecture$(NC)"
            exit 1
            ;;
    esac

    echo "${os}-${arch}"
}

PLATFORM=$(detect_platform)
echo -e "Detected platform: ${GREEN}${PLATFORM}${NC}"

# Check if running with sufficient privileges
check_privileges() {
    if [ ! -w "$INSTALL_DIR" ]; then
        echo -e "${YELLOW}Warning: No write permission for ${INSTALL_DIR}${NC}"
        echo -e "${YELLOW}You may need to run with sudo or choose a different install directory${NC}"
        echo ""
        read -p "Install to ~/bin instead? [Y/n] " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]] || [[ -z $REPLY ]]; then
            INSTALL_DIR="$HOME/bin"
            mkdir -p "$INSTALL_DIR"
        else
            echo -e "${RED}Installation cancelled${NC}"
            exit 1
        fi
    fi
}

check_privileges

# Download binary
BINARY_NAME="cntm-${PLATFORM}"
ARCHIVE_NAME="cntm-${VERSION}-${PLATFORM}.tar.gz"
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/v${VERSION}/${ARCHIVE_NAME}"

echo -e "Downloading from: ${BLUE}${DOWNLOAD_URL}${NC}"

if command -v curl &> /dev/null; then
    curl -fsSL -o "${TMP_DIR}/${ARCHIVE_NAME}" "${DOWNLOAD_URL}"
elif command -v wget &> /dev/null; then
    wget -q -O "${TMP_DIR}/${ARCHIVE_NAME}" "${DOWNLOAD_URL}"
else
    echo -e "${RED}Error: Neither curl nor wget found${NC}"
    exit 1
fi

if [ $? -ne 0 ]; then
    echo -e "${RED}Error: Failed to download cntm${NC}"
    echo -e "${YELLOW}Please check that version ${VERSION} exists at:${NC}"
    echo -e "  https://github.com/${REPO}/releases/tag/v${VERSION}"
    exit 1
fi

echo -e "${GREEN}  ✓ Downloaded${NC}"

# Download checksums
CHECKSUM_URL="https://github.com/${REPO}/releases/download/v${VERSION}/checksums.txt"

echo "Downloading checksums..."
if command -v curl &> /dev/null; then
    curl -fsSL -o "${TMP_DIR}/checksums.txt" "${CHECKSUM_URL}"
elif command -v wget &> /dev/null; then
    wget -q -O "${TMP_DIR}/checksums.txt" "${CHECKSUM_URL}"
fi

if [ $? -eq 0 ] && [ -f "${TMP_DIR}/checksums.txt" ]; then
    echo "Verifying checksum..."

    cd "$TMP_DIR"
    if command -v shasum &> /dev/null; then
        # Extract the checksum for our specific file
        EXPECTED_CHECKSUM=$(grep "${BINARY_NAME}" checksums.txt | awk '{print $1}')
        ACTUAL_CHECKSUM=$(shasum -a 256 "${ARCHIVE_NAME}" | awk '{print $1}')
    elif command -v sha256sum &> /dev/null; then
        EXPECTED_CHECKSUM=$(grep "${BINARY_NAME}" checksums.txt | awk '{print $1}')
        ACTUAL_CHECKSUM=$(sha256sum "${ARCHIVE_NAME}" | awk '{print $1}')
    else
        echo -e "${YELLOW}Warning: No SHA256 tool found, skipping checksum verification${NC}"
        EXPECTED_CHECKSUM=""
    fi

    if [ -n "$EXPECTED_CHECKSUM" ]; then
        if [ "$EXPECTED_CHECKSUM" = "$ACTUAL_CHECKSUM" ]; then
            echo -e "${GREEN}  ✓ Checksum verified${NC}"
        else
            echo -e "${RED}Error: Checksum verification failed${NC}"
            echo -e "Expected: ${EXPECTED_CHECKSUM}"
            echo -e "Actual:   ${ACTUAL_CHECKSUM}"
            exit 1
        fi
    fi
    cd - > /dev/null
else
    echo -e "${YELLOW}Warning: Could not download checksums, skipping verification${NC}"
fi

# Extract archive
echo "Extracting archive..."
tar -xzf "${TMP_DIR}/${ARCHIVE_NAME}" -C "$TMP_DIR"
echo -e "${GREEN}  ✓ Extracted${NC}"

# Install binary
echo "Installing to ${INSTALL_DIR}..."
mv "${TMP_DIR}/${BINARY_NAME}" "${INSTALL_DIR}/cntm"
chmod +x "${INSTALL_DIR}/cntm"
echo -e "${GREEN}  ✓ Installed${NC}"

# Verify installation
if command -v cntm &> /dev/null; then
    INSTALLED_VERSION=$(cntm version --output json 2>/dev/null | grep -o '"version":"[^"]*"' | cut -d'"' -f4 || echo "unknown")
    echo ""
    echo -e "${GREEN}================================================${NC}"
    echo -e "${GREEN}  cntm installed successfully!${NC}"
    echo -e "${GREEN}================================================${NC}"
    echo -e "Version: ${INSTALLED_VERSION}"
    echo -e "Location: ${INSTALL_DIR}/cntm"
    echo ""
    echo "Get started:"
    echo "  cntm init              # Initialize your project"
    echo "  cntm search <query>    # Search for tools"
    echo "  cntm install <name>    # Install a tool"
    echo ""
elif [ -f "${INSTALL_DIR}/cntm" ]; then
    echo ""
    echo -e "${GREEN}================================================${NC}"
    echo -e "${GREEN}  cntm installed successfully!${NC}"
    echo -e "${GREEN}================================================${NC}"
    echo -e "Location: ${INSTALL_DIR}/cntm"
    echo ""
    echo -e "${YELLOW}Note: ${INSTALL_DIR} is not in your PATH${NC}"
    echo "Add the following to your shell profile (~/.bashrc, ~/.zshrc, etc.):"
    echo ""
    echo "  export PATH=\"\$PATH:${INSTALL_DIR}\""
    echo ""
    echo "Then restart your shell or run: source ~/.bashrc"
    echo ""
else
    echo -e "${RED}Error: Installation failed${NC}"
    exit 1
fi
