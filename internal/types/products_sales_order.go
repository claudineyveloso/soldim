package types

import (
	"database/sql"
	"time"
)

type ProductSalesOrder struct {
	ID             int64
	SalesOrderID   int64
	ProductID      int64
	Quantidade     int32
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Numero         int32
	Numeroloja     string
	Data           time.Time
	Datasaida      time.Time
	Dataprevista   time.Time
	Totalprodutos  float64
	Totaldescontos float64
	SituationID    int64
	StoreID        int64
	SupplierID     int64
	Nome           string
	Codigo         string
	Preco          float64
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
	CreateProductSalesOrder(ProductSalesOrder) error
	GetProductSalesOrders() ([]*ProductSalesOrder, error)
	GetProductSalesOrdersBySupplierID(supplierID int64) ([]*ProductSalesOrder, error)
}
