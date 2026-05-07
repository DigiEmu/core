## Summary

This PR adds `docs/NEGATIVE_TEST_VECTORS_v0.9.md`.

The new public review draft defines intentional FAIL and ERROR test vectors for DigiEmu Core verification.

It includes:

- tampered snapshot FAIL cases
- malformed JSON ERROR case
- missing snapshot ERROR case
- missing expected hash ERROR case
- unsupported hash algorithm ERROR case
- unsupported canonicalization ERROR case
- invalid field type case
- expected reason codes
- conformance relevance

## Why this matters

DigiEmu Core should not only verify valid snapshots.

Implementations must also correctly distinguish:

```text
PASS  = computed hash matches expected hash
FAIL  = processable input with mismatched hash
ERROR = verification could not be completed