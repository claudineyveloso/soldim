-- name: CreateItemsSalesOrder :exec
INSERT INTO items_sales_orders (id, sales_order_id, codigo, unidade, quantidade, desconto, valor, aliquotaIPI, descricao, descricaoDetalhada, product_id, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);

