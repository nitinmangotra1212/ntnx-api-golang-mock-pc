# Simple Deployment Steps - Complete Guide

## Current Status ✅

- ✅ Build complete
- ✅ Binary ready: `golang-nexus-server`
- ✅ IDF script ready: `setup_nexus_idf.py`
- ✅ Service uses in-memory storage (no IDF integration yet)

---

## Complete Deployment Steps

### Step 1: Rebuild Prism-Service JAR (if gRPC client changed)

**Check if needed:**
- If you modified YAML models → **YES, rebuild needed**
- If you only changed Go service code → **NO, rebuild not needed**

**If rebuild needed:**
```bash
cd ~/ntnx-api-prism-service
mvn clean install -DskipTests -s settings.xml
```

**Note:** Always use `-s settings.xml` to avoid proxy/authentication issues.

**Verify JAR:**
```bash
ls -lh target/prism-service-*.jar

# Note: jar -tf doesn't work for nested JARs in Spring Boot fat JARs
# The controller is inside: BOOT-INF/lib/golang-nexus-grpc-client-*.jar

# Method 1: Check Maven repo JAR directly (easiest)
jar -tf ~/.m2/repository/com/nutanix/nutanix-core/ntnx-api/golang-nexus-pc/golang-nexus-grpc-client/17.0.0-SNAPSHOT/golang-nexus-grpc-client-17.0.0-SNAPSHOT.jar | grep Controller
# Expected: nexus/v4/config/server/controllers/NexusConfigItemController.class

# Method 2: Extract nested JAR (if needed)
cd /tmp && mkdir verify && cd verify
unzip -q ~/ntnx-api-prism-service/target/prism-service-*.jar 'BOOT-INF/lib/golang-nexus-grpc-client-*.jar'
unzip -q BOOT-INF/lib/golang-nexus-grpc-client-17.0.0-SNAPSHOT.jar
find . -name '*NexusConfigItemController*'
```

---

### Step 2: Run IDF Setup Script on PC

**Purpose:** Register entity type and columns in IDF (one-time setup)

**On your Mac:**
```bash
# Copy script to PC
scp -O ~/ntnx-api-golang-nexus-pc/setup_nexus_idf.py \
   nutanix@10.124.86.254:/home/nutanix/
```

**On PC:**
```bash
# SSH to PC
ssh nutanix@10.124.86.254

# Run script
cd /home/nutanix
python3 setup_nexus_idf.py
```

**Expected Output:**
```
✅ Successfully registered entity type: item
✅ Successfully registered metric types (attributes)
📝 Creating sample items in IDF...
  ✅ Created 10 items...
  ✅ Created 20 items...
  ...
✅ Successfully created 110 sample items in IDF!
```

**Verify IDF Data:**
- Access: `http://10.124.86.254:2027/`
- Should see `item` entity type with 110 items

---

### Step 3: Copy Artifacts to PC

**On your Mac:**

```bash
# 1. Copy Go binary
scp -O ~/ntnx-api-golang-nexus/golang-nexus-server \
   nutanix@10.124.86.254:/home/nutanix/golang-nexus-build/

# 2. Copy Prism-Service JAR (if rebuilt)
scp -O ~/ntnx-api-prism-service/target/prism-service-*.jar \
   nutanix@10.124.86.254:/home/nutanix/adonis/lib/

# 3. Copy API artifacts (proto files) - REQUIRED for Adonis routing
mkdir -p /tmp/nexus-artifacts/nexus/v4.r1.a1/golang-nexus-api-definitions-17.0.0-SNAPSHOT
cd ~/ntnx-api-golang-nexus-pc
cp generated-code/protobuf/swagger/nexus/v4/config/*.proto \
   /tmp/nexus-artifacts/nexus/v4.r1.a1/golang-nexus-api-definitions-17.0.0-SNAPSHOT/
scp -r -O /tmp/nexus-artifacts/nexus \
   nutanix@10.124.86.254:/home/nutanix/api_artifacts/
```

---

### Step 4: Verify Artifacts on PC

**On PC:**
```bash
# SSH to PC
ssh nutanix@10.124.86.254

# 1. Verify Go binary
ls -lh ~/golang-nexus-build/golang-nexus-server
file ~/golang-nexus-build/golang-nexus-server
# Expected: ELF 64-bit LSB executable, x86-64

# 2. Verify Adonis JAR
ls -lh ~/adonis/lib/prism-service-*.jar
# Backup old LKG-RELEASE.jar if exists
cd ~/adonis/lib
mv prism-service-*-LKG-RELEASE.jar prism-service-*-LKG-RELEASE.jar.backup 2>/dev/null || true

# Verify new JAR contains controller (extract nested JAR to check)
cd /tmp && mkdir verify-controller && cd verify-controller
unzip -q ~/adonis/lib/prism-service-*.jar 'BOOT-INF/lib/golang-nexus-grpc-client-*.jar'
unzip -q BOOT-INF/lib/golang-nexus-grpc-client-17.0.0-SNAPSHOT.jar
find . -name '*NexusConfigItemController*'
# Expected: ./nexus/v4/config/server/controllers/NexusConfigItemController.class
cd /tmp && rm -rf verify-controller

# 3. Verify API artifacts
ls -lh ~/api_artifacts/nexus/v4.r1.a1/golang-nexus-api-definitions-17.0.0-SNAPSHOT/*.proto
# Expected: item_service.proto, config.proto, etc.
```

---

### Step 5: Start Services on PC

**On PC:**
```bash
# Stop old server (if running)
pkill -f golang-nexus-server

# Start Go server
cd ~/golang-nexus-build
chmod +x golang-nexus-server
nohup ./golang-nexus-server -port 9090 -log-level debug > \
    golang-nexus-server.log 2>&1 &

# Check if running
ps aux | grep golang-nexus-server
tail -f golang-nexus-server.log

# Restart Adonis (to load new JAR)
genesis stop adonis mercury
cluster start 

# Wait for Adonis to start
# Check status
genesis status adonis
```

---

### Step 6: Verify API Works

**Test from your Mac:**
```bash
# Get authentication token (if needed)
# Then test API
curl -k -X GET "https://10.124.86.254:9440/api/nexus/v4.1/config/items" \
  -H "Authorization: Basic <base64_credentials>" \
  -H "Content-Type: application/json"
```

**Expected Response:**
```json
{
  "data": [
    {
      "itemId": 1,
      "itemName": "Item-1",
      "itemType": "TYPE1",
      "description": "A fluffy item"
    },
    ...
  ],
  "metadata": {
    "totalAvailableResults": 100,
    "isPaginated": false
  }
}
```

---

## Files That Need to Be Copied

### Required Files:
1. ✅ **golang-nexus-server** (Go binary) → `/home/nutanix/golang-nexus-build/`
2. ✅ **prism-service-*.jar** (Adonis JAR) → `/home/nutanix/adonis/lib/`
3. ✅ **API artifacts (proto files)** → `/home/nutanix/api_artifacts/nexus/v4.r1.a1/golang-nexus-api-definitions-17.0.0-SNAPSHOT/`

**Why API artifacts are needed:**
- Adonis uses these proto files for routing and API discovery
- `lookup_cache.json` references this path
- Without them, Adonis will return 404 errors

---

## When to Rebuild Prism-Service

**Rebuild if:**
- ✅ YAML models changed (itemModel.yaml)
- ✅ API endpoints changed (itemEndpoint.yaml)
- ✅ Proto definitions changed

**Don't rebuild if:**
- ❌ Only Go service code changed
- ❌ Only IDF repository code changed
- ❌ Only documentation changed

---

## Quick Checklist

- [ ] Step 1: Rebuild prism-service JAR (if YAML changed)
- [ ] Step 2: Run IDF script on PC
- [ ] Step 3: Copy artifacts to PC (binary + JAR)
- [ ] Step 4: Update Adonis JAR on PC
- [ ] Step 5: Start services on PC
- [ ] Step 6: Verify API works

---

## Troubleshooting

### Prism-Service JAR Missing Controller
- Rebuild: `cd ~/ntnx-api-prism-service && mvn clean install -DskipTests`
- Verify: `jar -tf target/prism-service-*.jar | grep ItemController`

### Adonis Not Loading New JAR
- Remove LKG-RELEASE.jar backup
- Restart Adonis: `genesis stop adonis && cluster start adonis`

### API Returns 404
- Check `lookup_cache.json` has `/nexus/v4.1/config` entry
- Verify Adonis is running
- Check Mercury config

---

**Last Updated:** November 27, 2025
