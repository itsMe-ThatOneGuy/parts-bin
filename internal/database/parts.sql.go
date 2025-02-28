// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: parts.sql

package database

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createPart = `-- name: CreatePart :one
INSERT INTO parts (id, created_at, updated_at, name, parent_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING part_id, id, created_at, updated_at, name, sku, parent_id
`

type CreatePartParams struct {
	Name     string
	ParentID uuid.UUID
}

func (q *Queries) CreatePart(ctx context.Context, arg CreatePartParams) (Part, error) {
	row := q.db.QueryRowContext(ctx, createPart, arg.Name, arg.ParentID)
	var i Part
	err := row.Scan(
		&i.PartID,
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Sku,
		&i.ParentID,
	)
	return i, err
}

const createSku = `-- name: CreateSku :exec
UPDATE parts SET sku = $2
WHERE part_id = $1
`

type CreateSkuParams struct {
	PartID int32
	Sku    sql.NullString
}

func (q *Queries) CreateSku(ctx context.Context, arg CreateSkuParams) error {
	_, err := q.db.ExecContext(ctx, createSku, arg.PartID, arg.Sku)
	return err
}

const deletePart = `-- name: DeletePart :exec
DELETE FROM parts
WHERE
    (name = $1 AND part_id = $2AND parent_id = $3) 
    OR sku = $4
`

type DeletePartParams struct {
	Name     string
	PartID   int32
	ParentID uuid.UUID
	Sku      sql.NullString
}

func (q *Queries) DeletePart(ctx context.Context, arg DeletePartParams) error {
	_, err := q.db.ExecContext(ctx, deletePart,
		arg.Name,
		arg.PartID,
		arg.ParentID,
		arg.Sku,
	)
	return err
}

const getPart = `-- name: GetPart :one
SELECT part_id, id, created_at, updated_at, name, sku, parent_id FROM parts
WHERE name = $1
AND (parent_id IS NOT DISTINCT FROM $2)
`

type GetPartParams struct {
	Name     string
	ParentID uuid.UUID
}

func (q *Queries) GetPart(ctx context.Context, arg GetPartParams) (Part, error) {
	row := q.db.QueryRowContext(ctx, getPart, arg.Name, arg.ParentID)
	var i Part
	err := row.Scan(
		&i.PartID,
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Sku,
		&i.ParentID,
	)
	return i, err
}

const getPartsByParent = `-- name: GetPartsByParent :many
SELECT part_id, id, created_at, updated_at, name, sku, parent_id FROM parts 
WHERE parent_id = $1
`

func (q *Queries) GetPartsByParent(ctx context.Context, parentID uuid.UUID) ([]Part, error) {
	rows, err := q.db.QueryContext(ctx, getPartsByParent, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Part
	for rows.Next() {
		var i Part
		if err := rows.Scan(
			&i.PartID,
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Name,
			&i.Sku,
			&i.ParentID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
