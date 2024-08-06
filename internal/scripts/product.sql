-- name: CreateProduct :exec
INSERT INTO products (ID, idProdutoPai, nome, codigo, preco, tipo, situacao, formato, descricao_curta, imagem_url, dataValidade, unidade, pesoLiquido, pesoBruto, volumes, itensPorCaixa, gtin, gtinEmbalagem, tipoProducao, condicao, freteGratis, marca, descricaoComplementar, linkExterno, observacoes, descricaoEmbalagemDiscreta, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28);

-- name: GetProduct :one
WITH aggregated_stocks AS (
    SELECT product_id,
           SUM(saldo_fisico_total) AS saldo_fisico_total,
           SUM(saldo_virtual_total) AS saldo_virtual_total
    FROM stocks
    GROUP BY product_id
),
aggregated_deposit_products AS (
    SELECT product_id,
           SUM(saldo_fisico) AS saldo_fisico,
           SUM(saldo_virtual) AS saldo_virtual
    FROM deposit_products
    GROUP BY product_id
),
aggregated_supplier_products AS (
    SELECT product_id,
           AVG(preco_custo) AS preco_custo,
           AVG(preco_compra) AS preco_compra,
           supplier_id
    FROM supplier_products
    GROUP BY product_id, supplier_id
)
SELECT
    p.ID,
    p.idProdutoPai,
    p.nome,
    p.codigo,
    p.preco,
    p.tipo,
    p.situacao,
    p.formato,
    p.descricao_curta,
    p.imagem_url,
    p.dataValidade,
    p.unidade,
    p.pesoLiquido,
    p.pesoBruto,
    p.volumes,
    p.itensPorCaixa,
    p.gtin,
    p.gtinEmbalagem,
    p.tipoProducao,
    p.condicao,
    p.freteGratis,
    p.marca,
    p.descricaoComplementar,
    p.linkExterno,
    p.observacoes,
    p.descricaoEmbalagemDiscreta,
    p.created_at,
    p.updated_at,
    COALESCE(s.saldo_fisico_total, 0) AS saldo_fisico_total,
    COALESCE(s.saldo_virtual_total, 0) AS saldo_virtual_total,
    COALESCE(dp.saldo_fisico, 0) AS saldo_fisico,
    COALESCE(dp.saldo_virtual, 0) AS saldo_virtual,
    COALESCE(sp.preco_custo, 0) AS preco_custo,
    COALESCE(sp.preco_compra, 0) AS preco_compra,
    sp.supplier_id
FROM
    products p
LEFT JOIN
    aggregated_stocks s
    ON p.id = s.product_id
LEFT JOIN
    aggregated_deposit_products dp
    ON p.id = dp.product_id
LEFT JOIN
    aggregated_supplier_products sp
    ON p.id = sp.product_id
WHERE p.id = $1;

-- name: GetProducts :many
WITH aggregated_stocks AS (
    SELECT product_id,
           SUM(saldo_fisico_total) AS saldo_fisico_total,
           SUM(saldo_virtual_total) AS saldo_virtual_total
    FROM stocks
    GROUP BY product_id
),
aggregated_deposit_products AS (
    SELECT product_id,
           SUM(saldo_fisico) AS saldo_fisico,
           SUM(saldo_virtual) AS saldo_virtual
    FROM deposit_products
    GROUP BY product_id
),
aggregated_supplier_products AS (
    SELECT product_id,
           AVG(preco_custo) AS preco_custo,
           AVG(preco_compra) AS preco_compra,
           supplier_id
    FROM supplier_products
    GROUP BY product_id, supplier_id
)
SELECT
    p.ID,
    p.idProdutoPai,
    p.nome,
    p.codigo,
    p.preco,
    p.tipo,
    p.situacao,
    p.formato,
    p.descricao_curta,
    p.imagem_url,
    p.dataValidade,
    p.unidade,
    p.pesoLiquido,
    p.pesoBruto,
    p.volumes,
    p.itensPorCaixa,
    p.gtin,
    p.gtinEmbalagem,
    p.tipoProducao,
    p.condicao,
    p.freteGratis,
    p.marca,
    p.descricaoComplementar,
    p.linkExterno,
    p.observacoes,
    p.descricaoEmbalagemDiscreta,
    p.created_at,
    p.updated_at,
    COALESCE(s.saldo_fisico_total, 0) AS saldo_fisico_total,
    COALESCE(s.saldo_virtual_total, 0) AS saldo_virtual_total,
    COALESCE(dp.saldo_fisico, 0) AS saldo_fisico,
    COALESCE(dp.saldo_virtual, 0) AS saldo_virtual,
    COALESCE(sp.preco_custo, 0) AS preco_custo,
    COALESCE(sp.preco_compra, 0) AS preco_compra,
    sp.supplier_id
FROM
    products p
LEFT JOIN
    aggregated_stocks s
    ON p.id = s.product_id
LEFT JOIN
    aggregated_deposit_products dp
    ON p.id = dp.product_id
LEFT JOIN
    aggregated_supplier_products sp
    ON p.id = sp.product_id
WHERE
    ($1::text IS NULL OR $1 = '' OR p.nome ILIKE '%' || $1 || '%')
    AND ($2::text IS NULL OR $2 = '' OR p.situacao = $2::text)
    AND ($3::text IS NULL OR $3 = '' OR sp.supplier_id = $3::int)
    ORDER BY p.nome;

-- name: GetProductByName :one
WITH aggregated_stocks AS (
    SELECT product_id,
           SUM(saldo_fisico_total) AS saldo_fisico_total,
           SUM(saldo_virtual_total) AS saldo_virtual_total
    FROM stocks
    GROUP BY product_id
),
aggregated_deposit_products AS (
    SELECT product_id,
           SUM(saldo_fisico) AS saldo_fisico,
           SUM(saldo_virtual) AS saldo_virtual
    FROM deposit_products
    GROUP BY product_id
),
aggregated_supplier_products AS (
    SELECT product_id,
           AVG(preco_custo) AS preco_custo,
           AVG(preco_compra) AS preco_compra,
           supplier_id
    FROM supplier_products
    GROUP BY product_id, supplier_id
)
SELECT
    p.ID,
    p.idProdutoPai,
    p.nome,
    p.codigo,
    p.preco,
    p.tipo,
    p.situacao,
    p.formato,
    p.descricao_curta,
    p.imagem_url,
    p.dataValidade,
    p.unidade,
    p.pesoLiquido,
    p.pesoBruto,
    p.volumes,
    p.itensPorCaixa,
    p.gtin,
    p.gtinEmbalagem,
    p.tipoProducao,
    p.condicao,
    p.freteGratis,
    p.marca,
    p.descricaoComplementar,
    p.linkExterno,
    p.observacoes,
    p.descricaoEmbalagemDiscreta,
    p.created_at,
    p.updated_at,
    COALESCE(s.saldo_fisico_total, 0) AS saldo_fisico_total,
    COALESCE(s.saldo_virtual_total, 0) AS saldo_virtual_total,
    COALESCE(dp.saldo_fisico, 0) AS saldo_fisico,
    COALESCE(dp.saldo_virtual, 0) AS saldo_virtual,
    COALESCE(sp.preco_custo, 0) AS preco_custo,
    COALESCE(sp.preco_compra, 0) AS preco_compra,
    sp.supplier_id
FROM
    products p
LEFT JOIN
    aggregated_stocks s
    ON p.id = s.product_id
LEFT JOIN
    aggregated_deposit_products dp
    ON p.id = dp.product_id
LEFT JOIN
    aggregated_supplier_products sp
    ON p.id = sp.product_id
WHERE p.nome = $1;

-- name: GetProductEmptyStock :many
WITH aggregated_stocks AS (
    SELECT product_id,
           SUM(saldo_fisico_total) AS saldo_fisico_total,
           SUM(saldo_virtual_total) AS saldo_virtual_total
    FROM stocks
    GROUP BY product_id
),
aggregated_deposit_products AS (
    SELECT product_id,
           SUM(saldo_fisico) AS saldo_fisico,
           SUM(saldo_virtual) AS saldo_virtual
    FROM deposit_products
    GROUP BY product_id
),
aggregated_supplier_products AS (
    SELECT product_id,
           AVG(preco_custo) AS preco_custo,
           AVG(preco_compra) AS preco_compra,
           supplier_id
    FROM supplier_products
    GROUP BY product_id, supplier_id
)
SELECT
    p.ID,
    p.idProdutoPai,
    p.nome,
    p.codigo,
    p.preco,
    p.tipo,
    p.situacao,
    p.formato,
    p.descricao_curta,
    p.imagem_url,
    p.dataValidade,
    p.unidade,
    p.pesoLiquido,
    p.pesoBruto,
    p.volumes,
    p.itensPorCaixa,
    p.gtin,
    p.gtinEmbalagem,
    p.tipoProducao,
    p.condicao,
    p.freteGratis,
    p.marca,
    p.descricaoComplementar,
    p.linkExterno,
    p.observacoes,
    p.descricaoEmbalagemDiscreta,
    p.created_at,
    p.updated_at,
    COALESCE(s.saldo_fisico_total, 0) AS saldo_fisico_total,
    COALESCE(s.saldo_virtual_total, 0) AS saldo_virtual_total,
    COALESCE(dp.saldo_fisico, 0) AS saldo_fisico,
    COALESCE(dp.saldo_virtual, 0) AS saldo_virtual,
    COALESCE(sp.preco_custo, 0) AS preco_custo,
    COALESCE(sp.preco_compra, 0) AS preco_compra,
    sp.supplier_id
FROM
    products p
LEFT JOIN
    aggregated_stocks s
    ON p.id = s.product_id
LEFT JOIN
    aggregated_deposit_products dp
    ON p.id = dp.product_id
LEFT JOIN
    aggregated_supplier_products sp
    ON p.id = sp.product_id
WHERE
    ($1::text IS NULL OR $1 = '' OR p.nome ILIKE '%' || $1 || '%')
    AND ($2::text IS NULL OR $2 = '' OR p.situacao = $2::text)
    AND s.saldo_fisico_total = 0
    AND s.saldo_virtual_total = 0
    AND dp.saldo_fisico = 0
    AND dp.saldo_virtual = 0
ORDER BY
    p.nome;

-- name: GetProductNoMovements :many
WITH aggregated_stocks AS (
    SELECT product_id,
           SUM(saldo_fisico_total) AS saldo_fisico_total,
           SUM(saldo_virtual_total) AS saldo_virtual_total
    FROM stocks
    GROUP BY product_id
),
aggregated_deposit_products AS (
    SELECT product_id,
           SUM(saldo_fisico) AS saldo_fisico,
           SUM(saldo_virtual) AS saldo_virtual
    FROM deposit_products
    GROUP BY product_id
),
aggregated_supplier_products AS (
    SELECT product_id,
           AVG(preco_custo) AS preco_custo,
           AVG(preco_compra) AS preco_compra,
           supplier_id,
           MAX(descricao) AS descricao,
           MAX(codigo) AS codigo
    FROM supplier_products
    GROUP BY product_id, supplier_id
),
aggregated_sales_orders AS (
    SELECT pso.product_id,
           so.numero,
           so.numeroloja,
           so.data,
           so.datasaida,
           so.dataprevista,
           COALESCE(so.totalprodutos, 0) AS totalprodutos,
           COALESCE(so.totaldescontos, 0) AS totaldescontos
    FROM products_sales_orders pso
    JOIN sales_orders so ON pso.sales_order_id = so.id
    WHERE so.datasaida < NOW() - INTERVAL '1 week'
)
SELECT
    p.ID,
    p.idProdutoPai,
    p.nome,
    p.codigo,
    p.preco,
    p.tipo,
    p.situacao,
    p.formato,
    p.descricao_curta,
    p.imagem_url,
    p.dataValidade,
    p.unidade,
    p.pesoLiquido,
    p.pesoBruto,
    p.volumes,
    p.itensPorCaixa,
    p.gtin,
    p.gtinEmbalagem,
    p.tipoProducao,
    p.condicao,
    p.freteGratis,
    p.marca,
    p.descricaoComplementar,
    p.linkExterno,
    p.observacoes,
    p.descricaoEmbalagemDiscreta,
    p.created_at,
    p.updated_at,
    COALESCE(s.saldo_fisico_total, 0) AS saldo_fisico_total,
    COALESCE(s.saldo_virtual_total, 0) AS saldo_virtual_total,
    COALESCE(dp.saldo_fisico, 0) AS saldo_fisico,
    COALESCE(dp.saldo_virtual, 0) AS saldo_virtual,
    COALESCE(sp.preco_custo, 0) AS preco_custo,
    COALESCE(sp.preco_compra, 0) AS preco_compra,
    sp.supplier_id,
    so.numero,
    so.numeroloja,
    so.data,
    so.datasaida,
    so.dataprevista,
    COALESCE(so.totalprodutos, 0) AS totalprodutos,
    COALESCE(so.totaldescontos, 0) AS totaldescontos
FROM
    products p
LEFT JOIN
    aggregated_stocks s ON p.id = s.product_id
LEFT JOIN
    aggregated_deposit_products dp ON p.id = dp.product_id
LEFT JOIN
    aggregated_supplier_products sp ON p.id = sp.product_id
LEFT JOIN
    aggregated_sales_orders so ON p.id = so.product_id
WHERE
    ($1::text IS NULL OR $1 = '' OR p.nome ILIKE '%' || $1 || '%')
    AND ($2::text IS NULL OR $2 = '' OR p.situacao = $2::text);

-- name: GetProductBySupplierID :one
SELECT
    p.ID,
    p.idProdutoPai,
    p.nome,
    p.codigo,
    p.preco,
    p.tipo,
    p.situacao,
    p.formato,
    p.descricao_curta,
    p.imagem_url,
    p.dataValidade,
    p.unidade,
    p.pesoLiquido,
    p.pesoBruto,
    p.volumes,
    p.itensPorCaixa,
    p.gtin,
    p.gtinEmbalagem,
    p.tipoProducao,
    p.condicao,
    p.freteGratis,
    p.marca,
    p.descricaoComplementar,
    p.linkExterno,
    p.observacoes,
    p.descricaoEmbalagemDiscreta,
    p.created_at,
    p.updated_at,
    COALESCE(s.saldo_fisico_total, 0) AS saldo_fisico_total,
    COALESCE(s.saldo_virtual_total, 0) AS saldo_virtual_total,
    COALESCE(dp.saldo_fisico, 0) AS saldo_fisico,
    COALESCE(dp.saldo_virtual, 0) AS saldo_virtual,
    COALESCE(sp.preco_custo, 0) AS preco_custo,
    COALESCE(sp.preco_compra, 0) AS preco_compra,
    sp.supplier_id
FROM
    products p
LEFT JOIN
    (SELECT product_id, SUM(saldo_fisico_total) as saldo_fisico_total, SUM(saldo_virtual_total) as saldo_virtual_total FROM stocks GROUP BY product_id) s
    ON p.id = s.product_id
LEFT JOIN
    (SELECT product_id, SUM(saldo_fisico) as saldo_fisico, SUM(saldo_virtual) as saldo_virtual FROM deposit_products GROUP BY product_id) dp
    ON p.id = dp.product_id
LEFT JOIN
    (SELECT product_id, AVG(preco_custo) as preco_custo, AVG(preco_compra) as preco_compra, supplier_id
     FROM supplier_products
     WHERE supplier_id = $1
     GROUP BY product_id, supplier_id) sp
    ON p.id = sp.product_id
WHERE sp.supplier_id IS NOT NULL;

-- name: UpdateProduct :exec
UPDATE products SET idProdutoPai = $2,
  nome = $3,
  codigo = $4,
  preco = $5,
  tipo  = $6,
  situacao = $7,
  formato = $8,
  descricao_curta = $9,
  imagem_url = $10,
  dataValidade = $11,
  unidade = $12,
  pesoLiquido = $13,
  pesoBruto = $14,
  volumes = $15,
  itensPorCaixa = $16,
  gtin = $17,
  gtinEmbalagem = $18,
  tipoProducao = $19,
  condicao = $20,
  freteGratis = $21,
  marca = $22,
  descricaoComplementar = $23,
  linkExterno = $24,
  observacoes = $25,
  descricaoEmbalagemDiscreta = $26,
  updated_at = $27
WHERE products.id = $1;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE products.id = $1;
