-- name: JoinSession :exec
INSERT INTO session_memberships (user_id, session_id)
VALUES ($1, $2);

-- name: LeaveSession :exec
DELETE FROM session_memberships
WHERE user_id = $1
  AND session_id = $2;