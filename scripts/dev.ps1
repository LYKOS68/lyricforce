# Runs the app in development mode (wails dev): Go changes are recompiled,
# changes in frontend\dist reload the UI.
# Usage: powershell -ExecutionPolicy Bypass -File scripts\dev.ps1
$ErrorActionPreference = "Stop"
$root = Split-Path -Parent $PSScriptRoot
Set-Location $root
. (Join-Path $PSScriptRoot "common.ps1")

& (Get-WailsCli) dev
