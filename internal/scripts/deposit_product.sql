-- name: CreateDepositProduct :exec
INSERT INTO deposit_products (deposit_id, product_id, saldo_fisico, saldo_virtual, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetDepositProducts :many
SELECT deposit_id,
        product_id,
        saldo_fisico,
        saldo_virtual,
        created_at,
        updated_at
FROM deposit_products;

-- name: UpdateDepositProduct :exec
UPDATE deposit_products SET saldo_fisico = $3, 
  saldo_virtual = $4, 
  updated_at = $5
WHERE deposit_products.deposit_id = $1 AND deposit_products.product_id = $2;


