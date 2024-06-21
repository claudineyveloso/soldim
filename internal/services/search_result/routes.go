package searchresult

import (
	"encoding/json"
	"fmt"
	"net/http"

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
	router.HandleFunc("/create_search_result", h.handleCreateSearchResult).Methods(http.MethodPost)
	router.HandleFunc("/delete_search_result/{searchID}", h.handleDeleteSearchResult).Methods(http.MethodDelete)
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
