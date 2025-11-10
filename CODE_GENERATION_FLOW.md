# 🔄 Code Generation Flow - Complete Guide

## 📊 Visual Flow Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                     START: YAML API Definitions                     │
│           (golang-mock-api-definitions/defs/.../catModel.yaml)      │
└────────────────────────────┬────────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────────┐
│                     Maven Build (mvn clean install)                 │
│                         (pom.xml executes)                          │
└────────────────────────────┬────────────────────────────────────────┘
                             │
                 ┌───────────┴───────────┐
                 │                       │
                 ▼                       ▼
    ┌─────────────────────┐   ┌─────────────────────┐
    │  Go DTOs            │   │  Proto Files        │
    │  (REST/JSON)        │   │  (gRPC)             │
    └──────────┬──────────┘   └──────────┬──────────┘
               │                         │
               ▼                         ▼
    ┌─────────────────────┐   ┌─────────────────────┐
    │ config_model.go     │   │ config.proto        │
    │ (NOT YET GENERATED) │   │ cat_service.proto   │
    │ Uses json tags      │   │ (GENERATED)         │
    └─────────────────────┘   └──────────┬──────────┘
                                         │
                                         ▼
                              ┌──────────────────────┐
                              │ ./generate-grpc.sh   │
                              │ (protoc compiler)    │
                              └──────────┬───────────┘
                                         │
                     ┌───────────────────┼───────────────────┐
                     │                   │                   │
                     ▼                   ▼                   ▼
         ┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐
         │ config.pb.go    │ │cat_service.pb.go│ │cat_service_grpc │
         │ (Messages)      │ │(Service Msgs)   │ │.pb.go           │
         │                 │ │                 │ │(gRPC Stubs)     │
         └────────┬────────┘ └────────┬────────┘ └────────┬────────┘
                  │                   │                   │
                  └───────────────────┴───────────────────┘
                                      │
                                      ▼
                          ┌───────────────────────┐
                          │   Used by Go Server   │
                          │   (grpc/cat_grpc_     │
                          │    service.go)        │
                          └───────────────────────┘
```

---

## 📁 Step-by-Step File Generation

### **Step 1: YAML API Definition (Source)**

**File:** `golang-mock-api-definitions/defs/namespaces/mock/versioned/v4/modules/config/released/models/catModel.yaml`

**Content Example:**
```yaml
Cat:
  description: Cat entity for mock REST API
  type: object
  properties:
    catId:
      type: integer
      format: int32
      description: Unique identifier for the cat
    catName:
      type: string
      description: Name of the cat
    catType:
      type: string
      description: Type of cat
    location:
      $ref: '#/components/schemas/Location'
```

**What it defines:** The data structure/schema for Cat entity

---

### **Step 2: Maven Build Triggers Code Generation**

**Command:**
```bash
cd /Users/nitin.mangotra/ntnx-api-golang-mock-pc
mvn clean install
```

**What Maven Does:**

#### **2a. Maven Plugin Configuration (pom.xml)**

```xml
<build>
  <plugins>
    <!-- Plugin 1: Generate Go DTOs -->
    <plugin>
      <groupId>io.swagger.codegen.v3</groupId>
      <artifactId>swagger-codegen-maven-plugin</artifactId>
      <executions>
        <execution>
          <id>generate-go-dto</id>
          <goals>
            <goal>generate</goal>
          </goals>
          <configuration>
            <inputSpec>defs/namespaces/mock/...</inputSpec>
            <language>go</language>
            <output>generated-code/dto</output>
          </configuration>
        </execution>
      </executions>
    </plugin>
    
    <!-- Plugin 2: Generate Proto Files -->
    <plugin>
      <groupId>com.nutanix</groupId>
      <artifactId>swagger-to-proto-maven-plugin</artifactId>
      <executions>
        <execution>
          <id>generate-proto</id>
          <goals>
            <goal>generate</goal>
          </goals>
          <configuration>
            <inputSpec>defs/namespaces/mock/...</inputSpec>
            <output>generated-code/protobuf/swagger</output>
          </configuration>
        </execution>
      </executions>
    </plugin>
  </plugins>
</build>
```

**Maven reads:** `catModel.yaml`  
**Maven generates:**
1. Go DTOs (for REST) → `generated-code/dto/models/config_model.go` (NOT YET CREATED)
2. Proto files (for gRPC) → `generated-code/protobuf/swagger/mock/v4/config/*.proto`

---

### **Step 3: Proto Files Generated**

**Files Created by Maven:**

#### **3a. config.proto** (Message Definitions)

**Location:** `generated-code/protobuf/swagger/mock/v4/config/config.proto`

**Content:**
```protobuf
syntax = "proto3";
package mock.v4.config;

message Cat {
  int32 cat_id = 1;
  string cat_name = 2;
  string cat_type = 3;
  string description = 4;
  Location location = 6;
  ObjectMapWrapper _reserved = 7;
}

message Location {
  Country country = 1;
  string city = 2;
  string zip = 3;
  ObjectMapWrapper _reserved = 4;
}

message Country {
  string state = 1;
  ObjectMapWrapper _reserved = 2;
}
```

**What it defines:** Protocol Buffer message structures

---

#### **3b. cat_service.proto** (Service Definition)

**Location:** `generated-code/protobuf/swagger/mock/v4/config/cat_service.proto`

**Content:**
```protobuf
syntax = "proto3";
package mock.v4.config;

import "mock/v4/config/config.proto";

// Request/Response messages
message ListCatsRequest {
  int32 page = 1;
  int32 limit = 2;
}

message ListCatsResponse {
  repeated Cat cats = 1;
  int32 total_available_results = 2;
}

// Service definition
service CatService {
  rpc ListCats (ListCatsRequest) returns (ListCatsResponse);
  rpc GetCat (GetCatRequest) returns (GetCatResponse);
  rpc CreateCat (CreateCatRequest) returns (CreateCatResponse);
  rpc UpdateCat (UpdateCatRequest) returns (UpdateCatResponse);
  rpc DeleteCat (DeleteCatRequest) returns (DeleteCatResponse);
  rpc GetCatAsync (GetCatAsyncRequest) returns (GetCatAsyncResponse);
}
```

**What it defines:** gRPC service interface and RPC methods

---

### **Step 4: Generate .pb.go Files (protoc)**

**Command:**
```bash
cd /Users/nitin.mangotra/ntnx-api-golang-mock-pc
./generate-grpc.sh
```

**Script Content (generate-grpc.sh):**
```bash
#!/bin/bash

# Install protoc plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Generate .pb.go from config.proto
protoc --go_out=generated-code/protobuf/mock/v4/config \
       --go_opt=module=github.com/nutanix/ntnx-api-golang-mock-pc \
       generated-code/protobuf/swagger/mock/v4/config/config.proto

# Generate .pb.go and _grpc.pb.go from cat_service.proto
protoc --go_out=generated-code/protobuf/mock/v4/config \
       --go_opt=module=github.com/nutanix/ntnx-api-golang-mock-pc \
       --go-grpc_out=generated-code/protobuf/mock/v4/config \
       --go-grpc_opt=module=github.com/nutanix/ntnx-api-golang-mock-pc \
       generated-code/protobuf/swagger/mock/v4/config/cat_service.proto
```

**What protoc does:**
1. Reads `.proto` files
2. Generates Go code using `protoc-gen-go` and `protoc-gen-go-grpc` plugins
3. Creates `.pb.go` files

---

### **Step 5: .pb.go Files Generated**

**Files Created:**

#### **5a. config.pb.go** (Message Implementations)

**Location:** `generated-code/protobuf/mock/v4/config/config.pb.go`

**Size:** ~11KB

**Content Example:**
```go
package config

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

// Cat message
type Cat struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CatId       int32                `protobuf:"varint,1,opt,name=cat_id,json=catId,proto3" json:"cat_id,omitempty"`
	CatName     string               `protobuf:"bytes,2,opt,name=cat_name,json=catName,proto3" json:"cat_name,omitempty"`
	CatType     string               `protobuf:"bytes,3,opt,name=cat_type,json=catType,proto3" json:"cat_type,omitempty"`
	Description string               `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	Location    *Location            `protobuf:"bytes,6,opt,name=location,proto3" json:"location,omitempty"`
	Reserved    *ObjectMapWrapper    `protobuf:"bytes,7,opt,name=_reserved,json=Reserved,proto3" json:"_reserved,omitempty"`
}

// Getter methods
func (x *Cat) GetCatId() int32 {
	if x != nil {
		return x.CatId
	}
	return 0
}

// ... more methods
```

**What it contains:**
- Go struct definitions with protobuf tags
- Getters/setters
- Marshaling/unmarshaling logic
- Protocol buffer metadata

---

#### **5b. cat_service.pb.go** (Service Message Implementations)

**Location:** `generated-code/protobuf/mock/v4/config/cat_service.pb.go`

**Size:** ~35KB

**Content Example:**
```go
package config

// ListCatsRequest message
type ListCatsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Page  int32 `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	Limit int32 `protobuf:"varint,2,opt,name=limit,proto3" json:"limit,omitempty"`
}

// ListCatsResponse message
type ListCatsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cats                  []*Cat `protobuf:"bytes,1,rep,name=cats,proto3" json:"cats,omitempty"`
	TotalAvailableResults int32  `protobuf:"varint,2,opt,name=total_available_results,json=totalAvailableResults,proto3" json:"total_available_results,omitempty"`
}

// ... more request/response messages
```

**What it contains:**
- Request/Response message structs
- Getters for all RPC messages

---

#### **5c. cat_service_grpc.pb.go** (gRPC Service Stubs)

**Location:** `generated-code/protobuf/mock/v4/config/cat_service_grpc.pb.go`

**Size:** ~19KB

**Content Example:**
```go
package config

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// CatServiceClient interface (for clients)
type CatServiceClient interface {
	ListCats(ctx context.Context, in *ListCatsRequest, opts ...grpc.CallOption) (*ListCatsResponse, error)
	GetCat(ctx context.Context, in *GetCatRequest, opts ...grpc.CallOption) (*GetCatResponse, error)
	CreateCat(ctx context.Context, in *CreateCatRequest, opts ...grpc.CallOption) (*CreateCatResponse, error)
	UpdateCat(ctx context.Context, in *UpdateCatRequest, opts ...grpc.CallOption) (*UpdateCatResponse, error)
	DeleteCat(ctx context.Context, in *DeleteCatRequest, opts ...grpc.CallOption) (*DeleteCatResponse, error)
	GetCatAsync(ctx context.Context, in *GetCatAsyncRequest, opts ...grpc.CallOption) (*GetCatAsyncResponse, error)
}

// CatServiceServer interface (for servers - YOU IMPLEMENT THIS)
type CatServiceServer interface {
	ListCats(context.Context, *ListCatsRequest) (*ListCatsResponse, error)
	GetCat(context.Context, *GetCatRequest) (*GetCatResponse, error)
	CreateCat(context.Context, *CreateCatRequest) (*CreateCatResponse, error)
	UpdateCat(context.Context, *UpdateCatRequest) (*UpdateCatResponse, error)
	DeleteCat(context.Context, *DeleteCatRequest) (*DeleteCatResponse, error)
	GetCatAsync(context.Context, *GetCatAsyncRequest) (*GetCatAsyncResponse, error)
	mustEmbedUnimplementedCatServiceServer()
}

// RegisterCatServiceServer function
func RegisterCatServiceServer(s grpc.ServiceRegistrar, srv CatServiceServer) {
	s.RegisterService(&CatService_ServiceDesc, srv)
}

// Client implementation
type catServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCatServiceClient(cc grpc.ClientConnInterface) CatServiceClient {
	return &catServiceClient{cc}
}

func (c *catServiceClient) ListCats(ctx context.Context, in *ListCatsRequest, opts ...grpc.CallOption) (*ListCatsResponse, error) {
	out := new(ListCatsResponse)
	err := c.cc.Invoke(ctx, "/mock.v4.config.CatService/ListCats", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ... more client methods

// Server stub implementation
type UnimplementedCatServiceServer struct {
}

func (UnimplementedCatServiceServer) ListCats(context.Context, *ListCatsRequest) (*ListCatsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListCats not implemented")
}

// ... more unimplemented methods
```

**What it contains:**
- `CatServiceClient` interface (for gRPC clients)
- `CatServiceServer` interface (for YOU to implement)
- `RegisterCatServiceServer` function (to register your implementation)
- Client stub implementations
- Unimplemented server stubs

---

### **Step 6: Your Implementation**

**File:** `grpc/cat_grpc_service.go`

**What YOU write:**
```go
package grpc

import (
	"context"
	pb "github.com/nutanix/ntnx-api-golang-mock-pc/generated-code/protobuf/mock/v4/config"
)

// CatGrpcService implements pb.CatServiceServer
type CatGrpcService struct {
	pb.UnimplementedCatServiceServer // Embed for forward compatibility
	cats map[int32]*pb.Cat
}

// ListCats implements the ListCats RPC
func (s *CatGrpcService) ListCats(ctx context.Context, req *pb.ListCatsRequest) (*pb.ListCatsResponse, error) {
	// YOUR BUSINESS LOGIC HERE
	var cats []*pb.Cat
	// ... populate cats ...
	
	return &pb.ListCatsResponse{
		Cats:                  cats,
		TotalAvailableResults: int32(len(s.cats)),
	}, nil
}

// ... implement other methods
```

---

### **Step 7: Server Setup**

**File:** `cmd/grpc-server/main.go`

**What YOU write:**
```go
package main

import (
	"net"
	"github.com/nutanix/ntnx-api-golang-mock/grpc"
	pb "github.com/nutanix/ntnx-api-golang-mock-pc/generated-code/protobuf/mock/v4/config"
	"google.golang.org/grpc"
)

func main() {
	lis, _ := net.Listen("tcp", ":50051")
	
	grpcServer := grpc.NewServer()
	catService := grpcService.NewCatGrpcService()
	
	// This function comes from cat_service_grpc.pb.go
	pb.RegisterCatServiceServer(grpcServer, catService)
	
	grpcServer.Serve(lis)
}
```

---

## 📊 Complete Generation Timeline

```
Time    | Step | File/Command                                    | Output
--------|------|------------------------------------------------|---------------------------
T0      | 1    | catModel.yaml                                  | (Manual creation)
        |      |                                                |
T1      | 2    | mvn clean install                              | Triggers Maven plugins
        |      |                                                |
T2      | 3a   | swagger-codegen-maven-plugin                   | config_model.go (NOT YET)
        |      |                                                |
T3      | 3b   | swagger-to-proto-maven-plugin                  | config.proto ✅
        |      |                                                | cat_service.proto ✅
        |      |                                                |
T4      | 4    | ./generate-grpc.sh                             | Runs protoc
        |      |                                                |
T5      | 5a   | protoc + protoc-gen-go                         | config.pb.go ✅
        |      |                                                |
T6      | 5b   | protoc + protoc-gen-go                         | cat_service.pb.go ✅
        |      |                                                |
T7      | 5c   | protoc + protoc-gen-go-grpc                    | cat_service_grpc.pb.go ✅
        |      |                                                |
T8      | 6    | (Manual) grpc/cat_grpc_service.go              | Your implementation
        |      |                                                |
T9      | 7    | go build -o bin/grpc-server ./cmd/grpc-server  | Binary executable
        |      |                                                |
T10     | 8    | ./bin/grpc-server                              | Running server! 🚀
```

---

## 🗂️ File Dependencies

```
catModel.yaml
    │
    └──> Maven Build (pom.xml)
            │
            ├──> config.proto
            │       │
            │       └──> protoc + protoc-gen-go
            │               │
            │               └──> config.pb.go
            │
            └──> cat_service.proto
                    │
                    └──> protoc + protoc-gen-go + protoc-gen-go-grpc
                            │
                            ├──> cat_service.pb.go
                            │
                            └──> cat_service_grpc.pb.go
                                    │
                                    └──> Used by:
                                            - cat_grpc_service.go (YOUR CODE)
                                            - cmd/grpc-server/main.go (YOUR CODE)
```

---

## 🎯 Key Takeaways

### **What Maven Generates:**
1. ✅ `.proto` files from YAML
2. ❌ NOT `.pb.go` files (you run `./generate-grpc.sh` separately)

### **What protoc Generates:**
1. ✅ `.pb.go` files (message implementations)
2. ✅ `_grpc.pb.go` files (gRPC service stubs)

### **What YOU Write:**
1. ✅ `cat_grpc_service.go` (implements `CatServiceServer`)
2. ✅ `cmd/grpc-server/main.go` (server setup)

### **What Gets Compiled:**
```
Your Code + Generated Code → Binary
    ↓
./bin/grpc-server
```

---

## 🔄 Quick Reference Commands

```bash
# Step 1: Generate .proto from YAML
cd /Users/nitin.mangotra/ntnx-api-golang-mock-pc
mvn clean install

# Step 2: Generate .pb.go from .proto
./generate-grpc.sh

# Step 3: Build gRPC server
cd /Users/nitin.mangotra/ntnx-api-golang-mock
go build -o bin/grpc-server ./cmd/grpc-server/main.go

# Step 4: Run server
./bin/grpc-server
```

---

## 📚 Files Summary Table

| File | Type | Generated By | Purpose |
|------|------|--------------|---------|
| `catModel.yaml` | YAML | Manual | API schema definition |
| `config.proto` | Proto | Maven | Message definitions |
| `cat_service.proto` | Proto | Maven | Service interface |
| `config.pb.go` | Go | protoc | Message implementations |
| `cat_service.pb.go` | Go | protoc | Service message implementations |
| `cat_service_grpc.pb.go` | Go | protoc | gRPC stubs (interfaces) |
| `cat_grpc_service.go` | Go | **Manual** | **Your business logic** |
| `cmd/grpc-server/main.go` | Go | **Manual** | **Your server setup** |

---

**🎉 Now you understand the complete flow from YAML to running gRPC server!**

