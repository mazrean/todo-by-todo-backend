-- name: CreateUser :execresult
INSERT INTO users (name) VALUES (?);

-- name: GetUser :one
SELECT id, name, created_at, updated_at FROM users WHERE id = ?;

-- name: UpdateUser :exec
UPDATE users SET name = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = ?;

-- name: ListUsers :many
SELECT id, name, created_at, updated_at FROM users ORDER BY created_at DESC;