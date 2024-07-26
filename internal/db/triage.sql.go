// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: triage.sql

package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createTriage = `-- name: CreateTriage :exec
INSERT INTO triages (id, type, grid, sku_sap, sku_wms, description, cust_id, seller, quantity_supplied, final_quantity, unitary_value, total_value_offered, final_total_value, category, sub_category, sent_to_batch, sent_to_bling, defect, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
`

type CreateTriageParams struct {
	ID                uuid.UUID `json:"id"`
	Type              string    `json:"type"`
	Grid              string    `json:"grid"`
	SkuSap            int32     `json:"sku_sap"`
	SkuWms            string    `json:"sku_wms"`
	Description       string    `json:"description"`
	CustID            int64     `json:"cust_id"`
	Seller            string    `json:"seller"`
	QuantitySupplied  int32     `json:"quantity_supplied"`
	FinalQuantity     int32     `json:"final_quantity"`
	UnitaryValue      float64   `json:"unitary_value"`
	TotalValueOffered float64   `json:"total_value_offered"`
	FinalTotalValue   float64   `json:"final_total_value"`
	Category          string    `json:"category"`
	SubCategory       string    `json:"sub_category"`
	SentToBatch       bool      `json:"sent_to_batch"`
	SentToBling       bool      `json:"sent_to_bling"`
	Defect            bool      `json:"defect"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (q *Queries) CreateTriage(ctx context.Context, arg CreateTriageParams) error {
	_, err := q.db.ExecContext(ctx, createTriage,
		arg.ID,
		arg.Type,
		arg.Grid,
		arg.SkuSap,
		arg.SkuWms,
		arg.Description,
		arg.CustID,
		arg.Seller,
		arg.QuantitySupplied,
		arg.FinalQuantity,
		arg.UnitaryValue,
		arg.TotalValueOffered,
		arg.FinalTotalValue,
		arg.Category,
		arg.SubCategory,
		arg.SentToBatch,
		arg.SentToBling,
		arg.Defect,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	return err
}

const getTotalTriages = `-- name: GetTotalTriages :one
SELECT COUNT(*)
FROM triages
WHERE ($1::text IS NULL OR $1 = '' OR description ILIKE '%' || $1 || '%')
  AND ($2::text IS NULL OR $2 = '' OR sku_wms = $2)
  AND ($3::int IS NULL OR $3 = 0 OR sku_sap = $3)
`

type GetTotalTriagesParams struct {
	Column1 string `json:"column_1"`
	Column2 string `json:"column_2"`
	Column3 int32  `json:"column_3"`
}

func (q *Queries) GetTotalTriages(ctx context.Context, arg GetTotalTriagesParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, getTotalTriages, arg.Column1, arg.Column2, arg.Column3)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getTriages = `-- name: GetTriages :many
SELECT id,
        type,
        grid,
        sku_sap,
        sku_wms,
        description,
        cust_id,
        seller,
        quantity_supplied,
        final_quantity,
        unitary_value,
        total_value_offered,
        final_total_value,
        category,
        sub_category,
        sent_to_batch,
        sent_to_bling,
        defect,
        created_at,
        updated_at
FROM triages
WHERE ($1::text IS NULL OR $1 = '' OR description ILIKE '%' || $1 || '%')
  AND ($2::text IS NULL OR $2 = '' OR sku_wms = $2)
  AND ($3::int IS NULL OR $3 = 0 OR sku_sap = $3)
LIMIT $4 OFFSET $5
`

type GetTriagesParams struct {
	Column1 string `json:"column_1"`
	Column2 string `json:"column_2"`
	Column3 int32  `json:"column_3"`
	Limit   int32  `json:"limit"`
	Offset  int32  `json:"offset"`
}

func (q *Queries) GetTriages(ctx context.Context, arg GetTriagesParams) ([]Triage, error) {
	rows, err := q.db.QueryContext(ctx, getTriages,
		arg.Column1,
		arg.Column2,
		arg.Column3,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Triage
	for rows.Next() {
		var i Triage
		if err := rows.Scan(
			&i.ID,
			&i.Type,
			&i.Grid,
			&i.SkuSap,
			&i.SkuWms,
			&i.Description,
			&i.CustID,
			&i.Seller,
			&i.QuantitySupplied,
			&i.FinalQuantity,
			&i.UnitaryValue,
			&i.TotalValueOffered,
			&i.FinalTotalValue,
			&i.Category,
			&i.SubCategory,
			&i.SentToBatch,
			&i.SentToBling,
			&i.Defect,
			&i.CreatedAt,
			&i.UpdatedAt,
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
