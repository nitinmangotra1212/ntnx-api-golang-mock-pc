# Nexus Service - Complete Deployment Guide

**Service Name:** Nexus (formerly Mock)  
**API Path:** `/api/nexus/v4.1/config`  
**gRPC Port:** 9090  
**Adonis Port:** 8888  
**Mercury Port:** 9440

---

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Build Steps](#build-steps)
3. [File Modifications](#file-modifications)
4. [Deployment to PC](#deployment-to-pc)
5. [Verification Steps](#verification-steps)
6. [Troubleshooting](#troubleshooting)

---

## Prerequisites

### Required Repositories
- `ntnx-api-golang-nexus-pc` - Code generation (Maven)
- `ntnx-api-golang-nexus` - Go gRPC server
- `ntnx-api-prism-service` - Adonis (REST-to-gRPC gateway)

### Required Tools
- Java 21
- Maven 3.8+
- Go 1.21+
- `protoc` (Protocol Buffer compiler)
- `protoc-gen-go` and `protoc-gen-go-grpc` plugins

### GitHub Authentication
**⚠️ IMPORTANT:** Maven build requires GitHub authentication to fetch common templates.  
See `GITHUB_AUTH_SETUP.md` for setup instructions.

---

## Build Steps

### Step 1: Build Code Generation (Maven)

```bash
cd ~/ntnx-api-golang-nexus-pc
mvn clean install -DskipTests -s settings.xml
```

**Expected Output:**
- `golang-nexus-api-definitions-17.0.0-SNAPSHOT.jar`
- `golang-nexus-grpc-client-17.0.0-SNAPSHOT.jar`
- Generated protobuf files in `generated-code/protobuf/swagger/nexus/v4/config/`

**Verification:**
```bash
ls -lh golang-nexus-api-codegen/golang-nexus-grpc-client/target/*.jar
jar -tf golang-nexus-api-codegen/golang-nexus-grpc-client/target/*.jar | \
  grep "NexusConfigItemController.class"
# Expected: nexus/v4/config/server/controllers/NexusConfigItemController.class
```

---

### Step 2: Generate gRPC Go Code

```bash
cd ~/ntnx-api-golang-nexus-pc
./generate-grpc.sh
```

**What This Does:**
- Generates `.pb.go` files from `.proto` definitions
- Fixes import paths for Go modules
- Ensures method names are lowercase (camelCase)

**Expected Output:**
- `generated-code/protobuf/nexus/v4/config/config.pb.go`
- `generated-code/protobuf/nexus/v4/config/item_service.pb.go`
- `generated-code/protobuf/nexus/v4/config/item_service_grpc.pb.go`

**Verification:**
```bash
grep "MethodName:" generated-code/protobuf/nexus/v4/config/item_service_grpc.pb.go
# Expected: MethodName: "listItems" (lowercase)
```

---

### Step 3: Build Go Server Binary

```bash
cd ~/ntnx-api-golang-nexus
make build
```

**Expected Output:**
- `golang-nexus-server` (Linux binary, ~15MB)

**Verification:**
```bash
file golang-nexus-server
# Expected: ELF 64-bit LSB executable, x86-64
ls -lh golang-nexus-server
# Expected: ~15MB
```

---

### Step 4: Build Adonis (prism-service)

```bash
cd ~/ntnx-api-prism-service
mvn clean install -DskipTests -s settings.xml
```

**Expected Output:**
- `target/prism-service-17.6.0-SNAPSHOT.jar` (~350MB)

**Verification:**
```bash
jar -tf target/prism-service-*.jar | grep "golang-nexus-grpc-client.*jar$"
# Expected: BOOT-INF/lib/golang-nexus-grpc-client-17.0.0-SNAPSHOT.jar
```

---

## File Modifications

### 1. PC: `application.yaml`

**File:** `/home/nutanix/adonis/config/application.yaml`

**Add to `adonis.controller.packages.onprem`:**
```yaml
adonis:
  controller:
    packages:
      onprem: |
        nexus.v4.server.configuration, \
        nexus.v4.config.server.controllers, \
        nexus.v4.config.server.services, \
        # ... other packages ...
```

**Add to `grpc:` section:**
```yaml
grpc:
  # ... other grpc services ...
  golangnexus:
    host: localhost
    port: 9090
```

---

### 2. PC: `lookup_cache.json`

**File:** `/home/nutanix/api_artifacts/lookup_cache.json`

**Add entry to `routeMappings` array:**
```json
{
  "apiPath": "/nexus/v4.1/config",
  "artifactPath": "nexus/v4.r1.a1/golang-nexus-api-definitions-17.0.0-SNAPSHOT"
}
```

**Command:**
```bash
# Backup original
cp /home/nutanix/api_artifacts/lookup_cache.json \
   /home/nutanix/api_artifacts/lookup_cache.json.backup

# Edit file
vi /home/nutanix/api_artifacts/lookup_cache.json
```

---

### 3. PC: Mercury Configuration

**File:** `~/config/mercury/mercury_request_handler_config_golangnexus.json`

**Create file with:**
```json
{
  "api_path_config_list" : [
    {
      "api_path" : "/api/nexus/v4.1",
      "handler_list" : [
        {
          "priority" : 1,
          "port" : 8888,
          "transport_options" : "kHttp",
          "external_request_auth_options" : "kAllowAnyAuthenticatedUserExt",
          "internal_request_auth_options" : "kAllowAnyAuthenticatedUserInt"
        }
      ]
    },
    {
      "api_path" : "/api/nexus/unversioned",
      "handler_list" : [
        {
          "priority" : 1,
          "port" : 8888,
          "transport_options" : "kHttp",
          "external_request_auth_options" : "kAllowAnyAuthenticatedUserExt",
          "internal_request_auth_options" : "kAllowAnyAuthenticatedUserInt"
        }
      ]
    }
  ]
}
```

**Command:**
```bash
mkdir -p ~/config/mercury
vi ~/config/mercury/mercury_request_handler_config_golangnexus.json
```

---

## Deployment to PC

### Step 1: Copy Artifacts

**From your Mac:**

```bash
# 1. Copy API artifacts (proto files)
mkdir -p /tmp/nexus-artifacts/nexus/v4.r1.a1/golang-nexus-api-definitions-17.0.0-SNAPSHOT
cd ~/ntnx-api-golang-nexus-pc
cp generated-code/protobuf/swagger/nexus/v4/config/*.proto \
   /tmp/nexus-artifacts/nexus/v4.r1.a1/golang-nexus-api-definitions-17.0.0-SNAPSHOT/
scp -r -O /tmp/nexus-artifacts/nexus nutanix@<PC_IP>:/home/nutanix/api_artifacts/

# 2. Copy Go binary
scp -O ~/ntnx-api-golang-nexus/golang-nexus-server \
   nutanix@<PC_IP>:/home/nutanix/golang-nexus-build/

# 3. Copy Adonis JAR
scp -O ~/ntnx-api-prism-service/target/prism-service-17.6.0-SNAPSHOT.jar \
   nutanix@<PC_IP>:/home/nutanix/adonis/lib/
```

---

### Step 2: Configure PC

**SSH to PC:**
```bash
ssh nutanix@<PC_IP>
```

**Update files as described in [File Modifications](#file-modifications) section above.**

---

### Step 3: Start Services

**On PC:**

```bash
# 1. Start golang-nexus-server
cd ~/golang-nexus-build
chmod +x golang-nexus-server
nohup ./golang-nexus-server -port 9090 -log-level debug > \
    golang-nexus-server.log 2>&1 &

# 2. Restart Adonis and Mercury
genesis stop adonis mercury
sleep 10
cluster start

# Wait for services to start (2-3 minutes)
sleep 120
```

---

## Verification Steps

### 1. Verify Services Are Running

```bash
# Check golang-nexus-server
ps aux | grep golang-nexus-server | grep -v grep
netstat -tlnp | grep 9090
# Expected: tcp6 ... :::9090 ... LISTEN

# Check Adonis
ps aux | grep adonis | grep java
netstat -tlnp | grep 8888
# Expected: tcp6 ... :::8888 ... LISTEN

# Check Mercury
ps aux | grep mercury | grep -v grep
netstat -tlnp | grep 9440
# Expected: tcp6 ... :::9440 ... LISTEN
```

---

### 2. Verify Configuration Files

```bash
# Check application.yaml has nexus packages
grep -A 5 "nexus.v4" /home/nutanix/adonis/config/application.yaml

# Check lookup_cache.json has nexus entry
grep -A 3 "nexus/v4.1" /home/nutanix/api_artifacts/lookup_cache.json

# Check Mercury config exists
ls -lh ~/config/mercury/mercury_request_handler_config_golangnexus.json
```

---

### 3. Verify Adonis Logs

```bash
# Check if nexus controller loaded
grep -i "nexus\|NexusConfigItem" /var/log/prism-service/*.log | tail -10

# Check for errors
tail -50 /var/log/prism-service/*.log | grep -i error | tail -10
```

**Expected:** No errors related to nexus or NexusConfigItemController.

---

### 4. Test API Endpoints

**Get Authentication Token:**
```bash
TOKEN=$(curl -k -s -X POST "https://<PC_IP>:9440/api/nutanix/v3/users/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"<YOUR_PASSWORD>"}' \
  | python3 -c "import sys, json; print(json.load(sys.stdin)['access_token'])")

echo "Token: ${TOKEN:0:20}..."
```

**Test Direct Adonis (Port 8888):**
```bash
curl -k "https://<PC_IP>:8888/api/nexus/v4.1/config/items" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" | python3 -m json.tool
```

**Expected:** JSON response with items array.

**Test via Mercury (Port 9440):**
```bash
curl -k "https://<PC_IP>:9440/api/nexus/v4.1/config/items" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" | python3 -m json.tool
```

**Expected:** Same JSON response.

---

### 5. Verify Response Structure

**Expected Response:**
```json
{
  "data": [
    {
      "itemId": 1,
      "name": "Item 1",
      "description": "Description 1"
    },
    ...
  ],
  "metadata": {
    "totalAvailableResults": 10,
    "isPaginated": true,
    "links": [...]
  }
}
```

---

## Troubleshooting

### Issue: `UNIMPLEMENTED: unknown method listItems`

**Cause:** Go server binary has wrong method names.

**Solution:**
1. Regenerate gRPC code: `cd ~/ntnx-api-golang-nexus-pc && ./generate-grpc.sh`
2. Rebuild Go server: `cd ~/ntnx-api-golang-nexus && make build`
3. Copy new binary to PC and restart server

---

### Issue: Mercury Not Listening on Port 9440

**Cause:** Mercury needs restart after config changes.

**Solution:**
```bash
genesis stop mercury
cluster start mercury
sleep 60
netstat -tlnp | grep 9440
```

---

### Issue: Adonis Returns 404 Not Found

**Check:**
1. `lookup_cache.json` has correct entry
2. API artifacts are in correct location
3. Adonis logs for routing errors

**Solution:**
```bash
# Verify lookup_cache.json
cat /home/nutanix/api_artifacts/lookup_cache.json | python3 -m json.tool

# Verify artifacts exist
ls -lh /home/nutanix/api_artifacts/nexus/v4.r1.a1/golang-nexus-api-definitions-17.0.0-SNAPSHOT/

# Restart Adonis
genesis stop adonis
cluster start adonis
```

---

### Issue: `git-upload-pack not permitted` During Maven Build

**Cause:** Missing GitHub authentication.

**Solution:** See `GITHUB_AUTH_SETUP.md` for detailed instructions.

**Quick Fix:**
1. Create GitHub Personal Access Token
2. Update `repositories.yaml` with token in URL:
   ```yaml
   uri: "https://<TOKEN>@github.com/nutanix-core/ntnx-api-dev-platform.git"
   ```

---

## Summary

### Files Modified on PC

| File | Location | Purpose |
|------|----------|---------|
| `application.yaml` | `/home/nutanix/adonis/config/` | Adonis controller packages and gRPC config |
| `lookup_cache.json` | `/home/nutanix/api_artifacts/` | API path routing |
| `mercury_request_handler_config_golangnexus.json` | `~/config/mercury/` | Mercury routing config |

### Artifacts Deployed

| Artifact | Location | Purpose |
|----------|----------|---------|
| Go Binary | `/home/nutanix/golang-nexus-build/golang-nexus-server` | gRPC server |
| Adonis JAR | `/home/nutanix/adonis/lib/prism-service-17.6.0-SNAPSHOT.jar` | REST-to-gRPC gateway |
| API Artifacts | `/home/nutanix/api_artifacts/nexus/v4.r1.a1/...` | Proto files for routing |

### Service Ports

| Service | Port | Protocol |
|---------|------|----------|
| golang-nexus-server | 9090 | gRPC |
| Adonis | 8888 | HTTP |
| Mercury | 9440 | HTTPS |

---

## Quick Reference

**Build Commands:**
```bash
# 1. Build code generation
cd ~/ntnx-api-golang-nexus-pc && mvn clean install -DskipTests -s settings.xml

# 2. Generate gRPC code
cd ~/ntnx-api-golang-nexus-pc && ./generate-grpc.sh

# 3. Build Go server
cd ~/ntnx-api-golang-nexus && make build

# 4. Build Adonis
cd ~/ntnx-api-prism-service && mvn clean install -DskipTests -s settings.xml
```

**Deployment Commands:**
```bash
# Copy artifacts
scp -r -O /tmp/nexus-artifacts/nexus nutanix@<PC_IP>:/home/nutanix/api_artifacts/
scp -O ~/ntnx-api-golang-nexus/golang-nexus-server nutanix@<PC_IP>:/home/nutanix/golang-nexus-build/
scp -O ~/ntnx-api-prism-service/target/prism-service-*.jar nutanix@<PC_IP>:/home/nutanix/adonis/lib/
```

**Service Management:**
```bash
# Start golang-nexus-server
cd ~/golang-nexus-build && nohup ./golang-nexus-server -port 9090 > golang-nexus-server.log 2>&1 &

# Restart Adonis and Mercury
genesis stop adonis mercury && cluster start
```

---

**Last Updated:** November 26, 2025  
**Service Version:** 17.0.0-SNAPSHOT

