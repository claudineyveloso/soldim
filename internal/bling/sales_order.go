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

func GetSalesOrdersIDInBling(bearerToken string, SalesOrderID int64) (types.SalesOrder, error) {
	var salesOrder types.SalesOrder

	url := fmt.Sprintf("https://bling.com.br/Api/v3/pedidos/vendas/%d", SalesOrderID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return salesOrder, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bearerToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return salesOrder, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return salesOrder, fmt.Errorf("failed to get salesOrder from Bling. Status: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return salesOrder, fmt.Errorf("error reading response body: %v", err)
	}

	var blingResponse struct {
		Data types.SalesOrder `json:"data"`
	}
	err = json.Unmarshal(body, &blingResponse)
	if err != nil {
		return salesOrder, fmt.Errorf("error unmarshalling Bling salesOrder: %v", err)
	}

	return blingResponse.Data, nil
}
