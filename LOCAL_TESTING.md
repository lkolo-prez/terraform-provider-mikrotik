# Local Testing Guide

This guide shows how to test the Terraform provider locally without waiting for GitHub Actions.

---

## 🚀 Quick Start

### Prerequisites

**Windows (PowerShell):**
- Docker Desktop installed and running
- Go 1.21+ installed (optional for build only)

**Linux/Mac:**
- Docker installed
- Go 1.21+ installed (optional)

---

## 📋 Testing Scenarios

### 1. **Build Test Only** (No RouterOS needed - 30 seconds)

Tests that the code compiles without errors.

#### Windows (PowerShell):
```powershell
# Navigate to project
cd "c:\Users\kolod\Desktop\LKP\00_AIRCLOUD\mikrotik-provider\terraform-provider-mikrotik"

# Method 1: Using Make (if available)
make build

# Method 2: Direct Go build
go build -v .

# Check if binary was created
ls terraform-provider-mikrotik.exe
```

#### Linux/Mac:
```bash
cd ~/path/to/terraform-provider-mikrotik

# Build
make build
# or
go build -v .

# Check binary
ls -lh terraform-provider-mikrotik
```

**Expected Output:**
```
✅ terraform-provider-mikrotik.exe created (or .so on Linux)
```

---

### 2. **Lint Test** (No RouterOS needed - 1 minute)

Tests code quality and style.

#### Windows (PowerShell):
```powershell
# Run linters
go vet ./client/...
go vet ./mikrotik/...

# Or using Make
make lint
```

#### Linux/Mac:
```bash
make lint
```

**Expected Output:**
```
✅ No issues found
```

---

### 3. **Unit Tests** (No RouterOS needed - 2-3 minutes)

Tests individual functions without RouterOS.

#### Windows (PowerShell):
```powershell
# Test client package only
cd client
go test -v -short ./...

# Go back
cd ..
```

#### Linux/Mac:
```bash
cd client
go test -v -short ./...
cd ..
```

**Expected Output:**
```
✅ PASS
```

---

### 4. **Full Integration Test** (RouterOS required - 10-30 minutes)

Complete test with real RouterOS container.

#### Step 1: Start RouterOS Container

**Windows (PowerShell):**
```powershell
# Start RouterOS 7.16.2 (stable)
docker run -d --name routeros `
  --cap-add=NET_ADMIN `
  --device=/dev/net/tun `
  -p 8728:8728 `
  -p 8080:80 `
  mnazarenko/docker-routeros:7.16.2

# Wait for RouterOS to start (2-3 minutes)
Write-Host "Waiting for RouterOS to start..."
for ($i = 1; $i -le 60; $i++) {
    try {
        $response = Invoke-WebRequest -Uri "http://127.0.0.1:8080" -TimeoutSec 1 -ErrorAction SilentlyContinue
        Write-Host "✅ RouterOS is ready!"
        break
    } catch {
        Write-Host -NoNewline "."
        Start-Sleep -Seconds 2
    }
}

# Check container status
docker ps | Select-String routeros
```

**Linux/Mac:**
```bash
# Using Make (recommended)
make routeros ROUTEROS_VERSION=7.16.2

# Or manually
docker run -d --name routeros \
  --cap-add=NET_ADMIN \
  --device=/dev/net/tun \
  -p 8728:8728 \
  -p 8080:80 \
  mnazarenko/docker-routeros:7.16.2

# Wait for RouterOS
echo "Waiting for RouterOS..."
for i in {1..60}; do
    if curl -s http://127.0.0.1:8080 > /dev/null 2>&1; then
        echo "✅ RouterOS is ready!"
        break
    fi
    echo -n "."
    sleep 2
done

# Check status
docker ps | grep routeros
```

#### Step 2: Set Environment Variables

**Windows (PowerShell):**
```powershell
$env:MIKROTIK_HOST = "127.0.0.1:8728"
$env:MIKROTIK_USER = "admin"
$env:MIKROTIK_PASSWORD = ""
$env:TF_ACC = "1"
```

**Linux/Mac:**
```bash
export MIKROTIK_HOST="127.0.0.1:8728"
export MIKROTIK_USER="admin"
export MIKROTIK_PASSWORD=""
export TF_ACC=1
```

#### Step 3: Run Tests

**Windows (PowerShell):**
```powershell
# Test client library
cd client
go test -v -timeout 30m ./...
cd ..

# Test provider (longer)
go test -v -timeout 40m ./...

# Or specific test
go test -v -run TestAccBridge ./...
```

**Linux/Mac:**
```bash
# Using Make
make testclient

# Or manually
cd client
go test -v -timeout 30m ./...
cd ..

# Full acceptance tests
make testacc

# Specific test
go test -v -run TestAccBridge ./...
```

#### Step 4: View Logs (if tests fail)

**Windows (PowerShell):**
```powershell
docker logs routeros
```

**Linux/Mac:**
```bash
# Follow logs
make routeros-logs

# Or view directly
docker logs routeros
```

#### Step 5: Cleanup

**Windows (PowerShell):**
```powershell
docker stop routeros
docker rm routeros
```

**Linux/Mac:**
```bash
make routeros-clean
# Or
docker stop routeros && docker rm routeros
```

---

## 🎯 Quick Test Workflows

### Workflow 1: **Fast Feedback** (3 minutes - no RouterOS)

```powershell
# Windows PowerShell
go build -v .
go vet ./client/...
go vet ./mikrotik/...
cd client ; go test -v -short ./... ; cd ..
```

```bash
# Linux/Mac
make build
make lint
cd client && go test -v -short ./... && cd ..
```

**Use when:** Quick syntax check, before commit

---

### Workflow 2: **RouterOS Integration** (15 minutes)

**Windows (PowerShell):**
```powershell
# 1. Start RouterOS
docker run -d --name routeros --cap-add=NET_ADMIN --device=/dev/net/tun -p 8728:8728 -p 8080:80 mnazarenko/docker-routeros:7.16.2

# 2. Wait
Start-Sleep -Seconds 120

# 3. Set env
$env:MIKROTIK_HOST = "127.0.0.1:8728"
$env:MIKROTIK_USER = "admin"
$env:MIKROTIK_PASSWORD = ""
$env:TF_ACC = "1"

# 4. Test
cd client ; go test -v -timeout 30m ./... ; cd ..

# 5. Cleanup
docker stop routeros ; docker rm routeros
```

**Linux/Mac:**
```bash
# All in one
make routeros ROUTEROS_VERSION=7.16.2 && \
export MIKROTIK_HOST=127.0.0.1:8728 MIKROTIK_USER=admin MIKROTIK_PASSWORD="" TF_ACC=1 && \
sleep 120 && \
make testclient && \
make routeros-clean
```

**Use when:** Testing RouterOS integration, before PR

---

## 🐛 Troubleshooting

### Problem: "docker: command not found"

**Solution:**
- Install Docker Desktop (Windows/Mac)
- Install Docker Engine (Linux)
- Start Docker service

### Problem: RouterOS container won't start

**Windows:**
```powershell
# Check if WSL2 is enabled
wsl --status

# Enable virtualization in BIOS if needed
```

**Solution:**
- Enable WSL2 backend in Docker Desktop
- Enable virtualization in BIOS
- Try different RouterOS version: `7.14.3` or `7.17`

### Problem: "dial tcp 127.0.0.1:8728: connect: connection refused"

**Solution:**
```powershell
# RouterOS not ready yet, wait longer
Start-Sleep -Seconds 60

# Check if container is running
docker ps

# Check logs
docker logs routeros

# Restart container
docker restart routeros
```

### Problem: Tests timeout

**Solution:**
```powershell
# Increase timeout
go test -v -timeout 60m ./...

# Or set in environment
$env:TIMEOUT = "60m"
make testacc
```

### Problem: "go: module requires Go 1.21"

**Solution:**
```powershell
# Check Go version
go version

# Should show: go version go1.21 or higher
# If not, download from: https://go.dev/dl/
```

---

## 📊 Test Matrix

| Test Type | Time | RouterOS? | Command |
|-----------|------|-----------|---------|
| Build | 30s | ❌ No | `make build` |
| Lint | 1m | ❌ No | `make lint` |
| Unit Tests | 3m | ❌ No | `cd client && go test -v -short ./...` |
| Client Integration | 15m | ✅ Yes | `make testclient` |
| Full Acceptance | 30m+ | ✅ Yes | `make testacc` |

---

## 🚀 CI/CD Simulation

Simulate exactly what GitHub Actions does:

**Windows (PowerShell):**
```powershell
# Phase 1: Build & Lint (mimics CI)
Write-Host "=== Phase 1: Build & Lint ===" -ForegroundColor Cyan
go mod download
go mod verify
go build -v .
go vet ./client/...
go vet ./mikrotik/...
cd client ; go test -v -short ./... ; cd ..

Write-Host "✅ Phase 1 Complete!" -ForegroundColor Green

# Phase 2: Integration (optional, mimics CI)
Write-Host "=== Phase 2: Integration Test ===" -ForegroundColor Cyan
docker run -d --name routeros --cap-add=NET_ADMIN --device=/dev/net/tun -p 8728:8728 -p 8080:80 mnazarenko/docker-routeros:7.16.2
Start-Sleep -Seconds 120

$env:MIKROTIK_HOST = "127.0.0.1:8728"
$env:MIKROTIK_USER = "admin"
$env:MIKROTIK_PASSWORD = ""
$env:TF_ACC = "1"

cd client ; go test -v -timeout 30m ./... ; cd ..

docker stop routeros
docker rm routeros

Write-Host "✅ Phase 2 Complete!" -ForegroundColor Green
```

**Linux/Mac:**
```bash
#!/bin/bash
set -e

echo "=== Phase 1: Build & Lint ==="
go mod download
go mod verify
go build -v .
go vet ./client/...
go vet ./mikrotik/...
cd client && go test -v -short ./... && cd ..
echo "✅ Phase 1 Complete!"

echo "=== Phase 2: Integration Test ==="
make routeros ROUTEROS_VERSION=7.16.2
sleep 120

export MIKROTIK_HOST=127.0.0.1:8728
export MIKROTIK_USER=admin
export MIKROTIK_PASSWORD=""
export TF_ACC=1

cd client && go test -v -timeout 30m ./... && cd ..
make routeros-clean
echo "✅ Phase 2 Complete!"
```

---

## 📝 Pro Tips

### 1. **Test Specific Resources**

```bash
# Only BGP tests
go test -v -run TestAccBgp ./...

# Only Bridge tests
go test -v -run TestAccBridge ./...

# Single test
go test -v -run TestAccBridgePort_basic ./...
```

### 2. **Keep RouterOS Running**

```bash
# Don't cleanup between test runs
docker ps  # Check it's running
# Just re-run tests without restarting container
```

### 3. **Use Different RouterOS Versions**

```bash
# Test against 7.14.3
make routeros ROUTEROS_VERSION=7.14.3

# Test against latest
make routeros ROUTEROS_VERSION=latest
```

### 4. **Debug Mode**

```bash
# Enable Terraform logs
export TF_LOG=DEBUG
export TF_LOG_PATH=./terraform.log

# Run test
go test -v -run TestAccBridge ./...

# Check logs
cat terraform.log
```

---

## ✅ Quick Checklist Before Commit

- [ ] `go build -v .` - passes
- [ ] `make lint` - no warnings
- [ ] `cd client && go test -v -short ./...` - passes
- [ ] Code formatted: `go fmt ./...`
- [ ] Documentation updated
- [ ] Examples work

---

## 🔗 Related Documentation

- [Makefile](./Makefile) - All available commands
- [DEVELOPMENT_ROADMAP.md](./DEVELOPMENT_ROADMAP.md) - Development phases
- [ROUTEROS7_SUPPORT.md](./ROUTEROS7_SUPPORT.md) - RouterOS 7 features
- [examples/](./examples/) - Usage examples

---

**Last Updated**: November 25, 2025  
**Tested With**: Docker Desktop 4.x, Go 1.21-1.23, RouterOS 7.14.3-7.17
