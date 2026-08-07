# ADR-P0-05 — Core Admission Constitution

- **Decision ID:** ADR-P0-05
- **Version:** 0.1
- **Status:** Proposed
- **Architecture layer:** P0
- **Target release:** DigiEmu Core 2.0.x
- **Date:** 2026-08-07

---

## 1. Context

DigiEmu Core 2.0 provides a deterministic decision-state verification kernel. It produces canonical decision-state artifacts, cryptographic hashes, replay evidence, and verification reports. These outputs are consumed by external trust, identity, attribution, compliance, and operational systems.

The Core 2.0 boundary is already defined in `docs/CORE_2_INTEROP_CONTRACT.md`, `docs/CORE_2_BOUNDARY_MODEL.md`, and `docs/CORE_2_EXTERNAL_REVIEW_NOTE.md`. The kernel remains responsible for deterministic state verification; external systems remain responsible for trust, legal, regulatory, ethical, and business decisions.

This ADR introduces a P0 admission constitution: a deterministic, constitution-governed boundary for state transitions that are admitted into the DigiEmu Core boundary. The constitution does not alter Core 2.0 runtime behavior, redefine state identity, or assume external authority.

---

## 2. Existing Core 2.0 Boundary

The existing boundary is described by:

- `docs/CORE_2_INTEROP_CONTRACT.md` — how external systems reference DigiEmu artifacts without redefining state identity.
- `docs/DIGIEMU_CORE_SPEC_v0.9.md` — the Core v0.9 specification for snapshots, bundles, replay, and verification.
- `docs/VERIFY_SPEC_v1.0.md` — the deterministic verification procedure.
- `docs/CLI_CONTRACT_v1.0.md` — the stable `digiemu` CLI contract.
- `docs/SNAPSHOT_HASH_v1.0.md` — the canonical snapshot hashing rules.
- `docs/CONFORMANCE.md` and `docs/CORE_2_CONFORMANCE_PACK.md` — conformance requirements and test pack.
- `internal/kernel/usecases/` and `internal/verify/` — the existing verification and replay use cases.

Core 2.0 does not perform agent identity verification, certification, trust-tier assignment, legal liability attribution, regulatory approval, or moral judgment.

---

## 3. Problem

As DigiEmu Core 2.0 is integrated into larger systems, there is a need to make explicit:

- Which state transitions are architecturally admissible into the DigiEmu Core boundary.
- What invariants must hold for a transition to be admitted.
- That admission is a boundary rule, not a business, legal, or governance decision.
- That the constitution does not redefine the existing canonicalization, hashing, replay, or verification semantics.

Without this ADR, the boundary between "what DigiEmu Core verifies" and "what a constitution allows to enter the Core boundary" remains implicit.

---

## 4. Decision

> DigiEmu Core 2.0 remains the deterministic state-verification kernel. The Architecture Constitution adds a deterministic admission boundary for constitution-governed state transitions without redefining Core state identity, verification semantics, or external authority.

Admission is the architectural gate by which a constitution-governed state transition is allowed to enter the DigiEmu Core boundary. DigiEmu Core continues to verify the resulting state deterministically. The constitution does not produce state identity; it governs whether a transition is admissible under declared rules.

---

## 5. Constitutional Invariants

The following invariants are defined in `docs/CORE_ARCHITECTURE_INTEGRATION_MANIFEST_v0.1.md` as IR-01 through IR-12.

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

## 6. KEEP

The following remain unchanged and authoritative:

- `digiemu` CLI v1 contract and exit codes (`docs/CLI_CONTRACT_v1.0.md`, `docs/CLI_VERIFY_v1.0.md`).
- Canonicalization profile `digiemu-canonical-json-v1` (`docs/SNAPSHOT_HASH_v1.0.md`, `docs/CORE_2_PROFILE_REGISTRY.md`).
- SHA-256 snapshot hashing (`docs/SNAPSHOT_HASH_v1.0.md`, `internal/verify/`).
- Deterministic replay and verification semantics (`docs/VERIFY_SPEC_v1.0.md`, `internal/verify/`).
- Verify Result v1 and v2 schemas (`schemas/VERIFY_RESULT_SCHEMA_v1.json`, `schemas/verify_result_v2.schema.json`).
- Conformance pack and runner behavior (`docs/CORE_2_CONFORMANCE_PACK.md`, `internal/conformance/runner.go`).
- Core 2.0 interop boundary and non-claims (`docs/CORE_2_INTEROP_CONTRACT.md`).

---

## 7. BIND

The following are bound by this ADR:

- **Admission is deterministic.** A transition is admissible if and only if it satisfies the constitutional invariants and the existing Core verification rules.
- **Admission is not approval.** The constitution does not make business, legal, regulatory, ethical, or trust-tier decisions. It governs architectural admissibility.
- **State identity remains internal to DigiEmu.** The constitution may carry, reference, or inspect state identity; it does not produce or overwrite it.
- **External authority remains external.** Business, legal, regulatory, and ethical authorities operate outside the DigiEmu Core boundary and consume DigiEmu artifacts as evidence.
- **All admitted transitions are replay-verifiable.** Any state admitted into the Core boundary must remain reproducible under the existing replay and verification procedure.

---

## 8. ADD

This ADR adds:

- The term **"Core Admission"** as the architectural gate for constitution-governed state transitions.
- The term **"Constitutional Invariants"** as the rules that admission must satisfy.
- The distinction between **admissibility** (architectural) and **approval** (business/governance).
- A P0 architecture decision record that does not modify Core 2.0 runtime code.

No new CLI commands, flags, schemas, or runtime behavior are added by this ADR.

---

## 9. Admission Semantics

### 9.0 Constitutional Boundary

Admission is distinct from the following:

- **Admission != Verification**  
  Verification is the deterministic replay and hash comparison performed by the existing DigiEmu Core kernel. Admission does not replace or redefine verification.

- **Admission != Business Approval**  
  Admission does not determine whether a business decision is correct, desirable, or authorized by an organization.

- **Admission != Regulatory Approval**  
  Admission does not grant or imply regulatory status, certification, or compliance.

- **Admission != Trust Assignment**  
  Admission does not assign trust tiers, agent trust status, or reliance metrics.

Admission determines **architectural admissibility** of constitution-governed Core transitions. It MUST NOT determine whether the substantive business decision is correct, legally valid, ethically correct, trustworthy, or approved by an external authority.

### 9.1 Admission as boundary rule

A state transition is admitted when:

1. The transition is well-formed under the existing Core snapshot/bundle schema.
2. The transition is produced or referenced under a declared canonicalization profile.
3. The transition satisfies the constitutional invariants (IR-01 through IR-12) defined in `docs/CORE_ARCHITECTURE_INTEGRATION_MANIFEST_v0.1.md`.
4. The transition does not conflict with the existing verification semantics.

### 9.2 Admission is not a substitute for verification

Admission governs whether a transition may enter the Core boundary. Verification is the deterministic replay and hash comparison performed by the Core kernel. An admitted transition may still fail verification (e.g., hash mismatch, reconstruction rule mismatch, transition invalid).

### 9.3 Admission preserves external authority

A constitution may require that an admitted transition carries external attestations (e.g., signatures, provenance, audit markers). DigiEmu Core treats these as **outside-hash metadata** and does not include them in state identity computation, consistent with `docs/SNAPSHOT_HASH_v1.0.md` and `docs/VERIFY_SPEC_v1.0.md`.

### 9.4 RS-001 Reference Scenario

**RS-001** is the first executable reference scenario for the Admission Constitution. It is a test and documentation artifact, not a runtime feature in released Core 2.0.0.

"Decision Approval" in RS-001 MUST NOT mean that DigiEmu substantively approves a business decision. RS-001 tests **architectural admissibility and transition integrity** of a decision-approval-shaped command under the Admission Constitution.

The conceptual flow is:

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

## 10. Consequences

- **Positive:** The boundary between "what a constitution admits" and "what DigiEmu verifies" becomes explicit.
- **Positive:** Integration partners can reason about admission rules without asking DigiEmu Core to make business or governance decisions.
- **Neutral:** No existing CLI, schema, or runtime behavior changes.
- **Risk:** If the constitutional invariants are not precisely defined, partners may misinterpret admission as governance approval. This is mitigated by the explicit distinction in §9 and the non-claims in `docs/CORE_2_INTEROP_CONTRACT.md`.

---

## 11. Non-Goals

This ADR explicitly does **not**:

- Redefine `digiemu-canonical-json-v1`, SHA-256 hashing, or snapshot hash boundary.
- Change the `digiemu` CLI v1 contract or exit codes.
- Modify the Core 2.0 conformance pack, runner, or test vectors.
- Make DigiEmu Core into a business, legal, regulatory, ethical, or trust decision maker.
- Introduce a new canonicalization profile.
- Replace or supersede `docs/CORE_2_INTEROP_CONTRACT.md`.

---

## 12. Migration

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

## 13. Conformance

This ADR is conformant with Core 2.0 if:

- The `digiemu verify` and `digiemu replay` commands continue to produce the same outputs for existing test vectors.
- Existing conformance test cases in `testdata/core_2_conformance/` continue to pass.
- The CLI v1 contract golden outputs remain byte-exact where already stable.
- No existing schema (`VERIFY_RESULT_SCHEMA_v1.json`, `verify_result_v2.schema.json`, `core_2_conformance_report.schema.json`) is modified.

---

## 14. Open Questions

The following questions are unresolved. Answers are not invented in this ADR.

1. Which existing Core use case should become the first real Admission-bound transition?
2. Is Decision a native Core aggregate, a reference aggregate, or an integration-level aggregate?
3. What exact evidence object represents Admission itself?
4. Which Architecture Baseline elements become normative and when?
5. How do existing CLI mutation paths coexist with constitution-governed paths during migration?

---

## 15. References

- `docs/CORE_2_INTEROP_CONTRACT.md`
- `docs/CORE_2_BOUNDARY_MODEL.md`
- `docs/CORE_2_EXTERNAL_REVIEW_NOTE.md`
- `docs/DIGIEMU_CORE_SPEC_v0.9.md`
- `docs/VERIFY_SPEC_v1.0.md`
- `docs/CLI_CONTRACT_v1.0.md`
- `docs/CLI_VERIFY_v1.0.md`
- `docs/SNAPSHOT_HASH_v1.0.md`
- `docs/CORE_2_PROFILE_REGISTRY.md`
- `docs/CORE_2_CONFORMANCE_PACK.md`
- `schemas/VERIFY_RESULT_SCHEMA_v1.json`
- `schemas/verify_result_v2.schema.json`
- `schemas/core_2_conformance_report.schema.json`
- `internal/verify/`
- `internal/conformance/runner.go`
- `01_decisions/DEC-2026-02-CLI-CONTRACT-LOCK-v1.md`
- `docs/CORE_ARCHITECTURE_INTEGRATION_MANIFEST_v0.1.md` — source of IR-01 through IR-12, integration classes, and RS-001 semantics
