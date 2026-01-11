@echo off
echo Starting EnvSend Local Environment...

:: Check if Docker is running (basic check)
docker info >nul 2>&1
if %errorlevel% neq 0 (
    echo [ERROR] Docker is not running! Please start Docker Desktop first.
    pause
    exit /b 1
)

:: Build latest binaries
echo [1/3] Building latest code...
go build -o bin/api-server.exe ./backend/cmd/server
go build -o bin/worker.exe ./backend/cmd/worker
if %errorlevel% neq 0 (
    echo [ERROR] Build failed! Fix errors and try again.
    pause
    exit /b 1
)

:: Start Worker in background
echo [2/3] Starting Worker (Background)...
start "EnvSend Worker" bin\worker.exe

:: Start API Server
echo [3/3] Starting API Server (Listening on :8080)...
echo.
echo    Server is READY! 
echo    - API: http://localhost:8080
echo    - Ctrl+C to stop
echo.
bin\api-server.exe
pause
