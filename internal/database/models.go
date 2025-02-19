// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Bin struct {
	ID             uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Name           string
	ParentID       uuid.NullUUID
	ParentIDOrNull uuid.NullUUID
}

type Part struct {
	PartID    int32
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Sku       sql.NullString
	ParentID  uuid.UUID
}
