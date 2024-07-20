package types

import "time"

type ProductSalesOrder struct {
	SalesOrderID int64     `json:"sales_order_id"`
	ProductID    int64     `json:"product_id"`
	Quantidade   int32     `json:"quantidade"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ProductSalesOrderStore interface {
	CreateProductSalesOrde(ProductSalesOrder) error
	GetProductSalesOrders() ([]*ProductSalesOrder, error)
	GetProductSalesOrdersBySupplierID(id int64) (*ProductSalesOrder, error)
}
