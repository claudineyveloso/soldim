package types

import (
	"database/sql"
	"encoding/json"
	"time"
)

type NullableInt struct {
	Value int32
	Valid bool
}

type Product struct {
	ID                         int64       `json:"id"`
	Idprodutopai               int64       `json:"id_produto_pai"`
	Nome                       string      `json:"nome"`
	Codigo                     string      `json:"codigo"`
	Preco                      float64     `json:"preco"`
	ImagemUrl                  string      `json:"imagem_url"`
	Tipo                       string      `json:"tipo"`
	Situacao                   string      `json:"situacao"`
	Formato                    string      `json:"formato"`
	DescricaoCurta             string      `json:"descricao_curta"`
	Datavalidade               time.Time   `json:"data_validade"`
	Unidade                    string      `json:"unidade"`
	Pesoliquido                float64     `json:"peso_liquido"`
	Pesobruto                  float64     `json:"peso_bruto"`
	Volumes                    int32       `json:"volumes"`
	Itensporcaixa              int32       `json:"itens_por_caixa"`
	Gtin                       string      `json:"gtin"`
	Gtinembalagem              string      `json:"gtin_embalagem"`
	Tipoproducao               string      `json:"tipo_producao"`
	Condicao                   int32       `json:"condicao"`
	Fretegratis                bool        `json:"frete_gratis"`
	Marca                      string      `json:"marca"`
	Descricaocomplementar      string      `json:"descricao_complementar"`
	Linkexterno                string      `json:"link_externo"`
	Observacoes                string      `json:"observacoes"`
	Descricaoembalagemdiscreta string      `json:"descricao_embalagem_discreta"`
	CreatedAt                  time.Time   `json:"created_at"`
	UpdatedAt                  time.Time   `json:"updated_at"`
	Saldofisicototal           NullableInt `json:"saldo_fisico_total"`
	Saldovirtualtotal          NullableInt `json:"saldo_virtual_total"`
	Saldofisico                NullableInt `json:"saldo_fisico"`
	Saldovirtual               NullableInt `json:"saldo_virtual"`
}

type ProductWrapper struct {
	Produto Product `json:"produto"`
}

type ProductPayload struct {
	ID                         int64         `json:"id"`
	Idprodutopai               int64         `json:"id_produto_pai"`
	Nome                       string        `json:"nome"`
	Codigo                     string        `json:"codigo"`
	Preco                      float64       `json:"preco"`
	ImagemUrl                  string        `json:"imagem_url"`
	Tipo                       string        `json:"tipo"`
	Situacao                   string        `json:"situacao"`
	Formato                    string        `json:"formato"`
	DescricaoCurta             string        `json:"descricao_curta"`
	Datavalidade               time.Time     `json:"data_validade"`
	Unidade                    string        `json:"unidade"`
	Pesoliquido                float64       `json:"peso_liquido"`
	Pesobruto                  float64       `json:"peso_bruto"`
	Volumes                    int32         `json:"volumes"`
	Itensporcaixa              int32         `json:"itens_por_caixa"`
	Gtin                       string        `json:"gtin"`
	Gtinembalagem              string        `json:"gtin_embalagem"`
	Tipoproducao               string        `json:"tipo_producao"`
	Condicao                   int32         `json:"condicao"`
	Fretegratis                bool          `json:"frete_gratis"`
	Marca                      string        `json:"marca"`
	Descricaocomplementar      string        `json:"descricao_complementar"`
	Linkexterno                string        `json:"link_externo"`
	Observacoes                string        `json:"observacoes"`
	Descricaoembalagemdiscreta string        `json:"descricao_embalagem_discreta"`
	CreatedAt                  time.Time     `json:"created_at"`
	UpdatedAt                  time.Time     `json:"updated_at"`
	Saldofisicototal           sql.NullInt32 `json:"saldo_fisico_total"`
	Saldovirtualtotal          sql.NullInt32 `json:"saldo_virtual_total"`
	Saldofisico                sql.NullInt32 `json:"saldo_fisico"`
	Saldovirtual               sql.NullInt32 `json:"saldo_virtual"`
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

func (ni NullableInt) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Value)
}

func NewNullableInt(ni sql.NullInt32) NullableInt {
	return NullableInt{
		Value: ni.Int32,
		Valid: ni.Valid,
	}
}
