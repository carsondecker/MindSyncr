-- name: InsertUser :one
INSERT INTO users (email, username, password_hash)
VALUES (
    $1,
    $2,
    $3
)
RETURNING id, email, username, role, created_at;

-- name: GetUserForLogin :one
SELECT id, email, username, role, password_hash
FROM users
WHERE email = $1;

-- name: GetUserById :one
SELECT id, email, username, role
FROM users
WHERE id = $1;