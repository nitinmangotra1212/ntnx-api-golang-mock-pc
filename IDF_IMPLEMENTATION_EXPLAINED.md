# IDF Implementation Explained

## Overview

The IDF (Insights Data Format) integration provides persistent storage, querying, and data management for the Nexus service. This document explains how IDF was implemented and the design decisions made.

---

## Architecture

### 1. **IDF Client Layer** (`external/idf/`)

**Purpose:** Abstraction layer for communicating with the IDF service.

**Files:**
- `interface.go` - Defines the IDF client interface
- `idf_client.go` - Implements IDF service method calls

**Key Components:**

```go
type IdfClientIfc interface {
    GetEntityRet(getArg *insights_interface.GetEntitiesArg) (*insights_interface.GetEntitiesRet, error)
    UpdateEntityRet(updateArg *insights_interface.UpdateEntityArg) (*insights_interface.UpdateEntityRet, error)
    GetEntitiesWithMetricsRet(getEntitiesWithMetricsArg *insights_interface.GetEntitiesWithMetricsArg) (*insights_interface.GetEntitiesWithMetricsRet, error)
    GetInsightsService() insights_interface.InsightsServiceInterface
}
```

**How it works:**
- Uses `github.com/nutanix-core/go-cache/insights/insights_interface` package
- Connects to IDF service via gRPC (host:port, typically `127.0.0.1:2027`)
- Wraps IDF service methods with proper error handling
- Uses `SendMsgWithTimeout()` to call IDF service methods

**Example:**
```go
idfClient := externalIdf.NewIdfClient("127.0.0.1", 2027)
response, err := idfClient.GetEntitiesWithMetricsRet(queryArg)
```

---

### 2. **Repository Pattern** (`idf/idf_item_repository.go`)

**Purpose:** Implements data access logic, mapping between protobuf models and IDF entities.

**Key Methods:**
- `ListItems()` - Query multiple items with pagination
- `GetItemById()` - Get single item by external ID
- `CreateItem()` - Create new item in IDF
- `UpdateItem()` - Update existing item
- `DeleteItem()` - Delete item (not yet implemented)

**How it works:**

#### **ListItems Implementation:**

```go
// 1. Build IDF query using query builder
query, err := idfQr.QUERY("itemListQuery").FROM("item").Proto()

// 2. Specify which columns to fetch (CRITICAL!)
query.GroupBy.RawColumns = []*insights_interface.QueryRawColumn{
    {Column: proto.String("item_id")},
    {Column: proto.String("item_name")},
    // ... other columns
}

// 3. Add pagination
query.GroupBy.RawLimit = &insights_interface.QueryLimit{
    Limit:  proto.Int64(int64(limit)),
    Offset: proto.Int64(int64(offset)),
}

// 4. Execute query
queryArg := &insights_interface.GetEntitiesWithMetricsArg{Query: query}
response, err := idfClient.GetEntitiesWithMetricsRet(queryArg)

// 5. Convert IDF entities to protobuf Items
entities := ConvertEntitiesWithMetricToEntities(response.GetGroupResultsList()[0].GetRawResults())
for _, entity := range entities {
    item := r.mapIdfAttributeToItem(entity)
    items = append(items, item)
}
```

**Key Points:**
- Uses `GetEntitiesWithMetrics` for list queries (supports metrics and attributes)
- Must explicitly specify `RawColumns` - without this, IDF returns entities but without attribute data!
- Converts `EntityWithMetric` to `Entity` using utility function
- Maps IDF attributes (snake_case) to protobuf fields (camelCase)

---

### 3. **Attribute Mapping** (`mapIdfAttributeToItem`)

**Purpose:** Converts IDF entity attributes to protobuf Item model.

**How it works:**

```go
func (r *ItemRepositoryImpl) mapIdfAttributeToItem(entity *insights_interface.Entity) *pb.Item {
    item := &pb.Item{}
    
    // Get extId from EntityGuid (UUID)
    if entity.GetEntityGuid() != nil {
        extId := entity.GetEntityGuid().GetEntityId()
        item.ExtId = &extId
    }
    
    // Map attributes (snake_case → camelCase)
    for _, attr := range entity.GetAttributeDataMap() {
        switch attr.GetName() {
        case "item_id":      // IDF → ItemId (protobuf)
            item.ItemId = int32(attr.GetValue().GetInt64Value())
        case "item_name":    // IDF → ItemName (protobuf)
            item.ItemName = attr.GetValue().GetStrValue()
        // ... etc
        }
    }
    return item
}
```

**Key Points:**
- IDF stores attributes in **snake_case** (`item_id`, `item_name`)
- Protobuf uses **camelCase** (`ItemId`, `ItemName`)
- `extId` comes from `EntityGuid.EntityId` (the UUID)
- Handles different data types (int64, string, bool)

---

### 4. **Query Parameter Extraction** (`utils/query/query_utils.go`)

**Purpose:** Extracts OData query parameters from HTTP requests.

**Current Implementation:**

```go
func ExtractQueryParamsFromContext(ctx context.Context) *models.QueryParams {
    // Get original HTTP path from gRPC context metadata
    path := responseUtils.GetPathFromGrpcContext(ctx)
    
    // Parse URL to extract query parameters
    parsedURL, err := url.Parse(path)
    values := parsedURL.Query()
    
    // Extract OData parameters
    queryParams.Filter = values.Get("$filter")
    queryParams.Orderby = values.Get("$orderby")
    queryParams.Select = values.Get("$select")
    queryParams.Expand = values.Get("$expand")
    queryParams.Page = parseInt(values.Get("$page"))
    queryParams.Limit = parseInt(values.Get("$limit"))
    
    return queryParams
}
```

**What it does:**
- Extracts OData query parameters from the HTTP request path
- Gets the path from gRPC context metadata (Adonis forwards it via `x-envoy-original-path` header)
- Parses basic query parameters manually

**What it does NOT do (yet):**
- ❌ Parse `$filter` expressions (e.g., `itemName eq 'test'`)
- ❌ Parse `$orderby` expressions (e.g., `itemName asc, itemId desc`)
- ❌ Convert OData field names to IDF attribute names
- ❌ Build IDF `WhereClause` from `$filter`

---

## OData Support: Current vs. Full Implementation

### **Current Implementation (Basic)**

✅ **Supported:**
- `$page` - Pagination page number
- `$limit` - Page size
- `$filter` - Stored but not parsed (passed as string)
- `$orderby` - Stored but not parsed (passed as string)
- `$select` - Stored but not parsed (passed as string)
- `$expand` - Stored but not parsed (passed as string)

❌ **Not Supported:**
- OData expression parsing (`$filter`, `$orderby`)
- Field name mapping (OData → IDF attributes)
- Complex queries (AND, OR, functions)
- IDF `WhereClause` generation from `$filter`

### **Full Implementation (like az-manager)**

**az-manager uses:**
1. **OData Parser** (`github.com/nutanix-core/ntnx-api-odata-go`)
   - Parses `$filter`, `$orderby`, `$select`, `$expand`
   - Converts OData expressions to IDF query clauses

2. **EDM (Entity Data Model) Provider**
   - Maps OData field names to IDF attribute names
   - Defines entity relationships

3. **IDF Query Evaluator**
   - Converts parsed OData to IDF `Query` objects
   - Generates `WhereClause` from `$filter`
   - Generates `GroupSortOrder` from `$orderby`

**Example from az-manager:**

```go
// 1. Create EDM provider with entity bindings
edmProvider := edm.NewCustomEdmProvider(entityBindingList)

// 2. Create OData parser
odataParser := parser.NewParser(edmProvider)

// 3. Parse OData query parameters
queryParam := parser.NewQueryParam()
queryParam.SetFilter(queryParams.Filter)  // e.g., "itemName eq 'test'"
queryParam.SetOrderBy(queryParams.Orderby) // e.g., "itemName asc"

uriInfo, err := odataParser.ParserWithQueryParam(queryParam, resourcePath)

// 4. Convert to IDF query
idfQueryEval := idf.IDFQueryEvaluator{}
idfQuery, err := idfQueryEval.GetQuery(uriInfo, resourcePath)

// 5. Use in IDF query
query.WhereClause = idfQuery.GetWhereClause()
query.GroupBy.RawSortOrder = idfQuery.GetGroupBy().GetGroupSortOrder()
```

---

## Why Not Full OData Yet?

**Current approach (manual extraction):**
- ✅ Simple and works for basic pagination
- ✅ No additional dependencies
- ✅ Fast to implement
- ❌ Limited query capabilities

**Full OData approach (like az-manager):**
- ✅ Full OData query support
- ✅ Complex filtering and sorting
- ✅ Field name mapping
- ❌ Requires `ntnx-api-odata-go` dependency
- ❌ Requires EDM bindings configuration
- ❌ More complex implementation

**Decision:** Start with basic implementation, add full OData support later if needed.

---

## IDF Query Builder

**Package:** `github.com/nutanix-core/go-cache/insights/insights_interface/query`

**Usage:**
```go
import idfQr "github.com/nutanix-core/go-cache/insights/insights_interface/query"

// Build query
query, err := idfQr.QUERY("itemListQuery").
    FROM("item").
    Proto()

// Query structure:
type Query struct {
    EntityList    []*EntityGuid      // For GetById queries
    GroupBy       *QueryGroupBy      // Columns, sorting, pagination
    WhereClause   *QueryWhereClause  // Filter conditions (from OData $filter)
}
```

**Key Components:**
- `RawColumns` - Which attributes to fetch
- `RawLimit` - Pagination (limit, offset)
- `RawSortOrder` - Sorting (from OData `$orderby`)
- `WhereClause` - Filter conditions (from OData `$filter`)

---

## Data Flow

### **List Items Request:**

```
1. HTTP Request
   GET /api/nexus/v4.1/config/items?$page=0&$limit=50
   ↓
2. Adonis (prism-service)
   - Extracts query params
   - Calls gRPC: ItemService.ListItems()
   - Forwards original path in metadata
   ↓
3. Go gRPC Service
   - ExtractQueryParamsFromContext() extracts $page, $limit
   - Calls itemRepository.ListItems(queryParams)
   ↓
4. IDF Repository
   - Builds IDF query with pagination
   - Specifies RawColumns (item_id, item_name, etc.)
   - Calls idfClient.GetEntitiesWithMetricsRet()
   ↓
5. IDF Service
   - Executes query
   - Returns EntityWithMetric list
   ↓
6. Repository
   - Converts EntityWithMetric → Entity
   - Maps IDF attributes → protobuf Item
   ↓
7. gRPC Service
   - Wraps in ListItemsRet with metadata
   ↓
8. Adonis
   - Converts to JSON response
   ↓
9. HTTP Response
   {
     "data": [
       {
         "itemId": 1,
         "itemName": "test item 0",
         "extId": "550e8400-...",
         ...
       }
     ],
     "metadata": {...}
   }
```

---

## Key Design Decisions

### 1. **Singleton Pattern for IDF Client**

```go
// external/singleton.go
func Interfaces() NexusInterfaces {
    singletonServiceOnce.Do(func() {
        singleton = &singletonService{}
    })
    return singleton
}

func (s *singletonService) IdfClient() idf.IdfClientIfc {
    idfClientOnce.Do(func() {
        s.idfClient = idf.NewIdfClient(constants.IdfHost, uint16(constants.IdfPort))
    })
    return s.idfClient
}
```

**Why:** Ensures single IDF connection, thread-safe initialization, matches az-manager pattern.

### 2. **Repository Pattern**

**Why:** Separates data access logic from business logic, makes testing easier, allows swapping implementations.

### 3. **Explicit Column Selection**

```go
query.GroupBy.RawColumns = []*insights_interface.QueryRawColumn{
    {Column: proto.String("item_id")},
    {Column: proto.String("item_name")},
    // ...
}
```

**Why:** Without this, IDF returns entities but without attribute data populated. This was a critical bug fix!

### 4. **EntityWithMetric → Entity Conversion**

```go
entities := ConvertEntitiesWithMetricToEntities(entitiesWithMetric)
```

**Why:** `GetEntitiesWithMetrics` returns `EntityWithMetric`, but mapping function expects `Entity`. Conversion extracts attributes from metrics.

---

## Dependencies

**Required Go packages:**
```go
github.com/nutanix-core/go-cache/insights/insights_interface
github.com/nutanix-core/go-cache/insights/insights_interface/query
```

**For full OData support (future):**
```go
github.com/nutanix-core/ntnx-api-odata-go
```

---

## Summary

**Current Implementation:**
- ✅ Basic IDF integration (CRUD operations)
- ✅ Pagination support
- ✅ Attribute mapping (snake_case ↔ camelCase)
- ✅ Singleton pattern for IDF client
- ✅ Repository pattern for data access
- ⚠️ Basic OData parameter extraction (not parsed)

**Future Enhancements:**
- 🔄 Full OData parser integration (like az-manager)
- 🔄 `$filter` expression parsing
- 🔄 `$orderby` expression parsing
- 🔄 Field name mapping via EDM
- 🔄 Complex query support

**The implementation follows the same patterns as `az-manager` but starts with a simpler approach, adding complexity only when needed.**

