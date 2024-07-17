package store

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/claudineyveloso/soldim.git/internal/types"
	"github.com/claudineyveloso/soldim.git/internal/utils"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

type Handler struct {
	storeStore types.StoreStore
}

func NewHandler(storeStore types.StoreStore) *Handler {
	return &Handler{storeStore: storeStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/create_store", h.handleCreateStore).Methods(http.MethodPost)
	router.HandleFunc("/get_stores", h.handleGetStores).Methods(http.MethodGet)
	router.HandleFunc("/get_store/{storeID}", h.handleGetStore).Methods(http.MethodGet)
	router.HandleFunc("/get_store_description/{description}", h.handleGetStoreByDescription).Methods(http.MethodGet)
}

func (h *Handler) handleCreateStore(w http.ResponseWriter, r *http.Request) {
	var store types.Store
	if err := utils.ParseJSON(r, &store); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := utils.Validate.Struct(store); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Payload inválido: %v", errors))
		return
	}
	err := h.storeStore.CreateStore(store)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]interface{}{
		"data":    store,
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

func (h *Handler) handleGetStores(w http.ResponseWriter, r *http.Request) {
	store, err := h.storeStore.GetStores()
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao obter o Rascunho: %v", err), http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, http.StatusOK, store)
}

func (h *Handler) handleGetStore(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	storeIDStr, ok := vars["storeID"]
	if !ok {
		http.Error(w, "ID da Loja ausente!", http.StatusBadRequest)
		return
	}

	storeID, err := strconv.ParseInt(storeIDStr, 10, 64)
	if err != nil {
		http.Error(w, "ID da Loja inválido!", http.StatusBadRequest)
		return
	}

	store, err := h.storeStore.GetStoreByID(storeID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao obter a Situação: %v", err), http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, http.StatusOK, store)
}

func (h *Handler) handleGetStoreByDescription(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	description, ok := vars["description"]
	if !ok {
		http.Error(w, "Descrição da Situação ausente!", http.StatusBadRequest)
		return
	}

	situation, err := h.storeStore.GetStoreByDescription(description)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao obter a Situação: %v", err), http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, http.StatusOK, situation)
}
