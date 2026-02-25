# DEC-2026-02-CLI-CONTRACT-LOCK-v1 — CLI Contract Locked (v1)

Status: LOCKED  
Scope: digiemu CLI JSON outputs (`verify --json`, `replay --json`)  
Applies to: Phase 5 contract gate

## Context
Phase 4 delivered deterministic replay + stable verification semantics.
Phase 5 introduces a formal CLI contract (docs + schemas) and contract tests.

We now have:
- Normative docs: `docs/CLI_CONTRACT_v1.0.md`
- Verify JSON schema: `docs/VERIFY_RESULT_SCHEMA_v1.json`
- Golden outputs with normalized trace:
  - `cmd/digiemu/testdata/replay_demo_golden.json`
  - `cmd/digiemu/testdata/verify_demo_golden.json`
- Enforced schema validation test (jsonschema/v5)

## Decision
We lock the CLI JSON contract v1 for:
- `digiemu verify --json`
- `digiemu replay --json`

The following become contract-governed:
1. Top-level JSON field presence and semantics
2. Stable `hash_alg` and `canonical_scope` strings
3. Trace rules:
   - includes file paths read
   - ends with exactly one `used:<bundleRoot>` marker
   - golden tests normalize to repo-relative + `/` for portability
4. Exit codes (normative) remain as documented in `docs/CLI_VERIFY_v1.0.md`

Any change to the above MUST be treated as a breaking change and requires:
- versioning policy compliance (`docs/VERSIONING_POLICY_v1.0.md`)
- updated schema + docs
- updated golden files
- explicit DEC amendment (new DEC or “Nachtrag” section)

## Consequences
- CI now fails on any accidental drift (golden/schema mismatch).
- Contract evolution becomes explicit and reviewable.

## Notes
- This lock is about *observable CLI behavior*, not internal refactors.
- Internal refactors are allowed if they preserve contract outputs.
