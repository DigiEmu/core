# P0 CreateUnit Production Admission Pilot Freeze v0.1

- **Status:** `FROZEN CREATEUNIT PRODUCTION PILOT BASELINE`
- **Frozen at HEAD:** `581fb82 Add CreateUnit production admission pilot`
- **Date:** 2026-08-19
- **Architecture baseline:** `architecture-baseline.yaml` revision `0.3`
- **ADR:** `01_decisions/ADR-P0-05-Core-Admission-Constitution.md`

## 1. About this freeze

This document freezes the historical meaning and evidence baseline of the first `CreateUnit` Production Admission pilot.

It does **not** make source files permanently immutable. Future changes may supersede the pilot through an explicit new revision, while this baseline remains historically reconstructable.

## 2. Baseline

| Artifact | Commit |
|---|---|
| Production pilot implementation | `581fb82` |
| Pilot design | `0461883` |
| Real-fs evidence | `f7fa892` |
| Phase E test orchestration freeze | `9d55d88` |

| Evidence | Status |
|---|---|
| CI at `581fb82` (Drift Guard, CI, Demo Verify Gate Windows) | `PASS` |
| Working tree at review | `clean` |
| `p0_architecture_constitution.enforced` | `false` |

## 3. Frozen pilot path

The demonstrated production path is:

```text
digiemu unit create --admission
→ acquire pilot lock
→ construct P0 Intent
→ Admission Engine Evaluate
→ system error / REJECT stop
→ ADMIT
→ existing CreateUnit use case
→ success or truthful failure/state inspection
```

| Concept | Value |
|---|---|
| Pilot mutation | `core.unit.create` |
| Aggregate | `unit` |
| Command | `unit.create` |
| Transition | `unit:created` |

## 4. Activation boundary

- `--admission` is an **explicit opt-in**.
- Default: `false`.
- Without `--admission`, the legacy `CreateUnit` path remains available and unchanged.
- The other four mutating capabilities are **not** production-gated by this pilot.
- HTTP is **not** Admission-gated by this pilot.
- Global Admission is **not** enabled.

## 5. Directly proven claims

The following claims are directly evidenced by `581fb82` and its supporting baselines:

1. A real production CLI mutation path can be Admission-gated.
2. Admission occurs before `CreateUnit` mutation on the gated path.
3. `REJECT` prevents `CreateUnit` invocation.
4. Admission Engine evaluation error prevents `CreateUnit` invocation.
5. `ADMIT` permits existing `CreateUnit` execution.
6. A later execution failure does not rewrite `ADMIT` into `REJECT`.
7. Known pre-persistence `CreateUnit` failures do not trigger unnecessary state observation.
8. Execution failures with uncertain persistence state trigger conservative state inspection.
9. State inspection reports only current observable facts.
10. A coherent `Unit` can remain observable after an `Audit.Append` failure under the deterministic real-fs evidence case from `f7fa892`.
11. Pilot lock conflict fails before Admission and mutation.
12. Normal pilot lock cleanup closes and removes the lock.
13. Existing Core state identity, hashing, canonicalization, replay, and verification semantics were not changed by the pilot.

## 6. Real-fs evidence

The frozen real-fs evidence lives in:

`internal/kernel/adapters/fs/create_unit_pilot_test.go` at commit `f7fa892`.

The three demonstrated cases are:

- **A.** `SaveUnit` failure before a coherent `Unit` is observable.
- **B.** `Audit.Append` failure after a coherent `Unit` becomes observable.
- **C.** `FindUnitByKey` observation/decode failure.

This evidence does **not** include, and the repository does **not** contain, deterministic portable evidence for the case where `SaveUnit` returns an error **and** a coherent final `Unit` is simultaneously observable. Do not infer this case from other evidence.

## 7. Failure semantics

The pilot failure semantics are frozen as follows.

For a known pre-persistence `CreateUnit` error, the CLI reports:

```text
persistence not attempted by this call
```

For an uncertain execution failure, the pilot performs a `FindUnitByKey` observation. Possible observational conclusions are:

- `no coherent Unit currently observable`
- `coherent Unit currently observable`
- `Unit observable; causation unresolved`
- `repository observation failed; mutation state unresolved`

These messages are **observations**. They do **not** establish:

- crash durability
- transactionality
- historical causation
- filesystem artifact absence

## 8. Admission stability

`ADMIT` is an architectural admissibility result. Once Admission returns `ADMIT`, later domain, persistence, or audit errors do **not** turn the Admission result into `REJECT`. Execution outcome and Admission outcome remain distinct facts.

## 9. Traceability baseline

Current pilot runtime traceability is:

```text
admission_id
→ capability_ref
→ command_ref
→ transition_ref
→ resulting Unit.ID
```

Limitations:

- `admission_id` is **not** persisted into `Unit` state.
- `admission_id` is **not** persisted into the `AuditEvent`.
- No durable Admission→Audit correlation exists.
- No `command_id` is produced.
- No `event_id` is produced.
- `AuditEvent.ID` is not exposed by `CreateUnitResponse`.
- `IntentID` remains a non-normative placeholder.

Therefore, **IR-12 remains `PARTIAL`**. This freeze record does not introduce new normative ID semantics.

## 10. Intent ID

Current non-normative placeholder:

```text
p0-pilot-unresolved-intent-id
```

This placeholder is:

- **not** a unique execution identity
- **not** a retry identity
- **not** an idempotency identity
- **not** a state identity
- **not** an evidence identity
- **not** final `IntentID` semantics

Existing Admission digest semantics already exclude `IntentID`, as previously established.

## 11. Pilot lock

Frozen mechanism:

```text
<data-dir>/.p0-admission-create-unit.lock
```

using `O_CREATE | O_EXCL | O_WRONLY`.

Limitations:

- Protects only Admission-enabled `CreateUnit` pilot instances using the same mechanism.
- Does not constrain legacy writers.
- Does not provide global Core multi-process safety.
- A stale lock may remain after abnormal termination.
- No stale-lock recovery is implemented.

Dedicated pilot data-directory isolation remains an **operational requirement**.

## 12. Controlled run status

| | |
|---|---|
| **READY FOR CONTROLLED LOCAL PILOT RUN** | **YES** |
| **READY FOR GENERAL PRODUCTION USE** | **NO** |
| **READY FOR GLOBAL ADMISSION ENFORCEMENT** | **NO** |

Conditions for a controlled local pilot run:

- Use a dedicated `--data` directory.
- Ensure no concurrent legacy or other writers against that directory.
- Explicitly pass `--admission`.
- Be aware of potential stale locks.
- No automatic retry.
- No automatic recovery.
- No crash-durability claim.
- Preserve relevant CLI output/evidence if later analysis is required.

## 13. Non-claims

This freeze explicitly does **not** prove:

- crash durability
- transactionality
- state+audit atomicity
- global Admission enforcement
- HTTP Admission
- all-five mutation production gating
- durable end-to-end Admission→Audit correlation
- global multi-process Core safety
- idempotency
- automatic retry
- automatic rollback
- automatic recovery
- regulatory compliance
- general production readiness

## 14. Conformance status

This is the conservative freeze-record interpretation of the P0 Conformance Matrix. It does **not** modify the matrix.

| IR | Status | Reason |
|---|---|---|
| IR-01 | `PARTIAL` | `581fb82` substantially strengthens regression evidence, but direct executable coverage of the legacy `CreateUnit` CLI branch is not complete enough to claim a fully proven no-runtime-regression invariant. |
| IR-02 | `PARTIAL` | The explicit `CreateUnit` pilot path satisfies Admission-before-mutation, but legacy `CreateUnit`, other mutating capabilities, HTTP, and global enforcement remain ungated. |
| IR-03 | `PASS` | No change. |
| IR-04 | `PASS` | No change. |
| IR-05 | `PASS` | No change. |
| IR-06 | `PASS` | No change. |
| IR-07 | `PASS` | No change. |
| IR-08 | `PASS` | No change. |
| IR-09 | `PASS` | No change. |
| IR-10 | `PASS` | No change. |
| IR-11 | `PASS` | No change. |
| IR-12 | `PARTIAL` | Runtime CLI traceability exists, but durable Admission→mutation→Audit correlation does not. |

## 15. Historical reconstructability

The minimum baseline chain is:

1. `9d55d88` — Phase E orchestration freeze
2. `f7fa892` — Real-fs failure evidence
3. `0461883` — Production pilot design
4. `581fb82` — Production pilot implementation

Plus the relevant normative sources:

- `ADR-P0-05`
- `architecture-baseline.yaml` revision `0.3`
- P0 registries
- Admission schemas
- `docs/P0_PRODUCTION_FAILURE_SEMANTICS_v0.1.md`
- `docs/P0_ATOMICITY_RECOVERY_CONTRACT_v0.1.md`

This chain defines the historical meaning of the first Production Admission pilot.

## 16. Future change rule

Future work may:

- revise the pilot,
- add another production mutation,
- add stronger traceability,
- or change enforcement policy,

only through a new explicit evidence/revision.

Do not silently reinterpret this frozen baseline, and do not mutate historical evidence merely to align it with future semantics.

## 17. Final summary

| Item | Value |
|---|---|
| Implementation HEAD | `581fb82` |
| Design | `0461883` |
| Real-FS evidence | `f7fa892` |
| Phase E | `9d55d88` |
| CI | `GREEN` |
| Pilot mutation | `core.unit.create` |
| Production pilot | `FROZEN` |
| Controlled local pilot | `READY` |
| General production | `NOT READY` |
| Global Admission | `DISABLED` |
| `architecture-baseline` enforced | `false` |
| Retry | `NONE` |
| Recovery | `NONE` |
| Crash durability | `NOT CLAIMED` |
| IR-01 | `PARTIAL` |
| IR-02 | `PARTIAL` |
| IR-12 | `PARTIAL` |
