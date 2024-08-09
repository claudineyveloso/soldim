-- name: CreateItemsSalesOrder :exec
INSERT INTO items_sales_orders (id, sales_order_id, codigo, unidade, quantidade, desconto, valor, aliquotaIPI, descricao, descricaoDetalhada, product_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13);

-- name: GetItemsSalesOrderByID :one
SELECT item.id,
       item.sales_order_id,
       item.codigo,
       item.unidade,
       item.quantidade,
       item.desconto,
       item.valor,
       item.aliquotaIPI,
       item.descricao,
       item.descricaoDetalhada,
       item.product_id,
       item.created_at,
       item.updated_at
FROM items_sales_orders item WHERE item.id = $1;

-- name: GetItemsSalesOrderByProductID :one
SELECT item.id,
       item.sales_order_id,
       item.codigo,
       item.unidade,
       item.quantidade,
       item.desconto,
       item.valor,
       item.aliquotaIPI,
       item.descricao,
       item.descricaoDetalhada,
       item.product_id,
       item.created_at,
       item.updated_at
FROM items_sales_orders item WHERE item.product_id = $1;

-- name: GetItemsSalesOrderBySalesOrderID :one
SELECT item.id,
       item.sales_order_id,
       item.codigo,
       item.unidade,
       item.quantidade,
       item.desconto,
       item.valor,
       item.aliquotaIPI,
       item.descricao,
       item.descricaoDetalhada,
       item.product_id,
       item.created_at,
       item.updated_at
FROM items_sales_orders item WHERE item.sales_order_id = $1;


