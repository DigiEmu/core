# Release v0.3.0 — Meaning Layer v1

Highlights
----------
- Meaning Layer v1: optional, version-scoped structured context attachments. Each meaning document is stored as a sidecar file next to unit JSON and is canonicalized and hashed for auditability.

What’s included
----------------
- New domain and usecase: `SetMeaning` — validates, canonicalizes, computes `meaning_hash`, persists the sidecar, updates the version record, and appends a `MEANING_SET` audit event.
- Storage: sidecar per version: `data/units/<unit-id>.<version-id>.meaning.json`. The version record contains an optional `meaning_hash`.
- CLI: `digiemu meaning set` and `digiemu meaning show` (thin adapters to the usecase and repository).
- HTTP API: `PUT /v1/units/{unit}/meaning` and `GET /v1/units/{unit}/meaning`.
- Verify-audit: enhanced to validate MEANING_SET events and to detect tampering when run with `--strict-hash`.

Upgrade notes
-------------
- Meaning is optional and non-breaking. Existing snapshots and workflows are unchanged when meaning sidecars are absent.
- To publish a meaning, run the CLI or call the HTTP PUT endpoint. Consider running `digiemu verify-audit --strict-hash` after setting meaning to ensure the audit journal contains a `MEANING_SET` event and the sidecar matches the recorded hash.

Security & Limits
------------------
- Maximum sidecar size: 64 KB.
- Do not store secrets or sensitive personal data in meaning documents.

Files changed (high-level)
--------------------------
- `internal/kernel/domain` — meaning types, version `MeaningHash` field, audit payloads
- `internal/kernel/usecases` — `SetMeaning` usecase, verify-audit enhancements
- `internal/kernel/adapters/fs` — version sidecar write/load, unit JSON updates
- `cmd/digiemu` — `meaning set` / `meaning show` CLI commands
- `internal/httpapi` — PUT/GET meaning endpoints
- `docs/meaning/README_Meaning_Layer_v1.md` — updated with CLI + HTTP examples
- `CHANGELOG.md` and this `RELEASE_NOTES_v0.3.0.md`

Suggested post-release checks
-----------------------------
1. Run `go test ./...` (already green locally).
2. Run `digiemu meaning set <unit> --file meaning.json` and `digiemu verify-audit --strict-hash` to confirm full flow.
3. Verify exported snapshots include meaning information if requested.

