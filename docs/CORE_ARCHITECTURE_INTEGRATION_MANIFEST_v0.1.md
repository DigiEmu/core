# DigiEmu Core — Architecture Integration Manifest v0.1

- **Status:** Proposed
- **Target:** DigiEmu Core 2.0.x
- **Base:** released Core 2.0.0
- **Architecture Layer:** P0 Integration
- **Date:** 2026-08-07

---

## PURPOSE

The manifest defines how the new DigiEmu Architecture Constitution integrates with the existing released Core 2.0 architecture.

Integration is organized into three classes:

- **KEEP** — Existing Core semantics and contracts remain unchanged.
- **BIND** — Existing Core capabilities are explicitly mapped/bound to architectural ownership, commands, transitions, and evidence.
- **ADD** — Missing constitutional and machine-readable architecture artifacts are introduced.

---

## CORE DECISION

> DigiEmu Core 2.0 remains the deterministic state-verification kernel. The Architecture Constitution adds a deterministic admission boundary for constitution-governed state transitions without redefining Core state identity, verification semantics, or external authority.

---

## CONSTITUTIONAL BOUNDARY

The Architecture Constitution defines **Admission**. Admission is distinct from the following:

- **Admission != Verification**  
  Verification is the deterministic replay and hash comparison performed by the existing DigiEmu Core kernel. Admission does not replace or redefine verification.

- **Admission != Business Approval**  
  Admission does not determine whether a business decision is correct, desirable, or authorized by an organization.

- **Admission != Regulatory Approval**  
  Admission does not grant or imply regulatory status, certification, or compliance.

- **Admission != Trust Assignment**  
  Admission does not assign trust tiers, agent trust status, or reliance metrics.

Admission determines **architectural admissibility** of constitution-governed Core transitions. It MUST NOT determine whether the substantive business decision is correct, legally valid, ethically correct, trustworthy, or approved by an external authority.

---

## TARGET FLOW

The following conceptual chain describes the proposed P0 integration flow. It is a target model, not a claim that the entire flow is already implemented in released Core 2.0.0.

```
External Actor / AI / Middleware
    → Intent
    → Core Admission Constitution
    → admitted Command
    → Existing Core 2.0 Kernel / Use Case
    → Transition
    → Event / Audit Evidence
    → Snapshot / Canonicalization
    → Hash / Replay / Verify
    → Evidence / Receipt
```

The chain explicitly separates:

- **Existing Core 2.0 behavior:** snapshot canonicalization, hashing, replay, verification, and evidence production (`docs/VERIFY_SPEC_v1.0.md`, `docs/SNAPSHOT_HASH_v1.0.md`, `internal/verify/`).
- **Proposed constitutional integration:** admission, intent/command/event envelopes, and aggregate ownership (`architecture-baseline.yaml`, `Core Capability Registry`, `Aggregate Ownership Registry`, `Command/Event Catalogue`).

---

## KEEP

The following existing Core 2.0 semantics and contracts remain unchanged and authoritative:

- Canonicalization (`docs/DIGIEMU_CORE_SPEC_v0.9.md`, `docs/SNAPSHOT_HASH_v1.0.md`, `docs/CORE_2_PROFILE_REGISTRY.md`).
- Snapshot formation (`docs/SNAPSHOT_HASH_v1.0.md`, `internal/snapshot/`).
- Snapshot hashing (`docs/SNAPSHOT_HASH_v1.0.md`, `internal/verify/`).
- Replay (`docs/VERIFY_SPEC_v1.0.md`, `internal/verify/`, `pkg/replay/`).
- Verification (`docs/VERIFY_SPEC_v1.0.md`, `internal/verify/`, `pkg/verify/`).
- Verify Result contracts (`docs/VERIFY_SPEC_v1.0.md`, `schemas/VERIFY_RESULT_SCHEMA_v1.json`, `schemas/verify_result_v2.schema.json`).
- Existing domain model (`internal/kernel/domain/`).
- Existing audit semantics (`internal/kernel/domain/audit_event.go`, `docs/VERIFY_SPEC_v1.0.md`).
- CLI v1 contract (`docs/CLI_CONTRACT_v1.0.md`, `docs/CLI_VERIFY_v1.0.md`, `cmd/digiemu/main.go`).
- Existing public APIs and schemas.
- Compatibility policy (`docs/VERSIONING_POLICY_v1.0.md`).
- Interop boundary (`docs/CORE_2_INTEROP_CONTRACT.md`, `docs/CORE_2_EXTERNAL_REVIEW_NOTE.md`).
- Security boundary (`SECURITY.md`, `docs/THREAT_MODEL.md`, `docs/GOVERNANCE.md`).

The new Architecture Constitution MUST NOT create a second definition of DigiEmu state identity. State identity is produced only under a declared DigiEmu canonicalization profile by an accountable state producer, as stated in `docs/CORE_2_INTEROP_CONTRACT.md`.

---

## BIND

This manifest proposes a future, controlled binding of the following existing Core artifacts to constitutional concepts:

- Existing kernel use cases (`internal/kernel/usecases/`).
- Existing domain concepts (`internal/kernel/domain/`).
- Existing `AuditEvent` / evidence mechanisms (`internal/kernel/domain/audit_event.go`, `schemas/VERIFY_RESULT_SCHEMA_v1.json`).
- Existing conformance runner (`internal/conformance/runner.go`, `docs/CORE_2_CONFORMANCE_PACK.md`).

The proposed conceptual mapping is:

```
Capability
    → Aggregate Owner
    → Command
    → Admission Rule
    → Existing Use Case
    → Transition
    → Event / Evidence
```

This mapping is descriptive and forward-looking. It does not rename or refactor existing Go packages. Existing package and command names remain stable under the CLI v1 contract.

**Important:** The CLI surface is an interface to Core capabilities, but the CLI surface is not identical to the architectural capability model. The capability model is a separate registry artifact.

---

## ADD

The following architecture artifacts are proposed for addition:

1. `architecture-baseline.yaml` — human- and machine-readable architecture layer map.
2. `01_decisions/ADR-P0.md` — ADR P0 registry.
3. `01_decisions/ADR-P0-05-Core-Admission-Constitution.md` — P0 admission constitution decision.
4. `docs/core-capability-registry.md` — registry of Core capabilities.
5. `docs/aggregate-ownership-registry.md` — mapping of aggregates to owners and use cases.
6. `docs/command-event-catalogue.md` — catalogue of commands, events, and transitions.
7. `schemas/intent-envelope.schema.json` — schema for external intent.
8. `schemas/command-envelope.schema.json` — schema for an admitted command.
9. `schemas/event-envelope.schema.json` — schema for evidence of a transition.
10. `docs/RS-001-Decision-Approval.md` — executable reference scenario for the Admission Constitution.

**ADR-P0-05** MUST remain the **Core Admission Constitution**. It MUST NOT be reinterpreted as a canonicalization ADR; canonicalization is already governed by `docs/SNAPSHOT_HASH_v1.0.md` and related specifications.

---

## ENVELOPE SEMANTICS

### Intent Envelope

The **Intent Envelope** represents an external request. An intent does **not** imply admission. It carries what an external actor wants to happen, not whether it is admissible under the Architecture Constitution.

### Command Envelope

The **Command Envelope** represents an operation that has been admitted for execution under the applicable Architecture Constitution. Admission is the gate between Intent and Command.

### Event Envelope

The **Event Envelope** represents or references evidence of the resulting transition. It MUST NOT replace existing Core verification-result schemas (`schemas/VERIFY_RESULT_SCHEMA_v1.json`, `schemas/verify_result_v2.schema.json`, `schemas/core_2_conformance_report.schema.json`). It is an additional wrapping or referencing mechanism.

---

## RS-001

**RS-001** is the first executable reference scenario for the Admission Constitution. It is a test and documentation artifact, not a runtime feature in released Core 2.0.0.

"Decision Approval" in RS-001 MUST NOT mean that DigiEmu substantively approves a business decision. RS-001 tests **architectural admissibility and transition integrity** of a decision-approval-shaped command under the Admission Constitution.

### Proposed RS-001 conceptual flow

```
Approval Intent
    → Envelope Validation
    → Admission Evaluation
    → Admitted Approval Command
    → Decision Aggregate
    → Defined Transition
    → Approval Event
    → Evidence
```

The existence of a "Decision Aggregate" in this scenario is **proposed**, not a claim that a native Decision aggregate exists in the current Core 2.0 repository. The current repository does not establish a runtime Decision aggregate. RS-001 is an integration-level reference scenario, not a claim of existing implementation.

---

## INTEGRATION INVARIANTS

### IR-01 — No Runtime Regression

Architecture integration MUST NOT silently alter existing Core 2.0 behavior or stable contracts.

### IR-02 — Admission Before Mutation

A mutation governed by the Architecture Constitution MUST pass Admission before execution.

### IR-03 — Explicit Capability

Every constitution-governed command MUST reference a registered capability.

### IR-04 — Explicit Ownership

Every constitution-governed mutation MUST resolve to an authorized aggregate owner.

### IR-05 — Command/Transition Correspondence

Every admitted mutating command MUST resolve to a defined transition.

### IR-06 — Transition Evidence

Every successfully executed constitution-governed transition MUST produce or reference deterministic evidence sufficient to identify the resulting transition.

### IR-07 — Verification Independence

Admission MUST NOT redefine canonicalization, hashing, reconstruction, or verification.

### IR-08 — Single State Identity

The architecture layer MUST NOT introduce an alternative DigiEmu state-identity mechanism.

### IR-09 — External Authority Preservation

Admission MUST NOT absorb business, legal, regulatory, ethical, or trust authority belonging outside DigiEmu Core.

### IR-10 — Fail Closed

Unknown capabilities, commands, ownership mappings, or required references MUST NOT silently become admissible.

### IR-11 — Deterministic Admission

Given the same normalized admission input, applicable architecture version, and rule set, Admission MUST return the same result.

### IR-12 — Traceability

An admission result SHOULD be traceable to:

```
architecture version
    → capability
    → aggregate
    → command
    → rule
    → transition
    → event/evidence
```

---

## MIGRATION

The integration is proposed to proceed in five phases:

### Phase A — Formalize

Introduce the architecture baseline, ADRs, and registries. No runtime changes.

### Phase B — Specify

Introduce envelope schemas and the RS-001 scenario. No general mutation interception.

### Phase C — Prove

Execute RS-001 and negative cases through a conformance path. Evidence is generated, not assumed.

### Phase D — Bind

Map selected existing Core use cases to constitutional commands and transitions. Existing package names remain unchanged.

### Phase E — Enforce

Only after successful conformance may Admission become mandatory for selected constitution-governed mutation paths.

### Migration principle

```
Document → Map → Test → Enforce
```

Admission MUST NOT be enforced in production until after Phase C and Phase D are complete and reviewed.

---

## EVIDENCE PRINCIPLE

Architecture claims are not proven merely because they are documented. Evidence is produced through the following chain:

```
Invariant
    → Specification
    → Test Vector
    → Execution
    → Expected Result
    → Actual Result
    → Evidence
```

The existing conformance runner (`internal/conformance/runner.go`) and conformance test data (`testdata/core_2_conformance/`) provide the mechanical pattern for this evidence chain, but the Admission Constitution evidence is future work.

---

## OPEN QUESTIONS

The following questions are unresolved. Answers are not invented in this manifest.

1. Which existing Core use case should become the first real Admission-bound transition?
2. Is Decision a native Core aggregate, a reference aggregate, or an integration-level aggregate?
3. What exact evidence object represents Admission itself?
4. Which Architecture Baseline elements become normative and when?
5. How do existing CLI mutation paths coexist with constitution-governed paths during migration?

---

## REFERENCES

- `core-2.0.0` release tag (released Core 2.0.0).
- `docs/DIGIEMU_CORE_SPEC_v0.9.md` — Core v0.9 specification.
- `docs/CLI_CONTRACT_v1.0.md` — stable `digiemu` CLI contract.
- `docs/VERIFY_SPEC_v1.0.md` — deterministic verification procedure.
- `docs/SNAPSHOT_HASH_v1.0.md` — canonical snapshot hashing rules.
- `docs/CORE_2_BOUNDARY_MODEL.md` — Core 2.0 boundary model.
- `docs/CORE_2_INTEROP_CONTRACT.md` — interop contract and non-claims.
- `docs/CORE_2_EXTERNAL_REVIEW_NOTE.md` — external review guidance.
- `docs/CORE_2_CONFORMANCE_PACK.md` — conformance pack.
- `docs/VERSIONING_POLICY_v1.0.md` — compatibility policy.
- `docs/GOVERNANCE.md`, `docs/ETHICS.md`, `SECURITY.md`, `docs/THREAT_MODEL.md` — governance and security boundaries.
- `internal/kernel/` — kernel use cases and domain model.
- `internal/conformance/runner.go` — conformance runner.
- `schemas/VERIFY_RESULT_SCHEMA_v1.json` — v1 verify result schema.
- `schemas/verify_result_v2.schema.json` — v2 verify result schema.
- `schemas/core_2_conformance_report.schema.json` — conformance report schema.

---

## GUARDRAILS

- Do not invent runtime behavior.
- Do not modify existing contracts.
- Do not claim Admission is currently enforced in Core 2.0.
- Do not turn DigiEmu into a business-decision authority.
- Do not create new canonicalization or state-identity semantics.
- Clearly label proposed/future integration separately from released Core behavior.
- Use normative MUST/SHALL language only for the proposed Architecture Constitution and its integration constraints, not to falsely describe unimplemented Core 2.0 runtime behavior.
