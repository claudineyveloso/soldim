package bling

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Product representa um produto na API do Bling
type Product struct {
	ID             int64       `json:"id"`
	Nome           string      `json:"nome"`
	Codigo         string      `json:"codigo"`
	Preco          interface{} `json:"preco"`
	Tipo           string      `json:"tipo"`
	Situacao       string      `json:"situacao"`
	Formato        string      `json:"formato"`
	DescricaoCurta string      `json:"descricaoCurta"`
	ImagemURL      string      `json:"imagemURL"`
}

// GetProductsFromBling acessa a API de produtos do Bling usando o Bearer Token
func GetProductsFromBling(bearerToken string) ([]Product, error) {
	req, err := http.NewRequest("GET", "https://www.bling.com.br/Api/v3/produtos", nil)
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
		Data []Product `json:"data"`
	}
	if err := json.Unmarshal(bodyBytes, &responseData); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta: %v", err)
	}

	log.Printf("Número de produtos retornados: %d\n", len(responseData.Data))

	return responseData.Data, nil
}
