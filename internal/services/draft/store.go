package draft

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

func (s *Store) CreateDraft(draft types.DraftPayload) error {
	queries := db.New(s.db)
	ctx := context.Background()

	draft.ID = uuid.New()
	now := time.Now()
	draft.CreatedAt = now
	draft.UpdatedAt = now

	createDraftParams := db.CreateDraftParams{
		ID:          draft.ID,
		ImageUrl:    draft.ImageURL,
		Source:      draft.Source,
		Price:       draft.Price,
		Description: draft.Description,
		Promotion:   draft.Promotion,
		Link:        draft.Link,
		SearchID:    draft.SearchID,
		CreatedAt:   draft.CreatedAt,
		UpdatedAt:   draft.UpdatedAt,
	}

	if err := queries.CreateDraft(ctx, createDraftParams); err != nil {
		fmt.Println("Erro ao criar um Rascunho:", err)
		return err
	}
	return nil
}

func (s *Store) GetDraft() ([]*types.Draft, error) {
	queries := db.New(s.db)
	ctx := context.Background()

	dbDrafts, err := queries.GetDrafts(ctx)
	if err != nil {
		return nil, err
	}

	var drafts []*types.Draft
	for _, dbDraft := range dbDrafts {
		draft := convertDBDraftToDraft(dbDraft)
		drafts = append(drafts, draft)
	}
	return drafts, nil
}

func (s *Store) UpdateDraft(draft types.DraftPayload) error {
	queries := db.New(s.db)
	ctx := context.Background()

	now := time.Now()
	draft.UpdatedAt = now

	updateDraftParams := db.UpdateDraftParams{
		ID:          draft.ID,
		ImageUrl:    draft.ImageURL,
		Source:      draft.Source,
		Price:       draft.Price,
		Description: draft.Description,
		Promotion:   draft.Promotion,
		Link:        draft.Link,
		UpdatedAt:   draft.UpdatedAt,
	}

	if err := queries.UpdateDraft(ctx, updateDraftParams); err != nil {
		fmt.Println("Erro ao atualizar um Rascunho:", err)
		return err
	}
	return nil
}

func (s *Store) GetDraftByID(draftID uuid.UUID) (*types.Draft, error) {
	queries := db.New(s.db)
	ctx := context.Background()
	dbDraft, err := queries.GetDraft(ctx, draftID)
	if err != nil {
		return nil, err
	}
	bucket := convertDBDraftToDraft(dbDraft)

	return bucket, nil
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

func convertDBDraftToDraft(dbDraft db.Draft) *types.Draft {
	draft := &types.Draft{
		ID:          dbDraft.ID,
		ImageURL:    dbDraft.ImageUrl,
		Description: dbDraft.Description,
		Source:      dbDraft.Source,
		Price:       dbDraft.Price,
		Promotion:   dbDraft.Promotion,
		Link:        dbDraft.Link,
		SearchID:    dbDraft.SearchID,
		CreatedAt:   dbDraft.CreatedAt,
		UpdatedAt:   dbDraft.UpdatedAt,
	}
	return draft
}
