# DigiEmu Core Abuse Cases

Version: 1.0  
Status: Draft  
Scope: CLI, bundle verification, deterministic replay, audit verification, fixture handling, release workflows

---

## 1. Purpose

This document lists realistic misuse, abuse, and failure cases relevant to DigiEmu Core.

DigiEmu Core is an integrity system. Therefore, abuse cases are not limited to classic intrusion scenarios. A serious abuse case can also be:

- producing misleading verification evidence
- hiding semantic state
- introducing nondeterministic outputs
- mutating trusted artifacts without clear provenance
- using the CLI in a way that appears valid but breaks contract assumptions

This document helps define guardrails for enterprise use and integration.

---

## 2. Abuse-Case Categories

Abuse cases are grouped into these categories:

- input abuse
- determinism abuse
- artifact abuse
- operator misuse
- audit misuse
- release/process misuse
- platform/encoding misuse

---

## 3. Input Abuse Cases

### AC-01: Malformed JSON bundle
**Scenario:** A user supplies invalid JSON in `snapshot.json` or collection files.  
**Risk:** Verification logic may crash, produce inconsistent error messages, or behave differently across environments.  
**Expected behavior:** Fail safely, deterministically, and explicitly.

### AC-02: Oversized semantic payload
**Scenario:** Extremely large `meaning`, `claims`, or `uncertainty` payloads are supplied.  
**Risk:** Resource pressure, inconsistent failure paths, denial-of-service style degradation.  
**Expected behavior:** Reject within documented size boundaries.

### AC-03: Partial bundle directory
**Scenario:** Required file exists, but collection subtrees are incomplete or corrupted.  
**Risk:** Integrator assumes verification covered more state than actually loaded.  
**Expected behavior:** Trace must make loaded scope explicit; missing optional directories must not be misrepresented.

### AC-04: BOM/encoding-tainted input
**Scenario:** JSON files contain BOM or were saved by platform tools with unexpected encoding quirks.  
**Risk:** Non-portable behavior or false mismatch.  
**Expected behavior:** BOM-safe read where supported; invalid encodings fail explicitly.

---

## 4. Determinism Abuse Cases

### AC-05: Map-order dependent hashing
**Scenario:** A code path ranges over a map and affects a public result, hash, or audit verdict.  
**Risk:** Same semantic input produces different outputs.  
**Severity:** Critical.  
**Expected behavior:** Never allow unordered map iteration in deterministic logic.

### AC-06: Hidden state via `omitempty`
**Scenario:** Two semantically different states serialize to the same output because one field disappears.  
**Risk:** Silent ambiguity in evidence or replay state.  
**Expected behavior:** Contract-critical structures must not collapse distinct states unintentionally.

### AC-07: Platform newline drift
**Scenario:** Expected reports differ only because one environment writes LF and another CRLF.  
**Risk:** Byte-exact artifact comparisons become unstable.  
**Expected behavior:** Normative fixtures must be controlled; tests should distinguish content drift from encoding drift.

### AC-08: Trace-order instability
**Scenario:** Files are loaded in nondeterministic order due to unsorted directory enumeration.  
**Risk:** Report drift and potential hash/input-order drift.  
**Expected behavior:** Explicit sort before load and before any public trace sequence dependent on traversal.

---

## 5. Artifact Abuse Cases

### AC-09: Unsafe write-back of expected hash
**Scenario:** Operator uses `write-expected` against a live or already trusted artifact.  
**Risk:** Existing expected value is overwritten or bundle trust is silently altered.  
**Expected behavior:** No overwrite unless policy explicitly allows; write reason must be stable and explicit.

### AC-10: Misleading expected fixture regeneration
**Scenario:** A maintainer regenerates expected verification output using shell tooling that changes encoding or appends extra newline bytes.  
**Risk:** Repository contains misleading normative artifacts.  
**Expected behavior:** Controlled fixture-writing procedure only.

### AC-11: Verifying the wrong bundle root
**Scenario:** A path resolves to a different snapshot root than intended.  
**Risk:** Operator believes one artifact was verified while another was actually used.  
**Expected behavior:** Trace must show root used; bundle root validation must be strict.

### AC-12: Placeholder expected hash shipped in release artifact
**Scenario:** A demo or reference bundle is published with placeholder expected hash when final expected value was intended.  
**Risk:** Consumers mistake demo placeholder for validated release evidence.  
**Expected behavior:** Placeholder bundles must be clearly marked and never confused with finalized reference artifacts.

---

## 6. Operator Misuse Cases

### AC-13: User treats warning-free guard run as full security approval
**Scenario:** Guard returns no findings, and operator concludes enterprise readiness.  
**Risk:** Remaining policy, support, or process gaps are ignored.  
**Expected behavior:** Docs must state that guard is one control, not a full assurance substitute.

### AC-14: User compares wrong example output
**Scenario:** User compares fixture `demo` with `snapshot_demo_v1` expected report.  
**Risk:** False conclusion that deterministic output is broken.  
**Expected behavior:** Example naming and test commands should make target artifact explicit.

### AC-15: User interprets non-zero verify exit as absence of useful JSON
**Scenario:** CLI correctly returns JSON on stdout and non-zero exit for verification failure, but tooling discards stdout.  
**Risk:** Evidence loss or incorrect test harness behavior.  
**Expected behavior:** Contract must document that structured JSON output remains authoritative even on verify-fail exit.

---

## 7. Audit Misuse Cases

### AC-16: Audit stream treated as ordered without canonicalization
**Scenario:** Integrator hashes or compares audit events in source arrival order only.  
**Risk:** Spurious drift or manipulated event presentation.  
**Expected behavior:** Audit verification must define canonical comparison semantics.

### AC-17: Optional audit data omitted and interpreted as “not applicable”
**Scenario:** Missing audit-adjacent fields are interpreted as intentionally absent rather than not captured.  
**Risk:** False trust in audit completeness.  
**Expected behavior:** Audit contract should distinguish explicit empty state from omission when semantically relevant.

---

## 8. Release and Process Abuse Cases

### AC-18: Broad ignore added to silence determinism findings
**Scenario:** Ignore file suppresses whole classes of findings without narrow rationale.  
**Risk:** Real integrity regressions become invisible.  
**Expected behavior:** Ignores must be narrow, justified, and documented.

### AC-19: Release tagged without rerunning byte-exact repro tests
**Scenario:** Code builds and most tests pass, but normative verification artifact drift is not checked.  
**Risk:** Published release breaks public evidence contract.  
**Expected behavior:** Repro tests are mandatory release gates.

### AC-20: Committed fixture differs from Git blob assumptions
**Scenario:** Working tree newline conversion changes file bytes relative to repository history.  
**Risk:** Local comparisons become misleading.  
**Expected behavior:** Byte-exact tests should prefer Git blob where appropriate.

---

## 9. Platform and Tooling Abuse Cases

### AC-21: PowerShell writes UTF-16 instead of UTF-8
**Scenario:** Expected report is regenerated with redirection that produces UTF-16 or BOM-affected output.  
**Risk:** Byte-exact tests fail for non-semantic reasons.  
**Expected behavior:** Use controlled file-write method; do not use uncontrolled redirection for normative fixtures.

### AC-22: External diff tool hides encoding mismatch
**Scenario:** Content looks identical in editor, but byte-level differences remain.  
**Risk:** Maintainer believes fixture is correct when it is not.  
**Expected behavior:** Use binary comparison when validating normative fixture bytes.

---

## 10. Mitigation Summary

Primary mitigations include:

- explicit sorting of deterministic sequences
- zero tolerance for map-order effects in deterministic logic
- narrow and documented ignore rules
- BOM-safe and strict input handling
- stable write-policy semantics
- controlled fixture generation
- byte-exact regression tests
- clear CLI contract documentation
- trace-based root transparency

---

## 11. Abuse Cases Requiring Ongoing Monitoring

These should remain on every release checklist:

- fixture regeneration via shell tooling
- accidental reintroduction of `omitempty` into contract-critical structures
- broadening of guard ignores
- new map iteration in hash, replay, audit, or report code
- CLI output changes that are not explicitly versioned

---

## 12. Decision Rule

A release should not be marketed as enterprise-grade if any unresolved abuse case can:

- falsify deterministic evidence
- hide semantic state
- mutate trusted artifacts ambiguously
- undermine audit reproducibility
- mislead a reasonable operator using documented commands