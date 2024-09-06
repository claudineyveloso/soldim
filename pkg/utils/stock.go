package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/claudineyveloso/soldim.git/internal/bling"
	"github.com/claudineyveloso/soldim.git/internal/types"
	"golang.org/x/time/rate"
)

func ProcessStocks(products []types.Product, bearerToken string, rateLimiter *time.Ticker) {
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

			processStockForProduct(product, bearerToken, rateLimiter)
		}(product)
	}

	wg.Wait()
}

func processStockForProduct(product types.Product, bearerToken string, rateLimiter *time.Ticker) {
	<-rateLimiter.C
	stockResponse, err := bling.GetStockProductFromBling(bearerToken, product.ID)
	if err != nil {
		fmt.Printf("Error fetching stock for product %d: %v\n", product.ID, err)
		LogErrorToFile(fmt.Sprintf("Error fetching stock for product %d: %v\n", product.ID, err))

		return
	}

	// fmt.Printf("Parsed stock response for product %d: %+v\n", product.ID, stockResponse)

	for _, stockData := range stockResponse.Data {
		// fmt.Printf("Processing stock data: %+v\n", stockData)

		// Criar o stock
		stock := types.Stock{
			ProductID:         stockData.Produto.ID,
			SaldoFisicoTotal:  int32(stockData.SaldoFisicoTotal),
			SaldoVirtualTotal: int32(stockData.SaldoVirtualTotal),
		}

		stockJSON, err := json.Marshal(stock)
		if err != nil {
			fmt.Printf("Error marshalling stock for product %d: %v\n", product.ID, err)
			continue
		}

		// fmt.Printf("Sending stock data for product %d: %s\n", product.ID, string(stockJSON))

		stockResp, err := http.Post(baseURL+"/create_stock", "application/json", bytes.NewBuffer(stockJSON))
		if err != nil {
			fmt.Printf("Error sending stock data for product %d: %v\n", product.ID, err)
			continue
		}
		defer stockResp.Body.Close()

		// Criar deposit products
		for _, deposito := range stockData.Depositos {
			// fmt.Printf("Processing deposit data: %+v\n", deposito)

			depositProduct := types.DepositProduct{
				ProductID:    stockData.Produto.ID,
				DepositID:    deposito.ID,
				SaldoFisico:  int32(deposito.SaldoFisico),
				SaldoVirtual: int32(deposito.SaldoVirtual),
			}

			// fmt.Printf("Values assigned for deposit product for product %d: %+v\n", product.ID, depositProduct)

			depositProductJSON, err := json.Marshal(depositProduct)
			if err != nil {
				fmt.Printf("Error marshalling deposit product for product %d: %v\n", product.ID, err)
				continue
			}

			// fmt.Printf("Sending deposit product data for product %d: %s\n", product.ID, string(depositProductJSON))

			depositResp, err := http.Post(baseURL+"/create_deposit_product", "application/json", bytes.NewBuffer(depositProductJSON))
			if err != nil {
				fmt.Printf("Error sending deposit product data for product %d: %v\n", product.ID, err)
				continue
			}
			defer depositResp.Body.Close()

			// depositRespBody, _ := io.ReadAll(depositResp.Body)
			// fmt.Printf("Response from create_deposit_product for product %d: %s\n", product.ID, string(depositRespBody))
		}
	}
}
