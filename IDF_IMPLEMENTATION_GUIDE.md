# IDF Implementation Guide for Nexus Service

## What is IDF?

**IDF (Insights Data Format)** is Nutanix's internal database/storage layer used for:
- **Persistent storage** of entities (items, domains, etc.)
- **Querying with OData filters** (`$filter`, `$orderby`, `$select`, `$expand`)
- **Pagination** support
- **Complex queries** with grouping, sorting, and metrics

---

## How IDF Works in az-manager (Reference Implementation)

### Architecture Pattern

```
Service Layer (item_grpc_service.go)
    ↓
Repository Interface (db.ItemRepository)
    ↓
IDF Implementation (idf.ItemRepositoryImpl)
    ↓
IDF Client (external/idf/idf_client.go)
    ↓
Insights Service (go-cache/insights/insights_interface)
```

### Key Components in az-manager

1. **Repository Interface** (`db/domain_repository.go`)
   - Defines CRUD operations interface
   - Used by service layer

2. **IDF Repository Implementation** (`idf/idf_domain_repository.go`)
   - Implements repository interface
   - Converts protobuf models ↔ IDF attributes
   - Handles Create, Read, List, Update, Delete operations

3. **IDF Client** (`external/idf/idf_client.go`)
   - Wraps `InsightsService` from `go-cache/insights/insights_interface`
   - Provides: `GetEntityRet`, `UpdateEntityRet`, `GetEntitiesWithMetricsRet`

4. **IDF Utils** (`idf/idf_utils.go`)
   - `AddAttribute()` - Converts Go types to IDF `AttributeDataArg`
   - `CreateDataArg()` - Type conversion (string, int, bool, arrays, etc.)
   - `constructIDFQuery()` - Builds IDF queries from OData params

---

## Implementation Steps for Nexus Service

### Step 1: Add Dependencies

**File:** `ntnx-api-golang-nexus/go.mod`

```go
require (
    // ... existing dependencies ...
    github.com/nutanix-core/go-cache v0.0.0-20240613003120-de1e4c3ed003
    github.com/nutanix-core/ntnx-api-odata-go v1.0.27
)
```

**Command:**
```bash
cd ~/ntnx-api-golang-nexus
go get github.com/nutanix-core/go-cache@v0.0.0-20240613003120-de1e4c3ed003
go get github.com/nutanix-core/ntnx-api-odata-go@v1.0.27
go mod tidy
```

---

### Step 2: Create Repository Interface

**File:** `golang-nexus-service/db/item_repository.go` (NEW)

```go
package db

import (
    pb "github.com/nutanix/ntnx-api-golang-nexus-pc/generated-code/protobuf/nexus/v4/config"
    "github.com/nutanix/ntnx-api-golang-nexus/golang-nexus-service/models"
)

type ItemRepository interface {
    CreateItem(itemEntity *models.ItemEntity) error
    GetItemById(extId string) (*models.ItemEntity, error)
    ListItems(queryParams *models.QueryParams) ([]*pb.Item, int64, error)
    UpdateItem(extId string, itemEntity *models.ItemEntity) error
    DeleteItem(extId string) error
}
```

---

### Step 3: Create IDF Client

**File:** `golang-nexus-service/external/idf/interface.go` (NEW)

```go
package idf

import (
    "github.com/nutanix-core/go-cache/insights/insights_interface"
)

type IdfClientIfc interface {
    GetEntityRet(getArg *insights_interface.GetEntitiesArg) (*insights_interface.GetEntitiesRet, error)
    UpdateEntityRet(updateArg *insights_interface.UpdateEntityArg) (*insights_interface.UpdateEntityRet, error)
    GetEntitiesWithMetricsRet(getEntitiesWithMetricsArg *insights_interface.GetEntitiesWithMetricsArg) (*insights_interface.GetEntitiesWithMetricsRet, error)
    GetInsightsService() insights_interface.InsightsServiceInterface
}

type IdfClientImpl struct {
    IdfSvc insights_interface.InsightsServiceInterface
}

func NewIdfClient(host string, port uint16) IdfClientIfc {
    IdfService := insights_interface.NewInsightsService(host, port)
    return &IdfClientImpl{
        IdfSvc: IdfService,
    }
}
```

**File:** `golang-nexus-service/external/idf/idf_client.go` (NEW)

```go
package idf

import (
    "fmt"
    "github.com/nutanix-core/go-cache/insights/insights_interface"
    log "github.com/sirupsen/logrus"
)

func (idf *IdfClientImpl) GetEntityRet(getArg *insights_interface.GetEntitiesArg) (*insights_interface.GetEntitiesRet, error) {
    if getArg == nil {
        log.Error("Nil get argument while trying to read from IDF")
        return nil, fmt.Errorf("invalid argument")
    }
    getResponse := &insights_interface.GetEntitiesRet{}
    err := idf.IdfSvc.SendMsgWithTimeout("GetOperationIdf", getArg, getResponse, nil, 30)
    return getResponse, err
}

func (idf *IdfClientImpl) UpdateEntityRet(updateArg *insights_interface.UpdateEntityArg) (*insights_interface.UpdateEntityRet, error) {
    if updateArg == nil {
        log.Error("Invalid update argument")
        return nil, fmt.Errorf("invalid argument")
    }
    updateResponse := &insights_interface.UpdateEntityRet{}
    err := idf.IdfSvc.SendMsgWithTimeout("UpdateOperationIdf", updateArg, updateResponse, nil, 30)
    return updateResponse, err
}

func (idf *IdfClientImpl) GetEntitiesWithMetricsRet(getEntitiesWithMetricsArg *insights_interface.GetEntitiesWithMetricsArg) (*insights_interface.GetEntitiesWithMetricsRet, error) {
    if getEntitiesWithMetricsArg == nil {
        log.Error("Invalid getEntitiesWithMetrics argument")
        return nil, fmt.Errorf("invalid argument")
    }
    getResponse := &insights_interface.GetEntitiesWithMetricsRet{}
    err := idf.IdfSvc.SendMsgWithTimeout("GetEntitiesWithMetricsOperationIdf", getEntitiesWithMetricsArg, getResponse, nil, 30)
    return getResponse, err
}

func (idf *IdfClientImpl) GetInsightsService() insights_interface.InsightsServiceInterface {
    return idf.IdfSvc
}
```

---

### Step 4: Create IDF Utils

**File:** `golang-nexus-service/idf/idf_utils.go` (NEW)

```go
package idf

import (
    idfIfc "github.com/nutanix-core/go-cache/insights/insights_interface"
    log "github.com/sirupsen/logrus"
    "google.golang.org/protobuf/proto"
)

// AddAttribute adds an attribute to the attribute data arg list
func AddAttribute(attributeDataArgList *[]*idfIfc.AttributeDataArg, name string, value interface{}) {
    dataArg := CreateDataArg(name, value)
    if dataArg == nil {
        log.Errorf("failed to create data arg for attribute %s", name)
        return
    }
    *attributeDataArgList = append(*attributeDataArgList, dataArg)
}

// CreateDataArg creates a data arg for the given name and value based on the value type
func CreateDataArg(name string, value interface{}) *idfIfc.AttributeDataArg {
    dataValue := &idfIfc.DataValue{}

    switch val := value.(type) {
    case string:
        dataValue.ValueType = &idfIfc.DataValue_StrValue{StrValue: val}
    case int32:
        dataValue.ValueType = &idfIfc.DataValue_Int64Value{Int64Value: int64(val)}
    case int64:
        dataValue.ValueType = &idfIfc.DataValue_Int64Value{Int64Value: val}
    case bool:
        dataValue.ValueType = &idfIfc.DataValue_BoolValue{BoolValue: val}
    case []string:
        dataValue.ValueType = &idfIfc.DataValue_StrList_{
            StrList: &idfIfc.DataValue_StrList{ValueList: val},
        }
    default:
        log.Errorf("Unsupported type for attribute %s: %T", name, value)
        return nil
    }

    dataArg := &idfIfc.AttributeDataArg{
        AttributeData: &idfIfc.AttributeData{
            Name:  proto.String(name),
            Value: dataValue,
        },
    }
    return dataArg
}
```

---

### Step 5: Create IDF Repository Implementation

**File:** `golang-nexus-service/idf/idf_item_repository.go` (NEW)

```go
package idf

import (
    "github.com/google/uuid"
    "github.com/nutanix-core/go-cache/insights/insights_interface"
    idfQr "github.com/nutanix-core/go-cache/insights/insights_interface/query"
    pb "github.com/nutanix/ntnx-api-golang-nexus-pc/generated-code/protobuf/nexus/v4/config"
    "github.com/nutanix/ntnx-api-golang-nexus/golang-nexus-service/db"
    "github.com/nutanix/ntnx-api-golang-nexus/golang-nexus-service/external/idf"
    "github.com/nutanix/ntnx-api-golang-nexus/golang-nexus-service/models"
    log "github.com/sirupsen/logrus"
    "google.golang.org/protobuf/proto"
)

type ItemRepositoryImpl struct {
    idfClient idf.IdfClientIfc
}

const (
    itemEntityTypeName = "nexus_item"
    itemIdAttr         = "item_id"
    itemNameAttr       = "item_name"
    itemTypeAttr       = "item_type"
    descriptionAttr    = "description"
)

func NewItemRepository(idfClient idf.IdfClientIfc) db.ItemRepository {
    return &ItemRepositoryImpl{
        idfClient: idfClient,
    }
}

func (r *ItemRepositoryImpl) CreateItem(itemEntity *models.ItemEntity) error {
    itemUuid := uuid.New().String()
    attributeDataArgList := []*insights_interface.AttributeDataArg{}

    // Add item attributes
    if itemEntity.Item.ItemId != nil {
        AddAttribute(&attributeDataArgList, itemIdAttr, *itemEntity.Item.ItemId)
    }
    if itemEntity.Item.ItemName != nil {
        AddAttribute(&attributeDataArgList, itemNameAttr, *itemEntity.Item.ItemName)
    }
    if itemEntity.Item.ItemType != nil {
        AddAttribute(&attributeDataArgList, itemTypeAttr, *itemEntity.Item.ItemType)
    }
    if itemEntity.Item.Description != nil {
        AddAttribute(&attributeDataArgList, descriptionAttr, *itemEntity.Item.Description)
    }

    updateArg := &insights_interface.UpdateEntityArg{
        EntityGuid: &insights_interface.EntityGuid{
            EntityTypeName: proto.String(itemEntityTypeName),
            EntityId:       &itemUuid,
        },
        AttributeDataArgList: attributeDataArgList,
    }

    _, err := r.idfClient.UpdateEntityRet(updateArg)
    if err != nil {
        log.Errorf("Failed to create item: %v", err)
        return err
    }

    // Set extId
    if itemEntity.Item.Base == nil {
        itemEntity.Item.Base = &response.ExternalizableAbstractModel{}
    }
    itemEntity.Item.Base.ExtId = &itemUuid

    return nil
}

func (r *ItemRepositoryImpl) ListItems(queryParams *models.QueryParams) ([]*pb.Item, int64, error) {
    // Build IDF query
    query, err := idfQr.QUERY(itemEntityTypeName + "ListQuery").
        FROM(itemEntityTypeName).Proto()
    if err != nil {
        return nil, 0, err
    }

    // Add pagination
    page := queryParams.Page
    limit := queryParams.Limit
    if limit <= 0 {
        limit = 50
    }
    offset := page * limit

    if query.GroupBy == nil {
        query.GroupBy = &insights_interface.QueryGroupBy{}
    }
    query.GroupBy.RawLimit = &insights_interface.QueryLimit{
        Limit:  proto.Int64(int64(limit)),
        Offset: proto.Int64(int64(offset)),
    }

    queryArg := &insights_interface.GetEntitiesWithMetricsArg{
        Query: query,
    }

    // Query IDF
    queryResponse, err := r.idfClient.GetEntitiesWithMetricsRet(queryArg)
    if err != nil {
        return nil, 0, err
    }

    // Convert IDF entities to Item protobufs
    var items []*pb.Item
    groupResults := queryResponse.GetGroupResultsList()
    if len(groupResults) == 0 {
        return []*pb.Item{}, 0, nil
    }

    entities := groupResults[0].GetRawResults()
    for _, entity := range entities {
        item := r.mapIdfAttributeToItem(entity)
        items = append(items, item)
    }

    totalCount := groupResults[0].GetTotalEntityCount()
    return items, totalCount, nil
}

func (r *ItemRepositoryImpl) mapIdfAttributeToItem(entity *insights_interface.Entity) *pb.Item {
    item := &pb.Item{}
    for _, attr := range entity.GetAttributeDataMap() {
        switch attr.GetName() {
        case itemIdAttr:
            val := int32(attr.GetValue().GetInt64Value())
            item.ItemId = &val
        case itemNameAttr:
            val := attr.GetValue().GetStrValue()
            item.ItemName = &val
        case itemTypeAttr:
            val := attr.GetValue().GetStrValue()
            item.ItemType = &val
        case descriptionAttr:
            val := attr.GetValue().GetStrValue()
            item.Description = &val
        }
    }
    return item
}
```

---

### Step 6: Update Service to Use Repository

**File:** `golang-nexus-service/grpc/item_grpc_service.go`

Replace in-memory map with repository pattern (similar to az-manager).

---

### Step 7: Initialize IDF Client in Main

**File:** `golang-nexus-service/server/main.go`

```go
import (
    "github.com/nutanix/ntnx-api-golang-nexus/golang-nexus-service/external/idf"
    idfRepo "github.com/nutanix/ntnx-api-golang-nexus/golang-nexus-service/idf"
    "github.com/nutanix/ntnx-api-golang-nexus/golang-nexus-service/db"
)

func main() {
    // ... existing code ...

    // Initialize IDF client
    idfClient := idf.NewIdfClient("localhost", 9876)  // Adjust host/port as needed
    itemRepository := idfRepo.NewItemRepository(idfClient)

    // Pass repository to service
    itemService := grpc.NewItemGrpcService(itemRepository)
    
    // ... rest of code ...
}
```

---

## IDF Setup on PC

### Entity Type Name

**Important:** The entity type name must match what's configured in IDF on PC.

**Example from az-manager:**
- Domain entity: `"nim_domain"`
- BOM entity: `"nim_bom"`

**For Nexus:**
- Item entity: `"nexus_item"` (or as configured in your IDF setup script)

---

## PC Setup Script

You mentioned you have a script to run on PC to create IDF tables/schema. This script should:

1. **Define Entity Type:** Create entity type `nexus_item` in IDF
2. **Define Attributes:** Define all attributes (item_id, item_name, item_type, description, etc.)
3. **Set Permissions:** Configure access permissions if needed

**Typical script structure:**
```bash
#!/bin/bash
# setup_nexus_idf.sh

# Create entity type
# Define attributes
# Set up indexes if needed
```

---

## Next Steps

1. **Review your PC setup script** - Share it and we can integrate it
2. **Implement repository pattern** - Follow az-manager pattern
3. **Test IDF connection** - Verify connectivity to IDF service
4. **Migrate data** - Move from in-memory to IDF storage

---

## Reference Files from az-manager

- `az-manager-service/external/idf/idf_client.go` - IDF client implementation
- `az-manager-service/external/idf/interface.go` - IDF client interface
- `az-manager-service/idf/idf_utils.go` - IDF utility functions
- `az-manager-service/idf/idf_domain_repository.go` - Domain repository implementation
- `az-manager-service/db/domain_repository.go` - Repository interface

---

**Last Updated:** November 27, 2025

