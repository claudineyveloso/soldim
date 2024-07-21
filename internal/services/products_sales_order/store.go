package productssalesorder

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

func (s *Store) CreateProductSalesOrder(productsalesorder types.ProductSalesOrderPayload) error {
	queries := db.New(s.db)
	ctx := context.Background()

	now := time.Now()
	productsalesorder.CreatedAt = now
	productsalesorder.UpdatedAt = now
	createProductSalesOrderParams := db.CreateProductSalesOrderParams{
		SalesOrderID: productsalesorder.SalesOrderID,
		ProductID:    productsalesorder.ProductID,
		Quantidade:   productsalesorder.Quantidade,
		CreatedAt:    productsalesorder.CreatedAt,
		UpdatedAt:    productsalesorder.UpdatedAt,
	}

	fmt.Println("Criando um Pedido de Vendas...", createProductSalesOrderParams)

	if err := queries.CreateProductSalesOrder(ctx, createProductSalesOrderParams); err != nil {
		fmt.Println("Erro ao criar um produto no pedido de vendas:", err)
		return err
	}
	return nil
}

func (s *Store) GetProductSalesOrders() ([]*types.ProductSalesOrder, error) {
	queries := db.New(s.db)
	ctx := context.Background()

	dbProductsSalesOrders, err := queries.GetProductSalesOrders(ctx)
	if err != nil {
		return nil, err
	}

	return convertDBProductsSalesOrders(dbProductsSalesOrders), nil
}

func (s *Store) GetProductSalesOrdersBySupplierID(supplierID int64) ([]*types.ProductSalesOrder, error) {
	queries := db.New(s.db)
	ctx := context.Background()

	dbSalesOrders, err := queries.GetProductSalesOrderBySupplierID(ctx, supplierID)
	if err != nil {
		return nil, err
	}

	return convertDBProductSalesOrdersBySupplierID(dbSalesOrders), nil
}

func convertDBProductsSalesOrderToProductsSalesOrder(row interface{}) *types.ProductSalesOrder {
	switch r := row.(type) {
	case db.GetProductSalesOrdersRow:
		return &types.ProductSalesOrder{
			ID:             types.ToInt64(r.ID),
			SalesOrderID:   r.SalesOrderID,
			ProductID:      r.ProductID,
			Quantidade:     r.Quantidade,
			CreatedAt:      r.CreatedAt,
			UpdatedAt:      r.UpdatedAt,
			Numero:         types.ToInt32(r.Numero),
			Numeroloja:     types.ToString(r.Numeroloja),
			Data:           types.ToTime(r.Data),
			Datasaida:      types.ToTime(r.Datasaida),
			Dataprevista:   types.ToTime(r.Dataprevista),
			Totalprodutos:  types.ToFloat64(r.Totalprodutos),
			Totaldescontos: types.ToFloat64(r.Totaldescontos),
			SituationID:    types.ToInt64(r.SituationID),
			StoreID:        types.ToInt64(r.StoreID),
			SupplierID:     types.ToInt64(r.SupplierID),
			Nome:           types.ToString(r.Nome),
			Codigo:         types.ToString(r.Codigo),
			Preco:          types.ToFloat64(r.Preco),
		}
	case db.GetProductSalesOrderBySupplierIDRow:
		return &types.ProductSalesOrder{
			ID:             types.ToInt64(r.ID),
			SalesOrderID:   r.SalesOrderID,
			ProductID:      r.ProductID,
			Quantidade:     r.Quantidade,
			CreatedAt:      r.CreatedAt,
			UpdatedAt:      r.UpdatedAt,
			Numero:         types.ToInt32(r.Numero),
			Numeroloja:     types.ToString(r.Numeroloja),
			Data:           types.ToTime(r.Data),
			Datasaida:      types.ToTime(r.Datasaida),
			Dataprevista:   types.ToTime(r.Dataprevista),
			Totalprodutos:  types.ToFloat64(r.Totalprodutos),
			Totaldescontos: types.ToFloat64(r.Totaldescontos),
			SituationID:    types.ToInt64(r.SituationID),
			StoreID:        types.ToInt64(r.StoreID),
			SupplierID:     types.ToInt64(r.SupplierID),
			Nome:           types.ToString(r.Nome),
			Codigo:         types.ToString(r.Codigo),
			Preco:          types.ToFloat64(r.Preco),
		}
	default:
		return nil
	}
}

func convertDBProductsSalesOrders(rows interface{}) []*types.ProductSalesOrder {
	var productssalesorders []*types.ProductSalesOrder
	switch rs := rows.(type) {
	case []db.GetProductSalesOrdersRow:
		for _, r := range rs {
			productsalesorder := convertDBProductsSalesOrderToProductsSalesOrder(r)
			productssalesorders = append(productssalesorders, productsalesorder)
		}
	case []db.GetProductSalesOrderBySupplierIDRow:
		for _, r := range rs {
			productsalesorder := convertDBProductsSalesOrderToProductsSalesOrder(r)
			productssalesorders = append(productssalesorders, productsalesorder)
		}
	}
	return productssalesorders
}

func convertDBProductSalesOrdersBySupplierID(rows []db.GetProductSalesOrderBySupplierIDRow) []*types.ProductSalesOrder {
	var productSalesOrders []*types.ProductSalesOrder
	for _, row := range rows {
		productSalesOrder := convertDBProductsSalesOrderToProductsSalesOrder(row)
		productSalesOrders = append(productSalesOrders, productSalesOrder)
	}
	return productSalesOrders
}
