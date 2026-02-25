## v0.1.1 — Summary

### Changes (4 commits)
1) CLI/Core: support --desc/--description end-to-end (domain/ports/usecase/fs/http)
2) Docs: improved quickstart + example responses
3) Release script: patch existing release by tag (avoid untagged drafts)
4) CI: add go vet (optional race via env)

### Verification
- go test ./... (pass)
- go fmt ./... applied
- Branch: feat/v0.1.1

### Notes
- Ensure CHANGELOG includes the v0.1.1 section before tagging.
