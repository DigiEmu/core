# Releasing (DigiEmu Core)

## Principle

A DigiEmu Core release is a **verifiable system state**, not just a code snapshot.

A valid release guarantees that:
- deterministic behavior is stable
- verification outputs are reproducible
- evidence semantics are preserved

## Policy
- Releases are produced **only** from annotated tags on main.
- Tag format: MAJOR.MINOR.PATCH (SemVer).
- A release must have green CI on main for the tagged commit.

## Checklist
1. Ensure main is clean, CI is green, and deterministic verification is stable.
2. Create annotated tag:
   - git tag -a vX.Y.Z -m "vX.Y.Z"
3. Push tag:
   - git push origin vX.Y.Z
4. Confirm the Release workflow published artifacts:
   - gh run list --workflow "Release (binaries)" --branch main

## Notes
- No releases from feature branches.
- Hotfix = PR -> green CI -> tag on main.

## Final Rule

A release is valid only if identical inputs produce identical verified outputs.