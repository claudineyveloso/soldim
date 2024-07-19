package bling

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/claudineyveloso/soldim.git/internal/types"
)

func GetSalesOrdersFromBling(bearerToken string, page int, limit int) ([]types.SalesOrder, int, error) {
	client := &http.Client{}

	baseURL := "https://bling.com.br/Api/v3/pedidos/vendas"
	params := url.Values{}
	params.Add("pagina", strconv.Itoa(page))
	params.Add("limite", strconv.Itoa(limit))

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

	fmt.Printf("Resposta da API: %s\n", string(bodyBytes))

	var responseData struct {
		Data  []types.SalesOrder `json:"data"`
		Total int                `json:"total"`
		Limit int                `json:"limit"`
		Page  int                `json:"pagina"`
	}

	if err := json.Unmarshal(bodyBytes, &responseData); err != nil {
		return nil, 0, fmt.Errorf("erro ao decodificar resposta: %v", err)
	}

	fmt.Printf("Número de vendas deserializados: %d\n", len(responseData.Data))

	totalPages := 0
	if responseData.Total > 0 && responseData.Limit > 0 {
		totalPages = (responseData.Total + responseData.Limit - 1) / responseData.Limit
	} else if len(responseData.Data) == limit {
		totalPages = page + 1
	}

	fmt.Printf("Total de páginas calculado: %d\n", totalPages)

	return responseData.Data, totalPages, nil
}

func GetSalesOrdersFromBling_old(bearerToken string, page int, limit int) ([]types.SalesOrder, int, error) {
	client := &http.Client{}

	// Construindo a URL com os parâmetros página, limite e nome (se fornecido)
	baseURL := "https://bling.com.br/Api/v3/pedidos/vendas"

	params := url.Values{}
	params.Add("pagina", fmt.Sprintf("%d", page))
	params.Add("limite", fmt.Sprintf("%d", limit))

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
		Data  []types.SalesOrder `json:"data"`
		Total int                `json:"total"`
		Limit int                `json:"limit"`
		Page  int                `json:"pagina"`
	}

	if err := json.Unmarshal(bodyBytes, &responseData); err != nil {
		return nil, 0, fmt.Errorf("erro ao decodificar resposta: %v", err)
	}

	// Adicionando logs para verificar a estrutura deserializada
	// fmt.Printf("Dados deserializados: %+v\n", responseData)

	sales := responseData.Data

	// Verificando a contagem de vendas
	fmt.Printf("Número de vendas deserializados: %d\n", len(sales))

	// Calculando o número total de páginas com base no total de vendas e no limite por página
	totalPages := 0
	if responseData.Total > 0 && responseData.Limit > 0 {
		totalPages = (responseData.Total + responseData.Limit - 1) / responseData.Limit
	} else if len(sales) == limit {
		// Se a API não retornar `total`, podemos inferir que há pelo menos mais uma página
		totalPages = page + 1
	}

	// Verificando o cálculo de totalPages
	fmt.Printf("Total de páginas calculado: %d\n", totalPages)

	return sales, totalPages, nil
}
