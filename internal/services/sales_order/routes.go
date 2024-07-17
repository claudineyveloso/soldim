package salesorder

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
	salesOrderStore types.SalesOrderStore
}

func NewHandler(salesOrderStore types.SalesOrderStore) *Handler {
	return &Handler{salesOrderStore: salesOrderStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/create_sales_order", h.handleCreateSalesOrder).Methods(http.MethodPost)
	router.HandleFunc("/get_sales_orders", h.handleGetSalesOrders).Methods(http.MethodGet)
	router.HandleFunc("/get_sales_order/{salesOrderID}", h.handleGetSalesOrder).Methods(http.MethodGet)
}

func (h *Handler) handleCreateSalesOrder(w http.ResponseWriter, r *http.Request) {
	var salesOrder types.SalesOrder
	if err := utils.ParseJSON(r, &salesOrder); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := utils.Validate.Struct(salesOrder); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Payload inválido: %v", errors))
		return
	}
	err := h.salesOrderStore.CreateSalesOrder(salesOrder)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]interface{}{
		"data":    salesOrder,
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

func (h *Handler) handleGetSalesOrders(w http.ResponseWriter, r *http.Request) {
	// bucketID := auth.GetUserIDFromContext(r.Context())
	// fmt.Println("Valor de userIDffsadfsda", bucketID)
	salesOrder, err := h.salesOrderStore.GetSalesOrders()
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao obter o Bucket: %v", err), http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, http.StatusOK, salesOrder)
}

func (h *Handler) handleGetSalesOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productIDStr, ok := vars["salesOrderID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("ID do Produto ausente!"))
		return
	}
	parsedSalesOrderID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("ID do Rascunho inválido!"))
		return
	}

	product, err := h.salesOrderStore.GetSalesOrderByID(parsedSalesOrderID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, product)
}
