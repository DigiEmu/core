# DigiEmu Core — Public API (pkg/*)

This repository exposes a small **stable public API surface** under `pkg/*`.
Everything under `internal/*` is **not** part of the public API and may change without notice.

## What is stable

The following packages are intended to remain backwards-compatible within the `v1.x` line:

- `pkg/meaning` — public meaning types (stable surface)
- `pkg/claims` — public claim types (stable surface)
- `pkg/uncertainty` — public uncertainty types (stable surface)
- `pkg/snapshot` — snapshot references + manifest shape (stable surface)
- `pkg/verify` — verifier interface + result shape (stable surface)

See: `docs/PUBLIC_API_v1.0.md` and `docs/COMPATIBILITY_POLICY_v1.0.md`.

## Import path

When using DigiEmu Core as a module:

- `digiemu-core/pkg/snapshot`
- `digiemu-core/pkg/verify`
- `digiemu-core/pkg/claims`
- `digiemu-core/pkg/meaning`
- `digiemu-core/pkg/uncertainty`

(The repository module path is defined in `go.mod`.)

## Minimal example: verification interface

```go
package main

import (
  "fmt"

  "digiemu-core/pkg/snapshot"
  "digiemu-core/pkg/verify"
)

type myVerifier struct{}

func (myVerifier) Verify(ref snapshot.Ref) (verify.Result, error) {
  return verify.Result{OK: true, Ref: ref}, nil
}

func main() {
  var v verify.Verifier = myVerifier{}
  ref := snapshot.Ref{Hash: snapshot.Hash("abc")}
  res, _ := v.Verify(ref)
  fmt.Println(res.OK)
}
```

## Compatibility notes (short)

- Patch releases may add new exported identifiers, but should not break existing ones.
- `v2` may introduce breaking changes.
- `internal/*` can change anytime.

## Contributing

If you propose changes to `pkg/*`, include:
- a compatibility note (why it does not break v1)
- tests for the new behavior

## Freeze marker

The v1.0 public API freeze marker is documented here:
- `docs/PUBLIC_API_FREEZE_v1.0.md`

## CLI note

A minimal CLI command `digiemu verify` is provided as a stub reference implementation for the `pkg/verify` contract.
See: `docs/CLI_VERIFY_v1.0.md`

## Genesis anchor
The normative genesis anchor for deterministic reconstruction is documented here:
- `docs/GENESIS_ANCHOR_v1.0.yaml`
