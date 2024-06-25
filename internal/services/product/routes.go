package product

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/claudineyveloso/soldim.git/internal/bling"
	"github.com/claudineyveloso/soldim.git/internal/types"
	"github.com/gorilla/mux"
)

const limitePorPagina = 100

func RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/get_products", handleGetProduct).Methods(http.MethodGet)
}

func handleGetProduct(w http.ResponseWriter, r *http.Request) {
	bearerToken := "399e7a79230b1491672c67d4191b72a94268be24" // r.Header.Get("Authorization")

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	name := r.URL.Query().Get("name")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = limitePorPagina
	}

	fmt.Printf("Requesting page: %d with limit: %d and name: %s\n", page, limit, name)

	products, totalPages, err := bling.GetProductsFromBling(bearerToken, page, limit, name)
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

// func handleGetProductiii(w http.ResponseWriter, r *http.Request) {
// 	bearerToken := "5d1874f7ee1c00d9f76bf980572c105b11b164bf"
// 	pageStr := r.URL.Query().Get("page")
// 	page, err := strconv.Atoi(pageStr)
//
// 	if err != nil || page < 1 {
// 		page = 1
// 	}
//
// 	products, totalPages, err := bling.GetProductsFromBling(bearerToken, page)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Erro ao obter produtos: %v", err), http.StatusInternalServerError)
// 		return
// 	}
// 	response := struct {
// 		Products   []Product `json:"products"`
// 		TotalPages int       `json:"totalPages"`
// 	}{
// 		Products:   products,
// 		TotalPages: totalPages,
// 	}
//
// 	json.NewEncoder(w).Encode(response)
// }
