package bling

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/claudineyveloso/soldim.git/internal/types"
)

// GetSalesChannelsFromBling  acessa a API de Canais de vendas do Bling usando o Bearer Token
func GetSalesChannelsFromBling(bearerToken string) ([]types.SalesChannel, error) {
	req, err := http.NewRequest("GET", "https://bling.com.br/Api/v3/canais-venda", nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+bearerToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao enviar requisição: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		return nil, fmt.Errorf("falha na requisição: %s", bodyString)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta: %v", err)
	}

	var responseData struct {
		Data []types.SalesChannel `json:"data"`
	}
	if err := json.Unmarshal(bodyBytes, &responseData); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta de etSalesChannelsFromBling: %v", err)
	}

	log.Printf("Número de canais de vendas retornados: %d\n", len(responseData.Data))

	return responseData.Data, nil
}
