# DigiEmu Core 2.0 Interoperability Examples

**Status:** Minimal illustrative example (non-production)

---

## Overview

This directory contains minimal, non-production interoperability examples showing how DigiEmu Core 2.0 decision-state artifacts may be referenced by external systems without redefining DigiEmu state identity.

These examples demonstrate the interop boundary defined in `docs/CORE_2_INTEROP_CONTRACT.md`.

---

## What DigiEmu Produces

DigiEmu produces:
- `digiemu_snapshot.json` – canonical decision-state snapshot with hashes
- `digiemu_verification_report.json` – structured verification report with PASS/FAIL outcome

These artifacts are produced under a declared DigiEmu canonicalization profile by an accountable state producer.

---

## How External Systems Reference DigiEmu

### TBN-Style Trust Systems

TBN-style systems may reference DigiEmu state identity through:
- Snapshot ID references
- Verification report ID references
- PASS/FAIL outcomes
- Decision-state hashes

TBN does **not** compute or redefine DigiEmu state identity. TBN provides its own provenance, signatures, and trust certification layers.

See: `tbn_receipt_reference_example.json`

### AntifragileOS-Style Remediation Systems

AntifragileOS-style systems may use DigiEmu PASS/FAIL outcomes before a splice operation, but must:
- Preserve the structured report reference
- Not erase diagnostic evidence
- Not claim DigiEmu certified deployment readiness

See: `antifragile_before_splice_example.json`

---

## Important Disclaimers

- **These are examples, not security credentials.**
- **Signatures are placeholders, not real cryptographic signatures.**
- **Hashes are illustrative, not real SHA-256 values.**
- **This is not a production receipt or trust credential.**
- **Do not use these examples in production systems.**

---

## Files

| File | Purpose |
|------|---------|
| `digiemu_snapshot.json` | Example DigiEmu decision-state snapshot |
| `digiemu_verification_report.json` | Example DigiEmu verification report |
| `tbn_receipt_reference_example.json` | Example TBN receipt referencing DigiEmu |
| `antifragile_before_splice_example.json` | Example AntifragileOS remediation flow |

---

## Boundary Summary

> One state identity is produced under one declared DigiEmu canonicalization profile by one accountable state producer.
>
> External systems may carry this identity.
>
> External systems do not compute, redefine, or overwrite this identity.
