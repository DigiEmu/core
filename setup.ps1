Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

function Ensure-Dir([string]$Path) {
	if (-not (Test-Path $Path)) { New-Item -ItemType Directory -Path $Path | Out-Null }
}

function Write-File([string]$Path, [string]$Content) {
	$dir = Split-Path $Path -Parent
	if ($dir -and -not (Test-Path $dir)) { Ensure-Dir $dir }
	$utf8NoBom = New-Object System.Text.UTF8Encoding($false)
	[System.IO.File]::WriteAllText($Path, $Content, $utf8NoBom)
	Write-Host "Wrote: $Path"
}

# 1) go.mod (only if missing)
if (-not (Test-Path "go.mod")) {
	$gm = @'
module digiemu-core

go 1.21
'@
	Write-File "go.mod" $gm
} else {
	Write-Host "Skipped: go.mod already exists"
}

# 2) directories
$dirs = @(
	"internal/kernel/domain",
	"internal/kernel/ports",
	"internal/kernel/usecases",
	"internal/kernel/adapters/memory",
	"internal/kernel/kernel_test"
)
foreach ($d in $dirs) { Ensure-Dir $d }

# 3) README
$readme = @'
# digiemu-core (Kernel)

## Run tests
```bash
go test ./...
```
'@
Write-File "README.md" $readme

# -------------------------
# Domain
# -------------------------
$errors = @'
package domain

import "errors"

var (
	ErrInvalidUnitKey      = errors.New("invalid unit key")
	ErrInvalidUnitTitle    = errors.New("invalid unit title")
	ErrUnitAlreadyExists   = errors.New("unit already exists")
	ErrUnitNotFound        = errors.New("unit not found")
	ErrInvalidVersionLabel = errors.New("invalid version label")
	ErrEmptyContent        = errors.New("empty content")
)
'@
Write-File "internal/kernel/domain/errors.go" $errors

$id = @'
package domain

import (
	"crypto/rand"
	"encoding/hex"
)

func NewID(prefix string) string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return prefix + "_" + hex.EncodeToString(b)
}
'@
Write-File "internal/kernel/domain/id.go" $id

$unit = @'
package domain

import "strings"

type Unit struct {
	ID    string
	Key   string
	Title string
}

func NewUnit(key, title string) (Unit, error) {
	key = strings.TrimSpace(key)
	title = strings.TrimSpace(title)

	if len(key) < 3 {
		return Unit{}, ErrInvalidUnitKey
	}
	if len(title) < 3 {
		return Unit{}, ErrInvalidUnitTitle
	}

	return Unit{
		ID:    NewID("unit"),
		Key:   key,
		Title: title,
	}, nil
}
'@
Write-File "internal/kernel/domain/unit.go" $unit

$version = @'
package domain

import "strings"

type Version struct {
	ID      string
	UnitID  string
	Label   string
	Content string
}

func NewVersion(unitID, label, content string) (Version, error) {
	unitID = strings.TrimSpace(unitID)
	label = strings.TrimSpace(label)
	content = strings.TrimSpace(content)

	if unitID == "" {
		return Version{}, ErrUnitNotFound
	}
	if len(label) < 1 {
		return Version{}, ErrInvalidVersionLabel
	}
	if content == "" {
		return Version{}, ErrEmptyContent
	}

	return Version{
		ID:      NewID("ver"),
		UnitID:  unitID,
		Label:   label,
		Content: content,
	}, nil
}
'@
Write-File "internal/kernel/domain/version.go" $version

# -------------------------
# Ports
# -------------------------
$ports = @'
package ports

import "digiemu-core/internal/kernel/domain"

type UnitRepository interface {
	ExistsByKey(key string) (bool, error)
	SaveUnit(u domain.Unit) error
	FindUnitByKey(key string) (domain.Unit, bool, error)
	FindUnitByID(id string) (domain.Unit, bool, error)

	SaveVersion(v domain.Version) error
	ListVersionsByUnitID(unitID string) ([]domain.Version, error)
}
'@
Write-File "internal/kernel/ports/repositories.go" $ports

# -------------------------
# Usecases
# -------------------------
$cu = @'
package usecases

import (
	"digiemu-core/internal/kernel/domain"
	"digiemu-core/internal/kernel/ports"
)

type CreateUnitInput struct {
	Key   string
	Title string
}

type CreateUnitOutput struct {
	Unit domain.Unit
}

type CreateUnit struct {
	Repo ports.UnitRepository
}

func (uc CreateUnit) Execute(in CreateUnitInput) (CreateUnitOutput, error) {
	u, err := domain.NewUnit(in.Key, in.Title, "")
	if err != nil {
		return CreateUnitOutput{}, err
	}

	exists, err := uc.Repo.ExistsByKey(u.Key)
	if err != nil {
		return CreateUnitOutput{}, err
	}
	if exists {
		return CreateUnitOutput{}, domain.ErrUnitAlreadyExists
	}

	if err := uc.Repo.SaveUnit(u); err != nil {
		return CreateUnitOutput{}, err
	}

	return CreateUnitOutput{Unit: u}, nil
}
'@
Write-File "internal/kernel/usecases/create_unit.go" $cu

$cv = @'
package usecases

import (
	"digiemu-core/internal/kernel/domain"
	"digiemu-core/internal/kernel/ports"
)

type CreateVersionInput struct {
	UnitKey string
	Label   string
	Content string
}

type CreateVersionOutput struct {
	Version domain.Version
}

type CreateVersion struct {
	Repo ports.UnitRepository
}

func (uc CreateVersion) Execute(in CreateVersionInput) (CreateVersionOutput, error) {
	unit, ok, err := uc.Repo.FindUnitByKey(in.UnitKey)
	if err != nil {
		return CreateVersionOutput{}, err
	}
	if !ok {
		return CreateVersionOutput{}, domain.ErrUnitNotFound
	}

	v, err := domain.NewVersion(unit.ID, in.Label, in.Content)
	if err != nil {
		return CreateVersionOutput{}, err
	}

	if err := uc.Repo.SaveVersion(v); err != nil {
		return CreateVersionOutput{}, err
	}

	return CreateVersionOutput{Version: v}, nil
}
'@
Write-File "internal/kernel/usecases/create_version.go" $cv

# -------------------------
# In-memory adapter
# -------------------------
$repo = @'
package memory

import (
	"digiemu-core/internal/kernel/domain"
	"sync"
)

type UnitRepo struct {
	mu sync.RWMutex

	unitsByID  map[string]domain.Unit
	unitsByKey map[string]string // key -> unitID

	versionsByUnitID map[string][]domain.Version
}

func NewUnitRepo() *UnitRepo {
	return &UnitRepo{
		unitsByID:        map[string]domain.Unit{},
		unitsByKey:       map[string]string{},
		versionsByUnitID: map[string][]domain.Version{},
	}
}

func (r *UnitRepo) ExistsByKey(key string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.unitsByKey[key]
	return ok, nil
}

func (r *UnitRepo) SaveUnit(u domain.Unit) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.unitsByID[u.ID] = u
	r.unitsByKey[u.Key] = u.ID
	return nil
}

func (r *UnitRepo) FindUnitByKey(key string) (domain.Unit, bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, ok := r.unitsByKey[key]
	if (!ok) {
		return domain.Unit{}, false, nil
	}
	u, ok := r.unitsByID[id]
	if (!ok) {
		return domain.Unit{}, false, nil
	}
	return u, true, nil
}

func (r *UnitRepo) FindUnitByID(id string) (domain.Unit, bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	u, ok := r.unitsByID[id]
	return u, ok, nil
}

func (r *UnitRepo) SaveVersion(v domain.Version) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.versionsByUnitID[v.UnitID] = append(r.versionsByUnitID[v.UnitID], v)
	return nil
}

func (r *UnitRepo) ListVersionsByUnitID(unitID string) ([]domain.Version, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	vs := r.versionsByUnitID[unitID]
	out := make([]domain.Version, len(vs))
	copy(out, vs)
	return out, nil
}
'@
Write-File "internal/kernel/adapters/memory/unit_repo.go" $repo

# -------------------------
# Tests
# -------------------------
$cutest = @'
package kernel_test

import (
	"digiemu-core/internal/kernel/adapters/memory"
	"digiemu-core/internal/kernel/domain"
	"digiemu-core/internal/kernel/usecases"
	"testing"
)

func TestCreateUnit_HappyPath(t *testing.T) {
	repo := memory.NewUnitRepo()
	uc := usecases.CreateUnit{Repo: repo}

	out, err := uc.Execute(usecases.CreateUnitInput{
		Key:   "reglement-bau",
		Title: "Bau-Reglement",
	})
	if err != nil {
		t.Fatalf("expected nil err, got %v", err)
	}
	if out.Unit.ID == "" {
		t.Fatalf("expected id")
	}
}

func TestCreateUnit_DuplicateKey(t *testing.T) {
	repo := memory.NewUnitRepo()
	uc := usecases.CreateUnit{Repo: repo}

	_, err := uc.Execute(usecases.CreateUnitInput{Key: "abc", Title: "Title"})
	if err != nil {
		t.Fatalf("expected nil err, got %v", err)
	}

	_, err = uc.Execute(usecases.CreateUnitInput{Key: "abc", Title: "Title 2"})
	if err != domain.ErrUnitAlreadyExists {
		t.Fatalf("expected ErrUnitAlreadyExists, got %v", err)
	}
}
'@
Write-File "internal/kernel/kernel_test/create_unit_test.go" $cutest

$cvtest = @'
package kernel_test

import (
	"digiemu-core/internal/kernel/adapters/memory"
	"digiemu-core/internal/kernel/domain"
	"digiemu-core/internal/kernel/usecases"
	"testing"
)

func TestCreateVersion_HappyPath(t *testing.T) {
	repo := memory.NewUnitRepo()

	createUnit := usecases.CreateUnit{Repo: repo}
	unitOut, err := createUnit.Execute(usecases.CreateUnitInput{
		Key:   "merkblatt-steuern",
		Title: "Merkblatt Steuern",
	})
	if err != nil {
		t.Fatalf("create unit err: %v", err)
	}

	uc := usecases.CreateVersion{Repo: repo}
	out, err := uc.Execute(usecases.CreateVersionInput{
		UnitKey: "merkblatt-steuern",
		Label:   "v1",
		Content: "Inhalt 1",
	})
	if err != nil {
		t.Fatalf("expected nil err, got %v", err)
	}
	if out.Version.UnitID != unitOut.Unit.ID {
		t.Fatalf("expected unitID match")
	}

	vs, _ := repo.ListVersionsByUnitID(unitOut.Unit.ID)
	if len(vs) != 1 {
		t.Fatalf("expected 1 version, got %d", len(vs))
	}
}

func TestCreateVersion_UnitNotFound(t *testing.T) {
	repo := memory.NewUnitRepo()
	uc := usecases.CreateVersion{Repo: repo}

	_, err := uc.Execute(usecases.CreateVersionInput{
		UnitKey: "does-not-exist",
		Label:   "v1",
		Content: "x",
	})
	if err != domain.ErrUnitNotFound {
		t.Fatalf("expected ErrUnitNotFound, got %v", err)
	}
}

func TestCreateVersion_Validation(t *testing.T) {
	repo := memory.NewUnitRepo()
	createUnit := usecases.CreateUnit{Repo: repo}
	_, _ = createUnit.Execute(usecases.CreateUnitInput{Key: "abc", Title: "Title"})

	uc := usecases.CreateVersion{Repo: repo}

	_, err := uc.Execute(usecases.CreateVersionInput{UnitKey: "abc", Label: "", Content: "x"})
	if err != domain.ErrInvalidVersionLabel {
		t.Fatalf("expected ErrInvalidVersionLabel, got %v", err)
	}

	_, err = uc.Execute(usecases.CreateVersionInput{UnitKey: "abc", Label: "v1", Content: ""})
	if err != domain.ErrEmptyContent {
		t.Fatalf("expected ErrEmptyContent, got %v", err)
	}
}
'@
Write-File "internal/kernel/kernel_test/create_version_test.go" $cvtest

Write-Host ""
Write-Host "✅ Scaffold done."
Write-Host "Next:"
Write-Host "  powershell -ExecutionPolicy Bypass -File .\setup.ps1"
Write-Host "  go fmt ./..."
Write-Host "  go test ./..."

