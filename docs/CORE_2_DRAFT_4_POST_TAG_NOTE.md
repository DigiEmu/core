# Core 2.0 Draft 4 Post-Tag Note

Status
------

- Tag: `v2.0.0-draft.4`
- Stability: draft / pre-release / not stable
- Note type: post-tag review-fix note
- Tag policy: Do not move or recreate the `v2.0.0-draft.4` tag.

Summary
-------

- `v2.0.0-draft.4` was tagged as a post-draft.3 runner hardening milestone.
- After tagging, a review found that `report_expected_basic.json` was still stale at 10 cases.
- A follow-up PR updated the fixture to the current 11-case conformance pack.
- A semantic regression assertion was added so the fixture does not drift again.

Current main after fix
----------------------

- Official conformance pack remains 11 cases.
- Expected conformance output remains total=11 passed=11 failed=0.
- Public JSON report shape remains unchanged.
- `report_version` remains `core-2-conformance-report-v1`.
- `v1.0.0` remains stable/latest.

Release note suggestion
-----------------------

- The GitHub release notes for `v2.0.0-draft.4` may mention that main now includes a post-tag fixture consistency fix.
- Do not retag `v2.0.0-draft.4`.
- Use main for the latest post-draft.4 review-fix state.

Validation
----------

- `go test ./...` passes.
- Experimental conformance CLI reports total=11 passed=11 failed=0.
- Experimental conformance CLI `--json` reports status PASS, total 11, passed 11, failed 0.

Compatibility statement
-----------------------

This post-tag note does not move, recreate, or replace `v2.0.0-draft.4`. Core 2.0 remains draft and pre-release, the public JSON report shape remains unchanged, and `v1.0.0` remains the stable/latest line.
