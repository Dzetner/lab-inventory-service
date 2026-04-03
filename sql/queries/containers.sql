-- name: ListContainers :many
SELECT id, chemical_id, room_id, label_code, quantity, unit, status, checked_out_by, created_at
FROM containers
ORDER BY id;

-- name: CreateContainer :one
INSERT INTO containers (chemical_id, room_id, label_code, quantity, unit, status, checked_out_by)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, chemical_id, room_id, label_code, quantity, unit, status, checked_out_by, created_at;

-- name: CheckoutContainer :one
UPDATE containers
SET status = 'checked_out',
    checked_out_by = $2
WHERE id = $1
RETURNING id, chemical_id, room_id, label_code, quantity, unit, status, checked_out_by, created_at;

-- name: ReturnContainer :one
UPDATE containers
SET status = 'available',
    checked_out_by = NULL
WHERE id = $1
RETURNING id, chemical_id, room_id, label_code, quantity, unit, status, checked_out_by, created_at;

-- name: FilterContainers :many
SELECT id, chemical_id, room_id, label_code, quantity, unit, status, checked_out_by, created_at
FROM containers
WHERE
    ($1 = '' OR status = $1)
  AND
    ($2 = 0 OR room_id = $2)
ORDER BY id;