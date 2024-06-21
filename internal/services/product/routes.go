package product

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/claudineyveloso/soldim.git/internal/bling"
	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/get_products", handleGetProduct).Methods(http.MethodGet)
}

func handleGetProduct(w http.ResponseWriter, r *http.Request) {
	bearerToken := "fe30f48f033b53f4fe262eb521e9554303e0235c"
	produtos, err := bling.GetProductsFromBling(bearerToken)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao obter produtos: %v", err), http.StatusInternalServerError)
		return
	}

	if len(produtos) == 0 {
		http.Error(w, "Nenhum produto encontrado.", http.StatusNotFound)
		return
	}

	// Definir o tipo de conteúdo da resposta como JSON
	w.Header().Set("Content-Type", "application/json")

	// Converter os produtos para JSON e enviar a resposta
	err = json.NewEncoder(w).Encode(produtos)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao converter produtos para JSON: %v", err), http.StatusInternalServerError)
	}
}
