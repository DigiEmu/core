# DigiEmu Core Conformance v0.9

Status: Public Review Draft  
Version: 0.9  
Repository: DigiEmu/core  
Scope: Conformance levels and implementation requirements for DigiEmu Core Specification v0.9.

---

## 1. Purpose

This document defines initial conformance levels for DigiEmu Core implementations.

The purpose of conformance is to make it clear what an implementation must support in order to be considered compatible with DigiEmu Core Specification v0.9.

DigiEmu Core conformance focuses on deterministic state handling, canonical snapshot processing, SHA-256 hash verification and machine-readable verification outcomes.

This document supports:

- `docs/DIGIEMU_CORE_SPEC_v0.9.md`
- `docs/TEST_VECTORS_v0.9.md`
- `docs/SNAPSHOT_HASH_v1.0.md`
- `docs/VERIFY_SPEC_v1.0.md`
- `docs/VERIFY_RESULT_SCHEMA_v1.json`

---

## 2. Terminology

The key words `MUST`, `MUST NOT`, `SHOULD`, `SHOULD NOT` and `MAY` are to be interpreted as requirement levels for DigiEmu Core conformance.

### MUST

A required behavior for the stated conformance level.

### MUST NOT

A prohibited behavior for the stated conformance level.

### SHOULD

A recommended behavior. An implementation may deviate, but should document the reason.

### SHOULD NOT

A discouraged behavior. An implementation may deviate, but should document the reason.

### MAY

An optional behavior.

---

## 3. Conformance Levels

DigiEmu Core defines three initial conformance levels:

1. DigiEmu Core Reader
2. DigiEmu Core Verifier
3. DigiEmu Core Producer

These levels are cumulative in practice, but an implementation may declare support for only one level.

For example:

- a documentation viewer may conform only as a Reader
- an audit tool may conform as a Verifier
- a full DigiEmu-compatible system may conform as Reader, Verifier and Producer

---

## 4. DigiEmu Core Reader

A DigiEmu Core Reader can inspect DigiEmu Core artifacts without necessarily verifying or producing them.

### 4.1 Reader Requirements

A conforming Reader MUST be able to:

- read a DigiEmu Core bundle or documented artifact location
- parse JSON snapshot files
- identify the declared expected hash field or expected hash file
- identify verification result files when present
- distinguish deterministic snapshot content from metadata when the bundle layout provides this separation

A conforming Reader SHOULD be able to:

- display the declared hash algorithm
- display the expected verification result
- display relevant version fields
- display links or references to supporting specification documents
- warn when required verification files are missing

A conforming Reader MAY be able to:

- render human-readable summaries
- display audit metadata
- export artifact summaries
- show differences between snapshot variants

### 4.2 Reader Non-Requirements

A Reader is not required to:

- compute hashes
- perform replay verification
- produce bundles
- sign artifacts
- validate transition chains

---

## 5. DigiEmu Core Verifier

A DigiEmu Core Verifier can independently verify whether a deterministic snapshot reconstructs to the expected hash.

### 5.1 Verifier Requirements

A conforming Verifier MUST be able to:

- parse a DigiEmu Core snapshot
- produce the canonical JSON representation required by the applicable snapshot hash specification
- compute a SHA-256 hash over the canonical JSON byte representation
- compare the computed hash with the expected hash
- return a machine-readable verification result
- distinguish `PASS`, `FAIL` and `ERROR`
- fail verification when the computed hash differs from the expected hash
- report invalid or missing input as `ERROR`, not as `FAIL`

A conforming Verifier MUST NOT:

- silently ignore deterministic fields
- include outside-hash metadata in the deterministic hash
- depend on current time for verification
- depend on random values for verification
- depend on mutable network responses for verification
- return `PASS` when the computed hash does not match the expected hash

A conforming Verifier SHOULD:

- include the computed hash in the verification result
- include the expected hash in the verification result
- include the hash algorithm identifier
- include the canonicalization version
- include a stable reason code
- include verifier version information
- provide deterministic output formatting

A conforming Verifier MAY:

- support multiple bundle formats
- support additional hash algorithms for non-normative experiments
- emit human-readable reports in addition to machine-readable results
- provide CLI, API or library interfaces

---

## 6. DigiEmu Core Producer

A DigiEmu Core Producer can create valid DigiEmu Core artifacts.

### 6.1 Producer Requirements

A conforming Producer MUST be able to:

- create a deterministic snapshot
- separate inside-hash data from outside-hash metadata
- serialize the deterministic snapshot according to the applicable canonical JSON rules
- compute the expected SHA-256 hash
- store or export the expected hash
- create a portable artifact or bundle that a conforming Verifier can inspect
- avoid non-deterministic fields inside the hash boundary

A conforming Producer MUST NOT:

- place timestamps inside the deterministic hash boundary unless explicitly part of the deterministic state
- place runtime notes inside the deterministic hash boundary
- place human comments inside the deterministic hash boundary
- create expected hash values from non-canonical JSON
- claim DigiEmu Core compatibility if independent verification cannot reproduce the expected hash

A conforming Producer SHOULD:

- include metadata outside the hash boundary
- include version information
- include the canonicalization identifier
- include the hash algorithm identifier
- include references to relevant specification documents
- provide a verification report after producing an artifact
- support the public test vectors

A conforming Producer MAY:

- sign generated bundles
- export bundles to external storage
- attach audit notes
- include human-readable summaries
- integrate with governance records or compliance systems

---

## 7. Verification Outcomes

A conforming DigiEmu Core implementation MUST use clear verification outcomes.

The minimal outcomes are:

```text
PASS
FAIL
ERROR
```

### 7.1 PASS

`PASS` means:

> Verification completed successfully and the computed deterministic hash matches the expected hash.

A result MUST NOT be marked `PASS` if:

- the computed hash is missing
- the expected hash is missing
- canonicalization failed
- the computed hash differs from the expected hash
- required deterministic fields were ignored

### 7.2 FAIL

`FAIL` means:

> Verification completed, but the computed deterministic hash does not match the expected hash.

A result SHOULD be marked `FAIL` when:

- the snapshot was tampered with
- a deterministic field was changed
- the canonical JSON differs from the expected canonical state
- the expected hash does not match the reconstructed state

### 7.3 ERROR

`ERROR` means:

> Verification could not be completed.

A result SHOULD be marked `ERROR` when:

- required files are missing
- input JSON is invalid
- expected hash format is invalid
- canonicalization cannot be completed
- the verifier encounters an unsupported version
- the bundle layout is incomplete

---

## 8. Hash Boundary Requirements

The hash boundary is the core conformance concept in DigiEmu Core.

A conforming implementation MUST clearly distinguish between:

- inside-hash deterministic state
- outside-hash metadata

Inside-hash data MUST be stable, deterministic and relevant to the reconstructed state.

Outside-hash metadata MUST NOT affect the snapshot hash.

Examples of outside-hash metadata include:

- timestamps
- comments
- runtime notes
- environment notes
- UI labels
- log formatting
- audit commentary

An implementation SHOULD document how it determines the hash boundary.

---

## 9. Canonical JSON Requirements

A conforming implementation MUST use deterministic canonicalization.

Canonicalization MUST produce the same byte representation for the same deterministic state.

Canonicalization MUST define stable handling of:

- object key ordering
- whitespace
- nested objects
- arrays
- strings
- booleans
- null values
- numbers

A conforming implementation MUST NOT compute the snapshot hash from pretty-printed JSON unless the pretty-printed representation is explicitly the normative canonical form.

---

## 10. Test Vector Requirements

A conforming Verifier SHOULD be able to reproduce the expected hash values in:

- `docs/TEST_VECTORS_v0.9.md`

A conforming Producer SHOULD be able to generate artifacts that are compatible with the same test vector model.

Before DigiEmu Core Specification v1.0, the test vector document should become part of the normative conformance process.

---

## 11. Conformance Declaration

An implementation may declare conformance using the following format:

```text
DigiEmu Core Conformance: Reader v0.9
DigiEmu Core Conformance: Verifier v0.9
DigiEmu Core Conformance: Producer v0.9
```

A full implementation may declare:

```text
DigiEmu Core Conformance: Reader + Verifier + Producer v0.9
```

The declaration SHOULD include:

- implementation name
- implementation version
- supported DigiEmu Core specification version
- supported canonicalization version
- supported hash algorithm
- supported verification result schema
- known limitations

---

## 12. Example Conformance Declaration

```text
Implementation: digiemu-core
Implementation Version: v1.0.0
Conformance Level: Reader + Verifier + Producer
DigiEmu Core Specification: v0.9
Canonicalization: canonical_json_v1
Hash Algorithm: SHA-256
Verification Outcomes: PASS / FAIL / ERROR
Known Limitations: Public review draft; conformance subject to test vector expansion.
```

---

## 13. Non-Conforming Behavior

The following behaviors are considered non-conforming for DigiEmu Core v0.9:

- returning `PASS` when hashes differ
- including timestamps in the hash boundary without explicit deterministic justification
- changing canonicalization behavior without versioning
- silently dropping deterministic fields
- treating missing expected hash values as successful verification
- mixing human-readable audit comments into deterministic state identity
- relying on non-pinned network responses during verification
- producing verification results without stable machine-readable status fields

---

## 14. Relationship to DigiEmu Secure

DigiEmu Core conformance does not imply full system security.

DigiEmu Core verifies deterministic state reconstruction and hash identity.

DigiEmu Secure may later define additional conformance requirements for:

- signatures
- key management
- tamper-evident storage
- trusted verifier execution
- artifact provenance
- verifier identity
- secure bundle transport

Until DigiEmu Secure is specified, DigiEmu Core conformance should not be presented as a complete security guarantee.

---

## 15. Current Status

DigiEmu Core Conformance v0.9 is a public review draft.

It is intended to support the DigiEmu Core Specification v0.9 by defining practical implementation levels and minimal compatibility requirements.

Future versions should add:

- normative test suite requirements
- machine-readable conformance manifests
- CLI conformance commands
- independent verifier expectations
- expanded producer requirements
- stricter bundle validation rules