-- name: CreateTask :one
INSERT INTO tasks (
  id,
  creator_id,
  name,
  deadline
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetTasks :many
SELECT * FROM tasks
WHERE 
  creator_id = $1
  AND (name ILIKE '%' || COALESCE(sqlc.arg('name'), '') || '%')
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
  name = COALESCE(sqlc.narg(name), name),
  deadline = COALESCE(sqlc.narg(deadline), deadline),
  completed = COALESCE(sqlc.narg(completed), completed)
WHERE
  id = $1
RETURNING *;