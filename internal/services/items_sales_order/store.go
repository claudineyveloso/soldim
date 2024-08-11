package itemssalesorder

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

func (s *Store) CreateItemsSalesOrder(itemssalesorder types.ItemsSalesOrder) error {
	queries := db.New(s.db)
	ctx := context.Background()

	now := time.Now()
	itemssalesorder.CreatedAt = now
	itemssalesorder.UpdatedAt = now
	createItemsSalesOrderParams := db.CreateItemsSalesOrderParams{
		ID:                 itemssalesorder.ID,
		SalesOrderID:       itemssalesorder.SalesOrderID,
		Codigo:             itemssalesorder.Codigo,
		Unidade:            itemssalesorder.Unidade,
		Quantidade:         itemssalesorder.Quantidade,
		Desconto:           itemssalesorder.Desconto,
		Valor:              itemssalesorder.Valor,
		Aliquotaipi:        itemssalesorder.Aliquotaipi,
		Descricao:          itemssalesorder.Descricao,
		Descricaodetalhada: itemssalesorder.Descricaodetalhada,
		ProductID:          itemssalesorder.ProductID,
		CreatedAt:          itemssalesorder.CreatedAt,
		UpdatedAt:          itemssalesorder.UpdatedAt,
	}

	fmt.Println("Criando um Item do Pedido de Vendas...", createItemsSalesOrderParams)

	if err := queries.CreateItemsSalesOrder(ctx, createItemsSalesOrderParams); err != nil {
		fmt.Println("Erro ao criar um item do pedido de vendas:", err)
		return err
	}
	return nil
}

func (s *Store) GetItemsSalesOrderByID(itemssalesorderID int64) (*types.ItemsSalesOrder, error) {
	queries := db.New(s.db)
	ctx := context.Background()
	dbItemsSalesOrder, err := queries.GetItemsSalesOrderByID(ctx, itemssalesorderID)
	if err != nil {
		return nil, err
	}
	product := convertDBItemsSalesOrderToItemsSalesOrder(dbItemsSalesOrder)

	return product, nil
}

func (s *Store) GetItemsSalesOrderByProductID(itemssalesorderID int64) (*types.ItemsSalesOrder, error) {
	queries := db.New(s.db)
	ctx := context.Background()
	dbItemsSalesOrder, err := queries.GetItemsSalesOrderByProductID(ctx, itemssalesorderID)
	if err != nil {
		return nil, err
	}
	product := convertDBItemsSalesOrderToItemsSalesOrder(dbItemsSalesOrder)

	return product, nil
}

func (s *Store) GetItemsSalesOrderBySalesOrderID(itemssalesorderID int64) (*types.ItemsSalesOrder, error) {
	queries := db.New(s.db)
	ctx := context.Background()
	dbItemsSalesOrder, err := queries.GetItemsSalesOrderBySalesOrderID(ctx, itemssalesorderID)
	if err != nil {
		return nil, err
	}
	product := convertDBItemsSalesOrderToItemsSalesOrder(dbItemsSalesOrder)

	return product, nil
}

func convertDBItemsSalesOrderToItemsSalesOrder(dbItemsSalesOrder db.ItemsSalesOrder) *types.ItemsSalesOrder {
	return &types.ItemsSalesOrder{
		ID:                 dbItemsSalesOrder.ID,
		SalesOrderID:       dbItemsSalesOrder.SalesOrderID,
		Codigo:             dbItemsSalesOrder.Codigo,
		Unidade:            dbItemsSalesOrder.Unidade,
		Quantidade:         dbItemsSalesOrder.Quantidade,
		Desconto:           dbItemsSalesOrder.Desconto,
		Valor:              dbItemsSalesOrder.Valor,
		Aliquotaipi:        dbItemsSalesOrder.Aliquotaipi,
		Descricao:          dbItemsSalesOrder.Descricao,
		Descricaodetalhada: dbItemsSalesOrder.Descricaodetalhada,
		ProductID:          dbItemsSalesOrder.ProductID,
		CreatedAt:          dbItemsSalesOrder.CreatedAt,
		UpdatedAt:          dbItemsSalesOrder.UpdatedAt,
	}
}
