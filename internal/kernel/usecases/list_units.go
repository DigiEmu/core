package usecases

import (
	"strings"

	"digiemu-core/internal/kernel/ports"
)

type ListUnits struct {
	Repo ports.UnitRepository
}

func (uc ListUnits) ListUnits(in ports.ListUnitsRequest) (ports.ListUnitsResponse, error) {
	us, err := uc.Repo.ListUnits()
	if err != nil {
		return ports.ListUnitsResponse{}, err
	}

	prefix := strings.TrimSpace(in.KeyPrefix)
	out := make([]ports.UnitDTO, 0, len(us))

	for i := 0; i < len(us); i++ {
		u := us[i]
		if prefix != "" && !strings.HasPrefix(u.Key, prefix) {
			continue
		}

		out = append(out, ports.UnitDTO{
			ID:            u.ID,
			Key:           u.Key,
			Title:         u.Title,
			Description:   u.Description,
			HeadVersionID: u.HeadVersionID,
		})
	}

	return ports.ListUnitsResponse{Units: out}, nil
}
