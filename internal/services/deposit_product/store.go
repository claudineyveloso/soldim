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
		DepositID:    depositproduct.DepositID,
		ProductID:    depositproduct.ProductID,
		SaldoFisico:  depositproduct.SaldoFisico,
		SaldoVirtual: depositproduct.SaldoVirtual,
		CreatedAt:    depositproduct.CreatedAt,
		UpdatedAt:    depositproduct.UpdatedAt,
	}

	if err := queries.CreateDepositProduct(ctx, createDepositProductParams); err != nil {
		fmt.Println("Erro ao criar um Deposito por produto:", err)
		return err
	}
	return nil
}

func (s *Store) UpdateDepositProduct(depositproduct types.DepositProduct) error {
	queries := db.New(s.db)
	ctx := context.Background()

	now := time.Now()
	depositproduct.UpdatedAt = now

	updateDepositProductParams := db.UpdateDepositProductParams{
		DepositID:    depositproduct.DepositID,
		ProductID:    depositproduct.ProductID,
		SaldoFisico:  depositproduct.SaldoFisico,
		SaldoVirtual: depositproduct.SaldoVirtual,
		UpdatedAt:    depositproduct.UpdatedAt,
	}

	if err := queries.UpdateDepositProduct(ctx, updateDepositProductParams); err != nil {
		fmt.Println("Erro ao atualizar um Rascunho:", err)
		return err
	}
	return nil
}

func (s *Store) GetDepositProductByProductID(productID int64) ([]*types.DepositProduct, error) {
	queries := db.New(s.db)
	ctx := context.Background()
	dbDepositProducts, err := queries.GetDepositProductsByProductID(ctx, productID)
	if err != nil {
		return nil, err
	}

	// Converter todos os produtos
	var products []*types.DepositProduct
	for _, dbDepositProduct := range dbDepositProducts {
		product := convertGetDepositProductRowToDepositProduct(dbDepositProduct)
		products = append(products, product)
	}

	return products, nil
}

func (s *Store) GetDepositProductByDepositID(depositID int64) ([]*types.DepositProduct, error) {
	queries := db.New(s.db)
	ctx := context.Background()
	dbDepositProducts, err := queries.GetDepositProductsByDepositID(ctx, depositID)
	if err != nil {
		return nil, err
	}

	// Converter todos os produtos
	var products []*types.DepositProduct
	for _, dbDepositProduct := range dbDepositProducts {
		product := convertGetDepositProductRowToDepositProduct(dbDepositProduct)
		products = append(products, product)
	}

	return products, nil
}

func convertGetDepositProductRowToDepositProduct(dbDepositProduct db.DepositProduct) *types.DepositProduct {
	DepositProduct := &types.DepositProduct{
		DepositID:    dbDepositProduct.DepositID,
		ProductID:    dbDepositProduct.ProductID,
		SaldoFisico:  dbDepositProduct.SaldoFisico,
		SaldoVirtual: dbDepositProduct.SaldoVirtual,
		CreatedAt:    dbDepositProduct.CreatedAt,
		UpdatedAt:    dbDepositProduct.UpdatedAt,
	}
	return DepositProduct
}
