# DigiEmu Core hardening drop-ins

This package contains targeted replacement files for the uploaded DigiEmu Core snapshot.

What these files do:
- restore `pkg/verify/verify.go` to the correct public result contract
- remove compile breaks from the deterministic `NewID(prefix, input)` change
- make verify/state bundle shapes more explicit for canonical serialization
- reduce silent-state-disappearance risk in the verify path

What these files do **not** guarantee:
- they do not make the codebase "wasserdicht" or fully enterprise-certified
- they do not replace threat modeling, independent pentesting, secure release signing, reproducible builds, SBOM, vuln scanning, or legal/compliance review
- they do not resolve every remaining `omitempty` warning in domain models

Recommended apply order:
1. replace the files in this package into the repo
2. run `go build ./...`
3. run `go run ./cmd/digiemu-guard --core-only --ignore-file .digiemu-guard-ignore.json`
4. review remaining warnings file-by-file
5. only then widen the hardening scope to adapters, HTTP boundary, signing, and release pipeline
