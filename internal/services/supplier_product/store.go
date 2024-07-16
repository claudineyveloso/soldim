package supplierproduct

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

func (s *Store) CreateSupplierProduct(supplierproduct types.SupplierProduct) error {
	queries := db.New(s.db)
	ctx := context.Background()

	now := time.Now()
	supplierproduct.CreatedAt = now
	supplierproduct.UpdatedAt = now

	createSupplierProductParams := db.CreateSupplierProductParams{
		ID:          supplierproduct.ID,
		Descricao:   supplierproduct.Descricao,
		PrecoCusto:  supplierproduct.PrecoCusto,
		PrecoCompra: supplierproduct.PrecoCompra,
		Padrao:      supplierproduct.Padrao,
		SupplierID:  supplierproduct.SupplierID,
		ProductID:   supplierproduct.ProductID,
		CreatedAt:   supplierproduct.CreatedAt,
		UpdatedAt:   supplierproduct.UpdatedAt,
	}

	if err := queries.CreateSupplierProduct(ctx, createSupplierProductParams); err != nil {
		fmt.Println("Erro ao criar um Fornecedor do Produto:", err)
		return err
	}
	return nil
}

func (s *Store) UpdateSupplierProduct(supplierproduct types.SupplierProduct) error {
	queries := db.New(s.db)
	ctx := context.Background()

	now := time.Now()
	supplierproduct.UpdatedAt = now

	updateSupplierProductParams := db.UpdateSupplierProductParams{
		ID:          supplierproduct.ID,
		Descricao:   supplierproduct.Descricao,
		PrecoCusto:  supplierproduct.PrecoCusto,
		PrecoCompra: supplierproduct.PrecoCompra,
		Padrao:      supplierproduct.Padrao,
		SupplierID:  supplierproduct.SupplierID,
		ProductID:   supplierproduct.ProductID,
		UpdatedAt:   supplierproduct.UpdatedAt,
	}

	if err := queries.UpdateSupplierProduct(ctx, updateSupplierProductParams); err != nil {
		fmt.Println("Erro ao criar um Fornecedor do Produto:", err)
		return err
	}
	return nil
}
