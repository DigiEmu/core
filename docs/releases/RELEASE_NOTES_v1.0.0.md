# Release Notes — v1.0.0 (Public API Freeze)

Date: 2026-02-25

This release marks the **v1.0 public API freeze** for DigiEmu Core.

## Highlights

- Public API surface stabilized under `pkg/*`:
  - `pkg/meaning`
  - `pkg/claims`
  - `pkg/uncertainty`
  - `pkg/snapshot`
  - `pkg/verify`
- Formal freeze marker added:
  - `docs/PUBLIC_API_FREEZE_v1.0.md`
- Compatibility policy & public API docs:
  - `docs/COMPATIBILITY_POLICY_v1.0.md`
  - `docs/PUBLIC_API_v1.0.md`
  - `docs/PUBLIC_API_README.md`

## Compatibility

Within the `v1.x` line:
- Existing exported identifiers in `pkg/*` MUST remain source-compatible.
- Breaking changes MUST wait for `v2`.

## Notes

Everything under `internal/*` is not part of the public API and may change without notice.
