# DigiEmu Core 2.0 Conformance

Status: draft / hardening phase

---

## Conformance Definition

An implementation is DigiEmu Core 2.0 conformant if, given the same canonical snapshot and verification inputs, it produces the same verification result as the reference implementation and passes the required test vectors.

- Determinism: same input -> same canonical JSON -> same SHA-256 snapshot hash -> same verify result.
- Independence: a third-party verifier can reproduce PASS/FAIL/ERROR with the same reason codes.

---

## Canonical JSON Requirement

- Implement `canonical_json_v1` semantics.
- Non-canonical serialization must be detected as a conformance failure when it changes the state identity.
- Canonicalization scope excludes self-referential hash fields (e.g., `expected_hash_v1`).

---

## Snapshot Hash Requirement

- Hash algorithm: `sha256(canonical_json_v1)`.
- The same canonicalized snapshot MUST produce the exact same 64-hex digest.
- Hash scope: explicitly excludes `expected_hash_v1` and includes all integrity-relevant fields.

---

## Replay Requirement

- Given a bundle, the verifier MUST reconstruct the decision state deterministically using the documented reconstruction profile.
- Any mismatch between reconstructed state and expected state MUST surface as a verification failure with a stable reason code.

---

## Verify Result Requirement

- The implementation MUST emit verify results that validate against `schemas/VERIFY_RESULT_SCHEMA_v1.json` (v1) or `schemas/verify_result_v2.schema.json` (v2 draft) when in the respective mode.
- PASS/FAIL/ERROR MUST be deterministically determined from canonicalization, hashing, and replay outcomes.

---

## PASS / FAIL / ERROR Semantics

- PASS: canonicalization and hashing succeed; computed digest equals expected; reconstruction matches profile; no schema violations.
- FAIL: verification completes but evidence indicates mismatch (e.g., hash mismatch, rule mismatch), with a stable reason code.
- ERROR: verification cannot complete due to malformed inputs, schema violations, or internal errors; includes a stable reason code.

---

## Required Test Vectors

Implementations MUST pass all vectors under `testdata/core_2_conformance/`, including at minimum:
- `basic_pass`
- `hash_mismatch_fail`
- `inside_payload_mutation_detected`
- `invalid_schema_error`
- `malformed_json_error`
- `missing_required_field_error`
- `outside_metadata_ignored`
- `unknown_reason_code_error`
- `unsupported_canonicalization_profile_error`
- `unsupported_hash_algorithm_error`
- `wrong_profile_fail`

---

## CLI Behavior

- The CLI MUST support `--json` output with `pretty` and `canonical` modes.
- In JSON mode, stdout is JSON-only; operational errors go to stderr.
- Exit codes follow the Phase A2 contract; VERIFY failures yield a non-zero exit code.

---

## JSON Report Output

- JSON verify results MUST validate against `schemas/VERIFY_RESULT_SCHEMA_v1.json` or `schemas/verify_result_v2.schema.json` (draft) when applicable.
- The `--json=canonical` mode MUST produce byte-stable canonical JSON suitable for test vectors and audit evidence.

---

## Version Compatibility

- Implementations MUST clearly declare the supported snapshot, canonicalization, and schema versions.
- Conformance is version-bound. A change in schema or canonicalization profile requires re-validation against test vectors.

---

## Non-Goals

- no agent identity certification
- no trust scoring
- no enterprise workflow
- no legal compliance claim
- no authorization framework

---

## CLI Examples

Example: verify a bundle with canonical JSON output to stdout:

```bash
go run ./cmd/digiemu verify --bundle examples/bundles/demo_ok_bundle_v1/snapshots/snapshot_demo_v1 --json=canonical
```

Example: run conformance and emit a canonical JSON report:

```bash
go run ./cmd/digiemu experimental conformance run testdata/core_2_conformance --json=canonical
```

The resulting JSON MUST validate against `schemas/VERIFY_RESULT_SCHEMA_v1.json` or the v2 draft schema where explicitly used.
