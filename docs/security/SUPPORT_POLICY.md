# DigiEmu Core Support Policy

Version: 1.0
Status: Draft

## Purpose

This document defines the support scope, expectations, and boundaries for DigiEmu Core.

DigiEmu Core is an integrity-focused infrastructure component. Support therefore prioritizes:

- deterministic behavior
- verification correctness
- contract stability
- reproducibility of evidence
- release safety

It does not imply unlimited consulting, custom integration work, or undefined compatibility promises.

---

## Support Scope

Supported areas include:

- build failures in supported release versions
- reproducible test failures in supported release versions
- verification contract regressions
- deterministic replay regressions
- documented CLI behavior regressions
- security issues within supported versions
- release artifact inconsistencies

Support may also include clarification of:

- expected CLI behavior
- fixture handling
- documented write-policy behavior
- documented release and verification procedures

---

## Out of Scope

The following are out of scope unless explicitly covered by a separate agreement:

- custom application integration
- project-specific modeling advice
- business logic consulting
- forensic investigation of third-party environments
- support for modified forks
- support for unpinned development snapshots
- support for undocumented shell/tooling workflows
- platform-specific editor misconfiguration outside documented procedures

---

## Supported Versions

Support is provided only for versions that are:

- tagged releases
- documented as supported
- built from clean release artifacts
- not end-of-support

By default, support applies to:

- the latest stable release
- the previous stable release, if not explicitly retired

Pre-release, experimental, draft, or internal builds are not guaranteed supported.

---

## Severity Levels

### Critical
A problem is Critical if it can:

- break deterministic guarantees
- falsify or corrupt verification evidence
- create ambiguous audit integrity
- cause security compromise of trusted artifacts
- invalidate a published release contract

Target handling: highest priority.

### High
A problem is High if it can:

- break normal verification workflows
- cause repeatable incorrect results
- break documented CLI usage
- produce materially misleading outputs

### Medium
A problem is Medium if it can:

- degrade usability
- cause documentation mismatch
- affect non-critical tooling behavior
- create avoidable operator confusion without corrupting core guarantees

### Low
A problem is Low if it is:

- cosmetic
- editorial
- non-contractual
- low-risk and easily worked around

---

## Response Principles

Support is governed by these principles:

- reproducibility first
- determinism first
- documented contract first
- no hidden hotfix logic
- no silent behavior changes

A reported issue should include:

- exact version or commit
- operating system
- exact command used
- exact stdout/stderr
- fixture or bundle path involved
- whether issue reproduces on a clean checkout

---

## Unsupported Issue Reports

An issue may be declined as unsupported if:

- it cannot be reproduced
- it relies on undocumented behavior
- it is caused by local shell encoding/redirection quirks outside documented procedures
- it affects a modified fork rather than official release state
- it comes from an unsupported version

Declining support does not mean the issue is unimportant. It means it is outside the defined support boundary.

---

## Security Issues

Security-relevant reports should be handled privately before public disclosure where appropriate.

Security support includes:

- assessment of reported issue
- severity classification
- fix decision or mitigation recommendation
- release note / advisory decision where applicable

No guarantee is made for immediate patching outside supported versions.

---

## Compatibility Expectations

Support assumes the user follows documented procedures for:

- verification commands
- fixture generation
- artifact comparison
- release gating
- deterministic test execution

If a workflow depends on undocumented shell behavior or ad hoc file rewriting, support may require moving first to the documented workflow.

---

## No Warranty Expansion

Support availability does not expand the warranty or guarantee scope beyond what is explicitly documented in the release and licensing terms.

Support does not mean:

- fitness for every environment
- custom architecture review
- guaranteed compatibility with all toolchains
- certification for regulated use by default

---

## Recommended Enterprise Add-On

For enterprise use, it is recommended to define separately:

- support response times
- supported platforms
- maintenance window
- security contact path
- long-term support window
- release notification mechanism
- escalation process