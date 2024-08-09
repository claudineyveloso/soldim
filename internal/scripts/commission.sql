-- name: CreateCommission :exec
INSERT INTO commissions (items_sales_order_id, base, aliquota, valor, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6);

