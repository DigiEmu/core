package usecases

import (
	"crypto/sha256"
	"encoding/hex"

	"digiemu-core/internal/kernel/domain"
	"digiemu-core/internal/kernel/ports"
)

// CreateVersion is the usecase implementation.
type CreateVersion struct {
	Repo  ports.UnitRepository
	Audit ports.AuditLog
	Clock ports.Clock
}

// CreateVersion implements ports.CreateVersionUsecase (strict audit).
func (uc CreateVersion) CreateVersion(in ports.CreateVersionRequest) (ports.CreateVersionResponse, error) {
	if uc.Audit == nil {
		return ports.CreateVersionResponse{}, domain.ErrAuditNotConfigured
	}
	if uc.Clock == nil {
		return ports.CreateVersionResponse{}, domain.ErrClockNotConfigured
	}

	unit, ok, err := uc.Repo.FindUnitByKey(in.UnitKey)
	if err != nil {
		return ports.CreateVersionResponse{}, err
	}
	if !ok {
		return ports.CreateVersionResponse{}, domain.ErrUnitNotFound
	}

	// optimistic locking (optional)
	if in.BaseVersionID != "" && in.BaseVersionID != unit.HeadVersionID {
		return ports.CreateVersionResponse{}, domain.ErrConflict
	}

	v, err := domain.NewVersion(unit.ID, in.Label, in.Content)
	if err != nil {
		return ports.CreateVersionResponse{}, err
	}

	v.PrevVersionID = unit.HeadVersionID
	v.ActorID = actorOrUnknown(in.ActorID)
	v.CreatedAtUnix = uc.Clock.NowUnix()

	// deterministic hash (content is already trimmed by domain.NewVersion)
	canonical := v.UnitID + "\n" + v.PrevVersionID + "\n" + v.Label + "\n" + v.Content
	sum := sha256.Sum256([]byte(canonical))
	v.ContentHash = hex.EncodeToString(sum[:])

	// state first
	if err := uc.Repo.SaveVersion(v); err != nil {
		return ports.CreateVersionResponse{}, err
	}
	if err := uc.Repo.UpdateUnitHead(unit.ID, v.ID); err != nil {
		return ports.CreateVersionResponse{}, err
	}

	// strict audit: no "success" without journal entry
	ev := domain.AuditEvent{
		Schema:    "digiemu.audit.v1",
		ID:        domain.NewIDParts("evt", "version.created", v.UnitID, v.ID),
		Type:      "version.created",
		AtUnix:    v.CreatedAtUnix,
		ActorID:   v.ActorID,
		UnitID:    v.UnitID,
		VersionID: v.ID,
		Data: domain.VersionCreatedData{
			PrevVersionID: v.PrevVersionID,
			ContentHash:   v.ContentHash,
			Label:         v.Label,
		},
	}
	if err := uc.Audit.Append(ev); err != nil {
		return ports.CreateVersionResponse{}, err
	}

	return ports.CreateVersionResponse{
		VersionID: v.ID,
		UnitID:    v.UnitID,
		Label:     v.Label,
		Content:   v.Content,
	}, nil
}
