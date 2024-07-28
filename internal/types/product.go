package types

import (
	"database/sql"
	"encoding/json"
	"time"
)

type NullableString struct {
	sql.NullString
}

func (n NullableString) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(n.String)
}

type NullableInt struct {
	sql.NullInt32
}

type NullableInt64 struct {
	sql.NullInt64
}

// NullableFloat represents a float64 that may be null.
type NullableFloat struct {
	sql.NullFloat64
}

type NullableTime struct {
	sql.NullTime
}

func (n NullableTime) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return json.Marshal(nil)
	}
	// Format the time as you prefer, for example, RFC3339 format
	return json.Marshal(n.Time.Format(time.RFC3339))
}

type GetTotalProductNoMovementsParams struct {
	Column1 string `json:"column1"`
	Column2 string `json:"column2"`
}

type Product struct {
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
	SaldoFisicoTotal           NullableInt   `json:"saldo_fisico_total"`
	SaldoVirtualTotal          NullableInt   `json:"saldo_virtual_total"`
	SaldoFisico                NullableInt   `json:"saldo_fisico"`
	SaldoVirtual               NullableInt   `json:"saldo_virtual"`
	PrecoCusto                 NullableFloat `json:"preco_custo"`
	PrecoCompra                NullableFloat `json:"preco_compra"`
	SupplierID                 NullableInt64 `json:"supplier_id"`
}

type ProductEmptyStock struct {
	ID                         int64         `json:"id"`
	Idprodutopai               int64         `json:"idprodutopai"`
	Nome                       string        `json:"nome"`
	Codigo                     string        `json:"codigo"`
	Preco                      float64       `json:"preco"`
	Tipo                       string        `json:"tipo"`
	Situacao                   string        `json:"situacao"`
	Formato                    string        `json:"formato"`
	DescricaoCurta             string        `json:"descricao_curta"`
	ImagemUrl                  string        `json:"imagem_url"`
	Datavalidade               time.Time     `json:"datavalidade"`
	Unidade                    string        `json:"unidade"`
	Pesoliquido                float64       `json:"pesoliquido"`
	Pesobruto                  float64       `json:"pesobruto"`
	Volumes                    int32         `json:"volumes"`
	Itensporcaixa              int32         `json:"itensporcaixa"`
	Gtin                       string        `json:"gtin"`
	Gtinembalagem              string        `json:"gtinembalagem"`
	Tipoproducao               string        `json:"tipoproducao"`
	Condicao                   int32         `json:"condicao"`
	Fretegratis                bool          `json:"fretegratis"`
	Marca                      string        `json:"marca"`
	Descricaocomplementar      string        `json:"descricaocomplementar"`
	Linkexterno                string        `json:"linkexterno"`
	Observacoes                string        `json:"observacoes"`
	Descricaoembalagemdiscreta string        `json:"descricaoembalagemdiscreta"`
	CreatedAt                  time.Time     `json:"created_at"`
	UpdatedAt                  time.Time     `json:"updated_at"`
	SaldoFisicoTotal           NullableInt   `json:"saldo_fisico_total"`
	SaldoVirtualTotal          NullableInt   `json:"saldo_virtual_total"`
	SaldoFisico                NullableInt   `json:"saldo_fisico"`
	SaldoVirtual               NullableInt   `json:"saldo_virtual"`
	PrecoCusto                 NullableFloat `json:"preco_custo"`
	PrecoCompra                NullableFloat `json:"preco_compra"`
	SupplierID                 NullableInt64 `json:"supplier_id"`
}

type ProductNoMovements struct {
	SalesOrderID               int64          `json:"sales_order_id"`
	ProductID                  int64          `json:"product_id"`
	Quantidade                 int32          `json:"quantidade"`
	ID                         NullableInt64  `json:"id"`
	Nome                       NullableString `json:"nome"`
	Codigo                     NullableString `json:"codigo"`
	Preco                      NullableFloat  `json:"preco"`
	Tipo                       NullableString `json:"tipo"`
	Situacao                   NullableString `json:"situacao"`
	Formato                    NullableString `json:"formato"`
	DescricaoCurta             NullableString `json:"descricao_curta"`
	ImagemUrl                  NullableString `json:"imagem_url"`
	Datavalidade               NullableTime   `json:"datavalidade"`
	Unidade                    NullableString `json:"unidade"`
	Pesoliquido                NullableFloat  `json:"pesoliquido"`
	Pesobruto                  NullableFloat  `json:"pesobruto"`
	Volumes                    NullableInt    `json:"volumes"`
	Itensporcaixa              NullableInt    `json:"itensporcaixa"`
	Gtin                       NullableString `json:"gtin"`
	Gtinembalagem              NullableString `json:"gtinembalagem"`
	Tipoproducao               NullableString `json:"tipoproducao"`
	Condicao                   NullableInt    `json:"condicao"`
	Fretegratis                sql.NullBool   `json:"fretegratis"`
	Marca                      NullableString `json:"marca"`
	Descricaocomplementar      NullableString `json:"descricaocomplementar"`
	Linkexterno                NullableString `json:"linkexterno"`
	Observacoes                NullableString `json:"observacoes"`
	Descricaoembalagemdiscreta NullableString `json:"descricaoembalagemdiscreta"`
	Numero                     NullableInt    `json:"numero"`
	Numeroloja                 NullableString `json:"numeroloja"`
	Data                       NullableTime   `json:"data"`
	Datasaida                  NullableTime   `json:"datasaida"`
	Dataprevista               NullableTime   `json:"dataprevista"`
	Totalprodutos              NullableFloat  `json:"totalprodutos"`
	Totaldescontos             NullableFloat  `json:"totaldescontos"`
	Descricao                  NullableString `json:"descricao"`
	Codigo_2                   NullableInt64  `json:"codigo_2"`
	PrecoCusto                 NullableFloat  `json:"preco_custo"`
	PrecoCompra                NullableFloat  `json:"preco_compra"`
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
	SaldoFisicoTotal           NullableInt   `json:"saldo_fisico_total"`
	SaldoVirtualTotal          NullableInt   `json:"saldo_virtual_total"`
	SaldoFisico                NullableInt   `json:"saldo_fisico"`
	SaldoVirtual               NullableInt   `json:"saldo_virtual"`
	PrecoCusto                 NullableFloat `json:"preco_custo"`
	PrecoCompra                NullableFloat `json:"preco_compra"`
	SupplierID                 NullableInt64 `json:"supplier_id"`
}

type ProductResponse struct {
	Data  []Product `json:"data"`
	Total int       `json:"total"`
	Limit int       `json:"limit"`
}

type ProductStore interface {
	CreateProduct(ProductPayload) error
	GetProducts(nome, situacao string) ([]*Product, error)
	GetProductByID(id int64) (*Product, error)
	GetProductNoMovements(nome, situacao string) ([]*ProductNoMovements, error)
	GetProductEmptyStock(nome, situacao string) ([]*ProductEmptyStock, error)
	UpdateProduct(ProductPayload) error
	DeleteProduct(id int64) error
}

// MarshalJSON customiza a serialização JSON para NullableInt.
func (ni NullableInt) MarshalJSON() ([]byte, error) {
	if ni.Valid {
		return json.Marshal(ni.Int32)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON customiza a desserialização JSON para NullableInt.
func (ni *NullableInt) UnmarshalJSON(data []byte) error {
	// Tenta desserializar como um inteiro primeiro
	var intValue int32
	if err := json.Unmarshal(data, &intValue); err == nil {
		ni.Int32 = intValue
		ni.Valid = true
		return nil
	}

	// Tenta desserializar como um null
	var nullValue interface{}
	if err := json.Unmarshal(data, &nullValue); err == nil {
		if nullValue == nil {
			ni.Valid = false
			return nil
		}
	}

	// Se falhar em ambos os casos, retorna um erro
	return json.Unmarshal(data, &ni.NullInt32)
}

// MarshalJSON for NullableInt64
func (ni NullableInt64) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Int64)
}

// UnmarshalJSON for NullableInt64
func (ni *NullableInt64) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		ni.Valid = false
		return nil
	}
	err := json.Unmarshal(b, &ni.Int64)
	if err == nil {
		ni.Valid = true
	}
	return err
}

// MarshalJSON for NullableFloat
func (nf NullableFloat) MarshalJSON() ([]byte, error) {
	if !nf.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nf.Float64)
}

// UnmarshalJSON for NullableFloat
func (nf *NullableFloat) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		nf.Valid = false
		return nil
	}
	err := json.Unmarshal(b, &nf.Float64)
	if err == nil {
		nf.Valid = true
	}
	return err
}
