# Releasing (DigiEmu Core)

## Policy
- Releases are produced **only** from annotated tags on main.
- Tag format: MAJOR.MINOR.PATCH (SemVer).
- A release must have green CI on main for the tagged commit.

## Checklist
1. Ensure main is clean and CI is green.
2. Create annotated tag:
   - git tag -a vX.Y.Z -m "vX.Y.Z"
3. Push tag:
   - git push origin vX.Y.Z
4. Confirm the Release workflow published artifacts:
   - gh run list --workflow "Release (binaries)" --branch main

## Notes
- No releases from feature branches.
- Hotfix = PR -> green CI -> tag on main.