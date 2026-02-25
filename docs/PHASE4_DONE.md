# DigiEmu Core — Phase 4 Completion Report

## Status
Phase 4 is formally completed.

All acceptance criteria, integration tests and CI gates are active and passing.

Commit reference: 888691c
(Fix: sync demo expected_hash_v1 with Verify/Replay v1)

---

## Scope of Phase 4

Phase 4 focused on:

- Deterministic Replay as single reconstruction path
- Verify built on Replay (no duplicated state assembly)
- Canonical hash scope: canonical_json_v1_excluding_expected_hash_v1
- Stable JSON output contract
- Explicit exit code governance
- Bundle resolver governance
- Cross-platform determinism (Linux + Windows)
- CI enforcement (gofmt + tests + acceptance script)

---

## Deterministic Guarantees

1. ReplayV1 is the single source of truth for state reconstruction.
2. VerifyV1 hashes only the replayed normalized state.
3. expected_hash_v1 is excluded from the canonical hash scope.
4. Replay --json is byte-identical across repeated executions.
5. CI enforces:
   - gofmt cleanliness
   - go test ./...
   - Windows acceptance script
6. Line endings are normalized via .gitattributes (LF enforced for Go/YAML/JSON/MD).

---

## Exit Code Governance (Locked)

0  → Verification OK
2  → Hash mismatch
3  → Write governance blocked
4  → Snapshot not found
5  → Invalid usage

Strict mode does not override exit codes.

---

## Acceptance Procedure

Local acceptance:

    gofmt -w .
    go test ./...
    powershell -ExecutionPolicy Bypass -File .\scripts\accept_phase4.ps1

CI matrix:

    ubuntu-latest
    windows-latest

Phase 4 is considered complete when:
    - CI is green on both platforms
    - Working tree is clean
    - Demo fixture hash matches canonical hash

---

## Architectural Outcome

DigiEmu Core now provides:

- Deterministic snapshot reconstruction
- Auditable hash validation
- Stable replay and verify semantics
- Cross-platform reproducibility
- Governance-safe expected hash handling

Phase 4 establishes the system as a deterministic verification tool,
not just a prototype.

---

End of Phase 4.
Further development continues in Phase 5.
