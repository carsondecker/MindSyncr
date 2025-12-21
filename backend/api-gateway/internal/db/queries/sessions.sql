-- name: InsertSession :one
INSERT INTO sessions (owner_id, room_id, name) 
SELECT $1, r.id, $3
FROM rooms r
WHERE room_id = $2
RETURNING id, room_id, owner_id, name, is_active, started_at, ended_at, created_at, updated_at;

-- name: GetSessionsByRoomId :many
SELECT id, room_id, owner_id, name, is_active, started_at, ended_at, created_at, updated_at
FROM sessions
WHERE room_id = $1;

-- name: GetSessionById :one
SELECT id, room_id, owner_id, name, is_active, started_at, ended_at, created_at, updated_at
FROM sessions
WHERE id = $1;

-- name: EndSession :exec
UPDATE sessions
SET is_active = FALSE,
    ended_at = NOW(),
    updated_at = NOW()
WHERE owner_id = $1
    AND id = $2;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE owner_id = $1
    AND id = $2;