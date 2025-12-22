-- name: InsertSession :one
INSERT INTO sessions (owner_id, room_id, name) 
VALUES ($1, $2, $3)
RETURNING id, room_id, owner_id, name, is_active, started_at, ended_at, created_at, updated_at;

-- name: GetSessionsByRoomId :many
SELECT id, room_id, owner_id, name, is_active, started_at, ended_at, created_at, updated_at
FROM sessions
WHERE room_id = $1
ORDER BY created_at DESC;

-- name: GetSessionById :one
SELECT id, room_id, owner_id, name, is_active, started_at, ended_at, created_at, updated_at
FROM sessions
WHERE id = $1;

-- name: CheckRoomMembershipBySessionId :one
SELECT 1
FROM rooms r
LEFT JOIN room_memberships rm
    ON r.id = rm.room_id
    AND rm.user_id = $2
JOIN sessions s
    ON r.id = s.room_id
WHERE s.id = $1
    AND (r.owner_id = $2 OR rm.user_id IS NOT NULL)
LIMIT 1;

-- name: CheckRoomOwnershipBySessionId :one
SELECT 1
FROM rooms r
JOIN sessions s
    ON r.id = s.room_id
WHERE s.id = $1
    AND r.owner_id = $2
LIMIT 1;

-- name: CheckSessionMembershipOnly :one
SELECT 1
FROM sessions s
LEFT JOIN session_memberships sm
    ON s.id = sm.session_id
    AND sm.user_id = $2
WHERE s.id = $1
    AND sm.user_id IS NOT NULL
LIMIT 1;

-- name: CheckSessionActive :one
SELECT 1
FROM sessions
WHERE id = $1
    AND is_active = TRUE;

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