# Core 2.0 — Feedback Intake Process

Use this document to understand how partners should submit structured feedback
for Core 2.0 draft artifacts (for example `v2.0.0-draft.2`). This is a lightweight
process to help the DigiEmu team triage and respond to partner input.

How to submit feedback
----------------------

- Use GitHub issues in this repository. Choose the appropriate template under
  `.github/ISSUE_TEMPLATE/`:
  - `core_2_conformance_feedback.md` — for conformance case results, reason codes, and expected/observed outcomes.
  - `core_2_integration_blocker.md` — for blocking issues related to Docker, CLI, JSON reports, or OpenAPI adoption.
  - `core_2_spec_clarification.md` — for documentation, spec, or boundary clarifications.

- Do NOT include secrets, private customer data, or sensitive logs in issues.

What to include
---------------

- Provide the Core 2.0 tag tested (e.g., `v2.0.0-draft.2`).
- State the test path used (local Go / Docker / CI).
- Give concise reproduction steps and sample commands.
- Include relevant excerpts of JSON output or logs (redact secrets).

Triaging and response
---------------------

- Issues filed via these templates will be triaged by the DigiEmu team.
- Conformance feedback is prioritized for test-vector updates and reason-code improvements.
- Integration blockers receive expedited attention; include severity and workarounds.
- Spec clarifications may result in doc updates or follow-up questions.

Expected timelines
------------------

- Initial triage within 3 business days for high-severity issues.
- Lower-severity or clarification requests will be acknowledged and scheduled.

Communication
-------------

- Use GitHub issues for the canonical record and follow-up discussion.
- For private or contract-bound feedback, coordinate with your DigiEmu contacts
  for secure channels — do not post secrets to public issues.

Change history
--------------

- 2026-06-02 — Initial feedback intake process created.
