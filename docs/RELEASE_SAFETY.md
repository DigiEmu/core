# Release Safety — Go Module and Tagging Strategy

Status: draft / hardening phase

---

## Current Module Path

- Detected module path from `go.mod`: `digiemu-core`

---

## Go Semantic Import Versioning Risk

- Go modules treat major versions >= v2 specially: the module path MUST end with `/v2` (or `/vN`) to safely publish a `vN.0.0` tag.
- If you tag `v2.0.0` without changing the module path to end with `/v2`, downstream `go get` may break or resolve incorrectly.

---

## Unsafe Tags

- Using real `v2.0.0` (or `v2.x.y`) tags with a module path that does not end with `/v2` is unsafe.

---

## Safe Milestone Names

Use pre-release milestone tags and names that do not trigger Go's major-version handling:

- core-2.0-draft4
- core-2.0-rc1
- core-2.0-spec-stable
- core-2-draft4-self-audit

These are safe as Git tags or GitHub Releases notes, because they are not semantic import versions to the Go toolchain.

---

## When Real v2.0.0 Is Allowed

A real `v2.0.0` tag is only safe if the module path is changed to `github.com/DigiEmu/core/v2` and all imports and installation instructions are updated accordingly.

- Update `module` directive in `go.mod` to `github.com/DigiEmu/core/v2`.
- Update all internal imports to the new path.
- Update README install instructions and examples.
- Verify `go get github.com/DigiEmu/core/v2@v2.0.0` works in a clean module cache.

---

## Recommended Current Release Strategy

- Do NOT publish real `v2.0.0` tags while the module path does not end with `/v2`.
- Use safe milestone names listed above for partner-testable drafts and RCs.
- Keep Core 2.0 in draft/RC tags until documentation and conformance are stable.
- When ready to cut `v2.0.0`, perform a dedicated PR that updates the module path, internal imports, CI, and docs, then tag.

---

## Warning

Do not use real Git tags starting with `v2.0.0` unless the Go module path intentionally ends in `/v2` and all internal imports are updated accordingly.
