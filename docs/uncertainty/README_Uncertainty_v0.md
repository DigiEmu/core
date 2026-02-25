# Uncertainty Minimum v0

This document describes the `uncertainty/v0` schema, storage, audit events and CLI/HTTP usage for Phase 3.1 Uncertainty Minimum v0.

Schema
------

- `schema_version`: `uncertainty/v0`
- `id`: string (required)
- `type`: one of `empirical`, `interpretative`, `incomplete` (required)
- `level`: one of `low`, `medium`, `high` (required)
- `text`: optional descriptive text
- `tags`: optional list of strings
- `applies_to`: object describing target
  - `scope`: `version` or `claim` (required)
  - `claim_id`: required when `scope == claim`

Sidecar storage
----------------

Uncertainty is stored as a version-scoped sidecar file when set:

```
data/units/<unit-id>.<version-id>.uncertainty.json
```

The `uncertainty_hash` (SHA-256 hex over canonicalized JSON) is stored in the version record inside the unit file.

Audit event
-----------

When uncertainty is set the system appends an audit event `UNCERTAINTY_SET` with payload `UncertaintySetData`:

- `unit_id`
- `version_id`
- `uncertainty_hash`
- `uncertainty_path`

Verify-audit (tamper detection)
------------------------------

`verify-audit --strict-hash` will for each version that has an `uncertainty_hash`:

- Require exactly one `UNCERTAINTY_SET` event for that version.
- Verify the event `uncertainty_hash` matches the recorded `version.uncertainty_hash`.
- Load the sidecar via repository `LoadUncertainty()` and recompute the canonical JSON hash. If the recomputed hash differs, the verifier reports a hash mismatch (tamper detected). Missing sidecar is also treated as a mismatch.

CLI examples
------------

Set uncertainty (reads JSON file and calls kernel usecase):

```
digiemu uncertainty set <unitKeyOrId> [--version <versionId>] --file uncertainty.json [--data ./data]
```

Show uncertainty (pretty-print JSON and prints `uncertainty_hash`):

```
digiemu uncertainty show <unitKeyOrId> [--version <versionId>] [--data ./data]
```

If no uncertainty exists for the version, `show` prints:

```
no uncertainty for this version
```

and exits with code `2`.

HTTP examples (curl)
--------------------

Set uncertainty via HTTP:

```
curl -X PUT "http://localhost:8080/v1/units/<unitKey>/uncertainty?version=<verId>" \
  -H 'Content-Type: application/json' \
  --data-binary @uncertainty.json
```

Response:

```
{ "unit_id": "unit_...", "version_id": "ver_...", "uncertainty_hash": "..." }
```

Get uncertainty via HTTP:

```
curl -X GET "http://localhost:8080/v1/units/<unitKey>/uncertainty?version=<verId>"
```

Response (200):

```
{ "uncertainty": { ... }, "uncertainty_hash": "..." }
```

Response (404): when no uncertainty sidecar exists for the version.

Notes
-----

- Uncertainty is optional — absence is a valid state and does not affect existing flows.
- The kernel enforces schema and minimal invariants on set.
- The canonicalization rules are the same as other sidecars (sorted keys, preserve arrays, minified) and a SHA-256 hex digest is used for the recorded hash.