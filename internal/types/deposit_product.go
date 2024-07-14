package types

import "time"

type DepositProduct struct {
	ID           int64     `json:"id"`
	DepositID    int64     `json:"deposit_id"`
	ProductID    int64     `json:"product_id"`
	SaldoFisico  int32     `json:"saldo_fisico"`
	SaldoVirtual int32     `json:"saldo_virtual"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type DepositProductStore interface {
	CreateDepositProduct(DepositProduct) error
	UpdateDepositProduct(DepositProduct) error
}
