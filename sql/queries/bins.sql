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
