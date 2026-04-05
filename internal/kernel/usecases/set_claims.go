package usecases

import (
	"encoding/json"
	"errors"

	"digiemu-core/internal/kernel/domain"
	"digiemu-core/internal/kernel/ports"
)

type SetClaims struct {
	Repo  ports.UnitRepository
	Audit ports.AuditLog
	Clock ports.Clock
}

func (uc SetClaims) SetClaims(in ports.SetClaimsRequest) (ports.SetClaimsResponse, error) {
	if uc.Repo == nil {
		return ports.SetClaimsResponse{}, domain.ErrUnitNotFound
	}
	if uc.Audit == nil {
		return ports.SetClaimsResponse{}, domain.ErrAuditNotConfigured
	}
	if uc.Clock == nil {
		return ports.SetClaimsResponse{}, domain.ErrClockNotConfigured
	}

	unit, ok, err := uc.Repo.FindUnitByKey(in.UnitKey)
	if err != nil {
		return ports.SetClaimsResponse{}, err
	}
	if !ok {
		return ports.SetClaimsResponse{}, domain.ErrUnitNotFound
	}

	verID := in.VersionID
	if verID == "" {
		verID = unit.HeadVersionID
	}
	if verID == "" {
		return ports.SetClaimsResponse{}, domain.ErrVersionNotFound
	}

	_, found, err := uc.Repo.FindVersionByID(verID)
	if err != nil {
		return ports.SetClaimsResponse{}, err
	}
	if !found {
		return ports.SetClaimsResponse{}, domain.ErrVersionNotFound
	}

	if len(in.BodyBytes) > 64*1024 {
		return ports.SetClaimsResponse{}, errors.New("claimset.json too large")
	}

	var cs domain.ClaimSet
	if err := json.Unmarshal(in.BodyBytes, &cs); err != nil {
		return ports.SetClaimsResponse{}, err
	}
	if cs.SchemaVersion != "claimset/v0" {
		return ports.SetClaimsResponse{}, errors.New("unsupported schema_version")
	}
	if err := cs.ValidateMinimal(); err != nil {
		return ports.SetClaimsResponse{}, err
	}

	ch, err := ComputeClaimSetHashFromStruct(cs)
	if err != nil {
		return ports.SetClaimsResponse{}, err
	}

	if err := uc.Repo.SaveClaimSet(unit.ID, verID, cs, ch); err != nil {
		return ports.SetClaimsResponse{}, err
	}

	ev := domain.AuditEvent{
		Schema:    "digiemu.audit.v1",
		ID:        domain.NewIDParts("evt", "CLAIM_SET", unit.ID, verID, ch),
		Type:      "CLAIM_SET",
		AtUnix:    uc.Clock.NowUnix(),
		ActorID:   in.ActorID,
		UnitID:    unit.ID,
		VersionID: verID,
		Data: domain.ClaimSetData{
			UnitID:       unit.ID,
			VersionID:    verID,
			ClaimSetHash: ch,
			ClaimSetPath: unit.ID + "." + verID + ".claimset.json",
		},
	}
	if err := uc.Audit.Append(ev); err != nil {
		return ports.SetClaimsResponse{}, err
	}

	return ports.SetClaimsResponse{UnitID: unit.ID, VersionID: verID, ClaimSetHash: ch}, nil
}
