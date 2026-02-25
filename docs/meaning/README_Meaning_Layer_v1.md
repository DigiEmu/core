README — Meaning Layer v1

Overview
--------
Meaning Layer v1 provides an optional, audit-friendly, and deterministic way to attach structured context to Units and Versions.
It is deliberately minimal: no ontologies, no inference, just a stable JSON schema and canonicalization rules so meaning can be hashed and verified.

File & location
---------------
- Sidecar filename: `data/units/<unit-id>.<version-id>.meaning.json`
  - Each meaning document is version-scoped and stored next to the unit JSON as a version-specific sidecar.
  - Example on disk: `data/units/unit_abcd.ver_12345.meaning.json` next to `data/units/unit_abcd.json`.
- Optional: Units/Versions without a meaning sidecar are fully supported and unchanged.

Schema (summary)
-----------------
Root object, `schema_version` required.
- `schema_version` (string) - e.g. "meaning/v1" (required)
- `title` (string) - short human label (optional)
- `purpose` (string) - why this Unit exists (optional)
- `scope` (object) - audience / jurisdiction / locale / timeframe (optional)
- `claims` (array) - small claim objects {text, strength, tags} (optional)
- `sources` (array) - source references {id, type, ref, quote?} (optional)
- `provenance` (object) - author/org/role/created_at/updated_at (optional)
- `integrity` (object) - narrative_id, supersedes, conflicts_with (optional)

Canonicalization & hashing
---------------------------
- Canonicalization: `json_c14n_minimal`
  - sort object keys
  - preserve array order
  - normalize newlines to `\n`
  - no extra whitespace
- Hashing: SHA-256 over canonicalized bytes
- Result field: `meaning_hash` (hex lowercase)
- Rules: meaning is canonicalized before inclusion in SnapshotHash and AuditHash; meaning_hash can be stored in audit events and manifests.

Audit integration
-----------------
- New audit event type: `MEANING_SET`
  - payload: `unit_id`, `version_id`, `meaning_hash`, `meaning_path`, `inline_preview` (optional: {title,purpose})
- When `meaning.json` is set via CLI or API, append `MEANING_SET` event (strict-audit semantics apply)
- Verify-audit will check: if meaning exists in snapshot, there must be a corresponding `MEANING_SET` event and its `meaning_hash` must match the file contents.

Export manifest extension
-------------------------
- Export manifest (snapshot) includes optional `meaning` object per unit/version:
  {
    "meaning": {
      "meaning_hash": "...",
      "meaning_path": "meaning.json",
      "schema_version": "meaning/v1"
    }
  }
- SnapshotHash includes the canonicalized meaning content when present; AuditHash includes MEANING_SET events like other events.

CLI
---
- `digiemu meaning set <unitKeyOrId> [--version <versionId>] --file path/to/meaning.json [--data ./data]`
  - Reads the file bytes and calls the `SetMeaning` usecase.
  - Validation, canonicalization and `meaning_hash` computation happen inside the usecase.
  - On success a `MEANING_SET` audit event is appended and the sidecar is written at `data/units/<unit-id>.<version-id>.meaning.json`.

- `digiemu meaning show <unitKeyOrId> [--version <versionId>] [--data ./data]`
  - Resolves the version (explicit or the unit head) and prints the canonicalized meaning JSON and the `meaning_hash` recorded on the version.

- `digiemu verify-audit --data ./data [--strict-hash]` will detect tampering of meaning sidecars when `--strict-hash` is used. If the sidecar contents do not canonicalize to the `meaning_hash` recorded on the version or the `MEANING_SET` event, `verify-audit` reports a hash mismatch.

HTTP API
--------
- PUT /v1/units/{unitKey}/meaning?version=<versionId>
  - Body: `meaning.json` payload
  - Response: `201 Created` with JSON `{ "unit_id": "...", "version_id": "...", "meaning_hash": "..." }`

- GET /v1/units/{unitKey}/meaning?version=<versionId>
  - Response: `200 OK` with JSON `{ "meaning": { ...canonicalized... }, "meaning_hash": "..." }`

Example curl (set meaning):

```bash
curl -X PUT \
  -H "Content-Type: application/json" \
  --data-binary @meaning.json \
  "http://localhost:8080/v1/units/merkblatt-steuern/meaning?version=ver_138415ce411a5990ea5e02bd3922e1c1"
```

Example curl (get meaning):

```bash
curl "http://localhost:8080/v1/units/merkblatt-steuern/meaning?version=ver_138415ce411a5990ea5e02bd3922e1c1"
```

Stability & compatibility
-------------------------
- Meaning is optional — existing snapshots and workflows are unchanged if no meaning is provided.
- No breaking changes to kernel ports or storage layout (meaning stored alongside existing files).

Limits & security
-----------------
- `meaning.json` max 64 KB
- No personal/sensitive data should be included in meaning (policy enforced externally)

Tests (developer guidance)
--------------------------
- Unit: canonicalization stability (shuffled keys -> same canonical + same hash)
- Unit: meaning optional (no change to SnapshotHash when absent)
- Unit: meaning alters SnapshotHash & AuditHash when present
- Integration: CLI set -> export -> verify-audit
- Integration: HTTP PUT/GET roundtrip

Next steps
----------
1. Implement domain types and canonical hashing helpers
2. Add ports and FS adapter support (store/load meaning.json)
3. Wire MEANING_SET audit events and verify-audit checks
4. Add CLI commands and HTTP endpoints

