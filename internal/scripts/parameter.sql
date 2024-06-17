-- name: CreateParameter :exec
INSERT INTO parameters ( ID, discount_percentage, created_at, updated_at)
VALUES ($1, $2, $3, $4);

-- name: GetParameters :many
SELECT *
FROM parameters;

-- name: GetParameter :one
SELECT *
FROM parameters
WHERE parameters.id = $1;

-- name: DeleteParameter :exec
DELETE FROM parameters
WHERE parameters.id = $1;
