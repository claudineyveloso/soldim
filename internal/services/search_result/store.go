package searchresult

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/claudineyveloso/soldim.git/internal/db"
	"github.com/claudineyveloso/soldim.git/internal/types"
	"github.com/google/uuid"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateSearchResult(searchresult types.SearchResultPayload) error {
	queries := db.New(s.db)
	ctx := context.Background()

	searchresult.ID = uuid.New()
	now := time.Now()
	searchresult.CreatedAt = now
	searchresult.UpdatedAt = now

	createSearchResultParams := db.CreateSearchResultParams{
		ID:          searchresult.ID,
		ImageUrl:    searchresult.ImageURL,
		Source:      searchresult.Source,
		Price:       searchresult.Price,
		Description: searchresult.Description,
		Promotion:   searchresult.Promotion,
		Link:        searchresult.Link,
		SearchID:    searchresult.SearchID,
		CreatedAt:   searchresult.CreatedAt,
		UpdatedAt:   searchresult.UpdatedAt,
	}

	if err := queries.CreateSearchResult(ctx, createSearchResultParams); err != nil {
		fmt.Println("Erro ao criar um Resultado da Busca:", err)
		return err
	}
	return nil
}

func (s *Store) DeleteSearchResult(searchID uuid.UUID) error {
	queries := db.New(s.db)
	ctx := context.Background()
	err := queries.DeleteSearchResult(ctx, searchID)
	if err != nil {
		return err
	}
	return nil
}
