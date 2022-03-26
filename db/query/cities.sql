-- name: GetCity :one
SELECT * FROM cities
WHERE id = $1;

-- name: ListCitiesByBoardID :many
SELECT * FROM cities
WHERE board_id = $1
ORDER BY id;

-- name: ListCityIDsByBoardID :many
SELECT id FROM cities
WHERE board_id = $1;

-- name: CreateCity :one
INSERT INTO cities (board_id, name, x, y, upgrade_offered, immediate_point)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateCity :one
UPDATE cities SET name = $2, x = $3, y = $4, upgrade_offered = $5, immediate_point = $6, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteCity :exec
DELETE FROM cities
WHERE id = $1;

-- name: DeleteMultipleCities :exec
DELETE FROM cities
WHERE id = ANY(@ids::BIGINT[]);
