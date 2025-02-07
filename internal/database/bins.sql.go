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

const updateBinName = `-- name: UpdateBinName :one
UPDATE bins SET name = $3, updated_at = NOW()
WHERE name = $1
AND (parent_bin IS NOT DISTINCT FROM $2)
RETURNING id, created_at, updated_at, name, parent_bin, parent_bin_or_null
`

type UpdateBinNameParams struct {
	Name      string
	ParentBin uuid.NullUUID
	Name_2    string
}

func (q *Queries) UpdateBinName(ctx context.Context, arg UpdateBinNameParams) (Bin, error) {
	row := q.db.QueryRowContext(ctx, updateBinName, arg.Name, arg.ParentBin, arg.Name_2)
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

const updateBinParent = `-- name: UpdateBinParent :exec
UPDATE bins SET parent_bin = $3, updated_at = NOW()
WHERE name = $1
AND (parent_bin IS NOT DISTINCT FROM $2)
`

type UpdateBinParentParams struct {
	Name        string
	ParentBin   uuid.NullUUID
	ParentBin_2 uuid.NullUUID
}

func (q *Queries) UpdateBinParent(ctx context.Context, arg UpdateBinParentParams) error {
	_, err := q.db.ExecContext(ctx, updateBinParent, arg.Name, arg.ParentBin, arg.ParentBin_2)
	return err
}
