

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

-- name: FindUserRolesWithPermissions :many
SELECT 
r.id AS role_id,
r.name AS role_name,
p.name AS permission_name
FROM user_roles ur
JOIN roles r ON r.id = ur.roles_id
LEFT JOIN roles_permissions rp ON rp.role_id = r.id
LEFT JOIN permissions p ON p.id = rp.permission_id
WHERE ur.user_id = $1
ORDER BY r.id, p.name;
