package models

import "github.com/google/uuid"

type Bin struct {
	Name     string
	ID       uuid.NullUUID
	ParentID uuid.NullUUID
	Path     string
}
