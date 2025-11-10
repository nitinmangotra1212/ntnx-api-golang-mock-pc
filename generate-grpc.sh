#!/bin/bash
#
# Generate gRPC protobuf code (.pb.go files)
# This creates the same .pb.go files that guru service has
#

set -e

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROTO_DIR="${SCRIPT_DIR}/generated-code/protobuf/swagger/mock/v4/config"
OUT_DIR="${SCRIPT_DIR}/generated-code/protobuf/mock/v4/config"

echo "🔧 Generating gRPC protobuf code..."
echo "================================================"

# Create output directory
mkdir -p "${OUT_DIR}"

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    echo "❌ protoc not found. Installing..."
    echo ""
    echo "Please install protoc:"
    echo "  macOS: brew install protobuf"
    echo "  Linux: apt-get install protobuf-compiler"
    exit 1
fi

# Check if protoc-gen-go is installed
if ! command -v protoc-gen-go &> /dev/null; then
    echo "📦 Installing protoc-gen-go..."
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
fi

# Check if protoc-gen-go-grpc is installed
if ! command -v protoc-gen-go-grpc &> /dev/null; then
    echo "📦 Installing protoc-gen-go-grpc..."
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
fi

echo ""
echo "✅ Prerequisites installed"
echo ""

# Generate protobuf Go code
echo "🔄 Generating .pb.go files from protobuf definitions..."
echo ""

cd "${PROTO_DIR}"

# Generate for config.proto
echo "  → config.proto"
protoc --go_out="${SCRIPT_DIR}" \
    --go_opt=module=github.com/nutanix/ntnx-api-golang-mock-pc \
    config.proto

# Generate for cat_service.proto (includes gRPC service definitions)
echo "  → cat_service.proto"
protoc --go_out="${SCRIPT_DIR}" \
    --go_opt=module=github.com/nutanix/ntnx-api-golang-mock-pc \
    --go-grpc_out="${SCRIPT_DIR}" \
    --go-grpc_opt=module=github.com/nutanix/ntnx-api-golang-mock-pc \
    cat_service.proto

echo ""
echo "================================================"
echo "✅ gRPC code generation complete!"
echo ""
echo "📁 Generated files:"
ls -lh "${OUT_DIR}"/*.pb.go 2>/dev/null || echo "  (checking...)"
echo ""
echo "📦 Generated files:"
echo "  - config.pb.go          (protobuf messages)"
echo "  - cat_service.pb.go     (service messages)"
echo "  - cat_service_grpc.pb.go (gRPC service stubs)"
echo ""
echo "🎉 Ready to implement gRPC servers!"

