-- name: CreateProduct :exec
INSERT INTO products ( ID, image_url, name, source, price, promotion, link, search_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);

-- name: GetProducts :many
SELECT *
FROM products;

-- name: GetProduct :one
SELECT *
FROM produts
WHERE products.id = $1;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE products.id = $1;
