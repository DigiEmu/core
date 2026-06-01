# DigiEmu Core 2.0 — Crypto Agility (Draft)

Status: DRAFT

Purpose: Describe hash agility, algorithm profiling, and Core-level constraints for future cryptographic migration.

Principles
----------
- Core MUST remain compatible with v1.0 proofs based on SHA-256.  
- Core MAY reference cryptographic profiles (algorithm identifiers, canonicalization profiles) but MUST NOT implement key custody or signature logic.
- Post-quantum migration is out-of-core: responsibility of DigiEmu Secure.

Current defaults
----------------
- `hash_algorithm`: `sha256` (default for v1.0 compatibility)
- `canonicalization_profile`: `digiemu-canonical-json-v1`

Agility mechanisms (draft)
--------------------------
- `crypto_profile` (metadata): Core SHALL include a `crypto_profile` reference in `Verify Result` outputs to declare the algorithm set used for computing or validating a snapshot hash.
- `multi_hash_support` (future): Core SHOULD accept a list of ancillary verification hashes (e.g., additional digests) but Core MUST treat the primary inside-hash snapshot hash as authoritative for state equivalence.
- `hash_algorithm` field: Core SHALL record the canonical hash algorithm identifier within any inside-hash payload.

Migration constraints
---------------------
- Any change to the canonical `hash_algorithm` for active snapshots MUST be documented and accompanied by a migration strategy that preserves historical proofs.
- New algorithm support SHOULD be additive (e.g., publish additional verification layers) rather than replacing historical proofs.

Interactions with Secure
------------------------
- Post-quantum signatures and key migration belong to DigiEmu Secure.  
- Secure SHALL be able to reference Core artifacts (snapshot hashes, canonicalization profile) and attach additional signatures or hybrid proofs without mutating Core outputs.
