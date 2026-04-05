# DEPENDENCY_AUDIT.md

## Purpose

This document records the dependency review posture for DigiEmu Core before deployment, tagging, or enterprise-facing distribution.

The goal is not to claim zero risk.  
The goal is to make dependency decisions explicit, reviewable, and repeatable.

---

## Audit Scope

This audit covers:

- direct Go module dependencies
- standard library usage where security relevance is material
- CLI/runtime dependencies that affect verification behavior
- test-only dependencies if they influence release confidence
- transitive dependencies that materially affect parsing, hashing, transport, or public contract behavior

This audit does **not** replace:
- OS patch management
- container base image review
- cloud runtime review
- network perimeter review
- workstation security review

---

## Core Dependency Principles

### 1. Minimalism
DigiEmu Core should prefer as few dependencies as reasonably possible.

### 2. Determinism First
No dependency may silently introduce nondeterministic behavior into canonicalization, hashing, replay, or verification output.

### 3. Auditability
Critical logic should stay understandable without deep dependency chains.

### 4. Controlled Upgrades
Dependency upgrades must be intentional, tested, and documented.

### 5. Prefer Standard Library
Where practical, prefer Go standard library implementations for:
- JSON handling
- hashing
- I/O
- path handling
- CLI parsing
- basic HTTP behavior

---

## Dependency Categories

## A. High-Sensitivity Dependencies

These require the highest scrutiny because they can affect trust directly:

- JSON serialization / deserialization
- canonicalization logic
- hashing logic
- bundle loading logic
- replay / verify output generation
- audit serialization logic

### Requirement
Any dependency touching these paths must be:
- stable
- well understood
- pinned by module version
- covered by tests

---

## B. Medium-Sensitivity Dependencies

These affect operational behavior but not necessarily trust semantics directly:

- HTTP helpers
- API routing helpers
- filesystem adapters
- logging helpers
- CLI convenience libraries

### Requirement
These should not be allowed to change canonical outputs or verification semantics.

---

## C. Low-Sensitivity Dependencies

These include:
- internal test helpers
- developer tooling
- non-runtime documentation helpers

### Requirement
These must not leak into release-critical behavior unintentionally.

---

## Review Checklist

For each dependency, review:

- module name
- current version
- why it exists
- whether it is used in deterministic paths
- whether the standard library could replace it
- maintenance health
- security history if known
- upgrade pressure
- test coverage around its usage

---

## Dependency Risk Questions

For every dependency, ask:

1. Can it affect byte-exact output?
2. Can it affect canonical JSON or hash results?
3. Can it affect file ordering or traversal behavior?
4. Can it affect error classification or audit behavior?
5. Can it introduce hidden network, time, or environment coupling?
6. Can it be removed without harming clarity or maintainability?

If the answer to any of 1–4 is yes, the dependency is security-relevant.

---

## Current Policy

### Allowed by default
- Go standard library
- minimal, well-bounded packages with clear runtime purpose
- dependencies already covered by stable tests

### Not allowed by default
- large framework-style packages without strong justification
- dependencies that obscure canonicalization logic
- dependencies that add reflection-heavy magic into deterministic paths
- dependencies that silently rewrite JSON structure or ordering assumptions

---

## Upgrade Policy

Dependency upgrades must follow this order:

1. review changelog
2. review public API changes
3. run full build
4. run full tests
5. run guard
6. verify reproducibility fixtures
7. document notable impact in release notes

### Mandatory caution
Any upgrade affecting:
- JSON behavior
- bundle loading
- CLI output
- hashing
- replay
- verification result structure

must be treated as potentially contract-affecting.

---

## Verification Commands

Recommended commands before release:

```txt
go mod tidy
go build ./...
go test ./...
go run ./cmd/digiemu-guard --core-only --ignore-file .digiemu-guard-ignore.json