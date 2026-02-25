# Phase 4 — Acceptance Check (Normative)

This repository provides a single acceptance command for Phase 4 work.

## Command (Windows / PowerShell)

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\accept_phase4.ps1
```

## What it guarantees

The acceptance script MUST:

1. Run `gofmt -w .`
2. Run `go test ./...` and require a green suite
3. Build the CLI binary (`digiemu.exe`) locally
4. Run smoke checks on the demo bundle:
   - `digiemu.exe verify --bundle .\data\test-fixtures\snapshots\demo --json`
   - `digiemu.exe replay --bundle .\data\test-fixtures\snapshots\demo --json`
5. Print `git status --porcelain` for human review (informational)

The script MUST fail fast (non-zero exit) if any step fails.
