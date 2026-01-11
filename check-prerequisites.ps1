# Quick Start Script for EnvSend
# This script helps you check prerequisites and provides download links

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  EnvSend Setup Checker" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Check Go
Write-Host "[1/3] Checking Go installation..." -ForegroundColor Yellow
try {
    $goVersion = go version 2>$null
    if ($goVersion) {
        Write-Host "  ✓ Go is installed: $goVersion" -ForegroundColor Green
    }
} catch {
    Write-Host "  ✗ Go is NOT installed" -ForegroundColor Red
    Write-Host "  → Download from: https://go.dev/dl/" -ForegroundColor Cyan
    Write-Host "  → Get the Windows installer (.msi file)" -ForegroundColor Cyan
}

Write-Host ""

# Check Docker
Write-Host "[2/3] Checking Docker installation..." -ForegroundColor Yellow
try {
    $dockerVersion = docker --version 2>$null
    if ($dockerVersion) {
        Write-Host "  ✓ Docker is installed: $dockerVersion" -ForegroundColor Green
        
        # Check if Docker is running
        try {
            docker ps 2>$null | Out-Null
            Write-Host "  ✓ Docker is running" -ForegroundColor Green
        } catch {
            Write-Host "  ⚠ Docker is installed but not running" -ForegroundColor Yellow
            Write-Host "  → Start Docker Desktop from the Start menu" -ForegroundColor Cyan
        }
    }
} catch {
    Write-Host "  ✗ Docker is NOT installed" -ForegroundColor Red
    Write-Host "  → Download from: https://www.docker.com/products/docker-desktop/" -ForegroundColor Cyan
    Write-Host "  → Get Docker Desktop for Windows" -ForegroundColor Cyan
}

Write-Host ""

# Check Make (optional)
Write-Host "[3/3] Checking Make installation (optional)..." -ForegroundColor Yellow
try {
    $makeVersion = make --version 2>$null
    if ($makeVersion) {
        Write-Host "  ✓ Make is installed" -ForegroundColor Green
    }
} catch {
    Write-Host "  ⚠ Make is NOT installed (optional)" -ForegroundColor Yellow
    Write-Host "  → You can install via Chocolatey or run commands manually" -ForegroundColor Cyan
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  Next Steps" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Determine what needs to be done
$needsGo = -not (Get-Command go -ErrorAction SilentlyContinue)
$needsDocker = -not (Get-Command docker -ErrorAction SilentlyContinue)

if ($needsGo -or $needsDocker) {
    Write-Host "You need to install:" -ForegroundColor Yellow
    if ($needsGo) {
        Write-Host "  1. Go: https://go.dev/dl/" -ForegroundColor Cyan
    }
    if ($needsDocker) {
        Write-Host "  2. Docker Desktop: https://www.docker.com/products/docker-desktop/" -ForegroundColor Cyan
    }
    Write-Host ""
    Write-Host "After installation:" -ForegroundColor Yellow
    Write-Host "  1. Restart PowerShell" -ForegroundColor Cyan
    Write-Host "  2. Run this script again: .\check-prerequisites.ps1" -ForegroundColor Cyan
} else {
    Write-Host "✓ All prerequisites are installed!" -ForegroundColor Green
    Write-Host ""
    Write-Host "Ready to set up EnvSend. Run:" -ForegroundColor Yellow
    Write-Host "  .\setup-local.ps1" -ForegroundColor Cyan
}

Write-Host ""
