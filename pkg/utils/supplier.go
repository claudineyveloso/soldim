package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/claudineyveloso/soldim.git/internal/bling"
	"github.com/claudineyveloso/soldim.git/internal/types"
	"golang.org/x/time/rate"
)

func ProcessSuppliers(products []types.Product, bearerToken string, rateLimiter *time.Ticker) {
	var wg sync.WaitGroup
	<-rateLimiter.C
	// Criar um rate limiter que permite 3 requisições por segundo
	limiter := rate.NewLimiter(rate.Every(time.Second/3), 1)

	for _, product := range products {
		wg.Add(1)
		go func(product types.Product) {
			defer wg.Done()

			// Aguardar até que uma requisição possa ser feita
			err := limiter.Wait(context.Background())
			if err != nil {
				fmt.Printf("Error waiting for rate limiter: %v\n", err)
				return
			}

			processSupplierForProduct(product, bearerToken, rateLimiter)
		}(product)
	}

	wg.Wait()
}

func processSupplierForProduct(product types.Product, bearerToken string, rateLimiter *time.Ticker) {
	<-rateLimiter.C
	supplierResponse, err := bling.GetSupplierProductFromBling(bearerToken, product.ID, rateLimiter)
	if err != nil {
		fmt.Printf("Error fetching supplier for product %d: %v\n", product.ID, err)
		return
	}

	for _, supplierData := range supplierResponse.Data {
		// Criar o supplier product
		supplierProduct := types.SupplierProduct{
			ID:          supplierData.ID,
			Descricao:   supplierData.Descricao,
			PrecoCusto:  supplierData.PrecoCusto,
			PrecoCompra: supplierData.PrecoCompra,
			Padrao:      supplierData.Padrao,
			SupplierID:  supplierData.Fornecedor.ID,
			ProductID:   supplierData.Produto.ID,
		}

		supplierProductJSON, err := json.Marshal(supplierProduct)
		if err != nil {
			fmt.Printf("Error marshalling supplier product for product %d: %v\n", product.ID, err)
			continue
		}

		fmt.Printf("Sending supplier product data for product %d: %s\n", product.ID, string(supplierProductJSON))
		<-rateLimiter.C
		supplierProductResp, err := http.Post(baseURL+"/create_supplier_product", "application/json", bytes.NewBuffer(supplierProductJSON))
		if err != nil {
			fmt.Printf("Error sending supplier product data for product %d: %v\n", product.ID, err)
			continue
		}
		defer supplierProductResp.Body.Close()

		supplierProductRespBody, _ := io.ReadAll(supplierProductResp.Body)
		fmt.Printf("Response from create_supplier_product for product %d: %s\n", product.ID, string(supplierProductRespBody))
	}
}
