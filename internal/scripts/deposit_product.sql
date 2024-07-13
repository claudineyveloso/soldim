-- name: CreateDepositProduct :exec
INSERT INTO deposit_products (ID, deposit_id, product_id, saldoFisico, saldoVirtual, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: GetDepositProducts :many
SELECT ID,
        deposit_id,
        product_id,
        saldoFisico,
        saldoVirtual,
        created_at,
        updated_at
FROM deposit_products;

-- name: UpdateDepositProduct :exec
UPDATE deposit_products SET saldoFisico = $3, 
  saldoVirtual = $4, 
  updated_at = $5
WHERE deposit_products.deposit_id = $1 AND deposit_products.product_id = $2;


