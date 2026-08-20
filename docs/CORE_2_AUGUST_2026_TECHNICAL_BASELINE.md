# DigiEmu Core 2.0 — August 2026 Technical Baseline

- **HEAD:** `8a75343 Add CreateUnit controlled local runtime evidence`
- **Branch:** `main`
- **Origin synchronization:** up to date
- **Working tree:** clean
- **CI:** Drift Guard PASS, CI PASS, Demo Verify Gate (Windows) PASS
- **Architecture baseline revision:** `0.3`
- **Global Admission:** `DISABLED`
- **`p0_architecture_constitution.enforced`:** `false`

## 1. Baseline

| Item | Value |
|---|---|
| HEAD | `8a75343` |
| Branch | `main` |
| CI | `GREEN` |
| Architecture baseline | `0.3` |
| Global Admission | `DISABLED` |
| `p0_architecture_constitution.enforced` | `false` |

## 2. August evidence chain

| Commit | What it established |
|---|---|
| `9d55d88` | Froze the Phase E test-only orchestration baseline. Demonstrated `ADMIT`/`REJECT` behavior against real Core handlers in a controlled test environment. |
| `f7fa892` | Added real-fs evidence for `CreateUnit` failure boundaries: `SaveUnit` fail-before-coherent, `Audit.Append` fail after coherent `Unit`, and `FindUnitByKey` decode failure. |
| `0461883` | Defined the formal `CreateUnit` Production Admission gate pilot design. Established the pilot scope, lock, intent mapping, and state-inspection contract. |
| `581fb82` | Implemented the first production `CreateUnit` Admission pilot in `cmd/digiemu/main.go`, with `--admission`, pilot lock, `admission.NewEngine`, and truthful failure reporting. |
| `f979647` | Froze the `CreateUnit` production pilot baseline, documenting proven claims, non-claims, controlled-run conditions, and conformance interpretation. |
| `18bdeaf` | Synchronized the P0 Conformance Matrix with the production pilot evidence, conservatively retaining `IR-01`, `IR-02`, and `IR-12` as `PARTIAL`. |
| `8a75343` | Added direct controlled local runtime evidence for the frozen pilot, including successful `ADMIT`, audit observation, `ADMIT`+domain failure, and lock-conflict cases. |

## 3. Current conformance

| IR | Status |
|---|---|
| IR-01 No Runtime Regression | `PARTIAL` |
| IR-02 Admission Before Mutation | `PARTIAL` |
| IR-03 Explicit Capability | `PASS` |
| IR-04 Explicit Ownership | `PASS` |
| IR-05 Command/Transition Correspondence | `PASS` |
| IR-06 Transition Evidence | `PASS` |
| IR-07 Verification Independence | `PASS` |
| IR-08 Single State Identity | `PASS` |
| IR-09 External Authority Preservation | `PASS` |
| IR-10 Fail Closed | `PASS` |
| IR-11 Deterministic Admission | `PASS` |
| IR-12 Traceability | `PARTIAL` |

No invariant has been promoted.

### Partial statuses explained

- **IR-01** — The legacy `CreateUnit` CLI branch is preserved and CI is green, but complete executable regression coverage of all legacy runtime behavior is not established.
- **IR-02** — The `digiemu unit create --admission` pilot path satisfies `Admission → mutation`. Legacy `CreateUnit`, the other four mutating capabilities, HTTP, and global `enforced: false` remain ungated.
- **IR-12** — The pilot runtime output provides `admission_id → transition_ref → Unit.ID`, but `admission_id` is not persisted into `Unit` state or `AuditEvent`, no `command_id`/`event_id` exists, and no durable Admission→Audit correlation is implemented.

## 4. Proven current capabilities

The repository now contains:

- P0 Admission Constitution (`ADR-P0-05`)
- `architecture-baseline.yaml` revision `0.3`
- core-capability registry
- aggregate-ownership registry
- command/event catalogue
- admission rule registry
- deterministic Admission Engine
- deterministic `admission_id` semantics
- Phase D bindings for all five mutating capabilities
- Phase E failure/partial-state evidence
- real-fs `CreateUnit` failure evidence
- the first real Production Admission CLI pilot
- controlled local runtime evidence

**Production pilot path:**

```text
digiemu unit create --admission
```

The pilot proves that a real production CLI mutation can be evaluated by the P0 Admission Engine before the existing `CreateUnit` Core use case is invoked.

## 5. Controlled runtime result

Reference: `docs/P0_CREATEUNIT_CONTROLLED_LOCAL_RUN_EVIDENCE_v0.1.md`

| Exercise | Result |
|---|---|
| Production success | `PASS` |
| Audit observation | `PASS` |
| `ADMIT` + domain failure | `PASS` |
| Lock conflict | `PASS` |
| Environment | `go1.25.7 windows/amd64` |

## 6. Intentional non-claims

August does **not** claim:

- global Production Admission
- HTTP Admission
- all-five production mutation gating
- general production readiness
- crash durability
- transactionality
- state+audit atomicity
- durable Admission→Audit correlation
- global multi-process safety
- idempotency
- automatic retry
- automatic rollback
- automatic recovery
- final `intent_id` semantics
- `command_id`/`event_id` semantics
- regulatory compliance

## 7. Do-not-touch baseline

Future work must **not** casually reopen:

- Core state identity
- canonicalization
- hashing
- deterministic replay
- verification
- existing deterministic Admission semantics
- frozen Phase D evidence
- frozen Phase E evidence
- frozen `CreateUnit` Production pilot baseline

Any semantic revision must be explicit, evidence-backed, and historically reconstructable.

## 8. Deferred items

The following are deliberately deferred:

- `architecture-baseline` `enforced: true`
- global Admission
- HTTP gating
- Production gating for `version.create`, `meaning.set`, `claim.set`, `uncertainty.set`
- persistent Admission→Audit correlation
- final `IntentID` semantics
- `command_id`/`event_id` derivation
- automatic recovery/retry
- persistence redesign
- transaction/WAL/database architecture
- Core 3.0

These are not described as defects unless the existing documentation does so.

## 9. Next technical entry point

**Next:** `CreateVersion` Production Admission Readiness Review.

**Not implementation yet.**

**Reason:** `CreateVersion` has a materially different persistence failure boundary (`SaveVersion → UpdateUnitHead`) and therefore provides a better second production case before introducing any shared production-gating abstraction.

A shared Production Orchestrator should **not** be created from the single `CreateUnit` pilot. A shared abstraction, if needed, must be derived only after multiple real production cases provide evidence.

## 10. August close

| Item | Status |
|---|---|
| DigiEmu Core 2.0 August hardening | `COMPLETE FOR PLANNED AUGUST SCOPE` |
| `CreateUnit` Production Admission pilot | `COMPLETE + FROZEN + RUNTIME EVIDENCED` |
| Core development | `PAUSED AT CLEAN BASELINE` |
| Next major implementation | `NOT STARTED` |
