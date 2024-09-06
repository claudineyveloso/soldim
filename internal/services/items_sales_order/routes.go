package itemssalesorder

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/claudineyveloso/soldim.git/internal/types"
	"github.com/claudineyveloso/soldim.git/pkg/utils"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

type Handler struct {
	itemsSalesOrderStore types.ItemsSalesOrderStore
}

func NewHandler(itemsSalesOrderStore types.ItemsSalesOrderStore) *Handler {
	return &Handler{itemsSalesOrderStore: itemsSalesOrderStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/create_items_sales_order", h.handleCreateItemsSalesOrder).Methods(http.MethodPost)
	router.HandleFunc("/get_items_sales_order_id/{itemsSalesOrderID}", h.handleGetItemsSalesOrder).Methods(http.MethodGet)
	router.HandleFunc("/get_items_sales_order_product_id/{productID}", h.handleGetItemsSalesOrderProductID).Methods(http.MethodGet)
	router.HandleFunc("/get_items_sales_order_sales_order_id/{salesOrderID}", h.handleGetItemsSalesOrderID).Methods(http.MethodGet)
}

func (h *Handler) handleCreateItemsSalesOrder(w http.ResponseWriter, r *http.Request) {
	var itemsSalesOrder types.ItemsSalesOrder
	if err := utils.ParseJSON(r, &itemsSalesOrder); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := utils.Validate.Struct(itemsSalesOrder); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("payload inválido: %v", errors))
		return
	}
	err := h.itemsSalesOrderStore.CreateItemsSalesOrder(itemsSalesOrder)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		// return
	}

	response := map[string]interface{}{
		"data":    itemsSalesOrder,
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

func (h *Handler) handleGetItemsSalesOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	itemsSalesOrderIDStr, ok := vars["itemsSalesOrderID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("id do item do pedido de venda ausente"))
		return
	}
	parsedItemsSalesOrderID, err := strconv.ParseInt(itemsSalesOrderIDStr, 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("id do rascunho inválido"))
		return
	}

	itemsSalesOrder, err := h.itemsSalesOrderStore.GetItemsSalesOrderByID(parsedItemsSalesOrderID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, itemsSalesOrder)
}

func (h *Handler) handleGetItemsSalesOrderProductID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	itemsSalesOrderIDStr, ok := vars["productID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("id do item do pedido de venda ausente"))
		return
	}
	parsedItemsSalesOrderID, err := strconv.ParseInt(itemsSalesOrderIDStr, 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("id do rascunho inválido"))
		return
	}

	itemsSalesOrder, err := h.itemsSalesOrderStore.GetItemsSalesOrderByProductID(parsedItemsSalesOrderID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, itemsSalesOrder)
}

func (h *Handler) handleGetItemsSalesOrderID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	itemsSalesOrderIDStr, ok := vars["salesOrderID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("id do item do pedido de venda ausente"))
		return
	}
	parsedItemsSalesOrderID, err := strconv.ParseInt(itemsSalesOrderIDStr, 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("id do rascunho inválido"))
		return
	}

	itemsSalesOrder, err := h.itemsSalesOrderStore.GetItemsSalesOrderBySalesOrderID(parsedItemsSalesOrderID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, itemsSalesOrder)
}
