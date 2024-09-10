package productbling

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/claudineyveloso/soldim.git/internal/bling"
	"github.com/claudineyveloso/soldim.git/internal/types"
	"github.com/claudineyveloso/soldim.git/pkg/utils"
	"github.com/gorilla/mux"
)

var (
	limitePorPagina = 100

	baseURL = utils.GetBaseURL()
)

func RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/import_products", handleImportBlingProductsToSoldim).Methods(http.MethodGet)
	router.HandleFunc("/get_products_bling", handleGetProduct).Methods(http.MethodGet)
	router.HandleFunc("/create_product_bling", handleCreateProduct).Methods(http.MethodPost)
	router.HandleFunc("/update_product_bling", handleUpdateProduct).Methods(http.MethodPut)
	router.HandleFunc("/delete_product_bling", handleDeleteProduct).Methods(http.MethodDelete)
	router.HandleFunc("/get_product_id_bling", handleGetProductId).Methods(http.MethodGet)
}

func handleImportBlingProductsToSoldim(w http.ResponseWriter, r *http.Request) {
	page := 1
	limit := 100 // Processa 100 produtos por vez

	// Coleta os parâmetros da query string
	params := r.URL.Query()

	pagina := params.Get("pagina")
	if pagina != "" {
		page, _ = strconv.Atoi(pagina) // Converte para int
	}

	limite := params.Get("limite")
	if limite != "" {
		limit, _ = strconv.Atoi(limite) // Converte para int
	}

	nome := params.Get("nome")
	criterio := params.Get("criterio")
	criterioInt, err := strconv.Atoi(criterio)
	if err != nil {
		// Tratar erro caso a conversão falhe
		fmt.Println("Erro ao converter criterio para int:", err)
		criterioInt = 0 // Valor padrão caso a conversão falhe
	}
	dataInclusaoInicial := params.Get("dataInclusaoInicial")
	dataInclusaoFinal := params.Get("dataInclusaoFinal")
	dataAlteracaoInicial := params.Get("dataAlteracaoInicial")
	dataAlteracaoFinal := params.Get("dataAlteracaoFinal")

	token, err := utils.FetchAccessToken()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching access token: %v", err), http.StatusInternalServerError)
		return
	}

	rateLimiter := time.NewTicker(333 * time.Millisecond) // 3 requisições por segundo
	defer rateLimiter.Stop()

	for {
		// Faz a requisição de uma página de produtos
		products, totalPages, err := bling.GetProductsFromBling(
			token, page, limit, nome, criterioInt, dataInclusaoInicial, dataInclusaoFinal, dataAlteracaoInicial, dataAlteracaoFinal)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(products) == 0 {
			break // Nenhum produto restante
		}

		// Processa os produtos em paralelo
		processProductsConcurrently(products, rateLimiter, token)

		if page >= totalPages {
			break // Última página alcançada
		}

		page++ // Próxima página
	}

	responseMessage := map[string]interface{}{
		"message": "Registros importados e atualizados com sucesso",
		"status":  http.StatusOK,
	}
	jsonResponse, err := json.Marshal(responseMessage)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error marshalling response: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonResponse)
}

func processProductsConcurrently(products []types.Product, rateLimiter *time.Ticker, token string) {
	var wg sync.WaitGroup
	sem := make(chan struct{}, 3) // Limitar a 3 requisições simultâneas

	for _, product := range products {
		wg.Add(1)
		go func(product types.Product) {
			defer wg.Done()

			sem <- struct{}{} // Limitar as goroutines em andamento
			processProduct(product, rateLimiter, token)
			<-sem
		}(product)
	}

	wg.Wait() // Aguarda o processamento de todos os produtos
}

func processProduct(product types.Product, rateLimiter *time.Ticker, token string) {
	// Atualizar preço
	<-rateLimiter.C
	utils.ProcessProducts([]types.Product{product})

	// Atualizar estoque
	<-rateLimiter.C
	utils.ProcessStocks([]types.Product{product}, token)

	// Atualizar fornecedores
	<-rateLimiter.C
	utils.ProcessSuppliers([]types.Product{product}, token)
}

func handleGetProduct(w http.ResponseWriter, r *http.Request) {
	token, err := utils.FetchAccessToken()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching access token: %v", err), http.StatusInternalServerError)
		return
	}

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	name := r.URL.Query().Get("name")
	criterioStr := r.URL.Query().Get("criterio")

	dataInclusaoInicial := r.URL.Query().Get("dataInclusaoInicial")
	dataInclusaoFinal := r.URL.Query().Get("dataInclusaoFinal")
	dataAlteracaoInicial := r.URL.Query().Get("dataAlteracaoInicial")
	dataAlteracaoFinal := r.URL.Query().Get("dataAlteracaoFinal")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = limitePorPagina
	}

	criterio, err := strconv.Atoi(criterioStr)
	if err != nil {
		criterio = 0
	}

	fmt.Printf("Requesting page: %d with limit: %d and name: %s\n", page, limit, name)
	products, totalPages, err := bling.GetProductsFromBling(token, page, limit, name, criterio, dataInclusaoInicial, dataInclusaoFinal, dataAlteracaoInicial, dataAlteracaoFinal)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Products   []types.Product `json:"products"`
		TotalPages int             `json:"totalPages"`
		Page       int             `json:"page"`
		Limit      int             `json:"limit"`
	}{
		Products:   products,
		TotalPages: totalPages,
		Page:       page,
		Limit:      limit,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Printf("Retornando %d produtos e %d páginas\n", len(products), totalPages)
}

func handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	token, err := utils.FetchAccessToken()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching access token: %v", err), http.StatusInternalServerError)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Loga o corpo da requisição para inspeção
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Erro ao ler corpo da requisição", http.StatusBadRequest)
		return
	}
	fmt.Printf("Corpo da requisição: %s\n", string(body))

	// Decodifica o JSON do corpo da requisição para a estrutura Product
	var newProduct types.ProductBlingPayload
	if err := json.Unmarshal(body, &newProduct); err != nil {
		http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}

	newProduct.DataValidade = time.Now().Format("2006-01-02")

	// Fecha o corpo da requisição após o processamento
	defer r.Body.Close()

	fmt.Printf("Novo produto: %+v\n", newProduct)

	// Chama a função para criar o produto no Bling
	err = bling.CreateProductInBling(token, newProduct)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao criar produto: %v", err), http.StatusInternalServerError)
		log.Fatalf("Erro ao criar produto: %v", err)
		return
	}

	// Responde com sucesso
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Produto criado com sucesso!")
}

func handleUpdateProduct(w http.ResponseWriter, r *http.Request) {
	token, err := utils.FetchAccessToken()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching access token: %v", err), http.StatusInternalServerError)
		return
	}

	if r.Method != http.MethodPut {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Extrai o productID dos parâmetros da URL
	productIDStr := r.URL.Query().Get("productID")
	if productIDStr == "" {
		http.Error(w, "productID é necessário", http.StatusBadRequest)
		return
	}
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		http.Error(w, "productID inválido", http.StatusBadRequest)
		return
	}

	// Decodifica o JSON do corpo da requisição para a estrutura Product
	var updatedProduct types.Product
	if err := json.NewDecoder(r.Body).Decode(&updatedProduct); err != nil {
		http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}

	// Fecha o corpo da requisição após o processamento
	defer r.Body.Close()

	// Chama a função para atualizar o produto no Bling
	err = bling.UpdateProductInBling(token, productID, updatedProduct)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao atualizar produto: %v", err), http.StatusInternalServerError)
		log.Fatalf("Erro ao atualizar produto: %v", err)
		return
	}

	// Responde com sucesso
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Produto atualizado com sucesso!")
}

func handleDeleteProduct(w http.ResponseWriter, r *http.Request) {
	token, err := utils.FetchAccessToken()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching access token: %v", err), http.StatusInternalServerError)
		return
	}

	if r.Method != http.MethodDelete {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Extrai o productID dos parâmetros da URL
	productIDStr := r.URL.Query().Get("productID")
	if productIDStr == "" {
		http.Error(w, "productID é necessário", http.StatusBadRequest)
		return
	}
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		http.Error(w, "productID inválido", http.StatusBadRequest)
		return
	}

	// Chama a função para deletar o produto no Bling
	err = bling.DeleteProductInBling(token, productID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao deletar produto: %v", err), http.StatusInternalServerError)
		log.Fatalf("Erro ao deletar produto: %v", err)
		return
	}

	// Responde com sucesso
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Produto deletado com sucesso!")
}

func handleGetProductId(w http.ResponseWriter, r *http.Request) {
	token, err := utils.FetchAccessToken()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching access token: %v", err), http.StatusInternalServerError)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Extrai o productID dos parâmetros da URL
	productIDStr := r.URL.Query().Get("productID")
	if productIDStr == "" {
		http.Error(w, "productID é necessário", http.StatusBadRequest)
		return
	}
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		http.Error(w, "productID inválido", http.StatusBadRequest)
		return
	}

	// Chama a função para obter os detalhes do produto no Bling
	product, err := bling.GetProductIDInBling(token, productID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao obter detalhes do produto: %v", err), http.StatusInternalServerError)
		log.Fatalf("Erro ao obter detalhes do produto: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func fetchWithRetries(url string, retries int) (*http.Response, error) {
	for i := 0; i < retries; i++ {
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode != http.StatusTooManyRequests {
			return resp, err
		}

		// Close the response body to prevent resource leakage
		if resp != nil {
			resp.Body.Close()
		}

		// Exponential backoff
		time.Sleep(time.Duration((1<<i)*100) * time.Millisecond)
	}

	return nil, fmt.Errorf("failed to get product from Bling after %d retries", retries)
}
