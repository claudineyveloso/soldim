-- name: CreateStock :exec
INSERT INTO stocks (product_id, saldo_fisico_total, saldo_virtual_total, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5);

-- name: GetStocks :many
SELECT product_id,
        saldo_fisico_total,
        saldo_virtual_total,
        created_at,
        updated_at
FROM stocks;

-- name: UpdateStock :exec
UPDATE stocks SET saldo_fisico_total = $2, 
  saldo_virtual_total = $3, 
  updated_at = $4
WHERE stocks.product_id = $1;


