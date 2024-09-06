package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/claudineyveloso/soldim.git/internal/types"
)

var limitePorPagina = 100

func ProcessProducts(products []types.Product, rateLimiter *time.Ticker) {
	for _, product := range products {
		<-rateLimiter.C
		productJSON, err := json.Marshal(product)
		if err != nil {
			fmt.Printf("Error marshalling product: %v\n", err)
			continue
		}

		req, err := http.NewRequest("POST", baseURL+"/create_product", bytes.NewBuffer(productJSON))
		if err != nil {
			fmt.Printf("Error creating request: %v\n", err)
			continue
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error sending request: %v\n", err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Failed to create product. Status: %v\n", resp.Status)
			continue
		}

		fmt.Printf("Product created successfully: %v\n", product)
	}
}
