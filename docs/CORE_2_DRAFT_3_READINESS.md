# Core 2.0 Draft 3 Readiness Notes

Status
------

- Milestone: planned `v2.0.0-draft.3`
- Stability: draft / pre-release / not stable
- Purpose: Prepare a future Draft 3 tag without creating it yet.

Purpose
-------

This document summarizes readiness for a future `v2.0.0-draft.3` milestone. It is documentation-only and does not create or move tags, promote Core 2.0 draft features to stable, change v1.0 behavior, change schemas, change testdata, or modify CLI behavior.

Changes since `v2.0.0-draft.2`
-------------------------------

- Partner handoff package
- Partner feedback intake templates and process
- Draft 2 review summary
- Expanded conformance pack from 3 to 10 cases
- Reason code review
- Runner real execution design
- Runner `input.json` parsing
- Missing conformance case file test coverage

Why Draft 3 is justified soon
-----------------------------

- The conformance path has moved from a small representative pack toward a broader 10-case draft pack.
- Partner-facing handoff and feedback intake are documented.
- Reason-code semantics have stabilization notes for future stable promotion decisions.
- The runner now reads and parses `input.json`, while malformed input behavior is covered by tests.
- Missing input and missing expected-result behavior has focused test coverage.
- Docker and CI checks provide a more reproducible partner evaluation path than Draft 2 alone.

Current strengths
-----------------

- Docker-CI validated conformance path
- Machine-readable JSON reports
- 10-case conformance pack
- Reason-code stabilization notes
- Partner-facing handoff and feedback process
- Runner now reads/parses `input.json`
- Missing input/expected result behavior has test coverage

Before tagging Draft 3
----------------------

- Run `go test ./...`
- Run experimental conformance CLI human-readable
- Run experimental conformance CLI `--json`
- Confirm GitHub Actions are green
- Confirm README still reflects current draft status
- Confirm no stable `v2.0.0` promotion
- Confirm v1.0 behavior remains unchanged
- Confirm no unexpected files are modified

Proposed tag
------------

- Name: `v2.0.0-draft.3`
- Title: DigiEmu Core 2.0 Draft 3 — Runner-Hardened Conformance Milestone
- Meaning: A draft/pre-release milestone that reflects the expanded conformance pack, reason-code review, runner input parsing, missing-file test coverage, and partner feedback readiness.

Non-goals
---------

- No stable Core 2.0 release
- No release candidate yet
- No HTTP API server
- No signing/key custody
- No authentication or tenancy
- No policy judgement

Recommended next after Draft 3
------------------------------

- Runner actual-vs-expected comparison model, documented in `docs/CORE_2_ACTUAL_EXPECTED_COMPARISON_DESIGN.md`
- Additional representable `UNSUPPORTED_PROFILE` and `MISSING_REFERENCE` cases
- OpenAPI contract refinement
- Partner feedback triage
- `v2.0.0-rc.1` planning

Compatibility statement
-----------------------

Draft 3 readiness does not promote any draft feature to stable. Experimental CLI behavior remains under the experimental namespace, v1.0 compatibility remains the baseline, and future Core 2.0 stabilization still requires explicit review and release decisions.
