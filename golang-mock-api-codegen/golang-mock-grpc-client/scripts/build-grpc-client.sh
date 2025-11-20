#!/bin/bash

echo "🚀 Building Golang Mock gRPC Client for Java/Adonis"
echo "=================================================="

cd "$(dirname "$0")/.."

echo "📦 Step 1: Building parent modules..."
cd ../..
mvn clean install -DskipTests -s settings.xml

echo ""
echo "🔧 Step 2: Building gRPC client module..."
cd golang-mock-api-codegen/golang-mock-grpc-client
mvn clean install -DskipTests

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ SUCCESS!"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "📦 Artifact created:"
    echo "   target/golang-mock-grpc-client-1.0.0-SNAPSHOT.jar"
    echo ""
    echo "📝 GroupId: com.nutanix.nutanix-core.ntnx-api.golang-mock-pc"
    echo "   ArtifactId: golang-mock-grpc-client"
    echo "   Version: 1.0.0-SNAPSHOT"
    echo ""
    echo "🔗 Next Steps:"
    echo "   1. Add dependency in ntnx-api-prism-service/pom.xml"
    echo "   2. Create Java controller to call gRPC service"
    echo "   3. Deploy to PC"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
else
    echo ""
    echo "❌ Build failed!"
    exit 1
fi

