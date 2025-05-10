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

-- name: GetPart :one
SELECT * FROM parts
WHERE (name = $1 AND (parent_id IS NOT DISTINCT FROM $2))
OR sku = $3
OR id = $4;

-- name: GetPartsByParent :many
SELECT * FROM parts 
WHERE parent_id = $1;

-- name: UpdatePartName :one
UPDATE parts SET name = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdatePartParent :exec
UPDATE parts SET parent_id = $2, updated_at = NOW()
WHERE id = $1;

-- name: UpdatePartSku :exec
UPDATE parts SET sku = $2
WHERE id = $1;

-- name: DeletePart :exec
DELETE FROM parts
WHERE (name = $1 AND serial_number = $2 AND parent_id = $3) 
OR sku = $4
OR id = $5;

-- name: DeleteManyParts :exec
WITH to_delete AS (
    SELECT p.id FROM parts p
    WHERE p.name = $1 AND p.parent_id = $2
    LIMIT $3
)
DELETE FROM parts
WHERE id IN (SELECT id FROM to_delete);

