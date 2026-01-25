#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

echo "Building all pdgo examples"
echo "=============================="
echo ""

EXAMPLES=(
    "3d_library"
    "spritegame"
    "life"
    "bouncing_square"
    "go_logo"
    "hello_world"
)

FAILED=()
SUCCEEDED=()

for example in "${EXAMPLES[@]}"; do
    echo "📦 Building $example..."
    echo "------------------------"
    
    if [ -d "$SCRIPT_DIR/$example" ] && [ -f "$SCRIPT_DIR/$example/build.sh" ]; then
        cd "$SCRIPT_DIR/$example"
        if ./build.sh; then
            SUCCEEDED+=("$example")
            echo "✅ $example built successfully"
        else
            FAILED+=("$example")
            echo "❌ $example failed to build"
        fi
        echo ""
    else
        echo "⚠️  Skipping $example (no build.sh found)"
        echo ""
    fi
done

echo "=============================="
echo "Build Summary"
echo "=============================="
echo "✅ Succeeded: ${#SUCCEEDED[@]}"
for ex in "${SUCCEEDED[@]}"; do
    echo "   - $ex"
done

if [ ${#FAILED[@]} -gt 0 ]; then
    echo "❌ Failed: ${#FAILED[@]}"
    for ex in "${FAILED[@]}"; do
        echo "   - $ex"
    done
    exit 1
fi

echo ""
echo "🎉 All examples built successfully!"
