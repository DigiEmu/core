# DigiEmu Core Conformance Declaration v0.9

Status: Public Review Draft  
Specification Layer: DigiEmu Core  
Version: v0.9  
Implementation Baseline: v1.0.0  

---

## 1. Purpose

This document defines a public review draft format for DigiEmu Core conformance declarations.

A conformance declaration allows an implementation, tool, verifier, producer, reader, library, or service to state which DigiEmu Core conformance level it claims to support.

The goal is to make implementation claims:

- explicit
- reviewable
- machine-readable
- auditable
- comparable across independent implementations

This document complements:

- DigiEmu Core Specification v0.9
- DigiEmu Core Conformance v0.9
- DigiEmu Core Test Vectors v0.9
- DigiEmu Core Verify Report Examples v0.9
- DigiEmu Core Verify Report Schema v0.9

---

## 2. Conformance Levels

A DigiEmu Core implementation MAY declare support for one or more of the following levels:

```text
DigiEmu Core Reader
DigiEmu Core Verifier
DigiEmu Core Producer
```

These levels are defined in:

```text
docs/CONFORMANCE_v0.9.md
```

---

## 3. Declaration Scope

A conformance declaration SHOULD identify:

- implementation name
- implementation version
- declared conformance level or levels
- supported specification version
- supported canonicalization method
- supported hash algorithm
- supported verification result types
- supported test vectors
- known limitations
- declaration date
- responsible party or maintainer

A declaration SHOULD distinguish between:

```text
claimed support
tested support
unsupported features
known limitations
```

---

## 4. Minimal Declaration Example

```json
{
  "declaration_type": "digiemu_core_conformance_declaration",
  "declaration_version": "v0.9",
  "implementation": {
    "name": "example-digiemu-verifier",
    "version": "0.1.0"
  },
  "conformance": {
    "specification": "DigiEmu Core Specification v0.9",
    "levels": ["DigiEmu Core Verifier"]
  }
}
```

---

## 5. Full Declaration Example

```json
{
  "declaration_type": "digiemu_core_conformance_declaration",
  "declaration_version": "v0.9",
  "implementation": {
    "name": "example-digiemu-verifier",
    "version": "0.1.0",
    "type": "cli",
    "language": "go",
    "repository": "https://example.org/example-digiemu-verifier"
  },
  "conformance": {
    "specification": "DigiEmu Core Specification v0.9",
    "conformance_document": "DigiEmu Core Conformance v0.9",
    "levels": [
      "DigiEmu Core Reader",
      "DigiEmu Core Verifier"
    ]
  },
  "capabilities": {
    "canonicalization": [
      "canonical_json_v1"
    ],
    "hash_algorithms": [
      "SHA-256"
    ],
    "verification_results": [
      "PASS",
      "FAIL",
      "ERROR"
    ],
    "verify_report_schema": "docs/VERIFY_REPORT_SCHEMA_v0.9.json"
  },
  "test_vectors": {
    "document": "docs/TEST_VECTORS_v0.9.md",
    "supported": [
      "CORE-TV-001",
      "GZC-CASE-001"
    ],
    "passed": [
      "CORE-TV-001",
      "GZC-CASE-001"
    ],
    "failed": []
  },
  "limitations": [
    "Does not produce signed bundles.",
    "Does not implement DigiEmu Core Producer level.",
    "Does not validate external identity metadata."
  ],
  "declaration": {
    "declared_by": "Example Maintainer",
    "declared_at": "2026-05-07T00:00:00Z",
    "status": "self-declared"
  }
}
```

---

## 6. Declaration Fields

Recommended top-level fields:

| Field | Description |
|---|---|
| `declaration_type` | Fixed declaration type identifier |
| `declaration_version` | Version of the declaration format |
| `implementation` | Information about the declaring implementation |
| `conformance` | Claimed DigiEmu Core conformance level or levels |
| `capabilities` | Supported technical capabilities |
| `test_vectors` | Supported and tested DigiEmu Core test vectors |
| `limitations` | Known unsupported features or constraints |
| `declaration` | Maintainer, timestamp, and declaration status |

---

## 7. Declaration Type

The `declaration_type` field SHOULD be:

```text
digiemu_core_conformance_declaration
```

This allows tools to distinguish DigiEmu Core conformance declarations from other JSON documents.

---

## 8. Declaration Version

The `declaration_version` field SHOULD identify the declaration format version.

For this document, the value SHOULD be:

```text
v0.9
```

---

## 9. Implementation Information

The `implementation` object SHOULD include:

```json
{
  "name": "example-digiemu-verifier",
  "version": "0.1.0",
  "type": "cli",
  "language": "go",
  "repository": "https://example.org/example-digiemu-verifier"
}
```

Recommended implementation fields:

| Field | Description |
|---|---|
| `name` | Implementation name |
| `version` | Implementation version |
| `type` | Implementation type, for example `cli`, `library`, `service`, `reference`, or `test_harness` |
| `language` | Main implementation language |
| `repository` | Optional repository or source location |

---

## 10. Conformance Information

The `conformance` object SHOULD include:

```json
{
  "specification": "DigiEmu Core Specification v0.9",
  "conformance_document": "DigiEmu Core Conformance v0.9",
  "levels": [
    "DigiEmu Core Verifier"
  ]
}
```

The `levels` array SHOULD contain one or more of:

```text
DigiEmu Core Reader
DigiEmu Core Verifier
DigiEmu Core Producer
```

---

## 11. Capability Information

The `capabilities` object SHOULD describe supported technical features.

Example:

```json
{
  "canonicalization": [
    "canonical_json_v1"
  ],
  "hash_algorithms": [
    "SHA-256"
  ],
  "verification_results": [
    "PASS",
    "FAIL",
    "ERROR"
  ],
  "verify_report_schema": "docs/VERIFY_REPORT_SCHEMA_v0.9.json"
}
```

A verifier claiming DigiEmu Core Verifier support SHOULD support:

```text
canonical_json_v1
SHA-256
PASS
FAIL
ERROR
```

---

## 12. Test Vector Information

The `test_vectors` object SHOULD identify which DigiEmu Core test vectors were tested.

Example:

```json
{
  "document": "docs/TEST_VECTORS_v0.9.md",
  "supported": [
    "CORE-TV-001",
    "GZC-CASE-001"
  ],
  "passed": [
    "CORE-TV-001",
    "GZC-CASE-001"
  ],
  "failed": []
}
```

Implementations SHOULD NOT claim tested support for a test vector unless the implementation has actually run that test vector successfully.

---

## 13. Limitations

The `limitations` array SHOULD list known limitations.

Example:

```json
[
  "Does not produce signed bundles.",
  "Does not implement DigiEmu Core Producer level.",
  "Does not validate external identity metadata."
]
```

Limitations SHOULD be explicit.

A conformance declaration SHOULD NOT imply support for features that are listed as limitations.

---

## 14. Declaration Status

The `declaration.status` field SHOULD identify how the declaration was made.

Recommended values:

```text
self-declared
third-party-reviewed
certified
experimental
deprecated
```

For v0.9 public review draft usage, most declarations are expected to be:

```text
self-declared
```

---

## 15. Human-Readable Declaration Example

An implementation MAY include a human-readable statement alongside the machine-readable declaration.

Example:

```text
example-digiemu-verifier version 0.1.0 self-declares support for DigiEmu Core Verifier level under DigiEmu Core Specification v0.9.

It supports canonical_json_v1, SHA-256, PASS/FAIL/ERROR verification outcomes, and has passed the CORE-TV-001 and GZC-CASE-001 test vectors.

It does not claim Producer-level support.
```

---

## 16. Non-Conforming Declarations

A declaration SHOULD NOT claim conformance if:

- required capabilities for the declared level are missing
- test vector support is claimed without successful execution
- the implementation cannot distinguish `FAIL` from `ERROR`
- the implementation cannot reproduce canonical JSON consistently
- the implementation uses an unsupported hash algorithm while claiming standard conformance
- the implementation changes deterministic snapshot content during verification
- known limitations contradict the claimed conformance level

---

## 17. Relationship to Conformance v0.9

This document does not redefine DigiEmu Core conformance requirements.

The normative conformance expectations are described in:

```text
docs/CONFORMANCE_v0.9.md
```

This document only defines how an implementation may declare its claimed conformance.

---

## 18. Relationship to Verify Report Schema v0.9

Conformance declarations and verification reports are different documents.

A verification report describes the result of one verification operation.

A conformance declaration describes the claimed capability of an implementation.

Relationship:

```text
Verify Report = result of a verification run.
Conformance Declaration = capability claim of an implementation.
```

---

## 19. Current Status

This document is a Public Review Draft.

Future versions may define:

- a formal JSON Schema for conformance declarations
- stricter required fields
- external review status
- certification profiles
- compatibility with DigiEmu Secure
- compatibility with DigiEmu Enterprise