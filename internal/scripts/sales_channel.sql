-- name: CreateSalesChannel :exec
INSERT INTO sales_channel (ID, descricao, tipo, situacao, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetSalesChannel :many
SELECT *
FROM sales_channel;

