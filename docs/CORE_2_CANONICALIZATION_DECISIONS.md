# Core 2.0 — Canonicalization Decisions (Decision Record)

**Last updated:** 2026-06-01

Purpose
-------

Record current, implementation-facing canonicalization decisions and open questions for Core 2.0. This is a documentation-only decision record: it does not change code, hashes, schemas, or CLI behavior. Any future change that affects canonicalization MUST be introduced as a new explicit profile.

Principles
----------

- Do not change v1 hash behavior retroactively.
- Canonicalization changes require a new explicit profile.
- Surprising current behavior should be documented and tested before any profile migration.
- Future profile decisions must not invalidate existing v1 evidence.

Decision topics
---------------

1) json.RawMessage

- Current behavior: `json.RawMessage` is currently treated as a typed byte slice by the custom canonical JSON encoder, not as embedded raw JSON.
- Decision for now: Preserve current behavior for v1 compatibility. Do not change existing hashes.
- Future question: Decide whether a future `digiemu-canonical-json-v2` profile should support explicitly embedded raw JSON and how to express embedded/raw semantics in the profile registry.

2) Unicode normalization

- Current behavior: Strings are quoted deterministically by the encoder, but no Unicode normalization (NFC/NFKC) is applied.
- Decision for now: Preserve current behavior for v1 compatibility; do not normalize Unicode implicitly.
- Future question: Decide whether a future profile should require NFC (or another normalization form) and document migration/compatibility steps.

3) Float encoding

- Current behavior: Floats are encoded using the current deterministic Go formatting behavior implemented by the canonical encoder.
- Decision for now: Preserve current behavior for v1 compatibility and continue to exercise it with focused tests.
- Future question: Decide whether future profiles should restrict floats (e.g., decimal strings), define bounded precision, or otherwise alter numeric encoding semantics.

Notes and migration guidance
---------------------------

- Any change to canonicalization that could change existing evidence or hashes must be introduced as a new profile and accompanied by test vectors, migration guidance, and a clear compatibility policy.
- The Core 2.0 Profile Registry (`docs/CORE_2_PROFILE_REGISTRY.md` and `schemas/core_2_profile_registry.schema.json`) should be used to register future canonicalization profiles (for example `digiemu-canonical-json-v2`) rather than modifying the default runtime behavior.

Appendix: Open profile questions
--------------------------------

- Should `digiemu-canonical-json-v2` treat `json.RawMessage` as embedded raw JSON blocks, or preserve typed-bytes behavior?
- Which Unicode normalization form (if any) should be normative for future profiles?
- What numeric restrictions (if any) should be introduced to avoid cross-language float formatting differences?
