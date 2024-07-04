package product

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/claudineyveloso/soldim.git/internal/bling"
	"github.com/claudineyveloso/soldim.git/internal/types"
	"github.com/gorilla/mux"
)

const limitePorPagina = 100

func RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/get_products", handleGetProduct).Methods(http.MethodGet)
	router.HandleFunc("/create_product", handleCreateProduct).Methods(http.MethodPost)
	router.HandleFunc("/update_product", handleUpdateProduct).Methods(http.MethodPut)
	router.HandleFunc("/delete_product", handleDeleteProduct).Methods(http.MethodDelete)
	router.HandleFunc("/get_product_id", handleGetProductId).Methods(http.MethodGet)
}

func handleGetProduct(w http.ResponseWriter, r *http.Request) {
	bearerToken := "bbe26de511a19fb57331313fd57049bc25fc4d84" // r.Header.Get("Authorization")

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	name := r.URL.Query().Get("name")
	criterioStr := r.URL.Query().Get("criterio")

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
		criterio = 0 // Valor padrão para criterio se não for fornecido ou inválido
	}

	fmt.Printf("Requesting page: %d with limit: %d and name: %s\n", page, limit, name)

	products, totalPages, err := bling.GetProductsFromBling(bearerToken, page, limit, name, criterio)
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
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Decodifica o JSON do corpo da requisição para a estrutura Product
	var newProduct types.Product
	if err := json.NewDecoder(r.Body).Decode(&newProduct); err != nil {
		http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}

	// Fecha o corpo da requisição após o processamento
	defer r.Body.Close()

	bearerToken := "bbe26de511a19fb57331313fd57049bc25fc4d84" // r.Header.Get("Authorization")
	// Chama a função para criar o produto no Bling
	err := bling.CreateProductInBling(bearerToken, newProduct)
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

	bearerToken := "bbe26de511a19fb57331313fd57049bc25fc4d84" // r.Header.Get("Authorization")
	// Chama a função para atualizar o produto no Bling
	err = bling.UpdateProductInBling(bearerToken, productID, updatedProduct)
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

	bearerToken := "bbe26de511a19fb57331313fd57049bc25fc4d84" // r.Header.Get("Authorization")
	// Chama a função para deletar o produto no Bling
	err = bling.DeleteProductInBling(bearerToken, productID)
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

	bearerToken := "bbe26de511a19fb57331313fd57049bc25fc4d84" // r.Header.Get("Authorization")
	// Chama a função para obter os detalhes do produto no Bling
	product, err := bling.GetProductIDInBling(bearerToken, productID)
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

	fmt.Fprintf(w, "Detalhes do produto obtidos com sucesso!")
}
