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

-- name: GetSalesOrderTotalByDay :many

SELECT so.id, 
       so.numero, 
       so.numeroloja, 
       so.data, 
       so.datasaida, 
       so.dataprevista,
       so.totalprodutos,
       so.totaldescontos,
       so.situation_id,
       so.store_id,
       (SELECT SUM(totalprodutos) 
        FROM sales_orders 
        WHERE sales_orders.datasaida >= $1) AS total_produtos_soma
FROM sales_orders so
WHERE so.datasaida = $1;


-- name: GetTotalSalesOrderTotalByWeek :many
SELECT so.id, 
        so.numero, 
        so.numeroloja, 
        so.data, 
        so.datasaida, 
        so.dataprevista,
        so.totalprodutos,
        so.totaldescontos,
        so.situation_id,
        so.store_id,
       (SELECT SUM(totalprodutos) 
        FROM sales_orders 
        WHERE DATE_TRUNC('week', datasaida) = DATE_TRUNC('week', $1::date)) AS total_produtos_soma
FROM sales_orders so
WHERE DATE_TRUNC('week', so.datasaida) = DATE_TRUNC('week', $1::date)
ORDER BY so.datasaida;

-- name: GetTotalSalesOrderLastThirtyDays :many
SELECT
    DATE_TRUNC('day', datasaida) AS dia,
    SUM(totalprodutos) AS total_produtos_soma,
    COUNT(*) AS total_vendas
FROM 
    sales_orders
WHERE 
    datasaida >= NOW() - INTERVAL '30 days' AND datasaida < NOW() 
GROUP BY 
    DATE_TRUNC('day', datasaida)
ORDER BY 
    dia;

