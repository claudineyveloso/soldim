package depositproduct

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
	depositproductStore types.DepositProductStore
}

func NewHandler(depositproductStore types.DepositProductStore) *Handler {
	return &Handler{depositproductStore: depositproductStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/create_deposit_product", h.handleCreateDepositProducts).Methods(http.MethodPost)
	router.HandleFunc("/update_deposit_product", h.handleUpdateDepositProducts).Methods(http.MethodPut)
	router.HandleFunc("/get_deposit_product_by_depositID/{depositID}", h.handleGetDepositProductByDepositID).Methods(http.MethodGet)
	router.HandleFunc("/get_deposit_product_by_productID/{productID}", h.handleGetDepositProductByProductID).Methods(http.MethodGet)
}

func (h *Handler) handleCreateDepositProducts(w http.ResponseWriter, r *http.Request) {
	var depositproduct types.DepositProduct
	if err := utils.ParseJSON(r, &depositproduct); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := utils.Validate.Struct(depositproduct); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Payload inválido: %v", errors))
		return
	}
	err := h.depositproductStore.CreateDepositProduct(depositproduct)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]interface{}{
		"data":    depositproduct,
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

func (h *Handler) handleUpdateDepositProducts(w http.ResponseWriter, r *http.Request) {
	var depositProduct types.DepositProduct
	if err := utils.ParseJSON(r, &depositProduct); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := utils.Validate.Struct(depositProduct); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("payload inválido: %v", errors))
		return
	}
	err := h.depositproductStore.UpdateDepositProduct(depositProduct)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]interface{}{
		"data":    depositProduct,
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

func (h *Handler) handleGetDepositProductByDepositID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	depositIDStr, ok := vars["depositID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("id do depósito ausente"))
		return
	}
	depositID, err := strconv.ParseInt(depositIDStr, 10, 64)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "room not found", http.StatusBadRequest)
			return
		}
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("id do depósito inválido"))
		return
	}

	product, err := h.depositproductStore.GetDepositProductByDepositID(depositID)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Error("Produto não encontrado dentro de handleGetProduct", slog.Int64("productID", depositID))
			utils.WriteError(w, http.StatusNotFound, fmt.Errorf("produto com ID %d não encontrado", depositID))
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, product)
}

func (h *Handler) handleGetDepositProductByProductID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productIDStr, ok := vars["productID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("id do depósito ausente"))
		return
	}
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "room not found", http.StatusBadRequest)
			return
		}
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("id do depósito inválido"))
		return
	}

	product, err := h.depositproductStore.GetDepositProductByProductID(productID)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Error("Produto não encontrado dentro de handleGetProduct", slog.Int64("productID", productID))
			utils.WriteError(w, http.StatusNotFound, fmt.Errorf("produto com ID %d não encontrado", productID))
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, product)
}
