package bling

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/claudineyveloso/soldim.git/internal/types"
)

func GetStockProductFromBling(bearerToken string, productID int64) (*types.StockResponse, error) {
	url := fmt.Sprintf("https://bling.com.br/Api/v3/estoques/saldos?idsProdutos%%5B%%5D=%d", productID)
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

	var stockResponse types.StockResponse
	err = json.Unmarshal(body, &stockResponse)
	if err != nil {
		return nil, fmt.Errorf("erro ao desserializar resposta: %v", err)
	}
	return &stockResponse, nil
}
