-- name: InsertUser :one
INSERT INTO users (email, username, password_hash)
VALUES (
    $1,
    $2,
    $3
)
RETURNING id, email, username, role, is_email_verified, status, created_at, updated_at;

-- name: GetUserForLogin :one
SELECT id, email, username, role, is_email_verified, status, created_at, updated_at, password_hash
FROM users
WHERE email = $1;

-- name: GetUserById :one
SELECT id, email, username, role, is_email_verified, status, created_at, updated_at
FROM users
WHERE id = $1;