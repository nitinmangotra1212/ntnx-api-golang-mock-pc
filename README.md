# ntnx-api-golang-mock-pc

API definitions and code generation for ntnx-api-golang-mock service.

## 📋 Overview

This repository contains:
- YAML API definitions (OpenAPI specs)
- Maven-based code generation
- Auto-generated Go DTOs with $objectType support
- Protocol Buffer definitions

## 🏗️ Structure

```
ntnx-api-golang-mock-pc/
├── pom.xml                              # Maven parent
├── settings.xml                         # Maven settings
├── golang-mock-api-definitions/         # YAML API definitions
│   └── defs/
│       └── namespaces/mock/v4/modules/config/released/
│           ├── models/catModel.yaml     # Cat schema
│           └── api/catEndpoint.yaml     # Cat endpoints
├── golang-mock-api-codegen/             # Code generators
│   ├── golang-mock-go-dto-definitions/  # DTO generator
│   └── pom.xml
└── generated-code/                      # Generated output
    └── dto/src/models/mock/v4/config/
        └── config_model.go              # Auto-generated DTOs
```

## 🚀 Build

### Prerequisites
- **Java 21+**
- **Maven 3.8+**
- **Go 1.23+**

### Generate Code

```bash
mvn clean install -s settings.xml
```

This generates:
- `generated-code/dto/src/models/mock/v4/config/config_model.go`
- Auto-generated constructors (NewCat(), NewLocation(), etc.)
- Auto-set $objectType and $reserved fields

## 📦 Usage

The service repository (`ntnx-api-golang-mock`) imports generated DTOs:

```go
// In ntnx-api-golang-mock/go.mod
require (
    github.com/nutanix/ntnx-api-golang-mock-pc/generated-code/dto v0.0.0
)
replace github.com/nutanix/ntnx-api-golang-mock-pc/generated-code/dto => 
    ../ntnx-api-golang-mock-pc/generated-code/dto/src
```

## 📝 Adding New APIs

1. Edit YAML:
   - Model: `defs/namespaces/mock/v4/modules/config/released/models/myModel.yaml`
   - API: `defs/namespaces/mock/v4/modules/config/released/api/myEndpoint.yaml`

2. Generate code:
   ```bash
   mvn clean install -s settings.xml
   ```

3. Generated DTOs will have auto-set $objectType!

## 🔗 Related

- **Service Implementation:** [ntnx-api-golang-mock](../ntnx-api-golang-mock)

## 📞 Contact

nitin.mangotra@nutanix.com
