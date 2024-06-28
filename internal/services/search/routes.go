package search

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/claudineyveloso/soldim.git/internal/crawler"
	"github.com/claudineyveloso/soldim.git/internal/types"
	"github.com/claudineyveloso/soldim.git/internal/utils"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Handler struct {
	searchStore       types.SearchStore
	searchResultStore types.SearchResultStore
}

func NewHandler(searchStore types.SearchStore, searchResultStore types.SearchResultStore) *Handler {
	return &Handler{searchStore: searchStore, searchResultStore: searchResultStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/create_search", h.handleCreateSearch).Methods(http.MethodPost)
	router.HandleFunc("/get_searches", h.handleGetSearches).Methods(http.MethodGet)
	router.HandleFunc("/get_search/{searchID}", h.handleGetSearch).Methods(http.MethodGet)
	router.HandleFunc("/delete_search/{searchID}", h.handleDeleteSearch).Methods(http.MethodDelete)
}

func (h *Handler) handleCreateSearch(w http.ResponseWriter, r *http.Request) {
	var search types.SearchPayload
	if err := utils.ParseJSON(r, &search); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := utils.Validate.Struct(search); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Payload inválido: %v", errors))
		return
	}

	createSearch, err := h.searchStore.CreateSearch(search)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	products, err := crawler.CrawlGoogle(search.Description)
	if err != nil {
		log.Fatalf("Erro ao coletar produtos: %v", err)
	}

	for _, product := range products {
		searchResult := types.SearchResultPayload{
			ID:          uuid.New(),
			SearchID:    createSearch,
			ImageURL:    product.ImageURL,
			Description: product.Description,
			Source:      product.Source,
			Price:       product.Price,
			Link:        product.Link,
		}
		if err := h.searchResultStore.CreateSearchResult(searchResult); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}
	}

	response := map[string]interface{}{
		"data":    search,
		"message": "Registro criado com sucesso",
		"status":  http.StatusOK,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonResponse)
}

func (h *Handler) handleGetSearches(w http.ResponseWriter, r *http.Request) {
	// bucketID := auth.GetUserIDFromContext(r.Context())
	// fmt.Println("Valor de userIDffsadfsda", bucketID)
	searches, err := h.searchStore.GetSearches()
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao obter o Bucket: %v", err), http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, http.StatusOK, searches)
}

func (h *Handler) handleGetSearch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["searchID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("ID da Busca ausente!"))
		return
	}
	parsedSearchesID, err := uuid.Parse(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("ID do Busca inválido!"))
		return
	}

	search, err := h.searchStore.GetSearchByID(parsedSearchesID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, search)
}

func (h *Handler) handleDeleteSearch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["searchID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("ID do Busca ausente!"))
		return
	}
	parsedSearchesID, err := uuid.Parse(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("ID do Busca inválido!"))
		return
	}

	err = h.searchStore.DeleteSearch(parsedSearchesID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]interface{}{
		"message": "Registro apagado com sucesso",
		"status":  http.StatusOK,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonResponse)
}
