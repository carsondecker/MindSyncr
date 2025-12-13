-- name: InsertRoom :exec
INSERT INTO rooms (owner_id, name, description, join_code) {
    $1,
    $2,
    $3,
    $4
}
RETURNING id, name, description, join_code, created_at, updated_at;