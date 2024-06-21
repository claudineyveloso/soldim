package bling

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Product representa um produto na API do Bling
type SalesChannels struct {
	ID        int64  `json:"id"`
	Descricao string `json:"descricao"`
	Tipo      string `json:"tipo"`
	Situacao  int32  `json:"situacao"`
}

// GetProductsFromBling acessa a API de produtos do Bling usando o Bearer Token
func GetSalesChannelsFromBling(bearerToken string) ([]SalesChannels, error) {
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

	// Log da resposta JSON bruta
	// log.Printf("Resposta da API: %s\n", string(bodyBytes))

	var responseData struct {
		Data []SalesChannels `json:"data"`
	}
	if err := json.Unmarshal(bodyBytes, &responseData); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta: %v", err)
	}

	log.Printf("Número de canais de vendas retornados: %d\n", len(responseData.Data))

	return responseData.Data, nil
}
