-- name: ListChemicals :many
SELECT id, name, cas_number, formula, sds_url, created_at
FROM chemicals
ORDER BY id;

-- name: CreateChemical :one
INSERT INTO chemicals (name, cas_number, formula, sds_url)
VALUES ($1, $2, $3, $4)
RETURNING id, name, cas_number, formula, sds_url, created_at;