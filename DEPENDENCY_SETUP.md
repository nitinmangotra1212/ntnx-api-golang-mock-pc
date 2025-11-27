# Setting Up Go Module Dependencies

## Problem
Go modules can't access private repositories (`github.com/nutanix-core/*`) because:
1. They require authentication
2. Go tries to verify checksums from public database (which fails for private repos)

## Solution: Use Personal Access Token

### Step 1: Create GitHub Personal Access Token

1. Go to: https://github.com/settings/tokens
2. Click "Generate new token" → "Generate new token (classic)"
3. Give it a name (e.g., "Go Modules Access")
4. Select scope: **`repo`** (full control of private repositories)
5. Click "Generate token"
6. **Copy the token** (you won't see it again!)

### Step 2: Set Environment Variables

Add to your `~/.zshrc` or `~/.bashrc`:

```bash
# GitHub Personal Access Token
export GITHUB_TOKEN=your_token_here

# Tell Go to skip checksum verification for private repos
export GOPRIVATE=github.com/nutanix-core/*
export GONOSUMDB=github.com/nutanix-core/*
```

Then reload your shell:
```bash
source ~/.zshrc  # or source ~/.bashrc
```

### Step 3: Configure Git to Use Token

```bash
# Remove any existing SSH rewrite (if you tried that)
git config --global --unset url."git@github.com:nutanix-core/".insteadOf

# Set up HTTPS with token
git config --global url."https://${GITHUB_TOKEN}@github.com/".insteadOf "https://github.com/"

# Verify it worked
git config --global --get url."https://${GITHUB_TOKEN}@github.com/".insteadOf
# Should output: https://github.com/
```

### Step 4: Download Dependencies

```bash
cd ~/ntnx-api-golang-nexus

# Download dependencies
go mod download

# Update go.mod and go.sum
go mod tidy
```

**Expected:** Should complete without errors.

---

## Alternative: Use SSH (If You Have SSH Keys Set Up)

If you already have SSH keys configured for GitHub:

### Step 1: Test SSH Access

```bash
ssh -T git@github.com
# Should say: "Hi username! You've successfully authenticated..."
```

If you get "Permission denied", you need to set up SSH keys first:
- https://docs.github.com/en/authentication/connecting-to-github-with-ssh

### Step 2: Configure Git

```bash
git config --global url."git@github.com:nutanix-core/".insteadOf "https://github.com/nutanix-core/"
```

### Step 3: Set GOPRIVATE

```bash
export GOPRIVATE=github.com/nutanix-core/*
export GONOSUMDB=github.com/nutanix-core/*
```

Add to `~/.zshrc` or `~/.bashrc` to make permanent.

### Step 4: Download Dependencies

```bash
cd ~/ntnx-api-golang-nexus
go mod download
go mod tidy
```

---

## Troubleshooting

### Issue: "fatal: could not read Username"
**Solution:** Make sure `GITHUB_TOKEN` is set and Git config is correct.

### Issue: "404 Not Found" from sum.golang.org
**Solution:** Set `GOPRIVATE` and `GONOSUMDB` environment variables.

### Issue: "Permission denied (publickey)"
**Solution:** Either set up SSH keys OR use Personal Access Token method.

### Issue: Still using HTTPS after Git config
**Solution:** 
1. Verify Git config: `git config --global --get-regexp url`
2. Make sure `GOPRIVATE` is set
3. Try: `go clean -modcache` then `go mod download`

---

## Quick Reference

```bash
# Personal Access Token method (recommended)
export GITHUB_TOKEN=your_token
export GOPRIVATE=github.com/nutanix-core/*
export GONOSUMDB=github.com/nutanix-core/*
git config --global url."https://${GITHUB_TOKEN}@github.com/".insteadOf "https://github.com/"
go mod download
```

---

**Last Updated:** November 27, 2025

