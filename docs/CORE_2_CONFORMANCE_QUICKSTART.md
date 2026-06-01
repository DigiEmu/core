Core 2.0 — Conformance Quickstart

Purpose
-------

This quickstart shows partners how to run the experimental Core 2.0 conformance CLI and how to author simple conformance cases compatible with the internal conformance runner MVP.

Prerequisites
-------------

...

Expected output
---------------

On a healthy checkout of the current draft conformance pack you should see a short summary like:

```
Conformance run summary: total=3 passed=3 failed=0
```

Conformance case structure
--------------------------

Each conformance case is a directory containing at minimum the following files:

- `input.json`
- `expected_verify_result.json`

The `expected_verify_result.json` file is a lightweight declaration used by the runner to assert the expected Verify Result for the case. The runner performs minimal structural validation of this file (fields required: `result`, `reason_code`, `verify_result_version`).

Meaning of PASS / FAIL / ERROR
-----------------------------

Verify Result values (PASS / FAIL / ERROR) describe the expected outcome of the verification logic for the case. These values are expectations about verification behavior — they are not the same as the conformance runner's notion of a passed case.

Semantic note: A negative conformance vector that declares `result: "FAIL"` or `result: "ERROR"` is a successful conformance case when the `expected_verify_result.json` is structurally valid and correctly declares that expected outcome. The conformance runner counts such cases as passed when the expectation declaration is valid.

How to add your own conformance case
-----------------------------------

1. Create a new directory under `testdata/core_2_conformance/` with a short, unique name.
2. Add `input.json` containing the fixture input to be verified.
3. Add `expected_verify_result.json` with the keys:
   - `result`: one of `PASS`, `FAIL`, or `ERROR`.
   - `reason_code`: a short reason code string from the draft reason-code registry.
   - `verify_result_version`: the version identifier for the declared verify result schema (e.g., `v2-draft`).
4. Run the experimental CLI against the directory to ensure the case is discovered and the expected result file validates structurally.

Current limitations
-------------------

- The CLI command is experimental.
- The command currently validates Core 2.0 draft conformance case structure.
- It is not yet a full certification suite.
- It does not perform signing, key custody, network API calls, or policy judgement.
- It does not change v1.0 CLI behavior.

Compatibility statement
-----------------------

This quickstart describes an experimental, partner-facing workflow for validating Core 2.0 draft conformance cases. It is documentation-only and does not modify any production code, v1.0 CLI behavior, or snapshot hashing semantics.

Partner integration notes
-------------------------

See `docs/CORE_2_PARTNER_INTEGRATION_NOTES.md` for partner-facing guidance
about integration testing, schemas, and recommended workflows.
# Core 2.0 — Conformance Quickstart

Purpose
-------

This quickstart shows partners how to run the experimental Core 2.0 conformance CLI and how to author simple conformance cases compatible with the internal conformance runner MVP.

Prerequisites
-------------

- A Go toolchain available on your PATH (Go 1.20+ recommended).
- A clone of the DigiEmu Core repository with the draft Core 2.0 artifacts present.

Run the test suite
------------------

Before running the experimental CLI, run the project tests to ensure local invariants hold:

```bash
go test ./...
```

Run the experimental conformance CLI
-----------------------------------

Use the repository-local experimental entrypoint to run the conformance pack shipped in `testdata`:

```bash
go run ./cmd/digiemu experimental conformance run testdata/core_2_conformance
```

Expected output
---------------

On a healthy checkout of the current draft conformance pack you should see a short summary like:

```
Conformance run summary: total=3 passed=3 failed=0
```

Conformance case structure
--------------------------

Each conformance case is a directory containing at minimum the following files:

- `input.json`
- `expected_verify_result.json`

The `expected_verify_result.json` file is a lightweight declaration used by the runner to assert the expected Verify Result for the case. The runner performs minimal structural validation of this file (fields required: `result`, `reason_code`, `verify_result_version`).

Meaning of PASS / FAIL / ERROR
-----------------------------

Verify Result values (PASS / FAIL / ERROR) describe the expected outcome of the verification logic for the case. These values are expectations about verification behavior — they are not the same as the conformance runner's notion of a passed case.

Semantic note: A negative conformance vector that declares `result: "FAIL"` or `result: "ERROR"` is a successful conformance case when the `expected_verify_result.json` is structurally valid and correctly declares that expected outcome. The conformance runner counts such cases as passed when the expectation declaration is valid.

How to add your own conformance case
-----------------------------------

1. Create a new directory under `testdata/core_2_conformance/` with a short, unique name.
2. Add `input.json` containing the fixture input to be verified.
3. Add `expected_verify_result.json` with the keys:
   - `result`: one of `PASS`, `FAIL`, or `ERROR`.
   - `reason_code`: a short reason code string from the draft reason-code registry.
   - `verify_result_version`: the version identifier for the declared verify result schema (e.g., `v2-draft`).
4. Run the experimental CLI against the directory to ensure the case is discovered and the expected result file validates structurally.

Current limitations
-------------------

- The CLI command is experimental.
- The command currently validates Core 2.0 draft conformance case structure.
- It is not yet a full certification suite.
- It does not perform signing, key custody, network API calls, or policy judgement.
- It does not change v1.0 CLI behavior.

Compatibility statement
-----------------------

This quickstart describes an experimental, partner-facing workflow for validating Core 2.0 draft conformance cases. It is documentation-only and does not modify any production code, v1.0 CLI behavior, or snapshot hashing semantics.
