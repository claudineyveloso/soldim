package searchresult

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/claudineyveloso/soldim.git/internal/types"
	"github.com/claudineyveloso/soldim.git/internal/utils"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Handler struct {
	searchresultStore types.SearchResultStore
}

func NewHandler(searchresultStore types.SearchResultStore) *Handler {
	return &Handler{searchresultStore: searchresultStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/get_searches_result", h.handleGetSearchesResult).Methods(http.MethodGet)
	router.HandleFunc("/get_search_result/{searchResultID}", h.handleGetSearchResult).Methods(http.MethodGet)
	router.HandleFunc("/get_search_result_sources/{searchID}", h.handleGetSearchResultSources).Methods(http.MethodGet)
	router.HandleFunc("/create_search_result", h.handleCreateSearchResult).Methods(http.MethodPost)
	router.HandleFunc("/delete_search_result/{searchID}", h.handleDeleteSearchResult).Methods(http.MethodDelete)
}

func (h *Handler) handleGetSearchesResult(w http.ResponseWriter, r *http.Request) {
	source := r.URL.Query().Get("source")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10 // Default limit
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0 // Default offset
	}

	searchesResults, totalCount, err := h.searchresultStore.GetSearchesResult(source, int32(limit), int32(offset))

	fmt.Println("Esse é o valor de totalCount: ", totalCount)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao obter o Resultado da Busca: %v", err), http.StatusInternalServerError)
		return
	}

	response := struct {
		SearchResults []*types.SearchResult `json:"search_results"`
		TotalCount    int64                 `json:"total_count"`
	}{
		SearchResults: searchesResults,
		TotalCount:    totalCount,
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) handleGetSearchResult(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["searchResultID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("ID do Resultado da Busca ausente!"))
		return
	}
	parsedSearchResultID, err := uuid.Parse(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("ID do Resultado da Busca inválido!"))
		return
	}

	searchResult, err := h.searchresultStore.GetSearchResultByID(parsedSearchResultID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, searchResult)
}

func (h *Handler) handleGetSearchResultSources(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["searchID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("ID da Busca ausente!"))
		return
	}
	parsedSearchID, err := uuid.Parse(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("ID da Busca inválido!"))
		return
	}
	searchResultSources, err := h.searchresultStore.GetSearchResultSources(parsedSearchID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, searchResultSources)
}

func (h *Handler) handleCreateSearchResult(w http.ResponseWriter, r *http.Request) {
	var search_result types.SearchResultPayload
	if err := utils.ParseJSON(r, &search_result); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := utils.Validate.Struct(search_result); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Payload inválido: %v", errors))
		return
	}
	err := h.searchresultStore.CreateSearchResult(search_result)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]interface{}{
		"data":    search_result,
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

func (h *Handler) handleDeleteSearchResult(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["searchID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("ID do Resultado da Busca ausente!"))
		return
	}
	parsedSearchesID, err := uuid.Parse(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("ID do Resultado da Busca inválido!"))
		return
	}

	err = h.searchresultStore.DeleteSearchResult(parsedSearchesID)
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
