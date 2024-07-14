-- name: CreateProduct :exec
INSERT INTO products (ID, idProdutoPai, nome, codigo, preco, tipo, situacao, formato, descricao_curta, imagem_url, dataValidade, unidade, pesoLiquido, pesoBruto, volumes, itensPorCaixa, gtin, gtinEmbalagem, tipoProducao, condicao, freteGratis, marca, descricaoComplementar, linkExterno, observacoes, descricaoEmbalagemDiscreta, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28);

-- name: GetProduct :one
SELECT p.ID,
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
       s.saldo_fisico_total,
       s.saldo_virtual_total,
       dp.saldo_fisico,
       dp.saldo_virtual
FROM 
    products p
LEFT JOIN 
    stocks s ON p.id = s.product_id
LEFT JOIN 
    deposit_products dp ON p.id = dp.product_id
WHERE p.id = $1;

-- name: GetProducts :many
SELECT p.ID,
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
       s.saldo_fisico_total,
       s.saldo_virtual_total,
       dp.saldo_fisico,
       dp.saldo_virtual
FROM 
    products p
LEFT JOIN 
    stocks s ON p.id = s.product_id
LEFT JOIN 
    deposit_products dp ON p.id = dp.product_id;

-- name: GetProductByName :one
SELECT p.ID,
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
       s.saldo_fisico_total,
       s.saldo_virtual_total,
       dp.saldo_fisico,
       dp.saldo_virtual
FROM 
    products p
LEFT JOIN 
    stocks s ON p.id = s.product_id
LEFT JOIN 
    deposit_products dp ON p.id = dp.product_id
WHERE p.nome = $1;

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
