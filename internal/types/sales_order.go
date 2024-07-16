package types

import "time"

type SalesOrder struct {
	ID             int64     `json:"id"`
	Numero         int32     `json:"numero"`
	Numeroloja     string    `json:"numeroloja"`
	Data           time.Time `json:"data"`
	Datasaida      time.Time `json:"datasaida"`
	Dataprevista   time.Time `json:"dataprevista"`
	Totalprodutos  float64   `json:"totalprodutos"`
	Totaldescontos float64   `json:"totaldescontos"`
	SituationID    int64     `json:"situation_id"`
	StoreID        int64     `json:"store_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type SalesOrderStore interface {
	CreateSalesOrder(SalesOrder) error
	UpdateSalesOrder(SalesOrder) error
}
