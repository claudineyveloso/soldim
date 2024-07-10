package types

import "time"

type SalesChannel struct {
	ID        int32     `json:"id"`
	Descricao string    `json:"descricao"`
	Tipo      string    `json:"tipo"`
	Situacao  int32     `json:"situacao"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SalesChannelStore interface {
	CreateSalesChannel(SalesChannel) error
	GetSalesChannel() ([]*SalesChannel, error)
}
