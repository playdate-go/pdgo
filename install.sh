#!/bin/bash
#
# pdgo Full Installer
#
# Installs everything needed to build Go games for Playdate:
#   - pdgoc (build tool)
#   - TinyGo with Playdate support (for device builds)
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/playdate-go/pdgo/main/install.sh | bash
#
# Options (via environment variables):
#   SKIP_TINYGO=1  - Skip TinyGo install (simulator-only)
#

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

TINYGO_DIR="$HOME/tinygo-playdate"
TINYGO_VERSION="0.40.1"

echo -e "${CYAN}╔══════════════════════════════════════════════════════════╗${NC}"
echo -e "${CYAN}║              PdGo (Golang for Playdate) - Full Installer ║${NC}"
echo -e "${CYAN}╚══════════════════════════════════════════════════════════╝${NC}"
echo ""

# Check OS
if [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "win32" ]]; then
    echo -e "${RED}ERROR: Windows is not supported. Use install.ps1 instead.${NC}"
    exit 1
fi

# Detect OS/arch for TinyGo download
if [[ "$OSTYPE" == "darwin"* ]]; then
    if [[ "$(uname -m)" == "arm64" ]]; then
        TINYGO_OS="darwin-arm64"
    else
        TINYGO_OS="darwin-amd64"
    fi
else
    TINYGO_OS="linux-amd64"
    if [[ "$(uname -m)" == "aarch64" ]]; then
        TINYGO_OS="linux-arm64"
    fi
fi

# ============================================================================
# Step 1: Check dependencies
# ============================================================================
echo -e "${YELLOW}[1/4] Checking dependencies...${NC}"

missing=""
command -v go &>/dev/null || missing="$missing go"
command -v git &>/dev/null || missing="$missing git"

if [ -n "$missing" ]; then
    echo -e "${RED}Missing dependencies:$missing${NC}"
    echo ""
    if [[ "$OSTYPE" == "darwin"* ]]; then
        echo "Install with: brew install$missing"
    elif command -v apt &>/dev/null; then
        echo "Install with: sudo apt install$missing"
    elif command -v dnf &>/dev/null; then
        echo "Install with: sudo dnf install$missing"
    elif command -v pacman &>/dev/null; then
        echo "Install with: sudo pacman -S$missing"
    fi
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo -e "  Go: ${GREEN}$GO_VERSION${NC}"

# Check Playdate SDK
if [ -z "$PLAYDATE_SDK_PATH" ]; then
    # Try default paths
    if [ -d "$HOME/Developer/PlaydateSDK" ]; then
        PLAYDATE_SDK_PATH="$HOME/Developer/PlaydateSDK"
    fi
fi

if [ -z "$PLAYDATE_SDK_PATH" ] || [ ! -d "$PLAYDATE_SDK_PATH" ]; then
    echo -e "${YELLOW}Playdate SDK not found.${NC}"
    DO_INSTALL=0
    if [ "$CI" = "1" ]; then
        DO_INSTALL=1
        echo -e "${YELLOW}CI detected — installing Playdate SDK automatically${NC}"
    else
        read -p "Download and install Playdate SDK automatically? [Y/n] " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Nn]$ ]]; then
            DO_INSTALL=1
        fi
    fi

    if [ "$DO_INSTALL" = "1" ]; then
        TMPDIR="${TMPDIR:-/tmp}"
        SDK_TMP="$TMPDIR/playdate-sdk-install"

        if [[ "$OSTYPE" == "darwin"* ]]; then
            echo -e "${CYAN}Downloading Playdate SDK for macOS...${NC}"
            mkdir -p "$SDK_TMP"
            curl -sL "https://download-cdn.panic.com/playdate_sdk/PlaydateSDK-latest.zip" -o "$SDK_TMP/PlaydateSDK.zip"
            unzip -q "$SDK_TMP/PlaydateSDK.zip" -d "$SDK_TMP"
            sudo installer -pkg "$SDK_TMP/PlaydateSDK.pkg" -target /
            rm -rf "$SDK_TMP"
            PLAYDATE_SDK_PATH="$HOME/Developer/PlaydateSDK"
        else
            echo -e "${CYAN}Downloading Playdate SDK for Linux...${NC}"
            mkdir -p "$SDK_TMP"
            curl -sL "https://download.panic.com/playdate_sdk/Linux/PlaydateSDK-latest.tar.gz" | tar xz -C "$SDK_TMP"
            SDK_EXTRACTED=$(ls -d "$SDK_TMP"/PlaydateSDK-* | head -1)
            PLAYDATE_SDK_PATH="$HOME/Developer/PlaydateSDK"
            mkdir -p "$HOME/Developer"
            mv "$SDK_EXTRACTED" "$PLAYDATE_SDK_PATH"
            echo -e "${YELLOW}Add this to your shell profile to persist across sessions:${NC}"
            echo -e "  export PLAYDATE_SDK_PATH=\"$PLAYDATE_SDK_PATH\""
        fi

        if [ ! -d "$PLAYDATE_SDK_PATH" ]; then
            echo -e "${RED}ERROR: Playdate SDK installation failed${NC}"
            exit 1
        fi

        export PLAYDATE_SDK_PATH
    else
        echo ""
        echo "1. Download Playdate SDK from https://play.date/dev/"
        echo "2. Add to your shell profile (~/.zshrc or ~/.bashrc):"
        echo "   export PLAYDATE_SDK_PATH=\"/path/to/PlaydateSDK\""
        exit 1
    fi
fi

SDK_VERSION=$(cat "$PLAYDATE_SDK_PATH/VERSION.txt" 2>/dev/null | tr -d '\n' || echo "unknown")
echo -e "  Playdate SDK: ${GREEN}$SDK_VERSION${NC}"

echo -e "${GREEN}All dependencies OK${NC}"

# ============================================================================
# Step 2: Install pdgoc
# ============================================================================
echo ""
echo -e "${YELLOW}[2/4] Installing pdgoc...${NC}"

# Detect if running from pdgo repo root
if [ -d "cmd/pdgoc" ] && [ -f "go.mod" ]; then
    echo "  Running from pdgo repository root - using local source"

    # Get version info from local git
    VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
    COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
    DATE=$(date -u '+%Y-%m-%d %H:%M:%S UTC')

    echo "  Building pdgoc from local source..."
    echo "    Version: $VERSION"
    echo "    Commit:  $COMMIT"
    echo "    Date:    $DATE"

    cd cmd/pdgoc

    echo "  Downloading Go dependencies..."
    go mod tidy

    GOBIN="${GOBIN:-$(go env GOPATH)/bin}"
    mkdir -p "$GOBIN"

    # Build with ldflags to inject version information
    go build -ldflags="-X 'main.Version=$VERSION' -X 'main.Commit=$COMMIT' -X 'main.Date=$DATE'" -o "$GOBIN/pdgoc" .

    echo -e "  pdgoc installed at: ${GREEN}$GOBIN/pdgoc${NC}"

    # Set flag for TinyGo patches to use local files
    USE_LOCAL_PATCHES=true
    LOCAL_REPO_ROOT=$(pwd)/../..
else
    # Running from curl - download from GitHub
    echo "  Running from curl - downloading from GitHub"

    # Create a temporary directory for building pdgoc
    BUILD_DIR=$(mktemp -d)
    trap "rm -rf $BUILD_DIR" EXIT

    echo "  Downloading pdgo source..."

    # Try to get version info from GitHub API
    GITHUB_API="https://api.github.com/repos/playdate-go/pdgo"
    LATEST_TAG=$(curl -s "${GITHUB_API}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"tag_name": ?"([^"]+)".*/\1/' 2>/dev/null || echo "")
    LATEST_COMMIT=$(curl -s "${GITHUB_API}/commits/main" | grep '"sha"' | head -1 | sed -E 's/.*"sha": ?"([^"]+)".*/\1/' 2>/dev/null | cut -c1-7 || echo "unknown")

    VERSION=${LATEST_TAG:-"latest"}
    COMMIT=${LATEST_COMMIT:-"unknown"}
    DATE=$(date -u '+%Y-%m-%d %H:%M:%S UTC')

    # Download and extract source
    if [ -n "$LATEST_TAG" ]; then
        echo "  Downloading release $LATEST_TAG..."
        curl -sL "https://github.com/playdate-go/pdgo/archive/refs/tags/${LATEST_TAG}.tar.gz" | tar -xz -C "$BUILD_DIR"
        SOURCE_DIR="$BUILD_DIR/pdgo-${LATEST_TAG#v}"
    else
        echo "  Downloading latest source..."
        curl -sL "https://github.com/playdate-go/pdgo/archive/refs/heads/main.tar.gz" | tar -xz -C "$BUILD_DIR"
        SOURCE_DIR="$BUILD_DIR/pdgo-main"
    fi

    cd "$SOURCE_DIR/cmd/pdgoc"

    echo "  Downloading Go dependencies..."
    go mod tidy

    echo "  Building pdgoc..."
    echo "    Version: $VERSION"
    echo "    Commit:  $COMMIT"
    echo "    Date:    $DATE"

    GOBIN="${GOBIN:-$(go env GOPATH)/bin}"
    mkdir -p "$GOBIN"

    # Build with ldflags to inject version information
    go build -ldflags="-X 'main.Version=$VERSION' -X 'main.Commit=$COMMIT' -X 'main.Date=$DATE'" -o "$GOBIN/pdgoc" .

    echo -e "  pdgoc installed at: ${GREEN}$GOBIN/pdgoc${NC}"

    # Set flag for TinyGo patches to download from GitHub
    USE_LOCAL_PATCHES=false
fi


# Check if GOBIN is in PATH
if ! command -v pdgoc &>/dev/null; then
    echo -e "${YELLOW}  Note: Add GOBIN to PATH: export PATH=\"\$PATH:$GOBIN\"${NC}"
fi

# ============================================================================
# Step 3: Install TinyGo with Playdate support (Pre-compiled)
# ============================================================================
if [ "$SKIP_TINYGO" = "1" ]; then
    echo ""
    echo -e "${YELLOW}[3/4] Skipping TinyGo install (SKIP_TINYGO=1)${NC}"
else
    echo ""
    echo -e "${YELLOW}[3/4] Installing TinyGo with Playdate support...${NC}"
    echo -e "  TinyGo version: ${GREEN}v$TINYGO_VERSION ($TINYGO_OS)${NC}"
    echo -e "  Install path:   ${GREEN}$TINYGO_DIR${NC}"
    echo ""

    TEMP_TGZ=$(mktemp /tmp/tinygo.XXXXXX.tar.gz)
    TINYGO_URL="https://github.com/tinygo-org/tinygo/releases/download/v${TINYGO_VERSION}/tinygo${TINYGO_VERSION}.${TINYGO_OS}.tar.gz"

    echo "  Downloading official ${TINYGO_OS} release..."
    curl -fSL -o "$TEMP_TGZ" "$TINYGO_URL"

    if [ -d "$TINYGO_DIR" ]; then
        echo "  Removing old installation..."
        rm -rf "$TINYGO_DIR"
    fi

    echo "  Extracting to target directory..."
    tar -xzf "$TEMP_TGZ" -C "$HOME"
    mv "$HOME/tinygo" "$TINYGO_DIR"
    rm -f "$TEMP_TGZ"

    echo "  Injecting Playdate support files..."

    TARGETS_DIR="$TINYGO_DIR/targets"
    RUNTIME_DIR="$TINYGO_DIR/src/runtime"

    if [ "$USE_LOCAL_PATCHES" = true ]; then
        echo "  Patching TinyGo with local files"
        LOCAL_PATCHES_DIR="${LOCAL_REPO_ROOT}/cmd/pdgoc/tinygo-patches"

        cp "${LOCAL_PATCHES_DIR}/playdate.json"   "$TARGETS_DIR/playdate.json"
        cp "${LOCAL_PATCHES_DIR}/playdate.ld"     "$TARGETS_DIR/playdate.ld"
        cp "${LOCAL_PATCHES_DIR}/runtime_playdate.go" "$RUNTIME_DIR/runtime_playdate.go"
        cp "${LOCAL_PATCHES_DIR}/gc_playdate.go"      "$RUNTIME_DIR/gc_playdate.go"
    else
        # Determine the tag/branch to use for downloading patches
        PDGO_REF=${LATEST_TAG:-"main"}
        PATCHES_BASE_URL="https://raw.githubusercontent.com/playdate-go/pdgo/${PDGO_REF}/cmd/pdgoc/tinygo-patches"

        curl -sL "${PATCHES_BASE_URL}/playdate.json"        -o "$TARGETS_DIR/playdate.json"
        curl -sL "${PATCHES_BASE_URL}/playdate.ld"          -o "$TARGETS_DIR/playdate.ld"
        curl -sL "${PATCHES_BASE_URL}/runtime_playdate.go"  -o "$RUNTIME_DIR/runtime_playdate.go"
        curl -sL "${PATCHES_BASE_URL}/gc_playdate.go"       -o "$RUNTIME_DIR/gc_playdate.go"
    fi

    echo -e "  ${GREEN}TinyGo with Playdate support ready!${NC}"
fi

# ============================================================================
# Step 4: Setup PATH
# ============================================================================
echo ""
echo -e "${YELLOW}[4/4] Configuring PATH...${NC}"

SHELL_RC=""
if [ -f "$HOME/.zshrc" ]; then
    SHELL_RC="$HOME/.zshrc"
elif [ -f "$HOME/.bashrc" ]; then
    SHELL_RC="$HOME/.bashrc"
fi

NEED_GOBIN=false
NEED_TINYGO=false

if ! echo "$PATH" | grep -q "$GOBIN"; then
    NEED_GOBIN=true
fi

if [ "$SKIP_TINYGO" != "1" ] && ! echo "$PATH" | grep -q "$TINYGO_DIR/bin"; then
    NEED_TINYGO=true
fi

if [ "$NEED_GOBIN" = true ] || [ "$NEED_TINYGO" = true ]; then
    echo ""
    echo -e "${YELLOW}Add these lines to $SHELL_RC:${NC}"
    echo ""
    if [ "$NEED_GOBIN" = true ]; then
        echo "  export PATH=\"\$PATH:$GOBIN\""
    fi
    if [ "$NEED_TINYGO" = true ]; then
        echo "  export PATH=\"\$PATH:$TINYGO_DIR/bin\""
    fi
    echo ""

    if [ "$CI" = "1" ]; then
        echo -e "${YELLOW}CI detected — skipping interactive PATH setup${NC}"
    else
        read -p "Add automatically? [Y/n] " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Nn]$ ]]; then
            echo "" >> "$SHELL_RC"
            echo "# pdgo - Playdate Go toolkit" >> "$SHELL_RC"
            if [ "$NEED_GOBIN" = true ]; then
                echo "export PATH=\"\$PATH:$GOBIN\"" >> "$SHELL_RC"
            fi
            if [ "$NEED_TINYGO" = true ]; then
                echo "export PATH=\"\$PATH:$TINYGO_DIR/bin\"" >> "$SHELL_RC"
            fi
            echo -e "${GREEN}Added to $SHELL_RC${NC}"
            echo "Run: source $SHELL_RC"
        fi
    fi
else
    echo -e "${GREEN}PATH already configured${NC}"
fi

# ============================================================================
# Done!
# ============================================================================
echo ""
echo -e "${GREEN}╔══════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║              Installation Complete!                      ║${NC}"
echo -e "${GREEN}╚══════════════════════════════════════════════════════════╝${NC}"
echo ""
echo "Installed:"
echo -e "  pdgoc:  ${GREEN}$GOBIN/pdgoc${NC}"
if [ "$SKIP_TINYGO" != "1" ]; then
echo -e "  TinyGo: ${GREEN}$TINYGO_DIR/bin/tinygo${NC}"
fi
echo ""
echo "Usage:"
echo ""
echo "  # Simulator build (run in simulator)"
echo "  pdgoc -sim -run -name MyGame -author Me -desc 'Game' \\"
echo "        -bundle-id com.me.game -version 1.0 -build-number 1"
echo ""
if [ "$SKIP_TINYGO" != "1" ]; then
echo "  # Device build (for real Playdate)"
echo "  pdgoc -device -name MyGame -author Me -desc 'Game' \\"
echo "        -bundle-id com.me.game -version 1.0 -build-number 1"
echo ""
echo "  # Device build + deploy to connected Playdate"
echo "  pdgoc -device -deploy -name MyGame -author Me -desc 'Game' \\"
echo "        -bundle-id com.me.game -version 1.0 -build-number 1"
echo ""
fi
echo "Documentation: https://github.com/playdate-go/pdgo"
echo ""
