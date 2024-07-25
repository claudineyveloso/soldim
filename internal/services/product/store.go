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

	//	fmt.Println("Esse é o valor de dbProducts", dbProducts)
	fmt.Println("Esse é o valor de params", params)

	var products []*types.Product
	for _, dbProduct := range dbProducts {
		product := convertDBProductToProduct(dbProduct)
		products = append(products, product)
	}
	return products, nil
}

func (s *Store) GetProductNoMovements(nome, situacao string, limit, offset int32) ([]*types.ProductNoMovements, int64, error) {
	queries := db.New(s.db)
	ctx := context.Background()
	params := db.GetProductNoMovementsParams{
		Column1: nome,
		Column2: situacao,
		Limit:   limit,
		Offset:  offset,
	}

	dbProducts, err := queries.GetProductNoMovements(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	totalCountParams := db.GetTotalProductNoMovementsParams{
		Column1: nome,
		Column2: situacao,
	}
	totalCounts, err := queries.GetTotalProductNoMovements(ctx, totalCountParams)
	if err != nil {
		return nil, 0, err
	}

	var totalCount int64
	if len(totalCounts) > 0 {
		totalCount = totalCounts[0]
	}

	var products []*types.ProductNoMovements
	for _, dbProduct := range dbProducts {
		search := convertGetProductNoMovementRowToProductNoMovementRow(dbProduct)
		products = append(products, search)
	}
	return products, totalCount, nil
}

func (s *Store) GetProductEmptyStock(nome, situacao string) ([]*types.ProductEmptyStock, error) {
	queries := db.New(s.db)
	ctx := context.Background()
	dbProducts, err := queries.GetProductEmptyStock(ctx)
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
		SaldoFisicoTotal:           types.NullableInt{NullInt32: dbProduct.SaldoFisicoTotal},
		SaldoVirtualTotal:          types.NullableInt{NullInt32: dbProduct.SaldoVirtualTotal},
		SaldoFisico:                types.NullableInt{NullInt32: dbProduct.SaldoFisico},
		SaldoVirtual:               types.NullableInt{NullInt32: dbProduct.SaldoVirtual},
		PrecoCusto:                 types.NullableFloat{NullFloat64: dbProduct.PrecoCusto},
		PrecoCompra:                types.NullableFloat{NullFloat64: dbProduct.PrecoCompra},
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
		SaldoFisicoTotal:           types.NullableInt{NullInt32: dbProduct.SaldoFisicoTotal},
		SaldoVirtualTotal:          types.NullableInt{NullInt32: dbProduct.SaldoVirtualTotal},
		SaldoFisico:                types.NullableInt{NullInt32: dbProduct.SaldoFisico},
		SaldoVirtual:               types.NullableInt{NullInt32: dbProduct.SaldoVirtual},
		PrecoCusto:                 types.NullableFloat{NullFloat64: dbProduct.PrecoCusto},
		PrecoCompra:                types.NullableFloat{NullFloat64: dbProduct.PrecoCompra},
		SupplierID:                 types.NullableInt64{NullInt64: dbProduct.SupplierID},
	}
	return product
}

func convertGetProductNoMovementRowToProductNoMovementRow(dbProductNoMovement db.GetProductNoMovementsRow) *types.ProductNoMovements {
	ProductNoMovements := &types.ProductNoMovements{
		SalesOrderID:               dbProductNoMovement.SalesOrderID,
		ProductID:                  dbProductNoMovement.ProductID,
		Quantidade:                 dbProductNoMovement.Quantidade,
		ID:                         types.NullableInt64{NullInt64: dbProductNoMovement.ID},
		Nome:                       types.NullableString{NullString: dbProductNoMovement.Nome},
		Codigo:                     types.NullableString{NullString: dbProductNoMovement.Codigo},
		Preco:                      types.NullableFloat{NullFloat64: dbProductNoMovement.Preco},
		Tipo:                       types.NullableString{NullString: dbProductNoMovement.Tipo},
		Situacao:                   types.NullableString{NullString: dbProductNoMovement.Situacao},
		Formato:                    types.NullableString{NullString: dbProductNoMovement.Formato},
		DescricaoCurta:             types.NullableString{NullString: dbProductNoMovement.DescricaoCurta},
		ImagemUrl:                  types.NullableString{NullString: dbProductNoMovement.ImagemUrl},
		Datavalidade:               types.NullableTime{NullTime: dbProductNoMovement.Datavalidade},
		Unidade:                    types.NullableString{NullString: dbProductNoMovement.Unidade},
		Pesoliquido:                types.NullableFloat{NullFloat64: dbProductNoMovement.Pesoliquido},
		Pesobruto:                  types.NullableFloat{NullFloat64: dbProductNoMovement.Pesobruto},
		Volumes:                    types.NullableInt{NullInt32: dbProductNoMovement.Volumes},
		Itensporcaixa:              types.NullableInt{NullInt32: dbProductNoMovement.Itensporcaixa},
		Gtin:                       types.NullableString{NullString: dbProductNoMovement.Gtin},
		Gtinembalagem:              types.NullableString{NullString: dbProductNoMovement.Gtinembalagem},
		Tipoproducao:               types.NullableString{NullString: dbProductNoMovement.Tipoproducao},
		Condicao:                   types.NullableInt{NullInt32: dbProductNoMovement.Condicao},
		Fretegratis:                dbProductNoMovement.Fretegratis,
		Marca:                      types.NullableString{NullString: dbProductNoMovement.Marca},
		Descricaocomplementar:      types.NullableString{NullString: dbProductNoMovement.Descricaocomplementar},
		Linkexterno:                types.NullableString{NullString: dbProductNoMovement.Linkexterno},
		Observacoes:                types.NullableString{NullString: dbProductNoMovement.Observacoes},
		Descricaoembalagemdiscreta: types.NullableString{NullString: dbProductNoMovement.Descricaoembalagemdiscreta},
		Numero:                     types.NullableInt{NullInt32: dbProductNoMovement.Numero},
		Numeroloja:                 types.NullableString{NullString: dbProductNoMovement.Numeroloja},
		Data:                       types.NullableTime{NullTime: dbProductNoMovement.Data},
		Datasaida:                  types.NullableTime{NullTime: dbProductNoMovement.Datasaida},
		Dataprevista:               types.NullableTime{NullTime: dbProductNoMovement.Dataprevista},
		Totalprodutos:              types.NullableFloat{NullFloat64: dbProductNoMovement.Totalprodutos},
		Totaldescontos:             types.NullableFloat{NullFloat64: dbProductNoMovement.Totaldescontos},
		Descricao:                  types.NullableString{NullString: dbProductNoMovement.Descricao},
		Codigo_2:                   types.NullableInt64{NullInt64: dbProductNoMovement.Codigo_2},
		PrecoCusto:                 types.NullableFloat{NullFloat64: dbProductNoMovement.PrecoCusto},
		PrecoCompra:                types.NullableFloat{NullFloat64: dbProductNoMovement.PrecoCompra},
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
		SaldoFisicoTotal:           types.NullableInt{NullInt32: dbProductEmptyStock.SaldoFisicoTotal},
		SaldoVirtualTotal:          types.NullableInt{NullInt32: dbProductEmptyStock.SaldoVirtualTotal},
		SaldoFisico:                types.NullableInt{NullInt32: dbProductEmptyStock.SaldoFisico},
		SaldoVirtual:               types.NullableInt{NullInt32: dbProductEmptyStock.SaldoVirtual},
		PrecoCusto:                 types.NullableFloat{NullFloat64: dbProductEmptyStock.PrecoCusto},
		PrecoCompra:                types.NullableFloat{NullFloat64: dbProductEmptyStock.PrecoCompra},
		SupplierID:                 types.NullableInt64{NullInt64: dbProductEmptyStock.SupplierID},
	}
	return ProductEmptyStock
}
