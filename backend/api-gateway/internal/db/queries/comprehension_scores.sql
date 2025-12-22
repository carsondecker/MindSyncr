-- name: InsertComprehensionScore :one
INSERT INTO comprehension_scores (user_id, session_id, score)
VALUES ($1, $2, $3)
RETURNING id, session_id, user_id, score, created_at;

-- name: GetComprehensionScoresBySession :many
SELECT id, session_id, user_id, score, created_at
FROM comprehension_scores
WHERE session_id = $1;