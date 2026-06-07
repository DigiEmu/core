# TBN / DigiEmu Boundary Note

## Purpose

Clarify the boundary between DigiEmu Core and TBN (or similar external trust/attestation systems), ensuring no scope overlap.

## Core distinction

- DigiEmu verifies decision-state reconstruction and deterministic verification evidence.
- TBN verifies agent identity, trust certification, and action attestation.

## DigiEmu Core responsibility

- canonical snapshot
- hash verification (SHA-256 over canonical JSON scope)
- replay / reconstruction
- PASS / FAIL / ERROR verify result
- conformance evidence (schemas, vectors, CLI contract)

## TBN responsibility

- agent identity
- trust state
- attestation
- certification
- external trust claims

## Complementary architecture

- DigiEmu verification reports can be used as deterministic evidence by external trust/responsibility systems.
- TBN consumes evidence but remains responsible for identity, trust, attestation, and responsibility claims.

## Non-overlap rules

- DigiEmu Core must not certify agent identity.
- DigiEmu Core must not perform trust scoring.
- DigiEmu Core must not decide action authorization.
- TBN must not replace DigiEmu deterministic decision-state reconstruction.

## Example integration boundary

- A DigiEmu `verify` JSON report is produced for a snapshot bundle.
- TBN ingests the report as evidence, attaches agent identity and attestation, and evaluates trust/responsibility.
- Decisions about identity, authorization, or trust remain in TBN; deterministic state integrity remains in DigiEmu.

## Summary

DigiEmu Core focuses on deterministic decision-state integrity and reproducible verification. TBN focuses on identity, attestation, and responsibility. Keep the layers separate and complementary.
