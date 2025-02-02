-- name: CreateBin :one
INSERT INTO bins (id, created_at, updated_at, name, parent_bin)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: DeleteAllBins :exec
DELETE FROM bins;

-- name: GetBinByID :one
SELECT * FROM bins
WHERE id = $1;

-- name: GetBinByName :one
SELECT * FROM bins
WHERE name = $1;

-- name: DeleteBinByID :one
DELETE FROM bins
WHERE id = $1
RETURNING *;

-- name: DeleteBinByName :one
DELETE FROM bins
WHERE name = $1
RETURNING *;

-- name: UpdateBinNameByID :one
UPDATE bins SET name = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateBinNameByName :one
UPDATE bins SET name = $2, updated_at = NOW()
WHERE name = $1
RETURNING *;

-- name: UpdateBinParentByID :one
UPDATE bins SET parent_bin = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateBinParentByName :one
UPDATE bins SET parent_bin = $2, updated_at = NOW()
WHERE name = $1
RETURNING *;
