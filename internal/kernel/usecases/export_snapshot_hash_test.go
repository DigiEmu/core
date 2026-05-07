package usecases

import (
	"testing"

	"digiemu-core/internal/kernel/ports"
)

func TestSnapshotHashIgnoresCreatedAtUnix(t *testing.T) {
	unit := ports.UnitDTO{
		ID:            "unit-1",
		Key:           "demo",
		Title:         "Demo",
		Description:   "Same deterministic state",
		HeadVersionID: "version-1",
	}

	versionsA := []ports.VersionDTO{
		{
			ID:            "version-1",
			Label:         "stable-label",
			PrevVersionID: "",
			ContentHash:   "CONTENT_HASH_ABC",
			ActorID:       "actor-1",
			CreatedAtUnix: 1000,
		},
	}

	versionsB := []ports.VersionDTO{
		{
			ID:            "version-1",
			Label:         "stable-label",
			PrevVersionID: "",
			ContentHash:   "CONTENT_HASH_ABC",
			ActorID:       "actor-1",
			CreatedAtUnix: 9999,
		},
	}

	hashA := sha256HexFromLines(snapshotCanonicalLines(unit, versionsA))
	hashB := sha256HexFromLines(snapshotCanonicalLines(unit, versionsB))

	if hashA != hashB {
		t.Fatalf("CreatedAtUnix must be outside the deterministic hash boundary: hashA=%s hashB=%s", hashA, hashB)
	}
}

func TestSnapshotHashIgnoresGeneratedLabel(t *testing.T) {
	unit := ports.UnitDTO{
		ID:            "unit-1",
		Key:           "demo",
		Title:         "Demo",
		Description:   "Same deterministic state",
		HeadVersionID: "version-1",
	}

	versionsA := []ports.VersionDTO{
		{
			ID:            "version-1",
			Label:         "20260507T220000.000Z",
			PrevVersionID: "",
			ContentHash:   "CONTENT_HASH_ABC",
			ActorID:       "actor-1",
			CreatedAtUnix: 1000,
		},
	}

	versionsB := []ports.VersionDTO{
		{
			ID:            "version-1",
			Label:         "20260507T230000.000Z",
			PrevVersionID: "",
			ContentHash:   "CONTENT_HASH_ABC",
			ActorID:       "actor-1",
			CreatedAtUnix: 1000,
		},
	}

	hashA := sha256HexFromLines(snapshotCanonicalLines(unit, versionsA))
	hashB := sha256HexFromLines(snapshotCanonicalLines(unit, versionsB))

	if hashA != hashB {
		t.Fatalf("generated Label must be outside the deterministic hash boundary: hashA=%s hashB=%s", hashA, hashB)
	}
}
