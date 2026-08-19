package fs_test

import (
	"os"
	"path/filepath"
	"testing"

	fsrepo "digiemu-core/internal/kernel/adapters/fs"
	"digiemu-core/internal/kernel/adapters/memory"
	"digiemu-core/internal/kernel/domain"
	"digiemu-core/internal/kernel/ports"
	"digiemu-core/internal/kernel/usecases"
)

// TestFSRepo_CreateUnit_Pilot_SaveUnit_FailsBeforeCoherentUnit proves that
// CreateUnit can return an error before a coherent Unit file is observable.
//
// The .json.tmp path is pre-created as a directory. SaveUnit writes the
// marshaled unit record to <u.ID>.json.tmp before os.Rename to <u.ID>.json.
// A directory at the .tmp path causes ioutil.WriteFile to fail before any
// coherent final .json file is produced.
func TestFSRepo_CreateUnit_Pilot_SaveUnit_FailsBeforeCoherentUnit(t *testing.T) {
	dir := t.TempDir()
	repo := fsrepo.NewUnitRepo(dir)
	audit := fsrepo.NewAuditLog(dir)
	clock := memory.FakeClock{Now: 1700000000}

	key := "pilot-save-fail"
	title := "Pilot Save Fail"
	description := "test"

	u, err := domain.NewUnit(key, title, description)
	if err != nil {
		t.Fatalf("domain.NewUnit: %v", err)
	}

	// Pre-create the .json.tmp path as a directory. SaveUnit writes to this
	// path before renaming to the final .json file. With the path as a
	// directory, WriteFile fails before a coherent Unit file is produced.
	if err := os.MkdirAll(filepath.Join(dir, "units", u.ID+".json.tmp"), 0o755); err != nil {
		t.Fatalf("mkdir .json.tmp dir: %v", err)
	}

	create := usecases.CreateUnit{Repo: repo, Audit: audit, Clock: clock}
	_, err = create.CreateUnit(ports.CreateUnitRequest{
		Key:         key,
		Title:       title,
		Description: description,
		ActorID:     "pilot",
	})
	if err == nil {
		t.Fatalf("expected CreateUnit error from SaveUnit")
	}

	_, ok, findErr := repo.FindUnitByKey(key)
	if findErr != nil {
		t.Fatalf("FindUnitByKey: %v", findErr)
	}
	if ok {
		t.Fatalf("expected no coherent Unit observable after SaveUnit failure")
	}
}

// TestFSRepo_CreateUnit_Pilot_AuditAppend_Fails_AfterUnitObservable proves
// that CreateUnit can return an audit error after a coherent Unit is already
// observable. The Unit exists and can be found; the audit append failed.
//
// The audit.ndjson path is pre-created as a directory. AuditLog.Append opens
// the path for write+append, which fails when the path is a directory.
// SaveUnit has already returned nil, so a coherent Unit .json file exists.
func TestFSRepo_CreateUnit_Pilot_AuditAppend_Fails_AfterUnitObservable(t *testing.T) {
	dir := t.TempDir()
	repo := fsrepo.NewUnitRepo(dir)
	audit := fsrepo.NewAuditLog(dir)
	clock := memory.FakeClock{Now: 1700000000}

	key := "pilot-audit-fail"
	title := "Pilot Audit Fail"
	description := "test"

	// Pre-create the audit.ndjson path as a directory. AuditLog.Append calls
	// os.OpenFile with O_CREATE|O_WRONLY|O_APPEND, which fails on a directory.
	if err := os.MkdirAll(filepath.Join(dir, "audit.ndjson"), 0o755); err != nil {
		t.Fatalf("mkdir audit.ndjson dir: %v", err)
	}

	create := usecases.CreateUnit{Repo: repo, Audit: audit, Clock: clock}
	_, err := create.CreateUnit(ports.CreateUnitRequest{
		Key:         key,
		Title:       title,
		Description: description,
		ActorID:     "pilot",
	})
	if err == nil {
		t.Fatalf("expected CreateUnit error from Audit.Append")
	}

	got, ok, findErr := repo.FindUnitByKey(key)
	if findErr != nil {
		t.Fatalf("FindUnitByKey: %v", findErr)
	}
	if !ok {
		t.Fatalf("expected coherent Unit to be currently observable after Audit.Append failure")
	}
	if got.Key != key || got.Title != title || got.Description != description {
		t.Fatalf("observed Unit does not match input: %+v", got)
	}
}

// TestFSRepo_CreateUnit_Pilot_FindUnitByKey_ObservationError proves that a
// corrupt .json file in the units directory can cause FindUnitByKey to return
// an observation/read error. In that case mutation certainty cannot be
// established through FindUnitByKey.
//
// This test is intentionally standalone: it does not require a preceding
// CreateUnit error because such an integrated sequence cannot be deterministically
// produced with the current adapter without production hooks.
func TestFSRepo_CreateUnit_Pilot_FindUnitByKey_ObservationError(t *testing.T) {
	dir := t.TempDir()
	repo := fsrepo.NewUnitRepo(dir)

	// Place a .json file that cannot be decoded. FindUnitByKey's scan will
	// attempt to read it and must return the decode error.
	if err := os.WriteFile(filepath.Join(dir, "units", "corrupt.json"), []byte("not-valid-json"), 0o644); err != nil {
		t.Fatalf("write corrupt json: %v", err)
	}

	_, ok, err := repo.FindUnitByKey("any-key")
	if err == nil {
		t.Fatalf("expected FindUnitByKey observation error")
	}
	if ok {
		t.Fatalf("expected ok=false")
	}
}
