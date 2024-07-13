package deposit

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

func (s *Store) CreateDeposit(deposit types.Deposit) error {
	queries := db.New(s.db)
	ctx := context.Background()

	now := time.Now()
	deposit.CreatedAt = now
	deposit.UpdatedAt = now

	createDepositParams := db.CreateDepositParams{
		ID:                 deposit.ID,
		Descricao:          deposit.Descricao,
		Situacao:           deposit.Situacao,
		Padrao:             deposit.Padrao,
		Desconsiderarsaldo: deposit.Desconsiderarsaldo,
		CreatedAt:          deposit.CreatedAt,
		UpdatedAt:          deposit.UpdatedAt,
	}

	if err := queries.CreateDeposit(ctx, createDepositParams); err != nil {
		fmt.Println("Erro ao criar um Deposito:", err)
		return err
	}
	return nil
}

func (s *Store) GetDeposits() ([]*types.Deposit, error) {
	queries := db.New(s.db)
	ctx := context.Background()

	dbDeposits, err := queries.GetDeposits(ctx)
	if err != nil {
		return nil, err
	}

	var deposits []*types.Deposit
	for _, dbDeposit := range dbDeposits {
		deposit := convertDBDepositToDeposit(dbDeposit)
		deposits = append(deposits, deposit)
	}
	return deposits, nil
}

func convertDBDepositToDeposit(dbDeposit db.Deposit) *types.Deposit {
	deposit := &types.Deposit{
		ID:                 dbDeposit.ID,
		Descricao:          dbDeposit.Descricao,
		Situacao:           dbDeposit.Situacao,
		Padrao:             dbDeposit.Padrao,
		Desconsiderarsaldo: dbDeposit.Desconsiderarsaldo,
		CreatedAt:          dbDeposit.CreatedAt,
		UpdatedAt:          dbDeposit.UpdatedAt,
	}
	return deposit
}
