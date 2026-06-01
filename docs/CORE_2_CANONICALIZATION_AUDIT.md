# Core 2.0 Canonicalization Audit (Draft)

Status: DRAFT

Purpose
-------
Audit the existing canonicalization and snapshot hashing paths to establish the current normative behavior, surface risks, and recommend minimal hardening steps for Core 2.0.

Scope
-----
Code paths inspected:
- `internal/canonicaljson/` (canonical JSON encoder)
- `internal/verify/` (ReplayV1 + StateV1 assembly)
- `pkg/snapshot` (HashV1FromState)
- CLI canonical output (`cmd/digiemu/json_output.go`)
- Existing tests around canonical JSON and snapshot hashing

Executive summary (short)
-------------------------
- There is a single canonical JSON implementation used by Core for deterministic hashing: `internal/canonicaljson.Marshal`.
- Snapshot hashing (v1) is computed by applying `canonicaljson.Marshal` to the reconstructed state and then computing SHA-256 (`pkg/snapshot.HashV1FromState`).
- The CLI's `--json=canonical` output normalizes through `encoding/json` before calling `canonicaljson.Marshal` to preserve `encoding/json` semantics for struct tags and `omitempty`.
- `ReplayV1` removes `expected_hash_v1` from the snapshot prior to reconstruction, ensuring outside-hash metadata is excluded from the hashed scope.

Answers to audit questions
-------------------------
- Is there exactly one normative canonicalization path?
  - Yes for hashing: `pkg/snapshot.HashV1FromState` → `internal/canonicaljson.Marshal` is the normative path used by verification and golden tests.
  - Note: CLI canonical output goes through `encoding/json` first (normalization for `omitempty`), which is intentionally different and used for human/CLI outputs only.

- Are object keys sorted deterministically?
  - Yes. `internal/canonicaljson` sorts `map[string]...` keys using `sort.Strings` before serialization.

- Is whitespace eliminated consistently?
  - Yes. The encoder emits compact JSON with no insignificant whitespace.

- Are arrays treated as order-sensitive?
  - Yes. Arrays and slices are serialized in original order and are treated as order-sensitive.

- How are numbers encoded?
  - Integers: encoded with `strconv.FormatInt`/`FormatUint` in base 10 (no extra formatting).
  - Floats: encoded using `strconv.FormatFloat` with format `'g', -1, 64` which produces a compact, stable textual representation but may not preserve some decimal formatting nuances (e.g., trailing zeros).

- How are strings and Unicode handled?
  - Strings are encoded using `strconv.Quote`, which applies Go's string quoting and escaping rules. Unicode characters are encoded according to Go's `strconv.Quote` behavior (printable Unicode typically left as-is; control characters escaped). This is deterministic but should be documented explicitly as the canonical string encoding.

- Are outside-hash metadata fields excluded consistently?
  - `expected_hash_v1` is explicitly removed by `ReplayV1` prior to hashing. Other outside-hash metadata must be excluded by the code that assembles the `StateV1` (the `StateV1FromBundle` composition). Consumers MUST ensure only inside-hash fields are present in the reconstructed state.

- Do existing tests cover field-order stability?
  - Yes. `internal/canonicaljson/canonicaljson_test.go` contains tests for map key ordering, nested maps, struct field order, and whitespace-agnostic canonicalization.

- Do existing tests cover inside-hash changes versus outside-hash changes?
  - There are tests around replay determinism and golden snapshot hashes (`internal/verify/replay_determinism_test.go`, `golden_snapshot_hash_test.go`) that exercise the canonicalization + hashing pipeline and implicitly validate that `expected_hash_v1` removal prevents self-reference. However, explicit tests that demonstrate a clear separation of an "inside-hash" vs arbitrary outside-hash metadata (other than `expected_hash_v1`) are limited; adding a small set of targeted conformance vectors is recommended.

Notable implementation details / risks
-----------------------------------
- json.RawMessage handling: `StateV1` and `BundleV1` fields are typed as `json.RawMessage` (alias `[]byte`). `internal/canonicaljson.Marshal` does not special-case `json.RawMessage` and will serialize `[]byte` data as a numeric array (each byte encoded as a number) if a RawMessage is passed directly as a field value. The codebase currently relies on the interplay between `ReplayV1` (which normalizes snapshot JSON into `json.RawMessage`) and the `StateV1` shape; tests and golden hashes currently reflect the repository's present behavior. This behavior is non-obvious and should be documented; it may be a candidate for hardening (e.g., treat `json.RawMessage` as pre-encoded JSON and embed it rather than encode as bytes).

- Struct vs JSON normalization: CLI canonical output intentionally runs data through `encoding/json` first to preserve `omitempty` and other struct-tag semantics, whereas hashing uses `canonicaljson.Marshal` directly on Go values reconstructed by ReplayV1. This can lead to subtle differences if callers compute canonical JSON differently — ensure canonicalization for hashing is always performed via the same reconstruct-then-hash path.

- Float formatting: using `'g'` with precision `-1` is compact and stable but may produce different textual forms for logically equal numeric values (e.g., `1.0` vs `1`). If exact textual digit patterns are a requirement, consider adding explicit tests and guidelines.

- Unicode normalization: `strconv.Quote` provides deterministic escaping but does not perform Unicode NFKC/NFC normalization. If canonicalization must be independent of Unicode normalization creases, add explicit normalization at a pre-processing step.

Recommendations (minimal, non-invasive)
------------------------------------
1. Document the canonicalization behavior clearly in a normative place (this file and `docs/SNAPSHOT_HASH_v1.0.md`): explain `json.RawMessage` handling, number formatting, Unicode behavior, and the canonical scope used by hashing (what is inside-hash vs outside-hash).

2. Add explicit conformance tests that demonstrate inside-hash vs outside-hash behavior (e.g., ensure adding `created_at` outside the snapshot does not change the hash; ensure `expected_hash_v1` is ignored). These live in `testdata/core_2_conformance/` (already added) and should be expanded with targeted cases.

3. Consider a small, opt-in change in a future minor release to treat `json.RawMessage` specially in `internal/canonicaljson.Marshal` when the underlying type is `json.RawMessage` (i.e., detect `[]byte` with the named type and unmarshal/embed rather than encode as bytes). This would be a breaking change for existing hashes and therefore MUST be gated behind a versioned canonicalization profile (e.g., `digiemu-canonical-json-v2`) and migration guidance.

4. Add tests for float and Unicode edge-cases so potential portability concerns are caught by CI across platforms.

5. Maintain a single canonicalization entrypoint for hashing (`pkg/snapshot.HashV1FromState`) and ensure any other module that needs to compute hashes calls that function rather than reimplementing encoding logic.

Suggested hardening roadmap (non-invasive)
-----------------------------------------
- Short term (docs + tests): publish this audit, add conformance vectors that demonstrate inside/outside-hash separation and float/Unicode edge-cases.
- Medium term (profiles + opt-in behavior): define `canonicalization_profile` metadata and support alternative profiles in a backward-compatible way (add profile field, support multiple profiles for validation, but keep `sha256`+`digiemu-canonical-json-v1` as default for v1 compatibility).
- Long term (opt-in runtime behavior): if needed, implement `json.RawMessage` embedding and Unicode normalization in a new `canonicaljson` profile and provide tooling to re-evaluate historical proofs.

Appendix — quick code references
-------------------------------
- Canonical encoder: `internal/canonicaljson/canonicaljson.go`  
- Hash computation: `pkg/snapshot/hashv1.go`  
- Replay/assembly: `internal/verify/replay.go`, `internal/verify/state_v1.go`  
- CLI canonical output (normalization): `cmd/digiemu/json_output.go`
