package types

import (
	"time"

	"github.com/google/uuid"
)

type Search struct {
	ID          uuid.UUID `json:"id"`
	Description string    `json:"Description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SearchPayload struct {
	ID          uuid.UUID `json:"id"`
	Description string    `json:"Description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SearchStore interface {
	CreateSearch(SearchPayload) error
	GetSearches() ([]*Search, error)
	GetSearchByID(id uuid.UUID) (*Search, error)
	DeleteSearch(id uuid.UUID) error
}
