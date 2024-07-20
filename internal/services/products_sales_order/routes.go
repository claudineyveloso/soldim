package productssalesorder

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
	productSalesOrderStore types.ProductSalesOrderStore
}

func NewHandler(productSalesOrderStore types.ProductSalesOrderStore) *Handler {
	return &Handler{productSalesOrderStore: productSalesOrderStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/create_products_sales_order", h.handleCreateProductSalesOrder).Methods(http.MethodPost)
	router.HandleFunc("/get_products_sales_orders", h.handleGetProductSalesOrders).Methods(http.MethodGet)
	router.HandleFunc("/get_products_sales_order/{supplierID}", h.handleGetProductSalesOrder).Methods(http.MethodGet)
}

func (h *Handler) handleCreateProductSalesOrder(w http.ResponseWriter, r *http.Request) {
	var productsalesOrder types.ProductSalesOrder
	if err := utils.ParseJSON(r, &productsalesOrder); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := utils.Validate.Struct(productsalesOrder); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Payload inválido: %v", errors))
		return
	}
	err := h.productSalesOrderStore.CreateProductSalesOrde(productsalesOrder)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]interface{}{
		"data":    productsalesOrder,
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

func (h *Handler) handleGetProductSalesOrders(w http.ResponseWriter, r *http.Request) {
	// bucketID := auth.GetUserIDFromContext(r.Context())
	// fmt.Println("Valor de userIDffsadfsda", bucketID)
	productsSalesOrder, err := h.productSalesOrderStore.GetProductSalesOrders()
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao obter o Pedido de vendas: %v", err), http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, http.StatusOK, productsSalesOrder)
}

func (h *Handler) handleGetProductSalesOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productSalesOrderIDStr, ok := vars["supplierID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("ID do Produto ausente!"))
		return
	}
	parsedProductSalesOrderID, err := strconv.ParseInt(productSalesOrderIDStr, 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("ID do Rascunho inválido!"))
		return
	}

	productsSalesOrder, err := h.productSalesOrderStore.GetProductSalesOrdersBySupplierID(parsedProductSalesOrderID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, productsSalesOrder)
}
