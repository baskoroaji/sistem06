CREATE TABLE IF NOT EXISTS user_roles (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL,
    roles_id BIGINT NOT NULL,

    CONSTRAINT fk_user
        FOREIGN KEY (user_id) REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_roles
        FOREIGN KEY (roles_id) REFERENCES roles(id)
        ON DELETE CASCADE
);

ALTER TABLE user_roles
    ADD CONSTRAINT unique_user_role UNIQUE (user_id, roles_id);