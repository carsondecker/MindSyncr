-- name: InsertRoom :one
INSERT INTO rooms (owner_id, name, description, join_code)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING id, owner_id, name, description, join_code, created_at, updated_at;

-- name: CheckNewJoinCode :many
SELECT id
FROM rooms
WHERE join_code = $1;

-- name: CheckRoomMembership :one
SELECT 1
FROM rooms r
LEFT JOIN room_memberships rm
    ON rm.room_id = r.id
    AND rm.user_id = $2
WHERE r.id = $1
    AND (r.owner_id = $2 OR rm.user_id IS NOT NULL)
LIMIT 1;

-- name: CheckRoomOwnership :one
SELECT 1
FROM rooms
WHERE id = $1
    AND owner_id = $2
LIMIT 1;

-- name: GetRoomById :one
SELECT id, owner_id, name, description, join_code, created_at, updated_at
FROM rooms
WHERE id = $1;

-- name: GetRoomsByOwner :many
SELECT id, owner_id, name, description, join_code, created_at, updated_at
FROM rooms
WHERE owner_id = $1;

-- name: GetRoomsByMembership :many
SELECT r.id, owner_id, r.name, r.description, r.join_code, r.created_at, r.updated_at
FROM rooms r
JOIN room_memberships rm
    ON rm.room_id = r.id
WHERE rm.user_id = $1;

-- name: GetRoomOwnerIdByJoinCode :one
SELECT owner_id
FROM rooms
WHERE join_code = $1;

-- name: UpdateRoom :one
UPDATE rooms
SET name = COALESCE(sqlc.narg(name), name),
    description = COALESCE(sqlc.narg(description), description),
    updated_at = NOW()
WHERE owner_id = $1
    AND id = $2
RETURNING id, owner_id, name, description, join_code, created_at, updated_at;

-- name: DeleteRoom :exec
DELETE FROM rooms
WHERE owner_id = $1
    AND id = $2;