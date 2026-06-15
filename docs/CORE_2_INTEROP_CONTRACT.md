# DigiEmu Core 2.0 Interop Contract

- **Status:** Draft
- **Scope:** DigiEmu Core 2.0
- **Purpose:** Defines how external systems may consume DigiEmu decision-state artifacts without redefining DigiEmu state identity.

---

## Purpose

DigiEmu Core 2.0 provides deterministic decision-state verification. It produces canonical decision-state artifacts, cryptographic hashes, replay evidence, and verification reports that can be referenced by external trust, identity, attribution, compliance, or operational systems.

External systems may consume DigiEmu artifacts as evidence of decision-state integrity. They may not compute, redefine, overwrite, or replace DigiEmu state identity.

---

## DigiEmu Produces

DigiEmu Core 2.0 produces the following artifacts:

- **Canonical decision-state snapshots** – deterministic captures of agent decision state
- **Snapshot identifiers** – unique references to specific decision-state snapshots
- **Canonicalization profile identifiers** – references to the declared profile under which state was produced
- **Decision-state hashes** – cryptographic integrity markers for decision state
- **Replay fixtures** – structured inputs and contexts for deterministic replay
- **Verification reports** – structured evidence of verification outcomes
- **PASS / FAIL verification outcomes** – binary results of verification
- **References to structured diagnostic evidence** – links to detailed diagnostic data
- **Cryptographic integrity metadata** – signatures, timestamps, and chain-of-custody markers

**State identity is produced only under a declared DigiEmu canonicalization profile by an accountable state producer.**

---

## External Systems May Consume

External systems may consume and reference DigiEmu artifacts for their own purposes:

- **Agent trust systems** – may reference DigiEmu verification outcomes as part of trust evaluation
- **Identity verification systems** – may correlate DigiEmu state identity with agent identity claims
- **Attribution systems** – may link DigiEmu artifacts to agent actions and outcomes
- **Compliance evidence systems** – may include DigiEmu reports in compliance documentation
- **Deployment validation systems** – may gate deployments on DigiEmu verification outcomes
- **Remediation loops** – may trigger operational workflows based on PASS / FAIL results
- **Audit layers** – may incorporate DigiEmu evidence into audit trails

External systems may carry DigiEmu state identity, snapshot identifiers, hashes, and verification outcomes as evidence. They may not produce, redefine, or overwrite this identity.

---

## External Systems Must Not Redefine

External systems must not:

- **Recompute DigiEmu state identity** under a different profile
- **Overwrite DigiEmu canonicalization semantics** with alternative interpretations
- **Replace structured DigiEmu reports** with unsupported or lossy summaries
- **Claim that DigiEmu verifies agent identity** – identity is outside DigiEmu scope
- **Claim that DigiEmu certifies trust tiers** – trust assignment is external
- **Claim that DigiEmu assigns legal responsibility** – liability is outside DigiEmu scope
- **Claim that DigiEmu performs regulatory approval** – regulatory status is external

External systems may reference DigiEmu verification results, but they must not become a second producer of the same DigiEmu state identity.

---

## Boundary Statement

> One state identity is produced under one declared DigiEmu canonicalization profile by one accountable state producer.
> 
> External systems may carry this identity.
> 
> External systems do not compute, redefine, or overwrite this identity.

---

## Non-Claims

DigiEmu Core 2.0 does not claim to provide:

- **Agent identity verification** – DigiEmu verifies decision state, not agent identity
- **Agent certification** – certification is outside DigiEmu scope
- **Trust-tier assignment** – trust tiers are determined by external systems
- **Legal liability attribution** – legal responsibility is outside DigiEmu scope
- **Regulatory approval** – regulatory status is determined by external authorities
- **Moral judgment** – ethical evaluation is outside DigiEmu scope
- **Model alignment guarantees** – alignment verification requires additional frameworks
- **Full system safety guarantees** – safety is a property of the complete deployed system

**DigiEmu verifies deterministic decision-state integrity and replay consistency within its declared scope.**

---

## TBN / Trust Boundary Example

A TBN-style receipt may reference:

- A DigiEmu snapshot ID
- A DigiEmu moment ID
- DigiEmu decision-state hashes
- A DigiEmu verification report ID
- A DigiEmu PASS / FAIL outcome

TBN may provide:
- Provenance metadata
- Digital signatures
- Agent identity status
- Trust certification

**TBN does not compute or redefine DigiEmu state identity.**

**DigiEmu does not compute or redefine TBN agent trust status.**

The two systems maintain distinct boundaries while interoperating through shared references.

---

## AntifragileOS / Before-Splice Example

An operational remediation flow using DigiEmu verification:

1. **Remediation task generated** – external system identifies a task requiring validation
2. **Task refined** – external system prepares the task with appropriate context
3. **DigiEmu runs deterministic replay** – DigiEmu executes replay against a declared fixture
4. **DigiEmu produces verification report** – structured evidence is generated
5. **PASS allows continuation** – external system may proceed with its own validation flow
6. **FAIL blocks or redirects** – task is returned for correction before proceeding

The external system may act on the DigiEmu result, but it must preserve reference to the structured verification report. The external system does not recompute the verification; it references DigiEmu's output.

---

## PASS / FAIL Preservation Rule

A reduced PASS / FAIL outcome is allowed for portability to systems that require simple binary results.

**However, the reduced outcome must preserve a reference to the structured DigiEmu verification report.**

The verdict must not erase the diagnostic evidence. Any system receiving a reduced outcome must be able to retrieve the full structured report using the preserved reference.

---

## Minimal Interop Fields

The following JSON structure represents the minimal fields for interoperability:

```json
{
  "digiemu_core_version": "2.0",
  "canonicalization_profile": "digiemu-core-2-profile",
  "moment_id": "agent-moment-001",
  "snapshot_id": "digiemu-snapshot-001",
  "verification_report_id": "digiemu-report-001",
  "verification_outcome": "PASS",
  "decision_state_hashes": [],
  "structured_report_ref": "reports/digiemu-report-001.json"
}
```

---

## Design Principle

DigiEmu is a decision-state verification layer.

It is designed to be consumed by trust, identity, attribution, compliance, and operational systems **without** collapsing those systems into DigiEmu and **without** allowing those systems to redefine DigiEmu artifacts.

**The boundary is the value.**
