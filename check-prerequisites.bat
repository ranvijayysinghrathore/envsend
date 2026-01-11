@echo off
REM Simple batch file to check prerequisites

echo ========================================
echo   EnvSend Setup Checker
echo ========================================
echo.

REM Check Go
echo [1/3] Checking Go installation...
where go >nul 2>&1
if %ERRORLEVEL% EQU 0 (
    echo   [OK] Go is installed
    go version
) else (
    echo   [X] Go is NOT installed
    echo   Download from: https://go.dev/dl/
    echo   Get the Windows installer (.msi file^)
)
echo.

REM Check Docker
echo [2/3] Checking Docker installation...
where docker >nul 2>&1
if %ERRORLEVEL% EQU 0 (
    echo   [OK] Docker is installed
    docker --version
    
    REM Check if Docker is running
    docker ps >nul 2>&1
    if %ERRORLEVEL% EQU 0 (
        echo   [OK] Docker is running
    ) else (
        echo   [!] Docker is installed but not running
        echo   Start Docker Desktop from the Start menu
    )
) else (
    echo   [X] Docker is NOT installed
    echo   Download from: https://www.docker.com/products/docker-desktop/
)
echo.

REM Check Make (optional)
echo [3/3] Checking Make installation (optional^)...
where make >nul 2>&1
if %ERRORLEVEL% EQU 0 (
    echo   [OK] Make is installed
) else (
    echo   [!] Make is NOT installed (optional^)
    echo   You can run commands manually
)
echo.

echo ========================================
echo   Next Steps
echo ========================================
echo.

REM Determine what to do next
where go >nul 2>&1
set GO_INSTALLED=%ERRORLEVEL%

where docker >nul 2>&1
set DOCKER_INSTALLED=%ERRORLEVEL%

if %GO_INSTALLED% NEQ 0 (
    echo You need to install:
    echo   1. Go: https://go.dev/dl/
    if %DOCKER_INSTALLED% NEQ 0 (
        echo   2. Docker Desktop: https://www.docker.com/products/docker-desktop/
    )
    echo.
    echo After installation:
    echo   1. Restart this terminal
    echo   2. Run this script again: check-prerequisites.bat
) else if %DOCKER_INSTALLED% NEQ 0 (
    echo You need to install:
    echo   1. Docker Desktop: https://www.docker.com/products/docker-desktop/
    echo.
    echo After installation:
    echo   1. Restart this terminal
    echo   2. Run this script again: check-prerequisites.bat
) else (
    echo [OK] All prerequisites are installed!
    echo.
    echo Ready to set up EnvSend. Run:
    echo   setup-local.bat
)

echo.
pause
