#!/bin/bash

# EchoWave Installation Script
# Part of the better-lyrics ecosystem

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

# Configuration
REPO="boidu/echowave"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="echowave"

# Detect OS and architecture
OS=$(uname -s)
ARCH=$(uname -m)

case $OS in
    Darwin)
        OS_NAME="darwin"
        ;;
    Linux)
        OS_NAME="linux"
        ;;
    *)
        echo -e "${RED}❌ Unsupported operating system: $OS${NC}"
        exit 1
        ;;
esac

case $ARCH in
    x86_64)
        ARCH_NAME="amd64"
        ;;
    arm64|aarch64)
        ARCH_NAME="arm64"
        ;;
    *)
        echo -e "${RED}❌ Unsupported architecture: $ARCH${NC}"
        exit 1
        ;;
esac

echo -e "${BLUE}🎵 EchoWave Installer${NC}"
echo -e "${PURPLE}Part of the better-lyrics ecosystem${NC}"
echo "================================="

# Get latest release version
echo -e "${BLUE}🔍 Fetching latest release...${NC}"
LATEST_VERSION=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_VERSION" ]; then
    echo -e "${RED}❌ Failed to fetch latest version${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Latest version: $LATEST_VERSION${NC}"

# Construct download URL
BINARY_FILE="echowave-${OS_NAME}-${ARCH_NAME}.tar.gz"
DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST_VERSION/$BINARY_FILE"

echo -e "${BLUE}📥 Downloading $BINARY_FILE...${NC}"

# Create temporary directory
TMP_DIR=$(mktemp -d)
cd "$TMP_DIR"

# Download and extract
curl -L -o "$BINARY_FILE" "$DOWNLOAD_URL"

if [ ! -f "$BINARY_FILE" ]; then
    echo -e "${RED}❌ Download failed${NC}"
    exit 1
fi

echo -e "${BLUE}📦 Extracting archive...${NC}"
tar -xzf "$BINARY_FILE"

# Get the extracted binary name
EXTRACTED_BINARY="echowave-${OS_NAME}-${ARCH_NAME}"

if [ ! -f "$EXTRACTED_BINARY" ]; then
    echo -e "${RED}❌ Binary not found in archive${NC}"
    exit 1
fi

# Install binary
echo -e "${BLUE}🔧 Installing to $INSTALL_DIR...${NC}"

if [ -w "$INSTALL_DIR" ]; then
    cp "$EXTRACTED_BINARY" "$INSTALL_DIR/$BINARY_NAME"
    chmod +x "$INSTALL_DIR/$BINARY_NAME"
else
    echo -e "${YELLOW}⚠️  Requires sudo to install to $INSTALL_DIR${NC}"
    sudo cp "$EXTRACTED_BINARY" "$INSTALL_DIR/$BINARY_NAME"
    sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"
fi

# Cleanup
cd /
rm -rf "$TMP_DIR"

echo -e "${GREEN}✅ EchoWave installed successfully!${NC}"
echo ""
echo -e "${BLUE}🚀 Quick Start:${NC}"
echo -e "${YELLOW}  echowave https://youtube.com/watch?v=xyz${NC}"
echo -e "${YELLOW}  echowave audio.mp3${NC}"
echo ""
echo -e "${BLUE}📚 Help:${NC}"
echo -e "${YELLOW}  echowave -help${NC}"
echo ""
echo -e "${BLUE}🌐 Part of the better-lyrics ecosystem${NC}"
echo -e "${BLUE}💻 Visit: https://better-lyrics.boidu.dev${NC}"

# Check if dependencies are available
echo ""
echo -e "${BLUE}🔍 Checking dependencies...${NC}"

check_dependency() {
    if command -v "$1" &> /dev/null; then
        echo -e "${GREEN}✅ $1 found${NC}"
        return 0
    else
        echo -e "${RED}❌ $1 not found${NC}"
        return 1
    fi
}

DEPS_OK=true

if ! check_dependency "ffmpeg"; then
    DEPS_OK=false
    echo -e "${YELLOW}   Install: brew install ffmpeg (macOS) or sudo apt-get install ffmpeg (Ubuntu)${NC}"
fi

if ! check_dependency "whisper"; then
    DEPS_OK=false
    echo -e "${YELLOW}   Install: pip install openai-whisper${NC}"
fi

if ! check_dependency "yt-dlp"; then
    DEPS_OK=false
    echo -e "${YELLOW}   Install: brew install yt-dlp (macOS) or pip install yt-dlp${NC}"
fi

if [ "$DEPS_OK" = true ]; then
    echo -e "${GREEN}🎉 All dependencies are installed! You're ready to go!${NC}"
else
    echo ""
    echo -e "${YELLOW}⚠️  Some dependencies are missing. EchoWave will show installation instructions when you run it.${NC}"
fi