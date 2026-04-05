# PRE_DEPLOY_CHECKLIST.md

## Purpose

This checklist defines the minimum release gate before deploying or tagging DigiEmu Core for external use, demos, partner review, or enterprise conversations.

---

## 1. Build and Test Gate

- [ ] `go build ./...` passes
- [ ] `go test ./...` passes
- [ ] critical reproducibility tests pass
- [ ] verification output contract tests pass
- [ ] no parser errors remain
- [ ] no unresolved compile errors remain

---

## 2. Guard Gate

- [ ] `go run ./cmd/digiemu-guard --core-only --ignore-file .digiemu-guard-ignore.json` passes
- [ ] every ignore entry in `.digiemu-guard-ignore.json` is intentional
- [ ] ignored findings are documented in release docs
- [ ] no newly introduced critical findings are ignored silently

---

## 3. Determinism Gate

- [ ] bundle loading order is deterministic
- [ ] map iteration does not affect hash, replay, or verification outputs
- [ ] canonical JSON behavior is unchanged or intentionally versioned
- [ ] snapshot hash output remains reproducible
- [ ] byte-exact verification fixture tests pass where required
- [ ] trace output remains stable enough for contract expectations

---

## 4. Serialization Gate

- [ ] all integrity-relevant fields are preserved in canonical or verification state
- [ ] `omitempty` use is reviewed
- [ ] any remaining `omitempty` usage is explicitly justified
- [ ] result JSON contract is stable
- [ ] audit payload fields are clear and consistent

---

## 5. Audit Gate

- [ ] strict-audit flows still require audit backend where intended
- [ ] unit creation audit event is emitted
- [ ] version creation audit event is emitted
- [ ] meaning / claims / uncertainty mutations emit expected audit events
- [ ] audit verification helpers still behave deterministically
- [ ] audit-related hashes or summaries are documented if used

---

## 6. Documentation Gate

- [ ] `README.md` is current
- [ ] `THREAT_MODEL.md` exists and is current
- [ ] `SECURITY_POLICY.md` exists and is current
- [ ] dependency audit document exists
- [ ] API hardening document exists if API is shipped
- [ ] support policy exists
- [ ] versioning policy exists
- [ ] release notes mention residual risks and ignores

---

## 7. Release Artifact Gate

- [ ] example bundles still verify as expected
- [ ] example verification report fixtures are current
- [ ] expected snapshot hash fixtures are current
- [ ] no accidental BOM / encoding drift in release fixtures
- [ ] no stale generated artifacts remain in repo
- [ ] tagged docs point to correct spec / contract version

---

## 8. API / Runtime Gate

Only if API/runtime deployment is part of the release:

- [ ] API endpoints have documented request size limits
- [ ] malformed JSON handling is tested
- [ ] error responses do not leak unnecessary internals
- [ ] write paths follow explicit policy
- [ ] audit-required operations fail closed if audit backend is missing
- [ ] deployment configuration is documented

---

## 9. Enterprise Readiness Gate

- [ ] threat model reviewed
- [ ] security policy reviewed
- [ ] support expectations documented
- [ ] known exceptions documented
- [ ] no open release-blocking integrity issues
- [ ] no unresolved determinism regressions in release path
- [ ] no known contract drift left undocumented

---

## 10. Final Human Review

Before deploy, confirm:

- [ ] this release can be explained clearly to a technical reviewer
- [ ] this release can be defended in an audit conversation
- [ ] this release has no hidden “temporary” behavior that matters for trust
- [ ] every accepted risk is written down somewhere visible
- [ ] release message matches actual repository state

---

## Deploy Decision

### Deploy only if all are true:
- build passes
- tests pass
- guard passes
- docs are in place
- residual risks are explicit

### Do not deploy if any of these are true:
- failing reproducibility test
- hidden guard ignore
- undocumented determinism exception
- undocumented contract change
- unresolved integrity ambiguity

---

## Release Command Log

Recommended to record before tag or deploy:

```txt
go build ./...
go test ./...
go run ./cmd/digiemu-guard --core-only --ignore-file .digiemu-guard-ignore.json