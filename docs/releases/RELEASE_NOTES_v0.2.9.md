# v0.2.9 — Audit & Snapshot Stabilization Release

Commit: 8893971
Tag: v0.2.9

Summary
-------
This release finalizes the audit and snapshot infrastructure and stabilizes
the kernel, adapters, and CLI export pipeline. It enables deterministic,
verifiable exports suitable for long-term archiving and compliance use cases.

Highlights
----------
- Deterministic Snapshot Hashing
	- `snapshotHash` is always included in export output
	- `auditHash` is included when exporting with `--audit`
	- Enables reproducible exports and integrity verification

- Audit System (Filesystem + Memory)
	- Unified audit readers and tails
	- Filesystem and in-memory audit log adapters
	- Supports unit-scoped and global audit traversal

- Kernel & Ports Refactor
	- Clear separation via explicit ports and usecases
	- Removal of legacy `repositories.go`
	- Improved test coverage across kernel and usecases

Changes
-------
- kernel_domain:
	- Added audit event domain model
	- Introduced structured error versions (v0.2.x lineage)
	- Snapshot hashing moved to dedicated usecase

- ports:
	- Added: audit_reader, audit_reader_by_unit, audit_tail, auditlog, export_snapshot_usecase, get_version_usecase, verify_audit_usecase, query_usecases
	- Removed: repositories.go

- usecases:
	- export_unit_snapshot, verify_audit, snapshot_hash, get_unit, get_version, get_head_version, list_units, list_versions

- adapters:
	- Filesystem: audit_reader, audit_reader_by_unit, audit_tail, auditlog, find_version_by_id, hardened index with rebuild support
	- Memory: audit_reader, audit_reader_by_unit, audit_tail, auditlog, clock, find_version_by_id

- cli_http:
	- Added CLI export command with snapshot and audit hashes
	- Updated HTTP API wiring

Verification
------------
- Tests: `go test ./...` — all passing
- CLI smoke: create unit, create version, export snapshot (snapshotHash), export snapshot with `--audit` (auditHash), audit verification

Upgrade notes
-------------
- Legacy repository ports must be migrated to the new explicit ports
- Export output now includes cryptographic hashes for downstream verification

Audience
--------
- developers, auditors, archivists, compliance and governance users
