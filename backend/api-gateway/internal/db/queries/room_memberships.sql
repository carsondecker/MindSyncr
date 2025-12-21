-- name: JoinRoom :one
INSERT INTO room_memberships (user_id, room_id)
SELECT $1, id
FROM rooms
WHERE join_code = $2
RETURNING room_id;

-- name: LeaveRoom :exec
DELETE FROM room_memberships rm
WHERE user_id = $1
  AND room_id = $2;