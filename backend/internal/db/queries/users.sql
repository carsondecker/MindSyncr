-- name: InsertUser :one
INSERT INTO users (email, username, password_hash)
VALUES (
    $1,
    $2,
    $3
)
RETURNING id, email, username;

-- name: GetUserForLogin :one
SELECT id, email, username, password_hash
FROM users
WHERE email = $1;