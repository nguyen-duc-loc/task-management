-- name: CreateTask :one
INSERT INTO tasks (
  id,
  creator_id,
  title,
  description,
  deadline
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetTasks :many
SELECT 
  *,
  COUNT(*) OVER() AS total
FROM tasks
WHERE 
  creator_id = $1
  AND (
    title ILIKE '%' || COALESCE(sqlc.arg('title'), '') || '%'
    OR 
    description ILIKE '%' || COALESCE(sqlc.arg('description'), '') || '%'
  )
  AND (
    sqlc.narg('start_deadline')::timestamptz IS NULL
    OR deadline >= sqlc.narg('start_deadline')
  )
  AND (
    sqlc.narg('end_deadline')::timestamptz IS NULL
    OR deadline <= sqlc.narg('end_deadline')
  )
  AND (
    sqlc.narg('completed')::bool IS NULL 
    OR completed = sqlc.narg('completed')::bool
  )
  ORDER BY completed ASC, deadline ASC
  LIMIT $2 OFFSET $3;

-- name: GetTaskByID :one
SELECT * FROM tasks
WHERE id = $1 LIMIT 1;

-- name: UpdateTask :one
UPDATE tasks
SET
  title = COALESCE(sqlc.narg(title), title),
  description = COALESCE(sqlc.narg(description), description),
  deadline = COALESCE(sqlc.narg(deadline), deadline),
  completed = COALESCE(sqlc.narg(completed), completed)
WHERE
  id = $1
RETURNING *;

-- name: DeleteTask :exec
DELETE FROM tasks
WHERE id = $1;