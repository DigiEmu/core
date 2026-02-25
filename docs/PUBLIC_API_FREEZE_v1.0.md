# Public API Freeze — v1.0 (pkg/*)

Effective date: 2026-02-25

This document marks the **v1.0 public API freeze** for DigiEmu Core.

## Scope

The **public API** is the set of exported identifiers in:

- `pkg/meaning`
- `pkg/claims`
- `pkg/uncertainty`
- `pkg/snapshot`
- `pkg/verify`

Everything under `internal/*` is **not** part of the public API and may change without notice.

## Compatibility promise (v1.x)

Within the `v1.x` line:

- Existing exported identifiers in `pkg/*` MUST remain source-compatible.
- Behavior MUST remain compatible unless explicitly documented as a bug fix.
- New exported identifiers MAY be added in minor/patch releases.
- Breaking changes MUST wait for `v2`.

See:
- `docs/COMPATIBILITY_POLICY_v1.0.md`
- `docs/PUBLIC_API_v1.0.md`
- `docs/PUBLIC_API_README.md`

## Notes

This freeze is about **stable interfaces**.
Internal implementations may evolve, and additional public packages may be introduced later,
but `pkg/*` listed above is treated as stable from this point onward.
