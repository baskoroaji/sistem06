-- name: GetRolesByUserID :many
SELECT r.id, r.name
FROM roles r
JOIN user_roles ur ON ur.roles_id = r.id
WHERE ur.user_id = $1;

-- name: GetRolesWithPermissionsByUserID :many
SELECT 
    r.id AS role_id,
    r.name AS role_name,
    p.id AS permission_id,
    p.name AS permission_name
FROM user_roles ur
JOIN roles r ON r.id = ur.roles_id
LEFT JOIN roles_permissions rp ON rp.role_id = r.id
LEFT JOIN permissions p ON p.id = rp.permission_id
WHERE ur.user_id = $1
ORDER BY r.id, p.name;

-- name: AssignRoleToUser :exec
INSERT INTO user_roles (user_id, roles_id)
VALUES ($1, $2);

-- name: RemoveRoleFromUser :exec
DELETE FROM user_roles
WHERE user_id = $1 AND roles_id = $2;
