package types

import (
	"database/sql"
	"fmt"
	"time"
)

type SalesOrder struct {
	ID                int64              `json:"id"`
	Numero            int32              `json:"numero"`
	Numeroloja        string             `json:"numeroloja"`
	Data              CustomDate         `json:"data"`
	Datasaida         CustomDate         `json:"datasaida"`
	Dataprevista      CustomDate         `json:"dataprevista"`
	Totalprodutos     float64            `json:"totalprodutos"`
	Totaldescontos    float64            `json:"totaldescontos"`
	SituationID       int64              `json:"situation_id"`
	StoreID           int64              `json:"store_id"`
	ContactID         int64              `json:"contact_id"`
	ItemsSalesOrderID int64              `json:"items_sales_order_id"`
	Contato           Contato            `json:"contato"`
	Situacao          Situacao           `json:"situacao"`
	Loja              Loja               `json:"loja"`
	CreatedAt         time.Time          `json:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at"`
	Items             []ItemsSalesOrders `json:"itens"`
}

type SalesOrderRow struct {
	ID                   int64     `json:"id"`
	Numero               int32     `json:"numero"`
	Numeroloja           string    `json:"numeroloja"`
	Data                 time.Time `json:"data"`
	Datasaida            time.Time `json:"datasaida"`
	Dataprevista         time.Time `json:"dataprevista"`
	Totalprodutos        float64   `json:"totalprodutos"`
	Totaldescontos       float64   `json:"totaldescontos"`
	SituationID          int64     `json:"situation_id"`
	SituationDescription string    `json:"situation_description"`
	StoreID              int64     `json:"store_id"`
	StoreDescription     string    `json:"store_description"`
	ContactID            int64     `json:"contact_id"`
	ContactName          string    `json:"contact_name"`
	ContactDocument      string    `json:"contact_document"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type SalesOrderItem struct {
	ID                 int64        `json:"id"`
	Codigo             string       `json:"codigo"`
	Unidade            string       `json:"unidade"`
	Quantidade         int32        `json:"quantidade"`
	Desconto           float64      `json:"desconto"`
	Valor              float64      `json:"valor"`
	AliquotaIPI        float64      `json:"aliquotaIPI"`
	Descricao          string       `json:"descricao"`
	DescricaoDetalhada string       `json:"descricaoDetalhada"`
	Produto            SalesProduct `json:"produto"`
	Comissao           Comissao     `json:"comissao"`
}

type SalesProduct struct {
	ID int64 `json:"id"`
}

type Contato struct {
	ID              int64  `json:"id"`
	Nome            string `json:"nome"`
	TipoPessoa      string `json:"tipoPessoa"`
	NumeroDocumento string `json:"numeroDocumento"`
}

type Comissao struct {
	Base     float64 `json:"base"`
	Aliquota float64 `json:"aliquota"`
	Valor    float64 `json:"valor"`
}

// Situacao representa a estrutura da situação no pedido de venda
type Situacao struct {
	ID              int64  `json:"id"`
	Nome            string `json:"nome"`
	TipoPessoa      string `json:"tipoPessoa"`
	NumeroDocumento string `json:"numeroDocumento"`
}

// Loja representa a estrutura da loja no pedido de venda
type Loja struct {
	ID int64 `json:"id"`
}

type ItemsSalesOrders struct {
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

type SalesOrderResponse struct {
	Data []SalesOrder `json:"data"`
}

// CustomDate é um tipo customizado para tratar o formato da data
type CustomDate struct {
	time.Time
}

// UnmarshalJSON sobrescreve a função padrão de unmarshal para o tipo CustomDate
func (cd *CustomDate) UnmarshalJSON(b []byte) error {
	s := string(b[1 : len(b)-1])
	if s == "0000-00-00" {
		cd.Time = time.Time{} // Zero value para time.Time
		return nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return fmt.Errorf("erro ao fazer parse da data: %v", err)
	}
	cd.Time = t
	return nil
}

// MarshalJSON sobrescreve a função padrão de marshal para o tipo CustomDate
func (cd CustomDate) MarshalJSON() ([]byte, error) {
	if cd.Time.IsZero() {
		return []byte(`"0000-00-00"`), nil
	}
	return []byte(fmt.Sprintf(`"%s"`, cd.Time.Format("2006-01-02"))), nil
}

func ConvertNullString(nullStr sql.NullString) string {
	if nullStr.Valid {
		return nullStr.String
	}
	return "" // Retorna uma string vazia se o valor for nulo
}

func ConvertNullInt32(nullInt sql.NullInt32) int32 {
	if nullInt.Valid {
		return nullInt.Int32
	}
	return 0 // Retorna 0 se o valor for nulo
}

func ConvertNullInt64(nullInt sql.NullInt64) int64 {
	if nullInt.Valid {
		return nullInt.Int64
	}
	return 0 // Retorna 0 se o valor for nulo
}

func ConvertNullFloat64(nullFloat sql.NullFloat64) float64 {
	if nullFloat.Valid {
		return nullFloat.Float64
	}
	return 0.0 // Retorna 0.0 se o valor for nulo
}

type SalesOrderStore interface {
	CreateSalesOrder(SalesOrder) error
	GetSalesOrders() ([]*SalesOrder, error)
	GetSalesOrderByID(id int64) (*SalesOrder, error)
	// GetSalesOrderByNumber(numero int32) (*SalesOrderRow, error)
}
