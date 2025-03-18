package models

import "github.com/google/uuid"

type Part struct {
	Name     string
	ID       uuid.NullUUID
	ParentID uuid.NullUUID
}
