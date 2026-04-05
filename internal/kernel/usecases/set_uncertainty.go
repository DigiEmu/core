package usecases

import (
	"encoding/json"
	"errors"

	"digiemu-core/internal/kernel/domain"
	"digiemu-core/internal/kernel/ports"
)

type SetUncertainty struct {
	Repo  ports.UnitRepository
	Audit ports.AuditLog
	Clock ports.Clock
}

func (uc SetUncertainty) SetUncertainty(in ports.SetUncertaintyRequest) (ports.SetUncertaintyResponse, error) {
	if uc.Repo == nil {
		return ports.SetUncertaintyResponse{}, domain.ErrUnitNotFound
	}
	if uc.Audit == nil {
		return ports.SetUncertaintyResponse{}, domain.ErrAuditNotConfigured
	}
	if uc.Clock == nil {
		return ports.SetUncertaintyResponse{}, domain.ErrClockNotConfigured
	}

	// resolve unit by key
	unit, ok, err := uc.Repo.FindUnitByKey(in.UnitKey)
	if err != nil {
		return ports.SetUncertaintyResponse{}, err
	}
	if !ok {
		return ports.SetUncertaintyResponse{}, domain.ErrUnitNotFound
	}

	// determine target version
	verID := in.VersionID
	if verID == "" {
		verID = unit.HeadVersionID
	}
	if verID == "" {
		return ports.SetUncertaintyResponse{}, domain.ErrVersionNotFound
	}

	// validate version exists
	_, found, err := uc.Repo.FindVersionByID(verID)
	if err != nil {
		return ports.SetUncertaintyResponse{}, err
	}
	if !found {
		return ports.SetUncertaintyResponse{}, domain.ErrVersionNotFound
	}

	if len(in.BodyBytes) > 64*1024 {
		return ports.SetUncertaintyResponse{}, errors.New("uncertainty.json too large")
	}

	var u domain.Uncertainty
	if err := json.Unmarshal(in.BodyBytes, &u); err != nil {
		return ports.SetUncertaintyResponse{}, err
	}
	if u.SchemaVersion != domain.UncertaintySchemaV0 {
		return ports.SetUncertaintyResponse{}, errors.New("unsupported schema_version")
	}

	if err := u.ValidateMinimal(); err != nil {
		return ports.SetUncertaintyResponse{}, err
	}

	uh, err := ComputeUncertaintyHashFromStruct(u)
	if err != nil {
		return ports.SetUncertaintyResponse{}, err
	}

	if err := uc.Repo.SaveUncertainty(unit.ID, verID, u, uh); err != nil {
		return ports.SetUncertaintyResponse{}, err
	}

	ev := domain.AuditEvent{
		Schema:    "digiemu.audit.v1",
		ID:        domain.NewIDParts("evt", "UNCERTAINTY_SET", unit.ID, verID, uh),
		Type:      "UNCERTAINTY_SET",
		AtUnix:    uc.Clock.NowUnix(),
		ActorID:   in.ActorID,
		UnitID:    unit.ID,
		VersionID: verID,
		Data: domain.UncertaintySetData{
			UnitID:          unit.ID,
			VersionID:       verID,
			UncertaintyHash: uh,
			UncertaintyPath: unit.ID + "." + verID + ".uncertainty.json",
		},
	}
	if err := uc.Audit.Append(ev); err != nil {
		return ports.SetUncertaintyResponse{}, err
	}

	return ports.SetUncertaintyResponse{UnitID: unit.ID, VersionID: verID, UncertaintyHash: uh}, nil
}
