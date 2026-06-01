# DigiEmu Core 2.0 — Verify Result Reason Codes (Draft Registry)

Status: DRAFT

Purpose: A normative registry of reason codes used by `Verify Result` v2 outputs. Codes are organized by `result` group and are intended to be stable, machine-readable identifiers.

Usage guidance
--------------
- Reason codes are canonical short identifiers describing why a verification `result` is `PASS`, `FAIL`, or `ERROR`.
- Consumers SHOULD use the `reason_code` field in `Verify Result` v2 to implement deterministic handling logic.
- New codes MUST be added here as draft entries and given a clear short description and semantic meaning.

PASS
----
- `STATE_RECONSTRUCTED` — The snapshot state was successfully reconstructed and validated against the expected proof.

FAIL
----
- `HASH_MISMATCH` — The computed snapshot hash does not match the expected snapshot hash.
- `RECONSTRUCTION_RULE_MISMATCH` — The reconstruction rules or profile used produced a differing reconstruction outcome.
- `CHAIN_BROKEN` — Evidence chain continuity is broken (missing or non-matching intermediate artifact).
- `TRANSITION_INVALID` — A state transition violated the declared transition rules or constraints.
- `INSIDE_HASH_FIELD_CHANGED` — A field that is part of the inside-hash changed unexpectedly.

ERROR
-----
- `CANONICALIZATION_ERROR` — An error occurred while canonicalizing the input for hashing.
- `INVALID_SNAPSHOT_SCHEMA` — The produced snapshot does not conform to the declared snapshot schema.
- `MISSING_REQUIRED_FIELD` — A required field for snapshot or verification is missing.
- `UNSUPPORTED_HASH_ALGORITHM` — The requested hash algorithm is not supported by this verifier.
- `UNSUPPORTED_CANONICALIZATION_PROFILE` — The canonicalization profile is not recognized or supported.
- `UNSUPPORTED_RECONSTRUCTION_PROFILE` — The reconstruction profile is not recognized or supported.
- `INTERNAL_ERROR` — An internal error occurred preventing verification.

Extensibility
-------------
New reason codes SHOULD be added under the appropriate `result` group and MUST include a short description and suggested remediation or consumer action.
