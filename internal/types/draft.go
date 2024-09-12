package types

import (
	"time"

	"github.com/google/uuid"
)

type Draft struct {
	ID                         uuid.UUID `json:"id"`
	Codigo                     string    `json:"codigo"`
	Tipo                       string    `json:"tipo"`
	Situacao                   string    `json:"situacao"`
	Formato                    string    `json:"formato"`
	ImageUrl                   string    `json:"image_url"`
	Description                string    `json:"description"`
	Datavalidade               time.Time `json:"datavalidade"`
	Unidade                    string    `json:"unidade"`
	Pesoliquido                float64   `json:"pesoliquido"`
	Pesobruto                  float64   `json:"pesobruto"`
	Volumes                    int32     `json:"volumes"`
	Itensporcaixa              int32     `json:"itensporcaixa"`
	Gtin                       string    `json:"gtin"`
	Gtinembalagem              string    `json:"gtinEmbalagem"`
	Tipoproducao               string    `json:"tipoproducao"`
	Condicao                   int32     `json:"condicao"`
	Fretegratis                bool      `json:"fretegratis"`
	Marca                      string    `json:"marca"`
	Descricaocomplementar      string    `json:"descricaocomplementar"`
	Linkexterno                string    `json:"linkexterno"`
	Observacoes                string    `json:"observacoes"`
	Descricaoembalagemdiscreta string    `json:"descricaoembalagemdiscreta"`
	Source                     string    `json:"source"`
	Price                      float64   `json:"price"`
	Promotion                  bool      `json:"promotion"`
	Link                       string    `json:"link"`
	SearchID                   uuid.UUID `json:"search_id"`
	CreatedAt                  time.Time `json:"created_at"`
	UpdatedAt                  time.Time `json:"updated_at"`
}

type DraftPayload struct {
	ID                         uuid.UUID `json:"id"`
	Codigo                     string    `json:"codigo"`
	Tipo                       string    `json:"tipo"`
	Situacao                   string    `json:"situacao"`
	Formato                    string    `json:"formato"`
	ImageUrl                   string    `json:"image_url"`
	Description                string    `json:"description"`
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
	Source                     string    `json:"source"`
	Price                      float64   `json:"price"`
	Promotion                  bool      `json:"promotion"`
	Link                       string    `json:"link"`
	SearchID                   uuid.UUID `json:"search_id"`
	CreatedAt                  time.Time `json:"created_at"`
	UpdatedAt                  time.Time `json:"updated_at"`
}

type DraftStore interface {
	CreateDraft(DraftPayload) error
	GetDrafts() ([]*Draft, error)
	GetDraftByID(id uuid.UUID) (*Draft, error)
	GetDraftBySearchID(id uuid.UUID) (*Draft, error)
	UpdateDraft(DraftPayload) error
	DeleteDraft(id uuid.UUID) error
	DeleteDraftBySearchID(id uuid.UUID) error
}
