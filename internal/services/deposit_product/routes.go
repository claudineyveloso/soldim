package depositproduct

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
	depositproductStore types.DepositProductStore
}

func NewHandler(depositproductStore types.DepositProductStore) *Handler {
	return &Handler{depositproductStore: depositproductStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/create_deposit_product", h.handleCreateDepositProducts).Methods(http.MethodPost)
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
