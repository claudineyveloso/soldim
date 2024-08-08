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
	salesOrderStore        types.SalesOrderStore
	productSalesOrderStore types.ProductSalesOrderStore
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
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("payload inválido: %v", errors))
		return
	}
	err := h.salesOrderStore.CreateSalesOrder(salesOrder)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		// return
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

	salesOrders, err := h.salesOrderStore.GetSalesOrders(int32(limit), int32(offset))
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao obter o Pedido de vendas: %v", err), http.StatusInternalServerError)
		return
	}

	response := struct {
		SalesOrders []*types.SalesOrder `json:"sales_orders"`
	}{
		SalesOrders: salesOrders,
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) handleGetSalesOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	salesOrderIDStr, ok := vars["salesOrderID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("id do produto ausente"))
		return
	}
	parsedSalesOrderID, err := strconv.ParseInt(salesOrderIDStr, 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("id do rascunho inválido"))
		return
	}

	salesOrder, err := h.salesOrderStore.GetSalesOrderByID(parsedSalesOrderID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, salesOrder)
}

func (h *Handler) handleGetSalesOrderNumber(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	salesOrderNumberStr, ok := vars["salesOrderNumber"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("número do produto ausente"))
		return
	}
	parsedSalesOrderNumber, err := strconv.ParseInt(salesOrderNumberStr, 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("ID do Produto inválido!"))
		return
	}

	// Converta parsedSalesOrderNumber de int64 para int32
	salesOrderNumber := int32(parsedSalesOrderNumber)

	salesOrder, err := h.salesOrderStore.GetSalesOrderByNumber(salesOrderNumber)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, salesOrder)
}
