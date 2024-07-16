-- name: CreateStore :exec
INSERT INTO sales_orders (id, descricao, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5);

-- name: GetStores :many
SELECT id,
        descricao,
        created_at,
        updated_at
FROM stores;

-- name: GetStore :one
SELECT id,
        descricao,
        created_at,
        updated_at
FROM stores
WHERE stores.id = $1;

-- name: GetStoreByDescription :one
SELECT id,
        descricao,
        created_at,
        updated_at
FROM stores
WHERE stores.descricao = $1;


