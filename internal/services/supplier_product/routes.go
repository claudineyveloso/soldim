package supplierproduct

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
	supplierproductStore types.SupplierProductStore
}

func NewHandler(supplierproductStore types.SupplierProductStore) *Handler {
	return &Handler{supplierproductStore: supplierproductStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/create_supplier_product", h.handleCreateSupplierProducts).Methods(http.MethodPost)
}

func (h *Handler) handleCreateSupplierProducts(w http.ResponseWriter, r *http.Request) {
	var supplierproduct types.SupplierProduct
	if err := utils.ParseJSON(r, &supplierproduct); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := utils.Validate.Struct(supplierproduct); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Payload inválido: %v", errors))
		return
	}
	err := h.supplierproductStore.CreateSupplierProduct(supplierproduct)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]interface{}{
		"data":    supplierproduct,
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
