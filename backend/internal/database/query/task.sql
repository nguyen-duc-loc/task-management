-- name: CreateTask :one
INSERT INTO tasks (
  id,
  creator_id,
  name,
  deadline
) VALUES (
  $1, $2, $3, $4
) RETURNING *;
