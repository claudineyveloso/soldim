package situation

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/claudineyveloso/soldim.git/internal/db"
	"github.com/claudineyveloso/soldim.git/internal/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateSituation(situation types.Situation) error {
	queries := db.New(s.db)
	ctx := context.Background()

	now := time.Now()
	situation.CreatedAt = now
	situation.UpdatedAt = now

	createSituationParams := db.CreateSituationParams{
		ID:        situation.ID,
		Descricao: situation.Descricao,
		CreatedAt: situation.CreatedAt,
		UpdatedAt: situation.UpdatedAt,
	}

	if err := queries.CreateSituation(ctx, createSituationParams); err != nil {
		fmt.Println("Erro ao criar uma Situação:", err)
		return err
	}
	return nil
}

func (s *Store) GetSituations() ([]*types.Situation, error) {
	queries := db.New(s.db)
	ctx := context.Background()

	dbSituations, err := queries.GetSituations(ctx)
	if err != nil {
		return nil, err
	}

	var situations []*types.Situation
	for _, dbSituation := range dbSituations {
		situation := convertDBSituationToSituation(dbSituation)
		situations = append(situations, situation)
	}
	return situations, nil
}

func (s *Store) GetSituationByID(situationID int64) (*types.Situation, error) {
	queries := db.New(s.db)
	ctx := context.Background()
	dbSituation, err := queries.GetSituation(ctx, situationID)
	if err != nil {
		return nil, err
	}
	product := convertDBSituationToSituation(dbSituation)

	return product, nil
}

func (s *Store) GetSituationByDescription(description string) (*types.Situation, error) {
	queries := db.New(s.db)
	ctx := context.Background()
	dbSituation, err := queries.GetSituationByDescroption(ctx, description)
	if err != nil {
		return nil, err
	}
	product := convertDBSituationToSituation(dbSituation)

	return product, nil
}

func convertDBSituationToSituation(dbSituation db.Situation) *types.Situation {
	situation := &types.Situation{
		ID:        dbSituation.ID,
		CreatedAt: dbSituation.CreatedAt,
		UpdatedAt: dbSituation.UpdatedAt,
	}
	return situation
}
