Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$script:Failures = New-Object System.Collections.Generic.List[string]

function Write-Pass {
  param([string]$Name)
  Write-Host "PASS: $Name"
}

function Write-Fail {
  param([string]$Name, [string]$Message)
  $script:Failures.Add($Name) | Out-Null
  Write-Host "FAIL: $Name - $Message"
}

function Invoke-Check {
  param([string]$Name, [scriptblock]$Block)
  try {
    & $Block
    Write-Pass $Name
  } catch {
    Write-Fail $Name $_.Exception.Message
  }
}

function Invoke-CommandText {
  param([string]$Name, [string[]]$Command)
  $exe = $Command[0]
  $args = @()
  if ($Command.Length -gt 1) {
    $args = $Command[1..($Command.Length - 1)]
  }
  $output = & $exe @args 2>&1
  $exitCode = $LASTEXITCODE
  if ($exitCode -ne 0) {
    throw "$Name failed with exit ${exitCode}: $($output -join [Environment]::NewLine)"
  }
  return ($output -join [Environment]::NewLine)
}

function Get-JsonObjectFromOutput {
  param([string]$Output)
  $start = $Output.IndexOf('{')
  if ($start -lt 0) {
    throw "JSON object was not found in output"
  }
  return $Output.Substring($start) | ConvertFrom-Json
}

$scriptPath = $MyInvocation.MyCommand.Path
if ([string]::IsNullOrWhiteSpace($scriptPath)) {
  Write-Host "FAIL: Resolve script path - cannot determine script path"
  exit 1
}

$scriptDir = Split-Path -Parent $scriptPath
$repoRoot = (Resolve-Path (Join-Path $scriptDir "..")).Path

Push-Location $repoRoot
try {
  Write-Host "Core 2.0 Draft 4 self-audit"
  Write-Host "Repository root: $repoRoot"
  Write-Host ""

  Invoke-Check "Git working tree clean" {
    $status = git status --porcelain
    if ($LASTEXITCODE -ne 0) {
      throw "git status failed"
    }
    if (($status | Measure-Object).Count -ne 0) {
      throw "working tree has uncommitted changes"
    }
  }

  Invoke-Check "Tag v2.0.0-draft.4 exists" {
    git rev-parse --verify --quiet "refs/tags/v2.0.0-draft.4" | Out-Null
    if ($LASTEXITCODE -ne 0) {
      throw "tag v2.0.0-draft.4 was not found"
    }
  }

  Invoke-Check "Go tests" {
    Invoke-CommandText "go test ./..." @("go", "test", "./...") | Out-Null
  }

  $humanOutput = $null
  Invoke-Check "Conformance CLI human-readable" {
    $humanOutput = Invoke-CommandText "conformance human-readable" @("go", "run", "./cmd/digiemu", "experimental", "conformance", "run", "testdata/core_2_conformance")
    if ($humanOutput -notmatch "total=11 passed=11 failed=0") {
      throw "expected total=11 passed=11 failed=0, got: $humanOutput"
    }
  }

  $jsonReport = $null
  Invoke-Check "Conformance CLI JSON" {
    $jsonOutput = Invoke-CommandText "conformance JSON" @("go", "run", "./cmd/digiemu", "experimental", "conformance", "run", "testdata/core_2_conformance", "--json")
    $jsonReport = Get-JsonObjectFromOutput $jsonOutput
    if ($jsonReport.status -ne "PASS") { throw "expected status PASS, got $($jsonReport.status)" }
    if ([int]$jsonReport.total -ne 11) { throw "expected total 11, got $($jsonReport.total)" }
    if ([int]$jsonReport.passed -ne 11) { throw "expected passed 11, got $($jsonReport.passed)" }
    if ([int]$jsonReport.failed -ne 0) { throw "expected failed 0, got $($jsonReport.failed)" }
    if ($jsonReport.report_version -ne "core-2-conformance-report-v1") { throw "unexpected report_version $($jsonReport.report_version)" }
  }

  Invoke-Check "Basic report fixture consistency" {
    $fixturePath = Join-Path $repoRoot "testdata/core_2_conformance/report_expected_basic.json"
    $fixture = Get-Content -Raw $fixturePath | ConvertFrom-Json
    if ([int]$fixture.total -ne 11) { throw "expected fixture total 11, got $($fixture.total)" }
    if ([int]$fixture.passed -ne 11) { throw "expected fixture passed 11, got $($fixture.passed)" }
    $case = $fixture.cases | Where-Object { $_.name -eq "unsupported_canonicalization_profile_error" } | Select-Object -First 1
    if ($null -eq $case) { throw "missing unsupported_canonicalization_profile_error" }
    if ($case.expected_result -ne "ERROR") { throw "expected unsupported case result ERROR, got $($case.expected_result)" }
    if ($case.reason_code -ne "UNSUPPORTED_CANONICALIZATION_PROFILE") { throw "expected UNSUPPORTED_CANONICALIZATION_PROFILE, got $($case.reason_code)" }
  }

  Invoke-Check "MISSING_REFERENCE not in Verify Result v2 schema" {
    $schemaText = Get-Content -Raw (Join-Path $repoRoot "schemas/verify_result_v2.schema.json")
    if ($schemaText -match '"MISSING_REFERENCE"') {
      throw "MISSING_REFERENCE appears in Verify Result v2 schema"
    }
  }

  Invoke-Check "MISSING_REFERENCE consideration doc exists" {
    $path = Join-Path $repoRoot "docs/CORE_2_MISSING_REFERENCE_SCHEMA_CONSIDERATION.md"
    if (-not (Test-Path $path)) { throw "missing $path" }
  }

  Invoke-Check "Draft 4 post-tag note exists" {
    $path = Join-Path $repoRoot "docs/CORE_2_DRAFT_4_POST_TAG_NOTE.md"
    if (-not (Test-Path $path)) { throw "missing $path" }
  }

  Invoke-Check "Docs keep v1.0.0 stable/latest guardrail" {
    $paths = @(
      "docs/CORE_2_DRAFT_4_POST_TAG_NOTE.md",
      "docs/CORE_2_DRAFT_4_READINESS.md",
      "docs/CORE_2_MISSING_REFERENCE_SCHEMA_CONSIDERATION.md"
    )
    foreach ($relativePath in $paths) {
      $text = Get-Content -Raw (Join-Path $repoRoot $relativePath)
      if ($text -notmatch "v1\.0\.0.*stable/latest") {
        throw "$relativePath does not mention v1.0.0 remains stable/latest"
      }
    }
  }

  Write-Host ""
  Write-Host "Audit summary"
  Write-Host "-------------"
  if ($script:Failures.Count -eq 0) {
    Write-Host "PASS: all Core 2.0 Draft 4 audit checks passed"
    exit 0
  }

  Write-Host "FAIL: $($script:Failures.Count) check(s) failed"
  foreach ($failure in $script:Failures) {
    Write-Host "- $failure"
  }
  exit 1
} finally {
  Pop-Location
}
