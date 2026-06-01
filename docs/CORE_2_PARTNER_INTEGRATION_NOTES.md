# Core 2.0 — Partner Integration Notes

Purpose
-------

This document provides concise guidance for external partners who want to test,
evaluate, and prepare integrations against DigiEmu Core 2.0 draft artifacts using
the experimental conformance tooling and existing schemas in this repository.

Who this document is for
------------------------

- Partner engineering teams exploring integration points with DigiEmu Core.
- Security and QA teams validating conformance vectors and canonicalization.
- Architects planning migration paths from v1.0 to Core 2.0 profiles.

Current integration surface
---------------------------

- JSON schemas (draft and v1 artifacts)
- Core 2.0 conformance testdata (`testdata/core_2_conformance/`)
- Experimental conformance CLI (`digiemu experimental conformance run ...`)
- Internal conformance runner (`internal/conformance`)
- Profile registry (draft schema under `schemas/`)
- Compatibility matrix (`docs/CORE_2_COMPATIBILITY_MATRIX.md`)

Recommended first test
----------------------

Run the repository-local experimental conformance CLI against the included
conformance pack to validate your environment and the case discovery logic:

```bash
go run ./cmd/digiemu experimental conformance run testdata/core_2_conformance
```

Expected output:

```
Conformance run summary: total=3 passed=3 failed=0
```

Using the experimental conformance CLI
-------------------------------------

The experimental command is intentionally non-stable and is provided to make
it easy for partners to validate conformance case structure locally. Invoke the
command from the repository root as shown above. The command exercises
discovery and lightweight structural validation of `expected_verify_result.json`.

Using schemas and profiles
--------------------------

- Schema-backed artifacts live in `schemas/` and include explicit versions in
  filenames. Use these schemas to validate machine-readable artifacts.
- Profiles and canonicalization decisions are draft and documented in the
  profile registry and canonicalization decision documents. Register a profile
  identifier before relying on it for canonicalization changes.

Expected partner workflow
-------------------------

1. Clone the repository and run `go test ./...` to verify local test expectations.
2. Run the experimental conformance CLI against the included testdata.
3. Author a small conformance case under `testdata/core_2_conformance/` and
   re-run the CLI to validate discovery and structural checks.
4. Use schema validators (JSON Schema tools) to verify `expected_verify_result.json`.
5. Coordinate with the DigiEmu team to register new profile identifiers when
   proposing changes to canonicalization or hashing behaviors.

What is stable, draft, and experimental
---------------------------------------

Stable:

- `v1.0` behavior including `Snapshot Hash v1` and `Canonical JSON v1`.

Draft:

- `Verify Result v2` (draft), `Core 2.0 Profile Registry`, and `Core 2.0 Conformance Pack`.

Experimental:

- The `experimental conformance CLI` under `digiemu experimental ...`.

What DigiEmu Core does not do
-----------------------------

- The experimental tooling does not perform signing, key custody, or network
  APIs. It does not perform policy judgement or act as a certification
  authority. Those features are intentionally out of scope for the draft runner.

Future integration paths
------------------------

Planned but not implemented in this draft:

- Promotion of stable CLI commands for conformance
- An HTTP API / OpenAPI specification for remote validation
- Docker images and SDKs for non-Go users
 - Docker images and SDKs for non-Go users

Docker usage (optional)
-----------------------

For partners who prefer not to install Go locally, the repository includes a
minimal Docker build and companion documentation (`docs/CORE_2_DOCKER_USAGE.md`)
that produce a small container with the `digiemu` CLI and the default conformance
testdata. This is an optional, partner-facing evaluation path and does not
change any runtime or CLI behavior.
- Secure layer signatures and post-quantum migration profiles

Compatibility statement
-----------------------

These notes describe draft and experimental artifacts for partner evaluation.
They do not modify production runtime behavior or the v1.0 CLI contract. Use
draft artifacts for testing and do not rely on them for production security
guarantees.
