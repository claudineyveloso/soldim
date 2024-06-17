package search

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

func (s *Store) CreateSearch(search types.SearchPayload) error {
	queries := db.New(s.db)
	ctx := context.Background()

	search.ID = uuid.New()
	now := time.Now()
	search.CreatedAt = now
	search.UpdatedAt = now

	createSearchParams := db.CreateSearchParams{
		ID:          search.ID,
		Description: search.Description,
		CreatedAt:   search.CreatedAt,
		UpdatedAt:   search.UpdatedAt,
	}

	if err := queries.CreateSearch(ctx, createSearchParams); err != nil {
		// http.Error(_, "Erro ao criar usuário", http.StatusInternalServerError)
		fmt.Println("Erro ao criar uma Busca:", err)
		return err
	}
	return nil
}

func (s *Store) GetSearches() ([]*types.Search, error) {
	queries := db.New(s.db)
	ctx := context.Background()

	dbSearches, err := queries.GetSearches(ctx)
	if err != nil {
		return nil, err
	}

	var searches []*types.Search
	for _, dbSearch := range dbSearches {
		search := convertDBSearchToSearch(dbSearch)
		searches = append(searches, search)
	}
	return searches, nil
}

func (s *Store) GetSearchByID(searchID uuid.UUID) (*types.Search, error) {
	queries := db.New(s.db)
	ctx := context.Background()
	dbSearch, err := queries.GetSearch(ctx, searchID)
	if err != nil {
		return nil, err
	}
	search := convertDBSearchToSearch(dbSearch)

	return search, nil
}

func (s *Store) DeleteSearch(searchID uuid.UUID) error {
	queries := db.New(s.db)
	ctx := context.Background()
	err := queries.DeleteSearch(ctx, searchID)
	if err != nil {
		return err
	}
	return nil
}

func convertDBSearchToSearch(dbSearch db.Search) *types.Search {
	search := &types.Search{
		ID:          dbSearch.ID,
		Description: dbSearch.Description,
		CreatedAt:   dbSearch.CreatedAt,
		UpdatedAt:   dbSearch.UpdatedAt,
	}
	return search
}
