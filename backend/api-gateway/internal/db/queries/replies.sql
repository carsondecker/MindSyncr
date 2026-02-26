-- name: InsertReply :one
INSERT INTO replies (user_id, question_id, parent_id, text)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING id, user_id, question_id, parent_id, text, created_at, updated_at;

-- name: GetRepliesBySession :many
SELECT r.id, r.user_id, r.question_id, r.parent_id, r.text, r.created_at, r.updated_at
FROM replies r
JOIN questions q
    ON r.question_id = q.id
    WHERE session_id = $1;

-- name: DeleteReply :one
DELETE FROM replies r
USING questions q
WHERE r.question_id = q.id
    AND r.user_id = $1
    AND r.id = $2
    AND q.session_id = $3
RETURNING r.id;

-- name: UpdateReply :one
UPDATE replies
SET text = COALESCE(sqlc.narg(text), text),
    updated_at = NOW()
WHERE user_id = $1
    AND id = $2
RETURNING id, user_id, question_id, parent_id, text, created_at, updated_at;
