#!/bin/bash
# Build script for Percepta release binaries
#
# NOTE: Due to CGO dependencies (webcam library), cross-compilation requires
# platform-specific toolchains. For alpha, build natively on each platform:
# - Linux: Run this script on Linux (builds Linux binaries)
# - macOS: Run this script on macOS (builds macOS binaries)
# - Windows: Run on Windows with Git Bash or WSL
#
# Future: Replace blackjack/webcam with pure-Go camera library for true cross-compilation

set -e

VERSION=${1:-dev}
echo "Building Percepta release: $VERSION"
echo "Platform: $(uname -s) / $(uname -m)"
echo ""

# Create dist directory
mkdir -p dist

# Detect current platform and build for it
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Normalize arch names
case "$ARCH" in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        echo "Warning: Unknown architecture $ARCH, defaulting to amd64"
        ARCH="amd64"
        ;;
esac

# Build for current platform
echo "Building for current platform: $OS-$ARCH"
if [ "$OS" = "darwin" ]; then
    go build -ldflags "-s -w" -o dist/percepta-darwin-${ARCH} ./cmd/percepta
    tar -czf dist/percepta-${VERSION}-darwin-${ARCH}.tar.gz -C dist percepta-darwin-${ARCH}
    echo "✓ Created: dist/percepta-${VERSION}-darwin-${ARCH}.tar.gz"
elif [ "$OS" = "linux" ]; then
    go build -ldflags "-s -w" -o dist/percepta-linux-${ARCH} ./cmd/percepta
    tar -czf dist/percepta-${VERSION}-linux-${ARCH}.tar.gz -C dist percepta-linux-${ARCH}
    echo "✓ Created: dist/percepta-${VERSION}-linux-${ARCH}.tar.gz"
elif [[ "$OS" =~ "mingw" ]] || [[ "$OS" =~ "msys" ]] || [[ "$OS" =~ "cygwin" ]]; then
    # Windows (Git Bash, MSYS2, Cygwin)
    go build -ldflags "-s -w" -o dist/percepta-windows-amd64.exe ./cmd/percepta
    if command -v zip &> /dev/null; then
        (cd dist && zip percepta-${VERSION}-windows-amd64.zip percepta-windows-amd64.exe)
        echo "✓ Created: dist/percepta-${VERSION}-windows-amd64.zip"
    else
        tar -czf dist/percepta-${VERSION}-windows-amd64.tar.gz -C dist percepta-windows-amd64.exe
        echo "✓ Created: dist/percepta-${VERSION}-windows-amd64.tar.gz"
    fi
else
    echo "Error: Unknown platform $OS"
    exit 1
fi


echo ""
echo "Build complete for version: $VERSION"
echo ""
echo "Release artifacts in dist/:"
ls -lh dist/*.tar.gz dist/*.zip 2>/dev/null || ls -lh dist/*.tar.gz
echo ""
echo "Next steps:"
echo "1. Test binary: ./dist/percepta-* --version"
echo "2. Build on other platforms (macOS, Windows) and collect binaries"
echo "3. Create GitHub release:"
echo "   gh release create v${VERSION} dist/percepta-${VERSION}-*.tar.gz dist/percepta-${VERSION}-*.zip"
echo "4. Or manually upload to: https://github.com/perceptumx/percepta/releases/new"
echo ""
echo "Note: Cross-platform builds require native toolchain on each platform."
echo "See docs/installation.md for platform-specific build instructions."
