package usecases

import (
	"digiemu-core/internal/kernel/domain"
	"digiemu-core/internal/kernel/ports"
)

// CreateUnit is the usecase implementation.
type CreateUnit struct {
	Repo  ports.UnitRepository
	Audit ports.AuditLog
	Clock ports.Clock
}

// CreateUnit implements ports.CreateUnitUsecase (strict audit).
func (uc CreateUnit) CreateUnit(in ports.CreateUnitRequest) (ports.CreateUnitResponse, error) {
	if uc.Audit == nil {
		return ports.CreateUnitResponse{}, domain.ErrAuditNotConfigured
	}
	if uc.Clock == nil {
		return ports.CreateUnitResponse{}, domain.ErrClockNotConfigured
	}

	u, err := domain.NewUnit(in.Key, in.Title, in.Description)
	if err != nil {
		return ports.CreateUnitResponse{}, err
	}

	exists, err := uc.Repo.ExistsByKey(u.Key)
	if err != nil {
		return ports.CreateUnitResponse{}, err
	}
	if exists {
		return ports.CreateUnitResponse{}, domain.ErrUnitAlreadyExists
	}

	if err := uc.Repo.SaveUnit(u); err != nil {
		return ports.CreateUnitResponse{}, err
	}

	ev := domain.AuditEvent{
		Schema:  "digiemu.audit.v1",
		ID:      domain.NewIDParts("evt", "unit.created", u.ID, actorOrUnknown(in.ActorID)),
		Type:    "unit.created",
		AtUnix:  uc.Clock.NowUnix(),
		ActorID: actorOrUnknown(in.ActorID),
		UnitID:  u.ID,
		Data: domain.UnitCreatedData{
			Key:   u.Key,
			Title: u.Title,
		},
	}
	if err := uc.Audit.Append(ev); err != nil {
		return ports.CreateUnitResponse{}, err
	}

	return ports.CreateUnitResponse{
		UnitID:      u.ID,
		Key:         u.Key,
		Title:       u.Title,
		Description: u.Description,
	}, nil
}

func actorOrUnknown(s string) string {
	if s == "" {
		return "unknown"
	}
	return s
}
