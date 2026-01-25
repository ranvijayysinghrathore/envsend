[Setup]
AppName=EnvSend
AppVersion=1.0.0
AppPublisher=EnvSend
AppPublisherURL=https://envsend.io
DefaultDirName={autopf}\EnvSend
DefaultGroupName=EnvSend
OutputDir=..\..\releases
OutputBaseFilename=EnvSendSetup
Compression=lzma2
SolidCompression=yes
ChangesEnvironment=yes
DisableProgramGroupPage=yes
WizardStyle=modern
UninstallDisplayIcon={app}\envsend.exe
SetupIconFile=compiler:SetupClassicIcon.ico
PrivilegesRequired=admin

[Run]
Filename: "{cmd}"; Parameters: "/c echo Installation complete! You can now use: envsend .env & pause"; Description: "Show completion message"; Flags: postinstall nowait skipifsilent

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"

[Files]
Source: "..\..\bin\envsend.exe"; DestDir: "{app}"; Flags: ignoreversion

[Tasks]
Name: "envPath"; Description: "Add EnvSend to system PATH (REQUIRED - Check this!)"; GroupDescription: "Installation options:"; Flags: checkedonce

[Registry]
Root: HKLM; Subkey: "SYSTEM\CurrentControlSet\Control\Session Manager\Environment"; ValueType: expandsz; ValueName: "Path"; ValueData: "{olddata};{app}"; Check: NeedsAddPath('{app}'); Tasks: envPath

[Code]
function NeedsAddPath(Param: string): boolean;
var
  OrigPath: string;
begin
  if not RegQueryStringValue(HKEY_LOCAL_MACHINE, 'SYSTEM\CurrentControlSet\Control\Session Manager\Environment', 'Path', OrigPath) then begin
    Result := True;
    exit;
  end;
  Result := Pos(';' + Param + ';', ';' + OrigPath + ';') = 0;
end;
