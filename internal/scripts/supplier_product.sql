-- name: CreateSupplierProduct :exec
INSERT INTO supplier_products (id, descricao, preco_custo, preco_compra, padrao, supplier_id, product_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);

-- name: GetSupplierProduct :many
SELECT  id,
          descricao,
          preco_custo,
          preco_compra,
          padrao,
          supplier_id,
          product_id,
          created_at,
          updated_at
FROM supplier_products;

-- name: UpdateSupplierProduct :exec
UPDATE supplier_products SET descricao = $2,
  preco_custo = $3,
  preco_compra = $4,
  padrao = $5,
  supplier_id = $6,
  product_id = $7,
  updated_at = $8
WHERE supplier_products.id = $1;


