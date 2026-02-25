# DigiEmu Core — Phase4 Acceptance Script (Windows / PowerShell)
# Guarantees:
# - gofmt applied
# - go test ./... passed
# - CLI builds
# - verify + replay smoke pass on demo fixture bundle
# - reports git status
#
# Usage:
#   powershell -ExecutionPolicy Bypass -File .\scripts\accept_phase4.ps1
#   (or) .\scripts\accept_phase4.ps1

$ErrorActionPreference = "Stop"

function Run($cmd) {
  Write-Host "`n> $cmd" -ForegroundColor Cyan
  iex $cmd
  if ($LASTEXITCODE -ne $null -and $LASTEXITCODE -ne 0) {
    throw "Command failed with exit code $($LASTEXITCODE): $cmd"
  }
}

# Verify we're in repo root
if (-not (Test-Path ".\go.mod")) {
  throw "go.mod not found. Please run from repo root."
}

# 1) Format
Run "gofmt -w ."

# 2) Unit tests
Run "go test ./..."

# 3) Build CLI (local artifact; must be ignored by git)
$exe = ".\digiemu.exe"
Run "go build -o $exe .\cmd\digiemu"

# 4) Smoke: verify + replay on demo bundle
$bundle = ".\data\test-fixtures\snapshots\demo"
if (-not (Test-Path $bundle)) {
  throw "Demo bundle not found at: $bundle"
}

Run "$exe verify --bundle $bundle --json | Out-Host"
Run "$exe replay  --bundle $bundle --json | Out-Host"

# 5) Git status (informational)
Write-Host "`n> git status --porcelain" -ForegroundColor Cyan
$st = git status --porcelain
if ($st) {
  Write-Host "WORKING TREE NOT CLEAN:" -ForegroundColor Yellow
  $st | Out-Host
  Write-Host "Note: This script does not auto-commit. Review changes before committing." -ForegroundColor Yellow
} else {
  Write-Host "WORKING TREE CLEAN" -ForegroundColor Green
}

Write-Host "`nACCEPTANCE OK" -ForegroundColor Green
