package usecases

import (
	"digiemu-core/internal/kernel/domain"
	"digiemu-core/internal/kernel/ports"
)

type ListVersions struct {
	Repo ports.UnitRepository
}

func (uc ListVersions) ListVersions(in ports.ListVersionsRequest) (ports.ListVersionsResponse, error) {
	unit, ok, err := uc.Repo.FindUnitByKey(in.UnitKey)
	if err != nil {
		return ports.ListVersionsResponse{}, err
	}
	if !ok {
		return ports.ListVersionsResponse{}, domain.ErrUnitNotFound
	}

	vs, err := uc.Repo.ListVersionsByUnitID(unit.ID)
	if err != nil {
		return ports.ListVersionsResponse{}, err
	}

	out := make([]ports.VersionDTO, 0, len(vs))
	for i := 0; i < len(vs); i++ {
		out = append(out, toVersionDTO(vs[i]))
	}

	return ports.ListVersionsResponse{Versions: out}, nil
}
