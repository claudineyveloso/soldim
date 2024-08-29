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
		ID:                         draft.ID,
		Codigo:                     draft.Codigo,
		Tipo:                       draft.Tipo,
		Situacao:                   draft.Situacao,
		Formato:                    draft.Formato,
		ImageUrl:                   draft.ImageUrl,
		Description:                draft.Description,
		Datavalidade:               draft.Datavalidade,
		Unidade:                    draft.Unidade,
		Pesoliquido:                draft.Pesoliquido,
		Pesobruto:                  draft.Pesobruto,
		Volumes:                    draft.Volumes,
		Itensporcaixa:              draft.Itensporcaixa,
		Gtin:                       draft.Gtin,
		Gtinembalagem:              draft.Gtinembalagem,
		Tipoproducao:               draft.Tipoproducao,
		Condicao:                   draft.Condicao,
		Fretegratis:                draft.Fretegratis,
		Marca:                      draft.Marca,
		Descricaocomplementar:      draft.Descricaocomplementar,
		Linkexterno:                draft.Linkexterno,
		Observacoes:                draft.Observacoes,
		Descricaoembalagemdiscreta: draft.Descricaoembalagemdiscreta,
		Source:                     draft.Source,
		Price:                      draft.Price,
		Promotion:                  draft.Promotion,
		Link:                       draft.Link,
		SearchID:                   draft.SearchID,
		CreatedAt:                  draft.CreatedAt,
		UpdatedAt:                  draft.UpdatedAt,
	}

	if err := queries.CreateDraft(ctx, createDraftParams); err != nil {
		fmt.Println("Erro ao criar um Rascunho:", err)
		return err
	}
	return nil
}

func (s *Store) GetDrafts() ([]*types.Draft, error) {
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
		ID:                         draft.ID,
		Codigo:                     draft.Codigo,
		Tipo:                       draft.Tipo,
		Situacao:                   draft.Situacao,
		Formato:                    draft.Formato,
		ImageUrl:                   draft.ImageUrl,
		Description:                draft.Description,
		Datavalidade:               draft.Datavalidade,
		Unidade:                    draft.Unidade,
		Pesoliquido:                draft.Pesoliquido,
		Pesobruto:                  draft.Pesobruto,
		Volumes:                    draft.Volumes,
		Itensporcaixa:              draft.Itensporcaixa,
		Gtin:                       draft.Gtin,
		Gtinembalagem:              draft.Gtinembalagem,
		Tipoproducao:               draft.Tipoproducao,
		Condicao:                   draft.Condicao,
		Fretegratis:                draft.Fretegratis,
		Marca:                      draft.Marca,
		Descricaocomplementar:      draft.Descricaocomplementar,
		Linkexterno:                draft.Linkexterno,
		Observacoes:                draft.Observacoes,
		Descricaoembalagemdiscreta: draft.Descricaoembalagemdiscreta,
		Source:                     draft.Source,
		Price:                      draft.Price,
		Promotion:                  draft.Promotion,
		Link:                       draft.Link,
		SearchID:                   draft.SearchID,
		UpdatedAt:                  draft.UpdatedAt,
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
	draft := convertDBDraftToDraft(dbDraft)

	return draft, nil
}

func (s *Store) GetDraftBySearchID(searchID uuid.UUID) (*types.Draft, error) {
	queries := db.New(s.db)
	ctx := context.Background()
	dbDraft, err := queries.GetDraftBySearchId(ctx, searchID)
	if err != nil {
		return nil, err
	}
	draft := convertDBDraftToDraft(dbDraft)

	return draft, nil
}

func (s *Store) DeleteDraft(draftID uuid.UUID) error {
	queries := db.New(s.db)
	ctx := context.Background()
	err := queries.DeleteDraft(ctx, draftID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) DeleteDraftBySearchID(searchID uuid.UUID) error {
	queries := db.New(s.db)
	ctx := context.Background()
	err := queries.DeleteDraftBySearchID(ctx, searchID)
	if err != nil {
		return err
	}
	return nil
}

func convertDBDraftToDraft(dbDraft db.Draft) *types.Draft {
	draft := &types.Draft{
		ID:                         dbDraft.ID,
		Codigo:                     dbDraft.Codigo,
		Tipo:                       dbDraft.Tipo,
		Situacao:                   dbDraft.Situacao,
		Formato:                    dbDraft.Formato,
		ImageUrl:                   dbDraft.ImageUrl,
		Description:                dbDraft.Description,
		Datavalidade:               dbDraft.Datavalidade,
		Unidade:                    dbDraft.Unidade,
		Pesoliquido:                dbDraft.Pesoliquido,
		Pesobruto:                  dbDraft.Pesobruto,
		Volumes:                    dbDraft.Volumes,
		Itensporcaixa:              dbDraft.Itensporcaixa,
		Gtin:                       dbDraft.Gtin,
		Gtinembalagem:              dbDraft.Gtinembalagem,
		Tipoproducao:               dbDraft.Tipoproducao,
		Condicao:                   dbDraft.Condicao,
		Fretegratis:                dbDraft.Fretegratis,
		Marca:                      dbDraft.Marca,
		Descricaocomplementar:      dbDraft.Descricaocomplementar,
		Linkexterno:                dbDraft.Linkexterno,
		Observacoes:                dbDraft.Observacoes,
		Descricaoembalagemdiscreta: dbDraft.Descricaoembalagemdiscreta,
		Source:                     dbDraft.Source,
		Price:                      dbDraft.Price,
		Promotion:                  dbDraft.Promotion,
		Link:                       dbDraft.Link,
		SearchID:                   dbDraft.SearchID,
		CreatedAt:                  dbDraft.CreatedAt,
		UpdatedAt:                  dbDraft.UpdatedAt,
	}
	return draft
}
