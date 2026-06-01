# Core 2.0 OpenAPI Contract Draft

Status: Draft — documentation-only

This document describes a draft OpenAPI (3.x) contract for an optional
HTTP/JSON surface that partners may use to run or validate Core 2.0
conformance artifacts. This is a contract draft only: there is no network
server implemented in this repository and no production behavior is changed
by including this draft.

Principles

- This is a contract draft only; it does not implement any server.
- The API MUST NOT change or replace the v1.0 CLI contract or runtime
  semantics of Core.
- The API surface is intended to wrap existing Core 2.0 artifacts (verify
  result v2 documents, conformance report format, and profile registry), not
  to redefine them.
- Authentication, tenancy, signing, key custody, or policy judgement are out
  of scope for Core; any production deployment should consider those as
  integration-layer responsibilities.

Schemas referenced

- `schemas/verify_result_v2.schema.json`
- `schemas/core_2_conformance_report.schema.json`
- `schemas/core_2_profile_registry.schema.json`

Intended endpoints (draft)

- `GET /v1/core2/version` — return draft API/version information.
- `GET /v1/core2/profiles` — list known profile identifiers from the registry.
- `POST /v1/core2/conformance/run` — accept a conformance case set or
  reference and return a machine-readable conformance report.
- `POST /v1/core2/verify-result/validate` — validate a Verify Result v2
  document structurally and return a validation response.
- `POST /v1/core2/conformance-report/validate` — validate a Core 2.0
  conformance report document structurally and return a validation
  response.

Usage

This draft is intended for discussion and to inform partner integration
designs. It is not a recommendation to ship an HTTP endpoint inside Core; a
separate integration service is expected to host such an API if desired.

Feedback

Open issues and feedback should be recorded as PR comments against the
`core-2-openapi-contract-draft` branch. When ready, maintainers may choose to
promote, adapt or split the contract into an integration repository.

Validation

This draft is structurally validated by a repository test that parses the
`openapi/core_2_conformance_api.yaml` file and asserts presence of the
`openapi` version, `info.title`, required paths and `components.schemas`.
The test is documentation-only and enforces basic structural correctness
but does NOT imply a running API server.
