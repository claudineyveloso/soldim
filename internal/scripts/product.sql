-- name: CreateProduct :exec
INSERT INTO products ( ID, nome, codigo, preco, tipo, situacao, formato, descricao_curta, unidade, condicao, gtin, imagem_url, data_validade, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15);

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

-- name: UpdateProduct :exec
UPDATE products SET nome = $2, 
  codigo = $3,
  preco = $4,
  tipo  = $5,
  situacao = $6,
  formato = $7,
  descricao_curta = $8,
  unidade = $9,
  condicao = $10,
  gtin = $11,
  imagem_url = $12,
  data_validade = $13,
  updated_at = $14
WHERE products.id = $1;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE products.id = $1;
