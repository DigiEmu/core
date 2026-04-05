package usecases

import (
	"digiemu-core/internal/kernel/domain"
	"digiemu-core/internal/kernel/ports"
)

func toVersionDTO(v domain.Version) ports.VersionDTO {
	return ports.VersionDTO{
		ID:            v.ID,
		UnitID:        v.UnitID,
		Label:         v.Label,
		Content:       v.Content,
		PrevVersionID: v.PrevVersionID,
		ContentHash:   v.ContentHash,
		CreatedAtUnix: v.CreatedAtUnix,
		ActorID:       v.ActorID,
	}
}
