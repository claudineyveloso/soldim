// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: items_sales_order.sql

package db

import (
	"context"
	"time"
)

const createItemsSalesOrder = `-- name: CreateItemsSalesOrder :exec
INSERT INTO items_sales_orders (id, sales_order_id, codigo, unidade, quantidade, desconto, valor, aliquotaIPI, descricao, descricaoDetalhada, product_id, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
`

type CreateItemsSalesOrderParams struct {
	ID                 int64     `json:"id"`
	SalesOrderID       int64     `json:"sales_order_id"`
	Codigo             string    `json:"codigo"`
	Unidade            string    `json:"unidade"`
	Quantidade         int32     `json:"quantidade"`
	Desconto           float64   `json:"desconto"`
	Valor              float64   `json:"valor"`
	Aliquotaipi        float64   `json:"aliquotaipi"`
	Descricao          string    `json:"descricao"`
	Descricaodetalhada string    `json:"descricaodetalhada"`
	ProductID          int64     `json:"product_id"`
	CreatedAt          time.Time `json:"created_at"`
}

func (q *Queries) CreateItemsSalesOrder(ctx context.Context, arg CreateItemsSalesOrderParams) error {
	_, err := q.db.ExecContext(ctx, createItemsSalesOrder,
		arg.ID,
		arg.SalesOrderID,
		arg.Codigo,
		arg.Unidade,
		arg.Quantidade,
		arg.Desconto,
		arg.Valor,
		arg.Aliquotaipi,
		arg.Descricao,
		arg.Descricaodetalhada,
		arg.ProductID,
		arg.CreatedAt,
	)
	return err
}