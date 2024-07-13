// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: product.sql

package db

import (
	"context"
	"database/sql"
	"time"
)

const createProduct = `-- name: CreateProduct :exec
INSERT INTO products (ID, idProdutoPai, nome, codigo, preco, tipo, situacao, formato, descricao_curta, imagem_url, dataValidade, unidade, pesoLiquido, pesoBruto, volumes, itensPorCaixa, gtin, gtinEmbalagem, tipoProducao, condicao, freteGratis, marca, descricaoComplementar, linkExterno, observacoes, descricaoEmbalagemDiscreta, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28)
`

type CreateProductParams struct {
	ID                         int64     `json:"id"`
	Idprodutopai               int64     `json:"idprodutopai"`
	Nome                       string    `json:"nome"`
	Codigo                     string    `json:"codigo"`
	Preco                      float64   `json:"preco"`
	Tipo                       string    `json:"tipo"`
	Situacao                   string    `json:"situacao"`
	Formato                    string    `json:"formato"`
	DescricaoCurta             string    `json:"descricao_curta"`
	ImagemUrl                  string    `json:"imagem_url"`
	Datavalidade               time.Time `json:"datavalidade"`
	Unidade                    string    `json:"unidade"`
	Pesoliquido                float64   `json:"pesoliquido"`
	Pesobruto                  float64   `json:"pesobruto"`
	Volumes                    int32     `json:"volumes"`
	Itensporcaixa              int32     `json:"itensporcaixa"`
	Gtin                       string    `json:"gtin"`
	Gtinembalagem              string    `json:"gtinembalagem"`
	Tipoproducao               string    `json:"tipoproducao"`
	Condicao                   int32     `json:"condicao"`
	Fretegratis                bool      `json:"fretegratis"`
	Marca                      string    `json:"marca"`
	Descricaocomplementar      string    `json:"descricaocomplementar"`
	Linkexterno                string    `json:"linkexterno"`
	Observacoes                string    `json:"observacoes"`
	Descricaoembalagemdiscreta string    `json:"descricaoembalagemdiscreta"`
	CreatedAt                  time.Time `json:"created_at"`
	UpdatedAt                  time.Time `json:"updated_at"`
}

func (q *Queries) CreateProduct(ctx context.Context, arg CreateProductParams) error {
	_, err := q.db.ExecContext(ctx, createProduct,
		arg.ID,
		arg.Idprodutopai,
		arg.Nome,
		arg.Codigo,
		arg.Preco,
		arg.Tipo,
		arg.Situacao,
		arg.Formato,
		arg.DescricaoCurta,
		arg.ImagemUrl,
		arg.Datavalidade,
		arg.Unidade,
		arg.Pesoliquido,
		arg.Pesobruto,
		arg.Volumes,
		arg.Itensporcaixa,
		arg.Gtin,
		arg.Gtinembalagem,
		arg.Tipoproducao,
		arg.Condicao,
		arg.Fretegratis,
		arg.Marca,
		arg.Descricaocomplementar,
		arg.Linkexterno,
		arg.Observacoes,
		arg.Descricaoembalagemdiscreta,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	return err
}

const deleteProduct = `-- name: DeleteProduct :exec
DELETE FROM products
WHERE products.id = $1
`

func (q *Queries) DeleteProduct(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteProduct, id)
	return err
}

const getProduct = `-- name: GetProduct :one
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
       s.saldofisicototal,
       s.saldovirtualtotal,
       dp.saldofisico,
       dp.saldovirtual
FROM 
    products p
LEFT JOIN 
    stocks s ON p.id = s.product_id
LEFT JOIN 
    deposit_products dp ON p.id = dp.product_id
WHERE p.id = $1
`

type GetProductRow struct {
	ID                         int64         `json:"id"`
	Idprodutopai               int64         `json:"idprodutopai"`
	Nome                       string        `json:"nome"`
	Codigo                     string        `json:"codigo"`
	Preco                      float64       `json:"preco"`
	Tipo                       string        `json:"tipo"`
	Situacao                   string        `json:"situacao"`
	Formato                    string        `json:"formato"`
	DescricaoCurta             string        `json:"descricao_curta"`
	ImagemUrl                  string        `json:"imagem_url"`
	Datavalidade               time.Time     `json:"datavalidade"`
	Unidade                    string        `json:"unidade"`
	Pesoliquido                float64       `json:"pesoliquido"`
	Pesobruto                  float64       `json:"pesobruto"`
	Volumes                    int32         `json:"volumes"`
	Itensporcaixa              int32         `json:"itensporcaixa"`
	Gtin                       string        `json:"gtin"`
	Gtinembalagem              string        `json:"gtinembalagem"`
	Tipoproducao               string        `json:"tipoproducao"`
	Condicao                   int32         `json:"condicao"`
	Fretegratis                bool          `json:"fretegratis"`
	Marca                      string        `json:"marca"`
	Descricaocomplementar      string        `json:"descricaocomplementar"`
	Linkexterno                string        `json:"linkexterno"`
	Observacoes                string        `json:"observacoes"`
	Descricaoembalagemdiscreta string        `json:"descricaoembalagemdiscreta"`
	CreatedAt                  time.Time     `json:"created_at"`
	UpdatedAt                  time.Time     `json:"updated_at"`
	Saldofisicototal           sql.NullInt32 `json:"saldofisicototal"`
	Saldovirtualtotal          sql.NullInt32 `json:"saldovirtualtotal"`
	Saldofisico                sql.NullInt32 `json:"saldofisico"`
	Saldovirtual               sql.NullInt32 `json:"saldovirtual"`
}

func (q *Queries) GetProduct(ctx context.Context, id int64) (GetProductRow, error) {
	row := q.db.QueryRowContext(ctx, getProduct, id)
	var i GetProductRow
	err := row.Scan(
		&i.ID,
		&i.Idprodutopai,
		&i.Nome,
		&i.Codigo,
		&i.Preco,
		&i.Tipo,
		&i.Situacao,
		&i.Formato,
		&i.DescricaoCurta,
		&i.ImagemUrl,
		&i.Datavalidade,
		&i.Unidade,
		&i.Pesoliquido,
		&i.Pesobruto,
		&i.Volumes,
		&i.Itensporcaixa,
		&i.Gtin,
		&i.Gtinembalagem,
		&i.Tipoproducao,
		&i.Condicao,
		&i.Fretegratis,
		&i.Marca,
		&i.Descricaocomplementar,
		&i.Linkexterno,
		&i.Observacoes,
		&i.Descricaoembalagemdiscreta,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Saldofisicototal,
		&i.Saldovirtualtotal,
		&i.Saldofisico,
		&i.Saldovirtual,
	)
	return i, err
}

const getProductByName = `-- name: GetProductByName :one
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
       s.saldofisicototal,
       s.saldovirtualtotal,
       dp.saldofisico,
       dp.saldovirtual
FROM 
    products p
LEFT JOIN 
    stocks s ON p.id = s.product_id
LEFT JOIN 
    deposit_products dp ON p.id = dp.product_id
WHERE p.nome = $1
`

type GetProductByNameRow struct {
	ID                         int64         `json:"id"`
	Idprodutopai               int64         `json:"idprodutopai"`
	Nome                       string        `json:"nome"`
	Codigo                     string        `json:"codigo"`
	Preco                      float64       `json:"preco"`
	Tipo                       string        `json:"tipo"`
	Situacao                   string        `json:"situacao"`
	Formato                    string        `json:"formato"`
	DescricaoCurta             string        `json:"descricao_curta"`
	ImagemUrl                  string        `json:"imagem_url"`
	Datavalidade               time.Time     `json:"datavalidade"`
	Unidade                    string        `json:"unidade"`
	Pesoliquido                float64       `json:"pesoliquido"`
	Pesobruto                  float64       `json:"pesobruto"`
	Volumes                    int32         `json:"volumes"`
	Itensporcaixa              int32         `json:"itensporcaixa"`
	Gtin                       string        `json:"gtin"`
	Gtinembalagem              string        `json:"gtinembalagem"`
	Tipoproducao               string        `json:"tipoproducao"`
	Condicao                   int32         `json:"condicao"`
	Fretegratis                bool          `json:"fretegratis"`
	Marca                      string        `json:"marca"`
	Descricaocomplementar      string        `json:"descricaocomplementar"`
	Linkexterno                string        `json:"linkexterno"`
	Observacoes                string        `json:"observacoes"`
	Descricaoembalagemdiscreta string        `json:"descricaoembalagemdiscreta"`
	CreatedAt                  time.Time     `json:"created_at"`
	UpdatedAt                  time.Time     `json:"updated_at"`
	Saldofisicototal           sql.NullInt32 `json:"saldofisicototal"`
	Saldovirtualtotal          sql.NullInt32 `json:"saldovirtualtotal"`
	Saldofisico                sql.NullInt32 `json:"saldofisico"`
	Saldovirtual               sql.NullInt32 `json:"saldovirtual"`
}

func (q *Queries) GetProductByName(ctx context.Context, nome string) (GetProductByNameRow, error) {
	row := q.db.QueryRowContext(ctx, getProductByName, nome)
	var i GetProductByNameRow
	err := row.Scan(
		&i.ID,
		&i.Idprodutopai,
		&i.Nome,
		&i.Codigo,
		&i.Preco,
		&i.Tipo,
		&i.Situacao,
		&i.Formato,
		&i.DescricaoCurta,
		&i.ImagemUrl,
		&i.Datavalidade,
		&i.Unidade,
		&i.Pesoliquido,
		&i.Pesobruto,
		&i.Volumes,
		&i.Itensporcaixa,
		&i.Gtin,
		&i.Gtinembalagem,
		&i.Tipoproducao,
		&i.Condicao,
		&i.Fretegratis,
		&i.Marca,
		&i.Descricaocomplementar,
		&i.Linkexterno,
		&i.Observacoes,
		&i.Descricaoembalagemdiscreta,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Saldofisicototal,
		&i.Saldovirtualtotal,
		&i.Saldofisico,
		&i.Saldovirtual,
	)
	return i, err
}

const getProducts = `-- name: GetProducts :many
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
       s.saldofisicototal,
       s.saldovirtualtotal,
       dp.saldofisico,
       dp.saldovirtual
FROM 
    products p
LEFT JOIN 
    stocks s ON p.id = s.product_id
LEFT JOIN 
    deposit_products dp ON p.id = dp.product_id
`

type GetProductsRow struct {
	ID                         int64         `json:"id"`
	Idprodutopai               int64         `json:"idprodutopai"`
	Nome                       string        `json:"nome"`
	Codigo                     string        `json:"codigo"`
	Preco                      float64       `json:"preco"`
	Tipo                       string        `json:"tipo"`
	Situacao                   string        `json:"situacao"`
	Formato                    string        `json:"formato"`
	DescricaoCurta             string        `json:"descricao_curta"`
	ImagemUrl                  string        `json:"imagem_url"`
	Datavalidade               time.Time     `json:"datavalidade"`
	Unidade                    string        `json:"unidade"`
	Pesoliquido                float64       `json:"pesoliquido"`
	Pesobruto                  float64       `json:"pesobruto"`
	Volumes                    int32         `json:"volumes"`
	Itensporcaixa              int32         `json:"itensporcaixa"`
	Gtin                       string        `json:"gtin"`
	Gtinembalagem              string        `json:"gtinembalagem"`
	Tipoproducao               string        `json:"tipoproducao"`
	Condicao                   int32         `json:"condicao"`
	Fretegratis                bool          `json:"fretegratis"`
	Marca                      string        `json:"marca"`
	Descricaocomplementar      string        `json:"descricaocomplementar"`
	Linkexterno                string        `json:"linkexterno"`
	Observacoes                string        `json:"observacoes"`
	Descricaoembalagemdiscreta string        `json:"descricaoembalagemdiscreta"`
	CreatedAt                  time.Time     `json:"created_at"`
	UpdatedAt                  time.Time     `json:"updated_at"`
	Saldofisicototal           sql.NullInt32 `json:"saldofisicototal"`
	Saldovirtualtotal          sql.NullInt32 `json:"saldovirtualtotal"`
	Saldofisico                sql.NullInt32 `json:"saldofisico"`
	Saldovirtual               sql.NullInt32 `json:"saldovirtual"`
}

func (q *Queries) GetProducts(ctx context.Context) ([]GetProductsRow, error) {
	rows, err := q.db.QueryContext(ctx, getProducts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetProductsRow
	for rows.Next() {
		var i GetProductsRow
		if err := rows.Scan(
			&i.ID,
			&i.Idprodutopai,
			&i.Nome,
			&i.Codigo,
			&i.Preco,
			&i.Tipo,
			&i.Situacao,
			&i.Formato,
			&i.DescricaoCurta,
			&i.ImagemUrl,
			&i.Datavalidade,
			&i.Unidade,
			&i.Pesoliquido,
			&i.Pesobruto,
			&i.Volumes,
			&i.Itensporcaixa,
			&i.Gtin,
			&i.Gtinembalagem,
			&i.Tipoproducao,
			&i.Condicao,
			&i.Fretegratis,
			&i.Marca,
			&i.Descricaocomplementar,
			&i.Linkexterno,
			&i.Observacoes,
			&i.Descricaoembalagemdiscreta,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Saldofisicototal,
			&i.Saldovirtualtotal,
			&i.Saldofisico,
			&i.Saldovirtual,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateProduct = `-- name: UpdateProduct :exec
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
WHERE products.id = $1
`

type UpdateProductParams struct {
	ID                         int64     `json:"id"`
	Idprodutopai               int64     `json:"idprodutopai"`
	Nome                       string    `json:"nome"`
	Codigo                     string    `json:"codigo"`
	Preco                      float64   `json:"preco"`
	Tipo                       string    `json:"tipo"`
	Situacao                   string    `json:"situacao"`
	Formato                    string    `json:"formato"`
	DescricaoCurta             string    `json:"descricao_curta"`
	ImagemUrl                  string    `json:"imagem_url"`
	Datavalidade               time.Time `json:"datavalidade"`
	Unidade                    string    `json:"unidade"`
	Pesoliquido                float64   `json:"pesoliquido"`
	Pesobruto                  float64   `json:"pesobruto"`
	Volumes                    int32     `json:"volumes"`
	Itensporcaixa              int32     `json:"itensporcaixa"`
	Gtin                       string    `json:"gtin"`
	Gtinembalagem              string    `json:"gtinembalagem"`
	Tipoproducao               string    `json:"tipoproducao"`
	Condicao                   int32     `json:"condicao"`
	Fretegratis                bool      `json:"fretegratis"`
	Marca                      string    `json:"marca"`
	Descricaocomplementar      string    `json:"descricaocomplementar"`
	Linkexterno                string    `json:"linkexterno"`
	Observacoes                string    `json:"observacoes"`
	Descricaoembalagemdiscreta string    `json:"descricaoembalagemdiscreta"`
	UpdatedAt                  time.Time `json:"updated_at"`
}

func (q *Queries) UpdateProduct(ctx context.Context, arg UpdateProductParams) error {
	_, err := q.db.ExecContext(ctx, updateProduct,
		arg.ID,
		arg.Idprodutopai,
		arg.Nome,
		arg.Codigo,
		arg.Preco,
		arg.Tipo,
		arg.Situacao,
		arg.Formato,
		arg.DescricaoCurta,
		arg.ImagemUrl,
		arg.Datavalidade,
		arg.Unidade,
		arg.Pesoliquido,
		arg.Pesobruto,
		arg.Volumes,
		arg.Itensporcaixa,
		arg.Gtin,
		arg.Gtinembalagem,
		arg.Tipoproducao,
		arg.Condicao,
		arg.Fretegratis,
		arg.Marca,
		arg.Descricaocomplementar,
		arg.Linkexterno,
		arg.Observacoes,
		arg.Descricaoembalagemdiscreta,
		arg.UpdatedAt,
	)
	return err
}
