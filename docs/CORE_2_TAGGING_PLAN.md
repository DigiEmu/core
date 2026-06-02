# Core 2.0 Draft Tagging Plan

This document explains the tagging scheme and expectations for Core 2.0 draft
releases. It is a lightweight plan to mark reproducible milestones for partner
testing and feedback. It does not promote experimental features to stable.

Draft tag example
-----------------

- Example tag: `v2.0.0-draft.1`

  Meaning: a partner-testable draft release marker. It signifies a reproducible
  milestone for review, testing, and partner feedback. It does NOT indicate
  stability or a final API/contract.

Future tags
-----------

- `v2.0.0-draft.2` — incremental improvements to conformance, reports, or API
  drafts.
- `v2.0.0-rc.1` — release candidate for broader review after experimental
  surfaces are stabilized.
- `v2.0.0` — stable release, only after explicit promotion and release notes
  that document any changes from draft artifacts.

Compatibility rules and constraints
-----------------------------------

- `v1.0` behavior must remain stable and compatible; any deviations must be
  documented and justified.
- Core 2.0 draft artifacts may evolve between draft tags and before stable
  release; they are intentionally draft and subject to change.
- Experimental CLI functionality must remain under the `experimental` namespace
  until an explicit promotion decision is made.
- Any breaking change in draft artifacts must be documented in the release notes
  and in the revision history of the relevant spec document.
- Promotion from a draft tag to a stable release requires an explicit decision
  and corresponding release notes documenting why the experimental surfaces
  are now considered stable.

Tagging guidance
----------------

- Create draft tags only from a branch that satisfies the checklist in
  `docs/CORE_2_RELEASE_CHECKLIST.md`.
- Use semver-like names with a `-draft.N` suffix for iterative draft releases.
- Keep changelogs and release notes up to date with each draft tag; include a
  clear summary of what changed from the previous draft tag.

Change history
--------------

- 2026-06-02 — Initial tagging plan added.
