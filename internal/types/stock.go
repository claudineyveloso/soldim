package types

import "time"

type Stock struct {
	ProductID         int64     `json:"product_id"`
	SaldoFisicoTotal  int32     `json:"saldo_fisico_total"`
	SaldoVirtualTotal int32     `json:"saldo_virtual_total"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type StockResponse struct {
	Data []struct {
		Produto struct {
			ID int64 `json:"id"`
		} `json:"produto"`
		SaldoFisicoTotal  int `json:"saldoFisicoTotal"`
		SaldoVirtualTotal int `json:"saldoVirtualTotal"`
		Depositos         []struct {
			ID           int64 `json:"id"`
			SaldoFisico  int   `json:"saldoFisico"`
			SaldoVirtual int   `json:"saldoVirtual"`
		} `json:"depositos"`
	} `json:"data"`
}

type StockStore interface {
	CreateStock(Stock) error
	UpdateStock(Stock) error
}
