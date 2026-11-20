# Creates the desktop shortcut "LYKOS REIMARCHIV" pointing to the built executable.
# Usage: powershell -ExecutionPolicy Bypass -File scripts\create-shortcut.ps1
$ErrorActionPreference = "Stop"
$root = Split-Path -Parent $PSScriptRoot
$exe = Join-Path $root "packaging\bin\LykosReimarchiv.exe"

if (-not (Test-Path $exe)) {
    Write-Host "Executable not found: $exe" -ForegroundColor Yellow
    Write-Host "Build first: powershell -ExecutionPolicy Bypass -File scripts\build.ps1"
    exit 1
}

# GetFolderPath also resolves a desktop redirected to OneDrive.
$desktop = [Environment]::GetFolderPath("Desktop")
$lnk = Join-Path $desktop "LYKOS REIMARCHIV.lnk"

$shell = New-Object -ComObject WScript.Shell
$sc = $shell.CreateShortcut($lnk)
$sc.TargetPath = $exe
$sc.WorkingDirectory = Split-Path -Parent $exe
$sc.IconLocation = "$exe,0"
$sc.Description = "LYKOS REIMARCHIV"
$sc.Save()

Write-Host "Shortcut: $lnk -> $exe" -ForegroundColor Green
