## v0.1.1

- CLI/Core: support --desc/--description end-to-end (domain/ports/usecase/fs/http)
- Docs: improved quickstart + example responses
- Release: create/patch release by tag (avoid untagged drafts), script supports -Tag
- CI: add go vet (optional race via RUN_RACE env)
# Changelog

## v0.1.0
- Kernel: Units & Versions with validation
- Ports & DTO contracts
- In-memory + FS JSON repository
- HTTP API with consistent error format
- CLI (cmd/digiemu) + HTTP API cmd (cmd/api)
- CI workflow (gofmt check + go test)
 
## v0.2.9 — Audit & Snapshot Stabilization Release

- Deterministic Snapshot Hashing: `snapshotHash` included in export; `auditHash` when `--audit`.
- Hardened audit + snapshot pipeline across kernel, FS and memory adapters.
- Added export CLI command and verified audit verification tooling.

## v0.3.0 — Meaning Layer v1

- Added Meaning Layer v1: optional, version-scoped structured context attached to Units/Versions.
- Persistence: version-scoped sidecar files: `data/units/<unit-id>.<version-id>.meaning.json` and the `meaning_hash` stored on the version record.
- New audit event: `MEANING_SET` appended when a meaning is set. Payload includes `meaning_hash` and `meaning_path`.
- CLI: `digiemu meaning set` and `digiemu meaning show` thin adapters to the `SetMeaning` usecase.
- HTTP: `PUT /v1/units/{unit}/meaning` and `GET /v1/units/{unit}/meaning` endpoints added.
- Verify-audit: `verify-audit --strict-hash` detects tampering of meaning sidecars by comparing canonicalized sidecar content to recorded `meaning_hash`.

## v0.5.0 — Uncertainty Minimum v0

- Added Uncertainty Minimum v0: optional, version-scoped uncertainty metadata attached to versions.
- Persistence: version-scoped sidecar files: `data/units/<unit-id>.<version-id>.uncertainty.json` and the `uncertainty_hash` stored on the version record.
- New audit event: `UNCERTAINTY_SET` appended when an uncertainty is set. Payload includes `uncertainty_hash` and `uncertainty_path`.
- Verify-audit: `verify-audit --strict-hash` now detects tampering of uncertainty sidecars by recomputing canonicalized uncertainty JSON and comparing to recorded `uncertainty_hash`.
- CLI: `digiemu uncertainty set` and `digiemu uncertainty show` thin adapters to the `SetUncertainty` usecase. `show` prints `no uncertainty for this version` and exits `2` when missing.
- HTTP: `PUT /v1/units/{unit}/uncertainty` and `GET /v1/units/{unit}/uncertainty` endpoints added.
- No breaking changes: uncertainty is fully optional and additive.


See RELEASE_NOTES_v0.2.9.md for full details.


## Unreleased

### Public API
- Freeze v1.0 public API surface for `pkg/*` (meaning, claims, uncertainty, snapshot, verify).
- Added `docs/PUBLIC_API_FREEZE_v1.0.md` as the formal freeze marker.

