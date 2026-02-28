param(
  # Path to a built digiemu executable. If omitted, we try to auto-detect one in ./_release.
  [string]$ExePath = "",

  [string]$Bundle = "examples/snapshot_v1_demo/snapshots/snapshot_demo_v1",
  [string]$Expected = "examples/snapshot_v1_demo/expected_verify_report.json",
  [string]$Out = "examples/snapshot_v1_demo/_out/release_report.json"
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

function Resolve-RepoRoot {
  $here = Split-Path -Parent $PSScriptRoot
  return $here
}

$root = Resolve-RepoRoot
Set-Location $root

if (-not $ExePath) {
  $candidate = Get-ChildItem -ErrorAction SilentlyContinue -File -Path "_release" -Filter "digiemu-*.exe" | Select-Object -First 1
  if ($candidate) {
    $ExePath = $candidate.FullName
  }
}

if (-not $ExePath) {
  throw "No -ExePath provided and no _release/digiemu-*.exe found. Build a release binary first, or pass -ExePath."
}

if (-not (Test-Path -LiteralPath $ExePath)) {
  throw "ExePath does not exist: $ExePath"
}

$outDir = Split-Path -Parent $Out
if ($outDir) {
  New-Item -ItemType Directory -Force -Path $outDir | Out-Null
}

$stderrPath = Join-Path $outDir "release_stderr.txt"

$psi = New-Object System.Diagnostics.ProcessStartInfo
$psi.FileName = $ExePath
$psi.Arguments = "verify --bundle `"$Bundle`" --json"
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

# Write stdout as UTF-8 bytes (avoid PowerShell's UTF-16LE redirection pitfalls)
[System.IO.File]::WriteAllBytes($Out, [System.Text.Encoding]::UTF8.GetBytes($stdout))
[System.IO.File]::WriteAllText($stderrPath, $stderr)

# Gate: canonical JSON compare
& "$PSScriptRoot/check_verify_report.ps1" -Expected $Expected -Got $Out

if ($p.ExitCode -ne 0) {
  # Verify can legitimately exit non-zero while still producing deterministic JSON.
  Write-Host "NOTE: digiemu exited with code $($p.ExitCode) (stdout was still gated)."
}
