-- name: CreateUser :exec
INSERT INTO users ( ID, email, password, is_active, user_type, session_version, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: GetUser :one
SELECT ID,
        email,
        password,
        is_active,
        user_type,
        session_version,
        created_at,
        updated_at
FROM users
WHERE users.id = $1;

-- name: GetUsers :many
SELECT ID,
        email,
        password,
        is_active,
        user_type,
        session_version,
        created_at,
        updated_at
FROM users
ORDER BY users.email ASC;

-- name: GetUserByEmail :one
SELECT ID,
        email,
        password,
        is_active,
        user_type,
        session_version,
        created_at,
        updated_at
FROM users
WHERE users.email = $1;

-- name: LoginUser :one
SELECT ID,
        email,
        password,
        is_active,
        user_type,
        session_version,
        created_at,
        updated_at
FROM users
WHERE users.email = $1 AND users.password = $2;

-- name: DisableUser :exec
UPDATE users SET is_active = $2, updated_at = $3 WHERE users.id = $1;

-- name: UpdateUser :exec
UPDATE users SET email = $2, is_active = $3, user_type = $4, updated_at = $5 WHERE users.id = $1;

-- name: UpdatePassword :exec
UPDATE users SET password = $2, updated_at = $3 WHERE users.id = $1;

-- name: UpdateSessionVersion :exec
UPDATE users SET session_version = $2, updated_at = $3 WHERE users.id = $1;
