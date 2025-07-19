-- name: CreateTodo :execresult
INSERT INTO todos (user_id, title, description, completed) VALUES (?, ?, ?, ?);

-- name: GetTodo :one
SELECT id, user_id, title, description, completed, created_at, updated_at FROM todos WHERE id = ?;

-- name: ListTodosByUser :many
SELECT id, user_id, title, description, completed, created_at, updated_at 
FROM todos 
WHERE user_id = ? 
ORDER BY created_at DESC;

-- name: UpdateTodo :exec
UPDATE todos 
SET title = ?, description = ?, completed = ?, updated_at = CURRENT_TIMESTAMP 
WHERE id = ?;

-- name: DeleteTodo :exec
DELETE FROM todos WHERE id = ?;

-- name: ListTodos :many
SELECT id, user_id, title, description, completed, created_at, updated_at 
FROM todos 
ORDER BY created_at DESC;

-- name: MarkTodoCompleted :exec
UPDATE todos SET completed = TRUE, updated_at = CURRENT_TIMESTAMP WHERE id = ?;

-- name: MarkTodoIncomplete :exec
UPDATE todos SET completed = FALSE, updated_at = CURRENT_TIMESTAMP WHERE id = ?;