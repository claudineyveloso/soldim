package saleschannel

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/claudineyveloso/soldim.git/internal/types"
	"github.com/claudineyveloso/soldim.git/internal/utils"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

type Handler struct {
	salesChannelStore types.SalesChannelStore
}

func NewHandler(salesChannelStore types.SalesChannelStore) *Handler {
	return &Handler{salesChannelStore: salesChannelStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/create_sales_channel", h.handleCreateSalesChannel).Methods(http.MethodPost)
	router.HandleFunc("/get_sales_channel", h.handleGetSalesChannel).Methods(http.MethodGet)
}

func (h *Handler) handleGetSalesChannel(w http.ResponseWriter, r *http.Request) {
	saleschannel, err := h.salesChannelStore.GetSalesChannel()
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao obter os canais de venda : %v", err), http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, http.StatusOK, saleschannel)
}

func (h *Handler) handleCreateSalesChannel(w http.ResponseWriter, r *http.Request) {
	var saleschannel types.SalesChannel
	if err := utils.ParseJSON(r, &saleschannel); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := utils.Validate.Struct(saleschannel); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Payload inválido: %v", errors))
		return
	}
	err := h.salesChannelStore.CreateSalesChannel(saleschannel)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]interface{}{
		"data":    saleschannel,
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
