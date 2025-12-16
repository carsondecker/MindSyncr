-- name: JoinRoom :exec
INSERT INTO room_memberships (user_id, room_id)
SELECT $1, id
FROM rooms
WHERE join_code = $2;

-- name: LeaveRoom :exec
DELETE FROM room_memberships rm
USING rooms r
WHERE rm.user_id = $1
  AND r.join_code = $2;