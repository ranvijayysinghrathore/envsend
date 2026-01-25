# EnvSend Windows Installer Builder
# This script builds the CLI and creates a Windows installer using Inno Setup

$ErrorActionPreference = "Stop"

Write-Host "Building EnvSend Windows Installer..." -ForegroundColor Cyan
Write-Host ""

# Configuration
$ServerURL = "https://envsend.onrender.com"
$Version = "1.0.0"

# Step 1: Build CLI
Write-Host "Building CLI with production URL..." -ForegroundColor Green
Write-Host "   Server: $ServerURL" -ForegroundColor Gray

Set-Location cli
go build -ldflags "-s -w -X 'github.com/ranvijayysinghrathore/envsend/cli/cmd.DefaultServerURL=$ServerURL'" -o ../bin/envsend.exe .
Set-Location ..

if (-not (Test-Path "bin\envsend.exe")) {
    Write-Host "Failed to build CLI" -ForegroundColor Red
    exit 1
}

Write-Host "CLI built successfully" -ForegroundColor Green
Write-Host ""

# Step 2: Build Inno Setup Installer
Write-Host "Creating Windows Installer..." -ForegroundColor Green

# Check if Inno Setup is installed
$InnoSetup = "C:\Program Files (x86)\Inno Setup 6\ISCC.exe"
if (-not (Test-Path $InnoSetup)) {
    Write-Host "Inno Setup not found at: $InnoSetup" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "Download Inno Setup from: https://jrsoftware.org/isdl.php" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "CLI binary ready at: bin\envsend.exe" -ForegroundColor Green
    Write-Host "You can distribute this file directly or install Inno Setup to create an installer" -ForegroundColor Gray
    exit 0
}

# Compile installer
& $InnoSetup "deployments\windows\setup.iss"

if ($LASTEXITCODE -eq 0) {
    Write-Host ""
    Write-Host "Installer created successfully!" -ForegroundColor Green
    Write-Host "Location: releases\EnvSendSetup.exe" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Distribution:" -ForegroundColor Cyan
    Write-Host "   1. Share EnvSendSetup.exe with users" -ForegroundColor White
    Write-Host "   2. Users run installer (requires admin)" -ForegroundColor White
    Write-Host "   3. Installer adds to PATH automatically" -ForegroundColor White
    Write-Host "   4. Users can use: envsend .env" -ForegroundColor White
} else {
    Write-Host "Failed to create installer" -ForegroundColor Red
    exit 1
}
