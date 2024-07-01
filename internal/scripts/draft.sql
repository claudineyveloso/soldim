-- name: CreateDraft :exec
INSERT INTO drafts ( ID, image_url, description, source, price, promotion, link, search_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);

-- name: GetDrafts :many
SELECT *
FROM drafts ORDER BY created_at DESC;

-- name: GetDraftByDescription :many
SELECT *
FROM drafts
WHERE drafts.description LIKE '%' || $1 || '%';

-- name: GetDraft :one
SELECT *
FROM drafts
WHERE drafts.id = $1;

-- name: UpdateDraft :exec
UPDATE drafts SET description = $2, 
  image_url = $3,
  source = $4,
  price  = $5,
  promotion = $6,
  link = $7,
  updated_at = $8
WHERE drafts.id = $1;

-- name: DeleteDraft :exec
DELETE FROM drafts
WHERE drafts.id = $1;
