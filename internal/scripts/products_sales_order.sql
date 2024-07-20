-- name: CreateProductSalesOrder :exec
INSERT INTO products_sales_orders (sales_order_id, product_id, quantidade, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5);

-- name: GetProductSalesOrder :many
SELECT pso.sales_order_id, 
        pso.product_id, 
        pso.quantidade, 
        pso.created_at, 
        pso.updated_at,
        so.id,
        so.numero,
        so.numeroloja,
        so.data,
        so.datasaida,
        so.dataprevista,
        so.totalprodutos,
        so.totaldescontos,
        so.situation_id,
        so.store_id,
        sp.supplier_id,
        p.nome,
        p.codigo,
        p.preco
FROM products_sales_orders pso
LEFT JOIN products p ON p.id = pso.product_id
LEFT JOIN sales_orders so ON so.id = pso.sales_order_id
LEFT JOIN supplier_products sp ON pso.product_id = sp.product_id
ORDER BY pso.product_id, pso.sales_order_id;

-- name: GetProductSalesOrderBySupplierID :many
SELECT pso.sales_order_id, 
        pso.product_id, 
        pso.quantidade, 
        pso.created_at, 
        pso.updated_at,
        so.id,
        so.numero,
        so.numeroloja,
        so.data,
        so.datasaida,
        so.dataprevista,
        so.totalprodutos,
        so.totaldescontos,
        so.situation_id,
        so.store_id,
        sp.supplier_id,
        p.nome,
        p.codigo,
        p.preco
FROM products_sales_orders pso
LEFT JOIN products p ON p.id = pso.product_id
LEFT JOIN sales_orders so ON so.id = pso.sales_order_id
LEFT JOIN supplier_products sp ON pso.product_id = sp.product_id
WHERE sp.supplier_id = $1
ORDER BY pso.product_id, pso.sales_order_id;

