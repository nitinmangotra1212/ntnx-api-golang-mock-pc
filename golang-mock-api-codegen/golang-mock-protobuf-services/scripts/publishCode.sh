#!/bin/bash
#
# Script to copy generated proto service files to the main protobuf directory

API_SERVER_SOURCE_PATH=../../../generated-code/protobuf/swagger
PROTO_SVC_PATH=$API_SERVER_SOURCE_PATH

mkdir -p $PROTO_SVC_PATH
cp -r target/generated-sources/swagger/* $PROTO_SVC_PATH/

echo "✅ Proto service files copied to $PROTO_SVC_PATH"

