#!/bin/bash
#
# CHANGE THIS FILE TO PUSH YOUR GO CODE TO A GO REPOSITORY
API_SERVER_SOURCE_PATH=../../generated-code
PROTO_MESSAGES_PATH=$API_SERVER_SOURCE_PATH/protobuf/swagger
 
mkdir -p $PROTO_MESSAGES_PATH
cp -r target/generated-sources/swagger/* $PROTO_MESSAGES_PATH/

echo "Start"
cd $PROTO_MESSAGES_PATH

# Check if protoc-gen-go is available
if command -v protoc-gen-go &> /dev/null; then
    # Generating go code from proto files
    find . -name "*.proto" -exec sh -c 'protoc --go_out=. --go-grpc_out=. "$0"' {} \;
    if [ $? -ne 0 ]; then
        echo "⚠️  Warning: Failed to compile proto files to Go (protoc-gen-go may not be installed)"
        echo "   Proto files copied successfully, but Go compilation skipped"
    fi
    
    if find "." -type f -name "*pb.go" | grep -q .; then
        echo "✅ Go files generated from protobuf messages"
    else
        echo "⚠️  Warning: No Go files found (may be generated later)"
    fi
else
    echo "⚠️  Warning: protoc-gen-go not found in PATH"
    echo "   Proto files copied successfully, but Go compilation skipped"
    echo "   Go files will be generated when protoc-gen-go is available"
fi

# Correcting import statements of go files generated from proto files
export old_path_common="common/"
export new_path_common="github.com/nutanix/ntnx-api-golang-mock-pc/generated-code/protobuf/common/"
export folder_path=$(pwd)

echo $folder_path
for file in $(find $folder_path -type f)
do
    if [[ -f "$file" && "$file" =~ \.go$ ]]; then
        echo " file name is $file"
        sed -i  "s#$old_path_common#$new_path_common#g" "$file"
        sed -i '' "s#$old_path_common#$new_path_common#g" "$file"
    fi
done
echo "Done"

