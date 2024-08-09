package types

import "time"

type ItemsSalesOrder struct {
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
	UpdatedAt          time.Time `json:"updated_at"`
}

type ItemsSalesOrderByIDRow struct {
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
	CreatedAt          int64     `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type ItemsSalesOrderByProductIDRow struct {
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
	CreatedAt          int64     `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type ItemsSalesOrderBySalesOrderIDRow struct {
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
	CreatedAt          int64     `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type ItemsSalesOrderStore interface {
	CreateItemsSalesOrder(ItemsSalesOrder) error
	GetItemsSalesOrderByID(id int64) (*ItemsSalesOrder, error)
	GetItemsSalesOrderByProductID(id int64) (*ItemsSalesOrder, error)
	GetItemsSalesOrderBySalesOrderID(id int64) (*ItemsSalesOrder, error)
}
