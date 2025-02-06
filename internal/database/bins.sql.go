// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: bins.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const createBin = `-- name: CreateBin :one
INSERT INTO bins (id, created_at, updated_at, name, parent_bin)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING id, created_at, updated_at, name, parent_bin, parent_bin_or_null
`

type CreateBinParams struct {
	Name      string
	ParentBin uuid.NullUUID
}

func (q *Queries) CreateBin(ctx context.Context, arg CreateBinParams) (Bin, error) {
	row := q.db.QueryRowContext(ctx, createBin, arg.Name, arg.ParentBin)
	var i Bin
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.ParentBin,
		&i.ParentBinOrNull,
	)
	return i, err
}

const deleteAllBins = `-- name: DeleteAllBins :exec
DELETE FROM bins
`

func (q *Queries) DeleteAllBins(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, deleteAllBins)
	return err
}

const deleteBin = `-- name: DeleteBin :one
DELETE FROM bins
WHERE name = $1
AND (parent_bin IS NOT DISTINCT FROM $2)
RETURNING id, created_at, updated_at, name, parent_bin, parent_bin_or_null
`

type DeleteBinParams struct {
	Name      string
	ParentBin uuid.NullUUID
}

func (q *Queries) DeleteBin(ctx context.Context, arg DeleteBinParams) (Bin, error) {
	row := q.db.QueryRowContext(ctx, deleteBin, arg.Name, arg.ParentBin)
	var i Bin
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.ParentBin,
		&i.ParentBinOrNull,
	)
	return i, err
}

const deleteBinByID = `-- name: DeleteBinByID :one
DELETE FROM bins
WHERE id = $1
RETURNING id, created_at, updated_at, name, parent_bin, parent_bin_or_null
`

func (q *Queries) DeleteBinByID(ctx context.Context, id uuid.UUID) (Bin, error) {
	row := q.db.QueryRowContext(ctx, deleteBinByID, id)
	var i Bin
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.ParentBin,
		&i.ParentBinOrNull,
	)
	return i, err
}

const getBin = `-- name: GetBin :one
SELECT id, created_at, updated_at, name, parent_bin, parent_bin_or_null FROM bins
WHERE name = $1
AND (parent_bin IS NOT DISTINCT FROM $2)
`

type GetBinParams struct {
	Name      string
	ParentBin uuid.NullUUID
}

func (q *Queries) GetBin(ctx context.Context, arg GetBinParams) (Bin, error) {
	row := q.db.QueryRowContext(ctx, getBin, arg.Name, arg.ParentBin)
	var i Bin
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.ParentBin,
		&i.ParentBinOrNull,
	)
	return i, err
}

const getBinByID = `-- name: GetBinByID :one
SELECT id, created_at, updated_at, name, parent_bin, parent_bin_or_null FROM bins
WHERE id = $1
`

func (q *Queries) GetBinByID(ctx context.Context, id uuid.UUID) (Bin, error) {
	row := q.db.QueryRowContext(ctx, getBinByID, id)
	var i Bin
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.ParentBin,
		&i.ParentBinOrNull,
	)
	return i, err
}

const getBinsByParent = `-- name: GetBinsByParent :many
SELECT id, name FROM bins
WHERE parent_bin = $1
`

type GetBinsByParentRow struct {
	ID   uuid.UUID
	Name string
}

func (q *Queries) GetBinsByParent(ctx context.Context, parentBin uuid.NullUUID) ([]GetBinsByParentRow, error) {
	rows, err := q.db.QueryContext(ctx, getBinsByParent, parentBin)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetBinsByParentRow
	for rows.Next() {
		var i GetBinsByParentRow
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
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

const updateBinNameByID = `-- name: UpdateBinNameByID :one
UPDATE bins SET name = $2, updated_at = NOW()
WHERE id = $1
RETURNING id, created_at, updated_at, name, parent_bin, parent_bin_or_null
`

type UpdateBinNameByIDParams struct {
	ID   uuid.UUID
	Name string
}

func (q *Queries) UpdateBinNameByID(ctx context.Context, arg UpdateBinNameByIDParams) (Bin, error) {
	row := q.db.QueryRowContext(ctx, updateBinNameByID, arg.ID, arg.Name)
	var i Bin
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.ParentBin,
		&i.ParentBinOrNull,
	)
	return i, err
}

const updateBinNameByName = `-- name: UpdateBinNameByName :one
UPDATE bins SET name = $2, updated_at = NOW()
WHERE name = $1
RETURNING id, created_at, updated_at, name, parent_bin, parent_bin_or_null
`

type UpdateBinNameByNameParams struct {
	Name   string
	Name_2 string
}

func (q *Queries) UpdateBinNameByName(ctx context.Context, arg UpdateBinNameByNameParams) (Bin, error) {
	row := q.db.QueryRowContext(ctx, updateBinNameByName, arg.Name, arg.Name_2)
	var i Bin
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.ParentBin,
		&i.ParentBinOrNull,
	)
	return i, err
}

const updateBinParentByID = `-- name: UpdateBinParentByID :one
UPDATE bins SET parent_bin = $2, updated_at = NOW()
WHERE id = $1
RETURNING id, created_at, updated_at, name, parent_bin, parent_bin_or_null
`

type UpdateBinParentByIDParams struct {
	ID        uuid.UUID
	ParentBin uuid.NullUUID
}

func (q *Queries) UpdateBinParentByID(ctx context.Context, arg UpdateBinParentByIDParams) (Bin, error) {
	row := q.db.QueryRowContext(ctx, updateBinParentByID, arg.ID, arg.ParentBin)
	var i Bin
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.ParentBin,
		&i.ParentBinOrNull,
	)
	return i, err
}

const updateBinParentByName = `-- name: UpdateBinParentByName :one
UPDATE bins SET parent_bin = $2, updated_at = NOW()
WHERE name = $1
RETURNING id, created_at, updated_at, name, parent_bin, parent_bin_or_null
`

type UpdateBinParentByNameParams struct {
	Name      string
	ParentBin uuid.NullUUID
}

func (q *Queries) UpdateBinParentByName(ctx context.Context, arg UpdateBinParentByNameParams) (Bin, error) {
	row := q.db.QueryRowContext(ctx, updateBinParentByName, arg.Name, arg.ParentBin)
	var i Bin
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.ParentBin,
		&i.ParentBinOrNull,
	)
	return i, err
}
