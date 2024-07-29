package types

import (
	"database/sql"
	"time"
)

type ProductSalesOrder struct {
	ID             int64     `json:"id"`
	SalesOrderID   int64     `json:"sales_order_id"`
	ProductID      int64     `json:"product_id"`
	Quantidade     int32     `json:"quantidade"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Numero         int32     `json:"numero"`
	Numeroloja     string    `json:"numeroloja"`
	Data           time.Time `json:"data"`
	Datasaida      time.Time `json:"datasaida"`
	Dataprevista   time.Time `json:"dataprevista"`
	Totalprodutos  float64   `json:"totalprodutos"`
	Totaldescontos float64   `json:"totaldescontos"`
	SituationID    int64     `json:"situation_id"`
	StoreID        int64     `json:"store_id"`
	SupplierID     int64     `json:"supplier_id"`
	Nome           string    `json:"nome"`
	Codigo         string    `json:"codigo"`
	Preco          float64   `json:"preco"`
}

type ProductSalesOrderPayload struct {
	SalesOrderID int64     `json:"sales_order_id"`
	ProductID    int64     `json:"product_id"`
	Quantidade   int32     `json:"quantidade"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func ToInt64(n sql.NullInt64) int64 {
	if n.Valid {
		return n.Int64
	}
	return 0
}

func ToInt32(n sql.NullInt32) int32 {
	if n.Valid {
		return n.Int32
	}
	return 0
}

func ToString(s sql.NullString) string {
	if s.Valid {
		return s.String
	}
	return ""
}

func ToFloat64(f sql.NullFloat64) float64 {
	if f.Valid {
		return f.Float64
	}
	return 0.0
}

func ToTime(t sql.NullTime) time.Time {
	if t.Valid {
		return t.Time
	}
	return time.Time{}
}

type ProductSalesOrderStore interface {
	CreateProductSalesOrder(ProductSalesOrderPayload) error
	GetProductSalesOrders() ([]*ProductSalesOrder, error)
	GetProductSalesOrdersBySupplierID(supplierID int64) ([]*ProductSalesOrder, error)
}
