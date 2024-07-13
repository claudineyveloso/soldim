package depositproduct

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

func (s *Store) CreateDepositProduct(depositproduct types.DepositProduct) error {
	queries := db.New(s.db)
	ctx := context.Background()

	now := time.Now()
	depositproduct.CreatedAt = now
	depositproduct.UpdatedAt = now

	createDepositProductParams := db.CreateDepositProductParams{
		ID:           depositproduct.ID,
		DepositID:    depositproduct.DepositID,
		ProductID:    depositproduct.ProductID,
		Saldofisico:  depositproduct.Saldofisico,
		Saldovirtual: depositproduct.Saldovirtual,
		CreatedAt:    depositproduct.CreatedAt,
		UpdatedAt:    depositproduct.UpdatedAt,
	}

	if err := queries.CreateDepositProduct(ctx, createDepositProductParams); err != nil {
		fmt.Println("Erro ao criar um Deposito por produto:", err)
		return err
	}
	return nil
}
