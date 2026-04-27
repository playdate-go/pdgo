#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

echo "Building all Tour of Go examples"
echo "=============================="
echo ""

FAILED=()
SUCCEEDED=()

for dir in "$SCRIPT_DIR"/[0-9]*/ "$SCRIPT_DIR"/generics_*/; do
    if [ -f "$dir/build.sh" ]; then
        name=$(basename "$dir")
        echo "📦 Building $name..."
        echo "------------------------"

        if (cd "$dir" && ./build.sh); then
            SUCCEEDED+=("$name")
            echo "✅ $name built successfully"
        else
            FAILED+=("$name")
            echo "❌ $name failed to build"
        fi
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
