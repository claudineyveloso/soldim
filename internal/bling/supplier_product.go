package bling

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/claudineyveloso/soldim.git/internal/types"
)

func GetSupplierProductFromBling(bearerToken string, productID int64, rateLimiter <-chan time.Time) (*types.SupplierResponse, error) {
	<-rateLimiter
	url := fmt.Sprintf("https://bling.com.br/Api/v3/produtos/fornecedores?idProduto=%d", productID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+bearerToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer requisição: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta: %v", err)
	}

	var supplierResponse types.SupplierResponse
	err = json.Unmarshal(body, &supplierResponse)
	if err != nil {
		return nil, fmt.Errorf("erro ao desserializar resposta: %v", err)
	}

	return &supplierResponse, nil
}
