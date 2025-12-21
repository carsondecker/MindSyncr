-- name: InsertSession :one
INSERT INTO sessions (owner_id, room_id, name) 
SELECT $1, r.id, $3
FROM rooms r
WHERE room_id = $2
RETURNING id, room_id, owner_id, name, is_active, started_at, ended_at, created_at, updated_at;

-- name: GetSessionsForOwner :many
SELECT s.id, s.room_id, s.owner_id, s.name, s.is_active, s.started_at, s.ended_at, s.created_at, s.updated_at
FROM sessions s
JOIN rooms r
    ON s.room_id = r.id
WHERE s.owner_id = $1
    AND r.join_code = $2;

-- name: EndSession :exec
UPDATE sessions
SET is_active = FALSE,
    ended_at = NOW(),
    updated_at = NOW();

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE id = $1;