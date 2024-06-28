package bling

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/claudineyveloso/soldim.git/internal/types"
)

const limitePorPagina = 100

func GetProductsFromBling(bearerToken string, page int, limit int, name string) ([]types.Product, int, error) {
	client := &http.Client{}

	// Construindo a URL com os parâmetros página, limite e nome (se fornecido)
	baseURL := "https://bling.com.br/Api/v3/produtos"
	params := url.Values{}
	params.Add("pagina", fmt.Sprintf("%d", page))
	params.Add("limite", fmt.Sprintf("%d", limit))
	if name != "" {
		params.Add("nome", name)
	}

	url := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	fmt.Printf("Enviando requisição para URL: %s\n", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("erro ao criar requisição: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+bearerToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("erro ao enviar requisição: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		return nil, 0, fmt.Errorf("falha na requisição: %s", bodyString)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("erro ao ler resposta: %v", err)
	}

	// Adicionando logs para verificar a resposta da API
	fmt.Printf("Resposta da API: %s\n", string(bodyBytes))

	var responseData struct {
		Data  []types.Product `json:"data"`
		Total int             `json:"total"`
		Limit int             `json:"limit"`
		Page  int             `json:"pagina"`
	}

	if err := json.Unmarshal(bodyBytes, &responseData); err != nil {
		return nil, 0, fmt.Errorf("erro ao decodificar resposta: %v", err)
	}

	// Adicionando logs para verificar a estrutura deserializada
	fmt.Printf("Dados deserializados: %+v\n", responseData)

	produtos := responseData.Data

	// Verificando a contagem de produtos
	fmt.Printf("Número de produtos deserializados: %d\n", len(produtos))

	// Calculando o número total de páginas com base no total de produtos e no limite por página
	totalPages := 0
	if responseData.Total > 0 && responseData.Limit > 0 {
		totalPages = (responseData.Total + responseData.Limit - 1) / responseData.Limit
	} else if len(produtos) == limit {
		// Se a API não retornar `total`, podemos inferir que há pelo menos mais uma página
		totalPages = page + 1
	}

	// Verificando o cálculo de totalPages
	fmt.Printf("Total de páginas calculado: %d\n", totalPages)

	return produtos, totalPages, nil
}

func CreateProductInBling(bearerToken string, product types.Product) error {
	client := &http.Client{}

	// Construindo a URL para a criação de produtos
	url := "https://bling.com.br/Api/v3/produtos"

	// Serializando o produto para JSON
	productData, err := json.Marshal(product)
	if err != nil {
		return fmt.Errorf("erro ao serializar produto: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(productData))
	if err != nil {
		return fmt.Errorf("erro ao criar requisição: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+bearerToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao enviar requisição: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		return fmt.Errorf("falha na requisição: %s", bodyString)
	}

	return nil
}

func UpdateProductInBling(bearerToken string, productID int64, product types.Product) error {
	client := &http.Client{}

	// Construindo a URL para a atualização de produtos
	url := fmt.Sprintf("https://bling.com.br/Api/v3/produtos/%d", productID)

	// Serializando o produto para JSON
	productData, err := json.Marshal(product)
	if err != nil {
		return fmt.Errorf("erro ao serializar produto: %v", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(productData))
	if err != nil {
		return fmt.Errorf("erro ao criar requisição: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+bearerToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao enviar requisição: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		return fmt.Errorf("falha na requisição: %s", bodyString)
	}

	return nil
}
