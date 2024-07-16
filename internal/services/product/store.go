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

func (s *Store) GetProducts() ([]*types.Product, error) {
	queries := db.New(s.db)
	ctx := context.Background()

	dbProducts, err := queries.GetProducts(ctx)
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
