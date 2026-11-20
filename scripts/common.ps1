# Shared helpers for the scripts in this folder (dot-sourced).

# Returns the path of the Wails CLI or exits with install instructions.
function Get-WailsCli {
    $cmd = Get-Command wails -ErrorAction SilentlyContinue
    if ($cmd) { return $cmd.Source }
    $candidate = Join-Path (go env GOPATH) "bin\wails.exe"
    if (Test-Path $candidate) { return $candidate }
    Write-Host "Wails CLI not found. Install with:" -ForegroundColor Yellow
    Write-Host "  go install github.com/wailsapp/wails/v2/cmd/wails@v2.17.1"
    exit 1
}
