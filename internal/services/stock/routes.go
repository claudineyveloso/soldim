package stock

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
	stockStore types.StockStore
}

func NewHandler(stockStore types.StockStore) *Handler {
	return &Handler{stockStore: stockStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/create_stock", h.handleCreateStocks).Methods(http.MethodPost)
}

func (h *Handler) handleCreateStocks(w http.ResponseWriter, r *http.Request) {
	var stock types.Stock
	if err := utils.ParseJSON(r, &stock); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := utils.Validate.Struct(stock); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Payload inválido: %v", errors))
		return
	}
	err := h.stockStore.CreateStock(stock)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]interface{}{
		"data":    stock,
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
