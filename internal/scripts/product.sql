-- name: CreateProduct :exec
INSERT INTO products ( ID, nome, codigo, preco, tipo, situacao, formato, descricaoCurta, imagemURL, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);

-- name: GetProduct :one
SELECT *
FROM products
WHERE products.id = $1;

-- name: GetProducts :many
SELECT *
FROM products;

-- name: GetProductByName :one
SELECT *
FROM products
WHERE products.nome = $1;

