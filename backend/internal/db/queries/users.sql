-- name: InsertUser :one
INSERT INTO users (email, username, password_hash)
VALUES (
    $1,
    $2,
    $3
)
RETURNING id, email, username;