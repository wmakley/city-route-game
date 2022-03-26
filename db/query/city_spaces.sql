-- name: ListCitySpacesByCityID :many
SELECT * FROM city_spaces
WHERE city_id = $1
ORDER BY "order";

-- name: GetCitySpaceByID :one
SELECT * FROM city_spaces
WHERE id = $1
LIMIT 1;

-- name: ListCitySpacesByMultipleCities :many
SELECT * FROM city_spaces
WHERE city_id = ANY($1::BIGINT[])
ORDER BY city_id, "order";

-- name: CreateCitySpace :one
INSERT INTO city_spaces (city_id, "order", space_type, required_privilege)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateCitySpace :one
UPDATE city_spaces SET "order" = $1, space_type = $2, required_privilege = $3, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteCitySpace :exec
DELETE FROM city_spaces WHERE id = $1;

-- name: DeleteCitySpacesWhereCityIDIn :exec
DELETE FROM city_spaces WHERE city_id = ANY(@cityIDs::bigint[]);
