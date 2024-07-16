package types

import "time"

type SupplierProduct struct {
	ID          int64     `json:"id"`
	Descricao   string    `json:"descricao"`
	PrecoCusto  float64   `json:"precoCusto"`
	PrecoCompra float64   `json:"precoCompra"`
	Padrao      bool      `json:"padrao"`
	SupplierID  int64     `json:"supplierID"`
	ProductID   int64     `json:"productID"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SupplierResponse struct {
	Data []struct {
		ID          int64   `json:"id"`
		Descricao   string  `json:"descricao"`
		PrecoCusto  float64 `json:"precoCusto"`
		PrecoCompra float64 `json:"precoCompra"`
		Padrao      bool    `json:"padrao"`
		Produto     struct {
			ID int64 `json:"id"`
		} `json:"produto"`
		Fornecedor struct {
			ID int64 `json:"id"`
		} `json:"fornecedor"`
	} `json:"data"`
}

type SupplierProductStore interface {
	CreateSupplierProduct(SupplierProduct) error
	UpdateSupplierProduct(SupplierProduct) error
}
