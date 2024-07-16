-- name: CreateSalesOrder :exec
INSERT INTO sales_orders (id, numero, numeroLoja, data, dataSaida, dataPrevista, totalProdutos, totalDescontos, situation_id, store_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);

-- name: GetSalesOrders :many
SELECT id,
        numero,
        numeroLoja,
        data,
        dataSaida,
        dataPrevista,
        totalProdutos,
        totalDescontos,
        situation_id,
        store_id,
        created_at,
        updated_at
FROM sales_orders;

-- name: GetSalesOrder :one
SELECT id,
        numero,
        numeroLoja,
        data,
        dataSaida,
        dataPrevista,
        totalProdutos,
        totalDescontos,
        situation_id,
        store_id,
        created_at,
        updated_at
FROM sales_orders
WHERE sales_orders.id = $1;

-- name: GetSalesOrderByNumber :one
SELECT id,
        numero,
        numeroLoja,
        data,
        dataSaida,
        dataPrevista,
        totalProdutos,
        totalDescontos,
        situation_id,
        store_id,
        created_at,
        updated_at
FROM sales_orders
WHERE sales_orders.numero = $1;


