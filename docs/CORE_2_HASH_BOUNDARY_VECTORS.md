# Core 2.0 — Inside/Outside Hash Boundary Vectors

This document describes the Core 2.0 conformance vectors used to demonstrate the boundary between fields that affect snapshot hashing (inside-hash) and metadata that is excluded from hashing (outside-hash) under current replay behavior.

Summary
-------

- Inside-hash fields: deterministic, evidence-relevant data that must change the snapshot canonical hash when mutated.
- Outside-hash metadata: auxiliary fields such as `expected_hash_v1`, notes, or UI metadata that can be excluded from the hashing scope by replay/profile rules.

Current behavior
----------------

- The current implementation (v1.x) and `ReplayV1` explicitly exclude `expected_hash_v1` from the hashed scope during reconstruction for deterministic hashing.
- Other metadata fields (timestamps, notes, arbitrary audit blobs) are currently not automatically excluded unless a profile declares them excluded; tests document the existing limitation.

Conformance vector layout
-------------------------

Repository path: `testdata/core_2_hash_boundary/`

- `inside_mutation/` — base and mutated snapshots that change an inside-hash field.
  - `base.json` — base snapshot object.
  - `mutated_inside.json` — identical to base except a canonicalized, inside-hash field changed.

- `outside_metadata/` — examples that exercise replay exclusion of `expected_hash_v1`.
  - `base.json` — base snapshot object.
  - `with_outside_metadata.json` — same as base but includes `expected_hash_v1` metadata which `ReplayV1` should exclude from hashing.

Test intent and expectations
---------------------------

- Inside-hash vector: changing a deterministic data field must change the computed `HashV1`.
- Outside-hash vector: adding `expected_hash_v1` should not change the replay-derived hashed state; tests will show this behavior and document that only `expected_hash_v1` is excluded today.

Migration note
--------------

Core 2.0 may define richer boundary rules (profiles) that exclude additional metadata fields; these vectors demonstrate current behavior and will be used to validate future profile implementations.
