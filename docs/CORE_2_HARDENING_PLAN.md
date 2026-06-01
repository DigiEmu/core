# DigiEmu Core 2.0 — Hardening Plan (Draft)

Status: DRAFT

Purpose: Roadmap and prioritized tasks to harden DigiEmu Core from v1.0 to a Core 2.0 hardening draft while preserving v1.0 contracts.

Principles (normative)
----------------------
- Preserve v1.0 compatibility: all changes MUST be opt-in and MUST NOT alter existing CLI behavior or document names.
- Minimal surface: Core SHALL remain small and deterministic; any non-core feature SHALL be deferred to Secure/Enterprise/Domain Apps.
- Auditability: changes SHOULD be documented with acceptance criteria and test vectors.
- Crypto agility: Core SHALL reference cryptographic profiles but MUST NOT perform key custody.

Phased roadmap (high level)
---------------------------
Phase 1 — Documentation & Boundaries (this draft)
- Deliver boundary model, verify-result v2 draft, migration guidance, and crypto-agility statement (this set of documents).  
- Acceptance: documentation files exist and SPEC_INDEX updated.

Phase 2 — Verify Result v2 Pilots (non-breaking)
- Implement optional v2 output mode behind a feature flag.  
- Add machine-readable conformance tests for v2 (SHOULD be separate from v1 vectors).

Phase 3 — Hash Agility Hooks
- Add support for multi-hash validation (e.g., accept additional verification layers) in a way that does NOT invalidate historical proofs.  
- Provide canonical migration patterns and examples.

Phase 4 — Secure Layer Integration (out-of-core)
- Define APIs and formats for DigiEmu Secure to provide signatures, sealing, and post-quantum migration strategies.  
- Ensure Secure can consume Core outputs without altering them.

Phase 5 — Acceptance and Release
- Publish acceptance tests, update SPEC_INDEX, and create a release note clarifying compatibility guarantees.

Risk & mitigations
------------------
- Risk: accidental change to v1.0 CLI behavior.  
  Mitigation: require CLI contract smoke tests and preserve behavior behind flags.
- Risk: mixing outside-hash metadata into inside-hash payloads.  
  Mitigation: enforce validation in test vectors and document canonical boundaries.

Prioritized next steps (short)
1. Publish these draft docs (done by creating files).  
2. Add v2 `Verify Result` schema draft and run separate conformance vectors.  
3. Implement an opt-in v2 output mode for experimentation in a feature branch.  
