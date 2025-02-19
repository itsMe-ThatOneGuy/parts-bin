-- name: CreatePart :one
INSERT INTO parts (id, created_at, updated_at, name, parent_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: CreateSku :exec
UPDATE parts SET sku = $2
WHERE part_id = $1;
