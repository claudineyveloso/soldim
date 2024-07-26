package salesorder

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

func (s *Store) CreateSalesOrder(salesorder types.SalesOrder) error {
	queries := db.New(s.db)
	ctx := context.Background()

	now := time.Now()
	salesorder.CreatedAt = now
	salesorder.UpdatedAt = now
	createSalesOrderParams := db.CreateSalesOrderParams{
		ID:             salesorder.ID,
		Numero:         salesorder.Numero,
		Numeroloja:     salesorder.Numeroloja,
		Data:           salesorder.Data.Time,
		Datasaida:      salesorder.Datasaida.Time,
		Dataprevista:   salesorder.Dataprevista.Time,
		Totalprodutos:  salesorder.Totalprodutos,
		Totaldescontos: salesorder.Totaldescontos,
		SituationID:    salesorder.SituationID,
		StoreID:        salesorder.StoreID,
		CreatedAt:      salesorder.CreatedAt,
		UpdatedAt:      salesorder.UpdatedAt,
	}

	fmt.Println("Criando um Pedido de Vendas...", createSalesOrderParams)

	if err := queries.CreateSalesOrder(ctx, createSalesOrderParams); err != nil {
		fmt.Println("Erro ao criar um pedido de vendas:", err)
		return err
	}
	return nil
}

func (s *Store) GetSalesOrders(limit, offset int32) ([]*types.SalesOrder, int64, error) {
	queries := db.New(s.db)
	ctx := context.Background()
	params := db.GetSalesOrdersParams{
		Limit:  limit,
		Offset: offset,
	}

	dbSalesOrders, err := queries.GetSalesOrders(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	totalCount, err := queries.GetTotalSalesOrders(ctx)
	if err != nil {
		return nil, 0, err
	}

	var salesorders []*types.SalesOrder
	for _, dbSalesOrder := range dbSalesOrders {
		salesorder := convertDBSalesOrderToSalesOrder(dbSalesOrder)
		salesorders = append(salesorders, salesorder)
	}
	return salesorders, totalCount, nil
}

func (s *Store) GetSalesOrderByID(salesorderID int64) (*types.SalesOrder, error) {
	queries := db.New(s.db)
	ctx := context.Background()
	dbSalesOrder, err := queries.GetSalesOrder(ctx, salesorderID)
	if err != nil {
		return nil, err
	}
	product := convertDBSalesOrderToSalesOrder(dbSalesOrder)

	return product, nil
}

func (s *Store) GetSalesOrderByNumber(number int32) (*types.SalesOrder, error) {
	queries := db.New(s.db)
	ctx := context.Background()
	dbSalesOrder, err := queries.GetSalesOrderByNumber(ctx, number)
	if err != nil {
		return nil, err
	}
	product := convertDBSalesOrderToSalesOrder(dbSalesOrder)

	return product, nil
}

func convertDBSalesOrderToSalesOrder(dbSalesOrder db.SalesOrder) *types.SalesOrder {
	salesorder := &types.SalesOrder{
		ID:             dbSalesOrder.ID,
		Numero:         dbSalesOrder.Numero,
		Numeroloja:     dbSalesOrder.Numeroloja,
		Data:           types.CustomDate{Time: dbSalesOrder.Data},
		Datasaida:      types.CustomDate{Time: dbSalesOrder.Datasaida},
		Dataprevista:   types.CustomDate{Time: dbSalesOrder.Dataprevista},
		Totalprodutos:  dbSalesOrder.Totalprodutos,
		Totaldescontos: dbSalesOrder.Totaldescontos,
		SituationID:    dbSalesOrder.SituationID,
		StoreID:        dbSalesOrder.StoreID,
		CreatedAt:      dbSalesOrder.CreatedAt,
		UpdatedAt:      dbSalesOrder.UpdatedAt,
	}
	return salesorder
}
