-- name: CreateSearchResult :exec
INSERT INTO searches_result ( ID, image_url, description, source, price, promotion, link, search_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);

-- name: GetSearchesResult :many
SELECT id, 
        image_url,
        description,
        source, price,
        promotion,
        link,
        search_id,
        created_at,
        updated_at
FROM searches_result ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetSearchResult :one
SELECT id,
        image_url,
        description,
        source, price,
        promotion,
        link,
        search_id,
        created_at,
        updated_at
FROM searches_result
WHERE searches_result.id = $1;

-- name: DeleteSearchResult :exec
DELETE FROM searches_result
WHERE searches_result.id = $1;


-- name: GetTotalSearchesResult :one
SELECT COUNT(*)
FROM searches_result;


