-- name: InsertQuestionLike :one
INSERT INTO question_likes (user_id, question_id)
VALUES (
    $1,
    $2
)
RETURNING id, user_id, question_id, created_at;

-- name: GetQuestionLikesBySession :many
SELECT ql.id, ql.user_id, ql.question_id, ql.created_at
FROM question_likes ql
JOIN questions q
    ON ql.question_id = q.id
    WHERE session_id = $1;

-- name: DeleteQuestionLike :one
DELETE FROM question_likes ql
USING questions q
WHERE ql.question_id = q.id
    AND ql.user_id = $1
    AND ql.question_id = $2
    AND q.session_id = $3
RETURNING ql.id;

-- name: CheckCanDeleteQuestionLike :one
SELECT EXISTS (
    SELECT 1
    FROM question_likes
    WHERE user_id = $1
);