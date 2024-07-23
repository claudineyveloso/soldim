package product

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
	productStore types.ProductStore
}

func NewHandler(productStore types.ProductStore) *Handler {
	return &Handler{productStore: productStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/create_product", h.handleCreateProduct).Methods(http.MethodPost)
	router.HandleFunc("/get_products", h.handleGetProducts).Methods(http.MethodGet)
	router.HandleFunc("/get_products_empty_stock", h.handleGetProductsEmptyStock).Methods(http.MethodGet)
	router.HandleFunc("/get_products_no_movements", h.handleGetProductsNoMovements).Methods(http.MethodGet)
	router.HandleFunc("/get_product/{productID}", h.handleGetProduct).Methods(http.MethodGet)
	router.HandleFunc("/update_product", h.handleUpdateProduct).Methods(http.MethodPut)
	router.HandleFunc("/delete_product/{productID}", h.handleDeleteProduct).Methods(http.MethodDelete)
}

func (h *Handler) handleGetProducts(w http.ResponseWriter, r *http.Request) {
	// bucketID := auth.GetUserIDFromContext(r.Context())
	// fmt.Println("Valor de userIDffsadfsda", bucketID)
	nome := r.URL.Query().Get("nome")
	situacao := r.URL.Query().Get("situacao")
	product, err := h.productStore.GetProducts(nome, situacao)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao obter os produtos : %v", err), http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, http.StatusOK, product)
}

func (h *Handler) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	var product types.ProductPayload
	if err := utils.ParseJSON(r, &product); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := utils.Validate.Struct(product); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Payload inválido: %v", errors))
		return
	}
	err := h.productStore.CreateProduct(product)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]interface{}{
		"data":    product,
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

func (h *Handler) handleGetProductsNoMovements(w http.ResponseWriter, r *http.Request) {
	nome := r.URL.Query().Get("nome")
	situacao := r.URL.Query().Get("situacao")
	product, err := h.productStore.GetProductNoMovements(nome, situacao)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao obter os produtos sem movimentação : %v", err), http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, http.StatusOK, product)
}

func (h *Handler) handleGetProductsEmptyStock(w http.ResponseWriter, r *http.Request) {
	nome := r.URL.Query().Get("nome")
	situacao := r.URL.Query().Get("situacao")
	product, err := h.productStore.GetProductEmptyStock(nome, situacao)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao obter os produtos com estoque vazio : %v", err), http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, http.StatusOK, product)
}

func (h *Handler) handleDeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productIDStr, ok := vars["productID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("ID do Product ausente!"))
		return
	}

	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("ID do Product inválido!"))
		return
	}

	err = h.productStore.DeleteProduct(productID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]interface{}{
		"message": "Registro apagado com sucesso",
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

func (h *Handler) handleGetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productIDStr, ok := vars["productID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("ID do Produto ausente!"))
		return
	}
	parsedDraftsID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("ID do Rascunho inválido!"))
		return
	}

	product, err := h.productStore.GetProductByID(parsedDraftsID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, product)
}

func (h *Handler) handleUpdateProduct(w http.ResponseWriter, r *http.Request) {
	var product types.ProductPayload
	if err := utils.ParseJSON(r, &product); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := utils.Validate.Struct(product); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Payload inválido: %v", errors))
		return
	}
	err := h.productStore.UpdateProduct(product)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]interface{}{
		"data":    product,
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
