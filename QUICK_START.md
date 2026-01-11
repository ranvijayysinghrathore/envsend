# 🚀 Quick Start - Install Prerequisites

You need to install 2 things before using EnvSend:

## 1️⃣ Install Go (Required)

**What is it?** The programming language EnvSend is written in.

**Steps:**
1. Go to: **https://go.dev/dl/**
2. Download: **go1.21.x.windows-amd64.msi** (the Windows installer)
3. Run the installer (just click "Next" through everything)
4. **Restart PowerShell** after installation

**Verify:**
```powershell
go version
```
Should show: `go version go1.21.x windows/amd64`

---

## 2️⃣ Install Docker Desktop (Required)

**What is it?** Runs PostgreSQL, Redis, and MinIO (the databases and storage).

**Steps:**
1. Go to: **https://www.docker.com/products/docker-desktop/**
2. Click **"Download for Windows"**
3. Run the installer
4. **Important**: Enable WSL 2 if prompted
5. **Restart your computer** if required
6. **Start Docker Desktop** from the Start menu (wait for it to fully start)

**Verify:**
```powershell
docker --version
docker ps
```
Should show Docker version and running containers.

---

## 3️⃣ After Installing Both

**Run the setup script:**
```powershell
cd "d:\Programs\Next js {projects}\envsend"
.\setup-local.ps1
```

This will:
- ✅ Download all Go dependencies
- ✅ Start PostgreSQL, Redis, and MinIO
- ✅ Create database tables
- ✅ Build the CLI and servers

**Then start the system:**
```powershell
# Terminal 1: Start API server
.\bin\api-server.exe

# Terminal 2: Test it
.\bin\envsend.exe --help
```

---

## ❓ Need Help?

**Check what's installed:**
```powershell
.\check-prerequisites.ps1
```

**Common Issues:**

**"go: command not found"**
→ Restart PowerShell after installing Go

**"docker: command not found"**
→ Make sure Docker Desktop is running (check system tray)

**"Port already in use"**
→ Something else is using port 8080. Change it in `.env` file:
```
SERVER_PORT=8081
```

---

## 📚 Full Documentation

- **SETUP_GUIDE.md** - Detailed step-by-step guide
- **README.md** - Project overview
- **docs/ARCHITECTURE.md** - How it works

---

## ⚡ Quick Test After Setup

```powershell
# Create a test secret
echo "API_KEY=test123" > test.env
.\bin\envsend.exe send test.env --server http://localhost:8080

# You'll get a URL like:
# http://localhost:8080/s/abc123#key

# Retrieve it
.\bin\envsend.exe receive "http://localhost:8080/s/abc123#key" > retrieved.env

# Check it worked
type retrieved.env
```

**That's it! 🎉**
