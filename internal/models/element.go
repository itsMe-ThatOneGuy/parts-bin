package models

import "github.com/google/uuid"

type Element struct {
	Type       string
	Name       string
	Sku        string
	ID         uuid.NullUUID
	ParentName string
	ParentID   uuid.NullUUID
	Path       string
}
