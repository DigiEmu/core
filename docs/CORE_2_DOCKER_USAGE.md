# Core 2.0 — Docker Usage (Partner-facing)

This guide describes a minimal Docker-based usage path for running the
experimental Core 2.0 conformance CLI without installing Go locally. The
container built by the provided `Dockerfile` is intended for partner
evaluation only.

Important constraints

- The image and docs are documentation-only and do not change CLI or v1.0
  behavior.
- The image is intentionally minimal and contains the `digiemu` CLI binary
  compiled from the repository, plus the default `testdata/core_2_conformance`
  pack for quick evaluation.
- The image does not implement or expose any network API.

Build

From the repository root:

```bash
docker build -t digiemu-core .
```

Run (human-readable summary)

```bash
docker run --rm digiemu-core experimental conformance run /opt/testdata/core_2_conformance
```

Run (JSON output)

```bash
docker run --rm digiemu-core experimental conformance run /opt/testdata/core_2_conformance --json
```

Expected sample outputs

- Human summary: `Conformance run summary: total=11 passed=11 failed=0`
- JSON report: top-level `status: PASS`, `total: 11`, `passed: 11`, `failed: 0`

Notes

- The Docker build requires network access at build time to download Go
  module dependencies but the runtime image is static and does not require
  network access.
- For faster iterative development, prefer building locally with `go build`
  and running the binary directly — the Docker path is primarily for
  partner evaluation where Go is not available.

CI validation
-------------

The repository's GitHub Actions CI includes a `docker-conformance` job that
builds the `digiemu-core` image and runs both the human-readable and `--json`
conformance checks inside the container to validate the Docker usage path.
This ensures the Docker evaluation path remains reproducible for partners.
