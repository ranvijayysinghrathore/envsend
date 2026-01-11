# ------------------------------------------------------------------
# CONFIGURATION: Change this to your real domain when you buy it!
# Example: "https://envsend.your-startup.com"
# ------------------------------------------------------------------
$ServerURL = "http://localhost:8080" 

echo "Building EnvSend with Server URL: $ServerURL"

go build -ldflags "-X 'github.com/ranvijayysinghrathore/envsend/cli/cmd.DefaultServerURL=$ServerURL'" -o bin/envsend.exe ./cli

# Create a Setup Script
$SetupContent = @"
@echo off
echo Installing EnvSend...

REM Copy binary to Program Files
mkdir "%ProgramFiles%\EnvSend" 2>nul
copy /Y envsend.exe "%ProgramFiles%\EnvSend\envsend.exe"

REM Add to System PATH (Requires Admin)
echo Adding to PATH...
setx /M PATH "%PATH%;%ProgramFiles%\EnvSend"

echo Installation Complete!
pause
"@

$SetupContent | Out-File -Encoding ASCII bin/install.bat

# Create ZIP package
Compress-Archive -Path bin/envsend.exe, bin/install.bat -DestinationPath bin/EnvSend-Windows-Installer.zip -Force

echo "Installer created at bin/EnvSend-Windows-Installer.zip"
echo "Distribute this ZIP file to your developers."
