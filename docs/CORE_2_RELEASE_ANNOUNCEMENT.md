# DigiEmu Core 2.0 Release Announcement

- **Status:** Released
- **Release date:** 15 June 2026
- **Tag:** `core-2.0.0`

---

## Summary

DigiEmu Core 2.0 defines a deterministic decision-state verification layer for AI governance, replay, interoperability, and audit evidence.

It provides canonical artifacts that external trust, identity, attribution, compliance, and operational systems may reference without redefining DigiEmu state identity.

---

## What Core 2.0 Provides

DigiEmu Core 2.0 provides:

- **Deterministic decision-state snapshots** – canonical captures of agent decision state
- **Canonicalization profile boundary** – explicit declared profile under which state is produced
- **Decision-state hash references** – cryptographic integrity markers
- **Replay consistency evidence** – proof that decision state can be reproduced deterministically
- **Structured verification reports** – machine-readable evidence of verification outcomes
- **PASS / FAIL outcome preservation** – portable binary results with preserved links to full evidence
- **Interop examples for external systems** – minimal examples for TBN, AntifragileOS, CLARIXO, compliance, and audit integrations
- **External review material** – documentation supporting external reviewer evaluation

---

## What Core 2.0 Does Not Claim

DigiEmu Core 2.0 does not claim to provide:

- **Agent identity verification** – DigiEmu verifies decision state, not who produced it
- **Agent certification** – Certification is outside DigiEmu scope
- **Trust-tier assignment** – Trust tiers are determined by external systems
- **Legal liability attribution** – Legal responsibility assignment is external
- **Regulatory approval** – Regulatory status is determined by external authorities
- **Full system safety guarantee** – Safety is a property of the complete deployed system

DigiEmu verifies deterministic decision-state integrity and replay consistency within its declared scope.

---

## Interoperability

DigiEmu Core 2.0 is designed to interoperate with external systems as a verification layer:

- **TBN-style trust systems** – May reference DigiEmu evidence without computing or redefining DigiEmu state identity
- **AntifragileOS-style operational validation** – May use PASS/FAIL outcomes for before-splice validation without DigiEmu becoming the deployment authority
- **CLARIXO-style attribution systems** – May use DigiEmu artifacts as upstream evidence for responsibility tracking
- **Compliance layers** – May incorporate DigiEmu reports into compliance documentation
- **Audit partners** – May reference DigiEmu evidence in audit trails

---

## Release Materials

- [Core 2.0 Interop Contract](CORE_2_INTEROP_CONTRACT.md) – Defines the interoperability boundary
- [Core 2.0 External Review Note](CORE_2_EXTERNAL_REVIEW_NOTE.md) – Guidance for external reviewers
- [Core 2.0 Interop Examples](../examples/interop/README.md) – Minimal illustrative examples

---

## Closing Principle

**The boundary is the value.**
