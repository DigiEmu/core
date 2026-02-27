# Snapshot v1 Demo (Golden)

This folder is a **reproducible public demo** for DigiEmu Core.

## Goal

A third party can:

1. Install the CLI
2. Run `verify` on the provided inputs
3. Reproduce the **expected snapshot hash**
4. Reproduce the **expected verify report** (deterministic JSON)

## Install

Use the locked CLI contract tag:

```bash
go install github.com/BrunoBaumgartner78/digiemu-core@cli-contract-v1.0.0
```

## Run (Verify)

From repository root:

```bash
digiemu verify \
  --snapshot "$(cat examples/snapshot_v1_demo/expected_snapshot_hash.txt)" \
  --inputs examples/snapshot_v1_demo/inputs \
  --out /tmp/report.json
```

## Expected Results

- Snapshot hash must equal `expected_snapshot_hash.txt`
- Verify report must be deterministic and match `expected_verify_report.json`

## Notes

- This demo is intentionally minimal.
- CI hardening and schema validation will be added in the next phase.
