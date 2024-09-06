package stock

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/claudineyveloso/soldim.git/internal/types"
	"github.com/claudineyveloso/soldim.git/pkg/utils"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
)

type Handler struct {
	stockStore types.StockStore
}

func NewHandler(stockStore types.StockStore) *Handler {
	return &Handler{stockStore: stockStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/get_stock/{productID}", h.handleGetStockProductID).Methods(http.MethodGet)
	router.HandleFunc("/create_stock", h.handleCreateStocks).Methods(http.MethodPost)
	router.HandleFunc("/update_stock", h.handleUpdateStock).Methods(http.MethodPut)
}

func (h *Handler) handleGetStockProductID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productIDStr, ok := vars["productID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("id do produto ausente"))
		return
	}
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "room not found", http.StatusBadRequest)
			return
		}
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("id do rascunho inválido"))
		return
	}

	stock, err := h.stockStore.GetStockByProductID(productID)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Error("Produto não encontrado dentro do Estoque", slog.Int64("productID", productID))
			utils.WriteError(w, http.StatusNotFound, fmt.Errorf("produto com ID %d não encontrado", productID))
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, stock)
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

func (h *Handler) handleUpdateStock(w http.ResponseWriter, r *http.Request) {
	var stock types.Stock
	if err := utils.ParseJSON(r, &stock); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := utils.Validate.Struct(stock); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("payload inválido: %v", errors))
		return
	}
	err := h.stockStore.UpdateStock(stock)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]interface{}{
		"data":    stock,
		"message": "Registro alterado com sucesso",
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
