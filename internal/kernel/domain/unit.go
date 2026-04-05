package domain

import "strings"

type Unit struct {
	ID          string
	Key         string
	Title       string
	Description string

	// v0.2: tracks current "head" version for optimistic locking and lineage
	HeadVersionID string
}

func NewUnit(key, title, description string) (Unit, error) {
	key = strings.TrimSpace(key)
	title = strings.TrimSpace(title)
	description = strings.TrimSpace(description)

	if len(key) < 3 {
		return Unit{}, ErrInvalidUnitKey
	}
	if len(title) < 3 {
		return Unit{}, ErrInvalidUnitTitle
	}

	return Unit{
		ID:            NewIDParts("unit", key, title, description),
		Key:           key,
		Title:         title,
		Description:   description,
		HeadVersionID: "",
	}, nil
}
