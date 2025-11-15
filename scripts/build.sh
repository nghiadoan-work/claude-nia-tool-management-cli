#!/bin/bash

# build.sh - Multi-platform build script for cntm
# Builds binaries for macOS, Linux, and Windows

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Version information
VERSION="1.0.0"
PACKAGE="github.com/nghiadt/claude-nia-tool-management-cli"
VERSION_PACKAGE="${PACKAGE}/pkg/version"

# Build information
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
GO_VERSION=$(go version | awk '{print $3}')

# Output directory
DIST_DIR="dist"
CHECKSUM_FILE="${DIST_DIR}/checksums.txt"

# Platforms to build
PLATFORMS=(
    "darwin/amd64"
    "darwin/arm64"
    "linux/amd64"
    "linux/arm64"
    "windows/amd64"
)

echo -e "${BLUE}================================================${NC}"
echo -e "${BLUE}  Building cntm v${VERSION}${NC}"
echo -e "${BLUE}================================================${NC}"
echo -e "Git commit: ${GIT_COMMIT}"
echo -e "Build date: ${BUILD_DATE}"
echo -e "Go version: ${GO_VERSION}"
echo ""

# Clean previous builds
echo -e "${YELLOW}Cleaning previous builds...${NC}"
rm -rf "${DIST_DIR}"
mkdir -p "${DIST_DIR}"

# Build flags
LDFLAGS="-X ${VERSION_PACKAGE}.Version=${VERSION}"
LDFLAGS="${LDFLAGS} -X ${VERSION_PACKAGE}.GitCommit=${GIT_COMMIT}"
LDFLAGS="${LDFLAGS} -X ${VERSION_PACKAGE}.BuildDate=${BUILD_DATE}"
LDFLAGS="${LDFLAGS} -X ${VERSION_PACKAGE}.GoVersion=${GO_VERSION}"
# Strip debug info to reduce binary size
LDFLAGS="${LDFLAGS} -s -w"

# Build for each platform
for PLATFORM in "${PLATFORMS[@]}"; do
    # Split platform into OS and ARCH
    OS=$(echo "$PLATFORM" | cut -d'/' -f1)
    ARCH=$(echo "$PLATFORM" | cut -d'/' -f2)

    # Output binary name
    OUTPUT_NAME="cntm-${OS}-${ARCH}"
    if [ "$OS" = "windows" ]; then
        OUTPUT_NAME="${OUTPUT_NAME}.exe"
    fi

    OUTPUT_PATH="${DIST_DIR}/${OUTPUT_NAME}"

    echo -e "${BLUE}Building for ${OS}/${ARCH}...${NC}"

    # Build
    GOOS=$OS GOARCH=$ARCH CGO_ENABLED=0 go build \
        -ldflags "${LDFLAGS}" \
        -o "${OUTPUT_PATH}" \
        .

    # Check if build succeeded
    if [ $? -eq 0 ]; then
        # Get file size
        FILE_SIZE=$(ls -lh "${OUTPUT_PATH}" | awk '{print $5}')
        echo -e "${GREEN}  ✓ Built ${OUTPUT_NAME} (${FILE_SIZE})${NC}"
    else
        echo -e "${RED}  ✗ Failed to build ${OUTPUT_NAME}${NC}"
        exit 1
    fi
done

echo ""
echo -e "${YELLOW}Creating checksums...${NC}"

# Create checksums file
cd "${DIST_DIR}"
if command -v shasum &> /dev/null; then
    shasum -a 256 cntm-* > checksums.txt
elif command -v sha256sum &> /dev/null; then
    sha256sum cntm-* > checksums.txt
else
    echo -e "${RED}Error: Neither shasum nor sha256sum found${NC}"
    exit 1
fi
cd ..

echo -e "${GREEN}  ✓ Created ${CHECKSUM_FILE}${NC}"

# Create archives
echo ""
echo -e "${YELLOW}Creating archives...${NC}"

cd "${DIST_DIR}"
for PLATFORM in "${PLATFORMS[@]}"; do
    OS=$(echo "$PLATFORM" | cut -d'/' -f1)
    ARCH=$(echo "$PLATFORM" | cut -d'/' -f2)

    BINARY_NAME="cntm-${OS}-${ARCH}"
    if [ "$OS" = "windows" ]; then
        BINARY_NAME="${BINARY_NAME}.exe"
        ARCHIVE_NAME="cntm-${VERSION}-${OS}-${ARCH}.zip"

        # Create ZIP for Windows
        zip -q "${ARCHIVE_NAME}" "${BINARY_NAME}"
        echo -e "${GREEN}  ✓ Created ${ARCHIVE_NAME}${NC}"
    else
        ARCHIVE_NAME="cntm-${VERSION}-${OS}-${ARCH}.tar.gz"

        # Create tar.gz for Unix
        tar -czf "${ARCHIVE_NAME}" "${BINARY_NAME}"
        echo -e "${GREEN}  ✓ Created ${ARCHIVE_NAME}${NC}"
    fi
done
cd ..

# Summary
echo ""
echo -e "${BLUE}================================================${NC}"
echo -e "${GREEN}Build completed successfully!${NC}"
echo -e "${BLUE}================================================${NC}"
echo -e "Output directory: ${DIST_DIR}/"
echo ""
echo "Binaries:"
ls -lh "${DIST_DIR}"/cntm-* | grep -v ".tar.gz" | grep -v ".zip" | awk '{printf "  %s  %s\n", $5, $9}'
echo ""
echo "Archives:"
ls -lh "${DIST_DIR}"/*.tar.gz "${DIST_DIR}"/*.zip 2>/dev/null | awk '{printf "  %s  %s\n", $5, $9}' || true
echo ""
echo -e "${YELLOW}Next steps:${NC}"
echo "  1. Test binaries: ./dist/cntm-<platform> version"
echo "  2. Create git tag: git tag -a v${VERSION} -m 'Release v${VERSION}'"
echo "  3. Push tag: git push origin v${VERSION}"
echo "  4. Create GitHub release with archives"
echo ""
