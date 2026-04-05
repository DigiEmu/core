package usecases

import (
	"digiemu-core/internal/kernel/domain"
	"digiemu-core/internal/kernel/ports"
)

type GetHeadVersion struct {
	Repo ports.UnitRepository
}

func (uc GetHeadVersion) GetHeadVersion(in ports.GetHeadVersionRequest) (ports.GetHeadVersionResponse, error) {
	unit, ok, err := uc.Repo.FindUnitByKey(in.UnitKey)
	if err != nil {
		return ports.GetHeadVersionResponse{}, err
	}
	if !ok {
		return ports.GetHeadVersionResponse{}, domain.ErrUnitNotFound
	}
	if unit.HeadVersionID == "" {
		return ports.GetHeadVersionResponse{}, domain.ErrVersionNotFound
	}

	v, found, err := uc.Repo.FindVersionByID(unit.HeadVersionID)
	if err != nil {
		return ports.GetHeadVersionResponse{}, err
	}
	if !found {
		return ports.GetHeadVersionResponse{}, domain.ErrVersionNotFound
	}

	return ports.GetHeadVersionResponse{
		Version: toVersionDTO(v),
	}, nil
}
