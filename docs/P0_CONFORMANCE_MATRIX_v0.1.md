# P0 Conformance Matrix v0.1

- **Architecture baseline:** `architecture-baseline.yaml` revision 0.3
- **ADR:** `01_decisions/ADR-P0-05-Core-Admission-Constitution.md`
- **Manifest:** `docs/CORE_ARCHITECTURE_INTEGRATION_MANIFEST_v0.1.md`
- **Date:** 2026-08-07
- **Status:** Proposed

## Methodology

Each integration invariant (IR-01 through IR-12) is evaluated against the actual repository evidence. A status is assigned only when the corresponding evidence exists; documentation alone is insufficient for `PASS`.

| Status | Meaning |
|--------|---------|
| `PASS` | Sufficient specification, implementation, and/or executable evidence exists. |
| `PARTIAL` | Evidence exists but is incomplete, scoped to RS-001 only, or has known limitations. |
| `OPEN` | No implementation or executable evidence; only documentation/specification. |
| `BLOCKED` | Progress is blocked by an external or environmental issue. |

## Conformance Matrix

### IR-01 — No Runtime Regression

- **ID:** IR-01
- **Title:** No Runtime Regression
- **Normative statement:** Architecture integration MUST NOT silently alter existing Core 2.0 behavior or stable contracts.
- **Current status:** PARTIAL
- **Evidence/artifacts:**
  - `architecture-baseline.yaml` lists the Core 2.0 boundary as authoritative and unchanged.
  - ADR-P0-05 and the manifest explicitly state that canonicalization, hashing, replay, verification, CLI v1, and state identity remain unchanged.
  - No files in `internal/`, `cmd/`, `pkg/`, or `schemas/` (existing schemas) were modified by P0 work.
  - `go test ./internal/kernel/... ./internal/conformance ./internal/verify -timeout 120s` passes.
- **Executable evidence/tests:** Targeted `go test` of kernel, conformance, and verify packages passes.
- **Remaining gap:** Full `go test ./...` is currently not green because `go run ./cmd/digiemu` in several test packages fails with `digiemu.exe` access denied in the local Windows temp environment. This is a file-lock/build-cache issue, not a code regression, so it should not be falsely reported as fully green.
- **Next action:** Resolve the temporary `go-build*` directory / antivirus / build-cache file lock, then rerun full `go test ./...` and confirm the existing Core 2.0 reproduction tests are byte-exact.

### IR-02 — Admission Before Mutation

- **ID:** IR-02
- **Title:** Admission Before Mutation
- **Normative statement:** A mutation governed by the Architecture Constitution MUST pass Admission before execution.
- **Current status:** OPEN
- **Evidence/artifacts:**
  - The conceptual target flow in the manifest shows `Intent → Admission → Command → Existing Core Kernel Use Case`.
  - `architecture-baseline.yaml` explicitly sets `p0_architecture_constitution.enforced: false`.
- **Executable evidence/tests:** None. RS-001 stops at the generated Command Envelope and does not invoke an existing Core use case.
- **Remaining gap:** No production or test path intercepts a Core mutation to require Admission first. The `digiemu` CLI and `internal/kernel/usecases/` still execute directly without an Admission gate.
- **Next action:** Define the first Admission-bound Core use case and a non-production wiring point (e.g., a conformance runner or a prototype middleware) that demonstrates the `Intent → Admission → Command → Use Case` sequence for `core.unit.create`.

### IR-03 — Explicit Capability

- **ID:** IR-03
- **Title:** Explicit Capability
- **Normative statement:** Every constitution-governed command MUST reference a registered capability.
- **Current status:** PARTIAL
- **Evidence/artifacts:**
  - `core-capability-registry.yaml` lists Core capabilities with `id`, `mutation`, and `source` fields.
  - `schemas/intent-envelope.schema.json` requires `capability_ref`.
  - `testdata/rs_001/run.ps1` reads the registry and evaluates `capability_ref` for each RS-001 case.
- **Executable evidence/tests:**
  - RS-001 `valid_admit` passes with `capability_ref: core.unit.create`.
  - RS-001 `invalid_unknown_capability` rejects an unregistered capability.
- **Remaining gap:** Only the `core.unit.create` and nearby capabilities are tested. The registry is not yet complete for all constitution-governed mutation paths, and the `digiemu` CLI does not require the registry.
- **Next action:** Complete the capability registry for all `mutation: true` Core capabilities and add corresponding RS-00x or negative conformance cases.

### IR-04 — Explicit Ownership

- **ID:** IR-04
- **Title:** Explicit Ownership
- **Normative statement:** Every constitution-governed mutation MUST resolve to an authorized aggregate owner.
- **Current status:** PARTIAL
- **Evidence/artifacts:**
  - `aggregate-ownership-registry.yaml` maps `unit` to its owned capabilities.
  - `schemas/intent-envelope.schema.json` requires `aggregate_ref`.
  - `testdata/rs_001/run.ps1` checks that the capability's owner matches the intent's `aggregate_ref`.
- **Executable evidence/tests:**
  - RS-001 `valid_admit` passes with `aggregate_ref: unit`.
  - RS-001 `invalid_ownership_mismatch` rejects a non-owning aggregate.
- **Remaining gap:** The registry only defines the `unit` aggregate. Other aggregates (if any) are not yet mapped. Ownership is not enforced in the runtime kernel.
- **Next action:** Expand `aggregate-ownership-registry.yaml` as new mutation paths are added, and ensure every `mutation: true` capability has a registered owner.

### IR-05 — Command/Transition Correspondence

- **ID:** IR-05
- **Title:** Command/Transition Correspondence
- **Normative statement:** Every admitted mutating command MUST resolve to a defined transition.
- **Current status:** PARTIAL
- **Evidence/artifacts:**
  - `command-event-catalogue.yaml` maps `unit.create` to `capability_id: core.unit.create`, `aggregate_id: unit`, `transition_id: unit:created`, and `event_id: unit.created`.
  - `schemas/command-envelope.schema.json` requires `transition_ref`.
  - `testdata/rs_001/run.ps1` resolves `command_ref` to `transition_id` and validates it.
- **Executable evidence/tests:**
  - RS-001 `valid_admit` produces a Command Envelope with `transition_ref: unit:created`.
  - RS-001 `invalid_unknown_command` and `invalid_command_capability_mismatch` reject mismatched commands.
- **Remaining gap:** Only one transition (`unit:created`) is catalogued. A formal transition registry or additional catalogue rows are needed for other mutation paths.
- **Next action:** Add the remaining Core mutation commands and transitions to `command-event-catalogue.yaml` and cover them with RS-00x cases.

### IR-06 — Transition Evidence

- **ID:** IR-06
- **Title:** Transition Evidence
- **Normative statement:** Every successfully executed constitution-governed transition MUST produce or reference deterministic evidence sufficient to identify the resulting transition.
- **Current status:** PARTIAL
- **Evidence/artifacts:**
  - `command-event-catalogue.yaml` references `event_id: unit.created`, which matches the runtime `AuditEvent.Type` emitted by `internal/kernel/usecases/create_unit.go`.
  - `schemas/event-envelope.schema.json` defines a wrapper for referencing event evidence without replacing the existing verify schemas.
  - ADR-P0-05 and `testdata/rs_001/README.md` explicitly state that RS-001 does **not** execute the `CreateUnit` handler.
- **Executable evidence/tests:**
  - RS-001 does not execute the runtime handler, so it cannot produce actual event evidence.
  - The catalogue link to `unit.created` is verified statically against `internal/kernel/domain/audit_event.go`.
- **Remaining gap:** No executable path demonstrates an admitted command being handed to a Core use case, producing an audit event, and wrapping that event in an Event Envelope. The event evidence is referenced, not generated.
- **Next action:** Create a second reference scenario (or extend a conformance runner) that invokes an existing Core use case after Admission and emits the matching audit event, producing a validated Event Envelope.

### IR-07 — Verification Independence

- **ID:** IR-07
- **Title:** Verification Independence
- **Normative statement:** Admission MUST NOT redefine canonicalization, hashing, reconstruction, or verification.
- **Current status:** PASS
- **Evidence/artifacts:**
  - `architecture-baseline.yaml` places canonicalization, snapshot hashing, replay, and verification under the authoritative Core 2.0 boundary, not the P0 constitution.
  - ADR-P0-05 `KEEP` and `Non-Goals` explicitly preserve `digiemu-canonical-json-v1`, SHA-256, verify/replay, and CLI v1.
  - P0 artifacts (`intent-envelope.schema.json`, `command-envelope.schema.json`, `event-envelope.schema.json`) carry references and metadata; they do not define new canonicalization profiles or hash algorithms.
  - `internal/verify` package targeted tests pass.
- **Executable evidence/tests:**
  - `go test ./internal/verify -timeout 120s` passes.
  - No existing `VERIFY_RESULT_SCHEMA_v1.json`, `verify_result_v2.schema.json`, or `core_2_conformance_report.schema.json` files were modified.
- **Remaining gap:** None identified.
- **Next action:** Maintain the separation: any future Admission artifacts must continue to reference, not redefine, Core verification semantics.

### IR-08 — Single State Identity

- **ID:** IR-08
- **Title:** Single State Identity
- **Normative statement:** The architecture layer MUST NOT introduce an alternative DigiEmu state-identity mechanism.
- **Current status:** PASS
- **Evidence/artifacts:**
  - P0 envelope schemas only introduce `intent_id`, `command_id`, and `event_id` as architectural envelope identifiers; these are explicitly documented as not yet normatively derived and not alternative state identity.
  - The envelope fields `capability_ref`, `aggregate_ref`, `command_ref`, and `transition_ref` are references into existing registries, not new identity producers.
  - `architecture-baseline.yaml` and the manifest state that state identity is produced only by the existing Core canonicalization profile.
- **Executable evidence/tests:**
  - `schemas/intent-envelope.schema.json`, `schemas/command-envelope.schema.json`, and `schemas/event-envelope.schema.json` contain no canonicalization or hashing fields.
- **Remaining gap:** `admission_id`, `intent_id`, `command_id`, and `event_id` derivation rules are still undefined; this must not drift into an alternative identity scheme.
- **Next action:** When normative identifier derivation is specified, explicitly document that it is an envelope-instance identifier, not a state-identity mechanism.

### IR-09 — External Authority Preservation

- **ID:** IR-09
- **Title:** External Authority Preservation
- **Normative statement:** Admission MUST NOT absorb business, legal, regulatory, ethical, or trust authority belonging outside DigiEmu Core.
- **Current status:** PASS
- **Evidence/artifacts:**
  - ADR-P0-05 §9 states `Admission != Business Approval`, `Admission != Regulatory Approval`, `Admission != Trust Assignment`.
  - `docs/CORE_2_INTEROP_CONTRACT.md`, `docs/CORE_2_BOUNDARY_MODEL.md`, and `docs/CORE_2_EXTERNAL_REVIEW_NOTE.md` preserve external authority boundaries.
  - `testdata/rs_001/README.md` and `architecture-baseline.yaml` external_authority section re-state that business/legal/trust authority remains outside Core.
- **Executable evidence/tests:**
  - No P0 code path claims to issue business approval or regulatory status.
  - RS-001 is explicitly labeled as testing architectural admissibility, not approving decisions.
- **Remaining gap:** None identified at the specification level. Future implementation should continue to reject language or schema fields that imply approval/authority.
- **Next action:** Add an ADR non-claim test or lint that flags any Admission schema field or documentation that could be read as business/regulatory/trust authority.

### IR-10 — Fail Closed

- **ID:** IR-10
- **Title:** Fail Closed
- **Normative statement:** Unknown capabilities, commands, ownership mappings, or required references MUST NOT silently become admissible.
- **Current status:** PARTIAL
- **Evidence/artifacts:**
  - `testdata/rs_001/run.ps1` rejects unknown capabilities, non-mutating capabilities, ownership mismatches, unknown commands, command/capability mismatches, missing references, and architecture revision mismatches.
  - `schemas/admission_result_v0.1.schema.json` requires `reason_codes` for `REJECT` decisions.
- **Executable evidence/tests:**
  - RS-001 negative cases:
    - `invalid_unknown_capability` → `UNKNOWN_CAPABILITY`
    - `invalid_capability_not_mutating` → `CAPABILITY_NOT_MUTATING`
    - `invalid_ownership_mismatch` → `OWNERSHIP_MISMATCH`
    - `invalid_unknown_command` → `UNKNOWN_COMMAND`
    - `invalid_command_capability_mismatch` → `COMMAND_CAPABILITY_MISMATCH`
    - `invalid_missing_required_field` → `MISSING_REQUIRED_FIELD`
    - `invalid_architecture_revision` → `ARCHITECTURE_REVISION_MISMATCH`
- **Remaining gap:** Fail-closed behavior is demonstrated only in the RS-001 conformance harness, not in a production or CLI-enforced path.
- **Next action:** Expand negative cases to cover additional mutation paths and eventually enforce the same checks at the point where a Core use case is invoked.

### IR-11 — Deterministic Admission

- **ID:** IR-11
- **Title:** Deterministic Admission
- **Normative statement:** Given the same normalized admission input, applicable architecture version, and rule set, Admission MUST return the same result.
- **Current status:** PARTIAL
- **Evidence/artifacts:**
  - `testdata/rs_001/run.ps1` reads static registries (`core-capability-registry.yaml`, `aggregate-ownership-registry.yaml`, `command-event-catalogue.yaml`, `architecture-baseline.yaml`) and produces deterministic decisions and reason codes.
  - `schemas/admission_result_v0.1.schema.json` enforces a stable result shape.
- **Executable evidence/tests:**
  - Repeated runs of `run.ps1` on the same inputs produce identical `decision` and `reason_codes`.
- **Remaining gap:**
  - The `admission_id` derivation rule is undefined; the harness uses a fixture-local `rs-001-<case>` identifier.
  - There is no machine-readable rule registry; `rule_refs` in RS-001 are hard-coded strings (`IR-01`, `IR-03`, `IR-04`, `IR-05`, `IR-10`).
  - No input-normalization rule is specified.
- **Next action:** Define a deterministic `admission_id` derivation and a machine-readable rule registry with stable rule identifiers. Specify input normalization for the Admission gate.

### IR-12 — Traceability

- **ID:** IR-12
- **Title:** Traceability
- **Normative statement:** An admission result SHOULD be traceable to: architecture version → capability → aggregate → command → rule → transition → event/evidence.
- **Current status:** PARTIAL
- **Evidence/artifacts:**
  - `architecture-baseline.yaml` provides the architecture revision.
  - `core-capability-registry.yaml`, `aggregate-ownership-registry.yaml`, and `command-event-catalogue.yaml` provide capability, aggregate, command, and transition mappings.
  - `command-event-catalogue.yaml` links to `event_id: unit.created`.
  - `schemas/admission_result_v0.1.schema.json` includes `rule_refs`, `capability_ref`, `aggregate_ref`, `command_ref`, and `transition_ref`.
- **Executable evidence/tests:**
  - RS-001 `valid_admit` admission result contains `architecture_revision`, `capability_ref`, `aggregate_ref`, `command_ref`, `transition_ref`, and `rule_refs`.
- **Remaining gap:**
  - There is no standalone `rule` registry or `rule_id` scheme; the `rule_refs` field currently uses hand-written strings.
  - The `event/evidence` end of the chain is referenced but not produced in RS-001.
- **Next action:** Create an `admission-rule-registry.yaml` (or equivalent) that binds rule identifiers to the invariant IDs and maps them into the `admission_result_v0.1.schema.json` `rule_refs` field.

## Summary

| Status | Count |
|--------|-------|
| PASS   | 3     |
| PARTIAL| 8     |
| OPEN   | 1     |
| BLOCKED| 0     |

- **Total invariants evaluated:** 12
- **Invariants with executable evidence:** 9 (IR-01, IR-03, IR-04, IR-05, IR-06, IR-07, IR-08, IR-10, IR-11, IR-12) — note that executable evidence for IR-06 is limited to static reference, and for IR-01 full suite is still blocked by environment.
- **Invariants with only documentation:** 1 (IR-02)

## Top 3 remaining P0 gaps by architectural importance

1. **IR-02 — Admission Before Mutation**  
   The constitutional gate is not yet inserted before any existing Core mutation. Without this, Admission remains a documentation/conformance concept rather than an architectural boundary.

2. **IR-12 / IR-11 — Missing rule registry and deterministic `admission_id`**  
   Traceability and deterministic admission both depend on a machine-readable rule registry and stable identifier derivation. This is a shared foundation for IR-11 and IR-12 and must be in place before any enforcement.

3. **IR-06 — Transition Evidence**  
   The chain currently ends at the Command Envelope. An admitted command has not yet been executed by a Core use case, so no runtime event evidence is produced or wrapped in an Event Envelope.

## Recommended next implementation step

Define and commit a machine-readable **P0 Admission Rule Registry** (e.g., `admission-rule-registry.yaml`) that assigns stable, versioned rule identifiers to IR-01 through IR-12 and any sub-rules, and specify the deterministic derivation of `admission_id` from normalized input. This unblocks the IR-11 and IR-12 gaps, and its identifiers then become the prerequisite for wiring Admission into a real Core use case (IR-02) and generating traceable event evidence (IR-06).

Do **not** proceed directly to a production Admission Runner or CLI enforcement before the rule registry and identifier derivation are stable, because the current `rule_refs` are hand-written strings and `admission_id` is fixture-local.
