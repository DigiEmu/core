This directory contains Core 2.0 hash boundary conformance vectors.

Structure:

- inside_mutation/: base and mutated snapshots that change an inside-hash field.
- outside_metadata/: base and variant snapshots that add `expected_hash_v1` metadata.

These vectors are used by focused unit tests that confirm current behavior: inside-hash changes affect `HashV1`, while `expected_hash_v1` is excluded by `ReplayV1` during reconstruction.
