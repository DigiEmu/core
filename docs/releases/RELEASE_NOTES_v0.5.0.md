DigiEmu Core v0.5.0 — Phase 3.1 Uncertainty Minimum

Summary
-------

This release introduces a minimal, auditable `uncertainty/v0` sidecar attached to versions. It is optional and additive.

Highlights
----------
- `UNCERTAINTY_SET` audit events when uncertainty is added.
- Version-scoped sidecars: `data/units/<unit-id>.<version-id>.uncertainty.json` with canonical JSON hashing (`uncertainty_hash`).
- `verify-audit --strict-hash` detects tampering or missing uncertainty sidecars.
- CLI: `digiemu uncertainty set` / `digiemu uncertainty show`.
- HTTP: `PUT /v1/units/{unit}/uncertainty` and `GET /v1/units/{unit}/uncertainty`.

Notes
-----

Uncertainty is an explicit metadata object for authors to record uncertainty about a version or a claim. It is not an inference engine and contains no probabilistic modeling.
