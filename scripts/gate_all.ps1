Param(
  [Parameter(Mandatory = $false)]
  [string]$ExePath = ""
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

function Resolve-PowerShellExe {
  foreach ($candidate in @("pwsh", "powershell")) {
    $cmd = Get-Command $candidate -ErrorAction SilentlyContinue
    if ($null -ne $cmd) {
      if ($cmd.Path) { return $cmd.Path }
      if ($cmd.Source) { return $cmd.Source }
      return $cmd.Name
    }
  }
  throw "Neither 'pwsh' nor 'powershell' was found in PATH"
}

$PowerShellExe = Resolve-PowerShellExe

function Run-Step {
  param([string]$Name, [scriptblock]$Block)
  Write-Host ""
  Write-Host "=== $Name ==="
  & $Block
  if ($LASTEXITCODE -ne 0) { throw "$Name failed with exit $LASTEXITCODE" }
}

# Resolve script dir + repo root reliably
$ScriptPath = $MyInvocation.MyCommand.Path
if ([string]::IsNullOrWhiteSpace($ScriptPath)) { throw "cannot determine script path" }
$ScriptDir = Split-Path -Parent $ScriptPath
$RepoRoot  = (Resolve-Path (Join-Path $ScriptDir "..")).Path

$DemoVerifyGate = Join-Path $RepoRoot "scripts\gate_demo_verify.ps1"
$DemoBundlesGate = Join-Path $RepoRoot "scripts\gate_demo_bundles.ps1"

Push-Location $RepoRoot
try {
  Run-Step "Go unit tests" {
    go test ./...
  }

  Run-Step "Demo verify gate (canonical JSON)" {
    if ([string]::IsNullOrWhiteSpace($ExePath)) {
      & $PowerShellExe -NoProfile -ExecutionPolicy Bypass -File $DemoVerifyGate
    } else {
      & $PowerShellExe -NoProfile -ExecutionPolicy Bypass -File $DemoVerifyGate -ExePath $ExePath
    }
  }

  Run-Step "Demo bundles gate (OK + FAIL)" {
    if ([string]::IsNullOrWhiteSpace($ExePath)) {
      & $PowerShellExe -NoProfile -ExecutionPolicy Bypass -File $DemoBundlesGate
    } else {
      & $PowerShellExe -NoProfile -ExecutionPolicy Bypass -File $DemoBundlesGate -ExePath $ExePath
    }
  }

  Run-Step "Schema lock test (explicit)" {
    go test ./internal/tests -run TestVerifyResultSchemaLocked -count=1
  }

  Write-Host ""
  Write-Host "OK: all gates passed"
} finally {
  Pop-Location
}
