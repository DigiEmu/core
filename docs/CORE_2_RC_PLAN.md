# DigiEmu Core 2.0 RC Plan

## Current checkpoint

- core-2.0-draft4-self-audit

## RC1 goal

Deliver a minimal, deterministic, partner-testable verification core with stable conformance rules, schemas, and vectors.

## Before RC1

- Confirm README boundary wording
- Confirm threat model completeness
- Confirm conformance rules
- Confirm test vectors
- Confirm Go module/tagging strategy
- Confirm no Core/Enterprise boundary leakage
- Confirm no TBN/DigiEmu responsibility overlap

## RC1 criteria

- All tests pass
- All conformance vectors pass
- No known module versioning risks
- README and docs do not overclaim compliance, identity, or trust certification
- External reviewer feedback incorporated or documented

## Not required for RC1

- Enterprise dashboard
- Key management
- Role system
- Legal compliance certification
- Agent identity certification
- Trust scoring

## Release/tagging safety

- Do not use real v2.0.0 tags unless the module path intentionally ends in /v2
- Use safe milestone names such as `core-2.0-rc1` until /v2 strategy is explicit

## External feedback handling

- Track Andrei feedback on conformance and standardization (see `docs/ANDREI_REVIEW_BRIEF.md`)
- Track Burhan feedback on the TBN boundary (see `docs/TBN_DIGIEMU_BOUNDARY_NOTE.md`)
- Convert accepted feedback into vectors or doc clarifications

## Recommended next milestone

- `core-2.0-rc1` (no /v2 module path change yet)
