param(
  # Path to a built digiemu executable. If omitted, we build a temp exe and use it.
  [string]$ExePath = ""
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

# Resolve repo root reliably from script location.
$ScriptPath = $MyInvocation.MyCommand.Path
if ([string]::IsNullOrWhiteSpace($ScriptPath)) { throw "cannot determine script path" }
$ScriptDir = Split-Path -Parent $ScriptPath
$RepoRoot  = (Resolve-Path (Join-Path $ScriptDir "..")).Path

function Build-TempExeIfNeeded {
  if (-not [string]::IsNullOrWhiteSpace($ExePath)) {
    if (-not (Test-Path -LiteralPath $ExePath)) { throw "ExePath does not exist: $ExePath" }
    return $ExePath
  }

  $tempRoot = $env:RUNNER_TEMP
  if ([string]::IsNullOrWhiteSpace($tempRoot)) { $tempRoot = $env:TEMP }
  if ([string]::IsNullOrWhiteSpace($tempRoot)) { $tempRoot = [System.IO.Path]::GetTempPath() }

  $outExe = Join-Path $tempRoot "digiemu-demo-bundles.exe"

  Push-Location $RepoRoot
  try {
    & go build -o $outExe ./cmd/digiemu | Out-Null
  } finally {
    Pop-Location
  }

  if (-not (Test-Path -LiteralPath $outExe)) { throw "failed to build temp exe: $outExe" }
  return $outExe
}

function Invoke-DigiemuVerifyJson {
  param(
    [Parameter(Mandatory = $true)][string]$Exe,
    [Parameter(Mandatory = $true)][string]$BundlePath
  )

  $psi = New-Object System.Diagnostics.ProcessStartInfo
  $psi.FileName = $Exe
  $psi.Arguments = "verify --bundle `"$BundlePath`" --json=canonical"
  $psi.RedirectStandardOutput = $true
  $psi.RedirectStandardError  = $true
  $psi.UseShellExecute        = $false
  $psi.CreateNoWindow         = $true

  $p = New-Object System.Diagnostics.Process
  $p.StartInfo = $psi
  [void]$p.Start()
  $stdout = $p.StandardOutput.ReadToEnd()
  $stderr = $p.StandardError.ReadToEnd()
  $p.WaitForExit()

  if (-not [string]::IsNullOrWhiteSpace($stderr)) {
    throw "stderr must be empty in JSON mode (bundle '$BundlePath')`nstderr=$stderr"
  }

  $obj = $null
  try {
    $obj = $stdout | ConvertFrom-Json
  } catch {
    throw "invalid JSON stdout for bundle '$BundlePath'`nstdout=$stdout`nstderr=$stderr"
  }

  return [pscustomobject]@{
    ExitCode = $p.ExitCode
    Stdout   = $stdout
    Stderr   = $stderr
    Result   = $obj
  }
}

$Exe = Build-TempExeIfNeeded

Push-Location $RepoRoot
try {
  $okBundle   = Join-Path $RepoRoot "examples\bundles\demo_ok_bundle_v1\snapshots\snapshot_demo_v1"
  $failBundle = Join-Path $RepoRoot "examples\bundles\demo_fail_bundle_v1\snapshots\snapshot_demo_v1"

  $okSnap   = Join-Path $okBundle "snapshot.json"
  $failSnap = Join-Path $failBundle "snapshot.json"

  if (-not (Test-Path -LiteralPath $okSnap))   { throw "missing OK bundle snapshot.json: $okSnap" }
  if (-not (Test-Path -LiteralPath $failSnap)) { throw "missing FAIL bundle snapshot.json: $failSnap" }

  Write-Host "Running verify on demo_ok_bundle_v1..."
  $ok = Invoke-DigiemuVerifyJson -Exe $Exe -BundlePath $okBundle
  if ($ok.ExitCode -ne 0) {
    throw "demo_ok_bundle_v1 expected exit 0, got $($ok.ExitCode)`nstdout=$($ok.Stdout)"
  }
  if (-not $ok.Result.ok) {
    throw "demo_ok_bundle_v1 expected ok=true`nstdout=$($ok.Stdout)"
  }

  Write-Host "Running verify on demo_fail_bundle_v1..."
  $fail = Invoke-DigiemuVerifyJson -Exe $Exe -BundlePath $failBundle
  if ($fail.ExitCode -eq 0) {
    throw "demo_fail_bundle_v1 expected non-zero exit code, got 0`nstdout=$($fail.Stdout)"
  }
  if ($fail.Result.ok) {
    throw "demo_fail_bundle_v1 expected ok=false`nstdout=$($fail.Stdout)"
  }

  Write-Host "OK: demo bundles gate passed"
} finally {
  Pop-Location
}