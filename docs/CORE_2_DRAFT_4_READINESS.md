# Core 2.0 Draft 4 Readiness Notes

Status
------

- Milestone: planned `v2.0.0-draft.4`
- Stability: draft / pre-release / not stable
- Purpose: Prepare a future Draft 4 tag without creating it yet.

Purpose
-------

This document summarizes readiness for a future `v2.0.0-draft.4` milestone. It is documentation-only and does not create or move tags, promote Core 2.0 draft features to stable, change v1.0 behavior, change production code, change tests, change schemas, change testdata, change CI behavior, or modify CLI behavior.

Changes since `v2.0.0-draft.3`
-------------------------------

- Actual-vs-Expected Comparison Design
- Observed-vs-Expected Comparison Helper
- Runner integration of observed-vs-expected comparison
- Unsupported canonicalization profile conformance case
- Conformance pack expanded from 10 to 11 cases
- Post-Draft.3 Hardening Summary
- `MISSING_REFERENCE` pending documentation
- `MISSING_REFERENCE` schema consideration note

Why Draft 4 is justified soon
-----------------------------

- Runner comparison semantics have moved from design-only planning into partial implementation where observed results exist.
- The official conformance pack now includes an unsupported canonicalization profile case and remains fully passing at 11 cases.
- The public JSON report shape and `report_version` remain unchanged while conformance coverage has improved.
- `MISSING_REFERENCE` has been explicitly kept out of the schema until semantics are stable, reducing premature partner-contract risk.
- Post-draft.3 hardening is summarized and ready for partner review before any future release-candidate planning.

Current strengths
-----------------

- Runner comparison semantics are documented and partially implemented.
- Observed-vs-expected comparison is centralized where observed results exist.
- Official conformance pack now has 11 passing cases.
- `UNSUPPORTED_CANONICALIZATION_PROFILE` is represented by an official conformance case.
- `MISSING_REFERENCE` is not forced into the schema and is documented as pending.
- Public JSON report shape remains unchanged.
- `report_version` remains `core-2-conformance-report-v1`.
- `v1.0.0` remains stable/latest.

Before tagging Draft 4
----------------------

- Run `go test ./...`
- Run experimental conformance CLI human-readable
- Run experimental conformance CLI `--json`
- Optionally run `powershell -ExecutionPolicy Bypass -File scripts/audit_core2_draft4.ps1` as a pre/post-tag verification helper.
- Confirm output remains total=11 passed=11 failed=0
- Confirm GitHub Actions are green
- Confirm README reflects 11-case draft status
- Confirm no stable `v2.0.0` promotion
- Confirm v1.0 behavior remains unchanged
- Confirm no unexpected files are modified

Post-tag review-fix note
------------------------

- After `v2.0.0-draft.4` was tagged, a follow-up review-fix on main updated the basic conformance report fixture to the current 11-case pack and added a semantic regression assertion.
- See `docs/CORE_2_DRAFT_4_POST_TAG_NOTE.md`.
- Do not move or recreate the `v2.0.0-draft.4` tag.

Proposed tag
------------

- Name: `v2.0.0-draft.4`
- Title: DigiEmu Core 2.0 Draft 4 — Post-Draft.3 Runner Hardening Milestone
- Meaning: A draft/pre-release milestone reflecting post-draft.3 runner comparison integration, 11-case conformance coverage, unsupported canonicalization profile handling, and `MISSING_REFERENCE` schema consideration.

Non-goals
---------

- No stable Core 2.0 release
- No release candidate yet
- No full verify execution yet
- No HTTP API server
- No signing or key custody
- No authentication or tenancy
- No policy judgement
- No forced `MISSING_REFERENCE` schema addition

Recommended next after Draft 4
------------------------------

- Collect partner feedback on the 11-case conformance pack.
- Decide whether `MISSING_REFERENCE` belongs in Verify Result v2 or bundle-level validation.
- Prepare schema-change proposal only if semantics are stable.
- Plan full verify execution phase after runner comparison semantics remain stable.
- Prepare `v2.0.0-rc.1` readiness only after partner feedback and schema decisions.

Compatibility statement
-----------------------

Draft 4 readiness does not promote any draft feature to stable. Experimental CLI behavior remains under the experimental namespace, v1.0 compatibility remains the baseline, and future Core 2.0 stabilization still requires explicit review and release decisions.
