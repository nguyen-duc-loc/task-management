-- name: CreateTask :one
INSERT INTO tasks (
  id,
  creator,
  name,
  deadline
) VALUES (
  $1, $2, $3, $4
) RETURNING *;
