# P0 CreateUnit Controlled Local Runtime Evidence v0.1

- **Status:** `CONTROLLED LOCAL RUNTIME EVIDENCE — PASS`
- **Run date:** 2026-08-19
- **Run baseline HEAD:** `18bdeaf18a7efccff464d55d9e1cdc911609ab21` (`18bdeaf`)
- **Frozen pilot implementation:** `581fb82`
- **Frozen pilot baseline:** `f979647`
- **Matrix synchronization:** `18bdeaf`

## 1. About this evidence

This document records direct runtime evidence from the compiled production `digiemu` CLI binary. It is **not** merely unit-test evidence. It does **not** change the frozen production pilot semantics.

## 2. Run baseline

| Item | Value |
|---|---|
| Full HEAD | `18bdeaf18a7efccff464d55d9e1cdc911609ab21` |
| Short HEAD | `18bdeaf` |
| Go version | `go version go1.25.7 windows/amd64` |
| OS/runtime environment | Windows (`%TEMP%`), compiled `digiemu.exe` for `windows/amd64` |
| Binary build command | `go build -o <pilot-root>\digiemu.exe ./cmd/digiemu` |
| `digiemu.exe` SHA-256 | `584173894657B70EEE797C97E7A8162B44426C2DBB820A5DB6DADABA4B8EA60F` |

The binary SHA-256 above is read directly from `binary-sha256.txt` in the controlled evidence directory.

## 3. Evidence isolation

The controlled run used an isolated temporary evidence root:

```text
%TEMP%\digiemu-p0-createunit-pilot-18bdeaf
```

Separate data directories were used for:

- successful mutation (`data-success`)
- invalid-key domain failure (`data-invalid-key`)
- lock-conflict case (`data-lock-conflict`)

The Git repository working tree remained clean after the runtime exercises. The absolute local path is **not** a normative path.

## 4. Successful production pilot run

Command semantics:

```text
digiemu unit create --admission \
  --key "p0-controlled-pilot-001" \
  --title "P0 Controlled Pilot 001" \
  --desc "Controlled local Production Admission pilot" \
  --data <isolated-success-data-dir>
```

Observed output (`run-success.txt`):

```text
admission_id=admission:sha256:0c93e480f89e975699aa9d5a96b71ba9ffb0c23a17dcef7a2cc30eb094619393
decision=ADMIT transition_ref=unit:created
OK: unit created id=unit_92ed464e27cd622f3f2a749054d6c9ae key=p0-controlled-pilot-001
exit_code=0
```

**Classification:** PASS

Directly observed:

- Admission result `ADMIT`
- `transition_ref = unit:created`
- real `Unit.ID` returned: `unit_92ed464e27cd622f3f2a749054d6c9ae`
- expected `Unit` key returned: `p0-controlled-pilot-001`
- process exit code `0`

This does **not** prove crash durability.

## 5. Audit observation

Command:

```text
digiemu audit tail --data <isolated-success-data-dir> --n 10
```

Observed output (`audit-tail-success.txt`):

```text
unit.created at=1787104877 actor=cli unit=unit_92ed464e27cd622f3f2a749054d6c9ae ver= id=evt_5809a8f489d10ad2d04e46d8abaf743e
audit_tail_exit_code=0
```

The `Unit.ID` observed in the successful `CreateUnit` result is also observed in the subsequent `unit.created` audit event. This proves runtime observation of the `CreateUnit` audit event for the same `Unit.ID`.

It does **not** prove durable Admission→Audit correlation because:

- `admission_id` is not stored in the `AuditEvent`.
- `AuditEvent.ID` is not linked to `admission_id`.
- `admission_id` is not persisted into `Unit` state.

No claim is made that `admission_id → AuditEvent.ID` is a durable relation.

## 6. ADMIT + domain failure

Command semantics:

```text
digiemu unit create --admission \
  --key "xx" \
  --title "P0 Controlled Invalid Key" \
  --data <isolated-invalid-key-data-dir>
```

Observed output (`run-invalid-key.txt`):

```text
admission_id=admission:sha256:92025be446d44b3cdc08dc4e3e259b6a5850210853220630dc5791d6b0ca9140
decision=ADMIT transition_ref=unit:created
execution/domain error: invalid unit key
persistence not attempted by this call
exit_code=1
```

**Classification:** `EXPECTED CONTROLLED FAILURE — PASS`

Directly demonstrates at runtime:

- Admission result: `ADMIT`
- Execution result: domain failure
- Persistence: `not attempted by this call`
- Process result: non-zero exit

Architectural admissibility and execution success remain separate facts. The later domain error does not rewrite `ADMIT` into `REJECT`.

## 7. Lock-conflict fail-closed case

Setup:

A pre-existing `.p0-admission-create-unit.lock` was placed in the isolated lock-conflict data directory.

Observed output (`run-lock-conflict.txt`):

```text
pilot lock already held: C:\Users\oondr\AppData\Local\Temp\digiemu-p0-createunit-pilot-18bdeaf\data-lock-conflict\.p0-admission-create-unit.lock
exit_code=1
```

No `admission_id=`, `decision=`, or `OK: unit created` was emitted.

**Classification:** `EXPECTED CONTROLLED FAILURE — PASS`

This is consistent with the frozen implementation's fail-before-Admission lock path. The ordering is:

```text
lock conflict
→ no Admission evaluation
→ no CreateUnit invocation
```

This is supported by the frozen implementation plus the executable pilot tests in `cmd/digiemu/unit_create_pilot_test.go`. The absence of CLI output alone is not treated as independent proof of internal non-execution.

## 8. Raw evidence hashes

Cross-checked from `runtime-evidence-sha256.txt`:

| File | SHA-256 |
|---|---|
| `run-success.txt` | `F7281568B10C5CAF1DD5B9FC7914FD6F9DFF8CB988B02CC57D69DC79BAEB04DF` |
| `audit-tail-success.txt` | `59AAB40A30A40CBF231DF924683751E3DD41E25D63463EB274424D152BE3513B` |
| `run-invalid-key.txt` | `C7B820265478FC7FC5D9B64728C11C719F769FDBFDB490A5C2CB76B549B60924` |
| `run-lock-conflict.txt` | `17FAA0CBB8AC09A384090457FAE96F1D9D758A0DC38DA3FE9F410AC70F81182D` |

The file hashes above match `runtime-evidence-sha256.txt`.

## 9. Evidence classification

Classified as **DIRECT CONTROLLED RUNTIME EVIDENCE**:

- A. Production binary executes Admission-gated `CreateUnit` successfully.
- B. `ADMIT` permits the existing production `CreateUnit` path.
- C. Successful `CreateUnit` returns a `Unit.ID`.
- D. Subsequent audit read observes a `unit.created` event for that same `Unit.ID`.
- E. `ADMIT` remains `ADMIT` when domain validation later fails.
- F. Known pre-persistence domain failure reports `persistence not attempted`.
- G. Lock-conflict run terminates non-zero before any Admission result is surfaced.
- H. The repository remains unchanged by the controlled run.

## 10. Combined evidence vs runtime-only evidence

**Runtime-only directly observes:**

- CLI outputs
- exit codes
- resulting `Unit.ID`
- subsequent `AuditEvent` for that `Unit.ID`
- lock-conflict failure output

**Frozen implementation/tests additionally establish:**

- lock is acquired before Admission
- `REJECT` prevents handler invocation
- Engine error prevents handler invocation
- lock conflict prevents Admission and mutation
- uncertain execution errors invoke state observation

These evidence classes are not collapsed.

## 11. Non-claims

This runtime run does **not** prove:

- crash durability
- transactionality
- state+audit atomicity
- global Admission
- HTTP Admission
- all-five mutation gating
- durable Admission→Audit correlation
- global multi-process safety
- idempotency
- automatic retry
- rollback
- recovery
- regulatory compliance
- general production readiness

## 12. Conformance impact

This document does **not** modify the P0 Conformance Matrix.

| IR | Evidence interpretation |
|---|---|
| IR-01 | Runtime evidence strengthens, but the status remains `PARTIAL`. |
| IR-02 | Production pilot runtime evidence strengthened, but the status remains `PARTIAL`. |
| IR-10 | Runtime evidence reinforces the existing `PASS`. |
| IR-11 | Runtime behavior is consistent with the existing `PASS`. |
| IR-12 | Runtime traceability strengthened, but the status remains `PARTIAL`. |

No status promotion is made by this document.

## 13. Final result

| Item | Value |
|---|---|
| Run baseline | `18bdeaf18a7efccff464d55d9e1cdc911609ab21` |
| Frozen implementation | `581fb82` |
| Frozen pilot | `f979647` |
| Environment | `go1.25.7 windows/amd64` |
| Production success | `PASS` |
| Audit observation | `PASS` |
| `ADMIT` + domain failure | `PASS` |
| Lock conflict | `PASS` |
| Repository remained clean | `YES` |
| Controlled local runtime evidence | `PASS` |
| General production readiness | `NOT CLAIMED` |
| Global Admission | `DISABLED` |
| Crash durability | `NOT CLAIMED` |
| Durable Admission→Audit correlation | `NOT PROVEN` |

## 14. Historical chain

| Commit/Artifact | Description |
|---|---|
| `9d55d88` | Phase E test orchestration freeze |
| `f7fa892` | `CreateUnit` real-fs failure evidence |
| `0461883` | Production pilot design |
| `581fb82` | Production pilot implementation |
| `f979647` | Production pilot freeze |
| `18bdeaf` | Conformance matrix synchronization |
| This document | 2026-08-19 controlled local runtime evidence |
