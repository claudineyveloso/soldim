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
		ProductID:         stock.ProductID,
		SaldoFisicoTotal:  stock.SaldoFisicoTotal,
		SaldoVirtualTotal: stock.SaldoVirtualTotal,
		CreatedAt:         stock.CreatedAt,
		UpdatedAt:         stock.UpdatedAt,
	}

	if err := queries.CreateStock(ctx, createStockParams); err != nil {
		fmt.Println("Erro ao criar um Estoque:", err)
		return err
	}
	return nil
}

func (s *Store) UpdateStock(stock types.Stock) error {
	queries := db.New(s.db)
	ctx := context.Background()

	now := time.Now()
	stock.UpdatedAt = now

	updateStockParams := db.UpdateStockParams{
		SaldoFisicoTotal:  stock.SaldoFisicoTotal,
		SaldoVirtualTotal: stock.SaldoVirtualTotal,
		UpdatedAt:         stock.UpdatedAt,
	}

	if err := queries.UpdateStock(ctx, updateStockParams); err != nil {
		fmt.Println("Erro ao criar um Estoque:", err)
		return err
	}
	return nil
}

func (s *Store) GetStockByProductID(productID int64) (*types.Stock, error) {
	queries := db.New(s.db)
	ctx := context.Background()
	dbStock, err := queries.GetStockByProductID(ctx, productID)
	if err != nil {
		return nil, err
	}
	stock := convertGetStocktToStock(dbStock)

	return stock, nil
}

func convertGetStocktToStock(dbStock db.Stock) *types.Stock {
	stock := &types.Stock{
		ProductID:         dbStock.ProductID,
		SaldoFisicoTotal:  dbStock.SaldoFisicoTotal,
		SaldoVirtualTotal: dbStock.SaldoVirtualTotal,
		CreatedAt:         dbStock.CreatedAt,
		UpdatedAt:         dbStock.UpdatedAt,
	}
	return stock
}
