# DigiEmu Core 2.0 — Migration from v1.0 (Draft)

Status: DRAFT

Purpose: Describe how Core 2.0 evolves from v1.0 while preserving compatibility and avoiding breaking changes.

Guiding constraints
-------------------
- DO NOT rewrite the existing architecture.  
- DO NOT break the existing v1.0 CLI contract.  
- DO NOT remove or rename existing v1.0 documents.

Compatibility strategy (normative)
---------------------------------
1. Non-breaking documentation first: Publish these drafts to record planned changes before any code is modified.  
2. Feature-flagged experimentation: Any runtime behavior changes (e.g., optional `Verify Result` v2) SHALL be behind opt-in flags or explicit new CLI subcommands.
3. Separate test vectors: v2 conformance vectors SHALL be stored separately from v1 vectors. v1 vectors MUST remain unchanged and executable.
4. Additive metadata: New fields (e.g., `crypto_profile`) SHALL be optional and MUST NOT be included in inside-hash payloads unless explicitly part of a new, versioned canonicalization profile.

Stepwise migration recommendations
---------------------------------
- Step A — Publish docs and SPEC_INDEX entries (this change).  
- Step B — Add v2 schema drafts and separate test vectors; run v1 test suite to ensure no regressions.  
- Step C — Implement opt-in v2 output mode in a feature branch; keep default behavior identical to v1.0.  
- Step D — If v2 proves stable, offer a documented upgrade path; core team SHOULD provide tooling to translate v2 optional fields for systems expecting v1 outputs.

Testing and acceptance
----------------------
- Acceptance tests MUST include the complete v1.0 test suite and a v2 smoke test set.  
- All changes that might affect CLI outputs MUST be gated behind automated regression checks.

Operational notes
-----------------
- Communication: Release notes SHALL clearly state that v1.0 proofs remain valid and how to enable v2 experimental outputs.  
- Backwards compatibility: Consumers SHALL be informed that v2 fields are optional; systems that rely on strict v1 outputs MUST continue to operate unchanged.
