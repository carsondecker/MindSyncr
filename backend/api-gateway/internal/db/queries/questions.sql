-- name: InsertQuestion :one
INSERT INTO questions (user_id, session_id, text)
VALUES (
    $1,
    $2,
    $3
)
RETURNING id, user_id, session_id, text, is_answered, answered_at, created_at, updated_at;