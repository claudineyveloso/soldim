package types

import (
	"time"

	"github.com/google/uuid"
)

type Draft struct {
	ID          uuid.UUID `json:"id"`
	ImageURL    string    `json:"image_url"`
	Description string    `json:"description"`
	Source      string    `json:"source"`
	Price       float64   `json:"price"`
	Promotion   bool      `json:"promotion"`
	Link        string    `json:"link"`
	SearchID    uuid.UUID `json:"search_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type DraftPayload struct {
	ID          uuid.UUID `json:"id"`
	ImageURL    string    `json:"image_url"`
	Description string    `json:"description"`
	Source      string    `json:"source"`
	Price       float64   `json:"price"`
	Promotion   bool      `json:"promotion"`
	Link        string    `json:"link"`
	SearchID    uuid.UUID `json:"search_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type DraftStore interface {
	CreateDraft(DraftPayload) error
	GetDrafts() ([]*Draft, error)
	GetDraftByID(id uuid.UUID) (*Draft, error)
	UpdateDraft(DraftPayload) error
	DeleteDraft(id uuid.UUID) error
}
