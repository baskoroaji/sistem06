CREATE TABLE IF NOT EXISTS roles_permissions(
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    role_id INT NOT NULL,
    permission_id INT NOT NULL
    Foreign Key role_id(roles) REFERENCES ()
)