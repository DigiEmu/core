# Core 2.0 — Conformance Runner (MVP)

This document describes a small internal conformance runner implemented as an MVP to validate Core 2.0 draft conformance vectors in `testdata/core_2_conformance`.

Purpose
-------

- Provide an internal, programmatic way to verify that conformance cases are structurally coherent.
- Keep scope intentionally small: this is NOT a public CLI command and does NOT implement verification logic beyond minimal structural checks.

Usage
-----

- The runner is implemented in `internal/conformance` and exposes `RunAll(root)` which returns a list of per-case `Result` objects.
- It discovers case directories under the provided `root` (for repository tests this is `testdata/core_2_conformance`).

Experimental CLI
----------------

An experimental CLI entrypoint is available under the main `digiemu` binary:

```
digiemu experimental conformance run <path-to-conformance-dir>
```

This command is intentionally experimental and draft-only. It invokes the internal conformance runner and prints a brief summary of total/passed/failed cases. It is not a public or stable CLI yet.

Limitations
-----------

- The MVP performs minimal structural validation of `expected_verify_result.json` (fields: `result`, `reason_code`, `verify_result_version`).
- Full JSON Schema validation is performed elsewhere in the test suite; this runner focuses on case discovery and lightweight checks to enable CI-friendly conformance smoke tests.
- No CLI integration, network APIs, or signing/custody features are included.

Future work
-----------

- Reuse the full Verify Result v2 schema validation where convenient.
- Add an executable test runner CLI that invokes this internal API and emits machine-readable test reports.
