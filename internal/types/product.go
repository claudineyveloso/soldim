package types

import (
	"time"
)

type Product struct {
	ID                int64     `json:"id"`
	Nome              string    `json:"nome"`
	Codigo            string    `json:"codigo"`
	Preco             float64   `json:"preco"`
	ImagemUrl         string    `json:"imagem_url"`
	Tipo              string    `json:"tipo"`
	Situacao          string    `json:"situacao"`
	Formato           string    `json:"formato"`
	DataValidade      time.Time `json:"data_validade"`
	Unidade           string    `json:"unidade"`
	Condicao          int32     `json:"condicao"`
	Gtin              string    `json:"gtin"`
	DescricaoCurta    string    `json:"descricao_curta"`
	SaldoFisicoTotal  int       `json:"saldo_fisico_total"`
	SaldoVirtualTotal int       `json:"saldo_virtual_total"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type ProductWrapper struct {
	Produto Product `json:"produto"`
}

type ProductPayload struct {
	ID                int64     `json:"id"`
	Nome              string    `json:"nome"`
	Codigo            string    `json:"codigo"`
	Preco             float64   `json:"preco"`
	ImagemUrl         string    `json:"imagem_url"`
	Tipo              string    `json:"tipo"`
	Situacao          string    `json:"situacao"`
	Formato           string    `json:"formato"`
	DataValidade      time.Time `json:"data_validade"`
	Unidade           string    `json:"unidade"`
	Condicao          int32     `json:"condicao"`
	Gtin              string    `json:"gtin"`
	DescricaoCurta    string    `json:"descricao_curta"`
	SaldoFisicoTotal  int       `json:"saldo_fisico_total"`
	SaldoVirtualTotal int       `json:"saldo_virtual_total"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type StockResponse struct {
	Data []struct {
		Produto struct {
			ID int `json:"id"`
		} `json:"produto"`
		SaldoFisicoTotal  int `json:"saldoFisicoTotal"`
		SaldoVirtualTotal int `json:"saldoVirtualTotal"`
	} `json:"data"`
}

//	type ProductResponse struct {
//		Retorno struct {
//			Products []ProductWrapper `json:"produtos"`
//			Total    int              `json:"total"`
//			Limit    int              `json:"limit"`
//		} `json:"retorno"`
//	}
type ProductResponse struct {
	Data  []Product `json:"data"`
	Total int       `json:"total"`
	Limit int       `json:"limit"`
}

type ProductStore interface {
	CreateProduct(ProductPayload) error
	GetProducts() ([]*Product, error)
	GetProductByID(id int64) (*Product, error)
	UpdateProduct(ProductPayload) error
	DeleteProduct(id int64) error
}
