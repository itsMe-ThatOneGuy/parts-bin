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

-- name: DeleteAllBins :exec
DELETE FROM bins;

-- name: GetBin :one
SELECT * FROM bins
WHERE name = $1
AND (parent_id IS NOT DISTINCT FROM $2);

-- name: GetBinsByParent :many
SELECT id, name, parent_id FROM bins
WHERE parent_id = $1;

-- name: DeleteBin :one
DELETE FROM bins
WHERE name = $1
AND (parent_id IS NOT DISTINCT FROM $2)
RETURNING *;

-- name: UpdateBinName :one
UPDATE bins SET name = $3, updated_at = NOW()
WHERE name = $1
AND (parent_id IS NOT DISTINCT FROM $2)
RETURNING *;

-- name: UpdateBinParent :exec
UPDATE bins SET parent_id = $3, updated_at = NOW()
WHERE name = $1
AND (parent_id IS NOT DISTINCT FROM $2);
