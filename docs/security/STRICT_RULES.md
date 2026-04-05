# DigiEmu Core – Strict Rules v1

## Zweck

Diese Regeln definieren die technischen Mindestanforderungen für DigiEmu Core als deterministische Wissens-Infrastruktur.

Ziel ist nicht nur funktionierende Software, sondern:

- reproduzierbare Zustände
- verifizierbare Replay-Pfade
- kryptographisch stabile Snapshots
- minimale Angriffsfläche gegen semantische und deterministische Drift

Diese Regeln gelten für alle Kernbereiche, insbesondere:

- `pkg/snapshot`
- `pkg/replay`
- `pkg/verify`
- `pkg/meaning`
- `pkg/uncertainty`

---

## 1. Grundprinzip

DigiEmu Core folgt dieser Regel:

> Gleicher Input muss immer denselben Zustand, dieselbe Rekonstruktion und dieselbe Verifikation erzeugen.

Das bedeutet:

- kein versteckter Zustand
- keine implizite Umgebung
- keine nicht-deterministischen Effekte
- keine semantisch relevanten Werte außerhalb des Snapshots

---

## 2. Harte Verbote im Core

Die folgenden Muster sind im Core verboten, sofern sie nicht ausdrücklich in einer kontrollierten Randkomponente isoliert und dokumentiert sind.

### 2.1 Zeitabhängigkeit verboten
Nicht erlaubt:

- `time.Now()`
- `time.Sleep()`
- `time.After()`
- `time.NewTicker()`
- jede indirekte Nutzung aktueller Systemzeit

Begründung:
Zeit ist nicht deterministisch und darf Kernlogik nicht beeinflussen.

---

### 2.2 Zufall verboten
Nicht erlaubt:

- `math/rand`
- kryptographischer Zufall im Core
- zufallsbasierte IDs
- implizit zufällige Reihenfolgen

Begründung:
Zufall zerstört Reproduzierbarkeit.

---

### 2.3 Nicht-deterministische Nebenläufigkeit verboten
Nicht erlaubt in hash-, snapshot-, replay- oder verify-kritischen Pfaden:

- unkontrollierte Goroutines
- `select` mit konkurrierenden Pfaden
- nicht geordnete Channel-Auswertung
- first-come-first-served Zusammenführung
- Logik, deren Ergebnis vom Scheduling abhängt

Begründung:
Concurrency ohne strenge Ordnung erzeugt nicht reproduzierbare Zustände.

---

### 2.4 Versteckter Zustand verboten
Nicht erlaubt:

- globale mutable Variablen
- Caches mit logischer Wirkung
- implizite Singleton-Zustände
- nicht serialisierte interne Steuerfelder

Begründung:
Alles, was Verhalten beeinflusst, muss explizit im Snapshot enthalten sein.

---

### 2.5 Ghost Fields verboten
Nicht erlaubt in zustandsrelevanten Structs:

- `json:"-"`
- private Felder mit Verhaltenswirkung
- `omitempty`, wenn dadurch semantisch relevante Zustände verschwinden
- Default-Verhalten, das nicht explizit im Snapshot erscheint

Begründung:
Ein Zustand darf nicht mehr wissen als der Snapshot zeigt.

---

### 2.6 `unsafe` verboten
Nicht erlaubt:

- `unsafe.Pointer`
- Speicherabbild-basierte Serialisierung
- direkte Hashbildung über Speicherlayout

Begründung:
Speicherlayout ist architektur- und compilerabhängig.

---

### 2.7 Float-Nutzung stark eingeschränkt
Nicht erlaubt in hash- oder entscheidungsrelevanter Kernlogik, solange nicht ausdrücklich begründet und abgesichert:

- `float32`
- `float64`

Bevorzugt:
- `int64`
- Fixed-Point-Repräsentation
- explizite Skalierung

Begründung:
Floats erzeugen Rundungs- und Plattformrisiken.

---

### 2.8 IO im Core verboten
Nicht erlaubt in Kernlogik:

- Dateisystemzugriff
- Netzwerkzugriff
- Umgebungsabfragen wie `os.Getenv`
- Host-/Prozesszustand als Eingabe

Begründung:
Der Core darf nur auf expliziten Inputs arbeiten.

---

## 3. Pflichtregeln

### 3.1 Vollständiger Zustandsbezug
Jeder Zustand, der Verhalten beeinflusst, muss:

- explizit modelliert
- serialisierbar
- reproduzierbar
- im Snapshot nachvollziehbar

sein.

---

### 3.2 Explizite Reihenfolge
Jede Reihenfolge muss definiert sein.

Nicht erlaubt:
- implizite Map-Reihenfolge
- ungeordnete Aggregation
- Reihenfolge durch Laufzeitverhalten

Pflicht:
- Sortierung vor Hashbildung
- stabile Traversierung
- definierte Serialisierungsreihenfolge

---

### 3.3 Kanonische Serialisierung
Alle hashrelevanten Zustände müssen in einem eindeutig definierten kanonischen Format serialisiert werden.

Pflicht:
- stabile Feldreihenfolge
- stabiles Zahlenformat
- definierte Behandlung von Null-/Leerwerten
- keine impliziten Defaults

---

### 3.4 Trennung von Rand und Kern
Nicht-deterministische Inputs dürfen nur an der Systemgrenze aufgenommen werden.

Pflicht:
- Ingress friert externe Inputs ein
- Core verarbeitet nur explizite, fixierte Daten
- Replay nutzt niemals Live-Umgebung

---

### 3.5 Unmögliche Zustände müssen scheitern
Semantisch unmögliche Zustände dürfen nicht still toleriert werden.

Beispiele:
- negative Werte, wo nur positive erlaubt sind
- Wahrscheinlichkeiten außerhalb definierter Grenzen
- inkonsistente Beziehungen
- widersprüchliche Zustände

Pflicht:
- Fehler
- Test-Fail
- oder expliziter Abbruch

---

## 4. Deterministische Invarianten

Diese Invarianten müssen jederzeit gelten.

### 4.1 Snapshot-Invariante
Gleicher Input erzeugt immer denselben Snapshot.

### 4.2 Replay-Invariante
Gleicher Snapshot erzeugt immer denselben Replay-Zustand.

### 4.3 Verify-Invariante
Die Verifikation muss denselben kanonischen Zustand und denselben Hash reproduzieren.

### 4.4 Zustandsvollständigkeit
Es darf keinen verhaltensrelevanten Zustand außerhalb des Snapshots geben.

### 4.5 Reihenfolgestabilität
Alle Aggregationen und Traversierungen müssen stabil und explizit geordnet sein.

### 4.6 Semantische Grenzinvariante
Definierte Entscheidungsgrenzen und mathematische Grenzen müssen durch Tests abgesichert sein.

---

## 5. Spezifische Prüffelder pro Package

### `pkg/snapshot`
Pflicht:
- keine versteckten Felder
- keine Map-Aggregation ohne Sortierung
- keine Live-Umgebungswerte
- kanonische Serialisierung

### `pkg/replay`
Pflicht:
- kein Scheduling-Effekt
- kein versteckter Cache
- kein zeitabhängiges Verhalten
- keine impliziten Side Effects

### `pkg/verify`
Pflicht:
- bitgenaue Reproduzierbarkeit
- gleiche Hash-Basis
- keine tolerant interpretierte Abweichung
- klare PASS/FAIL-Logik

### `pkg/meaning`
Pflicht:
- dokumentierte Regeln
- explizite Entscheidungslogik
- Grenzfalltests
- keine stillen Bedeutungsverschiebungen

### `pkg/uncertainty`
Pflicht:
- definierte mathematische Grenzen
- keine unkontrollierte Float-Drift
- klare Rundungsstrategie
- Testfälle für Randwerte

---

## 6. Testpflichten

Jede Änderung am Core muss mindestens gegen diese Prüfarten bestehen:

### 6.1 Golden Snapshot Regression
Bestehende Fixtures dürfen sich nicht unbegründet verändern.

### 6.2 Mehrfachlauf-Test
Mehrere identische Läufe müssen identische Ergebnisse erzeugen.

### 6.3 Grenzfall-Test
Prüfung von:
- Nullwerten
- leeren Inputs
- Maximalwerten
- Schwellwerten
- symmetrischen Fällen

### 6.4 Struktur-Test
Prüfung auf:
- Ghost Fields
- globale Zustände
- verbotene APIs
- implizite Defaults

---

## 7. Änderungsregeln

### 7.1 Jede semantische Änderung ist ausdrücklich
Wenn sich die Bedeutung eines Zustands ändert, muss das:

- dokumentiert
- begründet
- getestet
- bewusst freigegeben

werden.

---

### 7.2 Kein stilles Refactoring mit Bedeutungsänderung
Refactoring darf keine implizite semantische Änderung einführen.

Wenn Verhalten beeinflusst wird, ist es keine reine Strukturänderung mehr.

---

### 7.3 Sicherheit vor Bequemlichkeit
Bei Konflikt zwischen:
- Komfort
- Performance
- Determinismus
- Auditierbarkeit

haben Determinismus und Auditierbarkeit Vorrang.

---

## 8. Freigaberegel

Eine Änderung im Core gilt nur dann als akzeptabel, wenn:

1. die Strict Rules nicht verletzt werden
2. die Invarianten weiter gelten
3. Golden Snapshots bewusst unverändert bleiben oder explizit freigegeben wurden
4. keine neue versteckte Zustandsquelle eingeführt wurde
5. semantische Auswirkungen dokumentiert wurden

---

## 9. Kurzfassung

DigiEmu Core akzeptiert nur Logik, die:

- deterministisch
- explizit
- kanonisch
- reproduzierbar
- auditierbar

ist.

Alles andere gehört an den Rand des Systems oder ist unzulässig.