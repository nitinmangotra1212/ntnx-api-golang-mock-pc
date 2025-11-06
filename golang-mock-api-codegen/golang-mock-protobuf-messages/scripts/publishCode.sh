#!/bin/bash
# Script to publish generated protobuf code to the generated-code directory

set -e

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$(dirname "$(dirname "$(dirname "$SCRIPT_DIR")")")"
GENERATED_SRC="$SCRIPT_DIR/../target/generated-sources"
OUTPUT_DIR="$PROJECT_ROOT/generated-code/protobuf"

echo "Publishing protobuf code..."
echo "Source: $GENERATED_SRC"
echo "Output: $OUTPUT_DIR"

# Create output directory
mkdir -p "$OUTPUT_DIR"

# Copy generated proto files to output directory
if [ -d "$GENERATED_SRC" ]; then
  echo "Copying generated protobuf files..."
  cp -r "$GENERATED_SRC"/* "$OUTPUT_DIR/" || true
  echo "✅ Protobuf code published successfully"
else
  echo "⚠️  No generated sources found at $GENERATED_SRC"
fi

# Compile proto files to Go using protoc
echo "Compiling proto files to Go..."
cd "$OUTPUT_DIR"

# Find and compile all .proto files
find . -name "*.proto" -type f | while read proto_file; do
  dir=$(dirname "$proto_file")
  echo "Compiling $proto_file..."
  protoc --go_out=. --go_opt=paths=source_relative --proto_path=. "$proto_file" || echo "Failed to compile $proto_file"
done

echo "✅ Proto compilation complete"

