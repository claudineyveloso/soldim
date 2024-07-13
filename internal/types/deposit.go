package types

import "time"

type Deposit struct {
	ID                 int64     `json:"id"`
	Descricao          string    `json:"descricao"`
	Situacao           int32     `json:"situacao"`
	Padrao             bool      `json:"padrao"`
	Desconsiderarsaldo bool      `json:"desconsiderar_saldo"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type DepositStore interface {
	CreateDeposit(Deposit) error
	GetDeposits() ([]*Deposit, error)
}
