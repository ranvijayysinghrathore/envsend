# EnvSend Local Setup Script
# This script sets up EnvSend for local development

param(
    [switch]$SkipDocker = $false
)

$ErrorActionPreference = "Stop"

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  EnvSend Local Setup" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Step 1: Download Go dependencies
Write-Host "[Step 1/6] Downloading Go dependencies..." -ForegroundColor Yellow
try {
    go mod download
    go mod verify
    Write-Host "  ✓ Dependencies downloaded" -ForegroundColor Green
} catch {
    Write-Host "  ✗ Failed to download dependencies" -ForegroundColor Red
    Write-Host "  Error: $_" -ForegroundColor Red
    exit 1
}

Write-Host ""

# Step 2: Create .env file
Write-Host "[Step 2/6] Creating environment file..." -ForegroundColor Yellow
if (-not (Test-Path ".env")) {
    Copy-Item ".env.example" ".env"
    Write-Host "  ✓ Created .env file" -ForegroundColor Green
} else {
    Write-Host "  ⚠ .env file already exists, skipping" -ForegroundColor Yellow
}

Write-Host ""

# Step 3: Start Docker services
if (-not $SkipDocker) {
    Write-Host "[Step 3/6] Starting Docker services..." -ForegroundColor Yellow
    try {
        docker-compose up -d
        Write-Host "  ✓ Docker services started" -ForegroundColor Green
        Write-Host "  Waiting for services to be ready..." -ForegroundColor Yellow
        Start-Sleep -Seconds 30
        Write-Host "  ✓ Services should be ready" -ForegroundColor Green
    } catch {
        Write-Host "  ✗ Failed to start Docker services" -ForegroundColor Red
        Write-Host "  Error: $_" -ForegroundColor Red
        Write-Host "  Make sure Docker Desktop is running!" -ForegroundColor Yellow
        exit 1
    }
} else {
    Write-Host "[Step 3/6] Skipping Docker services (--SkipDocker flag)" -ForegroundColor Yellow
}

Write-Host ""

# Step 4: Install migration tool
Write-Host "[Step 4/6] Installing database migration tool..." -ForegroundColor Yellow
try {
    $env:GOBIN = "$PWD\bin"
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
    Write-Host "  ✓ Migration tool installed" -ForegroundColor Green
} catch {
    Write-Host "  ⚠ Failed to install migration tool (you can install manually later)" -ForegroundColor Yellow
}

Write-Host ""

# Step 5: Run migrations
if (-not $SkipDocker) {
    Write-Host "[Step 5/6] Running database migrations..." -ForegroundColor Yellow
    try {
        $migratePath = ".\bin\migrate.exe"
        if (Test-Path $migratePath) {
            & $migratePath -path migrations -database "postgres://envsend:envsend@localhost:5432/envsend?sslmode=disable" up
            Write-Host "  ✓ Migrations completed" -ForegroundColor Green
        } else {
            Write-Host "  ⚠ Migration tool not found, skipping" -ForegroundColor Yellow
            Write-Host "  You can run migrations manually later" -ForegroundColor Cyan
        }
    } catch {
        Write-Host "  ⚠ Failed to run migrations (database might not be ready yet)" -ForegroundColor Yellow
        Write-Host "  You can run migrations manually: .\bin\migrate.exe -path migrations -database 'postgres://envsend:envsend@localhost:5432/envsend?sslmode=disable' up" -ForegroundColor Cyan
    }
} else {
    Write-Host "[Step 5/6] Skipping migrations (--SkipDocker flag)" -ForegroundColor Yellow
}

Write-Host ""

# Step 6: Build binaries
Write-Host "[Step 6/6] Building EnvSend binaries..." -ForegroundColor Yellow
try {
    # Create bin directory
    if (-not (Test-Path "bin")) {
        New-Item -ItemType Directory -Force -Path "bin" | Out-Null
    }
    
    # Build CLI
    Write-Host "  Building CLI..." -ForegroundColor Cyan
    go build -o bin\envsend.exe .\cli
    
    # Build API server
    Write-Host "  Building API server..." -ForegroundColor Cyan
    go build -o bin\api-server.exe .\backend\cmd\server
    
    # Build worker
    Write-Host "  Building worker..." -ForegroundColor Cyan
    go build -o bin\worker.exe .\backend\cmd\worker
    
    Write-Host "  ✓ All binaries built successfully" -ForegroundColor Green
} catch {
    Write-Host "  ✗ Failed to build binaries" -ForegroundColor Red
    Write-Host "  Error: $_" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  Setup Complete!" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

Write-Host "Next steps:" -ForegroundColor Yellow
Write-Host ""
Write-Host "1. Start the API server:" -ForegroundColor Cyan
Write-Host "   .\bin\api-server.exe" -ForegroundColor White
Write-Host ""
Write-Host "2. In a new terminal, test the CLI:" -ForegroundColor Cyan
Write-Host "   .\bin\envsend.exe --help" -ForegroundColor White
Write-Host ""
Write-Host "3. Create a test secret:" -ForegroundColor Cyan
Write-Host "   echo 'API_KEY=test123' > test.env" -ForegroundColor White
Write-Host "   .\bin\envsend.exe send test.env --server http://localhost:8080" -ForegroundColor White
Write-Host ""
Write-Host "For more details, see SETUP_GUIDE.md" -ForegroundColor Yellow
Write-Host ""
