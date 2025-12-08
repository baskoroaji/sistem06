

-- name: CreateUser :one
INSERT INTO users (name, email, password, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;


-- name: CountUserByID :one
SELECT COUNT(*) FROM users WHERE id = $1;

-- name: CountUserByName :one
SELECT COUNT(*) FROM users WHERE name = $1;


-- name: FindUserByEmail :one
SELECT id, name, email, password, created_at, updated_at 
FROM users WHERE email = $1;

-- name: FindUserByID :one
SELECT id, name, email, password, created_at, updated_at
FROM users
WHERE id = $1;
