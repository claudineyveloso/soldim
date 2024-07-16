// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package db

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Address struct {
	ID              uuid.UUID      `json:"id"`
	PublicPlace     sql.NullString `json:"public_place"`
	Complement      sql.NullString `json:"complement"`
	Neighborhood    sql.NullString `json:"neighborhood"`
	City            sql.NullString `json:"city"`
	State           sql.NullString `json:"state"`
	ZipCode         sql.NullString `json:"zip_code"`
	AddressableID   uuid.UUID      `json:"addressable_id"`
	AddressableType string         `json:"addressable_type"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

type Deposit struct {
	ID                 int64     `json:"id"`
	Descricao          string    `json:"descricao"`
	Situacao           int32     `json:"situacao"`
	Padrao             bool      `json:"padrao"`
	Desconsiderarsaldo bool      `json:"desconsiderarsaldo"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type DepositProduct struct {
	DepositID    int64     `json:"deposit_id"`
	ProductID    int64     `json:"product_id"`
	SaldoFisico  int32     `json:"saldo_fisico"`
	SaldoVirtual int32     `json:"saldo_virtual"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Draft struct {
	ID          uuid.UUID `json:"id"`
	ImageUrl    string    `json:"image_url"`
	Description string    `json:"description"`
	Source      string    `json:"source"`
	Price       float64   `json:"price"`
	Promotion   bool      `json:"promotion"`
	Link        string    `json:"link"`
	SearchID    uuid.UUID `json:"search_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Parameter struct {
	ID                 uuid.UUID `json:"id"`
	DiscountPercentage int32     `json:"discount_percentage"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type Person struct {
	ID             uuid.UUID      `json:"id"`
	FirstName      string         `json:"first_name"`
	LastName       string         `json:"last_name"`
	Email          string         `json:"email"`
	Phone          sql.NullString `json:"phone"`
	CellPhone      string         `json:"cell_phone"`
	PersonableID   uuid.UUID      `json:"personable_id"`
	PersonableType string         `json:"personable_type"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

type Product struct {
	ID                         int64     `json:"id"`
	Idprodutopai               int64     `json:"idprodutopai"`
	Nome                       string    `json:"nome"`
	Codigo                     string    `json:"codigo"`
	Preco                      float64   `json:"preco"`
	Tipo                       string    `json:"tipo"`
	Situacao                   string    `json:"situacao"`
	Formato                    string    `json:"formato"`
	DescricaoCurta             string    `json:"descricao_curta"`
	ImagemUrl                  string    `json:"imagem_url"`
	Datavalidade               time.Time `json:"datavalidade"`
	Unidade                    string    `json:"unidade"`
	Pesoliquido                float64   `json:"pesoliquido"`
	Pesobruto                  float64   `json:"pesobruto"`
	Volumes                    int32     `json:"volumes"`
	Itensporcaixa              int32     `json:"itensporcaixa"`
	Gtin                       string    `json:"gtin"`
	Gtinembalagem              string    `json:"gtinembalagem"`
	Tipoproducao               string    `json:"tipoproducao"`
	Condicao                   int32     `json:"condicao"`
	Fretegratis                bool      `json:"fretegratis"`
	Marca                      string    `json:"marca"`
	Descricaocomplementar      string    `json:"descricaocomplementar"`
	Linkexterno                string    `json:"linkexterno"`
	Observacoes                string    `json:"observacoes"`
	Descricaoembalagemdiscreta string    `json:"descricaoembalagemdiscreta"`
	CreatedAt                  time.Time `json:"created_at"`
	UpdatedAt                  time.Time `json:"updated_at"`
}

type SalesChannel struct {
	ID        int32     `json:"id"`
	Descricao string    `json:"descricao"`
	Tipo      string    `json:"tipo"`
	Situacao  int32     `json:"situacao"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Search struct {
	ID          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SearchesResult struct {
	ID          uuid.UUID `json:"id"`
	ImageUrl    string    `json:"image_url"`
	Description string    `json:"description"`
	Source      string    `json:"source"`
	Price       float64   `json:"price"`
	Promotion   bool      `json:"promotion"`
	Link        string    `json:"link"`
	SearchID    uuid.UUID `json:"search_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Stock struct {
	ProductID         int64     `json:"product_id"`
	SaldoFisicoTotal  int32     `json:"saldo_fisico_total"`
	SaldoVirtualTotal int32     `json:"saldo_virtual_total"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type Supplier struct {
	ID          int64     `json:"id"`
	Descricao   string    `json:"descricao"`
	Codigo      string    `json:"codigo"`
	Precocusto  string    `json:"precocusto"`
	Precocompra string    `json:"precocompra"`
	Padrão      bool      `json:"padrão"`
	ProductID   int64     `json:"product_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SupplierProduct struct {
	ID          int64     `json:"id"`
	Descricao   string    `json:"descricao"`
	Codigo      int64     `json:"codigo"`
	PrecoCusto  float64   `json:"preco_custo"`
	PrecoCompra float64   `json:"preco_compra"`
	Padrao      bool      `json:"padrao"`
	SupplierID  int64     `json:"supplier_id"`
	ProductID   int64     `json:"product_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Token struct {
	ID           uuid.UUID `json:"id"`
	AccessToken  string    `json:"access_token"`
	ExpiresIn    int32     `json:"expires_in"`
	TokenType    string    `json:"token_type"`
	Scope        string    `json:"scope"`
	RefreshToken string    `json:"refresh_token"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	IsActive  bool      `json:"is_active"`
	UserType  string    `json:"user_type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
