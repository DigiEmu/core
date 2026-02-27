param(
  [string]$Expected = "examples/snapshot_v1_demo/expected_verify_report.json",
  [string]$Got      = "examples/snapshot_v1_demo/_out/release_report.json"
)

. "$PSScriptRoot/canon-json.ps1"
Assert-CanonJsonEqual $Expected $Got "verify report"
