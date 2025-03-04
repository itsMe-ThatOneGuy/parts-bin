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

-- name: GetPartsByParent :many
SELECT * FROM parts 
WHERE parent_id = $1;

-- name: DeletePart :exec
DELETE FROM parts
WHERE
    (name = $1 AND part_id = $2AND parent_id = $3) 
    OR sku = $4;

-- name: GetPart :one
SELECT * FROM parts
WHERE name = $1
AND (parent_id IS NOT DISTINCT FROM $2);
