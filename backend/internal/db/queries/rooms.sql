-- name: InsertRoom :one
INSERT INTO rooms (owner_id, name, description, join_code)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING id, name, description, join_code, created_at, updated_at;

-- name: CheckNewJoinCode :many
SELECT id
FROM rooms
WHERE join_code = $1;

-- name: GetRoomsByUser :many
SELECT id, name, description, join_code, created_at, updated_at
FROM rooms
WHERE owner_id = $1;

-- name: UpdateRoom :one
UPDATE rooms
SET 
    name = COALESCE(sqlc.narg(name), name),
    description = COALESCE(sqlc.narg(description), description),
    updated_at = NOW()
WHERE owner_id = $1
    AND join_code = $2
RETURNING id, name, description, join_code, created_at, updated_at;

-- name: DeleteRoom :exec
DELETE FROM rooms
WHERE owner_id = $1
    AND join_code = $2;