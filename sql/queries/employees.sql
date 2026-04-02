-- name: ListEmployees :many
SELECT id, full_name, role, created_at
FROM employees
ORDER BY id;

-- name: CreateEmployee :one
INSERT INTO employees (full_name, role)
VALUES ($1, $2)
RETURNING id, full_name, role, created_at;