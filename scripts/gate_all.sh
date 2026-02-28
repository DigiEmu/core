#!/usr/bin/env bash
set -euo pipefail

EXE_PATH="${1:-}"

echo ""
echo "=== Go unit tests ==="
go test ./...

echo ""
echo "=== Demo verify gate (canonical JSON) ==="
if [[ -z "${EXE_PATH}" ]]; then
  pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/gate_demo_verify.ps1
else
  pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/gate_demo_verify.ps1 -ExePath "${EXE_PATH}"
fi

echo ""
echo "=== Demo bundles gate (OK + FAIL) ==="
if [[ -z "${EXE_PATH}" ]]; then
  pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/gate_demo_bundles.ps1
else
  pwsh -NoProfile -ExecutionPolicy Bypass -File scripts/gate_demo_bundles.ps1 -ExePath "${EXE_PATH}"
fi

echo ""
echo "=== Schema lock test (explicit) ==="
go test ./internal/tests -run TestVerifyResultSchemaLocked -count=1

echo ""
echo "OK: all gates passed"
