package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/claudineyveloso/soldim.git/internal/types"
)

var limitePorPagina = 100

func ProcessProducts(products []types.Product) {
	for _, product := range products {
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

func CheckProductExists(productID int64) (bool, error) {
	url := fmt.Sprintf(baseURL+"/get_product/%d", productID)
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return false, nil // Produto não existe
	} else if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return true, nil // Produto existe
}
