# DigiEmu Core – Security Checklist v1

## Zweck

Diese Checkliste dient zur systematischen Überprüfung der Sicherheits- und Determinismus-Anforderungen von DigiEmu Core.

Sie wird verwendet für:

- Code Reviews
- Pull Requests
- Refactoring
- Release-Freigaben

---

# 🔴 Stufe 1 – Kritische Prüfungen (MUSS immer erfüllt sein)

## 1. Zeitabhängigkeit

- [ ] Wird irgendwo `time.Now()` verwendet?
- [ ] Gibt es `time.Sleep()` oder ähnliche Funktionen?
- [ ] Wird Zeit indirekt über Libraries eingebracht?

👉 Falls JA → BLOCKER

---

## 2. Zufall

- [ ] Wird `math/rand` verwendet?
- [ ] Gibt es zufallsbasierte IDs oder Seeds?
- [ ] Gibt es implizite Zufallsabhängigkeiten?

👉 Falls JA → BLOCKER

---

## 3. Map-Iteration

- [ ] Wird über Maps iteriert (`range map`)?
- [ ] Ist die Reihenfolge vor Nutzung sortiert?

👉 Falls NICHT sortiert → BLOCKER

---

## 4. Concurrency

- [ ] Gibt es Goroutines im Snapshot / Replay / Verify?
- [ ] Gibt es `select` ohne deterministische Reihenfolge?
- [ ] Wird Channel-Reihenfolge ausgewertet?

👉 Falls JA → sehr kritisch prüfen / meist BLOCKER

---

## 5. Ghost Fields

- [ ] Gibt es `json:"-"` in zustandsrelevanten Structs?
- [ ] Gibt es private Felder mit Logik-Einfluss?
- [ ] Gibt es `omitempty`, das Verhalten beeinflusst?

👉 Falls JA → BLOCKER

---

## 6. Globaler Zustand

- [ ] Gibt es globale Variablen?
- [ ] Gibt es Caches mit Einfluss auf Logik?
- [ ] Gibt es Singleton-Zustände?

👉 Falls JA → kritisch prüfen

---

## 7. IO / Environment

- [ ] Wird `os.Getenv` verwendet?
- [ ] Gibt es File- oder Netzwerkzugriffe im Core?
- [ ] Wird Host-Zustand verwendet?

👉 Falls JA → BLOCKER im Core

---

## 8. Floats

- [ ] Wird `float32` oder `float64` verwendet?
- [ ] Ist Float in Entscheidungslogik involviert?
- [ ] Gibt es Rundung ohne klare Definition?

👉 Falls JA → sehr kritisch prüfen

---

# 🟠 Stufe 2 – Semantische Risiken

## 9. Entscheidungslogik

- [ ] Gibt es Grenzwerte (z. B. 0.5)?
- [ ] Sind diese durch Tests abgesichert?
- [ ] Kann kleine Änderung Entscheidung kippen?

---

## 10. Meaning / Uncertainty

- [ ] Sind Regeln dokumentiert?
- [ ] Gibt es versteckte Gewichtungen?
- [ ] Gibt es implizite Annahmen?

---

## 11. Default-Werte

- [ ] Werden Defaults automatisch gesetzt?
- [ ] Sind Defaults im Snapshot sichtbar?
- [ ] Kann Default Verhalten verändern?

---

## 12. Struktur vs. Inhalt

- [ ] Bleibt nur Struktur gleich, aber Bedeutung ändert sich?
- [ ] Gibt es "Deep Logic Swap"-Risiko?

---

# 🟡 Stufe 3 – Architektur & System

## 13. Snapshot-Vollständigkeit

- [ ] Enthält Snapshot ALLE relevanten Daten?
- [ ] Gibt es versteckte Zustände außerhalb?

---

## 14. Replay

- [ ] Ist Replay vollständig deterministisch?
- [ ] Gibt es Nebenwirkungen?

---

## 15. Verify

- [ ] Ist Verify strikt (kein „fast gleich“)?
- [ ] Wird exakt derselbe Hash berechnet?

---

## 16. Serialisierung

- [ ] Ist Serialisierung kanonisch?
- [ ] Ist Reihenfolge stabil?
- [ ] Gibt es implizite Unterschiede?

---

# 🟢 Stufe 4 – Tests

## 17. Golden Snapshot

- [ ] Bestehende Snapshots bleiben identisch?
- [ ] Änderungen bewusst dokumentiert?

---

## 18. Mehrfachlauf

- [ ] Mehrere Runs liefern identische Ergebnisse?

---

## 19. Grenzfälle

- [ ] Nullwerte getestet?
- [ ] Maximalwerte getestet?
- [ ] Schwellwerte getestet?

---

## 20. Angreifer-Denken

- [ ] Kann ich Verhalten mit minimalem Code ändern?
- [ ] Kann ich etwas verstecken?
- [ ] Kann ich Tests umgehen?

---

# 🚨 Freigabe-Regel

Ein Change darf NICHT gemerged werden, wenn:

- ein Punkt aus Stufe 1 verletzt ist
- ein semantisches Risiko nicht verstanden ist
- ein Snapshot unerklärt abweicht

---

# 🧠 Nutzung

Diese Checkliste ist kein Formular.

Sie ist ein Denkwerkzeug.

👉 Ziel:
Nicht nur Fehler finden, sondern verhindern.