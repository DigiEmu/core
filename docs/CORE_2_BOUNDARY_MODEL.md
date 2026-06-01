# DigiEmu Core 2.0 — Boundary Model (Draft)

Status: DRAFT

Purpose: Define the responsibilities and clear boundaries between DigiEmu Core, Secure, Enterprise, and Domain Applications for Core 2.0 hardening.

Summary
-------
DigiEmu Core 2.0 is a minimal, deterministic verification kernel whose responsibilities are limited to canonicalization, snapshot formation, hash generation, deterministic reconstruction, hash comparison, and `Verify Result` output. Everything outside these responsibilities is explicitly assigned to higher-level layers.

Layer boundaries (normative)
---------------------------
- Core (DigiEmu Core): MUST perform deterministic canonicalization, snapshot formation, hashing, reconstruction, and produce `Verify Result` outputs. Core MUST NOT perform key custody, signature generation, policy judgement, user management, or long-term evidence storage.

- Secure (DigiEmu Secure): SHOULD be responsible for cryptographic custody, signatures, evidence sealing, post-quantum migration strategies, and long-term evidence storage. Secure MAY reference Core artifacts but MUST NOT change Core's deterministic behavior.

- Enterprise (DigiEmu Enterprise): SHOULD provide roles, workflows, dashboards, QMS integrations, and organization-level governance. Enterprise SHALL consume Core and Secure outputs and MUST NOT alter Core semantics.

- Domain Applications: MAY implement domain-specific user experiences, workflows, and business logic. Domain Applications SHALL treat Core outputs as authoritative evidence and SHALL NOT embed Core responsibilities.

Hash boundary (normative)
-------------------------
Core 2.0 MUST enforce an inside/outside hash boundary. Only deterministic, evidence-relevant data belongs inside the hash. Non-deterministic or administrative metadata belongs outside the hash.

Inside the hash (examples):
- `snapshot_version`
- `input_state`
- `policy_state` (when policy is part of deterministic reconstruction)
- `decision_state`
- `system_state`
- `meaning_state`
- `reconstruction_rules`
- `canonicalization_profile`
- `hash_algorithm`

Outside the hash (examples):
- `created_at`
- `created_by`
- `human_notes`
- `ui_context`
- `environment_metadata`
- `audit_comment`
- `signature`
- `storage_location`

Document invariants (MUST/SHOULD)
-------------------------------
- Core 2.0 MUST document the canonicalization profile(s) it supports and MUST include the `canonicalization_profile` field in any inside-hash payload.
- Core 2.0 SHALL keep outside-hash metadata separate from inside-hash data in outputs and MUST NOT incorporate outside-hash metadata into computed snapshot hashes.
- Any proposed changes to the boundary MUST be recorded in a public draft and SHOULD include migration guidance.

Compatibility
-------------
These boundary definitions are drafted to preserve v1.0 behavior. They MUST NOT rename or remove existing v1.0 documents or alter the v1.0 CLI contract.
