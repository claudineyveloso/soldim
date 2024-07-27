-- name: CreateSearchResult :exec
INSERT INTO searches_result ( ID, image_url, description, source, price, promotion, link, search_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);

-- name: GetSearchesResult :many
WITH LatestSearchID AS (
    SELECT search_id
    FROM searches_result
    ORDER BY created_at DESC
    LIMIT 1
)
SELECT id, 
       image_url,
       description,
       source,
       price,
       promotion,
       link,
       search_id,
       created_at,
       updated_at
FROM searches_result
WHERE search_id = (SELECT search_id FROM LatestSearchID)
AND ($1::text IS NULL OR $1 = '' OR source = $1)
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetSearchResult :one
SELECT id,
        image_url,
        description,
        source,
        price,
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
WITH LatestSearchID AS (
    SELECT search_id
    FROM searches_result
    ORDER BY created_at DESC
    LIMIT 1
)
SELECT COUNT(*)
FROM searches_result
WHERE search_id = (SELECT search_id FROM LatestSearchID);

-- name: GetSearchResultSources :many
WITH RankedResults AS (
    SELECT source,
           search_id,
           ROW_NUMBER() OVER (PARTITION BY source ORDER BY source DESC) AS rn
    FROM searches_result
    WHERE search_id = $1
)
SELECT source,
       search_id
FROM RankedResults
WHERE rn = 1
ORDER BY source ASC;
