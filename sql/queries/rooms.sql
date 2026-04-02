-- name: ListRooms :many
SELECT id, name, description, created_at
FROM rooms
ORDER BY id;

-- name: CreateRoom :one
INSERT INTO rooms (name, description)
VALUES ($1, $2)
RETURNING id, name, description, created_at;