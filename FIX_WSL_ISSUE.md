# Fix WSL Issue and Setup EnvSend

## The Problem

Docker needs WSL 2 (Windows Subsystem for Linux) to be updated. This is showing that error.

## Quick Fix - 2 Options:

### Option 1: Update WSL (Recommended for full setup)

**Run this command in PowerShell as Administrator:**

```powershell
wsl --update
```

Then restart Docker Desktop.

---

### Option 2: Skip Docker for Now (Simplest - Test CLI Only)

You can test the CLI without Docker by building it directly:

```cmd
cd "d:\Programs\Next js {projects}\envsend"

REM Download dependencies
go mod download

REM Build just the CLI
go build -o bin\envsend.exe .\cli

REM Test it
bin\envsend.exe --help
```

**Note**: Without Docker, you won't have the server running, but you can still:
- Build the CLI
- See the code
- Test encryption/decryption locally
- Deploy to a cloud server later

---

## Full Setup After Fixing WSL

Once WSL is updated:

1. **Restart Docker Desktop**
2. **Run setup again:**
   ```cmd
   .\setup-local.bat
   ```

---

## Alternative: Use SQLite Instead of Docker (Coming Soon)

I can modify the backend to use SQLite instead of PostgreSQL, which doesn't need Docker. This would let you run everything locally without Docker.

**Would you like me to create a SQLite version?**

---

## What to Do Right Now

**Choose one:**

**A) Fix WSL and use full Docker setup:**
```powershell
# Run as Administrator
wsl --update
# Then restart Docker Desktop
# Then run: .\setup-local.bat
```

**B) Just build and test the CLI (no server):**
```cmd
go mod download
go build -o bin\envsend.exe .\cli
bin\envsend.exe --help
```

**C) Wait for me to create a SQLite version (no Docker needed)**

Let me know which option you prefer!
