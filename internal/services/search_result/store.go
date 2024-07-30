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

func (s *Store) GetSearchesResult() ([]*types.SearchResult, error) {
	queries := db.New(s.db)
	ctx := context.Background()
	dbSearches, err := queries.GetSearchesResult(ctx)
	if err != nil {
		return nil, err
	}

	var searches []*types.SearchResult
	for _, dbSearch := range dbSearches {
		search := convertDBSearchResultToSearchResult(dbSearch)
		searches = append(searches, search)
	}
	return searches, nil
}

func (s *Store) GetSearchResultByID(searchID uuid.UUID) (*types.SearchResult, error) {
	queries := db.New(s.db)
	ctx := context.Background()
	dbSearch, err := queries.GetSearchResult(ctx, searchID)
	if err != nil {
		return nil, err
	}
	search := convertDBSearchResultToSearchResult(dbSearch)

	return search, nil
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

func (s *Store) GetSearchResultSources(searchID uuid.UUID) ([]*types.GetSearchResultSources, error) {
	queries := db.New(s.db)
	ctx := context.Background()

	// Chamada direta para a função com o searchID
	dbSearchResultSources, err := queries.GetSearchResultSources(ctx, searchID)
	if err != nil {
		return nil, err
	}

	var searchResultSources []*types.GetSearchResultSources
	for _, dbSearchResultSource := range dbSearchResultSources {
		searchResultSource := &types.GetSearchResultSources{
			Source:   dbSearchResultSource.Source,
			SearchID: dbSearchResultSource.SearchID,
		}
		searchResultSources = append(searchResultSources, searchResultSource)
	}
	return searchResultSources, nil
}

func convertDBSearchResultToSearchResult(dbSearchResult db.SearchesResult) *types.SearchResult {
	searchresult := &types.SearchResult{
		ID:          dbSearchResult.ID,
		ImageURL:    dbSearchResult.ImageUrl,
		Description: dbSearchResult.Description,
		Source:      dbSearchResult.Source,
		Price:       dbSearchResult.Price,
		Promotion:   dbSearchResult.Promotion,
		Link:        dbSearchResult.Link,
		SearchID:    dbSearchResult.SearchID,
		CreatedAt:   dbSearchResult.CreatedAt,
		UpdatedAt:   dbSearchResult.UpdatedAt,
	}
	return searchresult
}
