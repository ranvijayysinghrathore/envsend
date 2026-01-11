@echo off
REM EnvSend Local Setup Script
REM This script sets up EnvSend for local development

echo ========================================
echo   EnvSend Local Setup
echo ========================================
echo.

REM Step 1: Download Go dependencies
echo [Step 1/6] Downloading Go dependencies...
go mod download
if %ERRORLEVEL% NEQ 0 (
    echo   [X] Failed to download dependencies
    echo   Make sure Go is installed!
    pause
    exit /b 1
)
go mod verify
echo   [OK] Dependencies downloaded
echo.

REM Step 2: Create .env file
echo [Step 2/6] Creating environment file...
if not exist ".env" (
    copy ".env.example" ".env" >nul
    echo   [OK] Created .env file
) else (
    echo   [!] .env file already exists, skipping
)
echo.

REM Step 3: Start Docker services
echo [Step 3/6] Starting Docker services...
docker-compose up -d
if %ERRORLEVEL% NEQ 0 (
    echo   [X] Failed to start Docker services
    echo   Make sure Docker Desktop is running!
    pause
    exit /b 1
)
echo   [OK] Docker services started
echo   Waiting for services to be ready...
timeout /t 30 /nobreak >nul
echo   [OK] Services should be ready
echo.

REM Step 4: Install migration tool
echo [Step 4/6] Installing database migration tool...
set GOBIN=%CD%\bin
go install -tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate@latest
if %ERRORLEVEL% EQU 0 (
    echo   [OK] Migration tool installed
) else (
    echo   [!] Failed to install migration tool (you can install manually later^)
)
echo.

REM Step 5: Run migrations
echo [Step 5/6] Running database migrations...
if exist "bin\migrate.exe" (
    bin\migrate.exe -path migrations -database "postgres://envsend:envsend@localhost:5432/envsend?sslmode=disable" up
    if %ERRORLEVEL% EQU 0 (
        echo   [OK] Migrations completed
    ) else (
        echo   [!] Failed to run migrations (database might not be ready yet^)
        echo   You can run migrations manually later
    )
) else (
    echo   [!] Migration tool not found, skipping
    echo   You can run migrations manually later
)
echo.

REM Step 6: Build binaries
echo [Step 6/6] Building EnvSend binaries...

REM Create bin directory
if not exist "bin" mkdir bin

REM Build CLI
echo   Building CLI...
go build -o bin\envsend.exe .\cli
if %ERRORLEVEL% NEQ 0 (
    echo   [X] Failed to build CLI
    pause
    exit /b 1
)

REM Build API server
echo   Building API server...
go build -o bin\api-server.exe .\backend\cmd\server
if %ERRORLEVEL% NEQ 0 (
    echo   [X] Failed to build API server
    pause
    exit /b 1
)

REM Build worker
echo   Building worker...
go build -o bin\worker.exe .\backend\cmd\worker
if %ERRORLEVEL% NEQ 0 (
    echo   [X] Failed to build worker
    pause
    exit /b 1
)

echo   [OK] All binaries built successfully
echo.

echo ========================================
echo   Setup Complete!
echo ========================================
echo.
echo Next steps:
echo.
echo 1. Start the API server:
echo    bin\api-server.exe
echo.
echo 2. In a new terminal, test the CLI:
echo    bin\envsend.exe --help
echo.
echo 3. Create a test secret:
echo    echo API_KEY=test123 ^> test.env
echo    bin\envsend.exe send test.env --server http://localhost:8080
echo.
echo For more details, see SETUP_GUIDE.md
echo.
pause
