package usecases

import (
	"encoding/json"
	"errors"

	"digiemu-core/internal/kernel/domain"
	"digiemu-core/internal/kernel/ports"
)

// SetMeaning implements the semantics for setting a meaning document for a
// specific version. It is responsible for validation, hashing, persistence
// via the repository and emitting the MEANING_SET audit event.
type SetMeaning struct {
	Repo  ports.UnitRepository
	Audit ports.AuditLog
	Clock ports.Clock
}

func (uc SetMeaning) SetMeaning(in ports.SetMeaningRequest) (ports.SetMeaningResponse, error) {
	if uc.Repo == nil {
		return ports.SetMeaningResponse{}, domain.ErrUnitNotFound
	}
	if uc.Audit == nil {
		return ports.SetMeaningResponse{}, domain.ErrAuditNotConfigured
	}
	if uc.Clock == nil {
		return ports.SetMeaningResponse{}, domain.ErrClockNotConfigured
	}

	unit, ok, err := uc.Repo.FindUnitByKey(in.UnitKey)
	if err != nil {
		return ports.SetMeaningResponse{}, err
	}
	if !ok {
		return ports.SetMeaningResponse{}, domain.ErrUnitNotFound
	}

	verID := in.VersionID
	if verID == "" {
		verID = unit.HeadVersionID
	}
	if verID == "" {
		return ports.SetMeaningResponse{}, domain.ErrVersionNotFound
	}

	_, found, err := uc.Repo.FindVersionByID(verID)
	if err != nil {
		return ports.SetMeaningResponse{}, err
	}
	if !found {
		return ports.SetMeaningResponse{}, domain.ErrVersionNotFound
	}

	if len(in.MeaningJSON) > 64*1024 {
		return ports.SetMeaningResponse{}, errors.New("meaning.json too large")
	}

	var m domain.Meaning
	if err := json.Unmarshal(in.MeaningJSON, &m); err != nil {
		return ports.SetMeaningResponse{}, err
	}
	if m.SchemaVersion != "meaning/v1" {
		return ports.SetMeaningResponse{}, errors.New("unsupported schema_version")
	}

	mh, err := ComputeMeaningHash(m)
	if err != nil {
		return ports.SetMeaningResponse{}, err
	}

	if err := uc.Repo.SaveMeaning(unit.ID, verID, m, mh); err != nil {
		return ports.SetMeaningResponse{}, err
	}

	ev := domain.AuditEvent{
		Schema:    "digiemu.audit.v1",
		ID:        domain.NewIDParts("evt", "meaning_set", unit.ID, verID, mh),
		Type:      "MEANING_SET",
		AtUnix:    uc.Clock.NowUnix(),
		ActorID:   in.ActorID,
		UnitID:    unit.ID,
		VersionID: verID,
		Data: domain.MeaningSetData{
			MeaningHash:   mh,
			MeaningPath:   unit.ID + "." + verID + ".meaning.json",
			SchemaVersion: m.SchemaVersion,
			InlinePreview: &domain.MeaningInlinePreview{
				Title:   m.Title,
				Purpose: m.Purpose,
			},
		},
	}
	if err := uc.Audit.Append(ev); err != nil {
		return ports.SetMeaningResponse{}, err
	}

	return ports.SetMeaningResponse{
		UnitID:      unit.ID,
		VersionID:   verID,
		MeaningHash: mh,
	}, nil
}
