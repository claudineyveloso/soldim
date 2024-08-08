package types

import (
	"fmt"
	"time"
)

type SalesOrder struct {
	ID             int64            `json:"id"`
	Numero         int32            `json:"numero"`
	Numeroloja     string           `json:"numeroloja"`
	Data           CustomDate       `json:"data"`
	Datasaida      CustomDate       `json:"datasaida"`
	Dataprevista   CustomDate       `json:"dataprevista"`
	Totalprodutos  float64          `json:"totalprodutos"`
	Totaldescontos float64          `json:"totaldescontos"`
	SituationID    int64            `json:"situation_id"`
	StoreID        int64            `json:"store_id"`
	Contato        Contato          `json:"contato"`
	Situacao       Situacao         `json:"situacao"`
	Loja           Loja             `json:"loja"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
	Itens          []SalesOrderItem `json:"itens"`
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
	ID int64 `json:"id"`
}

type Comissao struct {
	Base     float64 `json:"base"`
	Aliquota float64 `json:"aliquota"`
	Valor    float64 `json:"valor"`
}

// Situacao representa a estrutura da situação no pedido de venda
type Situacao struct {
	ID int64 `json:"id"`
}

// Loja representa a estrutura da loja no pedido de venda
type Loja struct {
	ID int64 `json:"id"`
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

type SalesOrderStore interface {
	CreateSalesOrder(SalesOrder) error
	GetSalesOrders(limit, offset int32) ([]*SalesOrder, error)
	GetSalesOrderByID(id int64) (*SalesOrder, error)
	GetSalesOrderByNumber(numero int32) (*SalesOrder, error)
}
