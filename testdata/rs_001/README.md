# RS-001 — Registry-Backed Admission Conformance Reference Scenario

**Status:** Proposed conformance fixture set  
**Architecture baseline revision:** 0.3  
**Admission Result schema:** `schemas/admission_result_v0.2.schema.json`  
**Intent Envelope schema:** `schemas/intent-envelope.schema.json`  
**Command Envelope schema:** `schemas/command-envelope.schema.json`

## Scope

This directory contains the first executable, registry-backed conformance fixture set for the RS-001 reference scenario. RS-001 tests **architectural admissibility and transition integrity only** for a single real Core mutation path (`core.unit.create` → `unit` → `unit.create` → `unit:created`). It does **not** imply that DigiEmu substantively approves any business decision.

## Important non-claims

- The fixture identifiers (`rs-001/...`) are **conformance-scenario references**, not production runtime artifacts.
- The underlying mutation is `core.unit.create`; no native Core 2.0.0 Decision aggregate is claimed to exist.
- The harness does **not** invoke the actual Go `CreateUnit` handler, so no authoritative Core state is changed.
- The generated `actual_command_envelope.json` is architectural evidence, not a proof that Core 2.0 executed the command.
- No production Admission enforcement is claimed.

## Fixture convention

Each case is a sub-directory containing:

- `input.json` — a P0 Intent Envelope submitted for Admission
- `expected.json` — the expected `decision` and `reason_codes`
- `actual_admission_result.json` — produced by the harness
- `actual_command_envelope.json` — produced for ADMIT cases only

The harness is not hard-coded to the expected decisions. It reads the actual P0 registries and derives ADMIT/REJECT.

## Reason codes

| Code | Condition |
|------|-----------|
| `ARCHITECTURE_REVISION_MISMATCH` | `architecture_revision` in the intent does not match the authoritative baseline revision. |
| `UNKNOWN_CAPABILITY` | `capability_ref` is not present in `core-capability-registry.yaml`. |
| `CAPABILITY_NOT_MUTATING` | The referenced capability has `mutation: false`. |
| `OWNERSHIP_MISMATCH` | The capability is not owned by the `aggregate_ref` in `aggregate-ownership-registry.yaml`. |
| `UNKNOWN_COMMAND` | `command_ref` is not present in `command-event-catalogue.yaml`. |
| `COMMAND_CAPABILITY_MISMATCH` | The command's `capability_id` does not match the intent's `capability_ref`. |
| `COMMAND_AGGREGATE_MISMATCH` | The command's `aggregate_id` does not match the intent's `aggregate_ref`. |
| `UNDEFINED_TRANSITION` | The command is in the catalogue but has no `transition_id`. |
| `MISSING_REQUIRED_FIELD` | A required Intent Envelope property is missing. |
| `INVALID_INTENT` | The input is present but not a valid Intent Envelope for any other reason. |

## Cases

### valid_admit

Valid `core.unit.create` intent with correct `aggregate_ref: unit` and `command_ref: unit.create`. Expected `decision: ADMIT`, no reason codes, and a generated Command Envelope.

### invalid_architecture_revision

`architecture_revision` does not equal the authoritative `0.3` baseline. Expected `REJECT` with `ARCHITECTURE_REVISION_MISMATCH`.

### invalid_unknown_capability

`capability_ref` is not in the Core Capability Registry. Expected `REJECT` with `UNKNOWN_CAPABILITY`.

### invalid_capability_not_mutating

`capability_ref` is a real but read-only capability (`core.verify`). Expected `REJECT` with `CAPABILITY_NOT_MUTATING`.

### invalid_ownership_mismatch

`capability_ref: core.unit.create` is owned by `unit`, but `aggregate_ref` is not. Expected `REJECT` with `OWNERSHIP_MISMATCH`.

### invalid_unknown_command

`command_ref` is not in the Command/Event Catalogue. Expected `REJECT` with `UNKNOWN_COMMAND`.

### invalid_command_capability_mismatch

`command_ref: unit.create` is paired with `capability_ref: core.version.create`. Expected `REJECT` with `COMMAND_CAPABILITY_MISMATCH`.

### invalid_missing_required_field

The Intent Envelope is missing `command_ref`. Expected `REJECT` with `MISSING_REQUIRED_FIELD`.

### determinism_payload_alpha

Valid `core.unit.create` intent with payload `{"extra":"x","key":"alpha"}`. Expected `ADMIT`. Demonstrates the canonical `intent_digest` and `admission_id` for payload alpha.

### determinism_payload_beta

Same as `determinism_payload_alpha` except `payload.key` is `beta`. Expected `ADMIT` with a different `admission_id` from alpha, proving payload sensitivity.

### determinism_payload_order

Same canonical payload as `determinism_payload_alpha` but with `payload` object keys in a different input order (`{"key":"alpha","extra":"x"}`). Expected `ADMIT` with the same `admission_id` as `determinism_payload_alpha`, proving canonical key sorting.

## Schema validation fixtures

The `_schema_validation` directory contains reference `admission_result_v0.2` objects used to exercise the conditional `transition_ref` rule:

- `admit_no_transition.json` — ADMIT without `transition_ref`, expected to fail schema validation.
- `reject_no_transition.json` — REJECT without `transition_ref`, expected to pass schema validation.

## Execution

Run `run.ps1` with PowerShell 7 (`pwsh`) from this directory:

```powershell
pwsh -ExecutionPolicy Bypass -File .\run.ps1
```

The harness validates each `input.json` against `schemas/intent-envelope.schema.json`, evaluates it against the P0 registries using normative rule IDs, derives `intent_digest` (per `P0.ADMISSION.INTENT.v0.1`) and `admission_id` (per `P0.ADMISSION.ID.v0.1`), produces `actual_admission_result.json` (validated against `schemas/admission_result_v0.2.schema.json`), and for ADMIT cases produces `actual_command_envelope.json` (validated against `schemas/command-envelope.schema.json`). It exits `0` when all cases, including the schema validation checks, match expectations.
