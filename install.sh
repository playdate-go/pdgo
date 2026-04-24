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
#   SKIP_TINYGO=1  - Skip TinyGo build (simulator-only)
#   JOBS=N         - Number of parallel jobs for LLVM build
#

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

TINYGO_DIR="$HOME/tinygo-playdate"
TINYGO_VERSION="0.40.1"
LLVM_VERSION="20"
JOBS=${JOBS:-$(sysctl -n hw.ncpu 2>/dev/null || nproc 2>/dev/null || echo 4)}

echo -e "${CYAN}╔══════════════════════════════════════════════════════════╗${NC}"
echo -e "${CYAN}║              PdGo (Golang for Playdate) - Full Installer ║${NC}"
echo -e "${CYAN}╚══════════════════════════════════════════════════════════╝${NC}"
echo ""

# Check OS
if [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "win32" ]]; then
    echo -e "${RED}ERROR: Windows is not supported${NC}"
    echo "Supported: macOS, Linux"
    exit 1
fi

# ============================================================================
# Step 1: Check dependencies
# ============================================================================
echo -e "${YELLOW}[1/4] Checking dependencies...${NC}"

missing=""
command -v go &>/dev/null || missing="$missing go"
command -v git &>/dev/null || missing="$missing git"
command -v cmake &>/dev/null || missing="$missing cmake"
command -v ninja &>/dev/null || missing="$missing ninja"

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
    echo -e "${RED}ERROR: PLAYDATE_SDK_PATH is not set${NC}"
    echo ""
    echo "1. Download Playdate SDK from https://play.date/dev/"
    echo "2. Add to your shell profile (~/.zshrc or ~/.bashrc):"
    echo "   export PLAYDATE_SDK_PATH=\"/path/to/PlaydateSDK\""
    exit 1
fi

if [ ! -d "$PLAYDATE_SDK_PATH" ]; then
    echo -e "${RED}ERROR: PLAYDATE_SDK_PATH does not exist: $PLAYDATE_SDK_PATH${NC}"
    exit 1
fi

SDK_VERSION=$(cat "$PLAYDATE_SDK_PATH/VERSION.txt" 2>/dev/null | tr -d '\n' || echo "unknown")
echo -e "  Playdate SDK: ${GREEN}$SDK_VERSION${NC}"
echo -e "  cmake: ${GREEN}$(cmake --version | head -1 | awk '{print $3}')${NC}"
echo -e "  ninja: ${GREEN}$(ninja --version)${NC}"

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
# Step 3: Build TinyGo with Playdate support
# ============================================================================
if [ "$SKIP_TINYGO" = "1" ]; then
    echo ""
    echo -e "${YELLOW}[3/4] Skipping TinyGo build (SKIP_TINYGO=1)${NC}"
    echo "  You can build it later with:"
    echo "  git clone https://github.com/playdate-go/pdgo.git && cd pdgo"
    echo "  ./cmd/pdgoc/scripts/build-tinygo-playdate.sh"
else
    echo ""
    echo -e "${YELLOW}[3/4] Building TinyGo with Playdate support...${NC}"
    echo -e "  TinyGo version: ${GREEN}v$TINYGO_VERSION${NC}"
    echo -e "  Install path: ${GREEN}$TINYGO_DIR${NC}"
    echo -e "  Parallel jobs: ${GREEN}$JOBS${NC}"
    
    if [[ "$OSTYPE" == "darwin"* ]]; then
        echo -e "  LLVM: ${GREEN}Building from source (~25 minutes)${NC}"
    else
        echo -e "  LLVM: ${GREEN}Will use system LLVM if available${NC}"
    fi
    echo ""

    # Clone or update TinyGo
    if [ -d "$TINYGO_DIR/.git" ]; then
        echo "  Updating TinyGo repository..."
        cd "$TINYGO_DIR"
        git fetch --all --tags
    else
        echo "  Cloning TinyGo repository..."
        git clone --recursive https://github.com/tinygo-org/tinygo.git "$TINYGO_DIR"
        cd "$TINYGO_DIR"
    fi

    # Checkout specific version
    echo "  Checking out v$TINYGO_VERSION..."
    git checkout "v$TINYGO_VERSION"
    git submodule update --init --recursive

    # Clean previous build to ensure Playdate files are included
    if [ -d "$TINYGO_DIR/build" ]; then
        echo "  Cleaning previous build..."
        make clean 2>/dev/null || rm -rf "$TINYGO_DIR/build"
    fi

    # Add Playdate support files BEFORE building (required for gc.playdate)
    echo "  Adding Playdate support files..."

    # Define patch files by destination directory
    # Format: "source_file:destination_file"
    PATCHES_TARGETS=(
        "playdate.json:targets/playdate.json"
        "playdate.ld:targets/playdate.ld"
    )

    PATCHES_RUNTIME=(
        "runtime_playdate.go:src/runtime/runtime_playdate.go"
        "gc_playdate.go:src/runtime/gc_playdate.go"
    )

    if [ "$USE_LOCAL_PATCHES" = true ]; then
        # Use local patch files from the repository
        echo "  Using local Playdate patches from repository..."

        LOCAL_PATCHES_DIR="${LOCAL_REPO_ROOT}/cmd/pdgoc/tinygo-patches"

        # Copy all patch files
        for patch in "${PATCHES_TARGETS[@]}" "${PATCHES_RUNTIME[@]}"; do
            IFS=':' read -r src dst <<< "$patch"
            if ! cp "${LOCAL_PATCHES_DIR}/${src}" "$TINYGO_DIR/${dst}"; then
                echo -e "${RED}Failed to copy ${src}${NC}"
                exit 1
            fi
        done

        echo -e "  ${GREEN}Playdate patches applied from local files${NC}"
    else
        # Download patch files from GitHub
        echo "  Downloading Playdate patches from GitHub..."

        # Determine the branch/tag to use for downloading patches
        PDGO_BRANCH=${PDGO_BRANCH:-"main"}
        PATCHES_BASE_URL="https://raw.githubusercontent.com/playdate-go/pdgo/${PDGO_BRANCH}/cmd/pdgoc/tinygo-patches"

        # Download all patch files
        for patch in "${PATCHES_TARGETS[@]}" "${PATCHES_RUNTIME[@]}"; do
            IFS=':' read -r src dst <<< "$patch"
            if ! curl -sL "${PATCHES_BASE_URL}/${src}" -o "$TINYGO_DIR/${dst}"; then
                echo -e "${RED}Failed to download ${src}${NC}"
                exit 1
            fi
        done

        echo -e "  ${GREEN}Playdate patches applied successfully${NC}"
    fi

    # Check for system LLVM (Linux only)
    USE_SYSTEM_LLVM=false
    SYSTEM_LLVM_PATH=""

    if [[ "$OSTYPE" != "darwin"* ]]; then
        # Try to find system LLVM on Linux
        for ver in 20 19 18 17; do
            if [ -f "/usr/lib/llvm-$ver/bin/llvm-config" ]; then
                SYSTEM_LLVM_PATH="/usr/lib/llvm-$ver"
                USE_SYSTEM_LLVM=true
                break
            fi
        done
    fi

    # Build TinyGo
    if [ "$USE_SYSTEM_LLVM" = true ]; then
        echo -e "  Found system LLVM at: ${GREEN}$SYSTEM_LLVM_PATH${NC}"
        echo "  Building TinyGo (~1 minute)..."
        make LLVM_BUILDDIR="$SYSTEM_LLVM_PATH"
    else
        echo "  Downloading LLVM sources..."
        make llvm-source

        echo "  Building LLVM (~20-25 minutes)..."
        echo "  Started at: $(date)"
        
        # Use clang if available
        if command -v clang &>/dev/null; then
            export CC=clang
            export CXX=clang++
        fi

        LLVM_PARALLEL=$JOBS make llvm-build
        echo "  LLVM completed at: $(date)"

        echo "  Building TinyGo..."
        make
    fi

    echo -e "  ${GREEN}TinyGo built successfully!${NC}"
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

if [ "$SKIP_TINYGO" != "1" ] && ! echo "$PATH" | grep -q "$TINYGO_DIR/build"; then
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
        echo "  export PATH=\"\$PATH:$TINYGO_DIR/build\""
    fi
    echo ""
    
    read -p "Add automatically? [Y/n] " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Nn]$ ]]; then
        echo "" >> "$SHELL_RC"
        echo "# pdgo - Playdate Go toolkit" >> "$SHELL_RC"
        if [ "$NEED_GOBIN" = true ]; then
            echo "export PATH=\"\$PATH:$GOBIN\"" >> "$SHELL_RC"
        fi
        if [ "$NEED_TINYGO" = true ]; then
            echo "export PATH=\"\$PATH:$TINYGO_DIR/build\"" >> "$SHELL_RC"
        fi
        echo -e "${GREEN}Added to $SHELL_RC${NC}"
        echo "Run: source $SHELL_RC"
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
echo -e "  TinyGo: ${GREEN}$TINYGO_DIR/build/tinygo${NC}"
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
