# P0 Conformance Matrix v0.1

- **Architecture baseline:** `architecture-baseline.yaml` revision 0.3
- **ADR:** `01_decisions/ADR-P0-05-Core-Admission-Constitution.md`
- **Manifest:** `docs/CORE_ARCHITECTURE_INTEGRATION_MANIFEST_v0.1.md`
- **Date:** 2026-08-07
- **Reviewed at HEAD:** `4780d15 Add Phase D Uncertainty Set admission binding proof`
- **Review date:** 2026-08-17
- **Status:** Proposed

## Methodology

Each integration invariant (IR-01 through IR-12) is evaluated against the actual repository evidence. A status is assigned only when the corresponding evidence exists; documentation alone is insufficient for `PASS`.

`PASS` in this matrix means satisfied for the current defined P0 architecture/conformance scope. It does not imply production Admission enforcement unless the invariant itself requires production binding.

|| Status | Meaning |
||--------|---------|
|| `PASS` | Sufficient specification, implementation, and/or executable evidence exists. |
|| `PARTIAL` | Evidence exists but is incomplete, scoped to RS-001 only, or has known limitations. |
|| `OPEN` | No implementation or executable evidence; only documentation/specification. |
|| `BLOCKED` | Progress is blocked by an external or environmental issue. |

## Conformance Matrix

### IR-01 — No Runtime Regression

- **ID:** IR-01
- **Title:** No Runtime Regression
- **Normative statement:** Architecture integration MUST NOT silently alter existing Core 2.0 behavior or stable contracts.
- **Current status:** PARTIAL
- **Evidence/artifacts:**
  - `architecture-baseline.yaml` lists the Core 2.0 boundary as authoritative and unchanged.
  - ADR-P0-05 and the manifest explicitly state that canonicalization, hashing, replay, verification, CLI v1, and state identity remain unchanged.
  - `git diff --name-status e9876f2..HEAD` shows no existing production file in `cmd/`, `internal/kernel/`, `internal/verify/`, `internal/conformance/`, or `schemas/` (existing schemas) was modified by P0 work; only new P0 files were added.
  - Targeted Go test suites pass:
    - `go test ./internal/admission/... -v` — PASS
    - `go test ./internal/kernel/usecases/... -v` — PASS
    - `go test ./internal/kernel/... ./internal/conformance ./internal/verify ./internal/canonicaljson/... -timeout 120s` — PASS
  - Phase E increments 1–5 (`de9ec49`, `2e38388`, `67a6a7e`, `e3383bd`, `eefe604`) modify only `internal/kernel/usecases/orchestrate_test.go`; targeted `go test ./internal/admission/...` and `go test ./internal/kernel/usecases/...` continue to pass; no production Core runtime files were changed.
- **Executable evidence/tests:** Targeted `go test` of Admission Engine, kernel, conformance, verify, and canonicaljson packages passes. The Phase D binding tests execute real Core handlers and produce no regression in existing Core behavior.
- **Remaining gap:** Full `go test ./...` is not yet green because `go run ./cmd/digiemu` in several test packages fails with `digiemu.exe` access denied in the local Windows temp environment. This is a file-lock/build-cache issue, not a code regression, so it should not be falsely reported as fully green.
- **Next action:** Resolve the temporary `go-build*` directory / antivirus / build-cache file lock, then rerun full `go test ./...` and confirm the existing Core 2.0 reproduction tests are byte-exact.

### IR-02 — Admission Before Mutation

- **ID:** IR-02
- **Title:** Admission Before Mutation
- **Normative statement:** A mutation governed by the Architecture Constitution MUST pass Admission before execution.
- **Current status:** PARTIAL
- **Evidence/artifacts:**
  - `internal/admission/engine.go` evaluates Intents and returns ADMIT/REJECT with a resolved `transition_ref`.
  - The conceptual target flow in the manifest shows `Intent → Admission → Command → Existing Core Kernel Use Case`.
  - `architecture-baseline.yaml` explicitly sets `p0_architecture_constitution.enforced: false`.
- **Executable evidence/tests:**
  - All five Phase D binding tests demonstrate the `Intent → ADMIT → real Core handler → real state mutation` chain for `core.unit.create`, `core.version.create`, `core.meaning.set`, `core.claim.set`, and `core.uncertainty.set`.
  - Each binding test also demonstrates the `REJECT → handler not invoked → no state mutation` case.
  - Phase E test-only orchestration in `internal/kernel/usecases/orchestrate_test.go` demonstrates the representative `ADMIT → real Core handler` and `REJECT → no handler` chain for `core.version.create` and `core.unit.create`, including: `ADMIT +` success, `ADMIT +` domain error, `ADMIT +` audit failure, and `ADMIT +` confirmed partial mutation. In all cases the `admission.Result` remains `ADMIT`; the `REJECT` path blocks the execution closure.
- **Remaining gap:** No production or CLI/HTTP path intercepts a Core mutation to require Admission first. The `digiemu` CLI and `internal/httpapi/handlers.go` still call `CreateUnit`, `CreateVersion`, `SetMeaning`, `SetClaims`, and `SetUncertainty` directly without an Admission gate.
- **Next action:** Add a non-production admitted-command orchestrator or conformance runner that demonstrates the `Intent → Admission → Command → Use Case` sequence for the first Core mutation path, then later wire the gate into CLI/HTTP under a feature flag without modifying existing Core contracts.

### IR-03 — Explicit Capability

- **ID:** IR-03
- **Title:** Explicit Capability
- **Normative statement:** Every constitution-governed command MUST reference a registered capability.
- **Current status:** PASS
- **Evidence/artifacts:**
  - `core-capability-registry.yaml` v0.1 lists all current Core capabilities with `id` and `mutation` fields.
  - `internal/admission/registry_v0.1.go` compiles these capabilities into `V01Registry()`.
  - `internal/admission/registry_v0_1_parity_test.go` verifies that the compiled registry matches the YAML source.
  - `schemas/intent-envelope.schema.json` requires `capability_ref`.
- **Executable evidence/tests:**
  - The reusable Admission Engine enforces `P0.ADMISSION.CAPABILITY_EXISTS` and `P0.ADMISSION.CAPABILITY_MUTATES`.
  - All five Phase D binding tests reference a registered mutating capability (`core.unit.create`, `core.version.create`, `core.meaning.set`, `core.claim.set`, `core.uncertainty.set`) and succeed.
  - Negative tests prove unknown or non-mutating capabilities are rejected.
  - Phase E `REJECT` and `ADMIT` tests (`de9ec49`, `2e38388`, `67a6a7e`, `e3383bd`, `eefe604`) continue to use `core.version.create` and `core.unit.create` from the registered capability set.
- **Remaining gap:** None for the current defined P0 scope. New mutation paths must be registered before they are admissible.
- **Next action:** Maintain the registry and parity test as new capabilities are added.

### IR-04 — Explicit Ownership

- **ID:** IR-04
- **Title:** Explicit Ownership
- **Normative statement:** Every constitution-governed mutation MUST resolve to an authorized aggregate owner.
- **Current status:** PASS
- **Evidence/artifacts:**
  - `aggregate-ownership-registry.yaml` v0.1 maps the `unit` aggregate to all five mutating Core capabilities.
  - `internal/admission/registry_v0.1.go` includes the same ownership mapping.
  - `internal/admission/registry_v0_1_parity_test.go` verifies the mapping against the YAML source.
  - `schemas/intent-envelope.schema.json` requires `aggregate_ref`.
- **Executable evidence/tests:**
  - The reusable Admission Engine enforces `P0.ADMISSION.AGGREGATE_OWNS_CAPABILITY`.
  - All five Phase D binding tests use `aggregate_ref: unit` for the `unit` aggregate and are admitted.
  - `TestEngine_OwnershipMismatch` and the `invalid_ownership_mismatch` RS-001 case prove non-owning aggregates are rejected.
  - Phase E `ADMIT` tests for `core.version.create` and `core.unit.create` also use `aggregate_ref: unit`, consistent with the ownership registry.
- **Remaining gap:** None for the current defined P0 scope. `unit` is the only aggregate root in the current Core domain, and every constitution-governed mutating capability is explicitly owned by `unit`. Future aggregates must be registered before they can be targets.
- **Next action:** Expand `aggregate-ownership-registry.yaml` only as new aggregates are introduced to the architecture.

### IR-05 — Command/Transition Correspondence

- **ID:** IR-05
- **Title:** Command/Transition Correspondence
- **Normative statement:** Every admitted mutating command MUST resolve to a defined transition.
- **Current status:** PASS
- **Evidence/artifacts:**
  - `command-event-catalogue.yaml` v0.1 maps five mutating commands (`unit.create`, `version.create`, `meaning.set`, `claim.set`, `uncertainty.set`) to their capabilities, aggregate, and transition IDs.
  - `internal/admission/registry_v0.1.go` includes the same command and transition mappings.
  - `internal/admission/registry_v0_1_parity_test.go` verifies the mappings against the YAML source.
  - `schemas/command-envelope.schema.json` requires `transition_ref`.
- **Executable evidence/tests:**
  - The reusable Admission Engine enforces `P0.ADMISSION.COMMAND_EXISTS`, `P0.ADMISSION.COMMAND_CAPABILITY_MATCH`, `P0.ADMISSION.COMMAND_AGGREGATE_MATCH`, and `P0.ADMISSION.COMMAND_TRANSITION_DEFINED`.
  - All five Phase D binding tests receive an ADMIT result with a non-empty `transition_ref` (`unit:created`, `version:created`, `meaning:set`, `claim:set`, `uncertainty:set`).
  - `TestEngine_UndefinedTransition` proves a command with an empty `transition_id` is rejected.
  - Phase E `ADMIT` tests assert `transition_ref` values `unit:created` and `version:created` before invoking the real Core handler.
- **Remaining gap:** None for the current defined P0 scope.
- **Next action:** Maintain the catalogue and parity test as new commands are added.

### IR-06 — Transition Evidence

- **ID:** IR-06
- **Title:** Transition Evidence
- **Normative statement:** Every successfully executed constitution-governed transition MUST produce or reference deterministic evidence sufficient to identify the resulting transition.
- **Current status:** PASS
- **Evidence/artifacts:**
  - `command-event-catalogue.yaml` links each command to its exact runtime `event_id` (`unit.created`, `version.created`, `MEANING_SET`, `CLAIM_SET`, `UNCERTAINTY_SET`).
  - `schemas/event-envelope.schema.json` defines an evidence wrapper without replacing existing verify schemas.
  - `docs/ADMISSION_ID_SPEC_v0.1.md` defines deterministic `intent_digest` and `admission_id` profiles.
- **Executable evidence/tests:**
  - All five Phase D binding tests invoke the real Core handler after Admission, observe the real `AuditEvent`, and construct a validated Event Envelope:
    - `unit.created` for `core.unit.create`
    - `version.created` for `core.version.create`
    - `MEANING_SET` for `core.meaning.set`
    - `CLAIM_SET` for `core.claim.set`
    - `UNCERTAINTY_SET` for `core.uncertainty.set`
  - Each Event Envelope is validated against `schemas/event-envelope.schema.json` and contains `command_ref`, `transition_ref`, `runtime_event_type`, and state-derived `evidence`.
  - Phase E Increment 2 (`2e38388`) demonstrates the success chain `ADMIT → version:created → real CreateVersion → coherent VersionRecord, HeadVersionID, and version.created AuditEvent` in the test-only `admissionGateProbe` helper. Phase E does not introduce new `event_id` semantics.
- **Remaining gap:** None for the current defined P0 conformance scope. Event Envelope generation is currently test-only.
- **Next action:** When production wiring is added, ensure the same Event Envelope shape is generated or referenced by the admitted-command orchestrator.

### IR-07 — Verification Independence

- **ID:** IR-07
- **Title:** Verification Independence
- **Normative statement:** Admission MUST NOT redefine canonicalization, hashing, reconstruction, or verification.
- **Current status:** PASS
- **Evidence/artifacts:**
  - `architecture-baseline.yaml` places canonicalization, snapshot hashing, replay, and verification under the authoritative Core 2.0 boundary, not the P0 constitution.
  - ADR-P0-05 `KEEP` and `Non-Goals` explicitly preserve `digiemu-canonical-json-v1`, SHA-256, verify/replay, and CLI v1.
  - P0 artifacts (`intent-envelope.schema.json`, `command-envelope.schema.json`, `event-envelope.schema.json`) carry references and metadata; they do not define new canonicalization profiles or hash algorithms.
  - `internal/verify` package targeted tests pass.
- **Executable evidence/tests:**
  - `go test ./internal/verify -timeout 120s` passes.
  - No existing `VERIFY_RESULT_SCHEMA_v1.json`, `verify_result_v2.schema.json`, or `core_2_conformance_report.schema.json` files were modified.
  - Phase E tests invoke existing real Core use cases (`CreateVersion`, `CreateUnit`) and do not modify `internal/verify`, canonicalization, hashing, replay, or verification code.
- **Remaining gap:** None identified.
- **Next action:** Maintain the separation: any future Admission artifacts must continue to reference, not redefine, Core verification semantics.

### IR-08 — Single State Identity

- **ID:** IR-08
- **Title:** Single State Identity
- **Normative statement:** The architecture layer MUST NOT introduce an alternative DigiEmu state-identity mechanism.
- **Current status:** PASS
- **Evidence/artifacts:**
  - P0 envelope schemas only introduce `intent_id`, `command_id`, and `event_id` as architectural envelope identifiers; these are explicitly documented as not yet normatively derived and not alternative state identity.
  - The envelope fields `capability_ref`, `aggregate_ref`, `command_ref`, and `transition_ref` are references into existing registries, not new identity producers.
  - `architecture-baseline.yaml` and the manifest state that state identity is produced only by the existing Core canonicalization profile.
  - `docs/ADMISSION_ID_SPEC_v0.1.md` and `internal/admission/digest.go` use `p0-intent:sha256:` and `admission:sha256:` prefixes that are deliberately outside the Core state-identity namespace.
- **Executable evidence/tests:**
  - `schemas/intent-envelope.schema.json`, `schemas/command-envelope.schema.json`, and `schemas/event-envelope.schema.json` contain no canonicalization or hashing fields.
  - Golden tests in `internal/admission/engine_test.go` prove `admission_id` is derived only from normalized Admission inputs, not from Core state.
  - Phase E `IntentID` values in `orchestrate_test.go` are test-only fixtures; they are not used to derive Core state identity or envelope identifiers.
- **Remaining gap:** `intent_id`, `command_id`, and `event_id` derivation rules are still undefined; they must not drift into alternative state-identity schemes.
- **Next action:** When normative identifier derivation is specified, explicitly document that each is an envelope-instance identifier, not a state-identity mechanism.

### IR-09 — External Authority Preservation

- **ID:** IR-09
- **Title:** External Authority Preservation
- **Normative statement:** Admission MUST NOT absorb business, legal, regulatory, ethical, or trust authority belonging outside DigiEmu Core.
- **Current status:** PASS
- **Evidence/artifacts:**
  - ADR-P0-05 §9 states `Admission != Business Approval`, `Admission != Regulatory Approval`, `Admission != Trust Assignment`.
  - `docs/CORE_2_INTEROP_CONTRACT.md`, `docs/CORE_2_BOUNDARY_MODEL.md`, and `docs/CORE_2_EXTERNAL_REVIEW_NOTE.md` preserve external authority boundaries.
  - `testdata/rs_001/README.md` and `architecture-baseline.yaml` `external_authority` section re-state that business/legal/trust authority remains outside Core.
- **Executable evidence/tests:**
  - No P0 code path claims to issue business approval or regulatory status.
  - RS-001 and the Phase D binding tests are explicitly labeled as testing architectural admissibility, not approving decisions.
  - Phase E test comments explicitly state that successful paths, partial failures, and audit failures prove only adapter-level behavior, not business, regulatory, or trust authority.
- **Remaining gap:** None identified at the specification level. Future implementation should continue to reject language or schema fields that imply approval/authority.
- **Next action:** Add an ADR non-claim test or lint that flags any Admission schema field or documentation that could be read as business/regulatory/trust authority.

### IR-10 — Fail Closed

- **ID:** IR-10
- **Title:** Fail Closed
- **Normative statement:** Unknown capabilities, commands, ownership mappings, or required references MUST NOT silently become admissible.
- **Current status:** PASS
- **Evidence/artifacts:**
  - `admission-rule-registry.yaml` v0.1 enumerates nine executable fail-closed rules with deterministic failure reason codes.
  - `internal/admission/engine.go` implements all nine rules and returns the first failure immediately.
  - `internal/admission/engine_test.go` covers architecture mismatch, unknown capability, non-mutating capability, ownership mismatch, unknown command, command/capability mismatch, command/aggregate mismatch, undefined transition, and missing required fields.
- **Executable evidence/tests:**
  - The reusable Admission Engine rejects each invalid condition with a single normative reason code.
  - All five Phase D binding tests include an `UnknownCapability` REJECT case that proves the real Core handler is not invoked and target state is unchanged.
  - `TestRS001_Parity` and `TestV01Registry_Parity` ensure the Engine and registries are synchronized.
  - Phase E Increment 1 (`de9ec49`) adds `TestPhaseE_Reject_DoesNotInvokeHandler`, which proves an `UNKNOWN_CAPABILITY` `REJECT` causes the execution closure to be skipped entirely.
- **Remaining gap:** None for the current defined P0 scope. Fail-closed is proven for all catalogued capabilities.
- **Next action:** Maintain the rule set and parity tests as new capabilities are added.

### IR-11 — Deterministic Admission

- **ID:** IR-11
- **Title:** Deterministic Admission
- **Normative statement:** Given the same normalized admission input, applicable architecture version, and rule set, Admission MUST return the same result.
- **Current status:** PASS
- **Evidence/artifacts:**
  - `docs/ADMISSION_ID_SPEC_v0.1.md` defines the `P0.ADMISSION.INTENT.v0.1` and `P0.ADMISSION.ID.v0.1` canonical input profiles.
  - `internal/admission/digest.go` implements both profiles using `digiemu-core/internal/canonicaljson` and SHA-256.
  - `admission-rule-registry.yaml` and `internal/admission/registry_v0.1.go` provide a stable, ordered rule set.
- **Executable evidence/tests:**
  - `TestEngine_Determinism_*` in `internal/admission/engine_test.go` proves:
    - same input produces the same `admission_id`
    - different payloads produce different `admission_id`
    - payload key order does not affect the `admission_id`
    - `rule_refs` and `reason_codes` order does not affect the `admission_id`
    - REJECT before transition resolution produces a stable, no-transition `admission_id`
  - Golden digests in `TestEngine_Determinism_PayloadAlpha` and `TestEngine_Determinism_PayloadBeta` match the spec examples.
  - `TestRS001_Parity` runs fixture inputs repeatedly and compares deterministic outputs.
  - Phase E tests use `admission.NewEngine(admission.V01Registry())`, reusing the same stable `V01Registry()` and deterministic `Engine.Evaluate`.
- **Remaining gap:** None for the current defined P0 scope.
- **Next action:** Maintain the canonical input profiles; introduce new profile versions if the canonical field set or ordering ever changes.

### IR-12 — Traceability

- **ID:** IR-12
- **Title:** Traceability
- **Normative statement:** An admission result SHOULD be traceable to: architecture version → capability → aggregate → command → rule → transition → event/evidence.
- **Current status:** PARTIAL
- **Evidence/artifacts:**
  - `architecture-baseline.yaml` provides the architecture revision.
  - `core-capability-registry.yaml`, `aggregate-ownership-registry.yaml`, and `command-event-catalogue.yaml` provide capability, aggregate, command, and transition mappings.
  - `admission-rule-registry.yaml` provides machine-readable rule identifiers.
  - `internal/admission/engine.go` returns `admission_result_v0.2` objects containing `architecture_revision`, `capability_ref`, `aggregate_ref`, `command_ref`, `rule_refs`, and `transition_ref`.
  - `schemas/event-envelope.schema.json` provides the event/evidence wrapper.
- **Executable evidence/tests:**
  - The test-level traceability chain is complete for all five mutating capabilities:
    - `architecture_revision` (`0.3`)
    - → `capability_ref`
    - → `aggregate_ref` (`unit`)
    - → `command_ref`
    - → `rule_refs` (all evaluated P0.ADMISSION.* rules)
    - → `transition_ref`
    - → `runtime_event_type` (real `AuditEvent.Type`)
    - → `evidence` (state-derived Event Envelope)
  - All five Phase D binding tests validate this chain end-to-end.
  - Phase E strengthens test-level traceability for representative `ADMIT` cases: `ADMIT` is correlated with real `CreateVersion` and `CreateUnit` execution; successful execution is observable through coherent Core state and `AuditEvent`; audit failure is distinguishable from mutation failure; confirmed partial mutation can be identified through state inspection; `admission.Result` remains `ADMIT` across all execution outcomes.
- **Remaining gap:**
  - The traceability chain is not yet wired in production; `cmd/digiemu` and `internal/httpapi` do not produce or carry Admission results.
  - `intent_id`, `command_id`, and `event_id` derivations are still undefined, so the event and command instance ends of the chain lack stable identifiers beyond the runtime `AuditEvent.ID`.
- **Next action:** Define deterministic derivations for `intent_id`, `command_id`, and `event_id`; then add a non-production admitted-command orchestrator that closes the production traceability chain without modifying existing Core contracts.

## Summary

|| Status | Count |
||--------|-------|
|| PASS   | 9     |
|| PARTIAL| 3     |
|| OPEN   | 0     |
|| BLOCKED| 0     |

- **Total invariants evaluated:** 12
- **Invariants with executable evidence:** 12
- **Invariants with only documentation:** 0

## Top 3 remaining P0 gaps by architectural importance

1. **IR-02 — Admission Before Mutation (production)**
   The constitutional gate is proven in test-only Phase D bindings. The `digiemu` CLI and HTTP API still invoke Core handlers directly, so Admission is not yet an architectural boundary in production.

2. **IR-12 — Production Traceability and Envelope Identifiers**
   The test-level traceability chain is complete, but production wiring and the derivations of `intent_id`, `command_id`, and `event_id` are still undefined.

3. **IR-01 — Full Test Suite on Windows**
   Targeted suites pass, but the full `go test ./...` run is blocked by an independent Windows build-cache/lock issue, not by code regression.

## Phase E Test-Only Scope Note

Phase E evidence is TEST-ONLY orchestration/conformance evidence. It proves semantic behavior through existing real Core use cases and test-only adapters/fault wrappers in `internal/kernel/usecases/orchestrate_test.go`. Phase E does NOT prove:

- production Admission enforcement
- CLI/HTTP integration
- crash durability
- concurrent `AuditLog` safety
- filesystem transactionality
- automatic recovery
- idempotency
- production retry safety
- regulatory or legal compliance

## Recommended next implementation step

Define deterministic derivation rules for `intent_id`, `command_id`, and `event_id`, and add a non-production admitted-command orchestrator that demonstrates the complete `Intent → Admission → Command → Use Case → Event Envelope` chain for one Core mutation path. This closes the IR-12 production traceability gap and prepares a safe, reviewed wiring point for the eventual IR-02 production gate, without modifying existing CLI/HTTP contracts or Core state identity.
