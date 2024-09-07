package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/claudineyveloso/soldim.git/internal/bling"
	"github.com/claudineyveloso/soldim.git/internal/types"
)

func ProcessStocks(products []types.Product, bearerToken string) {
	var wg sync.WaitGroup
	rateLimiter := time.NewTicker(333 * time.Millisecond)
	defer rateLimiter.Stop()
	// Criar um rate limiter que permite 3 requisições por segundo
	// limiter := rate.NewLimiter(rate.Every(time.Second/3), 1)

	for _, product := range products {
		wg.Add(1)
		go func(product types.Product) {
			defer wg.Done()
			<-rateLimiter.C
			processStockForProduct(product, bearerToken)
		}(product)
	}

	wg.Wait()
}

func processStockForProduct(product types.Product, bearerToken string) {
	stockResponse, err := bling.GetStockProductFromBling(bearerToken, product.ID)
	if err != nil {
		fmt.Printf("Error fetching stock for product %d: %v\n", product.ID, err)
		LogErrorToFile(fmt.Sprintf("Error fetching stock for product %d: %v\n", product.ID, err))

		return
	}

	for _, stockData := range stockResponse.Data {
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

		stockResp, err := http.Post(baseURL+"/create_stock", "application/json", bytes.NewBuffer(stockJSON))
		if err != nil {
			fmt.Printf("Error sending stock data for product %d: %v\n", product.ID, err)
			continue
		}
		defer stockResp.Body.Close()

		// Criar deposit products
		for _, deposito := range stockData.Depositos {
			depositProduct := types.DepositProduct{
				ProductID:    stockData.Produto.ID,
				DepositID:    deposito.ID,
				SaldoFisico:  int32(deposito.SaldoFisico),
				SaldoVirtual: int32(deposito.SaldoVirtual),
			}

			depositProductJSON, err := json.Marshal(depositProduct)
			if err != nil {
				fmt.Printf("Error marshalling deposit product for product %d: %v\n", product.ID, err)
				continue
			}

			depositResp, err := http.Post(baseURL+"/create_deposit_product", "application/json", bytes.NewBuffer(depositProductJSON))
			if err != nil {
				fmt.Printf("Error sending deposit product data for product %d: %v\n", product.ID, err)
				continue
			}
			defer depositResp.Body.Close()
		}
	}
}
