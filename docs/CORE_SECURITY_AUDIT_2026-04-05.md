# DigiEmu Core security/core audit

## Executive summary

The uploaded repository is promising and already shows strong architectural intent around deterministic replay and verification. But it is **not yet ready** to be sold as a fully hardened enterprise product without further work.

The highest-value issues found in the uploaded snapshot were:

1. **Public verify contract file was corrupted**: `pkg/verify/verify.go` contains bundle-loading code instead of the public `Result` contract. That directly causes compile and embedding failures in `internal/verify`.
2. **Deterministic ID migration is incomplete**: `NewID(prefix, input)` exists, but multiple callers still use the old one-argument form. That breaks builds and creates inconsistent ID semantics across the core.
3. **Serialization semantics are still too permissive**: many core structs still use `omitempty`, which allows semantically relevant fields to silently disappear from JSON.
4. **Boundary time usage remains**: HTTP handlers and some adapters still use `time.Now()`. That is acceptable only at the boundary, not in hashed or replayed core logic.
5. **Release hardening is incomplete**: there is no evidence in the uploaded snapshot of reproducible build enforcement, signed release attestations, SBOM generation, dependency audit gating, or artifact provenance verification.

## What I would call the current state

- **Core concept**: strong
- **Deterministic verification layer**: good foundation
- **Operational hardening**: incomplete
- **Enterprise sale readiness**: not yet

## Immediate changes included in this drop-in pack

- restore `pkg/verify.Result`
- make `internal/verify.BundleV1` and `StateV1` explicit instead of omitting empty slices
- fix deterministic ID call sites in unit/version and major write usecases
- keep audit event IDs deterministic and derived from stable state
- fix `cmd/digiemu/verify.go` import set and keep it aligned with `internal/verify.ResultV1`

## Next hardening steps after these drop-ins

### Tier 1 — before selling to any serious customer
- remove `omitempty` from all snapshot / verify / audit-critical domain structs
- define explicit canonical serialization policy for every persisted/hashable struct
- add reproducible build workflow and release provenance attestations
- freeze dependency versions and add vulnerability scan gates
- add strict CI for `go build ./...`, `go test ./...`, guard, schema locks, and golden replay
- add security policy, disclosure process, key rotation policy, and signed release process

### Tier 2 — before regulated or high-assurance deployment
- SBOM generation and artifact signing
- reproducible build verification on at least two environments
- independent code review / pentest
- structured threat model for snapshot tampering, oracle poisoning, build compromise, and signing key misuse
- fuzzing for bundle loader and canonical JSON edge cases
- formalized compatibility/versioning policy for persisted state

### Tier 3 — true enterprise posture
- SLSA-style build provenance
- hardened release key management
- separate security review for identity/signature CLI path
- deterministic test matrix across Windows/Linux/macOS if platform parity matters commercially

## Honest conclusion

These drop-ins materially improve the repository and remove concrete stability/integration defects. They **do not** make DigiEmu Core automatically enterprise-certified or security-complete. They are the right next engineering step, not the end state.
