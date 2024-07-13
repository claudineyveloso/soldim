package stock

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

func (s *Store) CreateStock(stock types.Stock) error {
	queries := db.New(s.db)
	ctx := context.Background()

	now := time.Now()
	stock.CreatedAt = now
	stock.UpdatedAt = now

	createStockParams := db.CreateStockParams{
		ID:                stock.ID,
		ProductID:         stock.ProductID,
		Saldofisicototal:  stock.Saldofisicototal,
		Saldovirtualtotal: stock.Saldovirtualtotal,
		CreatedAt:         stock.CreatedAt,
		UpdatedAt:         stock.UpdatedAt,
	}

	if err := queries.CreateStock(ctx, createStockParams); err != nil {
		fmt.Println("Erro ao criar um Deposito:", err)
		return err
	}
	return nil
}
