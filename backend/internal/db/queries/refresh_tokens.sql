-- name: InsertRefreshToken :one
INSERT INTO refresh_tokens (token, expires_at)
VALUES (
    $1,
    $2
)
RETURNING token, expires_at, created_at;