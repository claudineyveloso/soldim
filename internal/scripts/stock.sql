-- name: CreateStock :exec
INSERT INTO stocks (ID, product_id, saldoFisicoTotal, saldoVirtualTotal, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetStocks :many
SELECT *
FROM stocks;

-- name: UpdateStock :exec
UPDATE stocks SET saldoFisicoTotal = $2, 
  saldoVirtualTotal = $3, 
  updated_at = $4
WHERE stocks.id = $1;


