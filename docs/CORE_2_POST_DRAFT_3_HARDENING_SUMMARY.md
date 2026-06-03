# Core 2.0 Post-Draft.3 Hardening Summary

**Milestone:** post-v2.0.0-draft.3  
**Stability:** draft / pre-release / not stable  
**Current conformance pack:** 11 cases, expected total=11 passed=11 failed=0

This document summarizes Core 2.0 hardening work completed after `v2.0.0-draft.3`. It is a documentation checkpoint only. It does not promote Core 2.0 to stable and does not change production code, schemas, testdata, CLI behavior, CI behavior, tags, or v1.0 behavior.

## Completed After Draft 3

- Actual-vs-Expected Comparison Design
- Observed-vs-Expected Comparison Helper
- Runner integration of observed-vs-expected comparison
- Unsupported canonicalization profile conformance case
- Conformance total updated from 10 to 11

## Current Strengths

- Runner comparison semantics are now documented and partially implemented.
- Observed-vs-expected comparison is centralized where observed results exist.
- Malformed JSON remains deterministic.
- Public JSON report shape remains unchanged.
- `report_version` remains `core-2-conformance-report-v1`.
- v1.0 behavior remains unchanged.
- `v1.0.0` remains stable/latest.

## Current Validation

- `go test ./...` passes.
- Experimental conformance CLI reports total=11 passed=11 failed=0.
- Experimental conformance CLI `--json` reports status PASS, total 11, passed 11, failed 0.

## Non-Goals

- No stable Core 2.0 release yet.
- No release candidate yet.
- No full verify execution yet.
- No HTTP API server.
- No signing or key custody.
- No authentication or tenancy.
- No policy judgement.

## Recommended Next Steps

- Use `docs/CORE_2_MISSING_REFERENCE_SCHEMA_CONSIDERATION.md` to plan whether `MISSING_REFERENCE` belongs in Verify Result v2 before adding schema support or a conformance case.
- Use `docs/CORE_2_DRAFT_4_READINESS.md` as the next milestone planning document before any future Draft 4 tag.
- Consider another draft milestone only after additional conformance coverage.
- Prepare full verify execution design/phase only after runner semantics remain stable.
- Collect partner feedback on the 11-case conformance pack.

## Compatibility Statement

Core 2.0 remains draft and pre-release. The current latest stable line remains `v1.0.0`, and this post-draft.3 hardening checkpoint preserves v1.0 compatibility and the existing experimental Core 2.0 report contract.
