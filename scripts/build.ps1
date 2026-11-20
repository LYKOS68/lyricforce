# Builds the Windows executable: packaging\bin\LykosReimarchiv.exe
# The full build output is also written to workspace\logs\build-<timestamp>.log.
# Usage: powershell -ExecutionPolicy Bypass -File scripts\build.ps1 [-Shortcut]
param([switch]$Shortcut)
$ErrorActionPreference = "Stop"
$root = Split-Path -Parent $PSScriptRoot
Set-Location $root
. (Join-Path $PSScriptRoot "common.ps1")

$wails = Get-WailsCli
if (Get-Process LykosReimarchiv -ErrorAction SilentlyContinue) {
    Write-Host "LYKOS REIMARCHIV is running - close it before building (the executable is locked)." -ForegroundColor Yellow
    exit 1
}
$logDir = Join-Path $root "workspace\logs"
New-Item -ItemType Directory -Force $logDir | Out-Null
$log = Join-Path $logDir ("build-{0}.log" -f (Get-Date -Format "yyyyMMdd-HHmmss"))

& $wails build -clean -platform windows/amd64 -trimpath 2>&1 | Tee-Object -FilePath $log
if ($LASTEXITCODE -ne 0) { throw "Build failed (exit $LASTEXITCODE), see $log" }

$exe = Join-Path $root "packaging\bin\LykosReimarchiv.exe"
Write-Host "Executable: $exe" -ForegroundColor Green
Write-Host "Build log:  $log"

if ($Shortcut) { & (Join-Path $PSScriptRoot "create-shortcut.ps1") }
