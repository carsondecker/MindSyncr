-- name: InsertRefreshToken :one
INSERT INTO refresh_tokens (user_id, token, expires_at)
VALUES (
    $1,
    $2,
    $3
)
RETURNING token, expires_at, created_at;