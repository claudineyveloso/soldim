package product

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

func (s *Store) CreateProduct(product types.ProductPayload) error {
	queries := db.New(s.db)
	ctx := context.Background()

	now := time.Now()
	product.CreatedAt = now
	product.UpdatedAt = now

	createProductParams := db.CreateProductParams{
		ID:                         product.ID,
		Nome:                       product.Nome,
		Codigo:                     product.Codigo,
		Preco:                      product.Preco,
		ImagemUrl:                  product.ImagemUrl,
		Tipo:                       product.Tipo,
		Situacao:                   product.Situacao,
		Formato:                    product.Formato,
		DescricaoCurta:             product.DescricaoCurta,
		Datavalidade:               product.Datavalidade,
		Unidade:                    product.Unidade,
		Pesoliquido:                product.Pesoliquido,
		Pesobruto:                  product.Pesobruto,
		Volumes:                    product.Volumes,
		Itensporcaixa:              product.Itensporcaixa,
		Gtin:                       product.Gtin,
		Gtinembalagem:              product.Gtinembalagem,
		Tipoproducao:               product.Tipoproducao,
		Condicao:                   product.Condicao,
		Fretegratis:                product.Fretegratis,
		Marca:                      product.Marca,
		Descricaocomplementar:      product.Descricaocomplementar,
		Linkexterno:                product.Linkexterno,
		Observacoes:                product.Observacoes,
		Descricaoembalagemdiscreta: product.Descricaoembalagemdiscreta,
		CreatedAt:                  product.CreatedAt,
		UpdatedAt:                  product.UpdatedAt,
	}

	if err := queries.CreateProduct(ctx, createProductParams); err != nil {
		fmt.Println("Erro ao criar um Rascunho:", err)
		return err
	}
	return nil
}

func (s *Store) GetProducts(nome, situacao string) ([]*types.Product, error) {
	queries := db.New(s.db)
	ctx := context.Background()
	params := db.GetProductsParams{
		Column1: nome,
		Column2: situacao,
	}

	dbProducts, err := queries.GetProducts(ctx, params)
	if err != nil {
		return nil, err
	}

	var products []*types.Product
	for _, dbProduct := range dbProducts {
		product := convertDBProductToProduct(dbProduct)
		products = append(products, product)
	}
	return products, nil
}

func (s *Store) GetProductNoMovements(nome, situacao string) ([]*types.ProductNoMovements, error) {
	queries := db.New(s.db)
	ctx := context.Background()
	params := db.GetProductNoMovementsParams{
		Column1: nome,
		Column2: situacao,
	}

	dbProducts, err := queries.GetProductNoMovements(ctx, params)
	if err != nil {
		return nil, err
	}

	var products []*types.ProductNoMovements
	for _, dbProduct := range dbProducts {
		product := convertGetProductNoMovementRowToProductNoMovementRow(dbProduct)
		products = append(products, product)
	}
	return products, nil
}

func (s *Store) GetProductEmptyStock(nome, situacao string) ([]*types.ProductEmptyStock, error) {
	queries := db.New(s.db)
	ctx := context.Background()
	params := db.GetProductEmptyStockParams{
		Column1: nome,
		Column2: situacao,
	}
	dbProducts, err := queries.GetProductEmptyStock(ctx, params)
	if err != nil {
		return nil, err
	}

	var products []*types.ProductEmptyStock
	for _, dbProduct := range dbProducts {
		product := convertGetProductEmptyStockRowToProductEmptyStockRow(dbProduct)
		products = append(products, product)
	}
	return products, nil
}

func (s *Store) UpdateProduct(product types.ProductPayload) error {
	queries := db.New(s.db)
	ctx := context.Background()

	now := time.Now()
	product.UpdatedAt = now

	updateProductParams := db.UpdateProductParams{
		ID:                         product.ID,
		Nome:                       product.Nome,
		Codigo:                     product.Codigo,
		Preco:                      product.Preco,
		ImagemUrl:                  product.ImagemUrl,
		Tipo:                       product.Tipo,
		Situacao:                   product.Situacao,
		Formato:                    product.Formato,
		DescricaoCurta:             product.DescricaoCurta,
		Datavalidade:               product.Datavalidade,
		Unidade:                    product.Unidade,
		Pesoliquido:                product.Pesoliquido,
		Pesobruto:                  product.Pesobruto,
		Volumes:                    product.Volumes,
		Itensporcaixa:              product.Itensporcaixa,
		Gtin:                       product.Gtin,
		Gtinembalagem:              product.Gtinembalagem,
		Tipoproducao:               product.Tipoproducao,
		Condicao:                   product.Condicao,
		Fretegratis:                product.Fretegratis,
		Marca:                      product.Marca,
		Descricaocomplementar:      product.Descricaocomplementar,
		Linkexterno:                product.Linkexterno,
		Observacoes:                product.Observacoes,
		Descricaoembalagemdiscreta: product.Descricaoembalagemdiscreta,
		UpdatedAt:                  product.UpdatedAt,
	}

	if err := queries.UpdateProduct(ctx, updateProductParams); err != nil {
		fmt.Println("Erro ao atualizar um Produto:", err)
		return err
	}
	return nil
}

func (s *Store) GetProductByID(productID int64) (*types.Product, error) {
	queries := db.New(s.db)
	ctx := context.Background()
	dbProduct, err := queries.GetProduct(ctx, productID)
	if err != nil {
		return nil, err
	}
	product := convertGetProductRowToProduct(dbProduct)

	return product, nil
}

func (s *Store) DeleteProduct(productID int64) error {
	queries := db.New(s.db)
	ctx := context.Background()
	err := queries.DeleteProduct(ctx, productID)
	if err != nil {
		return err
	}
	return nil
}

func convertDBProductToProduct(dbProduct db.GetProductsRow) *types.Product {
	product := &types.Product{
		ID:                         dbProduct.ID,
		Nome:                       dbProduct.Nome,
		Codigo:                     dbProduct.Codigo,
		Preco:                      dbProduct.Preco,
		ImagemUrl:                  dbProduct.ImagemUrl,
		Tipo:                       dbProduct.Tipo,
		Situacao:                   dbProduct.Situacao,
		Formato:                    dbProduct.Formato,
		DescricaoCurta:             dbProduct.DescricaoCurta,
		Datavalidade:               dbProduct.Datavalidade,
		Unidade:                    dbProduct.Unidade,
		Pesoliquido:                dbProduct.Pesoliquido,
		Pesobruto:                  dbProduct.Pesobruto,
		Volumes:                    dbProduct.Volumes,
		Itensporcaixa:              dbProduct.Itensporcaixa,
		Gtin:                       dbProduct.Gtin,
		Gtinembalagem:              dbProduct.Gtinembalagem,
		Tipoproducao:               dbProduct.Tipoproducao,
		Condicao:                   dbProduct.Condicao,
		Fretegratis:                dbProduct.Fretegratis,
		Marca:                      dbProduct.Marca,
		Descricaocomplementar:      dbProduct.Descricaocomplementar,
		Linkexterno:                dbProduct.Linkexterno,
		Observacoes:                dbProduct.Observacoes,
		Descricaoembalagemdiscreta: dbProduct.Descricaoembalagemdiscreta,
		CreatedAt:                  dbProduct.CreatedAt,
		UpdatedAt:                  dbProduct.UpdatedAt,
		SaldoFisicoTotal:           dbProduct.SaldoFisico,
		SaldoVirtualTotal:          dbProduct.SaldoVirtualTotal,
		SaldoFisico:                dbProduct.SaldoFisico,
		SaldoVirtual:               dbProduct.SaldoVirtual,
		PrecoCusto:                 dbProduct.PrecoCusto,
		PrecoCompra:                dbProduct.PrecoCompra,
		SupplierID:                 types.NullableInt64{NullInt64: dbProduct.SupplierID},
	}
	return product
}

func convertGetProductRowToProduct(dbProduct db.GetProductRow) *types.Product {
	product := &types.Product{
		ID:                         dbProduct.ID,
		Idprodutopai:               dbProduct.Idprodutopai,
		Nome:                       dbProduct.Nome,
		Codigo:                     dbProduct.Codigo,
		Preco:                      dbProduct.Preco,
		ImagemUrl:                  dbProduct.ImagemUrl,
		Tipo:                       dbProduct.Tipo,
		Situacao:                   dbProduct.Situacao,
		Formato:                    dbProduct.Formato,
		DescricaoCurta:             dbProduct.DescricaoCurta,
		Datavalidade:               dbProduct.Datavalidade,
		Unidade:                    dbProduct.Unidade,
		Pesoliquido:                dbProduct.Pesoliquido,
		Pesobruto:                  dbProduct.Pesobruto,
		Volumes:                    dbProduct.Volumes,
		Itensporcaixa:              dbProduct.Itensporcaixa,
		Gtin:                       dbProduct.Gtin,
		Gtinembalagem:              dbProduct.Gtinembalagem,
		Tipoproducao:               dbProduct.Tipoproducao,
		Condicao:                   dbProduct.Condicao,
		Fretegratis:                dbProduct.Fretegratis,
		Marca:                      dbProduct.Marca,
		Descricaocomplementar:      dbProduct.Descricaocomplementar,
		Linkexterno:                dbProduct.Linkexterno,
		Observacoes:                dbProduct.Observacoes,
		Descricaoembalagemdiscreta: dbProduct.Descricaoembalagemdiscreta,
		CreatedAt:                  dbProduct.CreatedAt,
		UpdatedAt:                  dbProduct.UpdatedAt,
		SaldoFisicoTotal:           dbProduct.SaldoFisicoTotal,
		SaldoVirtualTotal:          dbProduct.SaldoVirtualTotal,
		SaldoFisico:                dbProduct.SaldoFisico,
		SaldoVirtual:               dbProduct.SaldoVirtual,
		PrecoCusto:                 dbProduct.PrecoCusto,
		PrecoCompra:                dbProduct.PrecoCompra,
		SupplierID:                 types.NullableInt64{NullInt64: dbProduct.SupplierID},
	}
	return product
}

func convertGetProductNoMovementRowToProductNoMovementRow(dbProductNoMovement db.GetProductNoMovementsRow) *types.ProductNoMovements {
	ProductNoMovements := &types.ProductNoMovements{
		ID:                         dbProductNoMovement.ID,
		Nome:                       dbProductNoMovement.Nome,
		Codigo:                     dbProductNoMovement.Codigo,
		Preco:                      dbProductNoMovement.Preco,
		Tipo:                       dbProductNoMovement.Tipo,
		Situacao:                   dbProductNoMovement.Situacao,
		Formato:                    dbProductNoMovement.Formato,
		DescricaoCurta:             dbProductNoMovement.DescricaoCurta,
		ImagemUrl:                  dbProductNoMovement.ImagemUrl,
		Datavalidade:               dbProductNoMovement.Datavalidade,
		Unidade:                    dbProductNoMovement.Unidade,
		Pesoliquido:                dbProductNoMovement.Pesoliquido,
		Pesobruto:                  dbProductNoMovement.Pesobruto,
		Volumes:                    dbProductNoMovement.Volumes,
		Itensporcaixa:              dbProductNoMovement.Itensporcaixa,
		Gtin:                       dbProductNoMovement.Gtin,
		Gtinembalagem:              dbProductNoMovement.Gtinembalagem,
		Tipoproducao:               dbProductNoMovement.Tipoproducao,
		Condicao:                   dbProductNoMovement.Condicao,
		Fretegratis:                dbProductNoMovement.Fretegratis,
		Marca:                      dbProductNoMovement.Marca,
		Descricaocomplementar:      dbProductNoMovement.Descricaocomplementar,
		Linkexterno:                dbProductNoMovement.Linkexterno,
		Observacoes:                dbProductNoMovement.Observacoes,
		Descricaoembalagemdiscreta: dbProductNoMovement.Descricaoembalagemdiscreta,
		Numero:                     dbProductNoMovement.Numero,
		Numeroloja:                 dbProductNoMovement.Numeroloja,
		Data:                       dbProductNoMovement.Data,
		Datasaida:                  dbProductNoMovement.Datasaida,
		Dataprevista:               dbProductNoMovement.Dataprevista,
		Totalprodutos:              dbProductNoMovement.Totalprodutos,
		Totaldescontos:             dbProductNoMovement.Totaldescontos,
		PrecoCusto:                 dbProductNoMovement.PrecoCusto,
		PrecoCompra:                dbProductNoMovement.PrecoCompra,
	}
	return ProductNoMovements
}

func convertGetProductEmptyStockRowToProductEmptyStockRow(dbProductEmptyStock db.GetProductEmptyStockRow) *types.ProductEmptyStock {
	ProductEmptyStock := &types.ProductEmptyStock{
		ID:                         dbProductEmptyStock.ID,
		Idprodutopai:               dbProductEmptyStock.Idprodutopai,
		Nome:                       dbProductEmptyStock.Nome,
		Codigo:                     dbProductEmptyStock.Codigo,
		Preco:                      dbProductEmptyStock.Preco,
		Tipo:                       dbProductEmptyStock.Tipo,
		Situacao:                   dbProductEmptyStock.Situacao,
		Formato:                    dbProductEmptyStock.Formato,
		DescricaoCurta:             dbProductEmptyStock.DescricaoCurta,
		ImagemUrl:                  dbProductEmptyStock.ImagemUrl,
		Datavalidade:               dbProductEmptyStock.Datavalidade,
		Unidade:                    dbProductEmptyStock.Unidade,
		Pesoliquido:                dbProductEmptyStock.Pesoliquido,
		Pesobruto:                  dbProductEmptyStock.Pesobruto,
		Volumes:                    dbProductEmptyStock.Volumes,
		Itensporcaixa:              dbProductEmptyStock.Itensporcaixa,
		Gtin:                       dbProductEmptyStock.Gtin,
		Gtinembalagem:              dbProductEmptyStock.Gtinembalagem,
		Tipoproducao:               dbProductEmptyStock.Tipoproducao,
		Condicao:                   dbProductEmptyStock.Condicao,
		Fretegratis:                dbProductEmptyStock.Fretegratis,
		Marca:                      dbProductEmptyStock.Marca,
		Descricaocomplementar:      dbProductEmptyStock.Descricaocomplementar,
		Linkexterno:                dbProductEmptyStock.Linkexterno,
		Observacoes:                dbProductEmptyStock.Observacoes,
		Descricaoembalagemdiscreta: dbProductEmptyStock.Descricaoembalagemdiscreta,
		CreatedAt:                  dbProductEmptyStock.CreatedAt,
		UpdatedAt:                  dbProductEmptyStock.UpdatedAt,
		SaldoFisicoTotal:           dbProductEmptyStock.SaldoFisicoTotal,
		SaldoVirtualTotal:          dbProductEmptyStock.SaldoVirtualTotal,
		SaldoFisico:                dbProductEmptyStock.SaldoFisico,
		SaldoVirtual:               dbProductEmptyStock.SaldoVirtual,
		PrecoCusto:                 dbProductEmptyStock.PrecoCusto,
		PrecoCompra:                dbProductEmptyStock.PrecoCompra,
		SupplierID:                 dbProductEmptyStock.SupplierID,
	}
	return ProductEmptyStock
}
