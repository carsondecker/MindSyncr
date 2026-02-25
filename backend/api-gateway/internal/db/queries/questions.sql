-- name: InsertQuestion :one
INSERT INTO questions (user_id, session_id, text)
VALUES (
    $1,
    $2,
    $3
)
RETURNING id, user_id, session_id, text, is_answered, answered_at, created_at, updated_at;

-- name: GetQuestionsBySession :many
SELECT id, user_id, session_id, text, is_answered, answered_at, created_at, updated_at
FROM questions
WHERE session_id = $1;

-- name: DeleteQuestion :exec
DELETE FROM questions
WHERE user_id = $1
    AND id = $2
    AND session_id = $3;