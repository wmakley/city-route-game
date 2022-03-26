-- name: GetBoard :one
SELECT * FROM boards
WHERE id = $1 LIMIT 1;

-- name: ListBoards :many
SELECT * FROM boards
ORDER BY id;

-- name: CountBoards :one
SELECT count(*) as count FROM boards;

-- name: CreateBoard :one
INSERT INTO boards (name, width, height)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateBoard :one
UPDATE boards SET name = $2, width = $3, height = $4, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteBoard :exec
DELETE FROM boards
WHERE id = $1;
