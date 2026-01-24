[Setup]
AppName=EnvSend
AppVersion=1.0.0
AppPublisher=EnvSend
AppPublisherURL=https://envsend.io
DefaultDirName={autopf}\EnvSend
DefaultGroupName=EnvSend
OutputBaseFilename=EnvSendSetup
Compression=lzma2
SolidCompression=yes
ChangesEnvironment=yes
DisableProgramGroupPage=yes
WizardStyle=modern
UninstallDisplayIcon={app}\envsend.exe
SetupIconFile=compiler:SetupClassicIcon.ico

[Run]
Filename: "powershell.exe"; Parameters: "-NoExit -Command ""Write-Host 'Welcome to EnvSend! Type envsend --help to get started.' -ForegroundColor Cyan; envsend --help"""; Description: "Launch EnvSend (Open Terminal)"; Flags: postinstall nowait skipifsilent

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"

[Files]
; IMPORTANT: Run "go build -o bin/envsend.exe" before compiling this script!
Source: "..\..\bin\envsend.exe"; DestDir: "{app}"; Flags: ignoreversion

[Tasks]
Name: "envPath"; Description: "Add EnvSend to system PATH variable"; GroupDescription: "Additional icons:"; Flags: unchecked

[Registry]
Root: HKLM; Subkey: "SYSTEM\CurrentControlSet\Control\Session Manager\Environment"; \
    ValueType: expandsz; ValueName: "Path"; ValueData: "{olddata};{app}"; \
    Check: NeedsAddPath('{app}')

[Code]
function NeedsAddPath(Param: string): boolean;
var
  OrigPath: string;
begin
  if not RegQueryStringValue(HKEY_LOCAL_MACHINE,
    'SYSTEM\CurrentControlSet\Control\Session Manager\Environment',
    'Path', OrigPath)
  then begin
    Result := True;
    exit;
  end;
  // look for the path with leading and trailing semicolon
  // Pos() returns 0 if not found
  Result := Pos(';' + Param + ';', ';' + OrigPath + ';') = 0;
end;
