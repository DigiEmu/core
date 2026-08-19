# P0 CreateUnit Production Gate Pilot Design

- **Document:** P0 CreateUnit Production Gate Pilot Design
- **Version:** v0.1
- **Status:** PROPOSED — IMPLEMENTATION-BOUND PILOT DESIGN
- **Design commit:** `f7fa892 Add CreateUnit real fs pilot failure evidence`
- **Architecture baseline:** `0.3`
- **Date:** 2026-08-19
- **Scope:** First explicit `core.unit.create` Production Admission pilot for `digiemu unit create --admission`. This document does **not** claim production readiness.

## 1. Baseline

| Item | Value |
|---|---|
| **HEAD** | `f7fa892 Add CreateUnit real fs pilot failure evidence` |
| **Phase D** | FROZEN |
| **Phase E** | FROZEN |
| **CreateUnit real-fs evidence** | COMMITTED |
| **CI** | PASS |
| **Production Admission globally** | `p0_architecture_constitution.enforced: false` |

## 2. Pilot scope

**Pilot mutation:** `core.unit.create`

**Pilot entry point:** `digiemu unit create --admission`

The pilot covers exactly **one** existing CLI mutation path.

**Out of scope:**

- HTTP
- All other mutations (`CreateVersion`, `SetMeaning`, `SetClaims`, `SetUncertainty`)
- Global Admission enforcement
- Production Executor/Orchestrator
- Retry engine
- Recovery engine
- Rollback system
- Idempotency
- Transactionality
- Crash durability
- State+audit atomicity
- Global concurrency architecture
- Final `intent_id` semantics
- `command_id`/`event_id` infrastructure
- Core 3.0

## 3. Activation

| | |
|---|---|
| **Flag** | `--admission` |
| **Type** | command-local `bool` flag |
| **Default** | `false` |
| **When absent** | legacy `CreateUnit` path, unchanged behavior |
| **When present** | evaluate P0 Admission before `CreateUnit` |
| **Help text** | `Evaluate P0 Admission before creating the unit (experimental pilot).` |
| **Feature label** | `P0 ADMISSION PILOT — EXPERIMENTAL` |

The activation does **not** imply universal Production Admission enforcement.

## 4. Data isolation

The existing CLI supports `--data <dir>`.

**Pilot operational rule:**

- The pilot must use a dedicated pilot data directory.
- Only the Admission-enabled `CreateUnit` pilot may mutate that directory during a pilot run.
- Ungated/legacy mutation commands must not concurrently use the same pilot data directory.
- This is a **PILOT CONSTRAINT**, not a global Core concurrency guarantee.

## 5. Single-writer mechanism

| | |
|---|---|
| **Lock file** | `<data-dir>/.p0-admission-create-unit.lock` |
| **Mechanism** | `os.OpenFile(lockPath, os.O_CREATE\|os.O_EXCL, 0o644)` |
| **On success** | Remove lock in `defer` after `CreateUnit` completes. |
| **On existing lock** | Fail closed before Admission evaluation and mutation. |
| **Scope** | Only `digiemu unit create --admission`. |

**Explicit limitations:**

- A stale lock is possible after a process crash and requires manual removal.
- This is not a universal Core write lock.
- It does not protect legacy processes that ignore it.
- Dedicated pilot data isolation remains mandatory.

## 6. Intent mapping

```go
intent := admission.Intent{
    SchemaVersion:        "v0.1",
    ArchitectureRevision: "0.3",
    IntentID:             "p0-pilot-unresolved-intent-id",
    CapabilityRef:        "core.unit.create",
    AggregateRef:         "unit",
    CommandRef:           "unit.create",
    Payload: map[string]any{
        "key":         key,
        "title":       title,
        "description": description,
    },
}
```

`ActorID` is **not** part of the Admission `Payload` semantics. It is supplied to `CreateUnit` directly from the CLI.

`IntentID` is explicitly **NON-NORMATIVE**:

- not unique identity
- not state identity
- not evidence identity
- not retry identity
- not idempotency identity
- does not define future `IntentID` semantics

## 7. Gate control flow

```text
CLI input
→ dedicated pilot data directory
→ acquire pilot lock
→ normalize input
→ construct admission.Intent
→ engine := admission.NewEngine(admission.V01Registry())
→ result, err := engine.Evaluate(intent)

Engine error:
    → stop, handler not invoked

REJECT:
    → stop, handler not invoked
    → expose admission_id + reason_codes

ADMIT:
    → call existing CreateUnit use case
    → on success:
        → expose admission_id + Unit.ID/key
    → on error:
        → classify whether persistence may have begun
        → state inspection where required
        → truthful CLI outcome
        → NO automatic retry
        → NO rollback
        → NO automatic recovery

release pilot lock
```

No production `Executor`/`Orchestrator` is required.

## 8. CreateUnit error phases

Actual `CreateUnit` sequence:

```text
domain.NewUnit
→ Repo.ExistsByKey
→ Repo.SaveUnit
→ AuditEvent construction
→ Audit.Append
```

| Phase | Error | Handler invoked? | Persistence begun? |
|---|---|---|---|
| A | `domain.NewUnit` validation | Yes | No |
| B | `ExistsByKey` pre-persistence read or `ErrUnitAlreadyExists` | Yes | No |
| C | `SaveUnit` persistence | Yes | Yes, uncertain |
| D | `AuditEvent` construction | Yes | Yes |
| E | `Audit.Append` | Yes | Yes |

Handler invocation does **not** imply persistence.

## 9. State inspection rule

No post-error `FindUnitByKey` for errors provably before persistence:

- `domain.NewUnit` validation
- `ErrUnitAlreadyExists`
- `ExistsByKey` read failure for mutation certainty of this call
- `ErrAuditNotConfigured`
- `ErrClockNotConfigured`

**Required** after:

- `SaveUnit` error
- `Audit.Append` error

Use `FindUnitByKey(expectedKey)` and interpret conservatively:

| Observation | Inference |
|---|---|
| `!ok` and `err == nil` | No coherent Unit currently observable. |
| `ok` and fields match expected input | Coherent Unit currently observable. |
| `ok` and fields do not match | Unit observable; causation unresolved. |
| `err != nil` | Repository observation failed; mutation state unresolved. |

Never infer crash durability. Never equate filesystem artifacts with coherent Core `Unit` state.

## 10. Real-fs evidence

Referenced commit: `f7fa892`

Test file: `internal/kernel/adapters/fs/create_unit_pilot_test.go`

**Directly proven real-fs observations:**

A. **SaveUnit failure before coherent Unit**
- `.json.tmp` path pre-created as a directory.
- `ioutil.WriteFile` fails.
- `CreateUnit` returns error.
- `FindUnitByKey` returns `ok=false`.

B. **Audit.Append failure after SaveUnit success**
- `audit.ndjson` path pre-created as a directory.
- `os.OpenFile` fails.
- `CreateUnit` returns error.
- `FindUnitByKey` returns the matching `Unit`.

C. **Observation/decode failure**
- Corrupt `Unit` `.json` file.
- `FindUnitByKey` returns a read/decode error.
- Mutation certainty cannot be established through `FindUnitByKey`.

**Explicitly not proven:**

`SaveUnit` returns an error AND a coherent final `.json` Unit is observable. This case remains architecturally possible but is not deterministically injectable with the existing adapter.

## 11. Filesystem semantic limits

- `tmp + os.Rename` is platform-dependent.
- Cross-platform `Rename` atomicity is **not** claimed.
- Windows atomic replacement is **not** claimed.
- Crash durability is **not** claimed.
- Transactionality is **not** claimed.

Current tests prove observable process-level behavior only.

## 12. Audit failure

For `Audit.Append` failure after `SaveUnit` returned `nil`:

- `Admission` remains `ADMIT`.
- `CreateUnit` returns the audit error.
- Inspect `Unit`.
- If the `Unit` is currently observable, report `UNIT_OBSERVED_AUDIT_INCOMPLETE`.
- Do not delete the `Unit`.
- Do not retry the mutation.
- Do not rewrite the `Admission` decision.

Do not use phrasing such as `STATE_PERSISTED_AUDIT_FAILED`.

## 13. CLI outcome model

Conceptual categories only. No exported runtime enum.

| Category | When |
|---|---|
| `ADMISSION_REJECT` | `result.Decision == "REJECT"` |
| `ADMISSION_SYSTEM_ERROR` | `Engine.Evaluate` returns `error` |
| `EXECUTION_DOMAIN_ERROR` | `NewUnit` validation or `ErrUnitAlreadyExists` |
| `EXECUTION_REPOSITORY_READ_ERROR` | `ExistsByKey` read error |
| `PERSISTENCE_STATE_UNRESOLVED` | `SaveUnit`/`Audit` error + `FindUnitByKey` cannot determine state |
| `UNIT_OBSERVED_AUDIT_INCOMPLETE` | `SaveUnit` succeeded, `Audit.Append` failed, `Unit` observable |
| `SUCCESS` | `CreateUnit` succeeded |

A single non-zero exit code is sufficient for all failure categories unless later implementation evidence proves otherwise.

## 14. AdmissionID

`AdmissionID` is already deterministically computed by the `Engine`.

- Surface on `REJECT`.
- Surface on `ADMIT` success.
- Surface on `ADMIT +` handler failure.
- No additional persistent storage required for the pilot.

If `Engine.Evaluate` returns a system error before any `Result` is built, `AdmissionID` is unavailable.

## 15. Minimum traceability

```text
admission_id
→ "core.unit.create"
→ "unit.create"
→ "unit:created"
→ resulting Unit.ID
```

`AuditEvent.ID` is **not** in the minimum traceability chain because `CreateUnitResponse` does not expose it.

- `command_id`: not required.
- `event_id`: not required.
- `IntentID`: non-normative placeholder.

IR-12 remains `PARTIAL`.

## 16. Retry / rollback / recovery

| Mechanism | Decision |
|---|---|
| Automatic retry | **NO** |
| Rollback | **NO** |
| Automatic recovery | **NO** |

After uncertainty:

- state inspection
- truthful output
- operator/manual decision

No retry/idempotency/recovery infrastructure is added.

## 17. Legacy coexistence

- `digiemu unit create` without `--admission` retains existing behavior.
- `digiemu unit create --admission` activates the pilot gate.
- Other mutation commands remain unchanged.
- `architecture-baseline.yaml` `enforced: false` remains truthful.

## 18. HTTP

HTTP and `cmd/api` are **out of scope** for this pilot.

`cmd/api` has separate `Audit`/`Clock` wiring concerns that must not be addressed here.

## 19. Security / authority

`Admission` does **not** replace:

- domain validation
- authorization
- external authority
- truth verification
- regulatory approval

Admission remains deterministic architectural admissibility only.

## 20. Expected implementation footprint

| Item | Estimate |
|---|---|
| **Production modification** | `cmd/digiemu/main.go` |
| **Tests** | Narrow `cmd/digiemu` pilot tests |
| **Existing real-fs evidence** | `internal/kernel/adapters/fs/create_unit_pilot_test.go` |
| **Schema changes** | none |
| **Registry changes** | none |
| **ADR changes** | none |
| **State identity / hashing / replay / verification changes** | none |
| **Production orchestrator** | none |

## 21. Implementation preconditions

Before implementation:

- Phase D frozen
- Phase E frozen
- real-fs evidence committed
- current HEAD CI green
- pilot design documented
- dedicated pilot data-dir rule established
- single-writer lock behavior specified

## 22. Run preconditions

Before an actual pilot run:

- implementation tests passing
- real-fs evidence still passing
- CI green
- dedicated pilot data directory
- no concurrent legacy writers against that directory
- pilot lock verified
- automatic retry disabled
- crash durability explicitly not claimed

## 23. Non-claims

The design does **not** prove:

- Production readiness
- Global Admission
- HTTP Admission
- All-five production mutation gating
- Crash durability
- Transactional persistence
- State+audit atomicity
- Multi-process Core safety
- Automatic recovery
- Idempotency
- Complete IR-02
- Complete IR-12
- Regulatory compliance
- Core 3.0

## 24. Design freeze rule

This document is the implementation-bound design baseline.

Future implementation must either:

A. conform to this design, or
B. explicitly revise the design before changing semantics.

The design does not require source files to be permanently immutable, but the historical meaning must remain reconstructable.

## 25. Final summary

| Item | Value |
|---|---|
| **Baseline HEAD** | `f7fa892` |
| **Phase D** | FROZEN |
| **Phase E** | FROZEN |
| **Real-FS evidence** | COMMITTED + CI GREEN |
| **Pilot mutation** | `core.unit.create` |
| **Pilot command** | `digiemu unit create --admission` |
| **Global Admission** | DISABLED |
| **Production orchestrator** | NOT REQUIRED |
| **Automatic retry** | DISABLED |
| **Automatic recovery** | NONE |
| **Crash durability** | NOT CLAIMED |
| **Minimum traceability** | `admission_id` → capability → command → transition → `Unit.ID` |
| **Ready for implementation** | READY AFTER DESIGN REVIEW |
| **Ready for run** | NO |
