# P0 Atomicity & Recovery Contract

- **Document:** P0 Atomicity & Recovery Contract
- **Version:** v0.1
- **Status:** PROPOSED
- **Date:** 2026-08-17
- **Architecture baseline:** `0.3`
- **Phase D freeze:** `776ccf1`
- **Production Failure Semantics:** `fab9882`
- **Scope:** Minimum truthful atomicity and recovery contract for a future Production Admission gate. This document does not change any runtime behavior, schema, registry, or existing Core semantics.

## 1. Contract scope

This contract covers the current filesystem mutation semantics, mutation-certainty rules, Core-state acknowledgement, partial-failure classification, retry preconditions, and recovery decision boundaries required before Production Admission may protect Core mutations.

This contract does **not** define:

- production orchestrator implementation,
- CLI/HTTP mappings,
- database or distributed transactions,
- WAL, saga, or recovery journal,
- generic transaction abstraction,
- new persistence backend,
- new `AuditEvent.Type` values,
- idempotency-key design,
- `command_id` / `event_id` semantics,
- long-term historical archive format.

## 2. Four distinct persistence properties

The terms below are defined separately and are **not** interchangeable.

| Property | Definition |
|---|---|
| **Operation atomicity** | The entire replacement or update appears either completed or not from the perspective of an immediate observer. `os.Rename` is **not guaranteed to be atomic on non-Unix platforms** per the Go `os` package contract. |
| **Crash durability** | Data is preserved after a process or OS crash. The current implementation does **not** call `fsync` / `Sync` after `Save*` or `Audit.Append`. Durability is **not guaranteed**. |
| **Concurrent-write integrity** | Concurrent operations do not interleave or corrupt records. The filesystem adapter serializes per-`UnitRepo` operations. `AuditLog.Append` has **no application-level lock** and concurrent calls are **not safe**. |
| **Multi-step logical consistency** | All required steps of a logical mutation succeed together. Current use cases perform multiple `Save*` / `Audit.Append` calls in sequence; these are **not atomic** as a group. |

### Current proven behavior by step

| Step | Operation atomicity | Crash durability | Concurrent-write integrity | Multi-step logical consistency |
|---|---|---|---|---|
| `unit.json` replacement | Non-atomic on Windows; best-effort single-file replacement | No `fsync` | Per-`UnitRepo` mutex | Not applicable (single file) |
| `version` append + `unit.json` rewrite | Same as above | No `fsync` | Per-`UnitRepo` mutex | **Not consistent** with `UpdateUnitHead` as a single step |
| `HeadVersionID` update (`UpdateUnitHead`) | Same as above | No `fsync` | Per-`UnitRepo` mutex | Separate rewrite from `SaveVersion` |
| Sidecar replacement (`SaveMeaning` / `SaveClaimSet` / `SaveUncertainty`) | Same as above | No `fsync` | Per-`UnitRepo` mutex | Sidecar rename and unit-record rewrite are separate steps |
| `UnitRecord` hash/reference update | Same as above | No `fsync` | Per-`UnitRepo` mutex | Follows sidecar rename; may fail separately |
| Index update | Best-effort `os.Rename` of `unitsByKey.json` | No `fsync` | `indexStore` mutex | Not tied to `SaveUnit` return value |
| `AuditLog.Append` | One `os.File.Write`; no record framing | No `fsync` / `Sync` | **Not safe** under concurrent writers | Not tied to `Save*` success |
| State + audit together | No joint operation | No `fsync` | Not jointly protected | **Not atomic** |

## 3. Current file replacement contract

The filesystem adapter uses a temporary-file plus `os.Rename` pattern. The contract is:

- The Go `os.Rename` documentation states: *"Even within the same directory, on non-Unix platforms Rename is not an atomic operation."*
- DigiEmu Core does **not** guarantee cross-platform atomic file replacement.
- Successful method return does **not** imply crash durability.
- No `fsync` / `Sync` durability guarantee currently exists.
- No new filesystem implementation is introduced by this contract.
- No Unix-only behavior is required or assumed.

A `Save*` method that returns `error` means the coherent Core state is **not confirmed**. It does **not** prove that zero filesystem effects occurred.

## 4. Core state acknowledgement

Current Core code acknowledges the following as coherent Core state:

- `UnitRecord` stored in `unit.json`.
- `VersionRecord` stored inside `UnitRecord`.
- `HeadVersionID` stored in `UnitRecord`.
- `MeaningHash`, `ClaimSetHash`, and `UncertaintyHash` stored in `VersionRecord`.

**Acknowledgement rule:**

A sidecar file that exists on disk but is not referenced by the corresponding `VersionRecord` **must not** be treated as a confirmed coherent Core mutation. The sidecar is a persistent artifact; only the referenced `VersionRecord` field makes it Core-acknowledged.

The index is best-effort and **must not** be treated as authoritative Core state.

The `audit.ndjson` log **must not** be treated as the authority for whether Core state exists.

The handler return value is an in-process signal; it is not persistent and is not the authority for state.

## 5. Index contract

Actual current behavior:

- `SaveUnit` persists the `unit.json` record.
- `index.upsertUnitKey` is best-effort.
- Index persistence errors are **not propagated** by `SaveUnit`.
- `ExistsByKey` and `FindUnitByKey` have directory-scan fallback behavior.

Therefore:

- Index inconsistency alone must not be treated as proof that a `Unit` does not exist.
- A missing or stale index does not invalidate a persisted `unit.json`.
- Index repair is not required by this contract.

## 6. Mutation certainty contract

The following conceptual certainty classes are used. They are documentation categories, not public runtime identifiers.

| Class | Meaning |
|---|---|
| **A. MUTATION NOT STARTED** | The handler/persistence has not begun. Applicable to validation and pre-`Save*` errors. |
| **B. COHERENT MUTATION NOT CONFIRMED** | A persistence operation returned an error. Partial filesystem effects may exist. The intended coherent Core state is not confirmed. |
| **C. PARTIAL MUTATION CONFIRMED** | Repository evidence proves that some required steps completed while later required steps failed. |
| **D. COHERENT CORE STATE CONFIRMED** | The required Core state writes completed according to the current adapter behavior (`UnitRecord` / `VersionRecord` fields are updated). |
| **E. AUDIT COMPLETENESS NOT CONFIRMED** | Coherent Core state may be confirmed while the audit append failed or its durability is unknown. |

No new public enum is created by this contract.

## 7. Five mutation contracts

### `core.unit.create`

| Aspect | Contract |
|---|---|
| Required coherent-state condition | `unit.json` file exists with the expected `UnitRecord`. |
| Known partial-failure boundary | `unit.json` may exist while `unit.created` audit is missing; `SaveUnit` returns `nil` even if index update is best-effort. |
| Audit boundary | `Audit.Append` is separate; a failure after `SaveUnit` is `PARTIAL_FAILURE`. |
| Blind retry allowed? | **No.** A retry with the same key will fail with `ErrUnitAlreadyExists`. The caller cannot determine whether the unit was created by this or a prior attempt without inspecting state. |
| Required pre-retry inspection | Check whether `unit.json` exists for the deterministic `Unit.ID` and whether `unit.created` is present in the audit log. |

### `core.version.create`

| Aspect | Contract |
|---|---|
| Required coherent-state condition | `UnitRecord` contains the `VersionRecord` and `HeadVersionID` equals the new `Version.ID`. |
| Known partial-failure boundary | A `VersionRecord` may be appended while `HeadVersionID` is not updated; this is a partial mutation. A `VersionRecord` and `HeadVersionID` may be updated while `version.created` audit is missing. |
| Audit boundary | `Audit.Append` is separate. |
| Blind retry allowed? | **No.** `Version.ID` is deterministic from `unitID + label + content`; a retry can append a duplicate `VersionRecord`. |
| Required pre-retry inspection | Check whether the deterministic `VersionRecord` exists in `UnitRecord` and whether `HeadVersionID` references it. |

### `core.meaning.set`

| Aspect | Contract |
|---|---|
| Required coherent-state condition | The `MeaningHash` field in the `VersionRecord` is set and the `meaning.json` sidecar exists at the expected path. |
| Known partial-failure boundary | `meaning.json` sidecar may be present while `MeaningHash` in `UnitRecord` is not updated; audit may be missing after both are present. |
| Audit boundary | `Audit.Append` is separate. |
| Blind retry allowed? | **No.** A retry may rewrite the sidecar and append a duplicate `MEANING_SET` audit event (same `AuditEvent.ID` for the same content). |
| Required pre-retry inspection | Check `VersionRecord.MeaningHash` and the existence of the sidecar. Both must agree. |

### `core.claim.set`

| Aspect | Contract |
|---|---|
| Required coherent-state condition | `ClaimSetHash` field in `VersionRecord` is set and `claimset.json` sidecar exists. |
| Known partial-failure boundary | Same as `SetMeaning`, with `ClaimSetHash` and `claimset.json`. |
| Audit boundary | Same as `SetMeaning`. |
| Blind retry allowed? | **No.** Same as `SetMeaning`. |
| Required pre-retry inspection | Check `VersionRecord.ClaimSetHash` and `claimset.json` existence. |

### `core.uncertainty.set`

| Aspect | Contract |
|---|---|
| Required coherent-state condition | `UncertaintyHash` field in `VersionRecord` is set and `uncertainty.json` sidecar exists. |
| Known partial-failure boundary | Same as `SetMeaning`, with `UncertaintyHash` and `uncertainty.json`. |
| Audit boundary | Same as `SetMeaning`. |
| Blind retry allowed? | **No.** Same as `SetMeaning`. |
| Required pre-retry inspection | Check `VersionRecord.UncertaintyHash` and `uncertainty.json` existence. |

## 8. Audit contract

Current `AuditLog.Append` behavior:

- One JSON line is written through `os.File.Write`.
- The file is opened with `O_CREATE | O_WRONLY | O_APPEND`.
- No application-level lock is held.
- No `Sync` / `fsync` is called.
- `Close` is deferred; its error is not propagated.

Separate facts:

| Condition | Contract |
|---|---|
| **Write success** | `Write` returned no error. The record is likely in OS or file-system buffers. |
| **Crash durability** | **Not guaranteed.** A crash or process kill before OS flush may lose the record. |
| **Concurrent record integrity** | **Not guaranteed.** Concurrent `Append` calls may interleave or corrupt records on some platforms. |

A successful `AuditLog.Append` return does **not** prove the record is durably persisted or recoverable after a crash.

## 9. Success contract

`ADMIT` means only architectural admissibility. It is **not** execution success.

A future P0 production orchestration must not report successful end-to-end mutation solely because `admission.Result` is `ADMIT`. End-to-end success requires evidence that:

1. `Admission` returned `ADMIT`.
2. The Core handler executed to the extent that the operation semantics require.
3. Coherent Core state is confirmed according to the applicable mutation contract.
4. The required audit/evidence step completed at the adapter level (i.e., `AuditLog.Append` returned no error), with the crash-durability level explicitly chosen by the production policy.

This contract does not define transport response codes or redefine Core handler return semantics.

## 10. Partial failure contract

`PARTIAL_FAILURE` applies when repository evidence proves, or cannot safely exclude, persistent effects while the complete required operation was not confirmed.

A returned error does **not** imply that no state change occurred.

`PARTIAL_FAILURE` preserves earlier successful facts. For example:

- `ADMIT` + state persisted + audit failed
  - `Admission = ADMIT`
  - Core state = occurred / confirmed
  - Audit = incomplete
  - Overall processing = `PARTIAL_FAILURE`

A later failure must **not** convert an earlier `ADMIT` into `REJECT`.

## 11. State inspection before retry

Blind retry is forbidden whenever the previous attempt may have produced persistent effects.

| Mutation | Required pre-retry inspection |
|---|---|
| `CreateUnit` | `Unit` existence by key or deterministic ID; distinguish existence from absence as far as current information permits. |
| `CreateVersion` | Whether the deterministic `VersionRecord` already exists; whether `HeadVersionID` references it. |
| `SetMeaning` | Whether `VersionRecord.MeaningHash` matches the expected hash; whether `meaning.json` exists; whether they agree. |
| `SetClaims` | Whether `VersionRecord.ClaimSetHash` matches; whether `claimset.json` exists; whether they agree. |
| `SetUncertainty` | Whether `VersionRecord.UncertaintyHash` matches; whether `uncertainty.json` exists; whether they agree. |

This contract does not define automatic retry behavior, idempotency keys, or claim that inspection can always determine historical causation.

## 12. Recovery levels

| Level | Definition |
|---|---|
| **NO RECOVERY REQUIRED** | The coherent Core state and required audit are both confirmed. |
| **STATE INSPECTION REQUIRED** | The outcome is uncertain; inspect Core state to classify the failure. |
| **MANUAL REPAIR POSSIBLE** | A human or explicitly defined tool can repair a known partial state (e.g., reconstruct a missing audit record from confirmed state only when all deterministic inputs remain available). |
| **RECOVERY UNDEFINED** | The partial state cannot be repaired with current repository evidence. |

No recovery implementation is required by this contract.

## 13. Missing audit after confirmed state

If coherent Core state is confirmed but the corresponding `Audit.Append` failed:

- Do **not** roll back the historical fact in documentation.
- Do **not** claim the mutation never happened.
- Do **not** fabricate an `AuditEvent` silently.
- Do **not** blind-retry the mutation merely to obtain an audit record.

Classify as:

`PARTIAL_FAILURE / AUDIT INCOMPLETE`

A later recovery mechanism may reconstruct or repair audit evidence only if the required information is sufficient and the procedure is explicitly defined. This contract does not define that tool.

## 14. Orphaned sidecar

If a sidecar file exists on disk but the corresponding `VersionRecord` hash/reference is absent:

- A persistent artifact exists.
- The coherent Core mutation is **not** confirmed.
- The sidecar must not be silently adopted into Core state.
- The sidecar must not be silently deleted.
- Recovery requires inspection.

## 15. CreateVersion partial state

If the `VersionRecord` is appended to `UnitRecord` but `HeadVersionID` is not updated:

- This is a partial mutation, not successful `CreateVersion`.
- Blind retry is forbidden because the deterministic `Version.ID` may cause a duplicate `VersionRecord` insertion.
- Recovery must inspect the existing `VersionRecord` and `HeadVersionID` relationship before any retry or repair.

## 16. CreateUnit uncertain result

If `SaveUnit` outcome is uncertain and a later state check finds the `Unit` exists, the caller must not automatically infer whether the unit:

- pre-existed the request, or
- was persisted by the failed request,

unless available evidence proves causation.

This contract does not introduce execution-attempt identity.

## 17. Recovery must not rewrite identity

Any later recovery mechanism must not silently:

- change a `Unit.ID`,
- change a `Version.ID`,
- rewrite historical state identity,
- recanonicalize existing content,
- alter existing deterministic hashes merely to repair audit or evidence.

Recovery and identity are separate concerns.

## 18. Minimum production guarantee

The minimum truthful production guarantee supported by current evidence is:

> DigiEmu Core currently provides platform-dependent best-effort filesystem replacement and explicit deterministic state structures, but not cross-platform atomic mutation or crash durability.

Therefore a production Admission orchestration must:

1. Treat persistence errors conservatively.
2. Recognize possible partial effects.
3. Verify relevant Core state before retry.
4. Distinguish coherent state from incomplete sidecars, index, or audit.
5. Surface partial failure explicitly.
6. Never equate `ADMIT` with completed mutation.
7. Never rewrite earlier successful semantic stages because a later stage failed.

## 19. Test-only Phase E

This contract permits test-only Phase E orchestration if the tests explicitly cover:

- `ADMIT` + handler success (full success)
- `ADMIT` + handler error (execution error, no mutation or partial)
- `ADMIT` + persistence uncertainty (`Save*` error with possible partial side effects)
- `ADMIT` + audit failure (partial failure)
- `REJECT` → handler not invoked
- No blind-retry assumptions

No production wiring is permitted by this contract.

## 20. First production gate blockers

| Issue | Class |
|---|---|
| Actual production partial-failure detection | **BLOCKING** — the orchestration must classify `PARTIAL_FAILURE`. |
| State inspection / recovery procedure | **BLOCKING** — a deterministic state check before retry must be defined. |
| Audit completeness handling | **IMPORTANT** — missing audit after confirmed state must be surfaced. |
| Adapter error granularity | **IMPORTANT** — current `Save*` errors do not distinguish sub-step failures. |
| Audit durability (`Sync` / `fsync`) | **DEFERABLE** — not required for the first gate if partial-failure detection exists. |
| Concurrent audit safety | **DEFERABLE** only if the first production orchestration serializes audit appends; otherwise **IMPORTANT**. |
| Retry policy | **IMPORTANT** — at minimum, the "no blind retry" rule must be enforced. |

## 21. Current excellence filter

| Improvement | 2026 benefit | Class |
|---|---|---|
| `AuditLog` mutex | Medium — prevents concurrent corruption | **CONSIDER** |
| `AuditLog` `Sync` / `fsync` | Medium — improves durability | **CONSIDER** |
| `AuditLog` `Close` error handling | Low | **DEFER** |
| `AuditLog` write-length verification | Low | **CONSIDER** |
| More granular `Save*` error information | High — helps classify partial failures | **CONSIDER** |
| Sidecar / `UnitRecord` consistency verification | High — enables detection of orphan sidecars | **CONSIDER** |
| Duplicate `VersionRecord` prevention | Medium — would make `CreateVersion` retry safer | **CONSIDER** |
| Index rebuild | Medium — fixes best-effort index | **DEFER** |
| Recovery command / tool | Medium — enables manual repair | **DEFER** |
| Transaction-capable adapter | Low for 2026 | **REJECT** |

No implementation is required by this contract.

## 22. Explicit non-decisions

This contract does **not** decide:

- transaction architecture,
- WAL,
- recovery journal,
- database,
- new storage abstraction,
- saga,
- two-phase commit,
- new `AuditEvent.Type`,
- tombstones,
- `command_id` semantics,
- `event_id` semantics,
- `intent_id` semantics,
- idempotency-key semantics,
- execution-attempt identity,
- evidence archive format,
- historical verifier,
- HTTP/CLI failure mapping.

## 23. Contract principles

1. Never claim stronger persistence guarantees than the platform contract provides.
2. Admission success is not execution success.
3. Error does not prove absence of persistent side effects.
4. Coherent Core state and audit completeness are separate facts.
5. Unreferenced sidecars are not confirmed Core mutations.
6. Index state is not proof of `Unit` existence or absence.
7. Blind retry is forbidden after uncertain persistence.
8. Recovery begins with state inspection.
9. Recovery must not silently rewrite deterministic identity.
10. Stronger machinery is introduced only when demonstrated current need justifies it.

## 24. Non-claims

This document does **not** claim:

- cross-platform atomic filesystem writes,
- Windows atomic `os.Rename`,
- crash durability,
- concurrent `AuditLog` safety,
- state + audit atomicity,
- automatic rollback,
- automatic recovery,
- current idempotency,
- production Admission readiness,
- regulatory or legal evidentiary sufficiency.

## 25. Git verification statement

This document is a new proposed contract. No existing files have been modified.
