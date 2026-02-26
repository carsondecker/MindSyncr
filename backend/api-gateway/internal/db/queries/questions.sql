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

-- name: UpdateQuestion :one
UPDATE questions
SET text = COALESCE(sqlc.narg(text), text),
    updated_at = NOW()
WHERE user_id = $1
    AND id = $2
RETURNING id, user_id, session_id, text, is_answered, answered_at, created_at, updated_at;

-- name: CheckQuestionBelongsToSession :one
SELECT EXISTS (
    SELECT 1
    FROM questions
    WHERE id = $1
        AND session_id = $2
);

-- name: CheckCanDeleteQuestion :one
SELECT EXISTS (
    SELECT 1
    FROM questions q
    JOIN sessions s
        ON q.session_id = s.id
    WHERE q.id = $1
        AND (
            q.user_id = $2 OR
            s.owner_id = $2
        )
);

-- name: CheckOwnsQuestion :one
SELECT EXISTS (
    SELECT 1
    FROM questions
    WHERE id = $1
        AND user_id = $2
);