-- name: CreateBin :one
INSERT INTO bins (id, created_at, updated_at, name)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1
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
