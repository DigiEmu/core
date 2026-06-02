# Core 2.0 Draft — Partner-Testable Milestone (Draft Release Notes)

Milestone name: Core 2.0 Draft — Partner-Testable Milestone

Status: draft

Summary
-------

- Core 2.0 is now partner-testable via the experimental conformance CLI.
- Existing v1.0 behavior remains unchanged and is the baseline for
  compatibility.
- Core 2.0 artifacts remain draft unless explicitly promoted to stable.

Included in this milestone
--------------------------

- Boundary Model
- Hardening Plan
- Crypto Agility Draft
- Migration from v1 Guidance
- Verify Result v2 Draft
- Reason Codes Registry
- Verify Result v2 JSON Schema
- Verify Result v2 Examples
- Schema Validation Test
- Conformance Pack
- Test Vector Rules
- Canonicalization Audit
- Canonicalization Behavior Tests
- Hash Boundary Vectors
- Canonicalization Decision Record
- Profile Registry
- Compatibility Matrix
- Experimental Conformance CLI
- Conformance Quickstart
- Partner Integration Notes

Tested
------

- `go test ./...` passes.
- Experimental conformance CLI runs against `testdata/core_2_conformance`.
- Expected output when running the CLI locally:

```
Conformance run summary: total=10 passed=10 failed=0
```

Not included in this milestone
-----------------------------

- Stable Core 2.0 CLI
- HTTP API / OpenAPI specification
- Docker image
- SDKs for non-Go users
- Secure Layer signatures (MVP planned separately)
- Post-Quantum implementation (planned separately)
- Full certification suite

Guidance for partners
---------------------

- Use the experimental CLI and the test vectors in `testdata/` for evaluation and
  feedback.
- Do not treat draft artifacts as production guarantees.
- Coordinate profile registration and canonicalization decisions with the
  DigiEmu team before relying on non-default behavior.

Release checklist and tagging
----------------------------

Before creating a draft release tag (for example `v2.0.0-draft.1`) follow the
checklist in `docs/CORE_2_RELEASE_CHECKLIST.md` and the tagging guidance in
`docs/CORE_2_TAGGING_PLAN.md` to ensure the milestone is reproducible and
partner-testable.
