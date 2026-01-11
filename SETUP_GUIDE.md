# EnvSend Setup Guide for Windows

This guide will help you set up EnvSend from scratch on Windows, even if you're new to Docker and development tools.

## Step 1: Install Prerequisites

### 1.1 Install Go (Programming Language)

1. Download Go from: https://go.dev/dl/
2. Download the Windows installer (e.g., `go1.21.6.windows-amd64.msi`)
3. Run the installer and follow the prompts
4. Verify installation:
   ```powershell
   go version
   ```
   You should see something like: `go version go1.21.6 windows/amd64`

### 1.2 Install Docker Desktop (for PostgreSQL, Redis, MinIO)

1. Download Docker Desktop from: https://www.docker.com/products/docker-desktop/
2. Run the installer
3. **Important**: During installation, enable WSL 2 (Windows Subsystem for Linux) if prompted
4. Restart your computer if required
5. Start Docker Desktop from the Start menu
6. Verify installation:
   ```powershell
   docker --version
   docker-compose --version
   ```

### 1.3 Install Make (Build Tool) - Optional but Recommended

**Option A: Using Chocolatey (Recommended)**
1. Install Chocolatey (package manager for Windows):
   - Open PowerShell as Administrator
   - Run:
     ```powershell
     Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))
     ```
2. Install Make:
   ```powershell
   choco install make
   ```

**Option B: Without Make**
You can run commands manually instead of using `make`. I'll provide both options.

## Step 2: Set Up the Project

### 2.1 Navigate to Project Directory

```powershell
cd "d:\Programs\Next js {projects}\envsend"
```

### 2.2 Initialize Go Modules

```powershell
# Download all dependencies
go mod download

# Verify dependencies
go mod verify
```

This will download all required Go packages (might take a few minutes).

## Step 3: Start Infrastructure Services

### 3.1 Create Environment File

Copy the example environment file:
```powershell
copy .env.example .env
```

The default settings are already configured for local development.

### 3.2 Start Docker Services

```powershell
# Start PostgreSQL, Redis, and MinIO
docker-compose up -d

# Wait for services to be ready (about 30 seconds)
Start-Sleep -Seconds 30

# Check if services are running
docker-compose ps
```

You should see 3 services running:
- `envsend-postgres`
- `envsend-redis`
- `envsend-minio`

### 3.3 Verify Services

**PostgreSQL**:
```powershell
docker exec envsend-postgres psql -U envsend -c "SELECT version();"
```

**Redis**:
```powershell
docker exec envsend-redis redis-cli ping
```
Should return: `PONG`

**MinIO**:
Open browser to: http://localhost:9001
- Username: `minioadmin`
- Password: `minioadmin`

## Step 4: Set Up Database

### 4.1 Install Migration Tool

```powershell
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### 4.2 Run Database Migrations

```powershell
# Run migrations
migrate -path migrations -database "postgres://envsend:envsend@localhost:5432/envsend?sslmode=disable" up

# Verify tables were created
docker exec envsend-postgres psql -U envsend -d envsend -c "\dt"
```

You should see tables: `secrets`, `audit_logs`, etc.

## Step 5: Build the Project

### 5.1 Build All Components

**With Make**:
```powershell
make build
```

**Without Make** (manual build):
```powershell
# Create bin directory
New-Item -ItemType Directory -Force -Path bin

# Build CLI
go build -o bin/envsend.exe ./cli

# Build API server
go build -o bin/api-server.exe ./backend/cmd/server

# Build worker
go build -o bin/worker.exe ./backend/cmd/worker
```

### 5.2 Verify Builds

```powershell
# Check CLI
.\bin\envsend.exe --version

# Check API server
.\bin\api-server.exe --help

# Check worker
.\bin\worker.exe --help
```

## Step 6: Run the System

### 6.1 Start API Server

Open a new PowerShell terminal:
```powershell
cd "d:\Programs\Next js {projects}\envsend"
.\bin\api-server.exe
```

You should see:
```
Starting EnvSend API server in development mode...
✓ Connected to PostgreSQL
✓ Connected to Redis
✓ Connected to MinIO/S3
✓ Initialized services
✓ Configured router
🚀 Server listening on 0.0.0.0:8080
```

### 6.2 Start Worker (Optional)

Open another PowerShell terminal:
```powershell
cd "d:\Programs\Next js {projects}\envsend"
.\bin\worker.exe
```

You should see:
```
Starting EnvSend cleanup worker...
✓ Connected to PostgreSQL
✓ Connected to MinIO/S3
✓ Connected to Redis
🧹 Worker started (cleanup interval: 1m0s, batch size: 100)
```

## Step 7: Test the System

### 7.1 Test Health Check

Open a new PowerShell terminal:
```powershell
# Test API health
curl http://localhost:8080/health
```

Should return: `{"status":"healthy"}`

### 7.2 Create a Test Secret

```powershell
# Create a test .env file
echo "API_KEY=test123`nDATABASE_URL=postgres://localhost" > test.env

# Send the secret
.\bin\envsend.exe send test.env --server http://localhost:8080
```

You should see output like:
```
✅ Secret uploaded successfully!
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🔗 Share this link: /s/550e8400-e29b-41d4-a716-446655440000#<key>
⚠️  The key is in the URL fragment (after #) - it's never sent to the server
⏰ Expires: 2024-01-15T10:40:00+05:30 (10m)
👁️  Max views: 1
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

### 7.3 Retrieve the Secret

```powershell
# Copy the full URL from above and retrieve it
.\bin\envsend.exe receive "http://localhost:8080/s/<secret-id>#<key>" > retrieved.env

# Verify the content
type retrieved.env
```

Should show:
```
API_KEY=test123
DATABASE_URL=postgres://localhost
```

## Step 8: Test Different Encryption Modes

### 8.1 Passphrase-Protected Secret

```powershell
.\bin\envsend.exe send test.env --require-passphrase --server http://localhost:8080
```

You'll be prompted to enter a passphrase twice.

### 8.2 Custom Expiry and Views

```powershell
# Expires in 1 hour, max 3 views
.\bin\envsend.exe send test.env --expires 1h --max-views 3 --server http://localhost:8080
```

### 8.3 Pipe Support

```powershell
# Send from pipe
type test.env | .\bin\envsend.exe send --server http://localhost:8080
```

## Troubleshooting

### Issue: "go: command not found"
**Solution**: Restart PowerShell after installing Go, or add Go to PATH manually.

### Issue: "docker: command not found"
**Solution**: 
1. Make sure Docker Desktop is running
2. Restart PowerShell
3. Check Docker Desktop settings

### Issue: "Cannot connect to database"
**Solution**:
```powershell
# Check if PostgreSQL is running
docker ps | findstr postgres

# Restart if needed
docker-compose restart postgres
```

### Issue: "Port already in use"
**Solution**:
```powershell
# Find what's using port 8080
netstat -ano | findstr :8080

# Kill the process (replace PID with actual process ID)
taskkill /PID <PID> /F

# Or change the port in .env file
# Set SERVER_PORT=8081
```

### Issue: Migration fails
**Solution**:
```powershell
# Reset database
docker-compose down -v
docker-compose up -d
Start-Sleep -Seconds 30

# Run migrations again
migrate -path migrations -database "postgres://envsend:envsend@localhost:5432/envsend?sslmode=disable" up
```

## Next Steps

### View Logs
```powershell
# View Docker logs
docker-compose logs -f postgres
docker-compose logs -f redis
docker-compose logs -f minio
```

### Stop Services
```powershell
# Stop API server and worker: Press Ctrl+C in their terminals

# Stop Docker services
docker-compose down

# Stop and remove all data (reset everything)
docker-compose down -v
```

### Access MinIO Console
1. Open browser: http://localhost:9001
2. Login: `minioadmin` / `minioadmin`
3. View encrypted blobs in `envsend-secrets` bucket

### View Database
```powershell
# Connect to PostgreSQL
docker exec -it envsend-postgres psql -U envsend -d envsend

# View secrets
SELECT id, expires_at, max_views, view_count, destroyed FROM secrets;

# View audit logs
SELECT secret_id, action, timestamp FROM audit_logs ORDER BY timestamp DESC LIMIT 10;

# Exit
\q
```

## Summary

You now have a fully functional EnvSend system running locally! 

**What's running:**
- ✅ PostgreSQL (database for metadata)
- ✅ Redis (rate limiting)
- ✅ MinIO (encrypted blob storage)
- ✅ API Server (REST API on port 8080)
- ✅ Worker (background cleanup)
- ✅ CLI (command-line tool)

**You can:**
- Send secrets with `envsend send`
- Receive secrets with `envsend receive`
- Use different encryption modes
- View audit logs in the database
- See encrypted blobs in MinIO

For production deployment to cloud/Kubernetes, see `docs/DEPLOYMENT.md`.
