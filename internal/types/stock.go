package types

import "time"

type Stock struct {
	ID                int64     `json:"id"`
	ProductID         int64     `json:"product_id"`
	Saldofisicototal  int32     `json:"saldo_fisico_total"`
	Saldovirtualtotal int32     `json:"saldo_virtual_total"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type StockStore interface {
	CreateStock(Stock) error
	UpdateStock(Stock) error
}
