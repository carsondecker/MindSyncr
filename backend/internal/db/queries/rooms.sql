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