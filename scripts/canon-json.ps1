# Canonical JSON helpers (stable across CRLF/LF/whitespace/indentation)

function CanonJsonString([string]$json) {
  $obj = $json | ConvertFrom-Json
  ($obj | ConvertTo-Json -Depth 99 -Compress) + "`n"
}

function CanonJsonFile([string]$path) {
  $raw = Get-Content $path -Raw
  CanonJsonString $raw
}

function Assert-CanonJsonEqual(
  [string]$expectedPath,
  [string]$gotPath,
  [string]$label = "JSON"
) {
  $a = CanonJsonFile $expectedPath
  $b = CanonJsonFile $gotPath
  if ($a -ne $b) {
    throw "$label mismatch (canonical JSON). expected=$expectedPath got=$gotPath"
  }
  Write-Host "OK: $label canonical match"
}
