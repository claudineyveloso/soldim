package types

import (
	"time"

	"github.com/google/uuid"
)

type SearchResult struct {
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

type SearchResultPayload struct {
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

type GetSearchResultSources struct {
	Source   string    `json:"source"`
	SearchID uuid.UUID `json:"search_id"`
}

type SearchResultStore interface {
	CreateSearchResult(SearchResultPayload) error
	GetSearchesResult() ([]*SearchResult, error)
	GetSearchResultByID(id uuid.UUID) (*SearchResult, error)
	GetSearchResultSources(searchID uuid.UUID) ([]*GetSearchResultSources, error)
	DeleteSearchResult(id uuid.UUID) error
}
