# Release Notes — v0.5.2

- Release date: 2026-02-10
- Scope / Type: Maintenance / Hardening
- Breaking changes: None

## Kurzbeschreibung
Wartungs- und Härtungs-Release mit Fokus auf CI-Stabilität, Audit-/Tamper-Integrität und Test-Härte.

## Änderungen (ausgewählte Dateien)
- `.github/workflows/ci.yml`: Setup verbessert — `go.mod` direkt verwendet, `cache-dependency-path: go.sum` hinzugefügt und ein `go mod tidy`-Konsistenzcheck eingebaut.
- `go.sum`: wurde via `go mod tidy` (mit echten Prüf-Checksummen) hinzugefügt/aktualisiert, um CI-Caching zuverlässig zu machen.
- `internal/kernel/adapters/memory/auditlog.go`: Leseschnittstelle/Implementierung ergänzt (`Scan`) damit Adapter das `AuditLogReader`-Contract erfüllen.
- `internal/kernel/ports/audit_reader.go`: Audit-Reader Contract (Streaming API `Scan(fn func(ev AuditEvent) error) error`) wird konsistent genutzt.
- `internal/kernel/adapters/fs/unit_repo.go`: Listing von Versionen erweitert, damit `UncertaintyHash` in `ListVersionsByUnitID` zurückgegeben wird (ermöglicht Hash-Checks während VerifyAudit).
- `internal/kernel/kernel_test/uncertainty_tamper_fs_test.go`: Test-Fehlerursache behoben (korrekte Sidecar-Pfad-Verwendung) — Test bleibt strikt (deterministisches Fail bei Tamper, wenn `StrictHash=true`).
- `scripts/`: temporärer Debug-Helfer entfernt (`scripts/debug_tamper.go` gelöscht) — keine Tests/Build-Impact.

## Verifikation
Führen Sie lokal oder in CI die folgenden Befehle aus, um Konsistenz und Tests zu prüfen:

```
gofmt -w .
go vet ./...
go mod tidy
git diff --exit-code go.mod go.sum
go test ./...
```

CI-Hinweis: Der Workflow prüft `go mod tidy`-Konsistenz und verwendet `go.sum` als `cache-dependency-path`. Ein fehlschlagender `git diff --exit-code go.mod go.sum` schlägt die CI aus.

## Ergebnis
Dieses Release fixiert CI-Stabilität (echte `go.sum`-Checksummen, tidy-Check) und beseitigt Inkonsistenzen, die dazu führten, dass Tamper-Prüfungen in `VerifyAudit` nicht detektiert wurden. Keine API-Breaking-Änderungen.
