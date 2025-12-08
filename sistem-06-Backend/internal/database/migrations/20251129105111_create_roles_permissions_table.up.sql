CREATE TABLE IF NOT EXISTS roles_permissions (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    role_id BIGINT NOT NULL,
    permission_id BIGINT NOT NULL,
    
    CONSTRAINT fk_roles
        FOREIGN KEY (role_id) REFERENCES roles(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_permission
        FOREIGN KEY (permission_id) REFERENCES permissions(id)
        ON DELETE CASCADE
);

ALTER TABLE roles_permissions
    ADD CONSTRAINT unique_role_permission UNIQUE (role_id, permission_id);
