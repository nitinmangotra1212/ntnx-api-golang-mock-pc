# GitHub Authentication Setup for Maven Build

## Problem

When building `ntnx-api-golang-nexus-pc`, Maven fails with:
```
ERROR: git-upload-pack not permitted on 'https://github.com/nutanix-core/ntnx-api-dev-platform.git/'
```

**Root Cause:** The Maven plugin uses **JGit** (Eclipse JGit), which doesn't use Git's credential helper. Even though `git clone` works, Maven build fails because JGit can't authenticate.

---

## Solution: Use Personal Access Token

### Step 1: Create GitHub Personal Access Token

1. Go to GitHub → **Settings** → **Developer settings** → **Personal access tokens** → **Tokens (classic)**
2. Click **"Generate new token (classic)"**
3. Give it a name (e.g., "Maven Build - Nexus")
4. Select scope: **`repo`** (full control of private repositories)
5. Click **"Generate token"**
6. **Copy the token immediately** (you won't see it again!)

---

### Step 2: Update repositories.yaml

**File:** `golang-nexus-api-definitions/defs/metadata/repositories.yaml`

**Update the URI to include your token:**
```yaml
- name: "common"
  type: "git"
  uri: "https://<YOUR_TOKEN>@github.com/nutanix-core/ntnx-api-dev-platform.git"
  ref: "refs/tags/17.6.0.9581-RELEASE"
  baseDir: "ntnx-api-common/common-api-definitions/defs"
  authRequired: false
```

**Replace `<YOUR_TOKEN>` with your actual token.**

**Example:**
```yaml
- name: "common"
  type: "git"
  uri: "https://ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx@github.com/nutanix-core/ntnx-api-dev-platform.git"
  ref: "refs/tags/17.6.0.9581-RELEASE"
  baseDir: "ntnx-api-common/common-api-definitions/defs"
  authRequired: false
```

---

### Step 3: Test Build

```bash
cd ~/ntnx-api-golang-nexus-pc
mvn clean install -DskipTests -s settings.xml
```

**Expected:** Build succeeds without authentication errors.

---

## Alternative: Use SSH (If You Have SSH Keys)

### Step 1: Add SSH Key to GitHub

```bash
# Copy your public key
cat ~/.ssh/id_ed25519.pub
# Or
cat ~/.ssh/id_rsa.pub
```

1. Go to GitHub → **Settings** → **SSH and GPG keys**
2. Click **"New SSH key"**
3. Paste your public key
4. Save

### Step 2: Test SSH Connection

```bash
ssh -T git@github.com
# Should say: "Hi username! You've successfully authenticated..."
```

### Step 3: Update repositories.yaml

```yaml
- name: "common"
  type: "git"
  uri: "git@github.com:nutanix-core/ntnx-api-dev-platform.git"
  ref: "refs/tags/17.6.0.9581-RELEASE"
  baseDir: "ntnx-api-common/common-api-definitions/defs"
  authRequired: false
```

**Note:** SSH may not work with JGit. Personal Access Token (HTTPS) is more reliable.

---

## Security Best Practices

### ⚠️ Important Security Notes

1. **Token in File:** The token will be visible in `repositories.yaml`
   - Consider using `.gitignore` to exclude this file if it contains tokens
   - Or use environment variables (if supported by Maven plugin)

2. **Token Rotation:** Rotate tokens regularly (every 90 days recommended)

3. **Scope Limitation:** Use the minimum required scope (`repo` for private repos)

4. **Team Sharing:** For team environments, consider:
   - Using a shared service account token
   - Storing token in secure vault (HashiCorp Vault, AWS Secrets Manager, etc.)
   - Using CI/CD secrets management

---

## Troubleshooting

### Issue: Token Still Not Working

**Check:**
1. Token has `repo` scope
2. Token is not expired
3. Token format in URL is correct: `https://<TOKEN>@github.com/...`

**Test Token Manually:**
```bash
# Test if token works
curl -H "Authorization: token <YOUR_TOKEN>" \
  https://api.github.com/user
```

### Issue: Build Works Locally But Fails in CI/CD

**Solution:** Set up token as CI/CD secret and inject it during build.

**Example (GitHub Actions):**
```yaml
- name: Update repositories.yaml
  run: |
    sed -i "s|https://github.com|https://${{ secrets.GITHUB_TOKEN }}@github.com|g" \
      golang-nexus-api-definitions/defs/metadata/repositories.yaml
```

---

## Why This Happens

- **Regular Git** uses `~/.git-credentials` (stored when you run `git clone`)
- **JGit** (used by Maven plugin) doesn't read `~/.git-credentials` by default
- JGit needs credentials embedded in URL or configured via Java APIs
- Embedding token in URL is the simplest solution for JGit

---

## Quick Reference

**File to Modify:**
```
golang-nexus-api-definitions/defs/metadata/repositories.yaml
```

**Format:**
```yaml
uri: "https://<TOKEN>@github.com/nutanix-core/ntnx-api-dev-platform.git"
```

**Test:**
```bash
cd ~/ntnx-api-golang-nexus-pc
mvn clean install -DskipTests -s settings.xml
```

---

**Last Updated:** November 26, 2025
