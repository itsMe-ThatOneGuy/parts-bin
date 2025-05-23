-- name: CreateBin :one
INSERT INTO bins (id, created_at, updated_at, name, parent_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: GetBin :one
SELECT * FROM bins
WHERE (name = $1 AND (parent_id IS NOT DISTINCT FROM $2))
OR id = $3
OR sku = $4;

-- name: GetBinsByParent :many
SELECT * FROM bins
WHERE (parent_id = $1 OR (parent_id IS NULL AND $1 IS NULL));

-- name: UpdateBinName :one
UPDATE bins SET name = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateBinParent :exec
UPDATE bins SET parent_id = $2, updated_at = NOW()
WHERE id = $1;

-- name: UpdateBinSku :exec
UPDATE bins SET sku = $2
WHERE id = $1;

-- name: DeleteBin :exec
DELETE FROM bins
WHERE (name = $1 AND (parent_id IS NOT DISTINCT FROM $2))
OR id = $3;

