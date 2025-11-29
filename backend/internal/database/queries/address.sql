
-- name: CreateAddress :one
INSERT INTO address (jalan, rt, rw, kota, postal_code)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: FindAdressByID :one
SELECT id, jalan, rt, rw, kota, postal_code FROM address WHERE id = $1;

-- name: UpdateAddress :exec
UPDATE address 
SET
jalan = COALESCE(sqlc.narg('jalan'), jalan),
rt = COALESCE(sqlc.narg('rt'), rt),
rw = COALESCE(sqlc.narg('rw'), rw),
kota = COALESCE(sqlc.narg('kota'), kota),
postal_code = COALESCE(sqlc.narg('postal_code'), postal_code)
WHERE id = $1 ;
