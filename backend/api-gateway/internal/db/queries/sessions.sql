-- name: InsertSession :one
INSERT INTO sessions (owner_id, room_id, name) 
VALUES ($1, $2, $3)
RETURNING id, room_id, owner_id, name, is_active, started_at, ended_at, created_at, updated_at;

-- name: GetSessionsByRoomId :many
SELECT s.id, s.room_id, s.owner_id, s.name, s.is_active, s.started_at, s.ended_at, s.created_at, s.updated_at,
    (s.owner_id = $2) as is_owner,
    (sm.user_id IS NOT NULL)::boolean as is_member
FROM sessions s
LEFT JOIN session_memberships sm
    ON s.id = sm.session_id
    AND sm.user_id = $2
WHERE s.room_id = $1
ORDER BY s.created_at DESC;

-- name: GetSessionById :one
SELECT s.id, s.room_id, s.owner_id, s.name, s.is_active, s.started_at, s.ended_at, s.created_at, s.updated_at,
    (s.owner_id = $2) as is_owner,
    (sm.user_id IS NOT NULL)::boolean as is_member
FROM sessions s
LEFT JOIN session_memberships sm
    ON s.id = sm.session_id
    AND sm.user_id = $2
WHERE id = $1;

-- name: GetSessionOwnerById :one
SELECT owner_id
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