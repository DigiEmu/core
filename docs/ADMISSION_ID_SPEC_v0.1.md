# P0 Admission Identifier Derivation Specification v0.1

- **Specification ID:** P0.ADMISSION.ID.v0.1
- **Architecture baseline:** `architecture-baseline.yaml` revision 0.3
- **Date:** 2026-08-07
- **Status:** Proposed

## 1. Scope and non-claims

This document defines two P0-specific, versioned, deterministic derivation profiles:

1. **P0.ADMISSION.INTENT.v0.1** — a normalized Intent digest.
2. **P0.ADMISSION.ID.v0.1** — an Admission Result instance identifier.

- The `admission_id` is an **Admission evidence identifier** only.
- The `admission_id` is **not** a DigiEmu Core state hash, snapshot hash, bundle hash, or any alternative state-identity mechanism.
- The `admission_id` does **not** grant business, legal, regulatory, ethical, or trust authority.
- The `intent_digest` is an **Admission intent normalization artifact** only. It is not a Core state identity.
- Neither derivation depends on random UUIDs, timestamps, process IDs, hostnames, filesystem paths, local machine state, or unordered map serialization.

## 2. Profiles

```
Intent Envelope (normalized)
        ↓
P0.ADMISSION.INTENT.v0.1 canonical input
        ↓
UTF-8 compact canonical JSON
        ↓
SHA-256
        ↓
p0-intent:sha256:<lowercase-hex>
        ↓
Admission Result canonical input
        ↓
P0.ADMISSION.ID.v0.1 canonical input
        ↓
UTF-8 compact canonical JSON
        ↓
SHA-256
        ↓
admission:sha256:<lowercase-hex>
```

SHA-256 is consistent with the existing DigiEmu Core snapshot-hashing practice documented in `docs/SNAPSHOT_HASH_v1.0.md`. The `p0-intent:sha256:` and `admission:sha256:` prefixes keep these identifiers outside the Core state-identity namespace.

## 3. P0.ADMISSION.INTENT.v0.1 — Intent digest

### 3.1 Purpose

The `intent_digest` provides a deterministic, P0-specific, Admission-bound identity for the normalized content of an Intent Envelope. It is computed **without** using the `intent_id`, because the normative derivation of `intent_id` is not yet defined and must not contaminate deterministic intent identity.

### 3.2 Canonical input fields

Top-level field order, fixed:

1. `intent_digest_profile`
2. `schema_version`
3. `architecture_revision`
4. `capability_ref`
5. `aggregate_ref`
6. `command_ref`
7. `payload`

### 3.3 Exclusions

The `intent_id` is **excluded** from the canonical input. It may be arbitrary until its derivation is normatively defined.

### 3.4 Canonicalization rules

- **Encoding:** UTF-8 without a leading BOM.
- **Whitespace:** No insignificant whitespace, no line breaks, no indentation. One compact JSON line.
- **Top-level field order:** Fixed as listed in §3.2.
- **Object key ordering:** Object keys are sorted lexicographically **at every nesting level**, including inside `payload`.
- **Array order:** Preserved exactly as supplied.
- **Scalar types:** Preserved exactly.
- **Empty payload:** `payload` MUST be an object. An intent with no command-specific data uses:

```json
"payload": {}
```

- **No null payload:** Do not introduce `null` payload semantics in v0.1.

### 3.5 Hash and identifier format

- SHA-256 of the UTF-8 bytes of the canonical JSON string.
- 64-character lowercase hexadecimal digest.
- Identifier format:

```text
p0-intent:sha256:<64-char-hex>
```

## 4. P0.ADMISSION.ID.v0.1 — Admission ID

### 4.1 Purpose

The `admission_id` identifies one concrete, normalized Admission evaluation. It includes the normalized `intent_digest` plus the resolved Admission result fields, so materially different requested mutations produce different `admission_id` values.

### 4.2 Canonical input fields

Top-level field order, fixed:

1. `admission_id_profile`
2. `schema_version`
3. `architecture_revision`
4. `intent_digest`
5. `capability_ref`
6. `aggregate_ref`
7. `command_ref`
8. `transition_ref`
9. `decision`
10. `rule_refs`
11. `reason_codes`

### 4.3 `transition_ref` treatment

- **ADMIT:** `transition_ref` MUST be a resolved, non-empty transition identifier.
- **REJECT before transition resolution:** The `admission_id` canonical input uses `"transition_ref": null`. This `null` is part of the canonical canonical input only and means "no transition resolved".

The external Admission Result object may omit `transition_ref` for such a REJECT result (see `schemas/admission_result_v0.2.schema.json`). The canonical `admission_id` preimage normalizes that absence to `null`.

### 4.4 Canonicalization rules

- **Encoding:** UTF-8 without BOM.
- **Whitespace:** Compact, single-line JSON.
- **Field order:** Fixed as listed in §4.2.
- **`rule_refs`:** Sorted lexicographically by rule identifier before serialization.
- **`reason_codes`:** Sorted lexicographically by reason code before serialization.
- **Duplicates:** Prohibited by the Admission Result schema.
- **`admission_id` self-exclusion:** `admission_id` is never part of its own canonical input.

### 4.5 Hash and identifier format

- SHA-256 of the UTF-8 bytes of the canonical JSON string.
- 64-character lowercase hexadecimal digest.
- Identifier format:

```text
admission:sha256:<64-char-hex>
```

## 5. Versioning

Both profiles carry a `*_profile` field in their canonical input. If any of the canonical field set, ordering, hash algorithm, or prefix changes, a new profile value (e.g. `P0.ADMISSION.INTENT.v0.2` or `P0.ADMISSION.ID.v0.2`) must be used. The profile value is itself part of the canonical input, so identifiers produced by different profiles are inherently different.

## 6. Normative examples

### 6.1 Empty-payload ADMIT

**Canonical intent preimage:**

```text
{"intent_digest_profile":"P0.ADMISSION.INTENT.v0.1","schema_version":"v0.1","architecture_revision":"0.3","capability_ref":"core.unit.create","aggregate_ref":"unit","command_ref":"unit.create","payload":{}}
```

**Intent digest:**

```text
p0-intent:sha256:0dd513c96504e4a180ae0dcdab2a3746cf32cc2b6f697407f3cc04f18626ef8e
```

**Canonical admission preimage:**

```text
{"admission_id_profile":"P0.ADMISSION.ID.v0.1","schema_version":"v0.1","architecture_revision":"0.3","intent_digest":"p0-intent:sha256:0dd513c96504e4a180ae0dcdab2a3746cf32cc2b6f697407f3cc04f18626ef8e","capability_ref":"core.unit.create","aggregate_ref":"unit","command_ref":"unit.create","transition_ref":"unit:created","decision":"ADMIT","rule_refs":["P0.ADMISSION.AGGREGATE_OWNS_CAPABILITY","P0.ADMISSION.ARCHITECTURE_REVISION","P0.ADMISSION.CAPABILITY_EXISTS","P0.ADMISSION.CAPABILITY_MUTATES","P0.ADMISSION.COMMAND_AGGREGATE_MATCH","P0.ADMISSION.COMMAND_CAPABILITY_MATCH","P0.ADMISSION.COMMAND_EXISTS","P0.ADMISSION.COMMAND_TRANSITION_DEFINED","P0.ADMISSION.INTENT_REQUIRED_FIELDS"],"reason_codes":[]}
```

**Admission ID:**

```text
admission:sha256:f18917cd7d2dd6dbcbaf0217ba5df0f28cd3544c99eff792852961ff0327cee9
```

### 6.2 Payload "alpha" ADMIT

**Canonical intent preimage:**

```text
{"intent_digest_profile":"P0.ADMISSION.INTENT.v0.1","schema_version":"v0.1","architecture_revision":"0.3","capability_ref":"core.unit.create","aggregate_ref":"unit","command_ref":"unit.create","payload":{"extra":"x","key":"alpha"}}
```

**Intent digest:**

```text
p0-intent:sha256:a5bc1e21afb34ed1b358229935ef1e4ca4c99cdf5f3638688a72e108853e0294
```

**Admission ID:**

```text
admission:sha256:7f87382ff543c0b40539a05c70a2fad82fcd046a5989215296e8f8353f41ad4f
```

### 6.3 Payload "beta" ADMIT

**Canonical intent preimage:**

```text
{"intent_digest_profile":"P0.ADMISSION.INTENT.v0.1","schema_version":"v0.1","architecture_revision":"0.3","capability_ref":"core.unit.create","aggregate_ref":"unit","command_ref":"unit.create","payload":{"extra":"x","key":"beta"}}
```

**Intent digest:**

```text
p0-intent:sha256:2aaefddc9882b0894854c30abc0bff11ed7bd1008de0c5018b0fcbc2057112a0
```

**Admission ID:**

```text
admission:sha256:4ac7d1f1344b4621d6372a91ade7e81e46077ba6559a0f141e94e6a206bbd33c
```

### 6.4 REJECT (architecture revision mismatch)

**Canonical intent preimage:**

```text
{"intent_digest_profile":"P0.ADMISSION.INTENT.v0.1","schema_version":"v0.1","architecture_revision":"0.3","capability_ref":"core.unit.create","aggregate_ref":"unit","command_ref":"unit.create","payload":{}}
```

**Intent digest:**

```text
p0-intent:sha256:0dd513c96504e4a180ae0dcdab2a3746cf32cc2b6f697407f3cc04f18626ef8e
```

**Canonical admission preimage:**

```text
{"admission_id_profile":"P0.ADMISSION.ID.v0.1","schema_version":"v0.1","architecture_revision":"0.3","intent_digest":"p0-intent:sha256:0dd513c96504e4a180ae0dcdab2a3746cf32cc2b6f697407f3cc04f18626ef8e","capability_ref":"core.unit.create","aggregate_ref":"unit","command_ref":"unit.create","transition_ref":null,"decision":"REJECT","rule_refs":["P0.ADMISSION.ARCHITECTURE_REVISION"],"reason_codes":["ARCHITECTURE_REVISION_MISMATCH"]}
```

**Admission ID:**

```text
admission:sha256:c3d7f510ac9e36799aeecea82ce2fc33a478e2b71b7792358fdf00738cfcb433
```

## 7. Determinism validation

| Test | Result |
|------|--------|
| Same empty-payload ADMIT input twice | Identical `admission:sha256:f18917cd...` |
| Payload `alpha` vs `beta` | Different admission IDs (`7f87382f...` vs `4ac7d1f1...`) and different intent digests |
| Same payload with different object property ordering | Produces the same canonical intent digest, because object keys are sorted at every nesting level. |
| `rule_refs` in different incidental order | Sorted before hashing, producing the same admission ID. |
| `reason_codes` in different incidental order | Sorted before hashing, producing the same admission ID. |
| `architecture_revision` changed from `0.3` to `0.4` | Different admission ID (`f18917cd...` vs `admission:sha256:756a55e7bd8d616ef2703ce8c4c93810bd074f87dbe470dde1a9c782158f1b70`). |

## 8. Relation to DigiEmu Core state identity

The prefixes `p0-intent:sha256:` and `admission:sha256:` are deliberately outside the Core state-identity namespace. These identifiers are valid only as Admission evidence artifacts. They must never be used as:

- snapshot hashes,
- bundle hashes,
- Unit or Version identity,
- replay roots, or
- replacements for any canonicalization profile in `docs/SNAPSHOT_HASH_v1.0.md`.

## 9. Future integration notes

- `admission-rule-registry.yaml` is the source of `rule_refs` values.
- `admission_result_v0.2.schema.json` is the first Admission Result schema that supports `transition_ref` omission for unresolved REJECT cases.
- `testdata/rs_001/run.ps1` should derive `admission_id` from the canonical inputs above, not from fixture-local strings.
