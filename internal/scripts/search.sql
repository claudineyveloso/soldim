-- name: CreateSearch :exec
INSERT INTO searches ( ID, description, created_at, updated_at)
VALUES ($1, $2, $3, $4);

-- name: GetSearches :many
SELECT *
FROM searches;

-- name: GetLastSearch :one
SELECT id, description, created_at, updated_At
FROM searches
ORDER BY created_at DESC
LIMIT 1;


-- name: GetSearch :one
SELECT *
FROM searches
WHERE searches.id = $1;

-- name: UpdateSearch :exec
UPDATE searches SET description = $2, updated_at = $3 WHERE searches.id = $1;

-- name: DeleteSearch :exec
DELETE FROM searches
WHERE searches.id = $1;
