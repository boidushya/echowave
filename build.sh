#!/bin/bash

# EchoWave Build Script
# Part of the better-lyrics ecosystem

set -e

echo "🎵 Building EchoWave..."
echo "=================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Clean previous builds
echo -e "${BLUE}🧹 Cleaning previous builds...${NC}"
rm -rf dist/
mkdir -p dist/

# Build info
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
# Clean version by removing 'v' prefix if present
VERSION_CLEAN=$(echo $VERSION | sed 's/^v//')
BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ)
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

echo -e "${BLUE}📦 Building version: ${YELLOW}$VERSION_CLEAN${NC}"
echo -e "${BLUE}🕐 Build time: ${YELLOW}$BUILD_TIME${NC}"
echo -e "${BLUE}📝 Commit: ${YELLOW}$COMMIT${NC}"

# Build for current platform
echo -e "${BLUE}🔨 Building for current platform...${NC}"
go build -ldflags="-s -X main.VERSION=$VERSION_CLEAN" -o dist/echowave .

# Cross-compilation builds
echo -e "${BLUE}🌍 Cross-compiling for multiple platforms...${NC}"

# macOS (Intel)
echo -e "${BLUE}  📱 Building for macOS (Intel)...${NC}"
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -X main.VERSION=$VERSION_CLEAN" -o dist/echowave-macos-intel .

# macOS (Apple Silicon)
echo -e "${BLUE}  📱 Building for macOS (Apple Silicon)...${NC}"
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -X main.VERSION=$VERSION_CLEAN" -o dist/echowave-macos-arm64 .

# Linux (AMD64)
echo -e "${BLUE}  🐧 Building for Linux (AMD64)...${NC}"
GOOS=linux GOARCH=amd64 go build -ldflags="-s -X main.VERSION=$VERSION_CLEAN" -o dist/echowave-linux-amd64 .

# Linux (ARM64)
echo -e "${BLUE}  🐧 Building for Linux (ARM64)...${NC}"
GOOS=linux GOARCH=arm64 go build -ldflags="-s -X main.VERSION=$VERSION_CLEAN" -o dist/echowave-linux-arm64 .

# Windows (AMD64)
echo -e "${BLUE}  🪟 Building for Windows (AMD64)...${NC}"
GOOS=windows GOARCH=amd64 go build -ldflags="-s -X main.VERSION=$VERSION_CLEAN" -o dist/echowave-windows-amd64.exe .

# Windows (ARM64)
echo -e "${BLUE}  🪟 Building for Windows (ARM64)...${NC}"
GOOS=windows GOARCH=arm64 go build -ldflags="-s -X main.VERSION=$VERSION_CLEAN" -o dist/echowave-windows-arm64.exe .

# Create checksums
echo -e "${BLUE}🔐 Creating checksums...${NC}"
cd dist/
sha256sum * > checksums.txt
cd ..

# Show build results
echo -e "${GREEN}✅ Build completed successfully!${NC}"
echo -e "${GREEN}📂 Files created:${NC}"
ls -la dist/

# File sizes
echo -e "${GREEN}📊 File sizes:${NC}"
du -h dist/* | sort -hr

echo ""
echo -e "${GREEN}🎉 EchoWave build complete!${NC}"
echo -e "${BLUE}🌐 Part of the better-lyrics ecosystem${NC}"
echo -e "${BLUE}💻 Visit: https://better-lyrics.boidu.dev${NC}"