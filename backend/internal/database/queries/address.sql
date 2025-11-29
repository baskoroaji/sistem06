
-- name: CreateAddress :one
INSERT INTO address (jalan, rt, rw, kota, postal_code, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id;

SELECT id, jalan, rt, rw, kota, postal_code FROM address
WHERE id= $1