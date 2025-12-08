
-- name: GetPermissionsByRoleID :many
SELECT p.id, p.name
FROM permissions p
JOIN roles_permissions rp ON rp.permission_id = p.id
WHERE rp.role_id = $1;

-- name: GetPermissionsByUserID :many
SELECT DISTINCT p.id, p.name
FROM permissions p
JOIN roles_permissions rp ON rp.permission_id = p.id
JOIN user_roles ur ON ur.roles_id = rp.role_id
WHERE ur.user_id = $1;