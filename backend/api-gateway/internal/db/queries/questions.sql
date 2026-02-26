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

-- name: CheckQuestionBelongsToSession :one
SELECT 1
FROM questions
WHERE id = $1
    AND session_id = $2
LIMIT 1;

-- name: CheckCanDeleteQuestion :one
SELECT 1
FROM questions q
JOIN sessions s
    ON q.session_id = s.id
WHERE q.id = $1
    AND (
        q.user_id = $2 OR
        s.owner_id = $2
    );