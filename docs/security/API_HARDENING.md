
---

## `API_HARDENING.md`

```md
# API_HARDENING.md

## Purpose

This document defines API hardening expectations for DigiEmu Core when HTTP/API surfaces are exposed.

DigiEmu Core is primarily a deterministic integrity engine.  
If deployed behind an API, the API layer must not weaken determinism, auditability, or controlled failure behavior.

---

## Security Goals

The API layer should:

- reject malformed input safely
- preserve deterministic semantics
- avoid leaking unnecessary internal details
- enforce size and structure limits
- make audit-required operations fail closed when dependencies are missing
- produce stable, reviewable behavior under error conditions

---

## Non-Goals

This document does not by itself define:

- authentication product strategy
- SSO or identity federation
- reverse proxy configuration
- WAF strategy
- network segmentation
- secrets vaulting
- cloud IAM policy

Those belong to deployment architecture outside Core.

---

## API Design Principles

### 1. Explicit Input Contracts
All request bodies must have clear structure and size expectations.

### 2. Fail Closed
If audit logging is required for an operation and audit is unavailable, the operation must fail.

### 3. No Silent Mutation
The API must not silently normalize meaningful state in a way that changes trust semantics without documentation.

### 4. Stable Error Semantics
Errors should be machine-usable and human-reviewable.

### 5. No Hidden Determinism Drift
API wrappers must not reorder or mutate deterministic content unexpectedly.

---

## Minimum Hardening Controls

## 1. Request Size Limits

Apply explicit maximum payload sizes for:
- meaning documents
- claim sets
- uncertainty documents
- bundle uploads if supported

### Why
Prevents accidental or malicious oversized input from degrading service or masking logic errors.

---

## 2. Content-Type Validation

Accept only documented content types for JSON endpoints.

Example:
- `application/json`

Reject ambiguous or malformed content types where appropriate.

---

## 3. JSON Decode Discipline

- reject invalid JSON
- reject truncated JSON
- reject structurally incompatible payloads
- do not continue on partial parse success

Where practical:
- disallow unknown fields for strict endpoints
- keep parsing behavior explicit and documented

---

## 4. Error Output Discipline

Error responses should:
- identify the failure category
- avoid dumping internal stack traces
- avoid exposing local filesystem details unless intentionally part of a trusted debug flow
- avoid ambiguous “success-like” responses on failure paths

---

## 5. Audit-Critical Mutation Rules

For operations such as:
- create unit
- create version
- set meaning
- set claims
- set uncertainty

the API must preserve the same mutation guarantees as the underlying usecase.

If strict audit is required:
- no successful response without corresponding audit success

---

## 6. Deterministic Response Considerations

For responses used in verification or evidence contexts:
- keep field semantics stable
- avoid environment-dependent formatting
- avoid implicit omission of integrity-relevant fields
- document contract changes before release

---

## 7. Path and File Safety

If any API endpoint resolves files or bundle paths:

- do not allow path traversal
- validate bundle roots explicitly
- reject unexpected directory shapes
- avoid trusting caller-supplied relative paths blindly

---

## 8. Method Discipline

Use HTTP methods consistently:

- `GET` for read-only retrieval
- `POST` for mutations or verification requests with bodies
- avoid unsafe state changes through accidental `GET`

---

## 9. Timeout and Resource Boundaries

Recommended deployment controls:
- request timeout
- body size limit
- connection/read timeout
- bounded concurrency where needed

These may be enforced by application code, reverse proxy, or both.

---

## 10. Logging Hygiene

Operational logs should not leak:
- secrets
- raw sensitive payloads unless explicitly intended
- misleading success/failure status

Where request bodies are logged for debugging, that behavior must be explicit and controlled.

---

## Recommended Error Shape

Example structure:

```json
{
  "error": {
    "code": "invalid_request",
    "message": "unsupported schema_version"
  }
}