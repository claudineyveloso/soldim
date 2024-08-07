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
		return nil, 0, fmt.Errorf("erro ao decodificar resposta em GetSalesOrdersFromBling: %v", err)
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

func GetSalesOrdersIDInBling(bearerToken string, salesOrderID int64) ([]types.SalesOrder, error) {
	url := fmt.Sprintf("https://api.bling.com.br/Api/v2/pedido/%d/json", salesOrderID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bearerToken))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get sales order: %s", resp.Status)
	}

	var result struct {
		Retorno struct {
			Pedidos []struct {
				Pedido types.SalesOrder `json:"pedido"`
			} `json:"pedidos"`
		} `json:"retorno"`
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	var salesOrders []types.SalesOrder
	for _, p := range result.Retorno.Pedidos {
		salesOrders = append(salesOrders, p.Pedido)
	}

	return salesOrders, nil
}
