package bling

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/claudineyveloso/soldim.git/internal/types"
)

const limitePorPagina = 100

func GetProductsFromBling(bearerToken string, page int, limit int, name string, criterio int, dataInclusaoInicial, dataInclusaoFinal, dataAlteracaoInicial, dataAlteracaoFinal string) ([]types.Product, int, error) {
	client := &http.Client{}

	// Construindo a URL com os parâmetros página, limite e nome (se fornecido)
	baseURL := "https://bling.com.br/Api/v3/produtos"
	params := url.Values{}
	params.Add("pagina", fmt.Sprintf("%d", page))
	params.Add("limite", fmt.Sprintf("%d", limit))
	if name != "" {
		params.Add("nome", name)
	}
	params.Add("criterio", fmt.Sprintf("%d", criterio))

	if dataInclusaoInicial != "" {
		params.Add("dataInclusaoInicial", dataInclusaoInicial)
	}
	if dataInclusaoFinal != "" {
		params.Add("dataInclusaoFinal", dataInclusaoFinal)
	}
	if dataAlteracaoInicial != "" {
		params.Add("dataAlteracaoInicial", dataAlteracaoInicial)
	}
	if dataAlteracaoFinal != "" {
		params.Add("dataAlteracaoFinal", dataAlteracaoFinal)
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

func CreateProductInBling(bearerToken string, product types.ProductBlingPayload) error {
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

func DeleteProductInBling(bearerToken string, productID int64) error {
	client := &http.Client{}

	// Construindo a URL para a exclusão de produtos
	url := fmt.Sprintf("https://bling.com.br/Api/v3/produtos/%d", productID)
	fmt.Printf("URL de requisição: %s\n", url) // Adicionando log para imprimir a URL
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("erro ao criar requisição: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+bearerToken)

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

func GetProductIDInBling(bearerToken string, productID int64) (*types.Product, error) {
	var product *types.Product
	rateLimiter := time.NewTicker(1 * time.Second)
	defer rateLimiter.Stop()
	url := fmt.Sprintf("https://bling.com.br/Api/v3/produtos/%d", productID)
	<-rateLimiter.C

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return product, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bearerToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return product, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return product, fmt.Errorf("failed to get product from Bling. Status: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return product, fmt.Errorf("error reading response body: %v", err)
	}

	var blingResponse struct {
		Data *types.Product `json:"data"`
	}
	err = json.Unmarshal(body, &blingResponse)
	if err != nil {
		return product, fmt.Errorf("error unmarshalling Bling product: %v", err)
	}

	return blingResponse.Data, nil
}
