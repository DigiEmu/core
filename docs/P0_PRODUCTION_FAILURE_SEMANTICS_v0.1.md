# P0 Production Failure Semantics

- **Document:** P0 Production Failure Semantics
- **Version:** v0.1
- **Status:** PROPOSED
- **Date:** 2026-08-17
- **Architecture baseline:** `0.3`
- **Phase D freeze:** `776ccf1`
- **Scope:** Conceptual failure taxonomy for a future production Admission path. This document does not change any runtime behavior, schema, registry, or existing Core semantics.

## 1. Scope

This document defines a minimal, precise conceptual failure model for the future production path:

```
External Input
    ↓
Envelope / Input Validation
    ↓
Semantic/Profile Support Resolution
    ↓
Admission Evaluation
    ↓
ADMIT or REJECT
    ↓
Pre-execution Evidence / Command Preparation
    ↓
Execution
    ↓
Persistence / Audit
    ↓
Post-execution Evidence
    ↓
Response / Result
```

The categories defined here are conceptual. They may be represented by existing types, new internal types, or documentation only. The purpose is to prevent distinct technical meanings from being collapsed into a single `REJECT` or `ERROR` condition.

## 2. Failure stages

### 2.1 Input / Envelope Validation Failure

**Definition:** the received input cannot be parsed, decoded, or normalized into a syntactically valid `admission.Intent`.

**Examples:**
- Malformed JSON.
- Invalid JSON structure for an `Intent Envelope`.
- Field type that cannot be unmarshalled.
- Syntactically invalid envelope.

**Boundary with Admission:**
- The Admission Engine **is not reached**.
- No `admission.Result` is produced.
- The existing `P0.ADMISSION.INTENT_REQUIRED_FIELDS` rule is **not** an input validation failure. It is an Admission rule that evaluates a structurally valid `Intent` and returns `REJECT` with `MISSING_REQUIRED_FIELD` if a required field is empty.

**Distinction:**
- **A. Input so malformed that no valid `admission.Intent` exists** → `INPUT / ENVELOPE VALIDATION FAILURE`.
- **B. Structurally valid `admission.Intent` that reaches Admission and fails `P0.ADMISSION.INTENT_REQUIRED_FIELDS`** → `ADMISSION REJECT` with `MISSING_REQUIRED_FIELD`.

**Mutation:** no mutation occurs.

### 2.2 Unsupported Semantics

**Definition:** the implementation recognizes the artifact as a supported shape but does not understand the requested semantic profile or version.

**Examples:**
- Unknown `schema_version` whose structure cannot be interpreted.
- Unknown deterministic identity/profile version for which the implementation has no evaluation semantics (e.g., a future `p0-intent:sha3-256:` when only `sha256` is implemented).
- Unknown canonicalization or hashing profile.

**Boundary with Admission REJECT:**
- `P0.ADMISSION.ARCHITECTURE_REVISION` checks whether the `Intent`'s `architecture_revision` equals the **currently known** baseline. A mismatch against the known baseline is `ADMISSION REJECT` with `ARCHITECTURE_REVISION_MISMATCH`.
- `UNSUPPORTED` applies when the implementation cannot interpret the normative schema or profile semantics required to construct or evaluate the request at all (e.g., unknown `schema_version`, unknown deterministic identity profile, unknown canonicalization or hashing profile). It is **not** a normative `REJECT` merely because evaluation is impossible.

**Mutation:** no mutation occurs.

**Note:** this document does not introduce a public `UNSUPPORTED` enum. It is a conceptual classification.

### 2.3 Admission REJECT

**Definition:** the Admission Engine successfully evaluated a structurally valid `Intent` under a known rule/configuration set and determined that it is architecturally inadmissible.

**Handler:** the Core use-case handler **MUST NOT** run.

**Existing reason codes (from `admission-rule-registry.yaml` v0.1):**

| Rule | Reason code |
|---|---|
| `P0.ADMISSION.ARCHITECTURE_REVISION` | `ARCHITECTURE_REVISION_MISMATCH` |
| `P0.ADMISSION.CAPABILITY_EXISTS` | `UNKNOWN_CAPABILITY` |
| `P0.ADMISSION.CAPABILITY_MUTATES` | `CAPABILITY_NOT_MUTATING` |
| `P0.ADMISSION.AGGREGATE_OWNS_CAPABILITY` | `OWNERSHIP_MISMATCH` |
| `P0.ADMISSION.COMMAND_EXISTS` | `UNKNOWN_COMMAND` |
| `P0.ADMISSION.COMMAND_CAPABILITY_MATCH` | `COMMAND_CAPABILITY_MISMATCH` |
| `P0.ADMISSION.COMMAND_AGGREGATE_MATCH` | `COMMAND_AGGREGATE_MISMATCH` |
| `P0.ADMISSION.COMMAND_TRANSITION_DEFINED` | `UNDEFINED_TRANSITION` |
| `P0.ADMISSION.INTENT_REQUIRED_FIELDS` | `MISSING_REQUIRED_FIELD` |

**Invariants:**
- `REJECT` is a **normative architectural** outcome, not a runtime exception.
- `REJECT` means the input is inadmissible under known rules.
- `REJECT` does not imply implementation failure.
- The rejected `admission.Result` is stable and deterministic for the same `Intent` and `Config`.

### 2.4 System / Evaluation Error

**Definition:** the implementation cannot evaluate the `Intent` because of an internal or environmental condition that is not an architectural inadmissibility.

**Examples:**
- Required Admission `Config` is unavailable.
- Registry or `Config` initialization failed.
- `ComputeIntentDigest` or `ComputeAdmissionID` failed.
- Engine internal invariant violation.
- Implementation cannot reliably determine admissibility.

**Boundary with REJECT:**
- This is **not** a normative `REJECT` because the architecture did not successfully determine inadmissibility.
- The handler **MUST NOT** run.
- No new Admission reason code should be used for infrastructure failures.

**Result:** an error, not an `admission.Result` with `REJECT`.

### 2.5 Execution / Domain Error

**Definition:** the `Intent` was admitted, but the existing Core use case could not complete successfully.

**Examples:**
- `ErrUnitAlreadyExists` for `core.unit.create`.
- `ErrUnitNotFound` or `ErrVersionNotFound`.
- `ErrConflict` (optimistic locking).
- Invalid payload (`meaning.json too large`, unsupported `schema_version`).
- Repository failure before all required persistence steps succeed.

**Invariants:**
- Admission succeeded; the `admission.Result` remains `ADMIT`.
- A later execution error **does not retroactively turn `ADMIT` into `REJECT`**.
- The coherent Core mutation is **not confirmed** if the handler returns an error before all required persistence steps succeed; partial persistence side effects may remain.

### 2.6 Persistence / Audit Partial Failure

**Definition:** the state-mutating write succeeded, but a later step (audit append or related) failed.

**Current code evidence for all five mutating use cases:**

| Mutation | State write | Audit append | State before audit? |
|---|---|---|---|
| `core.unit.create` | `SaveUnit(u)` | `Audit.Append(ev)` | **Yes** |
| `core.version.create` | `SaveVersion(v)` + `UpdateUnitHead` | `Audit.Append(ev)` | **Yes** |
| `core.meaning.set` | `SaveMeaning(...)` | `Audit.Append(ev)` | **Yes** |
| `core.claim.set` | `SaveClaimSet(...)` | `Audit.Append(ev)` | **Yes** |
| `core.uncertainty.set` | `SaveUncertainty(...)` | `Audit.Append(ev)` | **Yes** |

**Key finding:** in every current mutating use case, the state sidecar or record is persisted **before** `Audit.Append`. If `Audit.Append` returns an error, the handler returns an error, but the state mutation has already occurred. The audit event is not appended.

**Partial-failure state:**
- Mutation: **may have occurred**.
- Audit: **missing**.
- Response: error.

**Critical principle:**
- An error return from the handler does **not** prove that no mutation occurred.
- A 2046 auditor cannot conclude the operation was rolled back simply because the handler returned an error.

**Non-decision:**
- This document does not prescribe transactions, tombstones, compensating records, retries, or a new storage abstraction.
- A separate **P0 Atomicity / Recovery Contract** must decide the minimum appropriate solution before production enforcement.

### 2.7 Evidence Construction Error

**Definition:** Admission succeeded, the mutation completed, but construction or validation of a `Command Envelope` or `Event Envelope` failed.

**Examples:**
- `Command Envelope` construction error after `ADMIT` but before execution.
- `Event Envelope` construction error after the use case produced an `AuditEvent`.
- Evidence serialization or schema validation failure.

**Invariants:**
- If the `Command Envelope` construction fails **before** the handler runs, the mutation has not occurred.
- If the `Event Envelope` construction fails **after** the `AuditEvent` is emitted, the mutation has still occurred and the `AuditEvent` remains the source-of-truth runtime event.
- Evidence failure **must not rewrite execution history**.

## 3. Conceptual state model

```
INPUT RECEIVED
    ↓
INPUT VALID?
    ├── NO  → INPUT / ENVELOPE VALIDATION FAILURE
    │           (no Admission, no handler)
    ↓ YES
SEMANTICS SUPPORTED?
    ├── NO  → UNSUPPORTED / SYSTEM CONDITION
    │           (no Admission, no handler)
    ↓ YES
ADMISSION EVALUATION
    ├── ERROR → SYSTEM / EVALUATION ERROR
    │           (no handler)
    ├── REJECT → ADMISSION REJECT
    │           (handler MUST NOT run)
    ↓ ADMIT
PRE-EXECUTION EVIDENCE / COMMAND PREPARATION
    ├── ERROR → PRE-EXECUTION EVIDENCE/SYSTEM ERROR
    │           (handler not yet invoked; no mutation)
    ↓
EXECUTION
    ├── DOMAIN/PERSISTENCE ERROR → EXECUTION ERROR
    │           (coherent Core mutation not confirmed; partial side effects possible)
    ├── PARTIAL FAILURE → MUTATION MAY HAVE OCCURRED
    ↓ SUCCESS
AUDIT APPEND
    ├── ERROR AFTER MUTATION → PARTIAL FAILURE
    │           (state persists; audit missing)
    ↓ SUCCESS
POST-EXECUTION EVIDENCE
    ├── ERROR → EVIDENCE ERROR
    │           (mutation complete; evidence missing)
    ↓
SUCCESS
```

## 4. Mutation certainty table

| Failure category | Admission outcome exists? | Handler invoked? | Mutation certainty | Audit certainty | Evidence certainty | Retry risk |
|---|---|---|---|---|---|---|
| Input / Envelope Validation | No | No | Not started | Not applicable | Not applicable | Safe to re-submit |
| Unsupported | No (or implementation-defined) | No | Not started | Not applicable | Not applicable | Requires new implementation |
| Admission REJECT | Yes (`REJECT`) | No | Not started | Not applicable | Not applicable | Re-submit corrected input |
| System / Evaluation Error | No (error) | No | Not started | Not applicable | Not applicable | Implementation may retry internally |
| Execution / Domain Error | Yes (`ADMIT`) | Attempted | Coherent mutation not confirmed; partial side effects possible | Not applicable | Not generated | Requires state check |
| Partial (audit after state) | Yes (`ADMIT`) | Yes | **May have occurred** | **Missing** | May not exist | **Requires state check** |
| Evidence Construction | Yes (`ADMIT`) | Yes | Confirmed | Confirmed | Missing/invalid | Evidence only; mutation complete |
| Success | Yes (`ADMIT`) | Yes | Confirmed | Confirmed | Confirmed | Safe to re-run only with idempotency |

## 5. Admission Result semantics

- **Input / Envelope Validation Failure:** no `admission.Result`.
- **Unsupported / System Condition:** no normative `admission.Result`.
- **Admission REJECT:** `admission.Result` with `decision: REJECT`.
- **System / Evaluation Error:** an error, not a `Result`.
- **Execution Error:** the existing `admission.Result` remains `ADMIT`.
- **Partial Failure:** the existing `admission.Result` remains `ADMIT`.
- **Evidence Construction Error:** the existing `admission.Result` remains `ADMIT`.

## 6. Error does not rewrite history

The following principles apply:

- `ADMIT` followed by an execution error does **not** become `REJECT`.
- `ADMIT` followed by a partial audit failure does **not** become `REJECT`.
- A successful mutation followed by an `Event Envelope` failure does **not** mean execution failed.
- An error return from a use case does **not** prove that no state mutation occurred.
- `REJECT` means architectural inadmissibility under known rules, not runtime failure.

## 7. Retry consequences

| Condition | Retry assessment |
|---|---|
| Response lost before Admission evaluation | Re-submit is safe; Admission is deterministic. |
| Response lost after `REJECT` | Re-submit corrected input; same `admission_id` if input unchanged. |
| Response lost after `ADMIT` | Re-submit may cause duplicate execution; requires idempotency policy. |
| Audit failure after state persistence | Blind retry may re-execute; **requires state check**. |
| Evidence construction failure | Mutation complete; re-constructing evidence is safe. |
| Client disconnect after execution | State and audit may be complete; client must reconcile. |
| Process restart after uncertain completion | No durable in-progress marker currently exists; **requires recovery policy**. |

**Note:** `admission_id` is not defined as an idempotency key. `command_id` and `event_id` are not defined in this document.

## 8. Identity non-decisions

This document does **not** decide:

- `intent_id` meaning or derivation.
- `command_id` meaning or derivation.
- `event_id` meaning or derivation.
- Whether `event_id` equals `AuditEvent.ID`.
- Idempotency-key semantics.
- Execution-attempt identity.

Failure semantics must be correct independently of later ID choices.

## 9. Evidence non-decisions

This document does **not** decide:

- Whether `Event Envelope` is persisted.
- Whether `Command Envelope` is persisted.
- Whether a compact evidence manifest is persisted.
- Whether a complete evidence bundle is persisted.
- The long-term archive format.

It only defines what an evidence system must not misrepresent.

## 10. Atomicity non-decision

This document identifies current atomicity/failure semantics. It does **not** choose between:

- transactions,
- compensating records,
- tombstones,
- a new `AuditEvent.Type`,
- recovery journal,
- retry protocol,
- rollback,
- atomic adapter,
- or other future mechanisms.

A separate **P0 Atomicity / Recovery Contract** will evaluate the minimum appropriate solution.

## 11. Transport independence

The failure taxonomy is transport-neutral. This document does **not** define:

- HTTP status codes,
- CLI exit codes,
- CLI text,
- HTTP JSON response shapes.

Those mappings are a future transport-layer concern.

## 12. Minimal failure taxonomy

| Category | Representation | Notes |
|---|---|---|
| `INPUT_ERROR` | Conceptual; may map to a transport/schema error | Pre-Admission |
| `UNSUPPORTED` | Conceptual | Implementation cannot interpret profile/version |
| `ADMISSION_REJECT` | Already public (`admission.Result` `decision: REJECT`) | Normative architectural inadmissibility |
| `SYSTEM_ERROR` | Conceptual; implementation may use `error` | Cannot evaluate Admission |
| `EXECUTION_ERROR` | Conceptual; implementation may use `error` | Use case failed before or during required persistence steps |
| `PARTIAL_FAILURE` | Conceptual | State persisted, audit/evidence may be missing |
| `EVIDENCE_ERROR` | Conceptual | Evidence construction failed after mutation |

## 13. Current five mutations

| Mutation | State-write order | Audit order | State before audit? | Handler return behavior | Current ambiguity | Production implication |
|---|---|---|---|---|---|---|
| `core.unit.create` | `NewUnit` → `ExistsByKey` → `SaveUnit` | `Audit.Append` after `SaveUnit` | Yes | Returns error if `SaveUnit` or `Audit.Append` fails | Cannot distinguish `SaveUnit` failure from `Audit.Append` failure by error alone | A `unit` may exist without `unit.created` audit |
| `core.version.create` | `SaveVersion` + `UpdateUnitHead` | `Audit.Append` after state updates | Yes | Returns error if state or audit fails | `UpdateUnitHead` may succeed while `Audit.Append` fails | Version and head may be updated without `version.created` audit |
| `core.meaning.set` | `SaveMeaning(...)` | `Audit.Append` after `SaveMeaning` | Yes | Returns error if `SaveMeaning` or `Audit.Append` fails | A `SaveMeaning` failure returns the same kind of error as an `Audit.Append` failure | Meaning may be saved without `MEANING_SET` audit |
| `core.claim.set` | `SaveClaimSet(...)` | `Audit.Append` after `SaveClaimSet` | Yes | Returns error if `SaveClaimSet` or `Audit.Append` fails | Same as above | Claim set may be saved without `CLAIM_SET` audit |
| `core.uncertainty.set` | `SaveUncertainty(...)` | `Audit.Append` after `SaveUncertainty` | Yes | Returns error if `SaveUncertainty` or `Audit.Append` fails | Same as above | Uncertainty may be saved without `UNCERTAINTY_SET` audit |

## 14. Current vs future responsibility

### A. Current 2026 semantic debt

- State and audit are not atomic.
- Error returns do not disambiguate pre-commit failures from audit failures.
- No in-progress or recovery marker exists for crashes.

### B. Phase E requirement

- Documented failure taxonomy (this document).
- `admission.Result` semantics preserved under all failure modes.
- `ADMIT` is not conflated with execution success.
- Partial-failure states are recognized in orchestration design.

### C. Future / deferred

- Concrete atomicity/recovery mechanism.
- Idempotency/retry policy.
- Evidence persistence format.
- `command_id` / `event_id` derivations.

## 15. Relation to Phase D

Phase D remains frozen. This document does not change:

- IR statuses in the conformance matrix.
- `internal/admission` Engine semantics.
- `admission-rule-registry.yaml` v0.1.
- Phase D binding results.
- `P0.ADMISSION.INTENT.v0.1` or `P0.ADMISSION.ID.v0.1`.
- Core state identity, canonicalization, hashing, replay, or verification.
- `AuditEvent` semantics.

This document clarifies behavior required before Phase E production wiring.

## 16. Relation to Phase E

| Failure semantics required before | Item |
|---|---|
| **Phase E architecture design** | This document (Part 2–12). |
| **Phase E test-only orchestration** | State/audit ordering table (Part 13); partial-failure recognition. |
| **First production Admission gate** | Atomicity / Recovery Contract (separate). |
| **General release** | Concrete transport mappings and operational runbook. |

## 17. Open questions

- Minimum `P0 Atomicity / Recovery` solution for state/audit partial failures.
- Exact mutation certainty after each `fs` adapter failure (not yet exhaustively tested).
- Retry / idempotency semantics.
- Evidence persistence model.
- `command` occurrence identity.
- `event_id` semantics and relationship to `AuditEvent.ID`.
- Transport mapping of conceptual categories to HTTP status codes / CLI exit codes.

## 18. Explicit non-claims

This document does **not** claim that:

- Production Admission is enabled.
- Mutation/audit atomicity currently exists.
- Retries are currently idempotent.
- `Event Envelope` persistence exists.
- Full evidence persistence exists.
- Production recovery exists.
- `command_id` semantics are defined.
- `event_id` semantics are defined.
- Historical verification is complete.
- Any new reason code beyond the nine v0.1 codes is defined.

## 19. Document principles

1. `REJECT` means known architectural inadmissibility, not generic failure.
2. Pre-Admission validation failure is distinct from `ADMISSION REJECT`.
3. Unsupported semantics are not guessed.
4. System inability to evaluate is distinct from architectural inadmissibility.
5. `ADMIT` remains `ADMIT` even if later execution fails.
6. Execution failure does not imply Admission failure.
7. An error after mutation must not falsely imply no mutation occurred.
8. Evidence failure must not rewrite execution history.
9. Retry safety depends on mutation certainty.
10. Failure semantics remain transport-independent.

## 20. Git verification statement

This document is a new proposed specification. No existing files have been modified.
