# P0 Phase D Binding Prototype — core.unit.create

- **Version:** 0.1
- **Date:** 2026-08-07
- **Status:** Proposed / test-only conformance artifact
- **Location of proof:** `internal/kernel/usecases/admission_binding_test.go`

## Purpose

This document describes a controlled, isolated Phase D binding prototype that proves
one existing Core mutation path can be invoked only after successful P0 Admission
and that the resulting runtime transition evidence can be captured in a P0 Event
Envelope.

This prototype does **not** enable production Admission enforcement, introduce
production CLI or HTTP API behavior, modify Core 2.0.0 state identity, or change
canonicalization, hashing, replay, or verification boundaries.

## Exact tested path

```text
Intent Envelope (core.unit.create)
  -> P0 Admission evaluation
  -> ADMIT
  -> Command Envelope (unit.create -> unit:created)
  -> existing usecases.CreateUnit
  -> domain.Unit persisted in repository
  -> AuditEvent Type == unit.created appended
  -> Event Envelope (runtime_event_type: unit.created)
```

## Architecture artifacts used

- `architecture-baseline.yaml` revision `0.3`
- `admission-rule-registry.yaml` v0.1
- `core-capability-registry.yaml` v0.1
- `aggregate-ownership-registry.yaml` v0.1
- `command-event-catalogue.yaml` v0.1
- `schemas/intent-envelope.schema.json`
- `schemas/event-envelope.schema.json`
- `internal/kernel/usecases/create_unit.go`
- `internal/kernel/domain/unit.go`
- `internal/kernel/adapters/memory/unit_repo.go`
- `internal/kernel/adapters/memory/auditlog.go`
- `internal/kernel/adapters/memory/audit_reader_by_unit.go`

## What executes for ADMIT

For a valid `core.unit.create` Intent:

- `architecture_revision` equals `0.3`
- `capability_ref` equals `core.unit.create`
- `aggregate_ref` equals `unit`
- `command_ref` equals `unit.create`
- the command maps to `transition_ref` `unit:created`

The test computes the deterministic `intent_digest` and `admission_id` using the
P0.ADMISSION.INTENT.v0.1 and P0.ADMISSION.ID.v0.1 profiles, then invokes the real
`usecases.CreateUnit` with the payload values from the Intent Envelope:

- `key`
- `title`
- `description`
- `actor_id`

## What does not execute for REJECT

For an Intent with an unknown `capability_ref` the Admission evaluator returns
`REJECT` with `reason_code` `UNKNOWN_CAPABILITY` and `transition_ref` is not
resolved. The test proves that `CreateUnit` is not called, no `domain.Unit` is
persisted, and no `unit.created` AuditEvent is emitted.

## Actual runtime event

`usecases.CreateUnit` emits an `AuditEvent` with:

- `Schema` = `digiemu.audit.v1`
- `Type` = `unit.created`
- `Data` = `UnitCreatedData{Key, Title}`
- `UnitID` = the created `domain.Unit.ID`
- `ActorID` = the `actor_id` from the payload

## Event Envelope mapping

From the actual runtime AuditEvent the test builds:

| Event Envelope field | Source |
|---|---|
| `schema_version` | `v0.1` |
| `architecture_revision` | `0.3` |
| `event_id` | `AuditEvent.ID` |
| `command_ref` | `unit.create` |
| `transition_ref` | `unit:created` |
| `runtime_event_type` | `unit.created` (exact, not normalized) |
| `evidence` | `{unit_id, key, title}` derived from the audit data |

The Event Envelope is validated against `schemas/event-envelope.schema.json` using
the same `jsonschema` library already used by the repository.

## Limitations

- The test is intentionally narrow: it binds one capability/command pair.
- The Admission evaluator is a test-only adapter; it does not load YAML at runtime.
- The test uses the in-memory adapters to avoid persistent side effects.
- The `admission_id` computed in the test is an Admission evidence identifier,
  not a Core state hash or bundle hash.
- Production Admission enforcement in `cmd/digiemu` or HTTP API remains outside
  the prototype.

## Explicit non-claims

This prototype does **not** claim that:

- P0 Admission is globally mandatory for Core 2.0.0
- the CLI enforces Admission before `unit create`
- the Admission gate is deployed to production
- the `admission_id` or `intent_digest` are Core state identities
- business, legal, regulatory, ethical, or trust authority is assigned to the
  architecture layer

## Invariant impact

The prototype provides test evidence toward the following constitutional
invariants, without automatically marking any of them fully satisfied:

- **IR-02 Admission Before Mutation:** exercised for `core.unit.create`; a REJECT
  path proves the mutation is not called.
- **IR-03 Explicit Capability:** the valid path uses `core.unit.create`, a
  registered mutating capability.
- **IR-04 Explicit Ownership:** the valid path resolves to `unit` aggregate.
- **IR-05 Command/Transition Correspondence:** `unit.create` resolves to
  `transition_ref` `unit:created`.
- **IR-06 Transition Evidence:** a real `unit.created` AuditEvent is produced and
  wrapped in an Event Envelope.
- **IR-10 Fail Closed:** the unknown-capability REJECT test proves the gate does
  not silently become admissible.
- **IR-12 Traceability:** the Event Envelope maps `transition_ref` to the exact
  `runtime_event_type` emitted by `usecases.CreateUnit`.

## Files

- `internal/kernel/usecases/admission_binding_test.go` — prototype Go test
- `docs/P0_PHASE_D_UNIT_CREATE_BINDING_v0.1.md` — this document
