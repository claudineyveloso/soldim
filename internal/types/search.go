package types

import (
	"time"

	"github.com/google/uuid"
)

type Search struct {
	ID          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SearchPayload struct {
	ID          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SearchStore interface {
	CreateSearch(SearchPayload) (uuid.UUID, error)
	GetSearches() ([]*Search, error)
	GetSearchByID(id uuid.UUID) (*Search, error)
	DeleteSearch(id uuid.UUID) error
}
