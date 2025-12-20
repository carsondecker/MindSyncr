-- name: InsertRefreshToken :one
INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
VALUES (
    $1,
    $2,
    $3
)
RETURNING expires_at, created_at;

-- name: RevokeUserTokens :exec
UPDATE refresh_tokens
SET is_revoked = TRUE
WHERE user_id = $1;

-- name: CheckValidRefreshToken :one
SELECT user_id
FROM refresh_tokens
WHERE token_hash = $1
    AND is_revoked = FALSE
    AND expires_at > NOW();