package models

import "github.com/google/uuid"

type Element struct {
	Type string
	Data interface{}
}

type Bin struct {
	Name     string
	ID       uuid.NullUUID
	ParentID uuid.NullUUID
}

type Part struct {
	Name     string
	ParentID uuid.NullUUID
}
