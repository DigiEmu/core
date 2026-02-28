param(
  # Path to a built digiemu executable. If omitted, we use `go run ./cmd/digiemu`.
  [string]$ExePath = ""
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

# Resolve repo root reliably from script location.
$ScriptPath = $MyInvocation.MyCommand.Path
if ([string]::IsNullOrWhiteSpace($ScriptPath)) { throw "cannot determine script path" }
$ScriptDir = Split-Path -Parent $ScriptPath
$RepoRoot  = (Resolve-Path (Join-Path $ScriptDir "..")).Path

function Invoke-DigiemuVerifyJson {
  param(
    [Parameter(Mandatory = $true)][string]$BundlePath
  )

  $psi = New-Object System.Diagnostics.ProcessStartInfo

  if ([string]::IsNullOrWhiteSpace($ExePath)) {
    $psi.FileName = "go"
    $psi.Arguments = "run ./cmd/digiemu verify --bundle `"$BundlePath`" --json=canonical"
  } else {
    if (-not (Test-Path -LiteralPath $ExePath)) {
      throw "ExePath does not exist: $ExePath"
    }
    $psi.FileName = $ExePath
    $psi.Arguments = "verify --bundle `"$BundlePath`" --json=canonical"
  }

  $psi.RedirectStandardOutput = $true
  $psi.RedirectStandardError = $true
  $psi.UseShellExecute = $false
  $psi.CreateNoWindow = $true

  $p = New-Object System.Diagnostics.Process
  $p.StartInfo = $psi
  [void]$p.Start()
  $stdout = $p.StandardOutput.ReadToEnd()
  $stderr = $p.StandardError.ReadToEnd()
  $p.WaitForExit()

  $obj = $null
  try {
    $obj = $stdout | ConvertFrom-Json
  } catch {
    throw "invalid JSON stdout for bundle '$BundlePath'\nstdout=$stdout\nstderr=$stderr"
  }

  return [pscustomobject]@{
    ExitCode = $p.ExitCode
    Stdout   = $stdout
    Stderr   = $stderr
    Result   = $obj
  }
}

Push-Location $RepoRoot
try {
  $okBundle = Join-Path $RepoRoot "examples\bundles\demo_ok_bundle_v1\snapshots\snapshot_demo_v1"
  $failBundle = Join-Path $RepoRoot "examples\bundles\demo_fail_bundle_v1\snapshots\snapshot_demo_v1"

  $okSnap = Join-Path $okBundle "snapshot.json"
  $failSnap = Join-Path $failBundle "snapshot.json"

  if (-not (Test-Path -LiteralPath $okSnap)) { throw "missing OK bundle snapshot.json: $okSnap" }
  if (-not (Test-Path -LiteralPath $failSnap)) { throw "missing FAIL bundle snapshot.json: $failSnap" }

  Write-Host "Running verify on demo_ok_bundle_v1..."
  $ok = Invoke-DigiemuVerifyJson -BundlePath $okBundle
  if (-not $ok.Result.ok) {
    throw "demo_ok_bundle_v1 expected ok=true\nstdout=$($ok.Stdout)\nstderr=$($ok.Stderr)"
  }
  if (-not [string]::IsNullOrWhiteSpace($ExePath) -and $ok.ExitCode -ne 0) {
    throw "demo_ok_bundle_v1 expected exit 0, got $($ok.ExitCode)"
  }

  Write-Host "Running verify on demo_fail_bundle_v1..."
  $fail = Invoke-DigiemuVerifyJson -BundlePath $failBundle
  if ($fail.Result.ok) {
    throw "demo_fail_bundle_v1 expected ok=false\nstdout=$($fail.Stdout)\nstderr=$($fail.Stderr)"
  }
  if (-not [string]::IsNullOrWhiteSpace($ExePath) -and $fail.ExitCode -eq 0) {
    throw "demo_fail_bundle_v1 expected non-zero exit code"
  }

  Write-Host "OK: demo bundles gate passed"
} finally {
  Pop-Location
}
